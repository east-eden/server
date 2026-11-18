package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/version"
	log "github.com/rs/zerolog/log"
	"golang.org/x/tools/imports"
)

var (
	relocatePath string // 重定位路径
	readJsonPath string // 读取 json 文件路径
	exportGoPath string // 导出 go 文件路径
)

func init() {
	flag.StringVar(&relocatePath, "relocatePath", "/server", "重定位到 east-eden/server 目录下")
	flag.StringVar(&readJsonPath, "readJsonPath", "config/json", "读取 json 配置表路径")
	flag.StringVar(&exportGoPath, "exportGoPath", "excel/config", "输出 go 文件路径")
}

func main() {
	utils.LDFlagsCheck(os.Args, version.Version, version.Help)
	utils.Setrlimit(math.MaxUint32)

	flag.Parse()
	log.Info().
		Str("relocatePath", relocatePath).
		Str("readJsonPath", readJsonPath).
		Str("exportGoPath", exportGoPath).
		Send()

	if err := utils.RelocatePath(relocatePath); err != nil {
		fmt.Println("relocate failed: ", err, relocatePath)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if !utils.ErrCheck(err, "os.Getwd() failed") {
		os.Exit(1)
	}

	logger.InitLogger("code_generator")

	mergedExportGoPath := fmt.Sprintf("%s/%s", dir, exportGoPath)

	// 删除旧的 *_entry.go，*_manual.go 保留（用于手写扩展逻辑）
	removeGoDirs, err := os.ReadDir(mergedExportGoPath)
	utils.ErrPrint(err, "read export dir failed", mergedExportGoPath)
	for _, d := range removeGoDirs {
		if strings.HasSuffix(d.Name(), "_entry.go") {
			_ = os.RemoveAll(fmt.Sprintf("%s/%s", mergedExportGoPath, d.Name()))
		}
	}

	// 根据 config/json 生成 entry 代码
	excel.GenerateFromJson(excel.ExportLanguage_Go, readJsonPath, mergedExportGoPath)

	formatGenerated(mergedExportGoPath)

	log.Info().Msg("generate all go code from json files success!")
}

// formatGenerated 对生成的文件执行 gofmt + 自动整理 import，
// 避免生成代码里残留未使用的 import 而导致编译失败。
func formatGenerated(dir string) {
	entries, err := os.ReadDir(dir)
	if !utils.ErrCheck(err, "read generated dir failed", dir) {
		return
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}

		path := filepath.Join(dir, e.Name())
		src, err := os.ReadFile(path)
		if !utils.ErrCheck(err, "read generated file failed", path) {
			continue
		}

		out, err := imports.Process(path, src, nil)
		if !utils.ErrCheck(err, "format generated file failed", path) {
			continue
		}

		_ = os.WriteFile(path, out, 0o644)
	}
}
