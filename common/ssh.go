package common

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"strings"
)

// 通过用户名和密码生成一个配置文件
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

// 通过*ssh.Client 执行命令
func SSHExecCmd(client *ssh.Client, cmd string) ([]byte, error) {
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return nil, fmt.Errorf("CombinedOutput fail %v", err)
	}
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

// 复制字节数组到远程服务器上
func CopyByteToRemote(client *ssh.Client, byteStream []byte, remoteFilePath string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Printf("sftp.NewClient occure err %s", err)
		return err
	}
	defer sftpClient.Close()
	dstFile, err := sftpClient.Create(remoteFilePath) //创建文件夹
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer dstFile.Close()
	dstFile.Write(byteStream)
	log.Println("copy byteStream to remote server finished!")
	return nil
}

//通过ssh.ClientConfig创建一个ssh连接
func GetSshClient(host, user, password string, port int) (*ssh.Client, error) {
	fmt.Printf("host:%s, user:%s, password:%s, port:%d", host, user, password, port)
	addr := fmt.Sprintf("%s:%d", host, port)
	sshConfig := GetSshConfigByPassword(user, password)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial:%v", err)
	} else if client == nil {
		panic("get client client failure")
	}
	return client, nil
}
