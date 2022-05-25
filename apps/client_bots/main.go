package main

import (
	"log"
	"math"
	"os"

	"github.com/east-eden/server/services/client"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/version"
)

func main() {
	utils.LDFlagsCheck(os.Args, version.Version, version.Help)
	utils.Setrlimit(math.MaxUint32)

	bots := client.NewClientBots()
	if err := bots.Run(os.Args); err != nil {
		log.Fatal("client_bots run error:", err)
		os.Exit(1)
	}

	bots.Stop()
}
