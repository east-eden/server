package main

import (
	"fmt"
	"os"

	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/logger"
	"bitbucket.org/east-eden/server/services/gate"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
)

type IfaceItem interface {
	GetType() int32
}

type Item struct {
	ItemField int32
}

func (i *Item) GetType() int32 {
	return i.ItemField
}

type Equip struct {
	*Item
	EquipField int32
}

func (e *Equip) GetType() int32 {
	return e.EquipField
}

func main() {
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("gate")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	e := &Equip{
		Item: &Item{
			ItemField: 1,
		},
		EquipField: 2,
	}

	func(i IfaceItem) {
		fmt.Println(i.(*Item).ItemField)
	}(e)

	// load xml entries
	// excel.ReadAllXmlEntries("config/entry")

	g := gate.New()
	if err := g.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("gate run failed")
	}

	g.Stop()
}
