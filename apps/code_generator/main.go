package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/logger"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	relocatePath  string // 重定位路径
	readExcelPath string // 读取excel文件路径
	exportPath    string // 导出路径
)

func init() {
	flag.StringVar(&relocatePath, "relocatePath", "/server", "重定位到east_eden/server/目录下")
	flag.StringVar(&readExcelPath, "readExcelPath", "config/excel/", "读取excel路径")
	flag.StringVar(&exportPath, "exportPath", "excel/auto/", "输出go文件路径")
}

func main() {
	flag.Parse()
	log.Info().
		Str("relocatePath", relocatePath).
		Str("readExcelPath", readExcelPath).
		Str("exportPath", exportPath).
		Send()

	if err := utils.RelocatePath(relocatePath); err != nil {
		fmt.Println("relocate failed: ", err, relocatePath)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if pass := utils.ErrCheck(err, "os.Getwd() failed"); !pass {
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("code_generator")

	mergedExportPath := fmt.Sprintf("%s/%s", dir, exportPath)
	err = os.RemoveAll(mergedExportPath)
	if !utils.ErrCheck(err, "remove all files in config/excel/auto/ failed", relocatePath, mergedExportPath) {
		os.Exit(1)
	}

	// generate go code with excel files
	err = os.MkdirAll(mergedExportPath, 0777)
	if !utils.ErrCheck(err, "make directory failed", dir, mergedExportPath) {
		os.Exit(1)
	}

	// generate from excel files
	excel.Generate(readExcelPath, mergedExportPath)

	log.Info().Msg("generate all go code from excel files success!")
}
