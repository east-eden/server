package excel

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"golang.org/x/exp/utf8string"
)

var (
	parsingFileRecord sync.Map // 正在处理的文件记录
)

type LuaCodeData struct {
	className string
	sheetName string
	excelRaw  *ExcelFileRaw
}

type LuaCodeGeneratorOptions struct {
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
	LuaData   *LuaCodeData

	FileNamePrefix string
	FileNameSuffix string
}

type LuaCodeGeneratorOption func(*LuaCodeGeneratorOptions)

// code generator
type LuaCodeGenerator struct {
	opts *LuaCodeGeneratorOptions
	buf  bytes.Buffer
}

func defaultLuaOptions() *LuaCodeGeneratorOptions {
	return &LuaCodeGeneratorOptions{
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

func LuaCodeDatas(data *LuaCodeData) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.LuaData = data
	}
}

func LuaCodePackageName(packageName string) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.PackageName = packageName
	}
}

func LuaCodeFilePath(filePath string) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.FilePath = filePath
	}
}

func LuaCodeFileName(fileName string) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.FileName = fileName
	}
}

func LuaCodeRequiredFileNames(fileNames []string) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.RequiredFilenames = fileNames
	}
}

func LuaCodeEnumRaws(raws []*ExcelEnumRaw) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.EnumRaws = raws
	}
}

func LuaCodeImportPath(importPath []string) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.ImportPath = importPath
	}
}

func LuaCodeStructs(structs []*CodeStruct) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.Structs = structs
	}
}

func LuaCodeVariables(variables []*CodeVariable) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.Variables = variables
	}
}

func LuaCodeConstants(constants []*CodeConstant) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.Constants = constants
	}
}

func LuaCodePolys(polys []*CodePoly) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.Polys = polys
	}
}

func LuaCodeFunctions(functions []*CodeFunction) LuaCodeGeneratorOption {
	return func(o *LuaCodeGeneratorOptions) {
		o.Functions = functions
	}
}

func toLuaType(tp string) string {
	switch tp {
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float", "float32", "float64":
		return "number"
	case "bool", "BOOL":
		return "boolean"
	default:
		if strings.HasPrefix(tp, "list2:") {
			subType := strings.TrimPrefix(tp, "list2:")
			return toLuaType(subType) + "[][]"

		} else if strings.HasPrefix(tp, "list:") {
			subType := strings.TrimPrefix(tp, "list:")
			return toLuaType(subType) + "[]"
		}

		return tp
	}
}

func genLuaValue(filename, tp string, value any) string {
	// 多态处理
	polyDealFunc := func(polyType string) string {
		var strVal strings.Builder
		var classType string
		if value == nil {
			value = map[string]any{}
		}
		mapValue, ok := value.(map[string]any)
		if !ok {
			value, found := parsingFileRecord.Load(filename)
			if !found {
				log.Panic().Msgf("cannot find filename:%s", filename)
			} else {
				mapValue := value.(map[string]string)
				log.Error().
					Str("文件名", filename).
					Str("主键名", mapValue["key_name"]).
					Str("主键值", mapValue["key_data"]).
					Str("字段名", polyType).
					Msg("请检查该字段填写是否有误?")
				os.Exit(0)
			}
		}

		if pointerType, found := mapValue["__pointer__type"]; found {
			classType = strings.TrimPrefix(pointerType.(string), "*")
		} else {
			classType = polyType
		}

		strVal.WriteString("{\n")
		strVal.WriteString(fmt.Sprintf("\t\t\tclass = \"%s\",\n", classType))

		if value, found := excelPolyTypeRaws.Load(classType); found {
			raw := value.(*ExcelPolyType)
			for _, polyData := range raw.Polys {
				polyValue := genLuaValue(filename, polyData.PolyType, mapValue[polyData.PolyName])
				strVal.WriteString(fmt.Sprintf("\t\t\t%s = %s,\n", utils.SimpleUntitle(polyData.PolyName), polyValue))
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
		return fmt.Sprintf("'%s'", strings.TrimSpace(cast.ToString(value)))
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		if value == nil {
			value = 0
		}
		return cast.ToString(cast.ToInt(value))
	case "float", "float32", "float64":
		if value == nil {
			value = 0
		}
		return cast.ToString(cast.ToFloat32(value))
	case "bool", "BOOL", "boolean", "Boolean":
		if value == nil {
			value = false
		}
		return cast.ToString(cast.ToBool(value))
	case "any":
		return polyDealFunc(tp)
	default:
		// 枚举类型
		if _, found := excelEnumTypeRaws.Load(tp); found {
			if value == nil {
				value = -1
			}
			return cast.ToString(cast.ToInt(value))
		}

		// 多态类型
		polyTp := strings.Title(tp)
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
			strVal.WriteString("{ ")

			nextType := strings.Join([]string{"list", subType}, ":")
			for _, v := range values {
				strVal.WriteString(fmt.Sprintf("%s, ", genLuaValue(filename, nextType, v)))
			}

			strVal.WriteString(" }")
			return strVal.String()
		} else if strings.HasPrefix(tp, "list:") {
			subType := strings.TrimPrefix(tp, "list:")
			if value == nil {
				value = []any{}
			}
			values := value.([]any)
			var strVal strings.Builder
			strVal.WriteString("{ ")

			for _, v := range values {
				strVal.WriteString(fmt.Sprintf("%s, ", genLuaValue(filename, subType, v)))
			}

			strVal.WriteString(" }")
			return strVal.String()
		}

		// default
		if value == nil {
			value = ""
		}
		return cast.ToString(value)
	}
}

func (g *LuaCodeGenerator) P(v ...any) {
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

// generate lua code starts here
func (g *LuaCodeGenerator) Generate() error {
	if len(g.opts.FilePath) == 0 {
		return errors.New("invalid file path")
	}

	if g.opts.LuaData != nil {
		// generate comments
		if len(g.opts.LuaData.className) > 0 {
			g.P("---@class ", g.opts.LuaData.className)
		}

		if len(g.opts.LuaData.sheetName) > 0 {
			g.P(fmt.Sprintf("---@file %s", g.opts.LuaData.sheetName))
		}

		// excel original data
		excelRaw := g.opts.LuaData.excelRaw

		// fields comments
		fieldCommentLines := make([]string, excelRaw.FieldRaw.Size())
		it := excelRaw.FieldRaw.Iterator()
		for it.Next() {
			field := it.Value().(*ExcelFieldRaw)
			if !field.export {
				continue
			}

			line := fmt.Sprintf("---@field %s %s @%s", field.originName, toLuaType(field.originTp), field.desc)
			fieldCommentLines[field.idx] = line
		}

		// print fields comment in sort
		for _, v := range fieldCommentLines {
			if len(v) == 0 {
				continue
			}

			g.P(v)
		}
		g.P()

		// generate lua data
		if excelRaw.Type != ExcelSheetType_One {
			g.P("---@type ", g.opts.LuaData.className, "[]")
		}
		g.P("local data =")
		g.P("{")

		// generate other rows
		for _, row := range excelRaw.RowData {
			if len(row) == 0 {
				continue
			}

			keyName := excelRaw.Keys[0]
			keyData := row[keyName]
			keyType, _ := excelRaw.FieldRaw.Get(keyName)
			keyOriginType := keyType.(*ExcelFieldRaw).originTp
			multiKey := (len(excelRaw.Keys) > 1)

			// record current parsing file, key, keyValue
			parsingFileRecord.Store(excelRaw.Filename, map[string]string{
				"key_name": cast.ToString(keyName),
				"key_data": cast.ToString(keyData),
			})

			if excelRaw.Type != ExcelSheetType_One {
				if multiKey {
					multiKeyValue := make([]string, 0)
					for _, keyName := range excelRaw.Keys {
						multiKeyValue = append(multiKeyValue, cast.ToString(row[keyName]))
					}
					g.P(fmt.Sprintf("\t[\"%s\"] = {", strings.Join(multiKeyValue, "_")))
				} else {
					g.P(fmt.Sprintf("\t[%s] = {", genLuaValue(excelRaw.Filename, keyOriginType, keyData)))
				}
			}

			classLine := fmt.Sprintf("\t\tclass = '%s',", g.opts.LuaData.className)
			g.P(classLine)

			// gen all fields
			mapFieldsString := treemap.NewWithIntComparator()
			for fieldName, fieldValue := range row {
				if len(fieldName) == 0 {
					continue
				}

				v, _ := excelRaw.FieldRaw.Get(fieldName)
				vv := v.(*ExcelFieldRaw)
				value := genLuaValue(excelRaw.Filename, vv.originTp, fieldValue)
				mapFieldsString.Put(vv.idx, fmt.Sprintf("\t\t%s = %s,", vv.originName, value))
			}

			// print all fields in sort
			it := mapFieldsString.Iterator()
			for it.Next() {
				str := it.Value().(string)
				g.P(str)
			}

			if excelRaw.Type != ExcelSheetType_One {
				g.P("\t},")
			}
		}

		g.P("}")
		g.P("return data")
	}

	// generate required files
	if len(g.opts.RequiredFilenames) > 0 {
		g.P("local os = require 'Config.Constants'")
		g.P("local cfgs = {}")

		tree := treemap.NewWithStringComparator()
		for _, filename := range g.opts.RequiredFilenames {
			tree.Put(filename, true)
		}

		it := tree.Iterator()
		for it.Next() {
			filename := it.Key().(string)
			g.P(fmt.Sprintf("cfgs['%s'] = require('Config.Data.%s')", filename, filename))
		}

		g.P("return cfgs")
	}

	// generate enum raws
	if len(g.opts.EnumRaws) > 0 {
		head := `
local insert = table.insert
local ipairs = ipairs
local setmetatable = setmetatable
local find = string.find
local sub = string.sub
local rawget = rawget
local bool = os.bool
local int = os.int
local long = os.long
local string = os.string
local float = os.float
local list = os.list
local set = os.set
local map = os.map

local function global_var(namespace)
	local t = _G
	local idx = 1
	while true do
		local start, ends = find(namespace, '.', idx, true)
		local subname = sub(namespace, idx, start and start - 1)
		local subt = t[subname]
		if not subt then
			subt = {}
			t[subname] = subt
		end
		t = subt
		if start then
			idx = ends + 1
		else
			return t
		end
	end
end`
		g.P(head)
		g.P()

		tree := treemap.NewWithStringComparator()
		for _, enumRaw := range g.opts.EnumRaws {
			tree.Put(enumRaw.enumOriginType, enumRaw)
		}

		it := tree.Iterator()
		for it.Next() {
			originType := it.Key().(string)
			enumRaw := it.Value().(*ExcelEnumRaw)
			g.P(fmt.Sprintf("---@file %s/%s", enumRaw.sheetName, enumRaw.confName))
			g.P(fmt.Sprintf("---@class Cfg%s", originType))
			g.P(fmt.Sprintf("global_var('cfg')['%s'] = {", originType))

			sortedEnum := treemap.NewWithIntComparator()
			for name, value := range enumRaw.nameValue {
				sortedEnum.Put(value, name)
			}

			it := sortedEnum.Iterator()
			for it.Next() {
				name := it.Value().(string)
				g.P(fmt.Sprintf("\t%s = %d,", name, enumRaw.nameValue[name]))
			}

			g.P("}")
		}

		g.P("return os")
	}

	return os.WriteFile(g.opts.FilePath, g.buf.Bytes(), 0666)
}

func NewLuaCodeGenerator(options ...LuaCodeGeneratorOption) *LuaCodeGenerator {
	g := &LuaCodeGenerator{
		opts: defaultLuaOptions(),
	}

	for _, o := range options {
		o(g.opts)
	}

	return g
}

// generate required lua code
func generateRequiredLuaCode(exportPath string, requiredFiles []string) error {
	filename := "Configs.lua"

	// required file names
	names := make([]string, 0, len(requiredFiles))
	for _, filename := range requiredFiles {
		names = append(names, strings.TrimSuffix(filename, ".json"))
	}

	g := NewLuaCodeGenerator(
		LuaCodeFilePath(filepath.Join(exportPath, filename)),
		LuaCodeFileName(filename),
		LuaCodeRequiredFileNames(names),
	)

	return g.Generate()
}

// generate enum lua code
func generateEnumLuaCode(exportPath string) error {
	filename := "Constants.lua"

	enumRaws := make([]*ExcelEnumRaw, 0, 32)
	excelEnumTypeRaws.Range(func(k, v any) bool {
		raw := v.(*ExcelEnumRaw)
		enumRaws = append(enumRaws, raw)
		return true
	})

	g := NewLuaCodeGenerator(
		LuaCodeFilePath(filepath.Join(exportPath, filename)),
		LuaCodeFileName(filename),
		LuaCodeEnumRaws(enumRaws),
	)

	return g.Generate()
}

// generateEntryLuaCode generates the contents of a .go file.
func generateEntryLuaCode(packageName, exportPath string, excelFileRaw *ExcelFileRaw) error {
	metaName := strings.Split(excelFileRaw.ConfFilename, ".")[0]

	// 文件名有中文时，使用Bean作为文件名
	for _, runeValue := range metaName {
		if unicode.Is(unicode.Han, runeValue) {
			metaName = excelFileRaw.Beanname
			break
		}
	}

	if !utf8string.NewString(metaName).IsASCII() {
		// py := pinyin.LazyPinyin(metaName, pinyin.NewArgs())
		// metaName = strings.Join(py, "_")
	}

	g := NewLuaCodeGenerator(
		LuaCodeDatas(&LuaCodeData{
			className: metaName,
			sheetName: excelFileRaw.Sheetname,
			excelRaw:  excelFileRaw,
		}),

		LuaCodeFileName(strings.Title(metaName)),
		LuaCodeFilePath(filepath.Join(exportPath, metaName+".lua")),
	)

	err := g.Generate()
	if !utils.ErrCheck(err, "generate lua code failed", g.opts.FilePath) {
		return err
	}

	return nil
}
