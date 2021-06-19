package docker

import (
	"cloud/models"
	"cloud/pkg/ssh"
	"fmt"

	"k8s.io/klog"
)

const (
	InstallDockerScript    = "/root/install_docker.sh"
	InstallDockerScriptTpl = "./template/install_docker.sh"
)

func InstallDocker(host models.Host) (err error) {
	client, err := ssh.GetClient(host)
	if err != nil {
		return fmt.Errorf("nodes: %s, %v", host.IP, err)
	}

	if err := ssh.ScpFile(InstallDockerScriptTpl,
		InstallDockerScript, client); err != nil {
		return fmt.Errorf("nodes: %s, %v", host.IP, err)
	}

	b, err := ssh.ExecCmd(client, "sh /root/install_docker.sh")
	if err != nil {
		return fmt.Errorf("nodes: %s, %v", host.IP, err)
	}
	klog.Infof("install docker resp: %s", string(b))
	klog.Infof("nodes: %s, install docker success", host.IP)
	return nil
}
