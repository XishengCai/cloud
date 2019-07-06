package kubernetes

import (
	"strings"
	"testing"
)

func TestInstallKube(t *testing.T) {
	testCase := []struct {
		Kube   KubeInstall
		Except interface{}
	}{
		{KubeInstall{MasterNodes: []string{"47.99.241.217"},Version:"1.13.5",Registry:"k8s.gcr.io"}, ""},
	}

	for index, unit := range testCase {
		t.Logf("case %d\r", index)
		err := unit.Kube.InstallKube()
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
