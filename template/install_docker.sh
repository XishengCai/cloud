#!/usr/bin/env bash
echo "clean env"
yum remove -y docker docker-common container-selinux docker-selinux docker-engine
rm -rf /var/lib/docker

echo "install docker 18.09.8"
yum install -y yum-utils

yum-config-manager \
    --add-repo \
    https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

yum clean packages
#查看docker-ce版本并且安装
yum list docker-ce --showduplicates | sort -r  
yum install -y docker-ce-19.03.14  docker-ce-cli containerd.io


echo "config docker daemon"
mkdir -p /etc/docker
cat > /etc/docker/daemon.json <<EOF
{
  "data-root": "/data/docker",
  "storage-driver": "overlay2",
  "exec-opts": [
    "native.cgroupdriver=systemd",
    "overlay2.override_kernel_check=true"
  ],
  "live-restore": true,
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  }
}
EOF

systemctl enable docker.service
systemctl daemon-reload
systemctl enable docker
systemctl restart docker

docker info
