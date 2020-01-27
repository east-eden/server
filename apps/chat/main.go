package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/chat"
)

func main() {
	b, err := chat.New()
	if err != nil {
		log.Fatal("chat new error:", err)
		os.Exit(1)
	}

	if err = b.Run(os.Args); err != nil {
		log.Fatal("chat run error:", err)
		os.Exit(1)
	}

	b.Stop()
}
