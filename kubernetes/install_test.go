package kubernetes

import (
	"os"
	"testing"
)

func TestInstallKube(t *testing.T) {
	testCase := []KubeInstall{
		{
			MasterNode:   []string{"47.52.20.128"},
			Version:      "1.15.3",
			Registry:     "k8s.gcr.io",
			ClusterName:  "master",
			NetWorkPlug:  "flannel",
			PodCidr:      "10.200.0.0/16",
			ServiceCidr:  "10.96.0.0/12",
			ControlPlane: "47.52.20.128:6443",
		},
	}

	for index, unit := range testCase {
		t.Logf("case %d\r", index)
		err := unit.Install()
		if err != nil {
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
