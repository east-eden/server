module e.coding.net/mmstudio/blade/server

go 1.16

require (
	e.coding.net/mmstudio/blade/gate v0.0.6
	e.coding.net/mmstudio/blade/golib v0.2.3
	e.coding.net/mmstudio/blade/kvs v0.0.6
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.3.1
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/asim/go-micro/plugins/registry/memory/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/transport/tcp/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3 v3.0.0-20210609093110-4af9e245fb62
	github.com/asim/go-micro/v3 v3.5.1
	github.com/aviddiviner/gin-limit v0.0.0-20170918012823-43b5f79762c1
	github.com/emirpasic/gods v1.12.0
	github.com/gammazero/workerpool v0.0.0-20191005194639-971bc780f6d7
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.5
	github.com/gorilla/websocket v1.4.2
	github.com/hellodudu/task v1.0.7
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/juju/ratelimit v1.0.2-0.20191002062651-f60b32039441
	github.com/klauspost/compress v1.9.7 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/miekg/dns v1.1.31 // indirect
	github.com/mitchellh/mapstructure v1.3.3
	github.com/msgpack/msgpack-go v0.0.0-20130625150338-8224460e6fa3
	github.com/nitishm/go-rejson v2.0.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.13.0 // indirect
	github.com/rs/zerolog v1.22.0
	github.com/shopspring/decimal v1.2.0
	github.com/sony/gobreaker v0.4.1
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cast v1.3.1
	github.com/thanhpk/randstr v1.0.4
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/urfave/cli/v2 v2.0.0
	github.com/valyala/bytebufferpool v1.0.0
	github.com/willf/bitset v1.1.11
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.mongodb.org/mongo-driver v1.5.1
	go.uber.org/atomic v1.7.0
	google.golang.org/grpc v1.38.0 // indirect
	google.golang.org/protobuf v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	stathat.com/c/consistent v1.0.0
)

// replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

// replace google.golang.org/grpc v1.38.0 => google.golang.org/grpc v1.26.0
