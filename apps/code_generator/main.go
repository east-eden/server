package main

import (
	"fmt"
	"os"

	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/logger"
	"e.coding.net/mmstudio/blade/server/utils"
	log "github.com/rs/zerolog/log"
)

func main() {
	if err := utils.RelocatePath(); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("code_generator")

	// remove all generated files in previous run
	dir, err := os.Getwd()
	if event, pass := utils.ErrCheck(err); !pass {
		event.Msg("get working directory failed")
		os.Exit(1)
	}
	err = os.RemoveAll(fmt.Sprintf("%s/excel/auto/", dir))
	if event, pass := utils.ErrCheck(err, dir); !pass {
		event.Msg("remove all file in config/excel/auto/ failed")
		os.Exit(1)
	}

	// generate go code with excel files
	err = os.MkdirAll(fmt.Sprintf("%s/excel/auto/", dir), 0777)
	if event, pass := utils.ErrCheck(err, dir); !pass {
		event.Msg("make directory config/excel/auto failed")
		os.Exit(1)
	}
	excel.Generate("config/excel")

	log.Info().Msg("generate all go code from excel files success!")
}
