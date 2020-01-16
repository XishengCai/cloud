package util

import (
	"encoding/base64"
	"fmt"
	"github.com/rs/xid"
	"io/ioutil"
	"k8s.io/klog"
	"net"
	"os"
	"unsafe"
)

func UuidProvide() string {
	return xid.New().String()
}

// 判断所给路径文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func ElementInArray(element string, array []string) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}

// path can be filename or directory name
func GetFilesByPath(path string) (paths []string, err error) {
	f, err := os.Stat(path)
	if err != nil {
		return
	}
	if !f.IsDir() {
		paths = append(paths, path)
	} else {
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			paths = append(paths, path+"/"+f.Name())
		}
	}
	return
}

func PanicRecover() {
	if err := recover(); err != nil {
		klog.Errorf("recover from panic: %v", err)
	}
}

func BuildTCPSocket(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	fmt.Println(*conn)
	conn.Close()
	return nil
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Base64Decode(code string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", err
	}
	return Bytes2Str(decodeBytes), nil
}
