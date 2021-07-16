Store提供分布式存储支持，Registry提供服务注册、发现支持。

reference：

Original Repo [libkv](https://github.com/docker/libkv)

- 保留libkv核心功能，移除部分依赖，调整构造参数
- 为服务发现修正watch部分实现,不只返回变化的事件，返回当前的全量服务列表
- etcd调整keepalive机制,根据ttl确定lease并自动keepalive