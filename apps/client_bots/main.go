package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/logger"
	"github.com/yokaiio/yokai_server/services/client"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
)

func main() {
	// check path
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// relocate project path
	if strings.Contains(path, "apps/") || strings.Contains(path, "apps\\") {
		if err := os.Chdir("../../"); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		newPath, _ := os.Getwd()
		fmt.Println("change current path to project root path:", newPath)
	}

	// logger init
	logger.InitLogger("game")

	// entries init
	entries.InitEntries()

	bots := client.NewClientBots()
	if err := bots.Run(os.Args); err != nil {
		log.Fatal("client_bots run error:", err)
		os.Exit(1)
	}

	bots.Stop()
}
