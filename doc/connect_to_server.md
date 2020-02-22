# 客户端连接过程

## 连接服务器
>完整连接过程需要三条消息

```graph
sequenceDiagram
    participant c as 客户端
    participant s as 服务端    
        c->>s: 连接请求（发送消息MC_ClientLogon）
        s->>c: 收到连接请求（返回消息MS_ClientLogon）
        c->>s: 确认收到服务器连接，发送消息（MC_ClientConnected
```


## 心跳验证
> 和服务器建立连接后需要每5秒发送一条心跳消息维持连接，一旦服务器20秒内没有收到心跳消息，将主动断开与客户端连接，客户端需要重新建立连接

```graph
sequenceDiagram
    participant c as 客户端
    participant s as 服务端    
        c->>s: 发送心跳（MC_HeartBeat)
        s->>c: 返回当前服务器时间戳（MS_HeartBeat
```

