### create dashboard by official yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml

### make ca
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout kube-dashboard.key -out kube-dashboard.crt -subj "/CN=zt.dashboard.io/O=zt.dashboard.io"
kubectl create secret tls kube-dasboard-ssl --key kube-dashboard.key --cert kube-dashboard.crt -n kube-system
```

### service expose
#### NodePort
使用补丁的方式，修改service type
```bash
kubectl patch svc kubernetes-dashboard -p '{"spec":{"type":"NodePort"}}' -n kube-system
```

#### ingress
```bash
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ingress-kube-dashboard
  annotations:
    # use the shared ingress-nginx
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
spec:
  tls:
  - hosts:
    - dashboard.kube.com
    secretName: kube-dasboard-ssl
  rules:
  - host: dashboard.kube.com
    http:
      paths:
      - path: /
        backend:
          serviceName: kubernetes-dashboard
          servicePort: 443
```
- kubernetes.io/ingress.class: "nginx"：Inginx Ingress Controller 根据该注解自动发现 Ingress；
- nginx.ingress.kubernetes.io/backend-protocol: Controller 向后端 Service 转发时使用 HTTPS 协议，这个注解必须添加，否则访问会报错，可以看到 Ingress Controller 报错日志：kubectl logs -f nginx-ingress-controller-mg8df
- kubernetes.io/ingress.class: "nginx"：Nginx Ingress Controller 根据该注解自动发现 Ingress；
- host: nginx.kube.com：对外访问的域名；
- serviceName: nginx：对外暴露的 Service 名称；
- servicePort: 80：nginx service 监听的端口；

#### access
- type (一): token
使用集群管理员token
```bash
kubectl -n kube-system describe secrets admin-token-s9d78 | grep token
```
使用自定义权限的token
```bash
kubectl create serviceaccount dashboard-admin -n kube-system

kubectl create clusterrolebinding dashboard-cluster-admin \
--clusterrole=cluster-admin \
--serviceaccount=kube-system:dashboard-admin

kubectl describe secret dashboard-admin-token-8bnk8 -n kube-system

```

- type (二)：kubeconfig<br>
只需要在kubeadm生成的admin.conf文件末尾加上刚刚获取的token就好了。
```bash
- name: kubernetes-admin
  user:
    client-certificate-data: xxxxxxxx
    client-key-data: xxxxxx
    token: "在这里加上token"
```

```bash
cd /etc/kubernetes/pki
export ADMIN="dev-ns-admin"
export NS="dev"
kubectl create serviceaccount $ADMIN -n $NS
kubectl create rolebinding $ADMIN --role=$ADMIN --serviceaccount=$NS:$ADMIN  --namespace=$NS
kubectl create clusterrolebinding $ADMIN --clusterrole=$ADMIN --serviceaccount=$NS:$ADMIN --namespace=$NS

export SECRET=$(kubectl get secret -n $NS |grep token | grep $ADMIN | awk '{print $1}')
kubectl describe secret $SECRET -n $NS

kubectl config set-cluster kubernetes --certificate-authority=./ca.crt \
--server=$(kubectl config view | grep https | awk '{print  $2}') \
--embed-certs=true --kubeconfig=/root/$ADMIN.conf

kubectl config set-context $ADMIN@kubernetes --cluster=kubernets --user=$ADMIN --kubeconfig=/root/$ADMIN.conf

kubectl config set-context $ADMIN@kubernetes --cluster=kubernets --user=$ADMIN --kubeconfig=/root/$ADMIN.conf
kubectl config use-context $ADMIN@kubernetes --kubeconfig=/root/$ADMIN.conf # switch user
kubectl config view

export ADMIN_TOKEN=$(kubectl get secret $SECRET -n $NS -o jsonpath={.data.token} | base64 -d)
echo $ADMIN_TOKEN

kubectl config set-credentials dev-ns-admin --token=$ADMIN_TOKEN  --kubeconfig=/root/dev-ns-admin.conf

kubectl config set-context $ADMIN@kubernetes --cluster=kubernets --user=$ADMIN --kubeconfig=/root/$ADMIN.conf
kubectl config view --kubeconfig=/root/$ADMIN.conf
```
- role 针对单个namespace api 的访问权限
- cluster role 针对所有namespace api 的访问权限


### link
[dashboard](https://github.com/kubernetes/dashboard)<br>
[token](https://blog.csdn.net/weixin_36171533/article/details/82726464)<br>
[暴露方式](https://blog.csdn.net/qianghaohao/article/details/99354304)<br>
