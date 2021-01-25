package main

import (
	"fmt"
	"os"

	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/logger"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

func main() {
	if err := utils.RelocatePath("/excel"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("code_generator")

	// remove all generated files in previous run
	dir, err := os.Getwd()
	if !utils.ErrCheck(err, "get working directory failed") {
		os.Exit(1)
	}

	err = os.RemoveAll(fmt.Sprintf("%s/../../server/excel/auto/", dir))
	if !utils.ErrCheck(err, "remove all file in config/excel/auto/ failed", dir) {
		os.Exit(1)
	}

	// generate go code with excel files
	err = os.MkdirAll(fmt.Sprintf("%s/../../server/excel/auto/", dir), 0777)
	if !utils.ErrCheck(err, "make directory config/excel/auto failed", dir) {
		os.Exit(1)
	}

	// generate from excel files
	excel.Generate(dir, "../server/excel/auto/")

	log.Info().Msg("generate all go code from excel files success!")
}
