package kubernetes

import (
	"cloud/models"
	"cloud/pkg/common"
	"cloud/pkg/e"
	"cloud/pkg/ssh"
	"cloud/service/docker"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocraft/work"

	"k8s.io/klog"
)

const (
	InstallKubeadmTpl    = "./template/install_kubeadm.sh"
	InstallKubeadmScript = "/root/install_kubeadm.sh"
)

type status struct {
	node   string
	stage  string
	LogRaw []byte
	event  []string
}

func NewStatus(host string) *status {
	return &status{
		node:   host,
		LogRaw: make([]byte, 0),
		event:  make([]string, 0),
	}
}

type InstallSlave struct {
	*models.KubernetesSlave
	status []*status
}

func (i InstallSlave) Export(job *work.Job) error {
	klog.Infof("export install k8s slave job: %v", job)
	for _, s := range i.status {
		job.Checkin(fmt.Sprintf("node: %s, stage: %s", s.node, s.stage))
	}
	return nil
}

func (i InstallSlave) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	klog.Infof("Starting job:%s, jobID: %s, install k8s slave  ", job.Name, job.ID)
	return next()
}

func (i InstallSlave) ConsumeJob(job *work.Job) error {
	if job.Args == nil {
		klog.Errorf("jobID:%s, job.Arg is nil", job.ID)
		return nil
	}
	b, _ := json.Marshal(job.Args)
	k := InstallSlave{}
	_ = json.Unmarshal(b, &k)
	return k.joinNodes()
}

func (i InstallSlave) Install() (err error) {
	arg, err := ConvertJobArg(i)
	if err != nil {
		return err
	}
	job, err := installK8sSlaveQueue.EnqueueUnique(installSlave, arg)
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

	var errorList []error
	for _, item := range i.Nodes {
		s := NewStatus(item.IP)
		err := s.joinNode(item, i.Version, string(joinCommand))
		if err != nil {
			errorList = append(errorList, err)
		}
		i.status = append(i.status, s)
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

func (s *status) joinNode(host models.Host, version, joinCmd string) (err error) {
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
		joinCmd,
	}
	for _, cmd := range commands {
		klog.Infof("exec cmd %s", cmd)
		b, err := ssh.SSHExecCmd(client, cmd)
		if err != nil {
			return err
		}
		klog.Infof("resp:  %s", string(b))
		klog.Infof("exec cmd: %s success", cmd)
		s.stage = cmd
		s.event = append(s.event, cmd)
		s.LogRaw = append(s.LogRaw, b...)
	}
	klog.Infof("install kubernetes slave node:%s success", host.IP)
	return nil
}
