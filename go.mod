module github.com/yokaiio/yokai_server

go 1.13

require (
	github.com/gammazero/workerpool v0.0.0-20191005194639-971bc780f6d7
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/websocket v1.4.2
	github.com/grafana/grafana v6.1.6+incompatible
	github.com/manifoldco/promptui v0.7.0
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/sirupsen/logrus v1.6.0
	github.com/sony/sonyflake v1.0.0
	github.com/urfave/cli/v2 v2.2.0
	github.com/yokaiio/yokai_combat v0.0.0-20200518070716-80f6312bec00 // indirect
	go.mongodb.org/mongo-driver v1.3.3
	golang.org/x/net v0.0.0-20191109021931-daa7c04131f5
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

replace github.com/gogo/protobuf v0.0.0-20190410021324-65acae22fc9 => github.com/gogo/protobuf v0.0.0-20190723190241-65acae22fc9d
