package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/combat"
	"github.com/yokaiio/yokai_server/entries"
)

func init() {
	// set working directory as yokai_combat
	os.Chdir("../../")
}

func main() {
	// entries init
	entries.InitEntries()

	c, err := combat.New()
	if err != nil {
		log.Fatal("combat new error:", err)
		os.Exit(1)
	}

	if err = c.Run(os.Args); err != nil {
		log.Fatal("combat run error:", err)
		os.Exit(1)
	}

	c.Stop()
}
