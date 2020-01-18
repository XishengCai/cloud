#!/usr/bin/env bash
echo "clean env"
sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine

rm -rf /var/lib/docker

echo "install docker 18.09.8"
yum install -y yum-utils   device-mapper-persistent-data   lvm2
sudo yum-config-manager     --add-repo     https://download.docker.com/linux/centos/docker-ce.repo
yum install containerd.io -y

echo "download docker rpm package"
# 如何选择最佳网络线路
if [ ! -f "/root/docker-ce-18.09.8-3.el7.x86_64.rpm" ]; then
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-18.09.8-3.el7.x86_64.rpm
fi

if [ ! -f "/root/docker-ce-cli-18.09.8-3.el7.x86_64.rpm" ]; then
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-cli-18.09.8-3.el7.x86_64.rpm
fi

echo "install rpm package"
rpm -ih docker-ce-cli-18.09.8-3.el7.x86_64.rpm
rpm -ih docker-ce-18.09.8-3.el7.x86_64.rpm

echo "config docker daemon"
mkdir /etc/docker
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

systemctl daemon-reload
systemctl enable docker
systemctl restart docker
docker info
