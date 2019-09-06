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
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-18.09.8-3.el7.x86_64.rpm
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-cli-18.09.8-3.el7.x86_64.rpm

echo "install rpm package"
rpm -ih docker-ce-cli-18.09.8-3.el7.x86_64.rpm
rpm -ih docker-ce-18.09.8-3.el7.x86_64.rpm

echo "config docker daemon"
mkdir -p /etc/docker
cat > /etc/docker/daemon.json <<EOF
{
    "data-root": "/data/docker",
    "storage-driver": "overlay2",
    "exec-opts": [
        "native.cgroupdriver=systemd"
    ],
    "live-restore": true,
    "log-driver":"json-file",
    "log-opts":{
        "max-size":"100m"
    }
}
EOF

cat >/usr/lib/systemd/system/docker.service <<EOF
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
BindsTo=containerd.service
After=network-online.target firewalld.service containerd.service
Wants=network-online.target
Requires=docker.socket

[Service]
Type=notify
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
#Environment="HTTP_PROXY=http://127.0.0.1:8118" "HTTPS_PROXY=http://127.0.0.1:8118"
ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock --insecure-registry=0.0.0.0/0
ExecReload=/bin/kill -s HUP $MAINPID
TimeoutSec=0
RestartSec=2
Restart=always

# Note that StartLimit* options were moved from "Service" to "Unit" in systemd 229.
# Both the old, and new location are accepted by systemd 229 and up, so using the old location
# to make them work for either version of systemd.
StartLimitBurst=3

# Note that StartLimitInterval was renamed to StartLimitIntervalSec in systemd 230.
# Both the old, and new name are accepted by systemd 230 and up, so using the old name to make
# this option work for either version of systemd.
StartLimitInterval=60s

# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity

# Comment TasksMax if your systemd version does not supports it.
# Only systemd 226 and above support this option.
TasksMax=infinity

# set delegate yes so that systemd does not reset the cgroups of docker containers
Delegate=yes

# kill only the docker process, not all processes in the cgroup
KillMode=process

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl restart docker
docker info


