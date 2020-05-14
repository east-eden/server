package main

import (
	"log"
	"os"

	"github.com/yokaiio/yokai_server/client"
	"github.com/yokaiio/yokai_server/internal/global"
)

func init() {
	// set working directory as yokai_server
	os.Chdir("../../")
}

func main() {
	// entries init
	global.InitEntries()

	c, err := client.NewClient()
	if err != nil {
		log.Fatal("client new error:", err)
		os.Exit(1)
	}

	if err = c.Run(os.Args); err != nil {
		log.Fatal("client run error:", err)
		os.Exit(1)
	}

	c.Stop()

}
