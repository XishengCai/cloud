package kubernetes

import (
	"cloud/models"
	"cloud/pkg/common"
	"cloud/pkg/e"
	"cloud/pkg/ssh"
	"cloud/service/docker"
	"encoding/json"
	"fmt"
	"github.com/gocraft/work"
	"strings"

	"k8s.io/klog"
)

const (
	InstallKubeadmTpl    = "./template/install_kubeadm.sh"
	InstallKubeadmScript = "/root/install_k8s_master.sh"
)

type InstallSlave struct{
	*models.KubernetesSlave
}

func(i InstallSlave) Export(job *work.Job) error{
	klog.Infof("export install k8s slave job: %v", job)
	return nil
}

func (i InstallSlave) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	klog.Infof("Starting job:%s, jobID: %s, install k8s slave  ",job.Name, job.ID)
	return next()
}

func(i InstallSlave) ConsumeJob(job *work.Job) error{
	if job.Args == nil{
		klog.Errorf("jobID:%s, job.Arg is nil", job.ID)
		return nil
	}
	b, err:= json.Marshal(job.Args)
	if err !=nil{
		panic(err)
	}

	k := InstallSlave{}
	err = json.Unmarshal(b, &k)
	if err != nil{
		panic(err)
	}
	return k.joinNodes()
}

func (i InstallSlave) Install() (err error) {
	arg, err := ConvertJobArg(i)
	if err != nil {
		return err
	}
	job, err := installK8sSlaveQueue.Enqueue(installSlave, arg)
	klog.Infof("enqueue job: %v", job)
	return err
}

func handCommandResult(result []byte) string {
	slice := strings.Split(string(result), "\n")
	var command string
	if len(slice) >= 1 {
		command = slice[len(slice)-2]
	}
	return command
}

func (i *InstallSlave) joinNodes() (err error) {
	joinCommand, err := getJoinNodeCommand(i.Master)
	if err != nil {
		return err
	}
	klog.Infof("joinCommand: %s", joinCommand)
	i.JoinSlaveCommand = string(joinCommand)
	var errorList []error
	for _, item := range i.Nodes {
		err := joinNode(item, i.Version, i.JoinSlaveCommand)
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	//return fmt.Errorf(errors.NewAggregate(errorList).Error())
	return e.MergeError(errorList)
}

func getJoinNodeCommand(host models.Host) ([]byte, error) {
	client, err := ssh.GetSshClient(host)
	if err != nil {
		return nil, err
	}
	return ssh.SSHExecCmd(client, "kubeadm token create --print-join-command")

}

func joinNode(host models.Host, version, joinCommand string) (err error) {

	client, err := ssh.GetSshClient(host)
	if err != nil {
		return
	}
	err = docker.InstallDocker(host, client)
	if err != nil {
		return err
	}

	buf, err := common.ParserTemplate(InstallKubeadmTpl,
		struct {
			Version string
		}{
			Version: version,
		})

	if err != nil {
		return
	}

	err = ssh.CopyByteToRemote(client, buf, InstallKubeadmScript)
	if err != nil {
		return
	}
	commands := []string{
		fmt.Sprintf(`sh %s`, InstallKubeadmScript),
		joinCommand,
	}
	for _, cmd := range commands {
		klog.Infof("exec cmd %s", cmd)
		b, err := ssh.SSHExecCmd(client, cmd)
		if err != nil {
			return err
		}
		klog.Infof("resp:  %s", string(b))
		klog.Infof("exec cmd: %s success", cmd)
	}
	klog.Infof("install kubernetes slave node:%s success", host.IP)
	return nil

}