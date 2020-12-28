package excel

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/rs/zerolog/log"
)

var (
	RowOffset int = 2 // 第一行数据偏移
	ColOffset int = 2 // 第一列数据偏移
)

var (
	allEntries    sync.Map                 // all auto generated entries
	excelFileRaws map[string]*ExcelFileRaw // all excel file raw data
)

type ExcelRowData map[string]interface{}

// Entries should implement Load function
type EntriesProto interface {
	load(excelFileRaw *ExcelFileRaw) error
}

// Excel field raw data
type ExcelFieldRaw struct {
	name string
	tp   string
	desc string
	tag  string
	def  string
	idx  int // field index in excel file
}

// Excel file raw data
type ExcelFileRaw struct {
	filename string
	fieldRaw *treemap.Map
	cellData []ExcelRowData
}

func init() {
	excelFileRaws = make(map[string]*ExcelFileRaw)
}

func AddEntries(name string, e EntriesProto) {
	allEntries.Store(name, e)
}

func loadOneExcelFile(dirPath, filename string) (*ExcelFileRaw, error) {
	filePath := fmt.Sprintf("%s/%s", dirPath, filename)
	xlsxFile, err := excelize.OpenFile(filePath)
	if utils.ErrCheck(err, "open faile failed", filePath) {
		return nil, err
	}

	rows, err := xlsxFile.GetRows(xlsxFile.GetSheetName(0))
	if utils.ErrCheck(err, "get rows failed", filePath) {
		return nil, err
	}

	fileRaw := &ExcelFileRaw{
		filename: filename,
		fieldRaw: treemap.NewWithStringComparator(),
		cellData: make([]ExcelRowData, 0),
	}
	parseExcelData(rows, fileRaw)
	return fileRaw, nil
}

func getAllExcelFileNames(dirPath string) []string {
	dir, err := ioutil.ReadDir(dirPath)
	if utils.ErrCheck(err, "read dir failed", dirPath) {
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
			defer utils.CaptureException()
			rowDatas, err := loadOneExcelFile(dirPath, name)
			if utils.ErrCheck(err, "loadOneExcelFile failed", name) {
				log.Fatal().Err(err).Send()
			}

			mu.Lock()
			excelFileRaws[name] = rowDatas
			mu.Unlock()
		})
	}
	wg.Wait()
}

// generate go code from excel file
func generateAllCodes(dirPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException()
			err := generateCode(dirPath, excelFileRaws[name])
			if utils.ErrCheck(err, "generateCode failed", dirPath, name) {
				log.Fatal().Err(err).Send()
			}
		})
	}

	wg.Wait()
}

func Generate(dirPath string) {
	fileNames := getAllExcelFileNames(dirPath)
	loadAllExcelFiles(dirPath, fileNames)
	generateAllCodes(dirPath, fileNames)
}

// read all excel entries
func ReadAllEntries(dirPath string) {
	fileNames := getAllExcelFileNames(dirPath)
	loadAllExcelFiles(dirPath, fileNames)

	wg := utils.WaitGroupWrapper{}
	allEntries.Range(func(k, v interface{}) bool {
		entryName := k.(string)
		entriesProto := v.(EntriesProto)
		log.Info().Str("excel_file", entryName).Msg("begin loading excel data")

		wg.Wrap(func() {
			err := entriesProto.load(excelFileRaws[entryName])
			if utils.ErrCheck(err, "gocode entry load failed", entryName) {
				log.Fatal().Err(err).Send()
			}
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
		// load type name
		if n == RowOffset {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := rows[n][m]
				raw := &ExcelFieldRaw{
					name: fieldName,
					tag:  fmt.Sprintf("`json:\"%s,omitempty\"`", fieldName),
					idx:  m - ColOffset,
				}
				fileRaw.fieldRaw.Put(fieldName, raw)
				typeNames[m-ColOffset] = fieldName
			}
		}

		// load type desc
		if n == RowOffset+1 {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := rows[n-1][m]
				desc := rows[n][m]
				value, ok := fileRaw.fieldRaw.Get(fieldName)
				if !ok {
					log.Fatal().
						Str("filename", fileRaw.filename).
						Str("fieldname", fieldName).
						Int("row", n).
						Int("col", m).
						Msg("parse excel data failed")
				}

				value.(*ExcelFieldRaw).desc = desc
			}
		}

		// load type value
		if n == RowOffset+2 {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := rows[n-2][m]
				typeValue := rows[n][m]

				value, ok := fileRaw.fieldRaw.Get(fieldName)
				if !ok {
					log.Fatal().
						Str("filename", fileRaw.filename).
						Str("fieldname", fieldName).
						Int("row", n).
						Int("col", m).
						Msg("parse excel data failed")
				}

				value.(*ExcelFieldRaw).tp = convertType(typeValue)
				typeValues[m-ColOffset] = typeValue
			}
		}

		// load default value
		if n == RowOffset+3 {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := rows[n-3][m]
				defaultValue := rows[n][m]

				value, ok := fileRaw.fieldRaw.Get(fieldName)
				if !ok {
					log.Fatal().
						Str("filename", fileRaw.filename).
						Str("fieldname", fieldName).
						Int("row", n).
						Int("col", m).
						Msg("parse excel data failed")
				}

				value.(*ExcelFieldRaw).def = defaultValue
			}
		}

		// there is no actual data before row:6
		if n < RowOffset+4 {
			continue
		}

		// empty data row
		if len(rows[n][2]) == 0 {
			continue
		}

		mapRowData := make(map[string]interface{})
		for m := ColOffset; m < len(rows[n]); m++ {
			cellColIdx := m - ColOffset
			cellValString := rows[n][m]

			fieldName := typeNames[cellColIdx]
			excelFieldRaw, ok := fileRaw.fieldRaw.Get(fieldName)
			if !ok {
				log.Fatal().
					Str("filename", fileRaw.filename).
					Str("fieldname", fieldName).
					Int("row", n).
					Int("col", m).
					Msg("parse excel data failed")
			}

			// set value
			var convertedVal interface{}
			if len(cellValString) == 0 {
				defValue := excelFieldRaw.(*ExcelFieldRaw).def
				if len(defValue) == 0 && typeValues[cellColIdx] != "string[]" && typeValues[cellColIdx] != "string" {
					log.Fatal().
						Str("filename", fileRaw.filename).
						Str("default_value", defValue).
						Int("row", n).
						Int("col", m).
						Msg("default value not assigned")
				}
				convertedVal = convertValue(typeValues[cellColIdx], excelFieldRaw.(*ExcelFieldRaw).def)
			} else {
				convertedVal = convertValue(typeValues[cellColIdx], cellValString)
			}
			mapRowData[typeNames[cellColIdx]] = convertedVal
		}

		fileRaw.cellData = append(fileRaw.cellData, mapRowData)
	}
}

func convertType(strType string) string {
	switch strType {
	case "float":
		return "float32"
	case "int[]":
		return "[]int"
	case "float[]":
		return "[]float32"
	case "string[]":
		return "[]string"
	default:
		return strType
	}
}

func convertValue(strType, strVal string) interface{} {
	var cellVal interface{}
	var err error

	switch strType {
	case "int":
		cellVal, err = strconv.Atoi(strVal)
		utils.ErrPrint(err, "convert cell value to int failed", strVal)

	case "float":
		cellVal, err = strconv.ParseFloat(strVal, 32)
		utils.ErrPrint(err, "convert cell value to float failed", strVal)

	case "int[]":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertValue("int", v)
		}
		cellVal = arrVals

	case "float[]":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertValue("float", v)
		}
		cellVal = arrVals

	case "string[]":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertValue("string", v)
		}
		cellVal = arrVals

	default:
		// default string value
		cellVal = strVal
	}

	return cellVal
}
