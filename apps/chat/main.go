package main

import (
	"fmt"
	"os"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/services/chat"
	"github.com/east-eden/server/utils"
)

func main() {
	// relocate path
	if err := utils.RelocatePath("/server_bin", "/server"); err != nil {
		fmt.Println("relocate path failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("game")

	// load excel entries
	excel.ReadAllEntries("config/csv/")

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
