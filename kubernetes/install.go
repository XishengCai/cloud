package kubernetes

import (
	"fmt"
	"github.com/cloud/common"
	"github.com/cloud/model"
	"io/ioutil"
	"strings"
)

type KubeInstall struct {
	MasterNodes []string `json:"master_nodes,omitempty"`
	SlaveNodes  []string `json:"slave_nodes,omitempty"`
	NetWorkPlug string   `json:"network_plug,omitempty"`
	Registry    string   `json:"registry,omitempty"`
	Version     string   `json:"version,omitempty"`
}

func (k *KubeInstall) Install() {
	/*

	 */

}

func (k *KubeInstall) installMaster() {

}

func (k *KubeInstall) checkEnv() {

}

func (k *KubeInstall) installDocker() {

}

func (k *KubeInstall) installProxy() {

}

func (k *KubeInstall)InstallKube() (err error) {
	// TODO: 当前脚本只支持单个master
	for _, node := range k.MasterNodes{
		hosts, _, err := model.GetHostList(0, 1, node)
		if err != nil {
			return err
		}
		host := hosts[0]
		sshClient, err := common.GetSshClient(host.IP, host.User, host.Password, host.Port)
		if err != nil{
			return err
		}
		b, err := ioutil.ReadFile("./shell/install_k8s_master.sh")
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
		for _, cmd := range cmds{
			b, err := common.SSHExecCmd(sshClient, cmd)
			fmt.Printf("%s , resp: \r\n %s", cmd, string(b))
			if err != nil{
				return err
			}
		}

	}

	return
}
