package main

import (
	"fmt"
	"os"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

func main() {
	if err := utils.RelocatePath(); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("code_generator")

	// generate go code with excel files
	excel.Generate("config/excel")

	log.Info().Msg("generate all go code from excel files success!")
}
