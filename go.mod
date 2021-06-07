module e.coding.net/mmstudio/blade/server

go 1.15

require (
	e.coding.net/mmstudio/blade/gate v0.0.3
	e.coding.net/mmstudio/blade/golib v0.2.0
	e.coding.net/mmstudio/blade/kvs v0.0.5
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.3.1
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/armon/go-metrics v0.3.4 // indirect
	github.com/aviddiviner/gin-limit v0.0.0-20170918012823-43b5f79762c1
	github.com/coreos/bbolt v1.3.5 // indirect
	github.com/coreos/etcd v3.3.24+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/gammazero/workerpool v0.0.0-20191005194639-971bc780f6d7
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.5 // indirect
	github.com/hashicorp/consul/api v1.6.0 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.9.4 // indirect
	github.com/hellodudu/task v1.0.7
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/juju/ratelimit v1.0.2-0.20191002062651-f60b32039441
	github.com/klauspost/compress v1.9.7 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/nsq/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/store/consul/v2 v2.9.1
	github.com/micro/go-plugins/transport/grpc/v2 v2.9.1
	github.com/micro/go-plugins/transport/tcp/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/breaker/gobreaker/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit/v2 v2.9.1
	github.com/miekg/dns v1.1.31 // indirect
	github.com/mitchellh/mapstructure v1.3.3
	github.com/msgpack/msgpack-go v0.0.0-20130625150338-8224460e6fa3
	github.com/nats-io/jwt v1.0.1 // indirect
	github.com/nats-io/nats.go v1.10.0 // indirect
	github.com/nats-io/nkeys v0.2.0 // indirect
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
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/urfave/cli/v2 v2.0.0
	github.com/valyala/bytebufferpool v1.0.0
	github.com/willf/bitset v1.1.11
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.etcd.io/etcd v3.3.24+incompatible // indirect
	go.mongodb.org/mongo-driver v1.5.1
	go.uber.org/atomic v1.7.0
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/protobuf v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
	stathat.com/c/consistent v1.0.0
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
