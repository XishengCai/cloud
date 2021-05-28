#!/usr/bin/env bash
echo "clean env"
sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine

rm -rf /var/lib/docker

echo "install docker 18.09.8"
sudo yum install -y yum-utils

sudo yum-config-manager \
    --add-repo \
    https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

sudo yum install -y docker-ce docker-ce-cli containerd.io

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

systemctl daemon-reload
systemctl enable docker
systemctl restart docker
docker info
