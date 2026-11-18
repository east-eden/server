package excel

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/exp/utf8string"

	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/spf13/cast"
)

type CodeEnumType int
type CodeEnumComment string

type CodeFieldType string
type CodeFieldTags string
type CodeFieldComment string

// single key load function
var singleKeyLoadFunctionBody string = `
	__lowerReplace__Entries = &__upperReplace__Entries{
		Rows: make(map[__keytype__]*__upperReplace__Entry),
		Arrays: make([]*__upperReplace__Entry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &__upperReplace__Entry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return __lowerReplace__Entries, err
		}

		__Replace__Poly__String

	 	__lowerReplace__Entries.Rows[entry.__keyname__] = entry
	 	__lowerReplace__Entries.Arrays = append(__lowerReplace__Entries.Arrays, entry)
	}

	sort.Slice(__lowerReplace__Entries.Arrays, func(i, j int) bool {
		return __lowerReplace__Entries.Arrays[i].__keyname__ < __lowerReplace__Entries.Arrays[j].__keyname__
	})

	excel.AddArrayGetter("__lowerReplace__.json", __lowerReplace__Entries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return __lowerReplace__Entries, nil
	`

// multi key load function
var multiKeyLoadFunctionBody string = `
	__lowerReplace__Entries = &__upperReplace__Entries{
		Rows: make(map[__keytype__]*__upperReplace__Entry),
		Arrays: make([]*__upperReplace__Entry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &__upperReplace__Entry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return __lowerReplace__Entries, err
		}

		key := fmt.Sprintf("%s", %s)
	 	__lowerReplace__Entries.Rows[key] = entry
	 	__lowerReplace__Entries.Arrays = append(__lowerReplace__Entries.Arrays, entry)
	}

	sort.Slice(__lowerReplace__Entries.Arrays, func(i, j int) bool {
		return __lowerReplace__Entries.Arrays[i].__keyname__ < __lowerReplace__Entries.Arrays[j].__keyname__
	})

	excel.AddArrayGetter("__lowerReplace__.json", __lowerReplace__Entries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return __lowerReplace__Entries, nil
	`

// get enum by string type function
var getEnumFunctionBody string = `switch tp {
__case__
	}
__default__`

// stringer enum by string type function
var stringerEnumFunctionBody string = `
__if__
__default__`

// code struct
type CodeStruct struct {
	name     string
	comment  string
	fieldRaw *treemap.Map
}

// code variable
type CodeVariable struct {
	name    string
	tp      string
	comment string
}

// code constants element
type CodeConstantElem struct {
	name    string
	value   string
	comment string
}

// code constants
type CodeConstant struct {
	tp           string
	sheetName    string
	confFileName string
	elems        []*CodeConstantElem
}

// code polys element
type CodePolyElem struct {
	name    string
	tp      string
	tag     string
	comment string
}

// code polys
type CodePoly struct {
	tp           string
	sep          string
	parent       string
	comment      string
	sheetName    string
	confFilename string
	elems        []*CodePolyElem
}

// code function
type CodeFunction struct {
	receiver    string
	receiverPtr bool
	name        string
	comment     string
	parameters  []string
	retType     string
	body        string
}

type GoCodeGeneratorOptions struct {
	PackageName string // name of this file's Go package
	FilePath    string
	FileName    string
	ImportPath  []string // import path of this file's Go package

	Structs   []*CodeStruct
	Variables []*CodeVariable
	Constants []*CodeConstant
	Polys     []*CodePoly
	Functions []*CodeFunction

	FileNamePrefix string
	FileNameSuffix string

	Comments string
}

type GoCodeGeneratorOption func(*GoCodeGeneratorOptions)

// code generator
type GoCodeGenerator struct {
	opts *GoCodeGeneratorOptions
	buf  bytes.Buffer
}

func defaultGoOptions() *GoCodeGeneratorOptions {
	return &GoCodeGeneratorOptions{
		ImportPath: []string{},
		Structs:    make([]*CodeStruct, 0),
		Variables:  make([]*CodeVariable, 0),
		Constants:  make([]*CodeConstant, 0),
		Polys:      make([]*CodePoly, 0),
		Functions:  make([]*CodeFunction, 0),
	}
}

func GoCodePackageName(packageName string) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.PackageName = packageName
	}
}

func GoCodeFilePath(filePath string) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.FilePath = filePath
	}
}

func GoCodeFileName(fileName string) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.FileName = fileName
	}
}

func GoCodeImportPath(importPath []string) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.ImportPath = importPath
	}
}

func GoCodeStructs(structs []*CodeStruct) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.Structs = make([]*CodeStruct, 0, len(structs))
		o.Structs = append(o.Structs, structs...)
	}
}

func GoCodeVariables(variables []*CodeVariable) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.Variables = make([]*CodeVariable, 0, len(variables))
		o.Variables = append(o.Variables, variables...)
	}
}

func GoCodeConstants(constants []*CodeConstant) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.Constants = make([]*CodeConstant, 0, len(constants))
		o.Constants = append(o.Constants, constants...)
	}
}

func GoCodePolys(polys []*CodePoly) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.Polys = make([]*CodePoly, 0, len(polys))
		o.Polys = append(o.Polys, polys...)
	}
}

func GoCodeFunctions(functions []*CodeFunction) GoCodeGeneratorOption {
	return func(o *GoCodeGeneratorOptions) {
		o.Functions = make([]*CodeFunction, 0, len(functions))
		o.Functions = append(o.Functions, functions...)
	}
}

func (g *GoCodeGenerator) P(v ...any) {
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

// generate go code starts here
func (g *GoCodeGenerator) Generate() error {
	if len(g.opts.FilePath) == 0 {
		return errors.New("invalid file path")
	}

	if len(g.opts.PackageName) == 0 {
		return errors.New("invalid package name")
	}

	// generate package
	g.P("package ", g.opts.PackageName)
	g.P()

	// generate import path
	g.P("import (")
	for _, path := range g.opts.ImportPath {
		g.P("\t\"", path, "\"")
	}
	g.P(")")
	g.P()

	// generate variables
	for _, v := range g.opts.Variables {
		variableLine := fmt.Sprintf("var\t%-15s\t%-15s\t//%-15s", v.name, v.tp, v.comment)
		g.P(variableLine)
		g.P()
	}

	// generate constants
	for _, constant := range g.opts.Constants {
		g.P(fmt.Sprintf("// %s/%s", constant.sheetName, constant.confFileName))

		constType := constant.tp
		typeLine := fmt.Sprintf("type %s int", constType)
		g.P(typeLine)

		g.P("const (")
		for _, elem := range constant.elems {
			constLine := fmt.Sprintf("\t%s_%s\t%-15s\t=\t%-15s\t//%-15s", constType, elem.name, constType, elem.value, elem.comment)
			g.P(constLine)
		}
		g.P(")")
		g.P()
	}

	// generate polys
	for _, poly := range g.opts.Polys {
		// struct comment
		if len(poly.comment) > 0 {
			g.P("// ", poly.comment)
		}

		polyType := poly.tp

		// struct begin
		g.P(fmt.Sprintf("// %s/%s", poly.sheetName, poly.confFilename))
		g.P("type ", polyType, " struct {")

		// parent
		// if len(poly.parent) > 0 {
		// 	tag := fmt.Sprintf("`json:\"%s,omitempty\"`", poly.parent)
		// 	fieldLine := fmt.Sprintf("\t%-15s\t%-20s\t//%-10s", poly.parent, tag, "父类")
		// 	g.P(fieldLine)

		// 	// add a blank
		// 	g.P()
		// }

		// struct fields
		for _, elem := range poly.elems {
			tp, _ := getFullTypeName(elem.tp)
			fieldLine := fmt.Sprintf("\t%-15s\t%-20s\t%-20s\t//%-10s", elem.name, tp, elem.tag, elem.comment)
			g.P(fieldLine)
		}

		// struct end
		g.P("}")
		g.P()
	}

	// generate structs
	for _, s := range g.opts.Structs {
		// struct comment
		if len(s.comment) > 0 {
			g.P("// ", s.comment)
		}

		// struct begin
		g.P("type ", s.name, " struct {")

		// struct fields
		fieldLines := make([]string, s.fieldRaw.Size())
		it := s.fieldRaw.Iterator()
		for it.Next() {
			fieldRaw := it.Value().(*ExcelFieldRaw)

			// don't need import
			if !fieldRaw.export {
				continue
			}

			tpName, _ := getFullTypeName(fieldRaw.tp)
			fieldLine := fmt.Sprintf("\t%-15s\t%-20s\t%-20s\t//%-10s", it.Key(), tpName, fieldRaw.tag, fieldRaw.desc)
			fieldLines[fieldRaw.idx] = fieldLine
		}

		// print struct field in sort
		for _, v := range fieldLines {
			if len(v) == 0 {
				continue
			}

			g.P(v)
		}

		// struct end
		g.P("}")
		g.P()
	}

	// generate functions
	for _, f := range g.opts.Functions {
		// function comment
		if len(f.comment) > 0 {
			g.P("// ", f.comment)
		}

		// function receiver
		var receiver string
		if len(f.receiver) > 0 {
			if f.receiverPtr {
				receiver = fmt.Sprintf("(e *%s)", f.receiver)
			} else {
				receiver = fmt.Sprintf("(e %s)", f.receiver)
			}
		}

		// function parameters
		parameters := strings.Join(f.parameters, ", ")

		// function begin
		g.P("func ", receiver, " ", f.name, "(", parameters, ") ", f.retType, " {")

		// function body
		g.P("\t", f.body)

		// function end
		g.P("}")
		g.P()
	}

	return os.WriteFile(g.opts.FilePath, g.buf.Bytes(), 0666)
}

func NewGoCodeGenerator(options ...GoCodeGeneratorOption) *GoCodeGenerator {
	g := &GoCodeGenerator{
		opts: defaultGoOptions(),
	}

	for _, o := range options {
		o(g.opts)
	}

	return g
}

// generate enum go code
func generateEnumGoCode(packageName, exportPath string, excelFileRaw *ExcelFileRaw) error {
	metaName := strings.Split(excelFileRaw.Filename, ".")[0]
	if utf8string.NewString(metaName).IsASCII() {
		for _, v := range metaName {
			metaName = string(unicode.ToLower(v)) + metaName[1:]
			break
		}
	}

	fileName := utils.SimpleTitle(metaName)

	constants := make([]*CodeConstant, 0, 8)
	codeFunctions := make([]*CodeFunction, 0)
	for _, enumType := range excelFileRaw.EnumData {
		constant := &CodeConstant{
			tp:           enumType.EnumType,
			sheetName:    excelFileRaw.Sheetname,
			confFileName: excelFileRaw.ConfFilename,
			elems:        make([]*CodeConstantElem, 0, len(enumType.Enums)),
		}

		minVal := 1
		maxVal := 0
		for _, enum := range enumType.Enums {
			constant.elems = append(constant.elems, &CodeConstantElem{
				name:    enum.EnumName,
				value:   enum.EnumValue,
				comment: enum.EnumDesc,
			})

			ev := cast.ToInt(enum.EnumValue)
			if ev < minVal {
				minVal = ev
			}

			if ev > maxVal {
				maxVal = ev
			}
		}

		if len(enumType.Enums) > 0 {
			// add enum header
			constant.elems = append([]*CodeConstantElem{
				{
					name:    "BEGIN",
					value:   cast.ToString(minVal),
					comment: "枚举开始",
				},
			}, constant.elems...)

			// add enum end
			constant.elems = append(constant.elems, &CodeConstantElem{
				name:    "END",
				value:   cast.ToString(maxVal + 1),
				comment: "枚举结束",
			})
		}

		constants = append(constants, constant)

		// stringer function
		stringerFunction := &CodeFunction{
			name:        "String",
			receiver:    enumType.EnumType,
			receiverPtr: false,
			retType:     "string",
			body:        stringerEnumFunctionBody,
		}

		var enumStringCaseBuffer bytes.Buffer
		for _, enum := range enumType.Enums {
			enumStringCaseBuffer.WriteString(fmt.Sprintf("\tif e == %s_%s {\n\t\treturn \"%s\"\n\t}\n", enumType.EnumType, enum.EnumName, enum.EnumName))
		}
		var enumStringDefaultBuffer bytes.Buffer
		enumStringDefaultBuffer.WriteString("\treturn \"\"\n")

		stringerFunction.body = strings.ReplaceAll(stringerFunction.body, "__if__", enumStringCaseBuffer.String())
		stringerFunction.body = strings.ReplaceAll(stringerFunction.body, "__default__", enumStringDefaultBuffer.String())
		codeFunctions = append(codeFunctions, stringerFunction)

		// get enum by string name function
		getEnumFunction := &CodeFunction{
			name:       fmt.Sprintf("Get%s", enumType.EnumType),
			parameters: []string{"tp string"},
			retType:    fmt.Sprintf("(%s, bool)", enumType.EnumType),
			body:       getEnumFunctionBody,
		}

		var caseBuffer bytes.Buffer
		for _, enum := range enumType.Enums {
			caseBuffer.WriteString(fmt.Sprintf("\tcase \"%s\":\n\t\treturn %s_%s, true\n", enum.EnumName, enumType.EnumType, enum.EnumName))
		}

		var defaultBuffer bytes.Buffer
		defaultBuffer.WriteString("\treturn -1, false\n")

		getEnumFunction.body = strings.ReplaceAll(getEnumFunction.body, "__case__", caseBuffer.String())
		getEnumFunction.body = strings.ReplaceAll(getEnumFunction.body, "__default__", defaultBuffer.String())

		codeFunctions = append(codeFunctions, getEnumFunction)

		// get enum by string alias function
		getEnumAliasFunction := &CodeFunction{
			name:       fmt.Sprintf("Get%s_ByAlias", enumType.EnumType),
			parameters: []string{"tp string"},
			retType:    fmt.Sprintf("(%s, bool)", enumType.EnumType),
			body:       getEnumFunctionBody,
		}

		var caseBuffer2 bytes.Buffer
		for _, enum := range enumType.Enums {
			caseBuffer2.WriteString(fmt.Sprintf("\tcase \"%s\":\n\t\treturn %s_%s, true\n", enum.EnumDesc, enumType.EnumType, enum.EnumName))
		}

		var defaultBuffer2 bytes.Buffer
		defaultBuffer2.WriteString("\treturn -1, false\n")

		getEnumAliasFunction.body = strings.ReplaceAll(getEnumAliasFunction.body, "__case__", caseBuffer2.String())
		getEnumAliasFunction.body = strings.ReplaceAll(getEnumAliasFunction.body, "__default__", defaultBuffer2.String())

		codeFunctions = append(codeFunctions, getEnumAliasFunction)

		// get enum string function
		getEnumStringFunction := &CodeFunction{
			name:       fmt.Sprintf("Get%s_EnumString", enumType.EnumType),
			parameters: []string{fmt.Sprintf("tp %s", enumType.EnumType)},
			retType:    "(string, bool)",
			body:       getEnumFunctionBody,
		}

		var caseBuffer3 bytes.Buffer
		for _, enum := range enumType.Enums {
			caseBuffer3.WriteString(fmt.Sprintf("\tcase %s_%s:\n\t\treturn \"%s\", true\n", enumType.EnumType, enum.EnumName, enum.EnumName))
		}

		var defaultBuffer3 bytes.Buffer
		defaultBuffer3.WriteString("\treturn \"\", false\n")

		getEnumStringFunction.body = strings.ReplaceAll(getEnumStringFunction.body, "__case__", caseBuffer3.String())
		getEnumStringFunction.body = strings.ReplaceAll(getEnumStringFunction.body, "__default__", defaultBuffer3.String())

		// codeFunctions = append(codeFunctions, getEnumStringFunction)
	}

	g := NewGoCodeGenerator(
		GoCodePackageName(packageName),
		GoCodeFilePath(filepath.Join(exportPath, metaName+"_enum.go")),
		GoCodeFileName(fileName),
		GoCodeImportPath([]string{}),
		GoCodeConstants(constants),
		GoCodeFunctions(codeFunctions),
	)

	err := g.Generate()
	if !utils.ErrCheck(err, "generate go code failed", g.opts.FilePath) {
		return err
	}

	return nil
}

// generate poly go code
func generatePolyGoCode(packageName, exportPath string, excelFileRaw *ExcelFileRaw) error {
	metaName := strings.Split(excelFileRaw.Filename, ".")[0]
	if utf8string.NewString(metaName).IsASCII() {
		for _, v := range metaName {
			metaName = string(unicode.ToLower(v)) + metaName[1:]
			break
		}
	}

	// init function
	initFunction := &CodeFunction{
		name: "init",
		// body: fmt.Sprintf("excel.AddEntryNewFunc(\"%s\", %s)", excelFileRaw.Filename, titleMetaName),
	}
	var initFunctionBodyBuilder strings.Builder

	// loot interface functions
	lootInterfaceFunctions := make([]*CodeFunction, 0, len(excelFileRaw.PolyData))

	// new object function
	newObjectFunctions := make([]*CodeFunction, 0, len(excelFileRaw.PolyData))

	polys := make([]*CodePoly, 0, 8)
	for _, polyType := range excelFileRaw.PolyData {
		poly := &CodePoly{
			tp:           polyType.PolyType,
			sep:          polyType.PolySeperator,
			parent:       polyType.PolyParentName,
			comment:      polyType.PolyComment,
			sheetName:    excelFileRaw.Sheetname,
			confFilename: excelFileRaw.ConfFilename,
			elems:        make([]*CodePolyElem, 0, len(polyType.Polys)),
		}

		for _, pl := range polyType.Polys {
			poly.elems = append(poly.elems, &CodePolyElem{
				name:    pl.PolyName,
				tp:      pl.PolyType,
				tag:     pl.PolyTag,
				comment: pl.PolyComment,
			})

		}

		polys = append(polys, poly)

		// init function body
		initBody := fmt.Sprintf("excel.AddEntryNewFunc(\"*%s\", new%s)", polyType.PolyType, polyType.PolyType)
		initFunctionBodyBuilder.WriteString(initBody)
		initFunctionBodyBuilder.WriteString("\n")

		// loot interface function
		// 包含掉落接口
		if polyType.PolyUnique == "loot" {
			lootInterfaceFunctions = append(lootInterfaceFunctions, &CodeFunction{
				name:        "GetId",
				receiver:    polyType.PolyType,
				receiverPtr: true,
				comment:     "掉落接口",
				retType:     polyType.PolyLootKeyType,
				body:        fmt.Sprintf("return %s(e.Id)", polyType.PolyLootKeyType),
			})

			lootInterfaceFunctions = append(lootInterfaceFunctions, &CodeFunction{
				name:        "GetWeight",
				receiver:    polyType.PolyType,
				receiverPtr: true,
				comment:     "掉落接口",
				retType:     "int32",
				body:        "return int32(e.Weight)",
			})
		}

		// new object function
		newObjectFunction := &CodeFunction{
			name:    fmt.Sprintf("new%s", polyType.PolyType),
			retType: "any",
			body:    fmt.Sprintf("return &%s{}", polyType.PolyType),
		}
		newObjectFunctions = append(newObjectFunctions, newObjectFunction)
	}

	initFunction.body = initFunctionBodyBuilder.String()
	codeFunctions := make([]*CodeFunction, 0, len(excelFileRaw.PolyData)+1)
	codeFunctions = append(codeFunctions, initFunction)
	codeFunctions = append(codeFunctions, lootInterfaceFunctions...)
	codeFunctions = append(codeFunctions, newObjectFunctions...)

	g := NewGoCodeGenerator(
		GoCodePackageName(packageName),
		GoCodeFilePath(filepath.Join(exportPath, metaName+"_poly.go")),
		GoCodeFileName(strings.Title(metaName)),
		GoCodeImportPath([]string{"github.com/east-eden/server/excel"}),
		GoCodePolys(polys),
		GoCodeFunctions(codeFunctions),
	)

	err := g.Generate()
	if !utils.ErrCheck(err, "generate go code failed", g.opts.FilePath) {
		return err
	}

	return nil
}

// generateEntryGoCode generates the contents of a .go file.
func generateEntryGoCode(packageName, exportPath string, excelFileRaw *ExcelFileRaw) error {
	// no need to export
	if !excelFileRaw.ShouldExport {
		return nil
	}

	metaName := strings.Split(excelFileRaw.Filename, ".")[0]

	if utf8string.NewString(metaName).IsASCII() {
		for _, v := range metaName {
			metaName = string(unicode.ToLower(v)) + metaName[1:]
			break
		}
	} else {
		// py := pinyin.LazyPinyin(metaName, pinyin.NewArgs())
		// metaName = strings.Join(py, "_")
	}

	titleMetaName := strings.Title(metaName)

	codeFunctions := make([]*CodeFunction, 0)

	// init function
	initFunction := &CodeFunction{
		name:       "init",
		parameters: []string{},
		body:       fmt.Sprintf("excel.AddEntryLoader(\"%s\", (*%sEntries)(nil))", excelFileRaw.Filename, titleMetaName),
	}

	// 动态配置特殊处理
	if packageName == "dynamic" {
		initFunction.body = fmt.Sprintf("excel.AddDynamicEntryLoader(\"%s\", (*%sEntries)(nil))", excelFileRaw.Filename, titleMetaName)
	}

	// load function
	loadFunction := &CodeFunction{
		receiver:    fmt.Sprintf("%sEntries", titleMetaName),
		receiverPtr: true,
		name:        "Load",
		parameters: []string{
			"excelFileRaw *excel.ExcelFileRaw",
		},
		retType: "(any, error)",
	}

	// single key
	var keyType string
	var keyName string
	if len(excelFileRaw.Keys) == 1 {
		loadFunction.body = singleKeyLoadFunctionBody
		keyName = excelFileRaw.Keys[0]

		if keyRaw, found := excelFileRaw.FieldRaw.Get(excelFileRaw.Keys[0]); found {
			keyRaw := keyRaw.(*ExcelFieldRaw)
			keyType, _ = getFullTypeName(keyRaw.tp)
		}
	} else {

		// multi key
		keyType = "string"
		keyName = excelFileRaw.Keys[0]
		loadFunction.body = func() string {
			keyName := make([]string, 0, len(excelFileRaw.Keys))
			keyValue := make([]string, 0, len(excelFileRaw.Keys))
			for _, key := range excelFileRaw.Keys {
				keyName = append(keyName, "%s")
				keyValue = append(keyValue, "cast.ToString(entry."+key+")")
			}

			finalKeyName := strings.Join(keyName, "_")
			finalKeyValue := strings.Join(keyValue, ", ")
			return fmt.Sprintf(multiKeyLoadFunctionBody, finalKeyName, finalKeyValue)
		}()
	}

	keyType, _ = getFullTypeName(keyType)
	loadFunction.body = strings.ReplaceAll(loadFunction.body, "__lowerReplace__", metaName)
	loadFunction.body = strings.ReplaceAll(loadFunction.body, "__upperReplace__", titleMetaName)
	loadFunction.body = strings.ReplaceAll(loadFunction.body, "__keytype__", keyType)
	loadFunction.body = strings.ReplaceAll(loadFunction.body, "__keyname__", keyName)

	// find all poly fields
	var polyReplaceBuilder strings.Builder
	it := excelFileRaw.FieldRaw.Iterator()
	for it.Next() {
		key := it.Key().(string)
		value := it.Value().(*ExcelFieldRaw)
		if value.tp != "any" {
			continue
		}

		// replace poly filed's any value to pointer value
		replaceOneString := fmt.Sprintf(`{
			pointerName := v["%s"].(map[string]any)["__pointer__type"]
			newFunc := excel.GetEntryNewFunc(pointerName.(string))
			if newFunc != nil {
				value := newFunc()
				err := mapstructure.Decode(v["%s"], value)
				if !utils.ErrCheck(err, "decode excel data to embeded struct failed", v) {
					return err
				}

				entry.%s = value
			}
			}`, key, key, key)
		polyReplaceBuilder.WriteString(replaceOneString)
		polyReplaceBuilder.WriteString("\n")
	}
	loadFunction.body = strings.ReplaceAll(loadFunction.body, "__Replace__Poly__String", polyReplaceBuilder.String())

	// GetRow function: single key GetRow and multi key GetRow
	singleKeyGetRowFn := &CodeFunction{
		name: fmt.Sprintf("Get%sEntry", titleMetaName),
		parameters: []string{
			fmt.Sprintf("id %s", keyType),
		},
		retType: fmt.Sprintf("(*%sEntry, bool)", titleMetaName),
		body:    fmt.Sprintf("entry, ok := %sEntries.Rows[id]\n\treturn entry, ok", metaName),
	}

	multiKeyGetRowFn := &CodeFunction{
		name: fmt.Sprintf("Get%sEntry", titleMetaName),
		parameters: []string{
			"keys ...string",
		},
		retType: fmt.Sprintf("(*%sEntry, bool)", titleMetaName),
		body: fmt.Sprintf(`keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := %sEntries.Rows[finalKey]
	return entry, ok `, metaName),
	}

	// one config
	if excelFileRaw.Type == ExcelSheetType_One {
		singleKeyGetRowFn.parameters = []string{}
		singleKeyGetRowFn.retType = fmt.Sprintf("(*%sEntry)", titleMetaName)
		singleKeyGetRowFn.body = fmt.Sprintf("entry := %sEntries.Rows[0]\n\treturn entry", metaName)
	}

	// GetSize function
	getSizeFunction := &CodeFunction{
		name:       fmt.Sprintf("Get%sSize", titleMetaName),
		parameters: []string{},
		retType:    "int",
		body:       fmt.Sprintf("return int(len(%sEntries.Rows))", metaName),
	}

	// GetRows function
	getRowsFunction := &CodeFunction{
		name:       fmt.Sprintf("Get%sRows", titleMetaName),
		parameters: []string{},
		body:       fmt.Sprintf("return %sEntries.Rows", metaName),
	}

	// receiver GetArrays function
	receiverGetArraysFunction := &CodeFunction{
		name:        "GetArrays",
		receiver:    fmt.Sprintf("%sEntries", titleMetaName),
		receiverPtr: true,
		parameters:  []string{},
		retType:     "any",
		body:        fmt.Sprintf("return %sEntries.Arrays", metaName),
	}

	// GetArrays function
	getArraysFunction := &CodeFunction{
		name:       fmt.Sprintf("Get%sArrays", titleMetaName),
		parameters: []string{},
		body:       fmt.Sprintf("return %sEntries.Arrays", metaName),
	}

	var getRowFunction *CodeFunction
	if len(excelFileRaw.Keys) == 1 {
		getRowFunction = singleKeyGetRowFn
		getRowsFunction.retType = fmt.Sprintf("map[%s]*%sEntry", keyType, titleMetaName)
		getArraysFunction.retType = fmt.Sprintf("[]*%sEntry", titleMetaName)
	} else {
		getRowFunction = multiKeyGetRowFn
		getRowsFunction.retType = fmt.Sprintf("map[%s]*%sEntry", keyType, titleMetaName)
		getArraysFunction.retType = fmt.Sprintf("[]*%sEntry", titleMetaName)
	}

	codeFunctions = append(codeFunctions, initFunction, loadFunction, receiverGetArraysFunction, getRowFunction, getSizeFunction, getRowsFunction, getArraysFunction)

	g := NewGoCodeGenerator(
		GoCodePackageName(packageName),
		GoCodeFileName(strings.Title(metaName)),
		GoCodeFilePath(filepath.Join(exportPath, metaName+"_entry.go")),

		GoCodeImportPath([]string{
			"sort",
			"container/list",
			"github.com/east-eden/server/excel",
			"github.com/east-eden/server/utils",
			"github.com/mitchellh/mapstructure",
			"github.com/rs/zerolog/log",
			"github.com/shopspring/decimal",
		}),

		GoCodeVariables([]*CodeVariable{
			{
				name:    fmt.Sprintf("%sEntries", metaName),
				tp:      fmt.Sprintf("*%sEntries", titleMetaName),
				comment: fmt.Sprintf("%s全局变量", excelFileRaw.Filename),
			},
		}),

		GoCodeFunctions(codeFunctions),
	)

	st := &CodeStruct{
		name:     fmt.Sprintf("%sEntry", titleMetaName),
		comment:  fmt.Sprintf("属性表 %s/%s", excelFileRaw.Sheetname, excelFileRaw.ConfFilename),
		fieldRaw: excelFileRaw.FieldRaw,
	}
	g.opts.Structs = append(g.opts.Structs, st)

	stRows := &CodeStruct{
		name:     fmt.Sprintf("%sEntries", titleMetaName),
		comment:  fmt.Sprintf("%s属性表集合", excelFileRaw.Filename),
		fieldRaw: treemap.NewWithStringComparator(),
	}

	// single key
	if len(excelFileRaw.Keys) == 1 {
		stRows.fieldRaw.Put("Rows", &ExcelFieldRaw{
			name:   "Rows",
			tp:     fmt.Sprintf("map[%s]*%sEntry", keyType, titleMetaName),
			tag:    "`json:\"rows,omitempty\" bson:\"rows\"`",
			export: true,
			idx:    0,
		})

		stRows.fieldRaw.Put("Arrays", &ExcelFieldRaw{
			name:   "Arrays",
			tp:     fmt.Sprintf("[]*%sEntry", titleMetaName),
			tag:    "`json:\"arrays,omitempty\" bson:\"arrays\"`",
			export: true,
			idx:    1,
		})
	} else {
		// multi key
		g.opts.ImportPath = append(g.opts.ImportPath, "github.com/spf13/cast", "fmt", "strings")
		stRows.fieldRaw.Put("Rows", &ExcelFieldRaw{
			name:   "Rows",
			tp:     fmt.Sprintf("map[%s]*%sEntry", keyType, titleMetaName),
			tag:    "`json:\"rows,omitempty\" bson:\"rows\"`",
			export: true,
			idx:    0,
		})

		stRows.fieldRaw.Put("Arrays", &ExcelFieldRaw{
			name:   "Arrays",
			tp:     fmt.Sprintf("[]*%sEntry", titleMetaName),
			tag:    "`json:\"arrays,omitempty\" bson:\"arrays\"`",
			export: true,
			idx:    1,
		})
	}

	// has map
	if excelFileRaw.HasMap {
		g.opts.ImportPath = append(g.opts.ImportPath, "github.com/emirpasic/gods/maps/treemap")
	}

	g.opts.Structs = append(g.opts.Structs, stRows)

	err := g.Generate()
	if !utils.ErrCheck(err, "generate go code failed", g.opts.FilePath) {
		return err
	}

	return nil
}
