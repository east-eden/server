package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
)

type CodeEnumType int
type CodeEnumComment string

type CodeFieldType string
type CodeFieldTags string
type CodeFieldComment string

// code enum
type CodeEnum struct {
	name          string
	enums         map[string]CodeEnumType
	enumsComments map[string]CodeEnumComment
}

// code struct
type CodeStruct struct {
	name           string
	comment        string
	fields         *treemap.Map
	fieldsTags     *treemap.Map
	fieldsComments *treemap.Map
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

	Enums     []*CodeEnum
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
		Enums:      make([]*CodeEnum, 0),
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
		// g.P("var\t", v.name, "\t", v.tp, "\t//", v.comment)
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
		it := s.fields.Iterator()
		for i := 0; it.Next(); i++ {
			fieldTag, okTag := s.fieldsTags.Get(it.Key())
			if !okTag {
				fieldTag = ""
			}

			fieldComment, okComment := s.fieldsComments.Get(it.Key())
			if !okComment {
				fieldComment = ""
			}

			fieldLine := fmt.Sprintf("\t%-10s\t%-10s\t%-10s\t//%-10s", it.Key(), it.Value(), fieldTag, fieldComment)
			g.P(fieldLine)
			// g.P("\t", it.Key(), "\t", it.Value(), "\t", fieldTag, "\t//", fieldComment)
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
		var parameters string
		for _, v := range f.parameters {
			parameters = strings.Join([]string{parameters, v}, ", ")
		}

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

// GenerateFile generates the contents of a .pb.go file.
func GenerateFile() error {
	g := NewCodeGenerator(
		CodePackageName("excel"),
		CodeFilePath("excel/hero_entry.go"),

		CodeImportPath([]string{
			"strconv",
			"strings",
			"github.com/360EntSecGroup-Skylar/excelize/v2",
			"github.com/east-eden/server/utils",
			"github.com/mitchellh/mapstructure",
			"github.com/rs/zerolog/log",
		}),

		CodeVariables([]*CodeVariable{
			{
				name:    "heroEntries",
				tp:      "*HeroEntries",
				comment: "英雄属性表全局变量",
			},
		}),

		CodeFunctions([]*CodeFunction{
			{
				name:       "init",
				parameters: []string{},
				body:       "AddEntries(heroEntries, \"HeroConfig.xlsx\")",
			},

			{
				receiver: "HeroEntries",
				name:     "Load",
				retType:  "error",
				body:     "return nil",
			},
		}),
	)

	st := &CodeStruct{
		name:           "HeroEntry",
		comment:        "英雄属性表",
		fields:         treemap.NewWithStringComparator(),
		fieldsTags:     treemap.NewWithStringComparator(),
		fieldsComments: treemap.NewWithStringComparator(),
	}
	st.fields.Put("ID", "int")
	st.fields.Put("Name", "string")
	st.fields.Put("AttID", "int")
	st.fields.Put("Quality", "int")
	st.fields.Put("AttList", "[]int")

	st.fieldsTags.Put("ID", "`json:\"Id\"`")
	st.fieldsTags.Put("Name", "`json:\"Name,omitempty\"`")
	st.fieldsTags.Put("AttID", "`json:\"AttID,omitempty\"`")
	st.fieldsTags.Put("Quality", "`json:\"Quality,omitempty\"`")
	st.fieldsTags.Put("AttList", "`json:\"AttList,omitempty\"`")

	st.fieldsComments.Put("ID", "id")
	st.fieldsComments.Put("Name", "名字")
	st.fieldsComments.Put("AttID", "属性id")
	st.fieldsComments.Put("Quality", "品质")
	st.fieldsComments.Put("AttList", "属性列表")
	g.opts.Structs = append(g.opts.Structs, st)

	stRows := &CodeStruct{
		name:           "HeroEntries",
		comment:        "英雄属性表集合",
		fields:         treemap.NewWithStringComparator(),
		fieldsTags:     treemap.NewWithStringComparator(),
		fieldsComments: treemap.NewWithStringComparator(),
	}
	stRows.fields.Put("Rows", "map[int]*HeroEntry")
	stRows.fieldsTags.Put("Rows", "`json:\"Rows\"`")
	g.opts.Structs = append(g.opts.Structs, stRows)

	err := g.Generate()
	if ErrCheck(err, "generate go code failed", g.opts.FilePath) {
		return err
	}

	return nil
}
