# 通信协议

通信协议基于tcp，使用protobuf或json来传输数据，每条消息的格式如下：

- `Message Header` 占用8个字节
	1. 头4个字节为二进制消息长度，表示`Message Body` + `Name Crc`的字节数组大小
	2. 后4个字节代表消息名字**Crc32**，例如**client.MC_ClientLogon**消息的**Crc32**为2997000906
	

- `Message Body` **protobuf** **marshal**后的字节数组
>注：所有二进制读取方式都为LittleEndian

---------


## 发送消息
举例说明一下，当客户端发送这条消息时：

```
message C2S_AccountLogon {
  string UserId = 1; // "1"
  int64 AccountId = 2; // 354313566561507648
  string AccountName = 3; // "dudu"
}
```

最后转换成二进制数据格式如下表：

| Header头4字节 | Header 4-8字节 |  Message Body |
| -- | -- | -- |
| 23 | 1971174571 | proto二进制数据 |

实际的二进制流数据如下：
`[23 0 0 0 171 188 125 117 10 1 49 16 192 178 128 144 252 201 177 245 4 26 4 100 117 100 117]` 总共27字节

## 接收消息
接收消息和发送消息过程一致：
- 先读取8个字节的`Message Header`，读取出头4个字节的`Message Body`+`Name Crc`大小
- 再读取出完整的`Message Body`，根据`Message Type`和`Message Name` unmarshal出body为protobuf结构体或者json
