module e.coding.net/mmstudio/blade/server

go 1.16

require (
	e.coding.net/mmstudio/blade/gate v0.0.6
	e.coding.net/mmstudio/blade/golib v0.2.3
	e.coding.net/mmstudio/blade/kvs v0.0.6
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.3.1
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/asim/go-micro/plugins/broker/nsq/v3 v3.0.0-20210706115128-7f1de77e8c9c
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.0.0-20210611085744-b892efa25f04
	github.com/asim/go-micro/plugins/registry/memory/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/transport/tcp/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/v3 v3.5.2
	github.com/aviddiviner/gin-limit v0.0.0-20170918012823-43b5f79762c1
	github.com/emirpasic/gods v1.12.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.6
	github.com/gorilla/websocket v1.4.2
	github.com/hellodudu/channelwriter v0.0.1
	github.com/hellodudu/task v1.1.5
	github.com/json-iterator/go v1.1.10
	github.com/juju/ratelimit v1.0.2-0.20191002062651-f60b32039441
	github.com/klauspost/compress v1.9.7 // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/msgpack/msgpack-go v0.0.0-20130625150338-8224460e6fa3
	github.com/nitishm/go-rejson v2.0.0+incompatible
	github.com/panjf2000/ants/v2 v2.4.5
	github.com/panjf2000/gnet v1.4.5
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.13.0 // indirect
	github.com/rs/zerolog v1.22.0
	github.com/shopspring/decimal v1.2.0
	github.com/sony/gobreaker v0.4.1
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cast v1.3.1
	github.com/thanhpk/randstr v1.0.4
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/urfave/cli/v2 v2.3.0
	github.com/valyala/bytebufferpool v1.0.0
	github.com/willf/bitset v1.1.11
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.mongodb.org/mongo-driver v1.5.1
	go.uber.org/atomic v1.7.0
	google.golang.org/genproto v0.0.0-20210726200206-e7812ac95cc0 // indirect
	google.golang.org/grpc/examples v0.0.0-20210610163306-6351a55c3895 // indirect
	google.golang.org/protobuf v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
	stathat.com/c/consistent v1.0.0
)
