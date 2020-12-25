package excel

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
)

type CodeEnumType int
type CodeEnumComment string

type CodeFieldType string
type CodeFieldTags string
type CodeFieldComment string

var defaultLoadFunctionBody string = `
	__lowerReplace__Entries = &__upperReplace__Entries{
		Rows: make(map[int]*__upperReplace__Entry),
	}

	for _, v := range excelFileRaw.cellData {
		entry := &__upperReplace__Entry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	__lowerReplace__Entries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.filename).Msg("excel load success")
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
		variableLine := fmt.Sprintf("var\t%-10s\t%-10s\t//%-10s", v.name, v.tp, v.comment)
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
			fieldLine := fmt.Sprintf("\t%-10s\t%-10s\t%-10s\t//%-10s", it.Key(), fieldRaw.tp, fieldRaw.tag, fieldRaw.desc)
			fieldLines[fieldRaw.idx] = fieldLine
		}

		// print struct field in sort
		for _, v := range fieldLines {
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
func generateCode(dirPath string, excelFileRaw *ExcelFileRaw) error {
	metaName := strings.Split(excelFileRaw.filename, ".")[0]
	titleMetaName := strings.Title(metaName)

	codeFunctions := make([]*CodeFunction, 0)

	// init function
	initFunction := &CodeFunction{
		name:       "init",
		parameters: []string{},
		body:       fmt.Sprintf("AddEntries(\"%s\", heroEntries)", excelFileRaw.filename),
	}

	// load function
	loadFunction := &CodeFunction{
		receiver: fmt.Sprintf("%sEntries", titleMetaName),
		name:     "load",
		parameters: []string{
			"excelFileRaw *ExcelFileRaw",
		},
		retType: "error",
	}
	loadFunction.body = defaultLoadFunctionBody
	loadFunction.body = strings.Replace(loadFunction.body, "__lowerReplace__", metaName, -1)
	loadFunction.body = strings.Replace(loadFunction.body, "__upperReplace__", titleMetaName, -1)

	// GetRow function
	getRowFunction := &CodeFunction{
		name: fmt.Sprintf("Get%sEntry", titleMetaName),
		parameters: []string{
			"id int",
		},
		retType: fmt.Sprintf("(*%sEntry, bool)", titleMetaName),
		body:    fmt.Sprintf("entry, ok := %sEntries.Rows[id]\n\treturn entry, ok", metaName),
	}

	codeFunctions = append(codeFunctions, initFunction, loadFunction, getRowFunction)

	g := NewCodeGenerator(
		CodePackageName("excel"),
		CodeFilePath(fmt.Sprintf("excel/%s_entry.go", metaName)),

		CodeImportPath([]string{
			"github.com/east-eden/server/utils",
			"github.com/mitchellh/mapstructure",
			"github.com/rs/zerolog/log",
		}),

		CodeVariables([]*CodeVariable{
			{
				name:    fmt.Sprintf("%sEntries", metaName),
				tp:      fmt.Sprintf("*%sEntries", titleMetaName),
				comment: fmt.Sprintf("%s全局变量", excelFileRaw.filename),
			},
		}),

		CodeFunctions(codeFunctions),
	)

	st := &CodeStruct{
		name:     fmt.Sprintf("%sEntry", titleMetaName),
		comment:  fmt.Sprintf("%s属性表", excelFileRaw.filename),
		fieldRaw: excelFileRaw.fieldRaw,
	}
	g.opts.Structs = append(g.opts.Structs, st)

	stRows := &CodeStruct{
		name:     fmt.Sprintf("%sEntries", titleMetaName),
		comment:  fmt.Sprintf("%s属性表集合", excelFileRaw.filename),
		fieldRaw: treemap.NewWithStringComparator(),
	}

	stRows.fieldRaw.Put("Rows", &ExcelFieldRaw{
		name: "Rows",
		tp:   fmt.Sprintf("map[int]*%sEntry", titleMetaName),
		tag:  "`json:\"Rows,omitempty\"`",
	})
	g.opts.Structs = append(g.opts.Structs, stRows)

	err := g.Generate()
	if utils.ErrCheck(err, "generate go code failed", g.opts.FilePath) {
		return err
	}

	return nil
}
