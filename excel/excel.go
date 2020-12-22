package excel

import (
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var (
	RowOffset int = 2 // 第一行数据偏移
	ColOffset int = 2 // 第一列数据偏移
)

type ExcelProto struct {
	ID      int    `json:"Id"`
	Name    string `json:"Name,omitempty"`
	AttID   int    `json:"AttID,omitempty"`
	Quality int    `json:"Quality,omitempty"`
	AttList []int  `json:"AttList,omitempty"`
}

type ExcelProtoConfig struct {
	Rows map[int]*ExcelProto
}

func GenerateExcelFile(path string, typeName string) error {
	xlsxFile, err := excelize.OpenFile(path)
	if utils.ErrCheck(err, "open faile failed", path, typeName) {
		return err
	}

	rows, err := xlsxFile.GetRows(xlsxFile.GetSheetName(0))
	if utils.ErrCheck(err, "get rows failed", path, typeName) {
		return err
	}

	excelProtoConfig := &ExcelProtoConfig{
		Rows: make(map[int]*ExcelProto),
	}

	typeDescs := make([]string, len(rows[0])-ColOffset)
	typeNames := make([]string, len(rows[0])-ColOffset)
	typeValues := make([]string, len(rows[0])-ColOffset)
	for n := 0; n < len(rows); n++ {
		// load type desc
		if n == RowOffset {
			for m := ColOffset; m < len(rows[n]); m++ {
				typeDescs[m-ColOffset] = rows[n][m]
			}
		}

		// load type name
		if n == RowOffset+1 {
			for m := ColOffset; m < len(rows[n]); m++ {
				typeNames[m-ColOffset] = rows[n][m]
			}
		}

		// load type value
		if n == RowOffset+2 {
			for m := ColOffset; m < len(rows[n]); m++ {
				typeValues[m-ColOffset] = rows[n][m]
			}
		}

		// there is no actual data before row:5
		if n < RowOffset+3 {
			continue
		}

		var id int
		mapRowData := make(map[string]interface{})
		for m := ColOffset; m < len(rows[n]); m++ {
			cellColIdx := m - ColOffset
			cellValString := rows[n][m]

			// set value
			convertedVal := convertValue(typeValues[cellColIdx], cellValString)
			mapRowData[typeNames[cellColIdx]] = convertedVal

			// set Id
			if typeNames[cellColIdx] == "Id" {
				id = convertedVal.(int)
			}

		}

		excelProto := &ExcelProto{}
		err := mapstructure.Decode(mapRowData, excelProto)
		if err != nil {
			log.Error().Err(err).Int("row", n).Msg("decode excel data to struct failed")
			continue
		}

		excelProtoConfig.Rows[id] = excelProto
	}

	return nil
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
