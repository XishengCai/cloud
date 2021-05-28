package ssh

import (
	"cloud/models"
	"fmt"
	"io"
	"net"
	"strings"

	"k8s.io/klog"

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
	defer session.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	out, err := session.CombinedOutput(cmd)
	return out, err
}

func asyncLog(reader io.Reader) error {
	cache := "" //缓存不足一行的日志信息
	buf := make([]byte, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n") //取出整行的日志
			fmt.Printf("%s%s\n", cache, line)
			cache = s[len(s)-1]
		} else {
			break
		}
	}
	return nil
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
	dstFile.Write(byteStream)
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
