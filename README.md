# function

- [x] install high available k8s masters
- [x] batch install slaves
- [] web terminal

## swagger url
- http://xxx:xxx/swagger/index.html

## post example
### install master
```json
 {
    "clusterName": "test",
    "user":"root",
    "controlPlaneEndpoint": "8.210.79.82:6443",
    "primaryMaster": 
      {
        "ip": "8.210.79.82",
        "password": "password",
        "port": 22,
        "user": "root"
      },
  "backendMasters":[
    {
      "ip":"47.242.36.172",
      "password":"password",
      "port":22
    },
    {
      "ip":"47.242.65.108",
      "password":"password",
      "port":22,
      "user":"root"
      }
   ],
  "netWorkPlug": "calico",
  "registry": "k8s.gcr.io",
  "podCidr": "10.244.0.0/16",
  "serviceCidr": "10.96.0.0/16",
  "version": "1.17.11"
}
```

### install slaves
```json

{
   "master":{
      "ip":"8.210.79.82",
      "password":"password",
      "port":22,
      "user":"root"
   },
   "nodes":[
      {
         "ip":"47.242.36.172",
         "password":"password",
         "port":22,
         "user":"root"
      },
      {
         "ip":"47.242.65.108",
         "password":"password",
         "port":22,
         "user":"root"
      }
   ]
}
```

## feature:
- [x] support centos version
    - [x] 8.0
    
- [x] support k8s version
    - [x] 1.17.11