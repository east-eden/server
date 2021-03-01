package main

import (
	"fmt"
	"os"

	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/services/chat"
	"bitbucket.org/funplus/server/utils"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
)

func main() {
	// relocate path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		fmt.Println("relocate path failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("game")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	c, err := chat.NewChat()
	if err != nil {
		fmt.Println("chat new error:", err)
		os.Exit(1)
	}

	if err = c.Run(os.Args); err != nil {
		fmt.Println("chat run error:", err)
		os.Exit(1)
	}

	c.Stop()
}
