#!/usr/bin/env bash
# made by Caixisheng  Fri Nov 9 CST 2018

#chec user
[[ $UID -ne 0 ]] && { echo "Must run in root user !";exit; }

isOuter=true

checkNet(){
  # -c: 表示次数，1 为1次
  # -w: 表示deadline, time out的时间，单位为秒，100为100秒。
  ping -c 1 -w 3 www.google.com
  if [[ $? != 0 ]];then
    echo "in inter"
    isOuter=false
  else
    echo "in out"
  fi
}

checkNet

set -e

kubeadm reset --force

cat > kubeadm-config.yaml <<EOF
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

echo "kubeadm init"
kubeadm init --config=/root/kubeadm-config.yaml --upload-certs

# To start using your cluster, you need to run the following as a regular user:
mkdir -p /root/.kube
\cp /etc/kubernetes/admin.conf /root/.kube/config


# install helm3
if $isOuter; then
  curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
  chmod 700 get_helm.sh
  ./get_helm.sh
fi

