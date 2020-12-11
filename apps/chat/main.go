package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/east-eden/server/entries"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/services/chat"

	// micro plugins
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/store/consul/v2"
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

	c, err := chat.NewChat()
	if err != nil {
		log.Fatal("chat new error:", err)
		os.Exit(1)
	}

	if err = c.Run(os.Args); err != nil {
		log.Fatal("chat run error:", err)
		os.Exit(1)
	}

	c.Stop()
}
