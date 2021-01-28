package main

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/logger"
	"bitbucket.org/east-eden/server/services/client"
	"bitbucket.org/east-eden/server/utils"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
)

func main() {
	// relocate path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		fmt.Println("relocate path failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("client_bots")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	bots := client.NewClientBots()
	if err := bots.Run(os.Args); err != nil {
		log.Fatal("client_bots run error:", err)
		os.Exit(1)
	}

	bots.Stop()
}
