package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/battle"
)

func main() {
	b, err := battle.New()
	if err != nil {
		log.Fatal("battle new error:", err)
		os.Exit(1)
	}

	if err = b.Run(os.Args); err != nil {
		log.Fatal("battle run error:", err)
		os.Exit(1)
	}

	b.Stop()
}
