# Retail Cloud Controller

Controller for aliyun retail cloud.

[![test](https://github.com/zhigang/arctl/actions/workflows/test.yml/badge.svg)](https://github.com/zhigang/arctl/actions)
[![codecov](https://codecov.io/gh/zhigang/arctl/branch/master/graph/badge.svg?token=I35IBQTJCW)](https://codecov.io/gh/zhigang/arctl)

## Howto

### Quick Start

#### Build

```bash
make build
```

#### Config File

Search config in home directory with name ".arctl" (without extension). Or use flag `--config`.

```bash
log:
  level: "info"
pool:
  size: 5
aksk:
  regionID: "cn-zhangjiakou"
  accessKeyID: "xxxxxxxxx"
  accessKeySecret: "xxxxxxxxxx"
```

#### CMD: create

```bash
> arctl create app --owner owner-id-xxxxxx --title test111 --bizTitle "测试"
> create application test111 succeed, application id: 35181

> arctl create config --name Test --appID 35181 --type 0 --fileType deploy -f ./test.yaml
> create application's deploy config Test succeed, config schema id: 44269

> arctl create env --name "UAT" --type 0 --replicas=1 --appID 35181 --schema 44269 --region "cn-zhangjiakou" --cluster "cluster-id-xxxxx"
> create environment UAT succeed, environment id: 39542
```

#### CMD: label

```bash
> arctl label -i node-id-1xxxxxx,node-id-2xxxxxx -c cluster-id-xxxxxx -l test11=true --add
> add label label.jst.com/test11=true on cluster cluster-id-xxxxxx instance node-id-1xxxxxx succeed
> add label label.jst.com/test11=true on cluster cluster-id-xxxxxx instance node-id-2xxxxxx succeed

> arctl label -i node-id-1xxxxxx,node-id-2xxxxxx -c cluster-id-xxxxxx -l test11=true --remove
> remove label label.jst.com/test11=true on cluster cluster-id-xxxxxx instance node-id-1xxxxxx succeed
> remove label label.jst.com/test11=true on cluster cluster-id-xxxxxx instance node-id-2xxxxxx succeed

```

#### CMD: deploy

```bash
> arctl deploy test111 UAT -t 0 --name arctltest1 --image "registry-vpc.cn-zhangjiakou.aliyuncs.com/test/test:1.4.0.290"
> deploy app (35181)test111 env (39542)UAT image to registry-vpc.cn-zhangjiakou.aliyuncs.com/test/test:1.4.0.290 succeed
```

#### CMD: scale

```bash
> arctl scale test111 UAT --replicas 2 -t 0
> scale app (35181)test111 env (39542)UAT size 1 to 2 succeed, deploy order id: 932489
```

#### CMD: get

```bash
> arctl get apps
> 
> APP ID APP TITLE  STATE TYPE BIZ NAME BIZ TITLE   DEPLOY TYPE         LANGUAGE OS
> 35181  test111    tateless   JST      测试         ContainerDeployment Java     Linux
> ...

> arctl get apps --id 35181 --show-envs
> 
> APP ID APP TITLE STATE TYPE ENV ID ENV NAME ENV TYPE REPLICAS CONFIG ID REGION
> 35181  test111   stateless  24399  UAT      测试      2/2      27590     cn-zhangjiakou
> 35181  test111   stateless  24495  生产      正式      48/48    27696     cn-zhangjiakou
> ...

> arctl get configs test111 -t 0
> 
> APP ID APP TITLE      CONFIG ID CONFIG NAME CONFIG TYPE CONTAINER YAML
> 35181  test111        22882     压测YAML     test        StatefulSet
> 35181  test111        22888     预发YAML     test        StatefulSet

> arctl get pods test111 -t 0
> 
> APP ID  APP TITLE  ENV ID  ENV NAME  ENV TYPE  INSTANCE ID                POD IP        HOST IP         RESTART  HEALTH   REQUESTS   LIMITS     CREATE TIME
> 35181   test111    24399   UAT       0         jck-deployment-yacs-xxxxx  172.20.1.147  172.26.123.235  0        RUNNING  3vcpu 4GB  4vcpu 8GB  2021-11-08T22:33:24
> 35181   test111    24399   UAT       0         jck-deployment-yacs-xxxxx  172.20.1.148  172.26.123.235  0        RUNNING  3vcpu 4GB  4vcpu 8GB  2021-11-08T22:33:24
> ...

> arctl get clusters
> 
> INDEX ID              TITLE   BUSINESS CODE STATUS  NODES NET    POD CIDR      SERVICE CIDR  ENV TYPE CPU   MEMERY VPC      REGION
> 0     cluster-id-xxx  UAT环境  JST           running 14    Terway 172.20.0.0/16 172.21.0.0/20 0        66.0% 44.0%  vpc-xxx  cn-zhangjiakou
> ...

> arctl get nodes --show-labels --id cluster-id-xxxxxx
> 
> CLUSTER ID                        CLUSTER TITLE ENV TYPE NODES REGION
> cluster-id-xxxxxx                 UAT环境        0        14    cn-zhangjiakou
> 
> INSTANCE ID            INSTANCE NAME    PRIVATE IP     CPU MEMORY OS    REGION         EXPIRED TIME LABELS
> i-8vb445kpberivgpswdro UAT-APP-001      172.26.124.74  16  32GB   linux cn-zhangjiakou 2021-12-11   label.jst.com/app-service=true
> ...

> arctl get users
> 
> USER ID        USER NAME     USER TYPE
> user-id-xxxxx  test111       TAOBAO
> ...
```

#### CMD: describe

```bash
> arctl describe app test111
> ...
> arctl describe config test111
> ...
```

#### CMD: delete

```bash
> arctl delete app --id 35181
> delete application 35181 succeed

> arctl delete config --id 44269
> delete application's deploy config 44269 succeed

> arctl delete env --appID 35181 --id 39542
> delete application's environment 39542 succeed

> arctl delete label --id cluster-id-xxxxxx -l test11=true
> delete label label.jst.com/test11=true on on cluster cluster-id-xxxxxx succeed
```
