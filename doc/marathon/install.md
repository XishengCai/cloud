# Install Mesos Cluster

# 1.环境准备
    systemctl stop firewalld.service
    systemctl disable firewalld.service
    setenforce 0

    vi /etc/selinux/config
    SELINUX=disabled

# 2.安装docker
    yum install epel-release -y
    yum install -y docker
    systemctl enable docker
    systemctl start docker

# 3.安装zookeeper
    rpm -Uvh http://repos.mesosphere.io/el/7/noarch/RPMS/mesosphere-el-repo-7-1.noarch.rpm
    rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-mesosphere
    yum install  mesosphere-zookeeper -y

   # 改成自己的内网ip
   HOST_IP=192.168.1.105
   echo 1 > /var/lib/zookeeper/myid
   echo "server.1=${HOST_IP}:2888:3888" >> /etc/zookeeper/conf/zoo.cfg
   systemctl start zookeeper

# 4.安装mesos-master
    rpm -Uvh http://repos.mesosphere.io/el/7/noarch/RPMS/mesosphere-el-repo-7-1.noarch.rpm
    rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-mesosphere
    yum install mesos  -y

    # 改成自己的zookeeper的ip
    HOST_IP=192.168.1.105
    echo "zk://${HOST_IP}:2181/mesos" > /etc/mesos/zk
    echo "${HOST_IP}" > /etc/mesos-master/hostname

    systemctl start mesos-master

# 5.配置mesos slave
    rpm -Uvh http://repos.mesosphere.io/el/7/noarch/RPMS/mesosphere-el-repo-7-1.noarch.rpm
    rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-mesosphere
    yum install mesos -y

    # 改成自己的zookeeper的ip
    HOST_IP=192.168.1.105
    echo "zk://${HOST_IP}:2181/mesos" > /etc/mesos/zk
    echo 'docker,mesos' > /etc/mesos-slave/containerizers
    echo '5mins' > /etc/mesos-slave/executor_registration_timeout
    # 改成本机IP
    HOST_IP=192.168.1.105
    echo "${HOST_IP}" > /etc/mesos-slave/hostname
    systemctl start mesos-slave

# 6.安装marathon
    rpm -Uvh http://repos.mesosphere.io/el/7/noarch/RPMS/mesosphere-el-repo-7-1.noarch.rpm
    rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-mesosphere
    yum install marathon -y

    systemctl start marathon
    docker run -d --restart=always --name marathon-lb --privileged -e PORTS=9090 --net=host ccr.ccs.tencentyun.com/mesos/marathon-lb:v1.11.1 sse -m http://192.168.0.105:8080   --group external

    HAPROXY_GROUP=external
    HAPROXY_0_VHOST=你的域名

    打开 http://你的IP:5050 即可看到mesos的web版控制台
    打开 http://你的IP:8080 即可看到marthon的web版控制台
    由于这样子配置免密码可以访问，所以不能直接如此放到生产环境

# 7.安装mesos-dns

`
mkdir /etc/mesos-dns
cat > /etc/mesos-dns/config.json << EOF
{
    "zk": "zk://172.16.71.242:2181/mesos",
    "masters": ["172.16.71.242:5050"],
    "refreshSeconds": 60,
    "mesosCredentials": {
        "principal": "admin",
        "secret": "password"
    },
    "mesosAuthentication": "basic",
    "ttl": 60,
    "domain": "mesos",
    "port": 53,
    "externalon": false,
    "timeout": 5
}
EOF
# Create Systemd service
cat > /usr/lib/systemd/system/mesos-dns.service << EOF
[Unit]
Description=Mesos-DNS
After=network.target
Wants=network.target

[Service]
ExecStart=/usr/local/bin/mesos-dns -config=/etc/mesos-dns/config.json
Restart=on-failure
RestartSec=20

[Install]
WantedBy=multi-user.target
EOF
# Start mesos-dns
systemctl daemon-reload
systemctl start mesos-dns
`




