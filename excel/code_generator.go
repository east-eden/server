package excel

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
)

type CodeEnumType int
type CodeEnumComment string

type CodeFieldType string
type CodeFieldTags string
type CodeFieldComment string

// single key load function
var singleKeyLoadFunctionBody string = `
	__lowerReplace__Entries = &__upperReplace__Entries{
		Rows: make(map[int32]*__upperReplace__Entry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &__upperReplace__Entry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	__lowerReplace__Entries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	`

// multi key load function
var multiKeyLoadFunctionBody string = `
	__lowerReplace__Entries = &__upperReplace__Entries{
		Rows: make(map[string]*__upperReplace__Entry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &__upperReplace__Entry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		key := fmt.Sprintf("%s", %s)
	 	__lowerReplace__Entries.Rows[key] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	`

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

// code function
type CodeFunction struct {
	receiver   string
	name       string
	comment    string
	parameters []string
	retType    string
	body       string
}

type CodeGeneratorOptions struct {
	PackageName string // name of this file's Go package
	FilePath    string
	ImportPath  []string // import path of this file's Go package

	Structs   []*CodeStruct
	Variables []*CodeVariable
	Functions []*CodeFunction

	FileNamePrefix string
	FileNameSuffix string

	Comments string
}

type CodeGeneratorOption func(*CodeGeneratorOptions)

// code generator
type CodeGenerator struct {
	opts *CodeGeneratorOptions
	buf  bytes.Buffer
}

func defaultOptions() *CodeGeneratorOptions {
	return &CodeGeneratorOptions{
		ImportPath: []string{},
		Structs:    make([]*CodeStruct, 0),
		Variables:  make([]*CodeVariable, 0),
		Functions:  make([]*CodeFunction, 0),
	}
}

func CodePackageName(packageName string) CodeGeneratorOption {
	return func(o *CodeGeneratorOptions) {
		o.PackageName = packageName
	}
}

func CodeFilePath(filePath string) CodeGeneratorOption {
	return func(o *CodeGeneratorOptions) {
		o.FilePath = filePath
	}
}

func CodeImportPath(importPath []string) CodeGeneratorOption {
	return func(o *CodeGeneratorOptions) {
		o.ImportPath = importPath
	}
}

func CodeStructs(structs []*CodeStruct) CodeGeneratorOption {
	return func(o *CodeGeneratorOptions) {
		o.Structs = structs
	}
}

func CodeVariables(variables []*CodeVariable) CodeGeneratorOption {
	return func(o *CodeGeneratorOptions) {
		o.Variables = variables
	}
}

func CodeFunctions(functions []*CodeFunction) CodeGeneratorOption {
	return func(o *CodeGeneratorOptions) {
		o.Functions = functions
	}
}

func (g *CodeGenerator) P(v ...interface{}) {
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
func (g *CodeGenerator) Generate() error {
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
			if !fieldRaw.imp {
				continue
			}

			fieldLine := fmt.Sprintf("\t%-15s\t%-20s\t%-20s\t//%-10s", it.Key(), fieldRaw.tp, fieldRaw.tag, fieldRaw.desc)
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
			receiver = fmt.Sprintf("(e *%s)", f.receiver)
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

	return ioutil.WriteFile(g.opts.FilePath, g.buf.Bytes(), 0666)
}

func NewCodeGenerator(options ...CodeGeneratorOption) *CodeGenerator {
	g := &CodeGenerator{
		opts: defaultOptions(),
	}

	for _, o := range options {
		o(g.opts)
	}

	return g
}

// generateCode generates the contents of a .go file.
func generateCode(exportPath string, excelFileRaw *ExcelFileRaw) error {
	metaName := strings.Split(excelFileRaw.Filename, ".")[0]
	for _, v := range metaName {
		metaName = string(unicode.ToLower(v)) + metaName[1:]
		break
	}

	titleMetaName := strings.Title(metaName)

	codeFunctions := make([]*CodeFunction, 0)

	// init function
	initFunction := &CodeFunction{
		name:       "init",
		parameters: []string{},
		body:       fmt.Sprintf("excel.AddEntryLoader(\"%s\", (*%sEntries)(nil))", excelFileRaw.Filename, titleMetaName),
	}

	// load function
	loadFunction := &CodeFunction{
		receiver: fmt.Sprintf("%sEntries", titleMetaName),
		name:     "Load",
		parameters: []string{
			"excelFileRaw *excel.ExcelFileRaw",
		},
		retType: "error",
	}

	// single key
	if len(excelFileRaw.Keys) == 1 {
		loadFunction.body = singleKeyLoadFunctionBody
	} else {

		// multi key
		loadFunction.body = func() string {
			keyName := make([]string, 0, len(excelFileRaw.Keys))
			keyValue := make([]string, 0, len(excelFileRaw.Keys))
			for _, key := range excelFileRaw.Keys {
				keyName = append(keyName, "%d")
				keyValue = append(keyValue, "entry."+key)
			}

			finalKeyName := strings.Join(keyName, "+")
			finalKeyValue := strings.Join(keyValue, ", ")
			return fmt.Sprintf(multiKeyLoadFunctionBody, finalKeyName, finalKeyValue)
		}()
	}
	loadFunction.body = strings.Replace(loadFunction.body, "__lowerReplace__", metaName, -1)
	loadFunction.body = strings.Replace(loadFunction.body, "__upperReplace__", titleMetaName, -1)

	// GetRow function: single key GetRow and multi key GetRow
	singleKeyGetRowFn := &CodeFunction{
		name: fmt.Sprintf("Get%sEntry", titleMetaName),
		parameters: []string{
			"id int32",
		},
		retType: fmt.Sprintf("(*%sEntry, bool)", titleMetaName),
		body:    fmt.Sprintf("entry, ok := %sEntries.Rows[id]\n\treturn entry, ok", metaName),
	}

	multiKeyGetRowFn := func() *CodeFunction {
		fn := &CodeFunction{
			name: fmt.Sprintf("Get%sEntry", titleMetaName),
			parameters: []string{
				"keys ...int32",
			},
			retType: fmt.Sprintf("(*%sEntry, bool)", titleMetaName),
		}

		fn.body = fmt.Sprintf(`keyName := make([]string, 0, len(keys))
	for _, key := range keys {
		keyName = append(keyName, cast.ToString(key))
	}

	finalKey := strings.Join(keyName, "+")
	entry, ok := %sEntries.Rows[finalKey]
	return entry, ok `, metaName)

		return fn
	}()

	// GetSize function
	getSizeFunction := &CodeFunction{
		name:       fmt.Sprintf("Get%sSize", titleMetaName),
		parameters: []string{},
		retType:    "int32",
		body:       fmt.Sprintf("return int32(len(%sEntries.Rows))", metaName),
	}

	// GetRows function
	getRowsFunction := &CodeFunction{
		name:       fmt.Sprintf("Get%sRows", titleMetaName),
		parameters: []string{},
		body:       fmt.Sprintf("return %sEntries.Rows", metaName),
	}

	var getRowFunction *CodeFunction
	if len(excelFileRaw.Keys) == 1 {
		getRowFunction = singleKeyGetRowFn
		getRowsFunction.retType = fmt.Sprintf("map[int32]*%sEntry", titleMetaName)
	} else {
		getRowFunction = multiKeyGetRowFn
		getRowsFunction.retType = fmt.Sprintf("map[string]*%sEntry", titleMetaName)
	}

	codeFunctions = append(codeFunctions, initFunction, loadFunction, getRowFunction, getSizeFunction, getRowsFunction)

	g := NewCodeGenerator(
		CodePackageName("auto"),
		CodeFilePath(fmt.Sprintf("%s/%s_entry.go", exportPath, metaName)),

		CodeImportPath([]string{
			"github.com/east-eden/server/excel",
			"github.com/east-eden/server/utils",
			"github.com/mitchellh/mapstructure",
			"github.com/rs/zerolog/log",
			"github.com/shopspring/decimal",
		}),

		CodeVariables([]*CodeVariable{
			{
				name:    fmt.Sprintf("%sEntries", metaName),
				tp:      fmt.Sprintf("*%sEntries", titleMetaName),
				comment: fmt.Sprintf("%s全局变量", excelFileRaw.Filename),
			},
		}),

		CodeFunctions(codeFunctions),
	)

	st := &CodeStruct{
		name:     fmt.Sprintf("%sEntry", titleMetaName),
		comment:  fmt.Sprintf("%s属性表", excelFileRaw.Filename),
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
			name: "Rows",
			tp:   fmt.Sprintf("map[int32]*%sEntry", titleMetaName),
			tag:  "`json:\"Rows,omitempty\"`",
			imp:  true,
		})
	} else {
		// multi key
		g.opts.ImportPath = append(g.opts.ImportPath, "github.com/spf13/cast", "fmt", "strings")
		stRows.fieldRaw.Put("Rows", &ExcelFieldRaw{
			name: "Rows",
			tp:   fmt.Sprintf("map[string]*%sEntry", titleMetaName),
			tag:  "`json:\"Rows,omitempty\"`",
			imp:  true,
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
