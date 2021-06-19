package kubernetes

import (
	"cloud/models"
	"cloud/pkg/app"
	"cloud/pkg/e"
	"cloud/pkg/ssh"
	"cloud/service/docker"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocraft/work"

	"k8s.io/klog"
)

// InstallSlave batch join slave to  k8s
type InstallSlave struct {
	*models.KubernetesSlave
}

// Export job interface implement
func (i InstallSlave) Export(job *work.Job) error {
	klog.Infof("export install k8s slave job: %v", job)
	return nil
}

// Log job interface implement
func (i InstallSlave) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	klog.Infof("Starting job:%s, jobID: %s, install k8s slave  ", job.Name, job.ID)
	return next()
}

// ConsumeJob job interface implement
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

// Install install slave by master config [hostIP, port root, password]
func (i *InstallSlave) Install() app.ServiceResponse {
	err := i.install()

	return app.ServiceResponse{
		Error:  err,
		Data:   i,
		Status: http.StatusCreated,
		Code:   -1,
	}

}

func (i *InstallSlave) install() error {
	if err := i.getVersion(i.Master); err != nil {
		return err
	}

	if err := i.setJoinCommand(); err != nil {
		return err
	}

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

func (i *InstallSlave) getVersion(host models.Host) error {
	client, err := ssh.GetClient(host)
	if err != nil {
		return err
	}
	b, err := ssh.ExecCmd(client, "kubectl version --short |grep Server")
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return fmt.Errorf("not find k8s version from master")
	}
	version := strings.Split(string(b), " ")[2]
	version = strings.Trim(version, "\n")
	version = strings.Trim(version, "v")
	klog.Info("find master version is ", version)
	i.Version = version
	return nil
}

func (i *InstallSlave) setJoinCommand() error {
	joinCommand, err := getJoinNodeCommand(i.Master)
	if err != nil {
		return err
	}

	cmdResp := string(joinCommand)
	i.JoinSlaveCommand = cmdResp[strings.Index(cmdResp, "kubeadm join"):]
	return nil
}

func (i *InstallSlave) joinNodes() (err error) {
	var errorList []error
	for _, item := range i.Nodes {
		err := joinNode(item, i.Version, i.JoinSlaveCommand)
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	//return fmt.Error(errors.NewAggregate(errorList).Error())
	return e.MergeError(errorList)
}

func getJoinNodeCommand(host models.Host) ([]byte, error) {
	client, err := ssh.GetClient(host)
	if err != nil {
		return nil, err
	}
	return ssh.ExecCmd(client, "kubeadm token create --print-join-command")
}

func joinNode(host models.Host, version, joinCmd string) (err error) {
	client, err := ssh.GetClient(host)
	if err != nil {
		return
	}
	err = docker.InstallDocker(host)
	if err != nil {
		return err
	}

	if err := scpData(client, models.Version{Version: version}, []string{installKubeletTpl}); err != nil {
		return err
	}

	commands := []string{
		fmt.Sprintf(`sh %s`, targetFile(installKubeletTpl)),
		joinCmd,
	}
	if err := executeCmd(client, commands); err != nil {
		return err
	}
	klog.Infof("join node:%s success", host.IP)
	return nil
}
