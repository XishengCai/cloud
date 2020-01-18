#!/usr/bin/env bash
# made by Caixisheng  Fri Nov 9 CST 2018

#chec user
[[ $UID -ne 0 ]] && { echo "Must run in root user !";exit; }


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

cat <<EOF >  /etc/sysctl.d/k8s.conf
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
vm.swappiness=0
EOF


sysctl --system
swapoff -a

hostnamectl set-hostname {{.Name}}
cat <<EOF> /etc/hosts
::1	localhost	localhost.localdomain	localhost6	localhost6.localdomain6
127.0.0.1	localhost	localhost.localdomain	localhost4	localhost4.localdomain4
{{.InternalIP}}　{{.Name}}
EOF
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


ifconfig cni0 down
ip link delete cni0
ifconfig flannel.1 down
ip link delete flannel.1
rm -rf /var/lib/cni/
rm -rf /var/lib/etcd/
rm -rf /etc/kubernetes


yum -y  remove kubeadm kubectl kubelet
yum -y install kubelet-{{.Version}} kubeadm-{{.Version}} kubectl-{{.Version}} kubernetes-cni


systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet

kubeadm reset --force

cat <<EOF> kubeadm-config.yaml
apiVersion: kubeadm.k8s.io/v1beta1
controlPlaneEndpoint: {{.ControlPlane}}
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


