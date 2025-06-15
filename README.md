# go-gb28181

一个 golang 写的分布式 GB28181 程序 

# 功能特性
-  跨平台服务，支持x86和arm,aarch64等
-  接入设备
-  支持主码流子码流切换
-  支持高标清流切换
-  云台控制，控制设备转向，拉近，拉远
-  支持摄像头语音对讲
-  回放控制，暂停，恢复，倍速，进度条
-  AI, 支持摄像头有无人物通知，火灾告警等，可以配置各种事件主动拉起流，算法因为考虑各个平台兼容性原因未开源


# 架构

## go-sip-gateway

sip-gateway是一个 无状态的 http 服务，主要用于代理ZLMediaKit，提供对外接口访问，可以自由水平扩展，避免单点故障

## go-sip-server

sip-server是一个 有状态的 tcp 服务， 主要是和信令客户端交互，把用户层的指令通过grpc传到信令客户端，主要的目的是为了解决公网网络状况不好导致的丢包问题，sip-server 也可以自由扩展，避免信令服务单点故障导致不可用


## go-sip-client

sip-client是部署在和IPC在同一个局域网的服务，兼容各种平台，主要是信令生成和对IPC设备做交互。



# 代码引用
信令部分使用了 [panjjo/gosip](https://github.com/panjjo/gosip) 回放控制，云台控制都是基于该项目做的扩展


# 开始使用



## go-sip-gateway
gateway配置
```
gateway_api: 0.0.0.0:8999 # 网关服务 restfulapi 端口
database:
  dialect: redis 
  host: 127.0.0.1:6379
  password: 
  DB: 1

```


```
cd cmd/go-sip-gateway
go run go-sip-gateway.go
```

## go-sip-server
server配置
```
sip_id: sip01 # sip服务id 起多个服务该ID必须唯一
Api:
  IP: 0.0.0.0
  Port: 18080
Sip:
  IP: 0.0.0.0
  Port: 12345
gateway: 127.0.0.1:8999
secret: z9hG4bK1233983766 # restful接口验证key 验证请求使用
sign: 3e80d1762a324d5b0ff636e0bd16f1e4
database:
  dialect: redis 
  host: 127.0.0.1:6379
  password: 
  DB: 1 
```
```
cd cmd/go-sip-server
go run cmd/go-sip-server.go
```

## go-sip-client
client配置
```
udp: 0.0.0.0:5060 # sip服务器udp端口
gateway: 127.0.0.1:8999
database:
  dialect: sqllite # mysql postgresql sqllite
gb28181: # gb28181 域，系统id，用户id，通道id，用户数量，初次运行使用配置，之后保存数据库，如果数据库不存在使用配置文件内容
  lid:    "37070000082008000001" # 系统ID
  region: 3707000008           # 系统域
  passwd: "admin123"

```
```
cd cmd/go-sip-client
go run cmd/go-sip-client.go
```