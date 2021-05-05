# 客户端连接过程

## 连接服务器
>完整连接过程需要三条消息

```graph
sequenceDiagram
    participant c as 客户端
    participant game as game_server    
    participant gate as gate_server    
        c-->>gate: 请求game服务器信息(发送http post /select_game_addr{userId, userName})
        gate-->>c: 返回game服务器信息(gameId, publicAddr)
        c->>game: 连接请求(发送消息C2S_ClientLogon)
        game->>c: 收到连接请求(返回消息S2C_ClientLogon)
```


## 心跳验证
> 和服务器建立连接后需要每5秒发送一条心跳消息维持连接，一旦服务器20秒内没有收到心跳消息，将主动断开与客户端连接，客户端需要重新建立连接

```graph
sequenceDiagram
    participant c as 客户端
    participant s as game_server    
        c->>s: 发送心跳(C2S_HeartBeat)
        s->>c: 返回当前服务器时间戳(S2C_HeartBeat)
```

