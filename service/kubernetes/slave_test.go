package kubernetes

import (
	"strings"
	"testing"
)

func TestString(t *testing.T){
	x := `W0602 16:03:03.056040   14987 validation.go:28] Cannot validate kube-proxy config - no validator is available
W0602 16:03:03.056091   14987 validation.go:28] Cannot validate kubelet config - no validator is available
kubeadm join 121.196.211.252:6443 --token pt0xqx.cmw3artl9tgbhm06     --discovery-token-ca-cert-hash sha256:12c0e4fd4c0dc7e69545988570ee383ef0cb1615e6a7ea5a85435208826b7621 `
	i := strings.Index(x, "kubeadm join")
	t.Log(x[i:])
}