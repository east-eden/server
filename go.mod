module github.com/east-eden/server

go 1.22

require (
	github.com/asim/go-micro/plugins/broker/nsq/v4 v4.7.0
	github.com/asim/go-micro/plugins/client/grpc/v4 v4.7.0
	github.com/asim/go-micro/plugins/logger/zerolog/v4 v4.7.0
	github.com/asim/go-micro/plugins/registry/consul/v4 v4.7.0
	github.com/asim/go-micro/plugins/registry/memory/v4 v4.7.0
	github.com/asim/go-micro/plugins/server/grpc/v4 v4.7.0
	github.com/asim/go-micro/plugins/transport/tcp/v4 v4.7.0
	github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v4 v4.7.0
	github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v4 v4.7.0
	github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v4 v4.7.0
	github.com/aviddiviner/gin-limit v0.0.0-20170918012823-43b5f79762c1
	github.com/bits-and-blooms/bitset v1.11.0
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/dengsgo/math-engine v0.0.0-20230823154425-78f211b48149
	github.com/emirpasic/gods v1.18.1
	github.com/gin-gonic/gin v1.9.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.1
	github.com/hashicorp/consul/api v1.9.0
	github.com/hellodudu/channelwriter v0.0.1
	github.com/hellodudu/task v1.2.2
	github.com/json-iterator/go v1.1.12
	github.com/juju/ratelimit v1.0.2-0.20191002062651-f60b32039441
	github.com/manifoldco/promptui v0.7.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/panjf2000/ants/v2 v2.8.2
	github.com/panjf2000/gnet/v2 v2.3.3
	github.com/pierrec/lz4/v4 v4.1.22
	github.com/prometheus/client_golang v1.13.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/zerolog v1.27.0
	github.com/shopspring/decimal v1.2.0
	github.com/sony/gobreaker v0.4.1
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cast v1.5.0
	github.com/thanhpk/randstr v1.0.4
	github.com/urfave/cli/v2 v2.11.0
	github.com/valyala/bytebufferpool v1.0.0
	github.com/xtaci/kcp-go v5.4.20+incompatible
	github.com/xuri/excelize/v2 v2.6.0
	go-micro.dev/v4 v4.7.0
	go.mongodb.org/mongo-driver v1.10.6
	go.uber.org/atomic v1.9.0
	golang.org/x/exp v0.0.0-20220706164943-b4a6d9510983
	golang.org/x/tools v0.24.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-hclog v0.12.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/serf v0.9.5 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/klauspost/reedsolomon v1.9.16 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/nsqio/go-nsq v1.0.8 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/xuri/efp v0.0.0-20220407160117-ad0f7a785be8 // indirect
	github.com/xuri/nfp v0.0.0-20220409054826-5e722a1d9e22 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto v0.0.0-20211020151524-b7c3a969101a // indirect
	google.golang.org/grpc v1.42.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
