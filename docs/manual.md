# 服务端开发手册

## 环境准备

### Go

- 需要 **Go 1.22 及以上**（`go.mod` 已声明 `go 1.22`）
- macOS 可用 Homebrew 安装：`brew install go`

推荐的环境变量（无需 GOPATH 目录结构，仓库使用 Go Modules）：

```shell
export GOPROXY=https://goproxy.cn,direct
export GOBIN=$HOME/go/bin
export PATH=$PATH:$GOBIN
```

### 依赖中间件

| 组件 | 用途 | 是否必需 |
| --- | --- | --- |
| MongoDB | 玩家数据、邮件、排行榜等 | 必需 |
| Consul | 服务注册；启用配置中心时同时提供 KV 配置与节点 id 锁 | release 模式必需 |
| NSQ | 服务间消息总线 | release 模式必需 |
| Redis | 存储层缓存（实现已保留，当前默认关闭） | 可选 |

本地开发也可以用 `docker-compose up -d mongo` 只起 MongoDB。

### 代码生成工具

```shell
brew install protobuf

# protobuf / go-micro 插件（版本需与 go.mod 中的 protobuf 版本匹配）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install go-micro.dev/v4/cmd/protoc-gen-micro@latest
```

## 仓库结构

服务端依赖两个仓库，**协议仓库需要与本仓库同级**：

```
.
├── proto    // .proto 定义与生成脚本
└── server   // 服务端全部代码
```

服务端内部结构：

```
server
├── apps/                // 各服务入口与 Dockerfile
├── cmd/code_generator/  // 配置代码生成器
├── config/
│   ├── json/            // 配置表源数据（新）
│   ├── csv/             // 旧版 CSV 配置（逐步迁移）
│   ├── cert/            // 本地开发证书（make cert 生成，不入库）
│   └── <服务>/config.toml
├── define/              // 常量与数据结构定义
├── excel/
│   ├── config/          // 由 config/json 生成的强类型访问代码
│   └── auto/            // 旧版 CSV 配置访问代码
├── internal/            // 通用组件（配置中心、启动器、属性、计时器等）
├── proto/               // protobuf 生成代码
├── services/            // 各服务实现
├── store/               // 存储层（MongoDB / 缓存）与 protobuf-BSON 编解码
├── transport/           // TCP / WebSocket / KCP 传输与消息注册
└── utils/               // 常用工具
```

## 生成代码

### 配置

配置表以 `config/json` 为源，生成强类型访问代码：

```shell
go run ./cmd/code_generator -readJsonPath=config/json -exportGoPath=excel/config
```

生成的代码位于 `excel/config`，服务启动时由 `excel.LoadAll("config/csv/", "config/json")` 加载。

### 协议

```shell
make proto            # 默认读取同级的 ../proto 仓库
```

会在 `proto/` 下生成 `*.pb.go` 与 `*.pb.micro.go`。新增消息时修改 proto 仓库后重新生成，不要手写序列化。

## 启动服务

各服务可单独启动，也可以 `docker-compose up -d` 本地构建镜像后整体启动。

```shell
# 网关
cd apps/gate && go run main.go

# 游戏逻辑
cd apps/game && go run main.go

# 场景战斗（本地战斗流程使用）
cd apps/combat && go run main.go

# 邮件 / 排行 / 评论 / 聊天
cd apps/mail && go run main.go
cd apps/rank && go run main.go
cd apps/comment && go run main.go
cd apps/chat && go run main.go

# 调试客户端
cd apps/client && go run main.go
```

服务配置位于 `config/<服务>/config.toml`；需要从 Consul 配置中心读取时设置
`CONSUL_HTTP_ADDR` 或传 `-config_center`，详见 README 的「服务启动与配置中心」。

## 测试

```shell
make test           # 单元测试（不依赖外部中间件）
make test_verbose   # 详细输出
make test_bench     # 基准测试
make test_integration  # 集成测试（-tags=integration，需要 MongoDB）
make lint           # go vet
```

## grpc 接口测试

```shell
brew install grpcurl

grpcurl -plaintext -proto server/rank/rank.proto \
  -d '{"RankId": 1, "ObjId": 5218938925633}' \
  localhost:49579 rank.RankService/QueryRankByObjId
```

`localhost:49579` 为 rank 服务的 gRPC 监听地址（以实际配置为准）。
