package kubernetes

import (
	"os"
	"strings"
	"testing"
)

func TestInstallKube(t *testing.T) {
	testCase := []struct {
		Kube   KubeInstall
		Except interface{}
	}{
		{KubeInstall{
			MasterNode: "47.99.241.217",
			Version:    "1.13.5",
			Registry:   "k8s.gcr.io",
			Name:       "test"},
			""},
	}

	for index, unit := range testCase {
		t.Logf("case %d\r", index)
		err := unit.Kube.Install()
		if err != nil {
			if unit.Except == "" {
				t.Fatal(err)
			}
			if strings.Contains(err.Error(), unit.Except.(string)) {
				continue
			}
			t.Fatal(err)
		}
	}
}

func TestCmd(t *testing.T) {
	file, err := os.Open("D:\\Go\\gopath\\src\\github.com\\cloud\\common\\conf.go")
	if err != nil {
		t.Fatalf("open file err: %v", err)
	}
	defer file.Close()
	b1 := make([]byte, 100)
	n, err := file.Read(b1)
	t.Log("n: ", n, "   ", string(b1))

}
