package kubernetes

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cloud/common"
	"github.com/cloud/model"
)

var kubernetesMasterServer map[string]int

func init() {
	kubernetesMasterServer = make(map[string]int)
	kubernetesMasterServer["api_server"] = 6443
	kubernetesMasterServer["kubelet"] = 6443
}

type KubeInstall struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	MasterNode  string   `json:"master_node,omitempty"`
	SlaveNodes  []string `json:"slave_nodes,omitempty"`
	NetWorkPlug string   `json:"network_plug"`
	Registry    string   `json:"registry"`
	Version     string   `json:"version"`
}

func (k *KubeInstall) Install() error {
	/*

	 */

	if err := k.installCheckArguments(); err != nil {
		return err
	}
	if err := k.saveCluster(); err != nil {
		return err
	}

	if err := k.installKube(); err != nil {
		return err
	}
	return nil
}

func (k *KubeInstall) installMaster() {

}

func (k *KubeInstall) checkEnv() {

}

func (k *KubeInstall) installDocker() {

}

func (k *KubeInstall) installProxy() {

}

func (k *KubeInstall) installKube() (err error) {
	// TODO: 当前脚本只支持单个master

	hosts, _, err := model.GetHostList(0, 1, k.MasterNode)
	if err != nil {
		return err
	}
	host := hosts[0]
	sshClient, err := common.GetSshClient(host.IP, host.User, host.Password, host.Port)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile("./kubernetes/shell/install_k8s_master.sh")
	if err != nil {
		return err
	}
	cmdList := strings.Replace(string(b), "{version}", k.Version, -1)
	cmdList = strings.Replace(cmdList, "{registry}", k.Registry, -1)
	cmdList = strings.Replace(cmdList, "{node_name}", host.Name, -1)
	kubeScriptFilePath := "/root/kube.sh"
	err = common.CopyByteToRemote(sshClient, []byte(cmdList), "/root/kube.sh")
	if err != nil {
		return fmt.Errorf("copy byte err: %v", err)
	}
	cmds := []string{
		fmt.Sprintf("chmod +x %s", kubeScriptFilePath),
		fmt.Sprintf(`echo -e "y"| sh %s`, kubeScriptFilePath),
	}
	for _, cmd := range cmds {
		b, err := common.SSHExecCmd(sshClient, cmd)
		fmt.Printf("%s , resp: \r\n %s", cmd, string(b))
		if err != nil {
			return err
		}
	}
	if err := k.saveServer(host.ID); err != nil {
		return err
	}

	return
}

func (k *KubeInstall) saveCluster() error {
	cluster := &model.Cluster{
		Name:        k.Name,
		Registry:    k.Registry,
		Version:     k.Version,
		NetWorkPlug: k.NetWorkPlug,
	}

	id, err := model.AddCluster(cluster)
	k.ID = id
	return err
}

func (k *KubeInstall) saveServer(hostID int) error {
	for key, v := range kubernetesMasterServer {
		hs := &model.HostServer{
			HostID:     hostID,
			ServerName: key,
			Port:       v,
			ClusterID:  k.ID,
		}
		if err := model.AddHostServer(hs); err != nil {
			return err
		}

	}
	return nil
}

func (k KubeInstall) installCheckArguments() error {
	if strings.TrimSpace(k.Name) == "" {
		return fmt.Errorf("name can't be null")
	}

	if len(k.MasterNode) == 0 {
		return fmt.Errorf("MasterNodes can't be null")
	}
	return nil
}
