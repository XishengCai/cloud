package ssh

import (
	"cloud/models"
	"fmt"
	"k8s.io/klog"
	"net"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// GetSshConfigByPassword 通过用户名和密码生成一个配置文件
func GetSshConfigByPassword(user string, password string) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return sshConfig
}

// SSHExecCmd 通过*ssh.Client 执行命令
func SSHExecCmd(client *ssh.Client, cmd string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil{
		return nil, err
	}
	defer session.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	out, err := session.CombinedOutput(cmd)
	return out, err
}


// CopyByteToRemote 复制字节数组到远程服务器上
func CopyByteToRemote(client *ssh.Client, byteStream []byte, remoteFilePath string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		klog.Errorf("sftp.NewClient  err %s", err)
		return err
	}
	defer sftpClient.Close()
	dstFile, err := sftpClient.Create(remoteFilePath) //创建文件夹
	if err != nil {
		klog.Error(err)
		return err
	}
	defer dstFile.Close()
	_,_ = dstFile.Write(byteStream)
	klog.Info("copy byteStream to remote server finished!")
	return nil
}

// GetSshClient 通过ssh.ClientConfig创建一个ssh连接
func GetSshClient(host models.Host) (*ssh.Client, error) {
	fmt.Printf("host:%s, user:%s, password:%s, port:%d",
		host.IP, host.User, host.Password, host.Port)
	addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
	sshConfig := GetSshConfigByPassword(host.User, host.Password)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial:%v", err)
	} else if client == nil {
		return nil, fmt.Errorf("get k8s k8s failure")
	}
	return client, nil
}
