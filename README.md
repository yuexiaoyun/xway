# API Gateway

----

## 开发计划

----

### 第一阶段

* 2017.12.30 ~ 2018.01.10 现有网关主体功能实现(路由, 认证授权, 代理)
* 2018.01.11 ~ 2018.01.18 分离认证授权服务器
* 2018.01.19 ~ 2018.01.26 集成测试/优化

### 第二阶段

* 2018.01.23 ~ 日志模块实现

## 构建应用

----

>- 构建

```bash
docker run --rm -it -v "$PWD":/go/src/xway -e CGO_ENABLED=0 -e GOOS=linux -w /go/src/xway golang:latest go build -a -installsuffix cgo -o app .
```

>- 启动

```bash
GOMAXPROCS=2 app -etcd xxx:2379 -dbHost xxx:3306 -dbUserName xxx -dbPassword xxx -redisHost xxx:6379 -redisPassword xxx
```

## etcdctl

----

### 添加测试数据

```powershell
etcdctl put /vulcand/backends/b1/backend '{\"Type\":\"http\",\"Settings\":{\"KeepAlive\":{\"MaxIdleConnsPerHost\":200,\"Period\":\"4s\"}}}'

etcdctl put /xway/frontends/f1/frontend '{\"routeId\":\"f1\",\"domainHost\":\"eapi.jiaofucloud.cn\",\"routeUrl\":\"/v5/user/\",\"redirectHost\":\"192.168.2.162:3038\",\"forwardUrl\":\"/user/\",\"type\":\"http\",\"config\":{\"auth\":[\"oauth\"],\"operation\":[{\"rate\":\"0\"}]}}'

etcdctl put /xway/frontends/f2/frontend '{\"routeId\":\"f2\",\"domainHost\":\"eapi.jiaofucloud.cn\",\"routeUrl\":\"/v5/epaperwork/\",\"redirectHost\":\"192.168.2.162:8895\",\"forwardUrl\":\"/epaperwork/\",\"type\":\"http\",\"config\":{\"auth\":[\"oauth\"],\"operation\":[{\"rate\":\"0\"}]}}'

etcdctl put /xway/frontends/f3/frontend '{\"routeId\":\"f3\",\"domainHost\":\"eapi.jiaofucloud.cn\",\"routeUrl\":\"/v5/epaperwork/getReceiveBookchapters/\",\"redirectHost\":\"192.168.2.162:8898\",\"forwardUrl\":\"/epaperwork/v2/getReceiveBookchapters/\",\"type\":\"http\",\"config\":{\"auth\":[\"oauth\"],\"operation\":[{\"rate\":\"0\"}]}}'

etcdctl put /xway/frontends/f4/frontend '{\"routeId\":\"f4\",\"domainHost\":\"eapi.jiaofucloud.cn\",\"routeUrl\":\"/nomux/hello/\",\"redirectHost\":\"192.168.2.102:8708\",\"forwardUrl\":\"/\",\"type\":\"http\",\"config\":{\"auth\":[\"oauth\"],\"operation\":[{\"rate\":\"0\"}]}}'

etcdctl put /xway/frontends/f5/frontend '{\"routeId\":\"f5\",\"domainHost\":\"eapi.jiaofucloud.cn\",\"routeUrl\":\"/oauth/\",\"redirectHost\":\"192.168.2.162:8000\",\"forwardUrl\":\"/oauth/\",\"type\":\"http\",\"config\":{\"auth\":[],\"operation\":[{\"rate\":\"0\"}]}}'
```
