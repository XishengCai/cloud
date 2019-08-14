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
mkdir /etc/docker
cat > /etc/docker/daemon.json <<EOF
{
    "data-root": "/data/docker-data",
    "storage-driver": "overlay2",
    "exec-opts": [
        "native.cgroupdriver=systemd"
    ]
    "live-restore": true,       # 保证 docker daemon重启，但容器不重启
    "log-driver":"json-file",
    "log-opts":{
        "max-size":"100m",
         "
    },
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


#{
#    "authorization-plugins": [],   //访问授权插件
#    "data-root": "",   //docker数据持久化存储的根目录
#    "dns": [],   //DNS服务器
#    "dns-opts": [],   //DNS配置选项，如端口等
#    "dns-search": [],   //DNS搜索域名
#    "exec-opts": [],   //执行选项
#    "exec-root": "",   //执行状态的文件的根目录
#    "experimental": false,   //是否开启试验性特性
#    "storage-driver": "",   //存储驱动器
#    "storage-opts": [],   //存储选项
#    "labels": [],   //键值对式标记docker元数据
#    "live-restore": true,   //dockerd挂掉是否保活容器（避免了docker服务异常而造成容器退出）
#    "log-driver": "",   //容器日志的驱动器
#    "log-opts": {},   //容器日志的选项
#    ,   //设置容器网络MTU（最大传输单元）
#    "pidfile": "",   //daemon PID文件的位置
#    "cluster-store": "",   //集群存储系统的URL
#    "cluster-store-opts": {},   //配置集群存储
#    "cluster-advertise": "",   //对外的地址名称
#    ,   //设置每个pull进程的最大并发
#    ,   //设置每个push进程的最大并发
#    "default-shm-size": "64M",   //设置默认共享内存的大小
#    ,   //设置关闭的超时时限(who?)
#    "debug": true,   //开启调试模式
#    "hosts": [],   //监听地址(?)
#    "log-level": "",   //日志级别
#    "tls": true,   //开启传输层安全协议TLS
#    "tlsverify": true,   //开启输层安全协议并验证远程地址
#    "tlscacert": "",   //CA签名文件路径
#    "tlscert": "",   //TLS证书文件路径
#    "tlskey": "",   //TLS密钥文件路径
#    "swarm-default-advertise-addr": "",   //swarm对外地址
#    "api-cors-header": "",   //设置CORS（跨域资源共享-Cross-origin resource sharing）头
#    "selinux-enabled": false,   //开启selinux(用户、进程、应用、文件的强制访问控制)
#    "userns-remap": "",   //给用户命名空间设置 用户/组
#    "group": "",   //docker所在组
#    "cgroup-parent": "",   //设置所有容器的cgroup的父类(?)
#    "default-ulimits": {},   //设置所有容器的ulimit
#    "init": false,   //容器执行初始化，来转发信号或控制(reap)进程
#    "init-path": "/usr/libexec/docker-init",   //docker-init文件的路径
#    "ipv6": false,   //开启IPV6网络
#    "iptables": false,   //开启防火墙规则
#    "ip-forward": false,   //开启net.ipv4.ip_forward
#    "ip-masq": false,   //开启ip掩蔽(IP封包通过路由器或防火墙时重写源IP地址或目的IP地址的技术)
#    "userland-proxy": false,   //用户空间代理
#    "userland-proxy-path": "/usr/libexec/docker-proxy",   //用户空间代理路径
#    "ip": "0.0.0.0",   //默认IP
#    "bridge": "",   //将容器依附(attach)到桥接网络上的桥标识
#    "bip": "",   //指定桥接ip
#    "fixed-cidr": "",   //(ipv4)子网划分，即限制ip地址分配范围，用以控制容器所属网段实现容器间(同一主机或不同主机间)的网络访问
#    "fixed-cidr-v6": "",   //（ipv6）子网划分
#    "default-gateway": "",   //默认网关
#    "default-gateway-v6": "",   //默认ipv6网关
#    "icc": false,   //容器间通信
#    "raw-logs": false,   //原始日志(无颜色、全时间戳)
#    "allow-nondistributable-artifacts": [],   //不对外分发的产品提交的registry仓库
#    "registry-mirrors": [],   //registry仓库镜像
#    "seccomp-profile": "",   //seccomp配置文件
#    "insecure-registries": [],   //非https的registry地址
#    "no-new-privileges": false,   //禁止新优先级(??)
#    "default-runtime": "runc",   //OCI联盟(The Open Container Initiative)默认运行时环境
#    ,   //内存溢出被杀死的优先级(-1000~1000)
#    "node-generic-resources": ["NVIDIA-GPU=UUID1", "NVIDIA-GPU=UUID2"],   //对外公布的资源节点
#    "runtimes": {   //运行时
#        "cc-runtime": {
#            "path": "/usr/bin/cc-runtime"
#        },
#        "custom": {
#            "path": "/usr/local/bin/my-runc-replacement",
#            "runtimeArgs": [
#                "--debug"
#            ]
#        }
#    }
#}
