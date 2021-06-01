package kubernetes

import (
	"bytes"
	"cloud/models"
	"cloud/pkg/e"
	"cloud/pkg/ssh"
	"cloud/service/docker"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/gocraft/work"
	ssh2 "golang.org/x/crypto/ssh"
	"k8s.io/klog"
)

const (
	installKubeletTpl         = "./template/install_kubeadm.sh"
	installK8sMasterScriptTpl = "./template/install_k8s_master.sh"
	calicoYamlTpl             = "./template/calico.yaml"
)

// InstallKuber implement install k8s master and slave
// ssh to nodes, run shell script
type InstallKuber struct {
	*models.Kubernetes
	status []*status
}

// ConsumeJob worker will call
func (i InstallKuber) ConsumeJob(job *work.Job) error {
	if job.Args == nil {
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
	klog.Infof("Starting job:%s, jobID: %s, install k8s master ", job.Name, job.ID)
	return next()
}

// Install export to API interface
func (i InstallKuber) Install() error {
	arg, err := ConvertJobArg(i)
	if err != nil {
		return err
	}
	// Enqueue a job named "install_k8s" with the specified parameters.
	job, err := installK8sQueue.EnqueueUnique(installMaster, arg)
	klog.Infof("enqueue job: %v", job)
	return err
}

// InstallMaster install k8s master
func (i *InstallKuber) install() error {
	client, err := ssh.GetSshClient(i.PrimaryMaster)
	if err != nil {
		return err
	}
	err = docker.InstallDocker(i.PrimaryMaster)
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
		s := NewStatus(item.IP)
		err = s.joinNode(item, i.Version, i.JoinMasterCommand)
		if err != nil {
			errs = append(errs, err)
		}
		i.status = append(i.status, s)
	}
	return e.MergeError(errs)
}

// installMaster kube init by kubeadm_config, or join k8s as master role
func (i *InstallKuber) installMaster(host models.Host) (err error) {
	client, err := ssh.GetSshClient(host)
	if err != nil {
		return err
	}

	if err := scpData(client, i, []string{installKubeletTpl, installK8sMasterScriptTpl, calicoYamlTpl}); err != nil {
		return err
	}

	commands := []string{
		fmt.Sprintf(`sh %s`, targetFile(installKubeletTpl)),
		fmt.Sprintf(`sh %s`, targetFile(installK8sMasterScriptTpl)),
		fmt.Sprintf(`kubectl create -f %s`, targetFile(calicoYamlTpl)),
		fmt.Sprintf(`cat %s`, "/root/.kube/config"),
	}
	if err := executeCmd(client, commands); err != nil {
		return err
	}
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

func parserTemplate(scriptTpl string, data interface{}) ([]byte, error) {
	t1, err := template.ParseFiles(scriptTpl)
	if err != nil {
		klog.Errorf("%s template parser failed, %v", scriptTpl, err)
		return nil, err
	}
	buff1 := new(bytes.Buffer)

	// 结构体数据映射到模版中
	err = t1.Execute(buff1, data)
	if err != nil {
		klog.Errorf("execute template failed, %v", err)
		return nil, err
	}
	return buff1.Bytes(), nil
}

func scpData(client *ssh2.Client, data interface{}, temp []string) error {
	for _, t := range temp {
		scriptBytes, err := parserTemplate(t, data)
		if err != nil {
			return err
		}
		if err := ssh.CopyByteToRemote(client, scriptBytes, targetFile(t)); err != nil {
			return err
		}
	}
	return nil
}

func executeCmd(client *ssh2.Client, commands []string) error {
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
	return nil
}

func targetFile(tmp string) string {
	t := strings.Split(tmp, "/")
	return "/root/" + t[len(t)-1]
}
