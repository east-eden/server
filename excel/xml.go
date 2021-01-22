package excel

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/rs/zerolog/log"
)

var (
	allXmlEntries sync.Map                 // all auto generated entries
	xmlFileRaws   map[string]*ExcelFileRaw // all xml file raw data
)

func init() {
	xmlFileRaws = make(map[string]*ExcelFileRaw)
}

func AddXmlEntries(name string, e EntriesProto) {
	allXmlEntries.Store(name, e)
}

func loadOneXmlFile(dirPath, filename string) (*ExcelFileRaw, error) {
	filePath := fmt.Sprintf("%s/%s", dirPath, filename)
	xmlFile, err := ioutil.ReadFile(filePath)
	if !utils.ErrCheck(err, "open file failed", filePath) {
		return nil, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(xmlFile))

	fileRaw := &ExcelFileRaw{
		Filename: filename,
		FieldRaw: treemap.NewWithStringComparator(),
		CellData: make([]ExcelRowData, 0),
	}
	parseXmlData(decoder, fileRaw)
	return fileRaw, nil
}

func getAllXmlFileNames(dirPath string) []string {
	dir, err := ioutil.ReadDir(dirPath)
	if !utils.ErrCheck(err, "read dir failed", dirPath) {
		return []string{}
	}

	// escape dir and ~$***.xlsx
	fileNames := make([]string, 0, len(dir))
	for _, fi := range dir {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".xml") && !strings.HasPrefix(fi.Name(), "~$") {
			fileNames = append(fileNames, fi.Name())
		}
	}

	return fileNames
}

// load all xml files
func loadAllXmlFiles(dirPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	mu := sync.Mutex{}
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException()
			rowDatas, err := loadOneXmlFile(dirPath, name)
			utils.ErrPrint(err, "loadOneXmlFile failed", name)

			mu.Lock()
			xmlFileRaws[name] = rowDatas
			mu.Unlock()
		})
	}
	wg.Wait()
}

// generate go code from xml file
func generateAllXmlCodes(dirPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException()
			err := generateCode(dirPath, xmlFileRaws[name])
			utils.ErrPrint(err, "generateCode failed", dirPath, name)
		})
	}

	wg.Wait()
}

func GenerateXml(dirPath string) {
	fileNames := getAllXmlFileNames(dirPath)
	loadAllXmlFiles(dirPath, fileNames)
	generateAllXmlCodes(dirPath, fileNames)
}

// read all xml entries
func ReadAllXmlEntries(dirPath string) {
	fileNames := getAllXmlFileNames(dirPath)
	loadAllXmlFiles(dirPath, fileNames)

	wg := utils.WaitGroupWrapper{}
	allXmlEntries.Range(func(k, v interface{}) bool {
		entryName := k.(string)
		entriesProto := v.(EntriesProto)

		wg.Wrap(func() {
			err := entriesProto.Load(xmlFileRaws[entryName])
			utils.ErrPrint(err, "gocode entry load failed", entryName)
		})

		return true
	})
	wg.Wait()

	log.Info().Msg("all xml entries reading completed!")
}

func parseXmlData(decoder *xml.Decoder, fileRaw *ExcelFileRaw) {
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}

		if !utils.ErrCheck(err, "get xml token failed") {
			break
		}

		if t, ok := token.(xml.StartElement); ok {
			if t.Name.Local == "root" {
				continue
			}

			if len(t.Attr) == 0 {
				continue
			}

			if t.Attr[0].Name.Local == "function" {
				// 字段注释
				if t.Attr[0].Value == "comment" {
					for n := 1; n < len(t.Attr); n++ {
						fieldName := t.Attr[n].Name.Local
						raw := &ExcelFieldRaw{
							name: fieldName,
							tag:  fmt.Sprintf("`json:\"%s,omitempty\" xml:\"%s,attr\"`", fieldName, fieldName),
							idx:  n - 1,
							desc: t.Attr[n].Value,
						}
						fileRaw.FieldRaw.Put(fieldName, raw)
					}
				}

				// 字段类型
				if t.Attr[0].Value == "type" {
					for n := 1; n < len(t.Attr); n++ {
						fieldName := t.Attr[n].Name.Local
						fieldValue := t.Attr[n].Value
						value, ok := fileRaw.FieldRaw.Get(fieldName)
						if !ok {
							log.Fatal().
								Str("filename", fileRaw.Filename).
								Str("fieldname", fieldName).
								Str("id", t.Attr[1].Value).
								Msg("parse xml data failed")
						}

						// import type: c->client, s->server
						needImport := true
						fieldValues := strings.Split(fieldValue, ":")
						if len(fieldValues) > 1 {
							needImport = false
							for k := 0; k < len(fieldValues)-1; k++ {
								if strings.Contains(fieldValues[k], "s") {
									needImport = true
								}
							}
						}

						value.(*ExcelFieldRaw).imp = needImport
						value.(*ExcelFieldRaw).tp = convertXmlType(fieldValues[len(fieldValues)-1])
						// typeValues[m-ColOffset] = fieldValues[len(fieldValues)-1]
					}
				}

				// 字段默认值
				if t.Attr[0].Value == "default" {
					for n := 1; n < len(t.Attr); n++ {
						fieldName := t.Attr[n].Name.Local
						fieldValue := t.Attr[n].Value
						value, ok := fileRaw.FieldRaw.Get(fieldName)
						if !ok {
							log.Fatal().
								Str("filename", fileRaw.Filename).
								Str("fieldname", fieldName).
								Str("id", t.Attr[1].Value).
								Msg("parse xml data failed")
						}

						value.(*ExcelFieldRaw).def = fieldValue
					}
				}

				continue
			}

			// xml element data
			mapRowData := make(map[string]interface{})
			for n := 0; n < len(t.Attr); n++ {
				fieldName := t.Attr[n].Name.Local
				fieldValue := t.Attr[n].Value
				value, ok := fileRaw.FieldRaw.Get(fieldName)
				if !ok {
					log.Fatal().
						Str("filename", fileRaw.Filename).
						Str("fieldname", fieldName).
						Str("id", t.Attr[1].Value).
						Msg("parse xml data failed")
				}

				xmlFieldRaw := value.(*ExcelFieldRaw)

				// set value
				var convertedVal interface{}
				if len(fieldValue) == 0 {
					if len(xmlFieldRaw.def) == 0 && xmlFieldRaw.tp != "string[]" && xmlFieldRaw.tp != "string" {
						log.Fatal().
							Str("filename", fileRaw.Filename).
							Str("fieldname", fieldName).
							Str("id", t.Attr[1].Value).
							Msg("default value not assigned")
					}
					convertedVal = convertXmlValue(xmlFieldRaw.tp, xmlFieldRaw.def)
				} else {
					convertedVal = convertXmlValue(xmlFieldRaw.tp, fieldValue)
				}
				mapRowData[fieldName] = convertedVal
			}

			fileRaw.CellData = append(fileRaw.CellData, mapRowData)
		}
	}
}

func convertXmlType(strType string) string {
	switch strType {
	case "int":
		return "int32"
	case "float":
		return "float32"
	case "int[]":
		return "[]int32"
	case "float[]":
		return "[]float32"
	case "string[]":
		return "[]string"
	case "manual":
		return "string"
	default:
		return strType
	}
}

func convertXmlValue(strType, strVal string) interface{} {
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
			arrVals[k] = convertXmlValue("int", v)
		}
		cellVal = arrVals

	case "float[]":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertXmlValue("float", v)
		}
		cellVal = arrVals

	case "string[]":
		cellVals := strings.Split(strVal, ",")
		arrVals := make([]interface{}, len(cellVals))
		for k, v := range cellVals {
			arrVals[k] = convertXmlValue("string", v)
		}
		cellVal = arrVals

	default:
		// default string value
		cellVal = strVal
	}

	return cellVal
}
