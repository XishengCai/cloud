package main

import (
	"bytes"
	"k8s.io/klog"
	"text/template"
)

type KubeInstall struct {
	Version      string `json:"version"`
	Registry     string `json:"registry"`
	ControlPlane string `json:"control_plane"`
	PodCidr      string `json:"pod_dir"`
	ServiceCidr  string `json:"service_cidr"`
}

func main() {
	k := KubeInstall{
		Version:      "1.15.3",
		Registry:     "k8s.gcr.io",
		PodCidr:      "10.200.0.0/16",
		ServiceCidr:  "10.96.0.0/12",
		ControlPlane: "47.56.114.4:6443",
	}

	t, err := template.ParseFiles("./install_k8s_master.sh")
	if err != nil {
		klog.Fatal(err)
	}

	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, k)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info(buffer.String())
}
