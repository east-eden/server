package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/logger"
	"e.coding.net/mmstudio/blade/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	relocatePath  string // 重定位路径
	readExcelPath string // 读取excel文件路径
	exportGoPath  string // 导出go文件路径
	exportCsvPath string // 导出csv文件路径
)

var (
	BinaryVersion string
	GoVersion     string
	GitLastLog    string
)

func init() {
	flag.StringVar(&relocatePath, "relocatePath", "/server", "重定位到east_eden/server/目录下")
	flag.StringVar(&readExcelPath, "readExcelPath", "../excel/global/", "读取excel路径")
	flag.StringVar(&exportGoPath, "exportGoPath", "excel/auto/", "输出go文件路径")
	flag.StringVar(&exportCsvPath, "exportCsvPath", "config/csv/", "输出csv文件路径")
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
	utils.Setrlimit(math.MaxUint32)

	flag.Parse()
	log.Info().
		Str("relocatePath", relocatePath).
		Str("readExcelPath", readExcelPath).
		Str("exportGoPath", exportGoPath).
		Str("exportCsvPath", exportCsvPath).
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
	mergedExportGoPath := fmt.Sprintf("%s/%s", dir, exportGoPath)
	removeGoDirs, err := os.ReadDir(mergedExportGoPath)
	utils.ErrPrint(err, "")
	for _, dir := range removeGoDirs {
		if strings.Contains(dir.Name(), "entry.go") {
			os.RemoveAll(fmt.Sprintf("%s%s", mergedExportGoPath, dir.Name()))
		}
	}

	// remove all *.csv
	mergedExportCsvPath := fmt.Sprintf("%s/%s", dir, exportCsvPath)
	removeCsvDirs, err := os.ReadDir(mergedExportCsvPath)
	utils.ErrPrint(err, "")
	for _, dir := range removeCsvDirs {
		if strings.Contains(dir.Name(), ".csv") {
			os.RemoveAll(fmt.Sprintf("%s%s", mergedExportCsvPath, dir.Name()))
		}
	}

	// generate go code and csv file from excel files
	excel.Generate(readExcelPath, mergedExportGoPath, mergedExportCsvPath)

	log.Info().Msg("generate all go code from excel files success!")
}
