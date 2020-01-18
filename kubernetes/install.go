package kubernetes

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strings"
	"text/template"

	"cloud/common"
	"cloud/model"
	"k8s.io/klog"
)

var kubernetesMasterServer map[string]int

func init() {
	kubernetesMasterServer = make(map[string]int)
	kubernetesMasterServer["api_server"] = 6443
	kubernetesMasterServer["kubelet"] = 6443
}

const (
	InstallK8sMasterScript = "/root/install_k8s_master.sh"
	InstallDockerScript = "/root/install_docker.sh"
	CalicoYaml             = "/root/calico.yaml"
	InstallDockerScriptTpl = "./template/install_docker.sh"
	InstallK8sMasterScriptTpl = "./template/install_k8s_master.sh"
	CalicoYamlTpl             = "./template/calico.yaml"
)

type KubeInstall struct {
	ID           int      `json:"id"`
	ClusterName  string   `json:"cluster_name"`
	MasterNode   []string `json:"master_node,omitempty"`
	SlaveNodes   []string `json:"slave_nodes,omitempty"`
	NetWorkPlug  string   `json:"network_plug"`
	Registry     string   `json:"registry"`
	Version      string   `json:"version"`
	ControlPlane string   `json:"control_plane"`
	PodCidr      string   `json:"pod_cidr"`
	ServiceCidr  string   `json:"service_cidr"`
}

type InstallK8sTemp struct {
	Registry     string `json:"registry"`
	Version      string `json:"version"`
	ControlPlane string `json:"control_plane"`
	PodCidr      string `json:"pod_cidr"`
	ServiceCidr  string `json:"service_cidr"`
	Name         string `json:"name"`
	InternalIP   string `json:"internal_ip"`
}

func (k *KubeInstall) Install() error {
	/* request body
	{
	  "master_node": [
	    "47.91.246.23"
	  ],
	  "version": "1.15.3",
	  "registry": "k8s.gcr.io",
	  "cluster_name": "master",
	  "network_plug": "flannel",
	  "pod_cidr": "10.200.0.0/16 ",
	  "service_cidr": "10.96.0.0/12",
	  "control_plan": "47.91.246.23:6443"
	}
	 */
	// 数据校验
	if err := k.installCheckArguments(); err != nil {
		return err
	}

	// 数据存储
	if err := k.saveCluster("installing"); err != nil {
		return err
	}

	hosts, err := model.GetHosts(k.MasterNode)
	if err != nil {
		return err
	}

	//安装docker
	if len(hosts) != len(k.MasterNode) {
		return errors.New("not found hosts in mysql")
	}

	if err := k.installDocker(hosts); err != nil {
		return err
	}

	// 安装kubernetes
	if err := k.installK8s(hosts); err != nil {
		return err
	}

	return nil
}

func (k *KubeInstall) installSlave() {

}

func (k *KubeInstall) installDocker(hosts []model.Host) error {

	for _, host := range hosts {

		sshClient, err := common.GetSshClient(host.IP, host.User, host.Password, host.Port)
		if err != nil {
			return err
		}

		if err := scpFile(InstallDockerScriptTpl,
			InstallDockerScript, sshClient); err != nil {
			return err
		}

		commands := []string{
			fmt.Sprintf(`sh /root/install_docker.sh`),
		}
		for _, cmd := range commands {
			b, err := common.SSHExecCmd(sshClient, cmd)
			fmt.Printf("%s , resp: \r\n %s", cmd, string(b))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *KubeInstall) installProxy() {

}

func (k *KubeInstall) installK8s(hosts []model.Host) (err error) {

	for _, host := range hosts {
		sshClient, err := common.GetSshClient(host.IP, host.User, host.Password, host.Port)

		t, err := template.ParseFiles(InstallK8sMasterScriptTpl)
		t2, err := template.ParseFiles(CalicoYamlTpl)
		if err != nil {
			return err
		}

		buff := new(bytes.Buffer)
		buff2 := new(bytes.Buffer)

		temStruct := InstallK8sTemp{
			Version:     k.Version,
			Registry:    k.Registry,
			PodCidr:     k.PodCidr,
			ServiceCidr: k.ServiceCidr,
			Name:        host.Name,
			InternalIP:  host.InternalIP,
		}
		err = t.Execute(buff, temStruct)
		if err != nil {
			return err
		}

		err = t2.Execute(buff2, temStruct)
		if err != nil {
			return err
		}

		err = common.CopyByteToRemote(sshClient, buff.Bytes(), InstallK8sMasterScript)
		err = common.CopyByteToRemote(sshClient, buff2.Bytes(), CalicoYaml)

		if err != nil {
			return fmt.Errorf("copy byte err: %v", err)
		}
		commands := []string{
			fmt.Sprintf(`sh %s`, InstallK8sMasterScript),
			fmt.Sprintf(`kubectl create -f %s`, CalicoYaml),
			fmt.Sprintf(`cat %s`, "/root/.kube/config"),
		}
		for _, cmd := range commands {
			b, err := common.SSHExecCmd(sshClient, cmd)
			fmt.Printf("%s , resp: \r\n %s", cmd, string(b))
			if err != nil {
				return err
			}
		}
		if err := k.saveServer(host.ID); err != nil {
			return err
		}
	}
	return
}

func (k *KubeInstall) saveCluster(action string) error {
	cluster := &model.Cluster{
		ClusterName: k.ClusterName,
		Registry:    k.Registry,
		Version:     k.Version,
		NetWorkPlug: k.NetWorkPlug,
		Status:      action,
	}
	klog.Infof("cluster: %+v", cluster)
	id, err := model.AddCluster(cluster)
	if err != nil {
		return err
	}
	k.ID = id
	return nil
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
	if strings.TrimSpace(k.ClusterName) == "" {
		return fmt.Errorf("name can't be null")
	}

	if len(k.MasterNode) == 0 {
		return fmt.Errorf("MasterNodes can't be null")
	}
	return nil
}

func scpFile(path, dest string, client *ssh.Client) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = common.CopyByteToRemote(client, b, dest)
	if err != nil {
		return fmt.Errorf("copy byte err: %v", err)
	}
	return nil
}
