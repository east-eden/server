package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/chat"
)

func main() {
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
