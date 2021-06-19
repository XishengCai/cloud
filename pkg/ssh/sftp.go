package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"k8s.io/klog"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// GetSftpClientByPassword get sftp client
func GetSftpClientByPassword(user, password, host string, port int) (*sftp.Client, error) {
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

	// connect to ssh
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

// CopyFileToRemote copy file to remote host
func CopyFileToRemote(sftpClient *sftp.Client, localFilePath string, remoteFilePath string) {
	defer func() {
		_ = sftpClient.Close()
	}()
	srcFile, err := os.Open(filepath.Clean(localFilePath))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = srcFile.Close()
	}()

	dstFile, err := sftpClient.Create(remoteFilePath) //创建文件夹
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = dstFile.Close()
	}()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		_, _ = dstFile.Write(buf)
	}

	klog.Info("copy file to remote server finished!")
}

// CopyRemoteToLocal copy remote file to local
func CopyRemoteToLocal(sftpClient *sftp.Client, localFilePath string, remoteFilePath string) {
	var err error
	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = srcFile.Close()
	}()

	dstFile, err := os.Create(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = dstFile.Close()
	}()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		log.Fatal(err)
	}

	klog.Infof("copy file from remote server finished!")
}

// ScpFile copy file local file to dest
func ScpFile(path, dest string, client *ssh.Client) error {
	b, err := ioutil.ReadFile(filepath.Clean(path))
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
