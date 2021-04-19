package excel

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"bitbucket.org/funplus/server/utils"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/emirpasic/gods/maps/treemap"
	map_utils "github.com/emirpasic/gods/utils"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/thanhpk/randstr"
)

var (
	RowOffset int = 2 // 第一行数据偏移
	ColOffset int = 2 // 第一列数据偏移
)

var (
	entryLoaders       sync.Map                 // all entry loaders
	entryManualLoaders sync.Map                 // all entry manual loaders
	excelFileRaws      map[string]*ExcelFileRaw // all excel file raw data
)

type ExcelRowData map[string]interface{}

// Entry should implement Load function
type EntryLoader interface {
	Load(*ExcelFileRaw) error
}

// Entry should implement ManualLoad
type EntryManualLoader interface {
	ManualLoad(*ExcelFileRaw) error
}

// Excel field raw data
type ExcelFieldRaw struct {
	name string
	tp   string
	desc string
	tag  string
	idx  int  // field index in excel file
	imp  bool // need import
}

// Excel file raw data
type ExcelFileRaw struct {
	Filename string
	Keys     []string
	HasMap   bool
	FieldRaw *treemap.Map
	CellData []ExcelRowData
}

func init() {
	excelFileRaws = make(map[string]*ExcelFileRaw, 200)
}

func AddEntryLoader(name string, e EntryLoader) {
	entryLoaders.Store(name, e)
}

func AddEntryManualLoader(name string, e EntryManualLoader) {
	entryManualLoaders.Store(name, e)
}

func loadOneExcelFile(dirPath, filename string) (*ExcelFileRaw, error) {
	filePath := fmt.Sprintf("%s%s", dirPath, filename)
	xlsxFile, err := excelize.OpenFile(filePath)
	if !utils.ErrCheck(err, "open file failed", filePath) {
		return nil, err
	}

	rows, err := xlsxFile.GetRows(xlsxFile.GetSheetName(0))
	if !utils.ErrCheck(err, "get rows failed", filePath) {
		return nil, err
	}

	fileRaw := &ExcelFileRaw{
		Filename: filename,
		FieldRaw: treemap.NewWithStringComparator(),
		CellData: make([]ExcelRowData, 0),
	}

	// for _, v := range filename {
	// 	fileRaw.Filename = string(unicode.ToLower(v)) + filename[1:]
	// 	break
	// }

	// rotate config excel files
	if strings.Contains(fileRaw.Filename, "Config") {
		newRows := make([][]string, len(rows[RowOffset]))
		for n := 0; n < len(newRows); n++ {
			newRows[n] = make([]string, len(rows))
		}

		for n := 0; n < len(rows); n++ {
			for m := 0; m < len(rows[RowOffset]); m++ {
				newRows[m][n] = rows[n][m]
			}
		}
		parseExcelData(newRows, fileRaw)
	} else {
		parseExcelData(rows, fileRaw)
	}

	return fileRaw, nil
}

func getAllExcelFileNames(readExcelPath string) []string {
	dir, err := ioutil.ReadDir(readExcelPath)
	if !utils.ErrCheck(err, "read dir failed", readExcelPath) {
		return []string{}
	}

	// escape dir and ~$***.xlsx
	fileNames := make([]string, 0, len(dir))
	for _, fi := range dir {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".xlsx") && !strings.HasPrefix(fi.Name(), "~$") {
			fileNames = append(fileNames, fi.Name())
		}
	}

	return fileNames
}

// load all excel files
func loadAllExcelFiles(dirPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	mu := sync.Mutex{}
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException(name)
			rowDatas, err := loadOneExcelFile(dirPath, name)
			utils.ErrPrint(err, "loadOneExcelFile failed", name)

			mu.Lock()
			excelFileRaws[name] = rowDatas
			mu.Unlock()
		})
	}
	wg.Wait()
}

// generate go code from excel file
func generateAllCodes(exportPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException(name)
			err := generateCode(exportPath, excelFileRaws[name])
			if !utils.ErrCheck(err, "generateCode failed", exportPath, name) {
				return
			}

			log.Info().Str("file_name", name).Str("export_dir", exportPath).Caller().Msg("generate go code success")
		})
	}

	wg.Wait()
}

func Generate(readExcelPath, exportPath string) {
	fileNames := getAllExcelFileNames(readExcelPath)
	loadAllExcelFiles(readExcelPath, fileNames)
	generateAllCodes(exportPath, fileNames)
}

// read all excel entries
func ReadAllEntries(dirPath string) {
	fileNames := getAllExcelFileNames(dirPath)
	loadAllExcelFiles(dirPath, fileNames)

	wg := utils.WaitGroupWrapper{}

	// read from excel files
	entryLoaders.Range(func(k, v interface{}) bool {
		entryName := k.(string)
		loader := v.(EntryLoader)

		wg.Wrap(func() {
			defer utils.CaptureException(entryName)
			err := loader.Load(excelFileRaws[entryName])
			utils.ErrPrint(err, "EntryLoader Load failed", entryName)
		})

		return true
	})
	wg.Wait()

	// load by manual
	entryManualLoaders.Range(func(k, v interface{}) bool {
		entryName := k.(string)
		loader := v.(EntryManualLoader)

		wg.Wrap(func() {
			defer utils.CaptureException(entryName)
			err := loader.ManualLoad(excelFileRaws[entryName])
			utils.ErrPrint(err, "EntryManualLoader Load failed", entryName)
		})

		return true
	})
	wg.Wait()

	log.Info().Msg("all excel entries reading completed!")
}

func parseExcelData(rows [][]string, fileRaw *ExcelFileRaw) {

	typeNames := make([]string, len(rows[2])-ColOffset)
	typeValues := make([]string, len(rows[2])-ColOffset)
	for n := 0; n < len(rows); n++ {
		if rows[n] == nil {
			break
		}

		// load type name
		if n == RowOffset {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := rows[n][m]
				raw := &ExcelFieldRaw{
					idx: m - ColOffset,
				}

				// 无字段名不导出，随机生成字段名字符串
				if len(fieldName) == 0 {
					raw.imp = false
					fieldName = randstr.String(16)
				}

				// 首字段
				if m == ColOffset {
					fieldName = "Id"
				}

				raw.name = strings.Title(fieldName)
				if _, found := fileRaw.FieldRaw.Get(raw.name); found {
					_ = utils.ErrCheck(errors.New("duplicate field name"), "parseExcelData failed", raw.name, fileRaw.Filename)
					continue
				}

				raw.tag = fmt.Sprintf("`json:\"%s,omitempty\"`", raw.name)
				fileRaw.FieldRaw.Put(raw.name, raw)
				typeNames[m-ColOffset] = raw.name
			}
		}

		// load type desc
		// if n == RowOffset+1 {

		// }

		// load type control
		if n == RowOffset+2 {
			var buffer bytes.Buffer
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := typeNames[m-ColOffset]
				fieldValue := rows[n][m]
				value, ok := fileRaw.FieldRaw.Get(fieldName)
				if !ok {
					log.Fatal().
						Caller().
						Str("filename", fileRaw.Filename).
						Str("fieldname", fieldName).
						Int("row", n).
						Int("col", m).
						Msg("parse excel data failed")
				}

				// 第一个字段默认主键
				if m == ColOffset {
					fileRaw.Keys = append(fileRaw.Keys, value.(*ExcelFieldRaw).name)

					// 去除换行
					desc := strings.Replace(value.(*ExcelFieldRaw).desc, "\n", ",", -1)

					buffer.Reset()
					buffer.WriteString(desc)
					buffer.WriteString(" 主键")
					value.(*ExcelFieldRaw).imp = true
					value.(*ExcelFieldRaw).desc = buffer.String()
					continue
				}

				// 带K标识的也是主键
				if strings.Contains(fieldValue, "K") {
					fileRaw.Keys = append(fileRaw.Keys, value.(*ExcelFieldRaw).name)

					// 去除换行
					desc := strings.Replace(value.(*ExcelFieldRaw).desc, "\n", ",", -1)
					buffer.Reset()
					buffer.WriteString(desc)
					buffer.WriteString(" 多主键之一")
					value.(*ExcelFieldRaw).desc = buffer.String()
				} else {
					// 去除换行
					desc := strings.Replace(rows[n-1][m], "\n", ",", -1)
					value.(*ExcelFieldRaw).desc = desc
				}

				if strings.Contains(fieldValue, "C") {
					value.(*ExcelFieldRaw).imp = false
				} else {
					value.(*ExcelFieldRaw).imp = true
				}
			}
		}

		// load type value
		if n == RowOffset+3 {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := typeNames[m-ColOffset]
				fieldValue := rows[n][m]
				convertType := convertType(fieldValue)

				value, ok := fileRaw.FieldRaw.Get(fieldName)
				if !ok {
					log.Fatal().
						Caller().
						Str("filename", fileRaw.Filename).
						Str("fieldname", fieldName).
						Int("row", n).
						Int("col", m).
						Msg("parse excel data failed")
				}

				if convertType == "*treemap.Map" {
					fileRaw.HasMap = true
				}

				if len(convertType) == 0 {
					value.(*ExcelFieldRaw).imp = false
				}

				value.(*ExcelFieldRaw).tp = convertType
				typeValues[m-ColOffset] = fieldValue
			}
		}

		// 客户端导出字段
		if n == RowOffset+2 {
			continue
		}

		// there is no actual data before row:7
		if n < RowOffset+4 {
			continue
		}

		// empty data row
		if len(rows[n][2]) == 0 {
			continue
		}

		// resize row
		if len(rows[n]) < len(rows[RowOffset]) {
			rows[n] = append(rows[n], make([]string, len(rows[RowOffset])-len(rows[n]))...)
		}
		rows[n] = rows[n][:len(rows[RowOffset])]

		mapRowData := make(map[string]interface{})
		for m := ColOffset; m < len(rows[n]); m++ {
			cellColIdx := m - ColOffset
			cellValString := rows[n][m]

			// set value
			convertedVal := convertValue(typeValues[cellColIdx], cellValString)
			mapRowData[typeNames[cellColIdx]] = convertedVal
		}

		fileRaw.CellData = append(fileRaw.CellData, mapRowData)
	}
}

// be tolerant with type names
func convertType(strType string) string {
	switch strType {
	case "String":
		fallthrough
	case "STRING":
		return "string"

	case "[]String":
		fallthrough
	case "String[]":
		fallthrough
	case "[]STRING":
		return "[]string"

	case "Int32":
		fallthrough
	case "Int":
		fallthrough
	case "INT":
		fallthrough
	case "int":
		return "int32"

	case "Number":
		fallthrough
	case "NUMBER":
		fallthrough
	case "number":
		return "decimal.Decimal"

	case "Float32":
		fallthrough
	case "Float":
		fallthrough
	case "FLOAT":
		fallthrough
	case "float":
		return "float32"

	case "[]Int32":
		fallthrough
	case "[]Int":
		fallthrough
	case "[]INT":
		fallthrough
	case "[]int":
		return "[]int32"

	case "[]Number":
		fallthrough
	case "[]NUMBER":
		fallthrough
	case "Number[]":
		fallthrough
	case "NUMBER[]":
		fallthrough
	case "number[]":
		fallthrough
	case "[]number":
		return "[]decimal.Decimal"

	case "Bool":
		fallthrough
	case "BOOL":
		return "bool"

	default:
		if strings.HasPrefix(strType, "map") || strings.HasPrefix(strType, "Map") {
			return "*treemap.Map"
		}

		return strType
	}
}

func convertValue(strType, strVal string) interface{} {
	var cellVal interface{}
	convertType := convertType(strType)

	switch convertType {
	case "int32":
		if len(strVal) == 0 {
			cellVal = int32(0)
		} else {
			cellVal = cast.ToInt32(strVal)
		}

	case "decimal.Decimal":
		if len(strVal) == 0 || strVal == "0" {
			cellVal = decimal.NewFromInt32(0)
		} else {
			cellVal, _ = decimal.NewFromString(strVal)
			// floatVal := cast.ToFloat64(strVal)
			// floatVal *= define.PercentBase
			// floatVal = math.Round(floatVal)
			// cellVal = int32(floatVal)
		}

	case "float32":
		if len(strVal) == 0 {
			cellVal = float32(0)
		} else {
			cellVal = cast.ToFloat32(strVal)
		}

	case "[]int32":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertValue("int32", v)
		}
		cellVal = arrVals

	case "[]decimal.Decimal":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertValue("number", v)
		}
		cellVal = arrVals

	case "[]float32":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertValue("float32", v)
		}
		cellVal = arrVals

	case "[]string":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertValue("string", v)
		}
		cellVal = arrVals

	case "*treemap.Map":
		cellVal = convertMapValue(strType, strVal)

	case "bool":
		if strings.Contains(strVal, "true") || strings.Contains(strVal, "True") || strings.Contains(strVal, "TRUE") {
			cellVal = true
			break
		}

		if strings.Contains(strVal, "false") || strings.Contains(strVal, "False") || strings.Contains(strVal, "FALSE") {
			cellVal = false
			break
		}

		cellVal = cast.ToBool(strVal)

	default:
		// default string value
		if len(strVal) == 0 {
			cellVal = ""
		} else {
			cellVal = strVal
		}
	}

	return cellVal
}

func convertMapValue(strType, strVal string) interface{} {
	// split type and value, example: map[int32]string => "int32" and "string"
	ts := strings.Split(strType, "[")
	t := ts[len(ts)-1]
	tt := strings.Split(t, "]")
	keyType := convertType(tt[0])
	valueType := convertType(tt[1])

	m := treemap.NewWith(func() map_utils.Comparator {
		switch keyType {
		case "int32":
			return map_utils.Int32Comparator
		case "string":
			return map_utils.StringComparator
		case "float32":
			return map_utils.Float32Comparator
		default:
			return map_utils.Int32Comparator
		}
	}())

	mapValues := strings.Split(strVal, ",")
	for _, oneMapValue := range mapValues {
		fields := strings.Split(oneMapValue, ":")
		if len(fields) < 2 {
			continue
		}

		k := convertValue(keyType, fields[0])
		v := convertValue(valueType, fields[1])
		m.Put(k, v)
	}

	return m
}
