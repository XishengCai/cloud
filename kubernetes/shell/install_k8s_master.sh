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

yum -y  remove kubeadm kubectl kubelet

yum -y install kubelet-{version} kubeadm-{version} kubectl-{version} kubernetes-cni

#sed -i 's/cgroup-driver=systemd/cgroup-driver=cgroupfs/g' /usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf
systemctl daemon-reload
#systemctl enable kubelet
systemctl start kubelet

----
kubeadm reset
ifconfig cni0 down
ip link delete cni0
ifconfig flannel.1 down
ip link delete flannel.1
rm -rf /var/lib/cni/
rm -rf /var/lib/etcd/

kubeadm init    --kubernetes-version={version} \
  --pod-network-cidr=10.244.0.0/16 \
  --image-repository={registry} \
  --node-name={node_name}

# To start using your cluster, you need to run the following as a regular user:
rm  -rf /root/.kube
mkdir /root/.kube
mv -f /etc/kubernetes/admin.conf /root/.kube/config

echo "install flannel"
kubectl create -f  ../yaml/flannel/kube-flannel.yml

# 因为是测试环境机器不够，所以去掉 NoSchedule 污点
kubectl taint nodes k8s-master  node-role.kubernetes.io/master-

echo -e "\033[32m查看K8S状态\033[0m"
sleep 5
kubectl get cs