#!/usr/bin/env bash

wget http://mirrors.aliyun.com/repo/Centos-7.repo
cp Centos-7.repo /etc/yum.repos.d/
yum install ipvsadm

yum install keepalived
echo 1 >/proc/sys/net/ipv4/ip_forward
