package ansible

//
//import "testing"
//
//func TestSftpFile(t *testing.T) {
//	connect, err := GetSftpConnectByPassword(user, password, host, port)
//	if err != nil {
//		t.Fatal(err)
//	}
//	CopyFileToRemote(connect, "id_rsa.pub", "/root/id_rsa.pub")
//
//}
//func TestSftpByte(t *testing.T) {
//	var (
//		hosts = "localhost sh.com"
//		cmd   = "yum install nginx -y"
//	)
//
//	playBook := GeneratePlayBook(hosts, cmd)
//	connect, err := GetSftpConnectByPassword(user, password, host, port)
//	if err != nil {
//		t.Fatal(err)
//	}
//	GeneratePlayBook(hosts, cmd)
//	CopyByteToRemote(connect, playBook, "/root/play.yml")
//}
