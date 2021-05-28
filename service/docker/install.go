package docker

import (
	"cloud/models"
	"cloud/pkg/ssh"
	"errors"
	"fmt"
	ssh2 "golang.org/x/crypto/ssh"
	"k8s.io/klog"
	"strings"
	"sync"
)

const (
	InstallDockerScript    = "/root/install_docker.sh"
	InstallDockerScriptTpl = "./template/install_docker.sh"
)

func Install(hosts []models.Host) error {
	var errMsg strings.Builder
	var wg sync.WaitGroup
	for _, host := range hosts {
		go func(host models.Host) {
			wg.Add(1)
			defer wg.Done()
			err := InstallDocker(host, nil)
			if err != nil {
				errMsg.WriteString(err.Error() + ";")
			}
		}(host)
	}

	wg.Wait()
	return errors.New(errMsg.String())
}

func InstallDocker(host models.Host, client *ssh2.Client) (err error) {
	if client == nil {
		client, err = ssh.GetSshClient(host)
		if err != nil {
			return fmt.Errorf("nodes: %s, %v", host.IP, err)
		}
	}

	if err := ssh.ScpFile(InstallDockerScriptTpl,
		InstallDockerScript, client); err != nil {
		return fmt.Errorf("nodes: %s, %v", host.IP, err)
	}

	commands := []string{
		fmt.Sprintf(`sh /root/install_docker.sh`),
	}
	for _, cmd := range commands {
		b, err := ssh.SSHExecCmd(client, cmd)
		klog.Infof("%s , resp: \r\n %s", cmd, string(b))
		if err != nil {
			return fmt.Errorf("nodes: %s, %v", host.IP, err)
		}
	}
	klog.Infof("nodes: %s, install docker success", host.IP)
	return nil
}
