## kubeadm-config
```
apiVersion: kubeadm.k8s.io/v1beta1
controlPlaneEndpoint: xxxxx:6443
imageRepository: k8s.gcr.io
#imageRepository: registry.aliyuncs.com/launcher
kind: ClusterConfiguration
kubernetesVersion: 1.15.3
networking:
  podSubnet: 10.96.0.0/16
  serviceSubnet: 10.244.0.0/24
```
  
# set hostname
```
hostnamectl set-hostname {{.Name}}
cat <<EOF> /etc/hosts
::1	localhost	localhost.localdomain	localhost6	localhost6.localdomain6
127.0.0.1	localhost	localhost.localdomain	localhost4	localhost4.localdomain4
{{.InternalIP}}　{{.Name}}
EOF
```

## k8s yum source of google
```
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kubelet kubeadm kubectl
EOF
```