package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func GetSftpConnectByPassword(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp k8s
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func CopyFileToRemote(sftpClient *sftp.Client, localFilePath string, remoteFilePath string) {
	defer sftpClient.Close()
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	dstFile, err := sftpClient.Create(remoteFilePath) //创建文件夹
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	fmt.Println("copy file to remote server finished!")
}

func CopyRemoteToLocal(sftpClient *sftp.Client, localFilePath string, remoteFilePath string) {
	var err error
	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		log.Fatal(err)
	}

	fmt.Println("copy file from remote server finished!")
}

func ScpFile(path, dest string, client *ssh.Client) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = CopyByteToRemote(client, b, dest)
	if err != nil {
		return fmt.Errorf("copy byte err: %v", err)
	}
	return nil
}

//func CopyByteToRemote(sftpClient *sftp.Client, byteStream []byte, remoteFilePath string) {
//	defer sftpClient.Close()
//	dstFile, err := sftpClient.Create(remoteFilePath) //创建文件夹
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer dstFile.Close()
//	dstFile.Write(byteStream)
//	fmt.Println("copy file to remote server finished!")
//}
