#!/usr/bin/env bash
# made by Caixisheng  Fri Nov 9 CST 2018

#chec user
[[ $UID -ne 0 ]] && { echo "Must run in root user !";exit; }

set -e

echo "添加kubernetes国内yum源"
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
       http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

#cat <<EOF > /etc/yum.repos.d/kubernetes.repo
#[kubernetes]
#name=Kubernetes
#baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
#enabled=1
#gpgcheck=1
#repo_gpgcheck=1
#gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
#exclude=kubelet kubeadm kubectl
#EOF

# Set SELinux in permissive mode (effectively disabling it)

cat <<EOF >  /etc/sysctl.d/k8s.conf
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
vm.swappiness=0
EOF

sysctl --system
swapoff -a

cat > /etc/sysconfig/modules/ipvs.modules <<EOF
#!/bin/bash
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- nf_conntrack_ipv4
EOF

chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4
yum install ipset ipvsadm -y


rm -rf /var/lib/cni/
rm -rf /var/lib/etcd/
rm -rf /etc/kubernetes


yum -y  remove kubeadm kubectl kubelet
yum -y install kubelet-{{.Version}} kubeadm-{{.Version}} kubectl-{{.Version}} --setopt=obsoletes=0


systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet

kubeadm reset --force

cat <<EOF> kubeadm-config.yaml
apiVersion: kubeadm.k8s.io/v1beta1
controlPlaneEndpoint: {{.ControlPlaneEndpoint}}
imageRepository: {{.Registry}}
kind: ClusterConfiguration
kubernetesVersion: {{.Version}}
networking:
  podSubnet: {{.PodCidr}}
  serviceSubnet: {{.ServiceCidr}}
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: ipvs
EOF


kubeadm init --config=/root/kubeadm-config.yaml --upload-certs

# To start using your cluster, you need to run the following as a regular user:
mkdir -p /root/.kube
\cp /etc/kubernetes/admin.conf $HOME/.kube/config

# https://helm.sh/docs/intro/install/
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh