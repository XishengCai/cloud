package common

import (
	"strings"
	"testing"
)

func TestGetSshConfigByPassword(t *testing.T) {

}

func TestSSHExecCmd(t *testing.T) {
	testCase := []struct {
		Host     string
		Port     int
		User     string
		Password string
		Except   interface{}
	}{
		{"47.99.241.217", 22, "cai", "abcd", ""},
		{"47.99.241.217", 22, "notExit", "abdcd", "failed to dial"},
	}

	for index, unit := range testCase {
		t.Logf("case %d\r", index)
		sshClient, err := GetSshClient(unit.Host, unit.User, unit.Password, unit.Port)
		if err != nil {
			if unit.Except == ""{
				t.Fatal(err)
			}
			if strings.Contains(err.Error(), unit.Except.(string)) {
				continue
			}
			t.Fatal(err)
		}
		cmd :="pwd"
		resp, err := SSHExecCmd(sshClient, cmd)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("cmd: %s ", cmd)
		t.Log(string(resp))
	}

}
