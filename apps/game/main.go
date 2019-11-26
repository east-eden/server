package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/game"
)

func main() {
	g, err := game.New()
	if err != nil {
		log.Fatal("game new error:", err)
		os.Exit(1)
	}

	if err = g.Run(os.Args); err != nil {
		log.Fatal("game run error:", err)
		os.Exit(1)
	}

	g.Stop()
}
