package kubernetes

import (
	"bytes"
	"cloud/models"
	"cloud/pkg/e"
	"cloud/pkg/ssh"
	"cloud/service/docker"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/gocraft/work"
	ssh2 "golang.org/x/crypto/ssh"
	"k8s.io/klog"
)

const (
	installK8sMasterScript    = "/root/install_k8s_master.sh"
	calicoYaml                = "/root/calico.yaml"
	installK8sMasterScriptTpl = "./template/install_k8s_master.sh"
	calicoYamlTpl             = "./template/calico.yaml"
)

// InstallKuber implement install k8s master and slave
// ssh to nodes, run shell script
type InstallKuber struct {
	*models.Kubernetes
	slave  *models.KubernetesSlave
}
// ConsumeJob worker will call
func (i InstallKuber) ConsumeJob(job *work.Job) error {
	if job.Args == nil{
		klog.Errorf("jobID:%s, job.Arg is nil", job.ID)
		return nil
	}
	b, _ := json.Marshal(job.Args)

	k := InstallKuber{}
	_ = json.Unmarshal(b, &k)
	return k.install()
}

func (i InstallKuber) Export(job *work.Job) error {
	klog.Infof("export install k8s master job: %v", job)
	return nil
}

func (i InstallKuber) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	klog.Infof("Starting job:%s, jobID: %s, install k8s master ",job.Name, job.ID)
	return next()
}



// Install export to API interface
func (i InstallKuber) Install() error {
	arg, err := ConvertJobArg(i)
	if err != nil {
		return err
	}
	// Enqueue a job named "install_k8s" with the specified parameters.
	job, err := installK8sQueue.Enqueue(installMaster, arg)
	klog.Infof("enqueue job: %v", job)
	return err
}

// InstallMaster install k8s master
func (i *InstallKuber) install() error {
	client, err := ssh.GetSshClient(i.PrimaryMaster)
	err = docker.InstallDocker(i.PrimaryMaster, client)
	if err != nil {
		klog.Errorf("install docker failed: %v", err)
		return fmt.Errorf("install docker failed: %v", err)
	}

	err = i.installMaster(i.PrimaryMaster)
	if err != nil {
		klog.Errorf("install master failed: %v", err)
		return fmt.Errorf("install master failed: %v", err)
	}

	// get joinMaster cmd
	i.JoinMasterCommand, err = getJoinMasterCommand(client)
	if err != nil {
		klog.Errorf("getJoinMasterCommand failed: %v", err)
		return fmt.Errorf("getJoinMasterCommand failed: %v", err)
	}

	klog.Infof("joinMasterCommand: %s", i.JoinMasterCommand)
	var errs []error
	for _, item := range i.BackendMasters {
		err = joinNode(item, i.Version, i.JoinMasterCommand)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return e.MergeError(errs)
}

func (i *InstallKuber) installMaster(host models.Host) (err error) {
	client, err := ssh.GetSshClient(host)

	t1, err := template.ParseFiles(installK8sMasterScriptTpl)
	if err != nil {
		klog.Errorf("%s template parser failed, %v", installK8sMasterScriptTpl, err)
		return err
	}

	t2, err := template.ParseFiles(calicoYamlTpl)
	if err != nil {
		klog.Errorf("%s template parser failed, %v", calicoYaml, err)
		return
	}

	buff1 := new(bytes.Buffer)
	buff2 := new(bytes.Buffer)

	err = t1.Execute(buff1, i)
	if err != nil {
		klog.Errorf("execute template failed, %v", err)
		return
	}

	err = t2.Execute(buff2, i)
	if err != nil {
		klog.Errorf("execute template failed, %v", err)
		return
	}

	err = ssh.CopyByteToRemote(client, buff1.Bytes(), installK8sMasterScript)
	if err != nil {
		klog.Errorf("copy byte err: %v", err)
		return
	}

	err = ssh.CopyByteToRemote(client, buff2.Bytes(), calicoYaml)
	if err != nil {
		klog.Errorf("copy byte err: %v", err)
		return
	}
	commands := []string{
		fmt.Sprintf(`sh %s`, installK8sMasterScript),
		fmt.Sprintf(`kubectl create -f %s`, calicoYaml),
		fmt.Sprintf(`cat %s`, "/root/.kube/config"),
	}
	for _, cmd := range commands {
		klog.Infof("exec cmd %s", cmd)
		b, err := ssh.SSHExecCmd(client, cmd)
		if err != nil {
			klog.Errorf("SSHExecCmd failed, %v", err)
			return err
		}
		klog.Infof("resp:  %s", string(b))
		klog.Infof("exec cmd: %s success", cmd)

	}
	klog.Infof("install kubernetes master node:%s success", host.IP)
	return
}

func getJoinMasterCommand(client *ssh2.Client) (string, error) {
	jointNodeCmd, err := ssh.SSHExecCmd(client, "kubeadm token create --print-join-command")
	if err != nil {
		return "", err
	}
	result, err := ssh.SSHExecCmd(client, "kubeadm init phase upload-certs --upload-certs | awk 'END {print}'")
	if err != nil {
		return "", err
	}
	certificateKey := handCommandResult(result)

	return handCommandResult(jointNodeCmd) + " --control-plane --certificate-key  " + certificateKey, nil

}
