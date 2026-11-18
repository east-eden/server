package excel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var (
	parsingTsFileRecord sync.Map // 正在处理的文件记录
)

type TsCodeData struct {
	className string
	sheetName string
}

type TsCodeGeneratorOptions struct {
	PackageName       string // name of this file's Go package
	FilePath          string
	FileName          string
	RequiredFilenames []string
	ImportPath        []string // import path of this file's Go package

	Structs   []*CodeStruct
	Variables []*CodeVariable
	Constants []*CodeConstant
	Polys     []*CodePoly
	Functions []*CodeFunction
	EnumRaws  []*ExcelEnumRaw
	TsData    *TsCodeData

	FileNamePrefix string
	FileNameSuffix string
}

type TsCodeGeneratorOption func(*TsCodeGeneratorOptions)

// code generator
type TsCodeGenerator struct {
	opts *TsCodeGeneratorOptions
	buf  bytes.Buffer
}

func defaultTsOptions() *TsCodeGeneratorOptions {
	return &TsCodeGeneratorOptions{
		RequiredFilenames: []string{},
		ImportPath:        []string{},
		Structs:           make([]*CodeStruct, 0),
		Variables:         make([]*CodeVariable, 0),
		Constants:         make([]*CodeConstant, 0),
		Polys:             make([]*CodePoly, 0),
		Functions:         make([]*CodeFunction, 0),
		EnumRaws:          make([]*ExcelEnumRaw, 0),
	}
}

func TsCodeDatas(data *TsCodeData) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.TsData = data
	}
}

func TsCodePackageName(packageName string) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.PackageName = packageName
	}
}

func TsCodeFilePath(filePath string) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.FilePath = filePath
	}
}

func TsCodeFileName(fileName string) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.FileName = fileName
	}
}

func TsCodeRequiredFileNames(fileNames []string) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.RequiredFilenames = make([]string, 0, len(fileNames))
		o.RequiredFilenames = append(o.RequiredFilenames, fileNames...)
	}
}

func TsCodeEnumRaws(raws []*ExcelEnumRaw) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.EnumRaws = make([]*ExcelEnumRaw, 0, len(raws))
		o.EnumRaws = append(o.EnumRaws, raws...)
	}
}

func TsCodeImportPath(importPath []string) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.ImportPath = make([]string, 0, len(importPath))
		o.ImportPath = append(o.ImportPath, importPath...)
	}
}

func TsCodeStructs(structs []*CodeStruct) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.Structs = make([]*CodeStruct, 0, len(structs))
		o.Structs = append(o.Structs, structs...)
	}
}

func TsCodeVariables(variables []*CodeVariable) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.Variables = make([]*CodeVariable, 0, len(variables))
		o.Variables = append(o.Variables, variables...)
	}
}

func TsCodeConstants(constants []*CodeConstant) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.Constants = make([]*CodeConstant, 0, len(constants))
		o.Constants = append(o.Constants, constants...)
	}
}

func TsCodePolys(polys []*CodePoly) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.Polys = make([]*CodePoly, 0, len(polys))
		o.Polys = append(o.Polys, polys...)
	}
}

func TsCodeFunctions(functions []*CodeFunction) TsCodeGeneratorOption {
	return func(o *TsCodeGeneratorOptions) {
		o.Functions = make([]*CodeFunction, 0, len(functions))
		o.Functions = append(o.Functions, functions...)
	}
}

func toTsType(tp string) string {
	switch tp {
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float", "float32", "float64":
		return "number"
	case "bool", "BOOL":
		return "boolean"
	case "string", "String":
		return "string"
	default:
		if strings.HasPrefix(tp, "list2:") {
			subType := strings.TrimPrefix(tp, "list2:")
			return toTsType(subType) + "[][]"

		} else if strings.HasPrefix(tp, "list:") {
			subType := strings.TrimPrefix(tp, "list:")
			return toTsType(subType) + "[]"
		} else if strings.HasPrefix(tp, "map[") {
			t := strings.Split(tp, "[")
			tt := strings.Split(t[len(t)-1], "]")

			keyType := toTsType(tt[0])
			valueType := toTsType(tt[1])
			return fmt.Sprintf("Map<%s, %s>", keyType, valueType)
		}

		if isPolyTypeName(tp) {
			confType := fmt.Sprintf("Conf%s", utils.SimpleTitle(tp))
			return confType
		} else {
			return utils.SimpleTitle(tp)
		}
	}
}

func genTsValue(filename, tp string, value any, jsonKey bool) string {
	// 多态处理
	polyDealFunc := func(polyType string) string {
		var strVal strings.Builder
		var classType string
		if value == nil {
			value = map[string]any{}
		}
		mapValue, ok := value.(map[string]any)
		if !ok {
			value, found := parsingTsFileRecord.Load(filename)
			if !found {
				log.Panic().Msgf("cannot find filename:%s", filename)
			} else {
				mapValue := value.(map[string]string)
				GlobalResultReport.AddError(filename, "", mapValue["key_name"], fmt.Sprintf("多态字段`%s` 类型转换失败,请检查该字段填写是否有误", polyType))
			}
		}

		if pointerType, found := mapValue["__pointer__type"]; found {
			classType = strings.TrimPrefix(pointerType.(string), "*")
		} else {
			classType = polyType
		}

		strVal.WriteString("{\n")
		if value, found := excelPolyTypeRaws.Load(classType); found {
			raw := value.(*ExcelPolyType)
			strVal.WriteString(fmt.Sprintf("\t\t\t\"className\": \"%s\"", classType))
			if len(raw.Polys) > 0 {
				strVal.WriteString(",")
			}
			strVal.WriteString("\n")

			for polyIdx, polyData := range raw.Polys {
				polyValue := genTsValue(filename, polyData.PolyType, mapValue[polyData.PolyName], false)
				polyDataStr := fmt.Sprintf("\t\t\t\"%s\": %s", utils.SimpleUntitle(polyData.PolyName), polyValue)
				if polyIdx == len(raw.Polys)-1 {
					polyDataStr += "\n"
				} else {
					polyDataStr += ",\n"
				}
				strVal.WriteString(polyDataStr)
			}

		}

		strVal.WriteString("\t\t}")
		return strVal.String()
	}

	switch tp {
	case "string":
		if value == nil {
			value = ""
		}

		if jsonKey {
			return fmt.Sprintf("%s", value)
		} else {
			jsonValue, _ := json.Marshal(value)
			return utils.ByteSlice2String(jsonValue)
		}
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		if value == nil {
			value = 0
		}
		return cast.ToString(cast.ToInt(value))
	case "float", "float32":
		if value == nil {
			value = 0
		}
		return cast.ToString(cast.ToFloat32(value))
	case "float64", "double":
		if value == nil {
			value = 0
		}
		return cast.ToString(cast.ToFloat64(value))
	case "bool", "BOOL", "boolean", "Boolean":
		if value == nil {
			value = false
		}
		return cast.ToString(cast.ToBool(value))
	case "any":
		return polyDealFunc(tp)
	default:
		// 枚举类型
		if _, found := excelEnumTypeRaws.Load(utils.SimpleTitle(tp)); found {
			if value == nil {
				value = -1
			}
			return cast.ToString(cast.ToInt(value))
		}

		// 多态类型
		polyTp := utils.SimpleTitle(tp)
		if _, found := excelPolyTypeRaws.Load(polyTp); found {
			return polyDealFunc(polyTp)
		}

		// list类型
		if strings.HasPrefix(tp, "list2:") {
			subType := strings.TrimPrefix(tp, "list2:")
			if value == nil {
				value = []any{}
			}
			values := value.([]any)
			var strVal strings.Builder
			strVal.WriteString("[ ")

			nextType := strings.Join([]string{"list", subType}, ":")
			for idx, v := range values {
				if idx == len(values)-1 {
					strVal.WriteString(genTsValue(filename, nextType, v, false))
				} else {
					strVal.WriteString(fmt.Sprintf("%s, ", genTsValue(filename, nextType, v, false)))
				}
			}

			strVal.WriteString(" ]")
			return strVal.String()

		} else if strings.HasPrefix(tp, "list:") {
			subType := strings.TrimPrefix(tp, "list:")
			if value == nil {
				value = []any{}
			}
			values := value.([]any)
			var strVal strings.Builder
			strVal.WriteString("[ ")

			for idx, v := range values {
				if idx == len(values)-1 {
					strVal.WriteString(genTsValue(filename, subType, v, false))
				} else {
					strVal.WriteString(fmt.Sprintf("%s, ", genTsValue(filename, subType, v, false)))
				}
			}

			strVal.WriteString(" ]")
			return strVal.String()
		} else if strings.HasPrefix(tp, "map[") {
			t := strings.Split(tp, "[")
			tt := strings.Split(t[len(t)-1], "]")
			if value == nil {
				value = []any{}
			}
			values, ok := value.(map[any]any)
			if !ok {
				GlobalResultReport.AddError(filename, "", tp, fmt.Sprintf("map类型`%s` 类型转换失败,请检查该字段填写是否有误", tp))
			}
			sortedKeys := make([]any, 0, len(values))
			for k := range values {
				sortedKeys = append(sortedKeys, k)
			}
			sort.Slice(sortedKeys, func(i, j int) bool {
				return cast.ToString(sortedKeys[i]) < cast.ToString(sortedKeys[j])
			})
			var strVal strings.Builder

			strVal.WriteString("{\n")
			length := 0
			for _, k := range sortedKeys {
				v := values[k]
				mapType := genTsValue(filename, tt[0], k, true)
				mapValue := genTsValue(filename, tt[1], v, false)
				polyDataStr := fmt.Sprintf("\t\t\t\"%s\": %s", mapType, mapValue)
				if length == len(values)-1 {
					polyDataStr += "\n"
				} else {
					polyDataStr += ",\n"
				}
				strVal.WriteString(polyDataStr)
				length++
			}

			strVal.WriteString("\t\t}")

			return strVal.String()
		}

		// default
		if value == nil {
			value = ""
		}
		return cast.ToString(value)
	}
}

func (g *TsCodeGenerator) P(v ...any) {
	for _, x := range v {
		switch x := x.(type) {
		// case GoIdent:
		// 	fmt.Fprint(&g.buf, g.QualifiedGoIdent(x))
		default:
			fmt.Fprint(&g.buf, x)
		}
	}
	fmt.Fprintln(&g.buf)
}

func (g *TsCodeGenerator) PBuf(buf *bytes.Buffer, v ...any) {
	for _, x := range v {
		switch x := x.(type) {
		// case GoIdent:
		// 	fmt.Fprint(&g.buf, g.QualifiedGoIdent(x))
		default:
			fmt.Fprint(buf, x)
		}
	}
	fmt.Fprintln(buf)
}

// generate one and map json buf
func (g *TsCodeGenerator) generateOneAndMapJsonBuf(excelRaw *ExcelFileRaw) (buf *bytes.Buffer) {
	buf = &bytes.Buffer{}

	// find last valid row idx
	lastValidIdx := 0
	for idx, row := range excelRaw.RowData {
		if len(row) == 0 {
			continue
		}
		lastValidIdx = idx
	}

	// generate other rows
	for idx, row := range excelRaw.RowData {
		if len(row) == 0 {
			continue
		}

		keyName := excelRaw.Keys[0]
		keyData := row[keyName]
		keyType, _ := excelRaw.FieldRaw.Get(keyName)
		keyOriginType := keyType.(*ExcelFieldRaw).originTp
		multiKey := (len(excelRaw.Keys) > 1)

		// record current parsing file, key, keyValue
		parsingTsFileRecord.Store(excelRaw.Filename, map[string]string{
			"key_name": cast.ToString(keyName),
			"key_data": cast.ToString(keyData),
		})

		if excelRaw.Type != ExcelSheetType_One {
			if multiKey {
				multiKeyValue := make([]string, 0)
				for _, keyName := range excelRaw.Keys {
					multiKeyValue = append(multiKeyValue, cast.ToString(row[keyName]))
				}
				g.PBuf(buf, fmt.Sprintf("\t\t\"%s\": {", strings.Join(multiKeyValue, "_")))
			} else {
				g.PBuf(buf, fmt.Sprintf("\t\t\"%s\": {", genTsValue(excelRaw.Filename, keyOriginType, keyData, true)))
			}
		}

		// gen all fields
		mapFieldsString := treemap.NewWithIntComparator()
		for fieldName, fieldValue := range row {
			if len(fieldName) == 0 {
				continue
			}

			v, _ := excelRaw.FieldRaw.Get(fieldName)
			vv := v.(*ExcelFieldRaw)
			value := genTsValue(excelRaw.Filename, vv.originTp, fieldValue, false)

			mapFieldsString.Put(vv.idx, fmt.Sprintf("\t\t\t\"%s\": %s", utils.SimpleUntitle(vv.originName), value))
		}

		// print all fields in sort
		it := mapFieldsString.Iterator()
		rowIdx := 0
		for it.Next() {
			rowIdx++
			str := it.Value().(string)
			if rowIdx != mapFieldsString.Size() {
				str += ","
			}
			g.PBuf(buf, str)
		}

		if excelRaw.Type != ExcelSheetType_One {
			if idx == lastValidIdx {
				g.PBuf(buf, "\t\t}")
			} else {
				g.PBuf(buf, "\t\t},")
			}
		}
	}

	return
}

// generate enum json buf
func (g *TsCodeGenerator) generateEnumJsonBuf(excelRaw *ExcelFileRaw) (buf *bytes.Buffer) {
	buf = &bytes.Buffer{}

	for _, excelEnumType := range excelRaw.EnumData {
		g.PBuf(buf, fmt.Sprintf("\t\"%s\": {", strings.Title(excelEnumType.EnumType)+"EnumInfo"))

		// generate enum rows
		for enumValueIdx, row := range excelEnumType.Enums {
			g.PBuf(buf, fmt.Sprintf("\t\t\"%s\": {", row.EnumValue))

			g.PBuf(buf, fmt.Sprintf("\t\t\t\"id\": %s,", row.EnumValue))
			g.PBuf(buf, fmt.Sprintf("\t\t\t\"key\": \"%s\",", row.EnumName))
			g.PBuf(buf, fmt.Sprintf("\t\t\t\"name\": \"%s\"", row.EnumDesc))

			if enumValueIdx == len(excelEnumType.Enums)-1 {
				g.PBuf(buf, "\t\t}")
			} else {
				g.PBuf(buf, "\t\t},")
			}
		}

		g.PBuf(buf, "\t},")
	}
	return
}

// generate ts code starts here
func (g *TsCodeGenerator) Generate() error {
	if len(g.opts.FilePath) == 0 {
		return errors.New("invalid file path")
	}

	if g.opts.TsData != nil {
		jsonOneAndMapTree := treemap.NewWithStringComparator()
		jsonEnumTree := treemap.NewWithStringComparator()

		excelFileRaws.Range(func(k, v any) bool {
			name := k.(string)
			excelRaw := v.(*ExcelFileRaw)
			if excelRaw.Type == ExcelSheetType_Poly {
				return true
			}

			if strings.HasSuffix(excelRaw.Sheetname, "client_skip") {
				return true
			}

			if excelRaw.Type == ExcelSheetType_Enum {
				buf := g.generateEnumJsonBuf(excelRaw)
				jsonEnumTree.Put(strings.TrimSuffix(name, ".json"), buf.String())
			}

			if excelRaw.Type == ExcelSheetType_One || excelRaw.Type == ExcelSheetType_Map {
				buf := g.generateOneAndMapJsonBuf(excelRaw)
				jsonOneAndMapTree.Put(strings.TrimSuffix(name, ".json"), buf.String())
			}
			return true
		})

		g.P("{")
		fileIdx := 0
		jsonEnumTree.All(func(k, v any) bool {
			fileIdx++
			buf := v.(string)
			builder := strings.Builder{}
			builder.Grow(len(buf) + 20)

			builder.WriteString(buf)
			g.P(builder.String())
			return true
		})

		fileIdx = 0
		jsonOneAndMapTree.All(func(k, v any) bool {
			fileIdx++
			name := k.(string)
			buf := v.(string)
			builder := strings.Builder{}
			builder.Grow(len(buf) + 20)

			builder.WriteString(fmt.Sprintf("\t\"%s\": {\n", name))
			builder.WriteString(buf)

			end := "\t}"
			if fileIdx != jsonOneAndMapTree.Size() {
				end += ","
			}

			builder.WriteString(end)
			g.P(builder.String())
			return true
		})
		g.P("}")
	}

	// generate required files
	if len(g.opts.RequiredFilenames) > 0 {

		// generate header
		g.P(`/***************************************
    This File is auto generated.
    Do on do any changes on the file !!!
    Any change will be override !!!
***************************************/
import { BaseConfig } from "../Config";
export interface EnumInfo {
	id: number;		// 枚举id
	key: string;	// 枚举key
	name: string;	// 枚举名称
}`)

		// generate fields body
		requiredFilenameTree := treemap.NewWithStringComparator()
		for _, filename := range g.opts.RequiredFilenames {
			requiredFilenameTree.Put(filename, true)
		}

		it := requiredFilenameTree.Iterator()
		for it.Next() {
			filename := it.Key().(string)
			value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
			if !found {
				continue
			}

			excelFileRaw := value.(*ExcelFileRaw)
			if strings.HasSuffix(excelFileRaw.Sheetname, "client_skip") {
				continue
			}

			confName := fmt.Sprintf("Conf%s", strings.Title(filename))
			g.P(fmt.Sprintf("export interface %s {	// %s %s", confName, excelFileRaw.Sheetname, excelFileRaw.ConfFilename))

			// export interface fields and type
			mapFieldsString := treemap.NewWithIntComparator()
			excelFileRaw.FieldRaw.All(func(k, v any) bool {
				vv := v.(*ExcelFieldRaw)
				if !vv.export {
					return true
				}

				if vv.option {
					mapFieldsString.Put(vv.idx, fmt.Sprintf("\t%s?: %s;\t\t// %s", utils.SimpleUntitle(vv.originName), toTsType(vv.originTp), vv.desc))
				} else {
					mapFieldsString.Put(vv.idx, fmt.Sprintf("\t%s: %s;\t\t// %s", utils.SimpleUntitle(vv.originName), toTsType(vv.originTp), vv.desc))
				}
				return true
			})

			// print fields in sort
			mapFieldsString.All(func(k, v any) bool {
				value := v.(string)
				g.P(value)
				return true
			})

			g.P("}")
		}

		// generate poly interface
		polyFilenameTree := treemap.NewWithStringComparator()
		excelPolyTypeRaws.Range(func(k, v any) bool {
			name := k.(string)
			poly := v.(*ExcelPolyType)
			polyFilenameTree.Put(name, poly)
			return true
		})

		polyFilenameTree.All(func(k, v any) bool {
			polyName := k.(string)
			excelPolyType := v.(*ExcelPolyType)
			value, _ := excelFileRaws.Load(v.(*ExcelPolyType).PolyFilename)
			excelFileRaw := value.(*ExcelFileRaw)
			confName := fmt.Sprintf("Conf%s", strings.Title(polyName))
			if len(excelPolyType.PolyParentName) > 0 {
				g.P(fmt.Sprintf("export interface %s extends %s { // %s %s",
					confName,
					fmt.Sprintf("Conf%s", strings.Title(excelPolyType.PolyParentName)),
					excelFileRaw.Sheetname,
					excelFileRaw.ConfFilename,
				))
			} else {
				g.P(fmt.Sprintf("export interface %s {	", confName))
			}

			// export interface fields and type
			for _, poly := range excelPolyType.OriginPolys {
				g.P(fmt.Sprintf("\t%s: %s;\t\t// %s", utils.SimpleUntitle(poly.PolyName), toTsType(poly.PolyType), poly.PolyComment))
			}

			g.P("}")
			return true
		})

		// generate get function
		g.P(`export default class ConfigData {
	public static getDynamicData<K,T>(tName: string, k: K):T{
		const t = ConfigData[tName];
		if(t){
			return t.get(k);
		}
	}
	public static __gd = ConfigData.getDynamicData;`)

		// one keepDict declaration
		oneFilesNameList := make([]string, 0, requiredFilenameTree.Size())
		requiredFilenameTree.All(func(k, _ any) bool {
			filename := k.(string)
			value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
			if !found {
				return true
			}

			excelFileRaw := value.(*ExcelFileRaw)
			if excelFileRaw.Type != ExcelSheetType_One {
				return true
			}

			oneFilesNameList = append(oneFilesNameList, fmt.Sprintf(`"%s"`, filename))
			return true
		})

		// generate _keepDict
		oneFilesStr := strings.Join(oneFilesNameList, ", ")
		g.P(fmt.Sprintf("\tpublic static _keepDict = new Set([%s]);", oneFilesStr))

		// one and map declaration
		requiredFilenameTree.All(func(k, _ any) bool {
			filename := k.(string)
			value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
			if !found {
				return true
			}

			excelFileRaw := value.(*ExcelFileRaw)
			if len(excelFileRaw.Keys) == 0 {
				// no data
				return true
			}

			if strings.HasSuffix(excelFileRaw.Sheetname, "client_skip") {
				return true
			}

			var keyType string
			var multiKey bool
			if excelFileRaw != nil {
				if len(excelFileRaw.Keys) > 1 {
					keyType = "string"
					multiKey = true
				} else {
					v, _ := excelFileRaw.FieldRaw.Get(excelFileRaw.Keys[0])
					keyType = toTsType(v.(*ExcelFieldRaw).originTp)
				}
			}

			confName := fmt.Sprintf("Conf%s", strings.Title(filename))
			if excelFileRaw.Type == ExcelSheetType_One {
				// if strings.Title(filename) == "Config" {
				g.P(fmt.Sprintf("\tpublic static %s: %s", utils.SimpleUntitle(filename), confName))
			} else {
				if multiKey {
					keyTypes := make([]string, 0, len(excelFileRaw.Keys))
					for _, keyName := range excelFileRaw.Keys {
						v, _ := excelFileRaw.FieldRaw.Get(keyName)
						vv := v.(*ExcelFieldRaw)
						keyTypes = append(keyTypes, toTsType(vv.originTp))
					}
					g.P(fmt.Sprintf("\tpublic static %s: %s<%s, [%s]>", utils.SimpleUntitle(filename), "BaseConfig", confName, strings.Join(keyTypes, ", ")))
				} else {
					g.P(fmt.Sprintf("\tpublic static %s: %s<%s, %s>", utils.SimpleUntitle(filename), "BaseConfig", confName, keyType))
				}
			}
			return true
		})

		enumFilenameTree := treemap.NewWithStringComparator()
		for _, enumRaw := range g.opts.EnumRaws {
			enumFilenameTree.Put(enumRaw.enumOriginType, enumRaw)
		}
		enumFilenameTree.All(func(k, _ any) bool {
			filename := k.(string)
			g.P(fmt.Sprintf("\tpublic static %s: %s<%s, %s>", fmt.Sprintf("%sEnumInfo", utils.SimpleUntitle(filename)), "BaseConfig", "EnumInfo", filename))
			return true
		})

		g.P("}")
	}

	// generate enum raws
	if len(g.opts.EnumRaws) > 0 {
		enumFilenameTree := treemap.NewWithStringComparator()
		for _, enumRaw := range g.opts.EnumRaws {
			enumFilenameTree.Put(enumRaw.enumOriginType, enumRaw)
		}

		it := enumFilenameTree.Iterator()
		for it.Next() {
			enumRaw := it.Value().(*ExcelEnumRaw)
			g.P(fmt.Sprintf("export enum %s {\t\t//%s %s", enumRaw.enumType, enumRaw.sheetName, enumRaw.confName))

			for _, enumValue := range enumRaw.orderedValues {
				g.P(fmt.Sprintf("\t%s = %d,\t\t// %s", enumValue.Name, enumValue.Value, enumValue.Alias))
			}

			g.P("}")
		}
	}

	return os.WriteFile(g.opts.FilePath, g.buf.Bytes(), 0666)
}

func NewTsCodeGenerator(options ...TsCodeGeneratorOption) *TsCodeGenerator {
	g := &TsCodeGenerator{
		opts: defaultTsOptions(),
	}

	for _, o := range options {
		o(g.opts)
	}

	return g
}

// generate required ts code
func generateRequiredTsCode(exportPath string, requiredFiles []string) error {
	filename := "Configs.ts"

	// required file names
	names := make([]string, 0, len(requiredFiles))
	for _, filename := range requiredFiles {
		names = append(names, strings.TrimSuffix(filename, ".json"))
	}

	// enum raws
	enumRaws := make([]*ExcelEnumRaw, 0, 32)
	excelEnumTypeRaws.Range(func(k, v any) bool {
		raw := v.(*ExcelEnumRaw)
		enumRaws = append(enumRaws, raw)
		return true
	})

	g := NewTsCodeGenerator(
		TsCodeFilePath(filepath.Join(exportPath, filename)),
		TsCodeFileName(filename),
		TsCodeRequiredFileNames(names),
		TsCodeEnumRaws(enumRaws),
	)

	return g.Generate()
}

// generateEntryTsCode generates the contents of a .go file.
func generateEntryTsCode(packageName, exportPath string) error {
	// metaName := strings.Split(excelFileRaw.ConfFilename, ".")[0]

	// 文件名有中文时，使用Bean作为文件名
	// for _, runeValue := range metaName {
	// 	if unicode.Is(unicode.Han, runeValue) {
	// 		metaName = excelFileRaw.Beanname
	// 		break
	// 	}
	// }

	// if !utf8string.NewString(metaName).IsASCII() {
	// 	py := pinyin.LazyPinyin(metaName, pinyin.NewArgs())
	// 	metaName = strings.Join(py, "_")
	// }

	filename := "data.json"
	g := NewTsCodeGenerator(
		TsCodeDatas(&TsCodeData{
			// className: metaName,
			// sheetName: excelFileRaw.Sheetname,
		}),

		TsCodeFilePath(filepath.Join(exportPath, filename)),
		TsCodeFileName(filename),
	)

	err := g.Generate()
	if !utils.ErrCheck(err, "generate ts code failed", g.opts.FilePath) {
		return err
	}

	return nil
}
