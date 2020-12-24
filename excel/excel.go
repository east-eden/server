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
	Load() error
}

// Excel field raw data
type ExcelFieldRaw struct {
	name string
	tp   string
	desc string
	tag  string
}

// Excel file raw data
type ExcelFileRaw struct {
	filename string
	fieldRaw *treemap.Map
	rawData  []ExcelRowData
}

// type ExcelProto struct {
// 	ID      int    `json:"Id"`
// 	Name    string `json:"Name,omitempty"`
// 	AttID   int    `json:"AttID,omitempty"`
// 	Quality int    `json:"Quality,omitempty"`
// 	AttList []int  `json:"AttList,omitempty"`
// }

// type ExcelProtoConfig struct {
// 	Rows ExcelRaws `json:"Rows"`
// }

// func (c *ExcelProtoConfig) Load() error {
// 	return nil
// }

func init() {
	excelFileRaws = make(map[string]*ExcelFileRaw)
}

func AddEntries(e EntriesProto, name string) {
	allEntries.Store(e, name)
}

func loadAllGocodeEntries() {
	wg := utils.WaitGroupWrapper{}
	allEntries.Range(func(k, v interface{}) bool {
		entriesProto := k.(EntriesProto)
		entryName := v.(string)
		log.Info().Str("entry_name", entryName).Msg("begin loading excel files")

		wg.Wrap(func() {
			err := entriesProto.Load()
			if utils.ErrCheck(err, "gocode entry load failed", entryName) {
				log.Fatal().Err(err).Send()
			}
		})

		return true
	})
	wg.Wait()

	log.Info().Msg("all excel entries reading completed!")
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
		rawData:  make([]ExcelRowData, 0),
	}
	parseExcelData(rows, fileRaw)
	return fileRaw, nil
}

func GenerateAndLoad(dirPath string) error {
	dir, err := ioutil.ReadDir(dirPath)
	if utils.ErrCheck(err, "read dir failed", dirPath) {
		return err
	}

	fileNames := make([]string, 0, len(dir))
	for _, fi := range dir {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".xlsx") {
			fileNames = append(fileNames, fi.Name())
		}
	}

	wg := utils.WaitGroupWrapper{}
	mu := sync.Mutex{}

	// load all excel raw data
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

	// generate go code from excel file
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException()
			err := GenerateExcelGocode(dirPath, excelFileRaws[name])
			if utils.ErrCheck(err, "GenerateExcelGocode failed", dirPath, name) {
				log.Fatal().Err(err).Send()
			}
		})
	}

	wg.Wait()

	// load all excel entries
	loadAllGocodeEntries()

	// excelProtoConfig := &ExcelProtoConfig{
	// 	Rows: make(ExcelRaws),
	// }
	// for _, v := range rowDatas {
	// 	excelProto := &ExcelProto{}
	// 	err := mapstructure.Decode(v, excelProto)
	// 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	// 		return err
	// 	}

	// 	excelProtoConfig.Rows[excelProto.ID] = excelProto
	// }

	// log.Info().Interface("excel proto", excelProtoConfig).Msg("parse excel data success")

	return nil
}

func parseExcelData(rows [][]string, fileRaw *ExcelFileRaw) {

	typeNames := make([]string, len(rows[0])-ColOffset)
	typeValues := make([]string, len(rows[0])-ColOffset)
	for n := 0; n < len(rows); n++ {
		// load type name
		if n == RowOffset {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := rows[n][m]
				raw := &ExcelFieldRaw{
					name: fieldName,
					tag:  fmt.Sprintf("`json:\"%s\"`", fieldName),
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

				value.(*ExcelFieldRaw).tp = typeValue
				typeValues[m-ColOffset] = rows[n][m]
			}
		}

		// there is no actual data before row:5
		if n < RowOffset+3 {
			continue
		}

		mapRowData := make(map[string]interface{})
		for m := ColOffset; m < len(rows[n]); m++ {
			cellColIdx := m - ColOffset
			cellValString := rows[n][m]

			// set value
			convertedVal := convertValue(typeValues[cellColIdx], cellValString)
			mapRowData[typeNames[cellColIdx]] = convertedVal
		}

		fileRaw.rawData = append(fileRaw.rawData, mapRowData)
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

	default:
		// default string value
		cellVal = strVal
	}

	return cellVal
}
