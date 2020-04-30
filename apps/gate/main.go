package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/gate"
)

func main() {
	g, err := gate.New()
	if err != nil {
		log.Fatal("gate new error:", err)
		os.Exit(1)
	}

	if err = g.Run(os.Args); err != nil {
		log.Fatal("gate run error:", err)
		os.Exit(1)
	}

	g.Stop()
}
