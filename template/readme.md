# kubeadm-config

apiVersion: kubeadm.k8s.io/v1beta1
controlPlaneEndpoint: xxxxx:6443
imageRepository: k8s.gcr.io
#imageRepository: registry.aliyuncs.com/launcher
kind: ClusterConfiguration
kubernetesVersion: 1.15.3
networking:
  podSubnet: 10.96.0.0/16
  serviceSubnet: 10.244.0.0/24
  
# set hostname
hostnamectl set-hostname {{.Name}}
cat <<EOF> /etc/hosts
::1	localhost	localhost.localdomain	localhost6	localhost6.localdomain6
127.0.0.1	localhost	localhost.localdomain	localhost4	localhost4.localdomain4
{{.InternalIP}}　{{.Name}}
EOF




Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

You can now join any number of the control-plane node running the following command on each as root:

  kubeadm join 8.210.79.82:6443 --token f06d6c.rks00oknesdpqp94 \
    --discovery-token-ca-cert-hash sha256:3b0216c4ddd6980f69e8d5b1d08bcf89924e68b75f1bf3b45ae296c2aee80afd \
    --control-plane --certificate-key 975da437e50aa3d02a5c0016ad6ffe66748e06307b9578e8b977b78093fccb10

Please note that the certificate-key gives access to cluster sensitive data, keep it secret!
As a safeguard, uploaded-certs will be deleted in two hours; If necessary, you can use
"kubeadm init phase upload-certs --upload-certs" to reload certs afterward.

