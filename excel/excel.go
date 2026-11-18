package excel

import (
	"container/list"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	map_utils "github.com/emirpasic/gods/utils"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/thanhpk/randstr"
	"github.com/xuri/excelize/v2"
)

var (
	RowOffset int = 2 // 第一行数据偏移
	ColOffset int = 2 // 第一列数据偏移
)

type KeyAliasData struct {
	Alias        string
	Filename     string
	RealKeyValue any
}

type KeyAliasMapping struct {
	mu                  sync.RWMutex
	conflictGlobalAlias map[string]bool
	globalAliasMappping map[string]*KeyAliasData
	fileAliasMapping    map[string]map[string]*KeyAliasData
}

func NewKeyAliasMapping() *KeyAliasMapping {
	return &KeyAliasMapping{
		conflictGlobalAlias: make(map[string]bool),
		globalAliasMappping: make(map[string]*KeyAliasData),
		fileAliasMapping:    make(map[string]map[string]*KeyAliasData),
	}
}

func (m *KeyAliasMapping) AddAlias(alias string, fileRaw *ExcelFileRaw, realKeyValue any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, found := m.globalAliasMappping[alias]; found {
		m.conflictGlobalAlias[alias] = true
		// log.Info().Str("filename", fileRaw.Filename).Str("sheetname", fileRaw.Sheetname).Str("alias", alias).Msg("有重复定义的全局别名，引用时请指定别名文件，否则使用此全局别名时会有冲突")
		// GlobalResultReport.AddWarnings(fileRaw.Filename, fileRaw.Sheetname, alias, fmt.Sprintf("%s有重复定义的全局别名%s，引用时请指定别名文件，否则使用此全局别名时会有冲突", data.Filename, data.Alias), "", "")
	}
	m.globalAliasMappping[alias] = &KeyAliasData{
		Alias:        alias,
		Filename:     fileRaw.ConfFilename,
		RealKeyValue: realKeyValue,
	}

	mapping, found := m.fileAliasMapping[fileRaw.ConfFilename]
	if !found {
		mapping = make(map[string]*KeyAliasData)
	}

	mapping[alias] = &KeyAliasData{
		Alias:        alias,
		Filename:     fileRaw.ConfFilename,
		RealKeyValue: realKeyValue,
	}

	m.fileAliasMapping[fileRaw.ConfFilename] = mapping
}

func (m *KeyAliasMapping) GetValue(aliasRefValue string, alias string) (*KeyAliasData, bool, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	aliasValue := strings.Split(aliasRefValue, "=")
	filenames := strings.Split(aliasValue[len(aliasValue)-1], ",")
	for _, filename := range filenames {
		if len(filename) == 0 || filename == "any" || filename == "global" {
			value, found := m.globalAliasMappping[alias]
			if m.conflictGlobalAlias[alias] {
				GlobalResultReport.AddError(value.Filename, "", "", fmt.Sprintf("引用了重复定义的全局别名%s", value.Alias))
			}
			return value, found, true
		}

		if mapping, found := m.fileAliasMapping[filename]; found {
			value, found := mapping[alias]
			if found {
				return value, found, false
			}
		}
	}

	return nil, false, false
}

type KeyRefMapping struct {
	mu             sync.RWMutex
	fileRefMapping map[string]map[string]any
}

func NewKeyRefMapping() *KeyRefMapping {
	return &KeyRefMapping{
		fileRefMapping: make(map[string]map[string]any),
	}
}

func (m *KeyRefMapping) AddRef(ref string, fileRaw *ExcelFileRaw, realKeyValue any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mapping, found := m.fileRefMapping[fileRaw.ConfFilename]
	if !found {
		mapping = make(map[string]any)
	}

	mapping[cast.ToString(realKeyValue)] = ref
	m.fileRefMapping[fileRaw.ConfFilename] = mapping
}

func (m *KeyRefMapping) GetValue(refFilenames string, ref string) (any, bool, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	refValue := strings.Split(refFilenames, "=")
	filenames := strings.Split(refValue[len(refValue)-1], ",")
	for _, filename := range filenames {
		if mapping, found := m.fileRefMapping[filename]; found {
			value, found := mapping[ref]
			if found {
				return value, found, false
			}
		}
	}

	return nil, false, false
}

var (
	entryArrayGetters    sync.Map
	dynamicEntryLoaders  sync.Map               // all dynamic entry loaders
	entryLoaders         sync.Map               // all entry loaders
	entryManualLoaders   sync.Map               // all entry manual loaders
	entryReflectTypes    sync.Map               // all entry reflect types
	excelFileRaws        sync.Map               // all excel file raw data string -> *ExcelFileRaw
	excelEnumTypeRaws    sync.Map               // all enum raw data: short_name -> enum raw data
	excelPolyTypeRaws    sync.Map               // all poly raw data: poly_name -> poly_data
	excelKeyAliasMapping = NewKeyAliasMapping() // key alias mapping
	excelKeyRefMapping   = NewKeyRefMapping()   // key ref mapping
)

type ExcelSheetType int

const (
	ExcelSheetType_Null ExcelSheetType = iota
	ExcelSheetType_Enum                // 枚举类型sheet
	ExcelSheetType_Poly                // 多态类型sheet
	ExcelSheetType_One                 // 单行数据sheet
	ExcelSheetType_Map                 // map数据sheet
)

type ExportLanguage_Type int

const (
	ExportLanguage_Go ExportLanguage_Type = iota
	ExportLanguage_GoDynamic
	ExportLanguage_Lua
	ExportLanguage_Ts
	ExportLanguage_Cs
	ExportLanguage_Diff
)

type ExcelRowData map[string]any

// enum data
type ExcelEnum struct {
	EnumName  string
	EnumDesc  string
	EnumValue string
}
type ExcelEnumType struct {
	EnumType       string
	EnumOriginType string
	Options        string
	Enums          []ExcelEnum
}

// poly data
type ExcelPoly struct {
	PolyName       string
	PolyOriginName string
	PolyType       string
	PolyTag        string
	PolySeperator  string
	PolyComment    string
	PolyAlias      string
}

type ExcelPolyType struct {
	PolyType        string
	PolyParentName  string
	PolySeperator   string
	PolyComment     string
	Options         string
	PolyUnique      string
	PolyLootKeyType string
	PolyFilename    string
	Polys           []ExcelPoly
	OriginPolys     []ExcelPoly
}

// Entry should implement Load function
type EntryLoader interface {
	Load(*ExcelFileRaw) (any, error)
}

// Dynamic Entry Loader
type DynamicEntryLoader interface {
	EntryLoader
	GetArrays() any
}

type EntryArrayGetter interface {
	GetArrays() any
}

// Entry should implement ManualLoad
type EntryManualLoader interface {
	ManualLoad(*ExcelFileRaw) error
}

// Excel field raw data
type ExcelFieldRaw struct {
	name       string
	originName string
	tp         string
	originTp   string
	desc       string
	tag        string
	idx        int  // field index in excel file
	export     bool // need export
	option     bool
}

type EnumOrderedValue struct {
	Name  string
	Value int
	Alias string
}

// Excel enum raw data
type ExcelEnumRaw struct {
	enumType       string
	enumOriginType string
	nameValue      map[string]int
	aliasValue     map[string]int
	orderedValues  []EnumOrderedValue
	sheetName      string
	confName       string
	options        string
}

// excel表格的json文件格式
type SheetsJsonFormat struct {
	SheetFile string     `json:"sheetFile"`
	ConfFile  string     `json:"confFile"`
	Bean      string     `json:"bean"`
	Mode      string     `json:"mode"`
	Options   string     `json:"options,omitempty"`
	Values    [][]string `json:"values"`
}

func normalizeSheetOptions(options string) string {
	items := strings.Split(options, ",")
	result := make([]string, 0, len(items))
	exists := make(map[string]struct{}, len(items))
	for _, item := range items {
		option := strings.TrimSpace(strings.ToLower(item))
		if len(option) == 0 {
			continue
		}
		if _, found := exists[option]; found {
			continue
		}
		exists[option] = struct{}{}
		result = append(result, option)
	}
	return strings.Join(result, ",")
}

func hasSheetOption(options, target string) bool {
	target = strings.TrimSpace(strings.ToLower(target))
	if len(target) == 0 {
		return false
	}
	for _, item := range strings.Split(options, ",") {
		if strings.TrimSpace(strings.ToLower(item)) == target {
			return true
		}
	}
	return false
}

func NewExcelEnumRaw(enumType, originType, sheetName, confName, options string) *ExcelEnumRaw {
	return &ExcelEnumRaw{
		enumType:       enumType,
		enumOriginType: originType,
		nameValue:      make(map[string]int),
		aliasValue:     make(map[string]int),
		orderedValues:  make([]EnumOrderedValue, 0),
		sheetName:      sheetName,
		confName:       confName,
		options:        normalizeSheetOptions(options),
	}
}

func (r *ExcelEnumRaw) HasOption(option string) bool {
	return hasSheetOption(r.options, option)
}

func (p *ExcelPolyType) HasOption(option string) bool {
	return hasSheetOption(p.Options, option)
}

func (r *ExcelEnumRaw) AddEnum(name string, alias string, value int) {
	r.nameValue[name] = value
	r.aliasValue[alias] = value
	r.orderedValues = append(r.orderedValues, EnumOrderedValue{Name: name, Value: value, Alias: alias})
}

func (r *ExcelEnumRaw) GetEnumValue(name string) (int, bool) {
	if val, ok := r.nameValue[name]; ok {
		return val, ok
	}

	if val, ok := r.aliasValue[name]; ok {
		return val, ok
	}

	// number type
	val, err := cast.ToIntE(name)
	if err != nil {
		return -1, false
	}

	return val, true
}

func (r *ExcelEnumRaw) GetEnumAlias(value int) (string, bool) {
	for name, v := range r.aliasValue {
		if v == value {
			return name, true
		}
	}

	return "", false
}

// Excel file raw data
type ExcelFileRaw struct {
	Type         ExcelSheetType      // sheet type
	LanType      ExportLanguage_Type // 导出语言类型
	ShouldExport bool                // 是否需要导出
	Filename     string
	Sheetname    string
	Beanname     string
	ConfFilename string
	Keys         []string
	HasMap       bool
	FieldRaw     *treemap.Map     // map[string]*ExcelFieldRaw
	EnumData     []*ExcelEnumType // enum data
	PolyData     []*ExcelPolyType // poly data
	RowData      []ExcelRowData   // sheet map data
}

func init() {
}

func AddDynamicEntryLoader(name string, e DynamicEntryLoader) {
	dynamicEntryLoaders.Store(name, e)
}

func RangeDynamicEntryLoaders(fn func(k, v any) bool) {
	dynamicEntryLoaders.Range(fn)
}

func AddEntryLoader(name string, e EntryLoader) {
	entryLoaders.Store(name, e)
}

func AddArrayGetter(name string, e EntryArrayGetter) {
	entryArrayGetters.Store(name, e)
}

func RangeArrayGetters(fn func(k, v any) bool) {
	entryArrayGetters.Range(fn)
}

func GetArrayGetter(name string) (EntryArrayGetter, bool) {
	value, found := entryArrayGetters.Load(name)
	if !found {
		return nil, false
	}

	return value.(EntryArrayGetter), true
}

func AddEntryManualLoader(name string, e EntryManualLoader) {
	entryManualLoaders.Store(name, e)
}

func AddEntryNewFunc(name string, fn func() any) {
	entryReflectTypes.Store(name, fn)
}

func GetEntryNewFunc(name string) func() any {
	v, ok := entryReflectTypes.Load(name)
	if !ok {
		return nil
	}

	return v.(func() any)
}

func parseCSVRows(filename string, rows [][]string) (*ExcelFileRaw, error) {

	fileRaw := &ExcelFileRaw{
		Type:     ExcelSheetType_Null,
		Filename: filename,
		FieldRaw: treemap.NewWithStringComparator(),
		EnumData: make([]*ExcelEnumType, 0),
		PolyData: make([]*ExcelPolyType, 0),
		RowData:  make([]ExcelRowData, 0),
		Keys:     make([]string, 0),
	}

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
	dir, err := os.ReadDir(readExcelPath)
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

func getAllJsonFileNames(readJsonPath string) []string {
	fileNames := make([]string, 0, 32)
	_ = filepath.WalkDir(readJsonPath, func(path string, d fs.DirEntry, err error) error {
		// 目录不存在或不可读时跳过，避免 nil 解引用
		if err != nil || d == nil {
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			fileNames = append(fileNames, d.Name())
		}

		return nil
	})

	return fileNames
}

func getAllCsvFileNames(readExcelPath string) []string {
	fileNames := make([]string, 0, 32)
	_ = filepath.WalkDir(readExcelPath, func(path string, d fs.DirEntry, err error) error {
		// escape dir and ~$***.csv
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".csv") && !strings.HasPrefix(d.Name(), "~$") {
			fileNames = append(fileNames, d.Name())
		}
		return nil
	})

	return fileNames
}

// load all excel files and translate to CSV
func loadExcelToCSV(dirPath string, exportCsvPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException(name)

			filePath := fmt.Sprintf("%s%s", dirPath, name)
			xlsxFile, err := excelize.OpenFile(filePath)
			if !utils.ErrCheck(err, "excelize.OpenFile failed", filePath) {
				return
			}

			// read rows from excel files
			rows, err := xlsxFile.GetRows(xlsxFile.GetSheetName(0))
			if !utils.ErrCheck(err, "xlsxFile.GetRows failed", filePath) {
				return
			}

			// parse rows to *ExcelFileRaw
			csvName := strings.Replace(name, ".xlsx", ".csv", -1)
			rowDatas, err := parseCSVRows(csvName, rows)
			if !utils.ErrCheck(err, "parseRows failed", name) {
				return
			}

			// write rows to csv files
			fiPath := fmt.Sprintf("%s%s", exportCsvPath, csvName)
			fi, err := os.OpenFile(fiPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if !utils.ErrCheck(err, "OpenFile failed", fiPath) {
				return
			}

			w := csv.NewWriter(fi)
			err = w.WriteAll(rows)
			if !utils.ErrCheck(err, "csv.WriteAll failed", fiPath) {
				return
			}

			excelFileRaws.Store(csvName, rowDatas)
		})
	}
	wg.Wait()
}

// load all excel files and translate to json
func LoadExcelToJson(dirPath string, exportJsonPath string) {
	wg := utils.WaitGroupWrapper{}

	_ = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		name := d.Name()
		if d.IsDir() || !strings.HasSuffix(name, ".xlsx") {
			return nil
		}

		filePath := filepath.Join(dirPath, name)
		xlsxFile, errOpen := excelize.OpenFile(filePath)
		if !utils.ErrCheck(errOpen, "excelize.OpenFile failed", filePath) {
			return nil
		}

		for _, sheetName := range xlsxFile.GetSheetList() {
			// read rows from excel files
			rows, err := xlsxFile.GetRows(sheetName, excelize.Options{RawCellValue: true})
			if !utils.ErrCheck(err, "xlsxFile.GetRows failed", filePath) {
				return nil
			}

			title, bean, mode, options := parseSheetMeta(sheetName)

			sheetJsonFormat := &SheetsJsonFormat{
				SheetFile: sheetName,
				ConfFile:  title,
				Bean:      bean,
				Mode:      mode,
				Options:   options,
				Values:    rows,
			}

			wg.Wrap(func() {
				defer utils.CaptureException(name)

				jsonData, err := json.MarshalIndent(sheetJsonFormat, "", "\t")
				if !utils.ErrCheck(err, "marshal sheetJsonFormat failed") {
					return
				}

				filePath := filepath.Join(exportJsonPath, title+".json")
				err = os.WriteFile(filePath, jsonData, 0666)
				if !utils.ErrCheck(err, "write file failed") {
					return
				}
			})
		}

		return nil
	})

	wg.Wait()
}

func readCSV(dirPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := v
		wg.Wrap(func() {
			defer utils.CaptureException(name)
			fiPath := filepath.Join(dirPath, name)
			fi, err := os.OpenFile(fiPath, os.O_RDONLY, 0666)
			if !utils.ErrCheck(err, "OpenFile failed", fiPath) {
				return
			}

			defer fi.Close()

			r := csv.NewReader(fi)
			r.FieldsPerRecord = -1
			rows, err := r.ReadAll()
			if !utils.ErrCheck(err, "ReadAll failed", fiPath) {
				return
			}

			rowDatas, err := parseCSVRows(name, rows)
			if !utils.ErrCheck(err, "parseRows failed", name) {
				return
			}

			excelFileRaws.Store(name, rowDatas)
		})
	}
	wg.Wait()
}

func parseSheetMeta(sheetName string) (title string, bean string, mode string, options string) {
	title = sheetName
	bean = sheetName
	mode = "map"

	parts := strings.Split(sheetName, "|")
	if len(parts) == 0 {
		return
	}

	title = parts[0]
	nonOptionParts := make([]string, 0, len(parts)-1)
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			continue
		}
		lowerPart := strings.ToLower(part)
		if strings.HasPrefix(lowerPart, "options=") {
			options = normalizeSheetOptions(strings.TrimSpace(part[len("options="):]))
			continue
		}
		nonOptionParts = append(nonOptionParts, part)
	}

	if len(nonOptionParts) > 0 {
		lastPart := strings.ToLower(nonOptionParts[len(nonOptionParts)-1])
		if lastPart == "one" {
			mode = "one"
			nonOptionParts = nonOptionParts[:len(nonOptionParts)-1]
		} else if lastPart == "map" {
			mode = "map"
			nonOptionParts = nonOptionParts[:len(nonOptionParts)-1]
		}
	}

	if mode == "one" {
		bean = title
	} else if len(nonOptionParts) > 0 {
		bean = nonOptionParts[len(nonOptionParts)-1]
	}

	return
}

// generate go code from excel file
func generateAllGoCodes(exportPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := strings.Replace(v, ".xlsx", ".csv", -1)
		wg.Wrap(func() {
			defer utils.CaptureException(name)
			value, found := excelFileRaws.Load(name)
			if !found {
				log.Warn().Str("name", name).Msg("cannot find excel file raw")
				return
			}

			excelFileRaw := value.(*ExcelFileRaw)
			if excelFileRaw.Type == ExcelSheetType_Enum {
				err := generateEnumGoCode("config", exportPath, excelFileRaw)
				if !utils.ErrCheck(err, "generateEnumGoCode failed", exportPath, name) {
					return
				}
			} else if excelFileRaw.Type == ExcelSheetType_Poly {
				err := generatePolyGoCode("config", exportPath, excelFileRaw)
				if !utils.ErrCheck(err, "generatePolyGoCode failed", exportPath, name) {
					return
				}
			} else {
				err := generateEntryGoCode("config", exportPath, excelFileRaw)
				if !utils.ErrCheck(err, "generateEntryGoCode failed", exportPath, name) {
					return
				}
			}

			log.Info().Str("fileName", name).Str("exportDir", exportPath).Caller().Msg("generate go code success")
		})
	}

	wg.Wait()

}

// generate go dynamic code from excel file
func generateAllGoDynamicCodes(exportPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := strings.Replace(v, ".xlsx", ".csv", -1)
		wg.Wrap(func() {
			defer utils.CaptureException(name)
			value, found := excelFileRaws.Load(name)
			if !found {
				log.Warn().Str("name", name).Msg("cannot find excel file raw")
				return
			}

			excelFileRaw := value.(*ExcelFileRaw)
			if excelFileRaw.Type == ExcelSheetType_Enum {
				err := generateEnumGoCode("dynamic", exportPath, excelFileRaw)
				if !utils.ErrCheck(err, "generateEnumGoCode failed", exportPath, name) {
					return
				}
			} else if excelFileRaw.Type == ExcelSheetType_Poly {
				err := generatePolyGoCode("dynamic", exportPath, excelFileRaw)
				if !utils.ErrCheck(err, "generatePolyGoCode failed", exportPath, name) {
					return
				}
			} else {
				err := generateEntryGoCode("dynamic", exportPath, excelFileRaw)
				if !utils.ErrCheck(err, "generateEntryGoCode failed", exportPath, name) {
					return
				}
			}

			log.Info().Str("fileName", name).Str("exportDir", exportPath).Caller().Msg("generate go code success")
		})
	}

	wg.Wait()

}

// generate lua code from excel file
func generateAllLuaCodes(exportPath string, fileNames []string) {
	wg := utils.WaitGroupWrapper{}
	for _, v := range fileNames {
		name := strings.Replace(v, ".xlsx", ".csv", -1)
		wg.Wrap(func() {
			value, found := excelFileRaws.Load(name)
			if !found {
				log.Warn().Str("name", name).Msg("cannot find excel file raw")
				return
			}

			excelFileRaw := value.(*ExcelFileRaw)
			if excelFileRaw.Type == ExcelSheetType_Enum {
				// err := generateEnumLuaCode("config", exportPath, excelFileRaw)
				// if !utils.ErrCheck(err, "generateEnumLuaCode failed", exportPath, name) {
				// 	return
				// }
			} else if excelFileRaw.Type == ExcelSheetType_Poly {
				// err := generatePolyGoCode("config", exportPath, excelFileRaw)
				// if !utils.ErrCheck(err, "generatePolyLuaCode failed", exportPath, name) {
				// 	return
				// }
			} else {
				err := generateEntryLuaCode("config", filepath.Join(exportPath, "Data"), excelFileRaw)
				if !utils.ErrCheck(err, "generateEntryLuaCode failed", exportPath, name) {
					return
				}
			}

			log.Info().Msgf("读取%s，生成lua代码到%s", name, exportPath)
			// log.Info().Str("file_name", name).Str("export_dir", exportPath).Caller().Msg("generate lua code success")
		})
	}

	// all required files
	requiredFilenames := make([]string, 0, len(fileNames))
	for _, v := range fileNames {
		name := strings.Replace(v, ".xlsx", ".csv", -1)
		value, found := excelFileRaws.Load(name)
		if !found {
			log.Warn().Str("name", name).Msg("cannot find excel file raw")
			return
		}

		excelFileRaw := value.(*ExcelFileRaw)
		if excelFileRaw.Type == ExcelSheetType_Enum || excelFileRaw.Type == ExcelSheetType_Poly {
			continue
		}

		requiredFilenames = append(requiredFilenames, name)
	}

	err := generateRequiredLuaCode(exportPath, requiredFilenames)
	if !utils.ErrCheck(err, "generateRequiredLuaCode failed", exportPath) {
		return
	}

	// all constants files
	err = generateEnumLuaCode(exportPath)
	if !utils.ErrCheck(err, "generateEnumLuaCode failed", exportPath) {
		return
	}

	wg.Wait()

}

// generate ts code from excel file
func generateAllTsCodes(exportPath string, fileNames []string) {

	// generate data.json
	var exportTsJsonPath string = "json"
	if f := flag.Lookup("exportTsJsonPath"); f != nil {
		exportTsJsonPath = f.Value.String()
	}
	_ = os.MkdirAll(exportTsJsonPath, os.ModePerm)
	err := generateEntryTsCode("config", exportTsJsonPath)
	if !utils.ErrCheck(err, "generateEntryTsCode failed", exportPath) {
		return
	}

	// all required files
	requiredFilenames := make([]string, 0, len(fileNames))

	for _, v := range fileNames {
		name := strings.Replace(v, ".xlsx", ".csv", -1)
		value, found := excelFileRaws.Load(name)
		if !found {
			log.Warn().Str("name", name).Msg("cannot find excel file raw")
			return
		}

		excelFileRaw := value.(*ExcelFileRaw)
		if excelFileRaw.Type == ExcelSheetType_Enum || excelFileRaw.Type == ExcelSheetType_Poly {
			continue
		}

		requiredFilenames = append(requiredFilenames, name)
	}

	// generate enum and require code
	_ = os.MkdirAll(exportPath, os.ModePerm)
	err = generateRequiredTsCode(exportPath, requiredFilenames)
	if !utils.ErrCheck(err, "generateRequiredTsCode failed", exportPath) {
		return
	}
}

func GenerateFromExcel(readExcelPath, exportGoPath, exportCsvPath string) {
	fileNames := getAllExcelFileNames(readExcelPath)
	loadExcelToCSV(readExcelPath, exportCsvPath, fileNames)
	generateAllGoCodes(exportGoPath, fileNames)
}

// read all excel entries
// LoadAll 依次加载 CSV 与 JSON 两套配置表。
//
// 迁移期间的约定：
//   - 尚未迁移的老师表仍在 config/csv，由 ReadAllCSVEntries 读取；
//   - 从 kof 移植的系统使用 config/json，由 ReadAllJsonEntries 读取。
//
// 未命中文件的 entry loader 会被自动跳过，因此两套并存是安全的。
func LoadAll(csvDir, jsonDir string) {
	ReadAllCSVEntries(csvDir)
	ReadAllJsonEntries(jsonDir)
}

func ReadAllCSVEntries(dirPath string) {
	fileNames := getAllCsvFileNames(dirPath)
	readCSV(dirPath, fileNames)

	wg := utils.WaitGroupWrapper{}

	// read from excel files
	entryLoaders.Range(func(k, v any) bool {
		entryName := k.(string)
		loader := v.(EntryLoader)

		wg.Wrap(func() {
			defer utils.CaptureException(entryName)
			if value, found := excelFileRaws.Load(entryName); found {
				excelFileRaw := value.(*ExcelFileRaw)
				_, err := loader.Load(excelFileRaw)
				utils.WarnPrint(err, "EntryLoader Load failed", entryName)
			} else {
				log.Error().Str("entryName", entryName).Msg("cannot find csv file")
			}
		})

		return true
	})
	wg.Wait()

	// load by manual
	entryManualLoaders.Range(func(k, v any) bool {
		entryName := k.(string)
		loader := v.(EntryManualLoader)

		wg.Wrap(func() {
			defer utils.CaptureException(entryName)
			if value, found := excelFileRaws.Load(entryName); found {
				excelFileRaw := value.(*ExcelFileRaw)
				err := loader.ManualLoad(excelFileRaw)
				utils.WarnPrint(err, "EntryManualLoader Load failed", entryName)
			}
		})

		return true
	})
	wg.Wait()

	log.Info().Msg("all excel entries reading completed!")
}

// read all json entries
func ReadAllJsonEntries(dirPath string) {
	{
		// path := filepath.Join("config", "csv")
		// fileNames := getAllCsvFileNames(path)
		// readCSV(path, fileNames)
	}

	{
		fileNames := getAllJsonFileNames(dirPath)
		loadJson(ExportLanguage_Go, dirPath, fileNames)
	}

	wg := utils.WaitGroupWrapper{}

	// read from excel files
	entryLoaders.Range(func(k, v any) bool {
		entryName := k.(string)
		loader := v.(EntryLoader)
		wg.Wrap(func() {
			// defer utils.CaptureException(entryName)
			if value, found := excelFileRaws.Load(entryName); found {
				excelFileRaw := value.(*ExcelFileRaw)
				_, err := loader.Load(excelFileRaw)
				utils.WarnPrint(err, "EntryLoader Load failed", entryName)
			}
		})

		return true
	})
	wg.Wait()

	// load by manual
	entryManualLoaders.Range(func(k, v any) bool {
		entryName := k.(string)
		loader := v.(EntryManualLoader)

		wg.Wrap(func() {
			if value, found := excelFileRaws.Load(entryName); found {
				excelFileRaw := value.(*ExcelFileRaw)
				err := loader.ManualLoad(excelFileRaw)
				utils.WarnPrint(err, "EntryManualLoader Load failed", entryName)
			}
		})

		return true
	})
	wg.Wait()

	// load dynamic entries
	dynamicEntryLoaders.Range(func(k, v any) bool {
		entryName := k.(string)
		loader := v.(DynamicEntryLoader)
		wg.Wrap(func() {
			// defer utils.CaptureException(entryName)
			if value, found := excelFileRaws.Load(entryName); found {
				excelFileRaw := value.(*ExcelFileRaw)
				row, err := loader.Load(excelFileRaw)
				if !utils.ErrCheck(err, "DynamicEntryLoader Load failed", entryName) {
					return
				}

				AddDynamicEntryLoader(entryName, row.(DynamicEntryLoader))
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
					raw.export = false
					fieldName = fmt.Sprintf("V%s", randstr.String(16))
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

				tagName := utils.SimpleUntitle(raw.name)
				raw.tag = fmt.Sprintf("`json:\"%s,omitempty\" bson:\"%s\"`", tagName, tagName)
				fileRaw.FieldRaw.Put(raw.name, raw)
				typeNames[m-ColOffset] = raw.name
			}
		}

		// load type desc
		// if n == RowOffset+1 {

		// }

		// load type control
		if n == RowOffset+2 {
			var strBuilder strings.Builder
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

					strBuilder.Reset()
					strBuilder.WriteString(desc)
					strBuilder.WriteString(" 主键")
					value.(*ExcelFieldRaw).export = true
					value.(*ExcelFieldRaw).desc = strBuilder.String()
					continue
				}

				// 带K标识的也是主键
				if strings.Contains(fieldValue, "K") {
					fileRaw.Keys = append(fileRaw.Keys, value.(*ExcelFieldRaw).name)

					// 去除换行
					desc := strings.Replace(value.(*ExcelFieldRaw).desc, "\n", ",", -1)
					strBuilder.Reset()
					strBuilder.WriteString(desc)
					strBuilder.WriteString(" 多主键之一")
					value.(*ExcelFieldRaw).desc = strBuilder.String()
				} else {
					// 去除换行
					desc := strings.Replace(rows[n-1][m], "\n", ",", -1)
					value.(*ExcelFieldRaw).desc = desc
				}

				if strings.Contains(fieldValue, "C") {
					value.(*ExcelFieldRaw).export = false
				} else {
					value.(*ExcelFieldRaw).export = true
				}
			}
		}

		// load type value
		if n == RowOffset+3 {
			for m := ColOffset; m < len(rows[n]); m++ {
				fieldName := typeNames[m-ColOffset]
				fieldValue := rows[n][m]
				convertType, _ := convertType(fieldValue)

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
					value.(*ExcelFieldRaw).export = false
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

		mapRowData := make(map[string]any)
		for m := ColOffset; m < len(rows[n]); m++ {
			cellColIdx := m - ColOffset
			cellValString := rows[n][m]
			fieldName := typeNames[m-ColOffset]

			// set value
			convertedVal, err := convertValue(fileRaw.Filename, fileRaw.Sheetname, fieldName, typeValues[cellColIdx], cellValString)
			if err != nil {
				log.Warn().Str("err", fmt.Sprintf("parse field failed, err:%+v, type:%+v, val:%+v", err, typeValues[cellColIdx], cellColIdx))
				GlobalResultReport.AddWarnings(fileRaw.Filename, fileRaw.Sheetname, fieldName, err.Error(), rows[m][0], cellValString)
			}
			mapRowData[typeNames[cellColIdx]] = convertedVal
		}

		fileRaw.RowData = append(fileRaw.RowData, mapRowData)
	}
}

// be tolerant with type names
func convertType(strType string) (string, bool) {
	switch strType {
	case "String", "STRING", "string":
		return "string", true

	case "[]String", "String[]", "[]STRING", "[]string":
		return "[]string", true

	case "Int32", "int32":
		return "int32", true

	case "Int", "INT", "int":
		return "int32", true

	// 注意：east-eden 早期导出的 CSV 用 "number" 表示十进制小数，
	// 而 kof 的表用 int/float 表达数值，不使用 "number"。
	// 这里保持 east-eden 的语义，避免既有配置表被当成百分比整数。
	case "Number", "NUMBER", "number", "decimal.Decimal":
		return "decimal.Decimal", true

	case "Float32", "Float", "FLOAT", "float", "float32":
		return "float32", true

	case "Float64", "Double", "double", "DOUBLE", "float64":
		return "float64", true

	case "[]Int32", "[]int32":
		return "[]int32", true

	case "[]Int", "[]INT", "[]int":
		return "[]int32", true

	case "[]Number", "[]NUMBER", "Number[]", "NUMBER[]", "number[]", "[]number", "[]decimal.Decimal":
		return "[]decimal.Decimal", true

	case "bool", "BOOL", "Bool", "boolean":
		return "bool", true

	default:
		if strings.HasPrefix(strType, "map") || strings.HasPrefix(strType, "Map") {
			ts := strings.Split(strType, "[")
			t := ts[len(ts)-1]
			tt := strings.Split(t, "]")
			keyType, _ := convertType(tt[0])
			valueType, _ := getFullTypeName(tt[1])

			return fmt.Sprintf("map[%s]%s", keyType, valueType), true
		}

		if strings.HasPrefix(strType, "list2") || strings.HasPrefix(strType, "List2") {
			types := strings.Split(strType, ":")
			tp, _ := getFullTypeName(types[1])
			return fmt.Sprintf("[][]%s", tp), true

		} else if strings.HasPrefix(strType, "list") || strings.HasPrefix(strType, "List") {
			types := strings.Split(strType, ":")
			tp, _ := getFullTypeName(types[1])
			return fmt.Sprintf("[]%s", tp), true
		}

		// not buildin type
		return strType, false
	}
}

func convertValue(filename, sheetname, keyname string, strType, strVal string, param ...string) (any, error) {
	// 分隔符
	seperator := ","

	// 主键别名
	aliasVal := ""

	if len(param) > 0 && len(param[0]) > 0 {
		seperator = param[0]
	}

	if len(param) > 1 {
		aliasVal = param[1]
	}

	var cellVal any
	var err error
	convertType, _ := getFullTypeName(strType)
	strVal = strings.TrimSpace(strVal)

	switch convertType {
	case "string":
		if len(aliasVal) > 0 {
			if value, found, _ := excelKeyAliasMapping.GetValue(aliasVal, strVal); found {
				strVal = cast.ToString(value.RealKeyValue)
			}
		}

		cellVal = strVal
	case "int64":
		if len(aliasVal) > 0 {
			if value, found, _ := excelKeyAliasMapping.GetValue(aliasVal, strVal); found {
				strVal = cast.ToString(value.RealKeyValue)
			}
		}

		if len(strVal) == 0 {
			cellVal = int64(0)
		} else {
			cellVal = cast.ToInt64(strVal)
		}
	case "int32":
		if len(aliasVal) > 0 {
			if value, found, _ := excelKeyAliasMapping.GetValue(aliasVal, strVal); found {
				strVal = cast.ToString(value.RealKeyValue)
			}
		}

		if len(strVal) == 0 {
			cellVal = int32(0)
		} else {
			cellVal = cast.ToInt32(strVal)
		}

	case "int":
		if len(aliasVal) > 0 {
			if value, found, _ := excelKeyAliasMapping.GetValue(aliasVal, strVal); found {
				strVal = cast.ToString(value.RealKeyValue)
			}
		}

		if len(strVal) == 0 {
			cellVal = int(0)
		} else {
			cellVal, err = cast.ToIntE(strVal)
			if err != nil {
				return 0, err
			}
		}

	case "decimal.Decimal":
		if len(strVal) == 0 || strVal == "0" {
			cellVal = decimal.NewFromInt32(0)
		} else {
			cellVal, err = decimal.NewFromString(strVal)
			if err != nil {
				return 0, err
			}
			// floatVal := cast.ToFloat64(strVal)
			// floatVal *= define.PercentBase
			// floatVal = math.Round(floatVal)
			// cellVal = int32(floatVal)
		}

	case "number":
		if len(aliasVal) > 0 {
			if value, found, _ := excelKeyAliasMapping.GetValue(aliasVal, strVal); found {
				strVal = cast.ToString(value.RealKeyValue)
			} else {
				err = fmt.Errorf("使用了未定义的主键别名：%s", strVal)
			}
		}

		if len(strVal) == 0 || strVal == "0" {
			cellVal = int32(0)
		} else {
			floatVal, err := strconv.ParseFloat(strVal, 32)
			utils.WarnPrint(err, "convert cell value to number failed", strVal)
			if err != nil {
				return 0, err
			}

			floatVal *= define.PercentBase
			floatVal = math.Round(floatVal)
			cellVal = int32(floatVal)
		}

	case "float32":
		if len(strVal) == 0 {
			cellVal = float32(0)
		} else {
			cellVal, err = cast.ToFloat32E(strVal)
			if err != nil {
				return 0, err
			}
		}

	case "float64":
		if len(strVal) == 0 {
			cellVal = float64(0)
		} else {
			cellVal, err = cast.ToFloat64E(strVal)
			if err != nil {
				return 0, err
			}
		}

	case "[]int64":
		cellVals := strings.Split(strVal, seperator)
		arrVals := make([]any, 0, len(cellVals))
		if len(strVal) > 0 {
			for _, v := range cellVals {
				var intVal any
				intVal, err = convertValue(filename, sheetname, keyname, "int64", v, "", aliasVal)
				if err != nil {
					return nil, err
				}
				arrVals = append(arrVals, intVal)
			}
		}
		cellVal = arrVals

	case "[]int32":
		cellVals := strings.Split(strVal, seperator)
		arrVals := make([]any, 0, len(cellVals))
		if len(strVal) > 0 {
			for _, v := range cellVals {
				var intVal any
				intVal, err = convertValue(filename, sheetname, keyname, "int32", v, "", aliasVal)
				if err != nil {
					return nil, err
				}
				arrVals = append(arrVals, intVal)
			}
		}
		cellVal = arrVals

	case "[]int":
		cellVals := strings.Split(strVal, seperator)
		arrVals := make([]any, 0, len(cellVals))
		if len(strVal) > 0 {
			for _, v := range cellVals {
				var intVal any
				intVal, err = convertValue(filename, sheetname, keyname, "int", v, "", aliasVal)
				if err != nil {
					return nil, err
				}
				arrVals = append(arrVals, intVal)
			}
		}
		cellVal = arrVals

	case "[]decimal.Decimal":
		cellVals := strings.Split(strVal, seperator)
		arrVals := make([]any, len(cellVals))
		for k, v := range cellVals {
			var decimalVal any
			decimalVal, err = convertValue(filename, sheetname, keyname, "number", v, "", aliasVal)
			if err != nil {
				return nil, err
			}
			arrVals[k] = decimalVal
		}
		cellVal = arrVals

	case "[]float32":
		cellVals := strings.Split(strVal, seperator)
		arrVals := make([]any, 0, len(cellVals))
		if len(strVal) > 0 {
			for _, v := range cellVals {
				var floatVal any
				floatVal, err = convertValue(filename, sheetname, keyname, "float32", v)
				if err != nil {
					return nil, err
				}
				arrVals = append(arrVals, floatVal)
			}
		}
		cellVal = arrVals

	case "[]float64":
		cellVals := strings.Split(strVal, seperator)
		arrVals := make([]any, 0, len(cellVals))
		if len(strVal) > 0 {
			for _, v := range cellVals {
				var floatVal any
				floatVal, err = convertValue(filename, sheetname, keyname, "float64", v)
				if err != nil {
					return nil, err
				}
				arrVals = append(arrVals, floatVal)
			}
		}
		cellVal = arrVals

	case "[]string":
		cellVals := strings.Split(strVal, seperator)
		arrVals := make([]any, 0, len(cellVals))
		if len(strVal) > 0 {
			for _, v := range cellVals {
				var stringVal any
				stringVal, err = convertValue(filename, sheetname, keyname, "string", v, "", aliasVal)
				if err != nil {
					return nil, err
				}
				arrVals = append(arrVals, stringVal)
			}
		}

		cellVal = arrVals

	case "*treemap.Map":
		cellVal, err = convertMapValue(filename, sheetname, keyname, strType, strVal, seperator)
		if err != nil {
			return nil, err
		}

	case "bool":
		if len(strVal) <= 0 {
			return false, nil
		}
		cellVal, err = cast.ToBoolE(strVal)
		if err != nil {
			return false, err
		}

	default:
		if strings.HasPrefix(convertType, "map[") {
			cellVal, err = convertBuildinMapValue(filename, sheetname, keyname, strType, strVal, seperator)
			if err != nil {
				return nil, err
			}
			break
		}

		// [][]any
		if strings.HasPrefix(convertType, "[][]*") {
			nextType, _ := strings.CutPrefix(convertType, "[]")
			cellVals := strings.Split(strVal, seperator)
			arrVals := make([]any, 0, len(cellVals))
			if len(strVal) > 0 {
				for _, v := range cellVals {
					var sliceVal any
					sliceVal, err = convertValue(filename, sheetname, keyname, nextType, v)
					if err != nil {
						return nil, err
					}
					arrVals = append(arrVals, sliceVal)
				}
			}
			cellVal = arrVals
			break
		}

		// list
		if strings.HasPrefix(convertType, "list") {
			types := strings.Split(convertType, ":")
			ls := list.New()
			if len(strVal) > 0 {
				values := strings.Split(strVal, ";")
				for _, v := range values {
					var typeVal any
					typeVal, err = convertValue(filename, sheetname, keyname, types[1], v)
					if err != nil {
						return nil, err
					}
					ls.PushBack(typeVal)
				}
			}

			cellVal = ls
			break
		}

		// array poly
		if strings.HasPrefix(convertType, "[]*") {
			types := strings.Split(convertType, "*")
			seperator := ";"
			if len(param) > 0 && len(param[0]) > 0 {
				seperator = param[0]
			}
			if value, found := excelPolyTypeRaws.Load(types[len(types)-1]); found {
				raw := value.(*ExcelPolyType)
				values := strings.Split(strVal, seperator)
				valueArr := make([]any, 0, len(values))
				if len(strVal) > 0 {
					for _, value := range values {
						polyType := raw.PolyType
						if strings.HasPrefix(convertType, "[]*") {
							polyType = fmt.Sprintf("*%s", polyType)
						}
						typeVal, err := convertValue(filename, sheetname, keyname, polyType, value)
						if err != nil {
							return nil, err
						}
						valueArr = append(valueArr, typeVal)
					}
				}

				cellVal = valueArr
				break
			}
		}

		// pointer poly
		if strings.HasPrefix(convertType, "*") {
			var err error
			tp := strings.TrimPrefix(convertType, "*")
			if value, found := excelPolyTypeRaws.Load(tp); found {
				raw := value.(*ExcelPolyType)
				values := strings.Split(strVal, raw.PolySeperator)
				mapVal := make(map[string]any)

				// default value
				for _, poly := range raw.Polys {
					mapVal[poly.PolyName], err = convertValue(filename, sheetname, keyname, poly.PolyType, "", poly.PolySeperator)
				}

				for idx, value := range values {
					var polyValue string
					tagName, splitValue, ok := splitPolyValue(value)
					if ok {
						// 指定tagname
						for _, poly := range raw.Polys {
							if poly.PolyOriginName != tagName {
								continue
							}

							polyValue = splitValue

							// 主键别名引用
							if len(poly.PolyAlias) > 0 {
								if value, found, _ := excelKeyAliasMapping.GetValue(poly.PolyAlias, polyValue); found {
									polyValue = cast.ToString(value.RealKeyValue)
								}
							}
							mapVal[poly.PolyName], err = convertValue(filename, sheetname, keyname, poly.PolyType, polyValue, poly.PolySeperator, poly.PolyAlias)
							break
						}
					} else {
						// 非指定tagname，按字段顺序
						polyValue = value
						if idx < len(raw.Polys) {
							// 主键别名引用
							if len(raw.Polys[idx].PolyAlias) > 0 {
								if value, found, _ := excelKeyAliasMapping.GetValue(raw.Polys[idx].PolyAlias, polyValue); found {
									polyValue = cast.ToString(value.RealKeyValue)
								}
							}
							mapVal[raw.Polys[idx].PolyName], err = convertValue(filename, sheetname, keyname, raw.Polys[idx].PolyType, polyValue, raw.Polys[idx].PolySeperator, raw.Polys[idx].PolyAlias)
						}
					}
				}

				if err != nil {
					return nil, err
				}

				cellVal = mapVal
				break
			}
		}

		// array
		if strings.HasPrefix(convertType, "[]") {
			tp := strings.TrimPrefix(convertType, "[]")
			seperator := ","
			if len(param) > 0 && len(param[0]) > 0 {
				seperator = param[0]
			}
			valueArr := make([]any, 0)
			if len(strVal) > 0 {
				values := strings.Split(strVal, seperator)
				for _, value := range values {
					var typeVal any
					typeVal, err = convertValue(filename, sheetname, keyname, tp, value)
					if err != nil {
						return nil, err
					}
					valueArr = append(valueArr, typeVal)
				}
			}

			cellVal = valueArr
			break
		}

		// enum type
		if value, found := excelEnumTypeRaws.Load(convertType); found {
			raw := value.(*ExcelEnumRaw)
			enumValue, success := raw.GetEnumValue(strings.Title(strVal))

			// 严格模式检测枚举
			if os.Getenv("STRICT") == "1" {
				if !success && convertType != "string" {
					switch strVal {
					case "", "{", "}":
						break
					default:
						GlobalResultReport.AddError(filename, sheetname, keyname, fmt.Sprintf("未能识别的枚举类型: %s", strVal))
					}
				}
			}
			cellVal = enumValue
			break
		}

		// string value
		if len(strVal) == 0 {
			cellVal = ""
			break
		}

		cellVal = strVal
	}

	return cellVal, nil
}

func convertMapValue(filename, sheetname, keyname, strType, strVal, separate string) (any, error) {
	// split type and value, example: map[int32]string => "int32" and "string"
	ts := strings.Split(strType, "[")
	t := ts[len(ts)-1]
	tt := strings.Split(t, "]")
	keyType, _ := convertType(tt[0])
	valueType, _ := convertType(tt[1])

	m := treemap.NewWith(func() map_utils.Comparator {
		switch keyType {
		case "int32":
			return map_utils.Int32Comparator
		case "string":
			return map_utils.StringComparator
		case "float32":
			return map_utils.Float32Comparator
		default:
			return map_utils.IntComparator
		}
	}())

	mapValues := strings.Split(strVal, separate)
	for _, oneMapValue := range mapValues {
		fields := strings.Split(oneMapValue, ",")
		if len(fields) < 2 {
			continue
		}

		k, err1 := convertValue(filename, sheetname, keyname, keyType, fields[0])
		if err1 != nil {
			return nil, err1
		}
		v, err2 := convertValue(filename, sheetname, keyname, valueType, fields[1])
		if err2 != nil {
			return nil, err2
		}
		m.Put(k, v)
	}

	return m, nil
}

// convert buildin map value
func convertBuildinMapValue(filename, sheetname, keyname, strType, strVal, separate string) (any, error) {
	// split type and value, example: map[int32]string => "int32" and "string"
	ts := strings.Split(strType, "[")
	t := ts[len(ts)-1]
	tt := strings.Split(t, "]")
	keyType, _ := convertType(tt[0])
	valueType, _ := getFullTypeName(tt[1])

	m := make(map[any]any)

	mapValues := strings.Split(strVal, separate)
	for _, oneMapValue := range mapValues {
		fields := strings.Split(oneMapValue, ",")
		if len(fields) < 2 {
			continue
		}

		k, err1 := convertValue(filename, sheetname, keyname, keyType, fields[0])
		if err1 != nil {
			return nil, err1
		}

		v, err2 := convertValue(filename, sheetname, keyname, valueType, fields[1], getPolySeparator(valueType))
		if err2 != nil {
			return nil, err2
		}
		m[k] = v
	}

	return m, nil
}

// get full enum or poly type name
func getFullTypeName(strType string) (string, error) {
	tp, buildin := convertType(strType)
	if !buildin {
		name := utils.SimpleTitle(tp)
		// enum
		if _, found := excelEnumTypeRaws.Load(name); found {
			return name, nil
		}

		// poly
		if _, found := excelPolyTypeRaws.Load(name); found {
			return fmt.Sprintf("*%s", name), nil
		}
	}

	return tp, nil
}

// 多态类型是否有parent
func PolyTypeNameHasParent(name string) bool {
	value, found := excelPolyTypeRaws.Load(name)
	if !found {
		return false
	}

	raw := value.(*ExcelPolyType)
	return len(raw.PolyParentName) > 0
}

// 名字是否注册为多态类型
func isPolyTypeName(name string) bool {
	_, found := excelPolyTypeRaws.Load(name)
	return found
}

// 获取多态类型分隔符
func getPolySeparator(name string) string {
	if value, found := excelPolyTypeRaws.Load(name); found {
		raw := value.(*ExcelPolyType)
		return raw.PolySeperator
	}

	// 默认分隔符为逗号
	return ","
}

// 获取多态类型元数据
func GetPolyRaw(name string) *ExcelPolyType {
	if value, found := excelPolyTypeRaws.Load(name); found {
		raw := value.(*ExcelPolyType)
		return raw
	}

	return nil
}

// 裁剪poly value为tagname, value, 是否包含tagname
func splitPolyValue(value string) (string, string, bool) {
	if !strings.Contains(value, ":") {
		return "", value, false
	}

	values := strings.Split(value, ":")
	if len(values) < 2 {
		return "", value, false
	}

	return values[0], strings.Join(values[1:], ":"), true
}
