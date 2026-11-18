package excel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"encoding/json"

	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/thanhpk/randstr"
)

var (
	GlobalResultReport = &ResultReport{Errs: []*ReportData{}, Warns: []*ReportData{}}
)

type ReportData struct {
	FileName  string
	SheetName string
	KeyName   string
	Reason    string
	MainId    string
	CellValue string
}

func (r *ReportData) String() string {
	return fmt.Sprintf("文件名:%s, 表名:%s, 字段名:%s, 原因:%s, 主键ID:%s, 有问题的值:%s", r.FileName, r.SheetName, r.KeyName, r.Reason, r.MainId, r.CellValue)
}

type ResultReport struct {
	Errs  []*ReportData
	Warns []*ReportData
	lock  sync.Mutex
}

func (r *ResultReport) AddError(fileName, SheetName, keyName, reason string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Errs = append(r.Errs, &ReportData{FileName: fileName, KeyName: keyName, Reason: reason, SheetName: SheetName})
}

func (r *ResultReport) Errors() []*ReportData {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.Errs
}

func (r *ResultReport) AddWarnings(fileName, SheetName, keyName, reason, mainId, cellValue string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Warns = append(r.Warns, &ReportData{FileName: fileName, SheetName: SheetName, KeyName: keyName, Reason: reason, MainId: mainId, CellValue: cellValue})
}

func (r *ResultReport) Warnings() []*ReportData {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.Warns
}

func (r *ResultReport) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Errs = make([]*ReportData, 0)
	r.Warns = make([]*ReportData, 0)
}

type JsonSheet struct {
	SheetFile string  `json:"sheetFile"`
	ConfFile  string  `json:"confFile"`
	Bean      string  `json:"bean"`
	Mode      string  `json:"mode"`
	Options   string  `json:"options,omitempty"`
	Values    [][]any `json:"values"`
}

// load all json files
func loadJson(lanType ExportLanguage_Type, dirPath string, fileNames []string) []string {
	GlobalResultReport.Clear()
	excelKeyAliasMapping = NewKeyAliasMapping() // key alias mapping
	excelKeyRefMapping = NewKeyRefMapping()     // key ref mapping

	replacedFilenames := make([]string, 0, len(fileNames))

	// first deal with enum and poly sheets
	enumAndPoly := make(map[string]*JsonSheet)

	// then deal with one and map sheets
	oneAndMap := make(map[string]*JsonSheet)

	wg := utils.WaitGroupWrapper{}
	mu := sync.Mutex{}
	for _, v := range fileNames {
		filename := v
		wg.Wrap(func() {
			defer utils.CaptureException(filename)

			filePath := filepath.Join(dirPath, filename)

			jsonFile, err := os.ReadFile(filePath)
			if !utils.ErrCheck(err, "OpenFile failed", filePath) {
				return
			}

			// read rows from excel files
			var sheet JsonSheet
			if err := json.Unmarshal(jsonFile, &sheet); err != nil {
				log.Fatal().Str("file", filename).Err(err).Msg("json unmarshal failed")
				return
			}

			sheet.SheetFile, _ = strings.CutSuffix(sheet.SheetFile, "正式")
			sheet.SheetFile, _ = strings.CutSuffix(sheet.SheetFile, "马甲")

			mu.Lock()
			if sheet.Bean == "SYS_ENUM" || sheet.Bean == "SYS_POLY" {
				// 文件名有中文时，使用Bean作为文件名
				for _, runeValue := range filename {
					if unicode.Is(unicode.Han, runeValue) {
						filename = fmt.Sprintf("%s.json", sheet.Bean)
						break
					}
				}
				enumAndPoly[filename] = &sheet
			} else {
				// 文件名有中文时，使用Bean作为文件名
				for _, runeValue := range filename {
					if unicode.Is(unicode.Han, runeValue) {
						filename = fmt.Sprintf("%s.json", sheet.Bean)
						break
					}
				}
				oneAndMap[filename] = &sheet
			}
			replacedFilenames = append(replacedFilenames, filename)
			mu.Unlock()
		})
	}
	wg.Wait()

	// deal with enum and poly sheets
	excelEnumTypeRaws = sync.Map{}
	excelPolyTypeRaws = sync.Map{}
	for name, s := range enumAndPoly {
		filename := name
		sheet := s
		wg.Wrap(func() {
			rowDatas, err := parseJsonRows(lanType, dirPath, filename, sheet)
			if !utils.ErrCheck(err, "parseJsonRows failed", filename) {
				return
			}

			// raw data
			excelFileRaws.Store(filename, rowDatas)

			// enum
			for _, enumData := range rowDatas.EnumData {
				enumRaw := NewExcelEnumRaw(enumData.EnumType, enumData.EnumOriginType, sheet.SheetFile, sheet.ConfFile, sheet.Options+","+enumData.Options)
				for _, enumItem := range enumData.Enums {
					enumRaw.AddEnum(enumItem.EnumName, enumItem.EnumDesc, cast.ToInt(enumItem.EnumValue))

					// 添加枚举到全局别名中
					excelKeyAliasMapping.AddAlias(enumItem.EnumDesc, rowDatas, cast.ToInt32(enumItem.EnumValue))
				}
				if _, ok := excelEnumTypeRaws.Load(enumData.EnumType); !ok {
					excelEnumTypeRaws.Store(enumData.EnumType, enumRaw)
				} else {
					log.Error().Str("enum_type", enumData.EnumType).
						Str("file_name", rowDatas.Filename).Msg("enum type already exists")
				}
			}

			// poly
			for _, polyData := range rowDatas.PolyData {
				if len(polyData.PolyParentName) > 0 {
					fullParentName := polyData.PolyParentName
					polyData.PolyParentName = strings.Title(fullParentName)
				}
				if _, ok := excelPolyTypeRaws.Load(polyData.PolyType); !ok {
					excelPolyTypeRaws.Store(polyData.PolyType, polyData)
				} else {
					log.Error().Str("poly_type", polyData.PolyType).
						Str("file_name", rowDatas.Filename).Msg("poly type already exists")
				}
			}
		})
	}
	wg.Wait()

	// pre parse one and map sheets
	for name, s := range oneAndMap {
		filename := name
		sheet := s
		wg.Wrap(func() {
			preParseJsonRows(lanType, dirPath, filename, sheet)
		})
	}
	wg.Wait()

	// deal with one and map sheets
	for name, s := range oneAndMap {
		filename := name
		sheet := s
		wg.Wrap(func() {
			rowDatas, err := parseJsonRows(lanType, dirPath, filename, sheet)
			if !utils.ErrCheck(err, "parseJsonRows failed", filename) {
				return
			}

			// save to global
			excelFileRaws.Store(filename, rowDatas)
		})
	}
	wg.Wait()

	// check whether main key conflicts
	excelFileRaws.Range(func(k, v any) bool {
		filename := k.(string)
		excelFileRaw := v.(*ExcelFileRaw)
		fn := filename
		raw := excelFileRaw
		wg.Wrap(func() {
			uniqueKey := make(map[string]bool)
			multiKey := len(raw.Keys) > 1
			for ind, rowData := range raw.RowData {
				if len(rowData) == 0 {
					continue
				}

				if multiKey {
					multiKeys := make([]string, 0, len(raw.Keys))
					for _, key := range raw.Keys {
						singleKey := cast.ToString(rowData[key])
						multiKeys = append(multiKeys, singleKey)
					}
					strMultiKeys := strings.Join(multiKeys, "_")
					if _, found := uniqueKey[strMultiKeys]; found {
						log.Error().Str("file_name", fn).Str("key_name", strMultiKeys).Msg("主键冲突")
						GlobalResultReport.AddError(fn, raw.Sheetname, strMultiKeys, "主键冲突"+fmt.Sprintf(" 行号:%d", ind))
					}
					uniqueKey[strMultiKeys] = true
				} else {
					strKey := cast.ToString(rowData[raw.Keys[0]])
					if _, found := uniqueKey[strKey]; found {
						log.Error().Str("file_name", fn).Str("key_name", strKey).Msg("主键冲突")
						GlobalResultReport.AddError(fn, raw.Sheetname, strKey, "主键冲突"+fmt.Sprintf(" 行号:%d", ind))
					}
					uniqueKey[strKey] = true
				}
			}
		})
		return true
	})
	wg.Wait()

	return replacedFilenames
}

// 预处理
func preParseJsonRows(lanType ExportLanguage_Type, dirPath, filename string, sheet *JsonSheet) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Str("file_name", filename).Str("sheet_name", sheet.SheetFile).Msg("preParseJsonRows failed")
		}
	}()

	fileRaw := &ExcelFileRaw{
		Type:         ExcelSheetType_Null,
		LanType:      lanType,
		Filename:     filename,
		Sheetname:    sheet.SheetFile,
		Beanname:     sheet.Bean,
		ConfFilename: sheet.ConfFile,
		Keys:         make([]string, 0, 16),
		FieldRaw:     treemap.NewWithStringComparator(),
		EnumData:     make([]*ExcelEnumType, 0, 16),
		RowData:      make([]ExcelRowData, 0, 16),
	}

	// enum file
	if sheet.Bean == "SYS_ENUM" {
		fileRaw.Type = ExcelSheetType_Enum
	} else if sheet.Bean == "SYS_POLY" {
		fileRaw.Type = ExcelSheetType_Poly
	} else if sheet.Mode == "one" {
		fileRaw.Type = ExcelSheetType_One
	} else if sheet.Mode == "map" {
		fileRaw.Type = ExcelSheetType_Map

		// 文件名有中文时，使用Bean作为文件名
		for _, runeValue := range fileRaw.Filename {
			if unicode.Is(unicode.Han, runeValue) {
				fileRaw.Filename = utils.SimpleTitle(sheet.Bean) + ".json"
				break
			}
		}
	}

	rows := make([][]any, len(sheet.Values))
	maxColLen := int(0)
	for n := 0; n < len(sheet.Values); n++ {
		rows[n] = make([]any, len(sheet.Values[n]))
		copy(rows[n], sheet.Values[n])
		if len(rows[n]) > maxColLen {
			maxColLen = len(rows[n])
		}
	}

	if sheet.Bean != "SYS_POLY" {
		for idx, row := range rows {
			for n := len(row); n < maxColLen; n++ {
				row = append(row, "")
			}
			rows[idx] = row
		}
	}

	if len(rows) == 0 {
		return
	}

	if sheet.Mode == "one" {
		// 矩阵转置后就是map类型
		newRows := make([][]any, len(rows[0]))
		for n := 0; n < len(newRows); n++ {
			newRows[n] = make([]any, len(rows))
		}

		for n := 0; n < len(rows); n++ {
			for m := 0; m < len(rows[0]); m++ {
				if n < len(rows) && rows[n] != nil && m < len(rows[n]) && rows[n][m] != nil {
					newRows[m][n] = rows[n][m]
				}
			}
		}

		rows = newRows
	}

	// 第一行确定是否导出字段
	canExports := make([]bool, len(rows[0]))
	var rowKeyIndex int

	// 第三行为导出字段名
	exportNames := make([]string, len(rows[2]))

	// 第四行为功能字段
	functionTypes := make([]map[string]string, len(rows[3]))

	for n := 0; n < len(rows); n++ {
		if rows[n] == nil {
			break
		}

		switch n {
		case 0:
			// 是否导出
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					continue
				}

				if _, ok := cell.(string); !ok {
					log.Error().Int32("行", int32(n+1)).Int32("列", int32(m+1)).Str("file_name", filename).Msg("第一行数据类型错误")
					return
				}

				if export, ok := cell.(string); ok {
					export = strings.ToLower(export)
					if lanType == ExportLanguage_Go {
						canExports[m] = strings.Contains(export, "s")
					} else if lanType == ExportLanguage_GoDynamic {
						canExports[m] = strings.Contains(export, "d")
					} else if lanType == ExportLanguage_Ts || lanType == ExportLanguage_Cs {
						canExports[m] = strings.Contains(export, "c")
					} else if lanType == ExportLanguage_Diff {
						canExports[m] = (strings.Contains(export, "c") ||
							strings.Contains(export, "s") ||
							strings.Contains(export, "d"))
					}
				}
			}

		case 1, 4:

		case 2:
			// 导出字段名
			idx := 0
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					continue
				}

				fieldName := cast.ToString(cell)

				// 去除space
				fieldName = strings.TrimSpace(fieldName)

				// 无字段名不导出，随机生成字段名字符串
				if len(fieldName) == 0 {
					fieldName = fmt.Sprintf("V%s", randstr.String(16))
				}

				exportNames[m] = utils.SimpleTitle(fieldName)
				idx++
			}

		case 3:
			// 功能字段
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				functionTypes[m] = make(map[string]string)

				if cell != nil {
					if _, ok := cell.(string); !ok {
						log.Error().Int32("行", int32(n+1)).Int32("列", int32(m+1)).Str("file_name", filename).Msg("第三行数据类型错误")
						return
					}

					for _, funcValue := range strings.Split(cell.(string), " ") {
						if len(funcValue) == 0 {
							continue
						}

						if strings.HasPrefix(funcValue, "main") {
							functionTypes[m]["main"] = "true"
						} else if strings.HasPrefix(funcValue, "keyAlias") {
							functionTypes[m]["keyAlias"] = "true"
						}
					}
				}

				// 默认第一列为主键
				if m == 0 || functionTypes[m]["main"] == "true" {
					fileRaw.Keys = append(fileRaw.Keys, exportNames[m])
					rowKeyIndex = m
				}

				if cell == nil {
					continue
				}
			}
			continue

		default:
			if len(rows[n]) == 0 {
				continue
			}

			// #开头为注释，跳过
			headerStr := cast.ToString(rows[n][0])
			if strings.HasPrefix(headerStr, "#") {
				continue
			}

			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					if m == 0 {
						break
					}

					cell = ""
				}

				cellValue := cast.ToString(cell)
				cellValue = strings.TrimSpace(cellValue)

				// 主键为空，跳过
				if m < len(functionTypes) && functionTypes[m]["main"] == "true" && len(cellValue) == 0 {
					continue
				}

				// 转换值参数
				keyAlias := false

				if len(functionTypes) > m && functionTypes[m] != nil {
					// 是否定义了主键别名
					if functionTypes[m]["keyAlias"] == "true" {
						keyAlias = true
					}
				}

				// 主键别名定义
				if keyAlias && len(cellValue) > 0 {
					realKeyValue := rows[n][rowKeyIndex]
					if value, found, global := excelKeyAliasMapping.GetValue(fileRaw.ConfFilename, cellValue); found && !global {
						GlobalResultReport.AddError(fileRaw.Filename, fileRaw.Sheetname, cellValue, fmt.Sprintf("主键别名%s冲突:%v, %v", cellValue, value, realKeyValue))
						continue
					}
					excelKeyAliasMapping.AddAlias(cellValue, fileRaw, realKeyValue)
					excelKeyRefMapping.AddRef(strings.Join(fileRaw.Keys, "_"), fileRaw, cellValue)
				}

				// 保存主键映射
				if m == 0 || (m < len(functionTypes) && functionTypes[m]["main"] == "true") {
					excelKeyRefMapping.AddRef(strings.Join(fileRaw.Keys, "_"), fileRaw, cellValue)
				}
			}
		}
	}
}

func parseJsonRows(lanType ExportLanguage_Type, dirPath, filename string, sheet *JsonSheet) (*ExcelFileRaw, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Str("file_name", filename).Str("sheet_name", sheet.SheetFile).Msg("parseJsonRows failed")
		}
	}()

	fileRaw := &ExcelFileRaw{
		Type:         ExcelSheetType_Null,
		LanType:      lanType,
		Filename:     filename,
		Sheetname:    sheet.SheetFile,
		Beanname:     sheet.Bean,
		ConfFilename: sheet.ConfFile,
		Keys:         make([]string, 0),
		FieldRaw:     treemap.NewWithStringComparator(),
		EnumData:     make([]*ExcelEnumType, 0, 16),
		RowData:      make([]ExcelRowData, 0, 16),
	}

	rows := make([][]any, len(sheet.Values))
	maxColLen := int(0)
	for n := 0; n < len(sheet.Values); n++ {
		rows[n] = make([]any, len(sheet.Values[n]))
		copy(rows[n], sheet.Values[n])
		if len(rows[n]) > maxColLen {
			maxColLen = len(rows[n])
		}
	}

	if sheet.Bean != "SYS_POLY" {
		for idx, row := range rows {
			for n := len(row); n < maxColLen; n++ {
				row = append(row, "")
			}
			rows[idx] = row
		}
	}

	// enum file
	if sheet.Bean == "SYS_ENUM" {
		fileRaw.Type = ExcelSheetType_Enum
		parseJsonEnumData(rows, fileRaw)

	} else if sheet.Bean == "SYS_POLY" {
		fileRaw.Type = ExcelSheetType_Poly
		parseJsonPolyData(rows, fileRaw)

	} else if sheet.Mode == "one" {
		fileRaw.Type = ExcelSheetType_One

		// 文件名有中文时，使用Bean作为文件名
		for _, runeValue := range fileRaw.Filename {
			if unicode.Is(unicode.Han, runeValue) {
				fileRaw.Filename = utils.SimpleTitle(sheet.Bean) + ".json"
				break
			}
		}

		// 矩阵转置后就是map类型
		newRows := make([][]any, len(rows[0]))
		for n := 0; n < len(newRows); n++ {
			newRows[n] = make([]any, len(rows))
		}

		for n := 0; n < len(rows); n++ {
			for m := 0; m < len(rows[0]); m++ {
				if n < len(rows) && rows[n] != nil && m < len(rows[n]) && rows[n][m] != nil {
					newRows[m][n] = rows[n][m]
				}
			}
		}
		parseJsonMapData(lanType, newRows, fileRaw)

	} else if sheet.Mode == "map" {
		fileRaw.Type = ExcelSheetType_Map

		// 文件名有中文时，使用Bean作为文件名
		for _, runeValue := range fileRaw.Filename {
			if unicode.Is(unicode.Han, runeValue) {
				fileRaw.Filename = utils.SimpleTitle(sheet.Bean) + ".json"
				break
			}
		}

		// if fileRaw.Filename == "IKActionAttributes.json" {
		parseJsonMapData(lanType, rows, fileRaw)
		// }

	} else {
		log.Fatal().Str("file", filename).Msg("invalid sheet type")
	}

	return fileRaw, nil
}

// parse enum sheet
func parseJsonEnumData(rows [][]any, fileRaw *ExcelFileRaw) {
	// skip first row
	var excelEnumType *ExcelEnumType
	for n := 1; n < len(rows); n++ {
		if rows[n] == nil {
			continue
		}

		var e ExcelEnum
		for m := 0; m < len(rows[n]); m++ {
			cell := rows[n][m]
			switch m {
			case 0:
				// 第0列为枚举类型名字
				if enumType, ok := cell.(string); ok && len(enumType) > 0 {
					// #开头为注释，跳过
					if strings.HasPrefix(enumType, "#") {
						continue
					}

					// 去除space
					enumType = strings.TrimSpace(enumType)
					if excelEnumType != nil && excelEnumType.EnumType != enumType {
						fileRaw.EnumData = append(fileRaw.EnumData, excelEnumType)
					}

					// 新的枚举类型
					if excelEnumType == nil || excelEnumType.EnumType != enumType {
						excelEnumType = &ExcelEnumType{
							EnumType:       utils.SimpleTitle(enumType),
							EnumOriginType: enumType,
							Options:        "",
							Enums:          make([]ExcelEnum, 0, 32),
						}

					}
				}
			case 1:
				// 第1列为枚举名
				enumName, ok := cell.(string)
				if !ok {
					break
				}
				// 去除space
				enumName = strings.TrimSpace(enumName)
				e.EnumName = strings.Title(enumName)
			case 2:
				// 第2列为枚举描述和别名
				enumDesc, ok := cell.(string)
				if !ok {
					break
				}

				// 去除换行
				enumDesc = strings.ReplaceAll(enumDesc, "\n", ",")
				e.EnumDesc = enumDesc
			case 3:
				// 第3列为枚举值
				enumValue, err := cast.ToStringE(cell)
				if err != nil {
					break
				}
				e.EnumValue = enumValue
			case 4:
				optionValue, ok := cell.(string)
				if !ok {
					break
				}
				if excelEnumType != nil {
					excelEnumType.Options = normalizeSheetOptions(excelEnumType.Options + "," + optionValue)
				}
			default:
				break
			}
		}

		if excelEnumType == nil {
			log.Fatal().Str("file", fileRaw.Filename).Msg("invalid excel enum type")
		}

		if len(e.EnumName) > 0 {
			excelEnumType.Enums = append(excelEnumType.Enums, e)
		}
	}

	if excelEnumType != nil {
		fileRaw.EnumData = append(fileRaw.EnumData, excelEnumType)
	}
}

// parse poly sheet
func parseJsonPolyData(rows [][]any, fileRaw *ExcelFileRaw) {
	// skip first row
	var excelPolyType *ExcelPolyType
	var excelPolyComment string
	hasOptionsColumn := len(rows) > 0 && len(rows[0]) > 9 && strings.EqualFold(cast.ToString(rows[0][9]), "options")
lineLoop:
	for n := 1; n < len(rows); n++ {
		if rows[n] == nil {
			continue
		}

		var e ExcelPoly
		for m := 0; m < len(rows[n]); m++ {
			cell := rows[n][m]
			switch m {
			case 0:
				// 第0列为poly类型名字
				if polyType, ok := cell.(string); ok && len(polyType) > 0 {
					// #开头为注释，跳过
					if strings.HasPrefix(polyType, "#") {
						excelPolyComment = polyType
						continue lineLoop
						// break
					}

					// 去除space
					polyType = strings.TrimSpace(polyType)
					if excelPolyType != nil && excelPolyType.PolyType != polyType {
						fileRaw.PolyData = append(fileRaw.PolyData, excelPolyType)
					}

					// 新的poly类型
					if excelPolyType == nil || excelPolyType.PolyType != polyType {
						excelPolyType = &ExcelPolyType{
							PolyType:      strings.Title(polyType),
							PolySeperator: ",",
							PolyComment:   excelPolyComment,
							Options:       "",
							PolyFilename:  fileRaw.Filename,
							Polys:         make([]ExcelPoly, 0, 16),
							OriginPolys:   make([]ExcelPoly, 0, 16),
						}

					}
				}
			case 1:
				// 第1列为父类
				polyParent, ok := cell.(string)
				if !ok || len(polyParent) == 0 {
					break
				}

				polyParent = strings.TrimSpace(polyParent)
				excelPolyType.PolyParentName = strings.Title(polyParent)
			case 2:
				// 第2列为分隔符
				polySep, ok := cell.(string)
				if !ok {
					break
				}

				polySep = strings.TrimSpace(polySep)

				// 默认分隔符为逗号
				if len(polySep) == 0 {
					continue
				}

				excelPolyType.PolySeperator = polySep
			case 3:
				// 第3列为子变量名
				polyName, ok := cell.(string)
				if !ok {
					break
				}
				// 去除space
				polyName = strings.TrimSpace(polyName)
				e.PolyName = strings.Title(polyName)
				e.PolyOriginName = polyName
				tagName := utils.SimpleUntitle(e.PolyName)
				e.PolyTag = fmt.Sprintf("`json:\"%s,omitempty\" bson:\"%s\"`", tagName, tagName)
			case 4:
				// 第4列为子变量类型
				polyType, ok := cell.(string)
				if !ok {
					break
				}

				// 去除space
				polyType = strings.TrimSpace(polyType)
				e.PolyType = polyType
			case 5:
				// 第5列为子变量分隔符
				polySep, ok := cell.(string)
				if !ok {
					break
				}

				polySep = strings.TrimSpace(polySep)

				// 默认分隔符为逗号
				if len(polySep) == 0 {
					polySep = ","
				}

				e.PolySeperator = polySep
			case 6:
				polyAlias, ok := cell.(string)
				if !ok {
					continue
				}

				e.PolyAlias = polyAlias

			case 7:
				polyUnique, ok := cell.(string)
				if !ok {
					break
				}

				if len(excelPolyType.PolyUnique) > 0 {
					break
				}

				polyUnique = strings.TrimSpace(polyUnique)
				excelPolyType.PolyUnique = polyUnique

				if polyUnique == "loot" {
					excelPolyType.PolyLootKeyType = rows[n][4].(string)
				}

			case 8:
				// 8列暂时跳过
				break

			case 9:
				if hasOptionsColumn {
					polyOptions, ok := cell.(string)
					if !ok || len(polyOptions) == 0 {
						break
					}
					if excelPolyType != nil && len(e.PolyName) == 0 {
						excelPolyType.Options = normalizeSheetOptions(excelPolyType.Options + "," + polyOptions)
					}
					break
				}

				// 第9列为注释
				polyComment, ok := cell.(string)
				if !ok || len(polyComment) == 0 {
					break
				}

				// 去除space
				e.PolyComment = strings.TrimSpace(polyComment)
			case 10:
				if !hasOptionsColumn {
					break
				}
				polyComment, ok := cell.(string)
				if !ok || len(polyComment) == 0 {
					break
				}

				e.PolyComment = strings.TrimSpace(polyComment)
			default:
				break
			}
		}

		if excelPolyType == nil {
			continue
		}

		if len(e.PolyName) > 0 {
			excelPolyType.Polys = append(excelPolyType.Polys, e)
			excelPolyType.OriginPolys = append(excelPolyType.OriginPolys, e)
		}
	}

	if excelPolyType != nil {
		fileRaw.PolyData = append(fileRaw.PolyData, excelPolyType)
	}

	// 多态类型加上父类型的所有字段
	for _, childPolyData := range fileRaw.PolyData {
		// 没有父类型不处理
		if len(childPolyData.PolyParentName) == 0 {
			continue
		}

		// 有父类型的找到父类型多态属性
		for _, parentPolyData := range fileRaw.PolyData {
			if parentPolyData.PolyType == childPolyData.PolyParentName {
				combinePolys := make([]ExcelPoly, 0, len(childPolyData.Polys)+len(parentPolyData.Polys))
				combinePolys = append(combinePolys, parentPolyData.Polys...)
				combinePolys = append(combinePolys, childPolyData.Polys...)
				childPolyData.Polys = combinePolys
				break
			}
		}
	}
}

// 内嵌的多态结构
type embedPolyData struct {
	TypeName    string   // json中的类型名
	ExportName  string   // 导出字段名
	PointerName string   // 生成的对象类型名
	SubValues   []string // 子类型值
}

// parse map config sheet
func parseJsonMapData(lanType ExportLanguage_Type, rows [][]any, fileRaw *ExcelFileRaw) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().
				Str("filename", fileRaw.Filename).
				Str("sheetname", fileRaw.Sheetname).
				Msgf("解析json配置表失败：%v", err)
		}
	}()

	if len(rows) == 0 {
		return
	}

	if len(rows) < 5 {
		log.Fatal().
			Str("filename", fileRaw.Filename).
			Str("sheetname", fileRaw.Sheetname).
			Msg("表格基础行数缺失：第一行为导出客户端或服务器控制，第二行为导出类型，第三行为导出字段名，第四行为功能字段，第五行为字段描述")
		return
	}

	// 第一行确定是否导出字段
	canExports := make([]bool, len(rows[0]))

	// 第二行为导出类型
	originalExportTypes := make([]string, len(rows[1]))
	exportTypes := make([]string, len(rows[1]))
	isPolyTypes := make([]bool, len(rows[1]))
	isPolyArrayTypes := make([]bool, len(rows[1]))

	// 第三行为导出字段名
	exportNames := make([]string, len(rows[2]))

	// 第四行为功能字段
	functionTypes := make([]map[string]string, len(rows[3]))

	// 第五行为字段描述
	exportDescs := make([]string, len(rows[4]))

	for n := 0; n < len(rows); n++ {
		if rows[n] == nil {
			break
		}

		switch n {
		case 0:
			shouldExport := false

			// 是否导出
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					continue
				}

				if export, ok := cell.(string); ok {
					export = strings.ToLower(export)
					if lanType == ExportLanguage_Go {
						canExports[m] = strings.Contains(export, "s")
					} else if lanType == ExportLanguage_GoDynamic {
						canExports[m] = strings.Contains(export, "d")
					} else if lanType == ExportLanguage_Ts || lanType == ExportLanguage_Cs {
						canExports[m] = strings.Contains(export, "c")
					} else if lanType == ExportLanguage_Diff {
						canExports[m] = (strings.Contains(export, "c") ||
							strings.Contains(export, "s") ||
							strings.Contains(export, "d"))
					}

					if canExports[m] {
						shouldExport = true
					}
				}
			}

			fileRaw.ShouldExport = shouldExport

		case 1:
			// 导出类型
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					continue
				}

				exportType := cast.ToString(cell)

				// 去除space
				exportType = strings.TrimSpace(exportType)
				originalExportTypes[m] = exportType
				tp, _ := getFullTypeName(exportType)
				exportTypes[m] = tp

				// 多态类型类型处理
				if isPolyTypeName(exportType) && PolyTypeNameHasParent(exportType) {
					isPolyTypes[m] = true
					exportTypes[m] = "any"
				}

				// 多态数组类型
				trimListExportType := strings.TrimPrefix(exportType, "list:")
				if isPolyTypeName(trimListExportType) && PolyTypeNameHasParent(trimListExportType) {
					isPolyArrayTypes[m] = true
					exportTypes[m] = "[]any"
				}
			}

		case 2:
			// 导出字段名
			idx := 0
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					continue
				}

				if m >= len(exportTypes) {
					continue
				}

				fieldName := cast.ToString(cell)

				// 去除space
				fieldName = strings.TrimSpace(fieldName)

				// 字段属性
				raw := &ExcelFieldRaw{
					originName: fieldName,
					tp:         exportTypes[m],
					originTp:   originalExportTypes[m],
					idx:        idx,
					export:     true,
				}

				// 无字段名不导出，随机生成字段名字符串
				if len(fieldName) == 0 {
					raw.export = false
					fieldName = fmt.Sprintf("V%s", randstr.String(16))
				}

				// 无字段类型不导出
				if len(exportTypes[m]) == 0 {
					raw.export = false
				}

				// 控制是否导出
				if m < len(canExports) && !canExports[m] {
					raw.export = false
				}

				raw.name = utils.SimpleTitle(fieldName)
				exportNames[m] = raw.name

				if _, found := fileRaw.FieldRaw.Get(raw.name); found {
					continue
				}

				tagName := utils.SimpleUntitle(raw.name)
				raw.tag = fmt.Sprintf("`json:\"%s,omitempty\" bson:\"%s\"`", tagName, tagName)
				fileRaw.FieldRaw.Put(raw.name, raw)
				idx++
			}

			// 导出类型连续声明为多态类型的，也需要将类型置为any
			var fieldName string
			for idx, name := range originalExportTypes {
				trimListName := strings.TrimPrefix(name, "list:")
				if isPolyTypeName(name) {
					if fieldName == exportNames[idx] {
						isPolyTypes[idx-1] = true
						isPolyTypes[idx] = true
						exportTypes[idx-1] = "any"
						exportTypes[idx] = "any"

						// 修正字段类型为any
						// if v, found := fileRaw.FieldRaw.Get(fieldName); found {
						// 	v.(*ExcelFieldRaw).tp = "any"
						// }
					}
					fieldName = exportNames[idx]
				} else if isPolyTypeName(trimListName) {
					if fieldName == exportNames[idx] {
						isPolyArrayTypes[idx-1] = true
						isPolyArrayTypes[idx] = true
						exportTypes[idx-1] = "[]any"
						exportTypes[idx] = "[]any"

						// 修正字段类型为any
						// if v, found := fileRaw.FieldRaw.Get(fieldName); found {
						// 	v.(*ExcelFieldRaw).tp = "[]any"
						// }
					}
					fieldName = exportNames[idx]
				} else {
					fieldName = ""
				}
			}
		case 3:
			// 功能字段
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				functionTypes[m] = make(map[string]string)

				if cell != nil {
					if _, ok := cell.(string); !ok {
						log.Error().Str("单元格内容", cast.ToString(cell)).Int32("行", int32(n+1)).Int32("列", int32(m+1)).Str("file_name", fileRaw.Filename).Msg("第三行数据类型错误")
						return
					}

					for _, funcValue := range strings.Split(cell.(string), " ") {
						if len(funcValue) == 0 {
							continue
						}

						if strings.HasPrefix(funcValue, "main") {
							functionTypes[m]["main"] = "true"
						} else if strings.HasPrefix(funcValue, "option") {
							functionTypes[m]["option"] = "true"
						} else if strings.HasPrefix(funcValue, "sep") {
							sepValues := strings.Split(funcValue, "=")
							functionTypes[m]["sep"] = sepValues[len(sepValues)-1]
						} else if strings.HasPrefix(funcValue, "def") {
							defValues := strings.Split(funcValue, "=")
							functionTypes[m]["def"] = defValues[len(defValues)-1]
						} else if strings.HasPrefix(funcValue, "ref") {
							refValues := strings.Split(funcValue, "=")
							functionTypes[m]["ref"] = refValues[len(refValues)-1]
						} else if strings.HasPrefix(funcValue, "keyAlias") {
							functionTypes[m]["keyAlias"] = "true"
						} else if strings.HasPrefix(funcValue, "alias") {
							aliasValues := strings.Split(funcValue, "=")
							functionTypes[m]["alias"] = aliasValues[len(aliasValues)-1]
						}
					}
				}

				// 默认第一列为主键
				if m == 0 || functionTypes[m]["main"] == "true" {
					fileRaw.Keys = append(fileRaw.Keys, exportNames[m])
				}

				if cell == nil {
					continue
				}

				if functionTypes[m]["option"] == "true" {
					if value, found := fileRaw.FieldRaw.Get(exportNames[m]); found {
						value.(*ExcelFieldRaw).option = true
					}
				}
			}
			continue
		case 4:
			// 字段描述
			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					continue
				}

				if m >= len(exportNames) {
					continue
				}

				fieldDesc := cast.ToString(cell)

				// 去除space
				fieldDesc = strings.TrimSpace(fieldDesc)
				fieldDesc = strings.ReplaceAll(fieldDesc, "\n", "")

				// 如果是多态类型，则在描述中添加指针类型
				if m < len(isPolyTypes) && isPolyTypes[m] &&
					(fileRaw.LanType == ExportLanguage_Go || fileRaw.LanType == ExportLanguage_GoDynamic) {
					fieldDesc = fmt.Sprintf("%s 指针类型 使用时需要强转", fieldDesc)
				}

				// 如果是多态列表类型，则在描述中添加指针数组类型
				if m < len(isPolyArrayTypes) && isPolyArrayTypes[m] &&
					(fileRaw.LanType == ExportLanguage_Go || fileRaw.LanType == ExportLanguage_GoDynamic) {
					fieldDesc = fmt.Sprintf("%s 数组指针类型 使用时需要强转", fieldDesc)
				}

				exportDescs[m] = fieldDesc

				if value, found := fileRaw.FieldRaw.Get(exportNames[m]); found {
					value.(*ExcelFieldRaw).desc = fieldDesc
				}
			}

		default:
			if len(rows[n]) == 0 {
				continue
			}

			// #开头为注释，跳过
			headerStr := cast.ToString(rows[n][0])
			if strings.HasPrefix(headerStr, "#") {
				continue
			}

			// 普通数据行处理
			mapRowData := make(map[string]any)

			// 多态数据行处理
			embedPolyDatas := make(map[string]*embedPolyData)
			embedPolyArrayDatas := make(map[string]*embedPolyData)

			for m := 0; m < len(rows[n]); m++ {
				cell := rows[n][m]
				if cell == nil {
					if m == 0 {
						break
					}

					cell = ""
				}

				if m >= len(canExports) || !canExports[m] {
					continue
				}

				cellValue := cast.ToString(cell)
				cellValue = strings.TrimSpace(cellValue)

				if len(cellValue) == 0 && m == 0 {
					break
				}

				// 主键为空，跳过
				if m < len(functionTypes) && functionTypes[m]["main"] == "true" && len(cellValue) == 0 {
					continue
				}

				if isPolyTypes[m] {
					// 内嵌多态类型处理
					exportName := exportNames[m]
					if len(exportName) > 0 {
						polyData, found := embedPolyDatas[exportName]
						if !found {
							polyData = &embedPolyData{
								TypeName:   originalExportTypes[m],
								ExportName: exportName,
								SubValues:  make([]string, 0),
							}

							polyData.PointerName, _ = getFullTypeName(cellValue)
							embedPolyDatas[exportName] = polyData
						} else {
							polyData.SubValues = append(polyData.SubValues, cellValue)
						}
					}

				} else if isPolyArrayTypes[m] {
					// 内嵌数组多态类型处理
					exportName := exportNames[m]
					if len(exportName) > 0 {
						polyArrayData, found := embedPolyArrayDatas[exportName]
						if !found {
							polyArrayData = &embedPolyData{
								TypeName:    originalExportTypes[m],
								ExportName:  exportNames[m],
								PointerName: fmt.Sprintf("*%s", exportNames[m]),
								SubValues:   make([]string, 0),
							}
							embedPolyArrayDatas[exportName] = polyArrayData
						}
						polyArrayData.SubValues = append(polyArrayData.SubValues, cellValue)
					}

				} else {

					var err error
					// 转换值参数
					seperator := ""
					option := true
					defval := ""
					refVal := ""
					aliasVal := ""

					if len(functionTypes) > m && functionTypes[m] != nil {
						// 是否定义了分隔符
						if functionTypes[m]["sep"] != "" {
							seperator = functionTypes[m]["sep"]
						}

						// 是否定义了option
						if functionTypes[m]["option"] == "true" {
							option = false
						}

						// 是否定义了默认值
						if functionTypes[m]["def"] != "" {
							defval = functionTypes[m]["def"]
						}

						// 是否定义了引用表，定义了引用默认定义主键别名
						if functionTypes[m]["ref"] != "" {
							refVal = functionTypes[m]["ref"]
							aliasVal = functionTypes[m]["ref"]
						}

						// 是否引用主键别名
						if functionTypes[m]["alias"] != "" {
							aliasVal = functionTypes[m]["alias"]
						}
					}

					// 普通类型处理
					if option || len(cellValue) > 0 || len(defval) > 0 {
						exportType := exportTypes[m]
						beConvertValue := cellValue
						if len(defval) > 0 && len(cellValue) == 0 {
							beConvertValue = defval
						}

						mapRowData[exportNames[m]], err = convertValue(fileRaw.Filename, fileRaw.Sheetname, exportNames[m], exportType, beConvertValue, seperator, aliasVal)
						if err != nil {
							GlobalResultReport.AddWarnings(fileRaw.Filename, fileRaw.Sheetname, exportNames[m], err.Error(), cast.ToString(rows[n][0]), cellValue)
						}

						// 检查引用
						if len(refVal) > 0 {
							curKey := cast.ToString(rows[n][m])
							if len(curKey) > 0 {
								value, found := excelEnumTypeRaws.Load(refVal)
								if found {
									// 枚举引用
									enumRaw := value.(*ExcelEnumRaw)
									_, found = enumRaw.GetEnumValue(curKey)
									if !found {
										GlobalResultReport.AddError(fileRaw.Filename, fileRaw.Sheetname, curKey, fmt.Sprintf("未找到引用的枚举类型:%s, 枚举值:%s", refVal, curKey))
										continue
									}
								} else {
									// 表引用
									_, found, _ := excelKeyRefMapping.GetValue(refVal, curKey)
									if !found {
										GlobalResultReport.AddError(fileRaw.Filename, fileRaw.Sheetname, curKey, fmt.Sprintf("未找到引用的表格:%s", refVal))
										continue
									}
								}
							}
						}
					}
				}
			}

			// 处理多态数据
			if len(embedPolyDatas) > 0 {
				for _, embedPolyData := range embedPolyDatas {
					// 所有多态数据合并为一个value数据传入convertValue解析
					seperator := getPolySeparator(embedPolyData.TypeName)
					value, err := convertValue(fileRaw.Filename, fileRaw.Sheetname, embedPolyData.ExportName, embedPolyData.PointerName, strings.Join(embedPolyData.SubValues, seperator))
					if err != nil {
						GlobalResultReport.AddWarnings(fileRaw.Filename, fileRaw.Sheetname, embedPolyData.ExportName, err.Error(), "", strings.Join(embedPolyData.SubValues, seperator))
					}
					if v, ok := value.(map[string]any); ok {
						// value多塞一个多态指针类型名
						v["__pointer__type"] = embedPolyData.PointerName
					}
					mapRowData[embedPolyData.ExportName] = value
				}
			}

			// 处理数组多态数据
			if len(embedPolyArrayDatas) > 0 {
				for _, embedPolyData := range embedPolyArrayDatas {
					if len(embedPolyData.SubValues) == 0 {
						continue
					}

					type SubValue struct {
						Key    string
						Values []string
					}

					var subValuesList []SubValue

					for _, subValue := range embedPolyData.SubValues {
						if len(subValue) == 0 {
							continue
						}

						if strings.HasPrefix(subValue, "{") || strings.HasPrefix(subValue, "}") {
							continue
						}

						if isPolyTypeName(subValue) {
							newKey := subValue
							subValuesList = append(subValuesList, SubValue{Key: newKey, Values: []string{}})
						} else {
							if len(subValuesList) == 0 {
								GlobalResultReport.AddError(fileRaw.Filename, fileRaw.Sheetname, headerStr, "内嵌多态类型缺少子类类型声明")
								continue
							}

							lastValue := &subValuesList[len(subValuesList)-1]
							lastValue.Values = append(lastValue.Values, subValue)
						}
					}

					rowValues := make([]any, 0)
					for _, subValue := range subValuesList {
						// 分隔符
						seperator := "@"
						if value, found := excelPolyTypeRaws.Load(subValue.Key); found {
							raw := value.(*ExcelPolyType)
							seperator = raw.PolySeperator
						}

						// 所有多态数据合并为一个value数据传入convertValue解析
						pointerName := fmt.Sprintf("*%s", subValue.Key)
						value, err := convertValue(fileRaw.Filename, fileRaw.Sheetname, embedPolyData.ExportName, pointerName, strings.Join(subValue.Values, seperator))
						if err != nil {
							GlobalResultReport.AddWarnings(fileRaw.Filename, fileRaw.Sheetname, embedPolyData.ExportName, err.Error(), "", strings.Join(subValue.Values, seperator))
						}
						if v, ok := value.(map[string]any); ok {
							// value多塞一个多态指针类型名
							v["__pointer__type"] = pointerName
						}
						rowValues = append(rowValues, value)
					}
					mapRowData[embedPolyData.ExportName] = rowValues
				}
			}

			if len(mapRowData) > 0 {
				fileRaw.RowData = append(fileRaw.RowData, mapRowData)
			}

			// 保存主键映射
			// keyValues := make([]string, 0, len(fileRaw.Keys))
			// for _, key := range fileRaw.Keys {
			// 	keyValues = append(keyValues, cast.ToString(mapRowData[key]))
			// }
			// fileRaw.KeyMappings[strings.Join(keyValues, "_")] = struct{}{}
		}
	}
}

func formatAndImport(exportGoPath string) {
	// execute gofmt -w ../server/excel/auto/*.go
	cmd := exec.Command("bash", "-c", "gofmt", "-w", fmt.Sprintf("%s/*.go", exportGoPath))
	ret, err := cmd.Output()
	if err != nil {
		log.Fatal().Str("output", string(ret)).Err(err).Msg("gofmt failed")
		return
	}

	// execute goimports -w ../server/excel/auto/*.go
	cmd = exec.Command("bash", "-c", "goimports", "-w", fmt.Sprintf("%s/*.go", exportGoPath))
	ret, err = cmd.Output()
	if err != nil {
		log.Fatal().Str("output", string(ret)).Err(err).Msg("goimports failed")
	}
}

func GenerateFromJson(lanType ExportLanguage_Type, readJsonPath, exportPath string) {
	fileNames := getAllJsonFileNames(readJsonPath)
	fileNames = loadJson(lanType, readJsonPath, fileNames)

	switch lanType {
	case ExportLanguage_Go:
		// generate go files
		generateAllGoCodes(exportPath, fileNames)

		// format and import
		// formatAndImport(exportGoPath)
	case ExportLanguage_GoDynamic:
		// generate go files
		generateAllGoDynamicCodes(exportPath, fileNames)

	case ExportLanguage_Lua:
		// generate lua files
		generateAllLuaCodes(exportPath, fileNames)

		// todo format

	case ExportLanguage_Ts, ExportLanguage_Diff:
		// generate ts files
		generateAllTsCodes(exportPath, fileNames)

	case ExportLanguage_Cs:
		// generate cs files
		generateAllCsCodes(exportPath, fileNames)

	}
}
