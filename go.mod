module github.com/east-eden/server

go 1.18

require (
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.3.1
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
	github.com/bits-and-blooms/bitset v1.2.2
	github.com/emirpasic/gods v1.12.0
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/mock v1.4.4
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.6
	github.com/gorilla/websocket v1.4.2
	github.com/hellodudu/channelwriter v0.0.1
	github.com/hellodudu/task v1.2.2
	github.com/json-iterator/go v1.1.10
	github.com/juju/ratelimit v1.0.2-0.20191002062651-f60b32039441
	github.com/manifoldco/promptui v0.7.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/msgpack/msgpack-go v0.0.0-20130625150338-8224460e6fa3
	github.com/nitishm/go-rejson v2.0.0+incompatible
	github.com/panjf2000/ants/v2 v2.4.5
	github.com/panjf2000/gnet v1.4.5
	github.com/prometheus/client_golang v1.7.1
	github.com/rs/zerolog v1.22.0
	github.com/shopspring/decimal v1.2.0
	github.com/sony/gobreaker v0.4.1
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cast v1.3.1
	github.com/thanhpk/randstr v1.0.4
	github.com/urfave/cli/v2 v2.3.0
	github.com/valyala/bytebufferpool v1.0.0
	github.com/willf/bitset v1.1.11
	github.com/xtaci/kcp-go v5.4.20+incompatible
	go.mongodb.org/mongo-driver v1.5.1
	go.uber.org/atomic v1.7.0
	golang.org/x/exp v0.0.0-20220328175248-053ad81199eb
	google.golang.org/protobuf v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	stathat.com/c/consistent v1.0.0
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da // indirect
	github.com/aws/aws-sdk-go v1.37.27 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/hashicorp/consul/api v1.3.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.0 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/serf v0.8.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/compress v1.9.7 // indirect
	github.com/klauspost/cpuid/v2 v2.0.6 // indirect
	github.com/klauspost/reedsolomon v1.9.16 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/micro/cli/v2 v2.1.2 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/nsqio/go-nsq v1.0.8 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.13.0 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	github.com/richardlehane/mscfb v1.0.3 // indirect
	github.com/richardlehane/msoleps v1.0.1 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/xuri/efp v0.0.0-20200605144744-ba689101faaf // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210726200206-e7812ac95cc0 // indirect
	google.golang.org/grpc v1.39.0 // indirect
	google.golang.org/grpc/examples v0.0.0-20210610163306-6351a55c3895 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)
