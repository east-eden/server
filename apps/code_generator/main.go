package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	relocatePath  string // 重定位路径
	readExcelPath string // 读取excel文件路径
	exportPath    string // 导出路径
)

var (
	BinaryVersion string
	GoVersion     string
	GitLastLog    string
)

func init() {
	flag.StringVar(&relocatePath, "relocatePath", "/server", "重定位到east_eden/server/目录下")
	flag.StringVar(&readExcelPath, "readExcelPath", "config/excel/", "读取excel路径")
	flag.StringVar(&exportPath, "exportPath", "excel/auto/", "输出go文件路径")
}

func main() {
	utils.LDFlagsCheck(
		os.Args,

		// version
		func() {
			fmt.Println("BinaryVersion:", BinaryVersion)
			fmt.Println("GoVersion:", GoVersion)
			fmt.Println("GitLastLog:", GitLastLog)
			os.Exit(0)
		},

		// help
		func() {
			fmt.Println("The commands are:")
			fmt.Println("version       see all versions")
			os.Exit(0)
		},
	)

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

	// remove all *_entry.go
	mergedExportPath := fmt.Sprintf("%s/%s", dir, exportPath)
	removeDirs, err := ioutil.ReadDir(mergedExportPath)
	utils.ErrPrint(err, "")
	for _, dir := range removeDirs {
		if strings.Contains(dir.Name(), "entry.go") {
			os.RemoveAll(fmt.Sprintf("%s%s", mergedExportPath, dir.Name()))
		}
	}

	// generate from excel files
	excel.Generate(readExcelPath, mergedExportPath)

	log.Info().Msg("generate all go code from excel files success!")
}
