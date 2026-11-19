# 通信协议

服务端与客户端（以及服务之间）使用 **protobuf** 传输消息。TCP、WebSocket、KCP 三种传输共用同一套报文格式：

- `Message Header` 固定 8 字节（小端序）
  1. 前 4 字节：`Message Body` 字节数 + 4，即 `4 字节消息名 CRC + protobuf 字节数`
  2. 后 4 字节：消息名（protobuf 消息的短名，如 `C2S_AccountLogon`）的 **CRC32**
- `Message Body`：protobuf marshal 后的字节数组

> CRC 取的是消息短名（`Descriptor().Name()`），不是包名限定的全名。

## 发送消息

以登录消息为例：

```protobuf
message C2S_AccountLogon {
  string UserId = 1;      // "1"
  int64 AccountId = 2;    // 354313566561507648
  string AccountName = 3; // "dudu"
}
```

按上面的字段值序列化后：

| 字段 | 值 |
| --- | --- |
| 消息名 | `C2S_AccountLogon` |
| 消息名 CRC32 | `1971174571` |
| protobuf 字节数 | 19 |
| Header 前 4 字节（长度 = 19 + 4） | `23` |
| Header 后 4 字节（CRC32） | `1971174571` |

完整二进制流（共 27 字节）：

```
[23 0 0 0 171 188 125 117 10 1 49 16 192 178 128 144 252 201 177 245 4 26 4 100 117 100 117]
```

前 8 字节是 header，后面是 protobuf body。

## 接收消息

与发送过程一致：

1. 先读 8 字节 header，取出长度与消息名 CRC
2. 按长度读取 body（长度字段包含末尾 4 字节 CRC，所以 body 实际是 `长度 - 4` 字节）
3. 用 CRC 在消息注册表里找到消息类型并反序列化

## 消息注册

- 消息通过 `transport.RegisterProtobufMessage(msg, handler)` 注册，注册表以消息名 CRC32 为键
- 传输层保留 JSON codec 的注册接口，但当前所有消息都走 protobuf
- 新增消息时在 proto 仓库定义并重新生成，不要手写序列化
