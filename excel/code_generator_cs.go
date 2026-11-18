package excel

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/spf13/cast"
)

var (
	parsingCsFileRecord sync.Map // 正在处理的文件记录
)

type CsCodeData struct{}

type CsCodeGeneratorOptions struct {
	FilePath          string
	FileName          string
	RequiredFilenames []string
	EnumRaws          []*ExcelEnumRaw
	CsData            *CsCodeData
}

type CsCodeGeneratorOption func(*CsCodeGeneratorOptions)

type CsCodeGenerator struct {
	opts *CsCodeGeneratorOptions
	buf  bytes.Buffer
}

func defaultCsOptions() *CsCodeGeneratorOptions {
	return &CsCodeGeneratorOptions{
		RequiredFilenames: []string{},
		EnumRaws:          make([]*ExcelEnumRaw, 0),
	}
}

func CsCodeDatas(data *CsCodeData) CsCodeGeneratorOption {
	return func(o *CsCodeGeneratorOptions) {
		o.CsData = data
	}
}

func CsCodeFilePath(filePath string) CsCodeGeneratorOption {
	return func(o *CsCodeGeneratorOptions) {
		o.FilePath = filePath
	}
}

func CsCodeFileName(fileName string) CsCodeGeneratorOption {
	return func(o *CsCodeGeneratorOptions) {
		o.FileName = fileName
	}
}

func CsCodeRequiredFileNames(fileNames []string) CsCodeGeneratorOption {
	return func(o *CsCodeGeneratorOptions) {
		o.RequiredFilenames = make([]string, 0, len(fileNames))
		o.RequiredFilenames = append(o.RequiredFilenames, fileNames...)
	}
}

func CsCodeEnumRaws(raws []*ExcelEnumRaw) CsCodeGeneratorOption {
	return func(o *CsCodeGeneratorOptions) {
		o.EnumRaws = make([]*ExcelEnumRaw, 0, len(raws))
		o.EnumRaws = append(o.EnumRaws, raws...)
	}
}

func toCsType(tp string) string {
	switch tp {
	case "int", "int32":
		return "int"
	case "int64":
		return "long"
	case "uint", "uint32":
		return "uint"
	case "uint64":
		return "ulong"
	case "float", "float32":
		return "float"
	case "float64", "double":
		return "double"
	case "bool", "BOOL", "boolean", "Boolean":
		return "bool"
	case "string", "String":
		return "string"
	default:
		if strings.HasPrefix(tp, "list2:") {
			subType := strings.TrimPrefix(tp, "list2:")
			return fmt.Sprintf("List<List<%s>>", toCsType(subType))
		}

		if strings.HasPrefix(tp, "list:") {
			subType := strings.TrimPrefix(tp, "list:")
			return fmt.Sprintf("List<%s>", toCsType(subType))
		}

		if strings.HasPrefix(tp, "[]") {
			subType := strings.TrimPrefix(tp, "[]")
			return fmt.Sprintf("List<%s>", toCsType(subType))
		}

		if strings.HasPrefix(tp, "map[") {
			t := strings.Split(tp, "[")
			tt := strings.Split(t[len(t)-1], "]")
			keyType := toCsType(tt[0])
			valueType := toCsType(tt[1])
			return fmt.Sprintf("Dictionary<%s, %s>", keyType, valueType)
		}

		trimPtrType := strings.TrimPrefix(tp, "*")
		if isPolyTypeName(trimPtrType) {
			return fmt.Sprintf("Conf%s", utils.SimpleTitle(trimPtrType))
		}

		if _, found := excelEnumTypeRaws.Load(utils.SimpleTitle(trimPtrType)); found {
			return fmt.Sprintf("E%s", utils.SimpleTitle(trimPtrType))
		}

		return fmt.Sprintf("Conf%s", utils.SimpleTitle(trimPtrType))
	}
}

func csParseExpr(originTp, jsonExpr string) string {
	csType := toCsType(originTp)
	trimType := strings.TrimPrefix(originTp, "*")
	if _, found := excelEnumTypeRaws.Load(utils.SimpleTitle(trimType)); found {
		return fmt.Sprintf("(%s)int.Parse(%s.ToString())", csType, jsonExpr)
	}

	switch csType {
	case "string":
		return fmt.Sprintf("%s.ToString()", jsonExpr)
	case "int":
		return fmt.Sprintf("int.Parse(%s.ToString())", jsonExpr)
	case "long":
		return fmt.Sprintf("long.Parse(%s.ToString())", jsonExpr)
	case "uint":
		return fmt.Sprintf("uint.Parse(%s.ToString())", jsonExpr)
	case "ulong":
		return fmt.Sprintf("ulong.Parse(%s.ToString())", jsonExpr)
	case "float":
		return fmt.Sprintf("float.Parse(%s.ToString())", jsonExpr)
	case "double":
		return fmt.Sprintf("double.Parse(%s.ToString())", jsonExpr)
	case "bool":
		return fmt.Sprintf("bool.Parse(%s.ToString())", jsonExpr)
	default:
		return fmt.Sprintf("JsonMapper.ToObject<%s>(%s.ToJson())", csType, jsonExpr)
	}
}

func csKeyToStringExpr(originTp, valueExpr string) string {
	trimType := strings.TrimPrefix(originTp, "*")
	if _, found := excelEnumTypeRaws.Load(utils.SimpleTitle(trimType)); found {
		return fmt.Sprintf("((int)%s).ToString()", valueExpr)
	}

	switch toCsType(originTp) {
	case "string":
		return valueExpr
	case "float":
		return fmt.Sprintf("%s.ToString(System.Globalization.CultureInfo.InvariantCulture)", valueExpr)
	case "double":
		return fmt.Sprintf("%s.ToString(System.Globalization.CultureInfo.InvariantCulture)", valueExpr)
	case "bool":
		return fmt.Sprintf("%s.ToString().ToLower()", valueExpr)
	default:
		return fmt.Sprintf("%s.ToString()", valueExpr)
	}
}

func shouldExportEnumInfo(enumRaw *ExcelEnumRaw) bool {
	return enumRaw != nil && enumRaw.HasOption("export_info")
}

func hasExportableEnumInfo(enumRaws []*ExcelEnumRaw) bool {
	for _, enumRaw := range enumRaws {
		if shouldExportEnumInfo(enumRaw) {
			return true
		}
	}
	return false
}

func shouldExportPolyClassName(polyType string) bool {
	rawValue, found := excelPolyTypeRaws.Load(polyType)
	if !found {
		return false
	}

	raw := rawValue.(*ExcelPolyType)
	if raw.HasOption("export_class_name") {
		return true
	}
	if len(raw.PolyParentName) == 0 {
		return false
	}
	return shouldExportPolyClassName(raw.PolyParentName)
}

func shouldDeclarePolyClassName(raw *ExcelPolyType) bool {
	if raw == nil || !raw.HasOption("export_class_name") {
		return false
	}
	if len(raw.PolyParentName) == 0 {
		return true
	}
	return !shouldExportPolyClassName(raw.PolyParentName)
}

func buildCsJsonValue(filename, tp string, value any, jsonKey bool) any {
	polyDealFunc := func(polyType string) any {
		if value == nil {
			value = map[string]any{}
		}

		mapValue, ok := value.(map[string]any)
		if !ok {
			if record, found := parsingCsFileRecord.Load(filename); found {
				info := record.(map[string]string)
				GlobalResultReport.AddError(filename, "", info["key_name"], fmt.Sprintf("多态字段`%s` 类型转换失败,请检查该字段填写是否有误", polyType))
			}
			mapValue = map[string]any{}
		}

		classType := polyType
		if pointerType, found := mapValue["__pointer__type"]; found {
			classType = strings.TrimPrefix(cast.ToString(pointerType), "*")
		}

		result := make(map[string]any)
		if shouldExportPolyClassName(classType) {
			result["ClassName"] = classType
		}

		if rawValue, found := excelPolyTypeRaws.Load(classType); found {
			raw := rawValue.(*ExcelPolyType)
			for _, polyData := range raw.Polys {
				result[utils.SimpleTitle(polyData.PolyName)] = buildCsJsonValue(filename, polyData.PolyType, mapValue[polyData.PolyName], false)
			}
		}

		return result
	}

	switch tp {
	case "string":
		if value == nil {
			value = ""
		}
		return cast.ToString(value)
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		if value == nil {
			return 0
		}
		return cast.ToInt(value)
	case "float", "float32":
		if value == nil {
			return float32(0)
		}
		return cast.ToFloat32(value)
	case "float64", "double":
		if value == nil {
			return float64(0)
		}
		return cast.ToFloat64(value)
	case "bool", "BOOL", "boolean", "Boolean":
		if value == nil {
			return false
		}
		return cast.ToBool(value)
	case "any":
		return polyDealFunc(tp)
	default:
		if _, found := excelEnumTypeRaws.Load(utils.SimpleTitle(tp)); found {
			if value == nil {
				return -1
			}
			return cast.ToInt(value)
		}

		polyTp := utils.SimpleTitle(strings.TrimPrefix(tp, "*"))
		if isPolyTypeName(polyTp) {
			return polyDealFunc(polyTp)
		}

		if strings.HasPrefix(tp, "list2:") {
			subType := strings.TrimPrefix(tp, "list2:")
			if value == nil {
				return []any{}
			}
			values := value.([]any)
			result := make([]any, 0, len(values))
			nextType := strings.Join([]string{"list", subType}, ":")
			for _, v := range values {
				result = append(result, buildCsJsonValue(filename, nextType, v, false))
			}
			return result
		}

		if strings.HasPrefix(tp, "list:") {
			subType := strings.TrimPrefix(tp, "list:")
			if value == nil {
				return []any{}
			}
			values := value.([]any)
			result := make([]any, 0, len(values))
			for _, v := range values {
				result = append(result, buildCsJsonValue(filename, subType, v, false))
			}
			return result
		}

		if strings.HasPrefix(tp, "map[") {
			if value == nil {
				return map[string]any{}
			}

			t := strings.Split(tp, "[")
			tt := strings.Split(t[len(t)-1], "]")
			result := make(map[string]any)
			values, ok := value.(map[any]any)
			if !ok {
				GlobalResultReport.AddError(filename, "", tp, fmt.Sprintf("map类型`%s` 类型转换失败,请检查该字段填写是否有误", tp))
				return result
			}

			for k, v := range values {
				key := cast.ToString(buildCsJsonValue(filename, tt[0], k, true))
				result[key] = buildCsJsonValue(filename, tt[1], v, false)
			}
			return result
		}

		if strings.HasPrefix(tp, "[]") {
			subType := strings.TrimPrefix(tp, "[]")
			if value == nil {
				return []any{}
			}
			values := value.([]any)
			result := make([]any, 0, len(values))
			for _, v := range values {
				result = append(result, buildCsJsonValue(filename, subType, v, false))
			}
			return result
		}

		if value == nil {
			if jsonKey {
				return ""
			}
			return ""
		}

		return value
	}
}

func buildCsRowObject(excelRaw *ExcelFileRaw, row ExcelRowData) map[string]any {
	result := make(map[string]any)
	fieldTree := treemap.NewWithIntComparator()

	excelRaw.FieldRaw.All(func(_, v any) bool {
		fieldRaw := v.(*ExcelFieldRaw)
		if !fieldRaw.export {
			return true
		}
		fieldTree.Put(fieldRaw.idx, fieldRaw)
		return true
	})

	fieldTree.All(func(_, v any) bool {
		fieldRaw := v.(*ExcelFieldRaw)
		result[utils.SimpleTitle(fieldRaw.originName)] = buildCsJsonValue(excelRaw.Filename, fieldRaw.originTp, row[fieldRaw.name], false)
		return true
	})

	return result
}

func buildCsOneAndMapObject(excelRaw *ExcelFileRaw) any {
	if excelRaw.Type == ExcelSheetType_One {
		for _, row := range excelRaw.RowData {
			if len(row) == 0 {
				continue
			}
			return buildCsRowObject(excelRaw, row)
		}
		return map[string]any{}
	}

	result := make(map[string]any)
	for _, row := range excelRaw.RowData {
		if len(row) == 0 || len(excelRaw.Keys) == 0 {
			continue
		}

		keyName := excelRaw.Keys[0]
		keyData := row[keyName]
		keyType, _ := excelRaw.FieldRaw.Get(keyName)
		keyOriginType := keyType.(*ExcelFieldRaw).originTp
		multiKey := len(excelRaw.Keys) > 1

		parsingCsFileRecord.Store(excelRaw.Filename, map[string]string{
			"key_name": cast.ToString(keyName),
			"key_data": cast.ToString(keyData),
		})

		var rowKey string
		if multiKey {
			multiKeyValue := make([]string, 0, len(excelRaw.Keys))
			for _, currentKeyName := range excelRaw.Keys {
				multiKeyValue = append(multiKeyValue, cast.ToString(row[currentKeyName]))
			}
			rowKey = strings.Join(multiKeyValue, "_")
		} else {
			rowKey = cast.ToString(buildCsJsonValue(excelRaw.Filename, keyOriginType, keyData, true))
		}

		result[rowKey] = buildCsRowObject(excelRaw, row)
	}

	return result
}

func buildCsEnumObjects(excelRaw *ExcelFileRaw) map[string]any {
	result := make(map[string]any)
	for _, excelEnumType := range excelRaw.EnumData {
		rawValue, found := excelEnumTypeRaws.Load(excelEnumType.EnumType)
		if !found || !shouldExportEnumInfo(rawValue.(*ExcelEnumRaw)) {
			continue
		}
		enumInfoName := utils.SimpleTitle(excelEnumType.EnumType) + "EnumInfo"
		enumInfo := make(map[string]any)
		for _, row := range excelEnumType.Enums {
			enumInfo[row.EnumValue] = map[string]any{
				"Id":   cast.ToInt(row.EnumValue),
				"Key":  row.EnumName,
				"Name": row.EnumDesc,
			}
		}
		result[enumInfoName] = enumInfo
	}
	return result
}

func buildCsRootObject() map[string]any {
	root := make(map[string]any)
	excelFileRaws.Range(func(k, v any) bool {
		name := k.(string)
		excelRaw := v.(*ExcelFileRaw)
		if excelRaw.Type == ExcelSheetType_Poly || strings.HasSuffix(excelRaw.Sheetname, "client_skip") {
			return true
		}

		if excelRaw.Type == ExcelSheetType_Enum {
			for enumName, enumData := range buildCsEnumObjects(excelRaw) {
				root[enumName] = enumData
			}
			return true
		}

		if excelRaw.Type == ExcelSheetType_One || excelRaw.Type == ExcelSheetType_Map {
			root[utils.SimpleTitle(strings.TrimSuffix(name, ".json"))] = buildCsOneAndMapObject(excelRaw)
		}
		return true
	})
	return root
}

func (g *CsCodeGenerator) P(v ...any) {
	for _, x := range v {
		fmt.Fprint(&g.buf, x)
	}
	fmt.Fprintln(&g.buf)
}

func (g *CsCodeGenerator) writeFields(fieldRaws []*ExcelFieldRaw) {
	for _, fieldRaw := range fieldRaws {
		line := fmt.Sprintf("        public %s %s;", toCsType(fieldRaw.originTp), utils.SimpleTitle(fieldRaw.originName))
		if len(fieldRaw.desc) > 0 {
			line += fmt.Sprintf(" // %s", fieldRaw.desc)
		}
		g.P(line)
	}
}

func (g *CsCodeGenerator) writeParseMethod(fieldRaws []*ExcelFieldRaw) {
	g.P("")
	g.P("        public void Parse(JsonData tableNode)")
	g.P("        {")
	for _, fieldRaw := range fieldRaws {
		fieldName := utils.SimpleTitle(fieldRaw.originName)
		g.P(fmt.Sprintf("            %s = %s;", fieldName, csParseExpr(fieldRaw.originTp, fmt.Sprintf("tableNode[\"%s\"]", fieldName))))
	}
	g.P("        }")
}

func getSortedExportFieldRaws(excelFileRaw *ExcelFileRaw) []*ExcelFieldRaw {
	fieldTree := treemap.NewWithIntComparator()
	excelFileRaw.FieldRaw.All(func(_, v any) bool {
		fieldRaw := v.(*ExcelFieldRaw)
		if !fieldRaw.export {
			return true
		}
		fieldTree.Put(fieldRaw.idx, fieldRaw)
		return true
	})

	result := make([]*ExcelFieldRaw, 0, fieldTree.Size())
	fieldTree.All(func(_, v any) bool {
		result = append(result, v.(*ExcelFieldRaw))
		return true
	})
	return result
}

func getMultiKeyFieldRaws(excelFileRaw *ExcelFileRaw) []*ExcelFieldRaw {
	result := make([]*ExcelFieldRaw, 0, len(excelFileRaw.Keys))
	for _, keyName := range excelFileRaw.Keys {
		value, found := excelFileRaw.FieldRaw.Get(keyName)
		if !found {
			continue
		}
		result = append(result, value.(*ExcelFieldRaw))
	}
	return result
}

func shouldGenerateTableWrapper(excelFileRaw *ExcelFileRaw) bool {
	if excelFileRaw == nil {
		return false
	}

	// 多主键一定需要 wrapper：ConfigTable<...> 的 Get 不能直接支持多参数 key
	if len(excelFileRaw.Keys) > 1 {
		return true
	}

	// 单主键：主键类型不是 string 时，需要 wrapper 以便传入类型化的 key
	if len(excelFileRaw.Keys) != 1 {
		return false
	}

	value, found := excelFileRaw.FieldRaw.Get(excelFileRaw.Keys[0])
	if !found {
		return false
	}

	fieldRaw := value.(*ExcelFieldRaw)
	return toCsType(fieldRaw.originTp) != "string"
}

func (g *CsCodeGenerator) generateCode() error {
	g.P(`/***************************************
    This File is auto generated.
    Do on do any changes on the file !!!
    Any change will be override !!!
***************************************/`)
	g.P("using System.Collections.Generic;")
	g.P("using LitJson;")
	g.P("")
	g.P("namespace Game")
	g.P("{")
	if hasExportableEnumInfo(g.opts.EnumRaws) {
		g.P("    public class ConfEnumInfo")
		g.P("    {")
		g.P("        public int Id; // 枚举id")
		g.P("        public string Key; // 枚举key")
		g.P("        public string Name; // 枚举名称")
		g.P("")
		g.P("        public void Parse(JsonData tableNode)")
		g.P("        {")
		g.P("            Id = int.Parse(tableNode[\"Id\"].ToString());")
		g.P("            Key = tableNode[\"Key\"].ToString();")
		g.P("            Name = tableNode[\"Name\"].ToString();")
		g.P("        }")
		g.P("    }")
		g.P("")
	}

	requiredFilenameTree := treemap.NewWithStringComparator()
	for _, filename := range g.opts.RequiredFilenames {
		requiredFilenameTree.Put(filename, true)
	}

	requiredFilenameTree.All(func(k, _ any) bool {
		filename := k.(string)
		value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
		if !found {
			return true
		}

		excelFileRaw := value.(*ExcelFileRaw)
		if strings.HasSuffix(excelFileRaw.Sheetname, "client_skip") {
			return true
		}

		confName := fmt.Sprintf("Conf%s", utils.SimpleTitle(filename))
		classHeader := fmt.Sprintf("    public class %s", confName)
		if excelFileRaw.Type == ExcelSheetType_One {
			classHeader += " : IConfigDict"
		}
		classHeader += fmt.Sprintf(" // %s %s", excelFileRaw.Sheetname, excelFileRaw.ConfFilename)
		g.P(classHeader)
		g.P("    {")
		fieldRaws := getSortedExportFieldRaws(excelFileRaw)
		g.writeFields(fieldRaws)
		if excelFileRaw.Type == ExcelSheetType_One {
			g.writeParseMethod(fieldRaws)
		}
		g.P("    }")
		g.P("")
		return true
	})

	polyFilenameTree := treemap.NewWithStringComparator()
	excelPolyTypeRaws.Range(func(k, v any) bool {
		polyFilenameTree.Put(k.(string), v.(*ExcelPolyType))
		return true
	})

	polyFilenameTree.All(func(k, v any) bool {
		polyName := k.(string)
		excelPolyType := v.(*ExcelPolyType)
		value, found := excelFileRaws.Load(excelPolyType.PolyFilename)
		if !found {
			return true
		}
		excelFileRaw := value.(*ExcelFileRaw)
		confName := fmt.Sprintf("Conf%s", utils.SimpleTitle(polyName))
		classHeader := fmt.Sprintf("    public class %s", confName)
		if len(excelPolyType.PolyParentName) > 0 {
			classHeader += fmt.Sprintf(" : Conf%s", utils.SimpleTitle(excelPolyType.PolyParentName))
		}
		classHeader += fmt.Sprintf(" // %s %s", excelFileRaw.Sheetname, excelFileRaw.ConfFilename)
		g.P(classHeader)
		g.P("    {")
		if shouldDeclarePolyClassName(excelPolyType) {
			g.P("        public string ClassName; // 多态类型名")
		}
		fieldRaws := make([]*ExcelFieldRaw, 0, len(excelPolyType.OriginPolys))
		for idx, poly := range excelPolyType.OriginPolys {
			fieldRaws = append(fieldRaws, &ExcelFieldRaw{
				originName: poly.PolyName,
				originTp:   poly.PolyType,
				desc:       poly.PolyComment,
				idx:        idx,
				export:     true,
			})
		}
		g.writeFields(fieldRaws)
		g.P("    }")
		g.P("")
		return true
	})

	requiredFilenameTree.All(func(k, _ any) bool {
		filename := k.(string)
		value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
		if !found {
			return true
		}

		excelFileRaw := value.(*ExcelFileRaw)
		if strings.HasSuffix(excelFileRaw.Sheetname, "client_skip") ||
			excelFileRaw.Type == ExcelSheetType_One ||
			!shouldGenerateTableWrapper(excelFileRaw) {
			return true
		}

		confName := fmt.Sprintf("Conf%s", utils.SimpleTitle(filename))
		tableName := fmt.Sprintf("%sTable", confName)
		keyFieldRaws := getMultiKeyFieldRaws(excelFileRaw)
		paramList := make([]string, 0, len(keyFieldRaws))
		formatArgs := make([]string, 0, len(keyFieldRaws))
		formatHolders := make([]string, 0, len(keyFieldRaws))
		for idx, fieldRaw := range keyFieldRaws {
			paramName := fmt.Sprintf("key%d", idx+1)
			paramList = append(paramList, fmt.Sprintf("%s %s", toCsType(fieldRaw.originTp), paramName))
			formatArgs = append(formatArgs, csKeyToStringExpr(fieldRaw.originTp, paramName))
			formatHolders = append(formatHolders, fmt.Sprintf("{%d}", idx))
		}

		g.P(fmt.Sprintf("    public class %s : ConfigTable<%s>", tableName, confName))
		g.P("    {")
		g.P(fmt.Sprintf("        public %s Get(%s)", confName, strings.Join(paramList, ", ")))
		g.P("        {")
		if len(keyFieldRaws) == 1 {
			// 单主键场景：csKeyToStringExpr 已经返回字符串表达式，可直接透传给 base.Get
			g.P(fmt.Sprintf("            return base.Get(%s);", formatArgs[0]))
		} else {
			// 多主键场景：拼接成 `key1_key2_...` 作为字符串主键
			g.P(fmt.Sprintf("            return base.Get(string.Format(\"%s\", %s));", strings.Join(formatHolders, "_"), strings.Join(formatArgs, ", ")))
		}
		g.P("        }")
		g.P("    }")
		g.P("")
		return true
	})

	enumFilenameTree := treemap.NewWithStringComparator()
	for _, enumRaw := range g.opts.EnumRaws {
		enumFilenameTree.Put(enumRaw.enumOriginType, enumRaw)
	}

	enumFilenameTree.All(func(_, v any) bool {
		enumRaw := v.(*ExcelEnumRaw)
		g.P(fmt.Sprintf("    public enum E%s // %s %s", enumRaw.enumType, enumRaw.sheetName, enumRaw.confName))
		g.P("    {")
		for _, enumValue := range enumRaw.orderedValues {
			g.P(fmt.Sprintf("        %s = %d, // %s", enumValue.Name, enumValue.Value, enumValue.Alias))
		}
		g.P("    }")
		g.P("")
		return true
	})

	g.P("    public static class ConfigData")
	g.P("    {")

	g.P("        static ConfigData()")
	g.P("        {")
	g.P("            Apply();")
	g.P("        }")
	g.P("")

	g.P("        public static bool IsDict(string key)")
	g.P("        {")
	g.P("            return dicts.ContainsKey(key);")
	g.P("        }")
	g.P("")
	g.P("        public static IConfigDict GetDict(string key)")
	g.P("        {")
	g.P("            return dicts.TryGetValue(key, out var dict) ? dict : null;")
	g.P("        }")
	g.P("")
	g.P("        public static IConfigTable GetTable(string key)")
	g.P("        {")
	g.P("            return tables.TryGetValue(key, out var table) ? table : null;")
	g.P("        }")
	g.P("")

	oneFiles := make([]string, 0)
	tableFiles := make([]string, 0)
	requiredFilenameTree.All(func(k, _ any) bool {
		filename := k.(string)
		value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
		if !found {
			return true
		}
		excelFileRaw := value.(*ExcelFileRaw)
		if strings.HasSuffix(excelFileRaw.Sheetname, "client_skip") || len(excelFileRaw.Keys) == 0 {
			return true
		}
		if excelFileRaw.Type == ExcelSheetType_One {
			oneFiles = append(oneFiles, filename)
		} else {
			tableFiles = append(tableFiles, filename)
		}
		return true
	})

	g.P("        private static readonly Dictionary<string, IConfigDict> dicts = new Dictionary<string, IConfigDict>()")
	g.P("        {")
	for _, filename := range oneFiles {
		g.P(fmt.Sprintf("            {\"%s\", new Conf%s()},", utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
	}
	g.P("        };")
	g.P("")

	g.P("        private static readonly Dictionary<string, IConfigTable> tables = new Dictionary<string, IConfigTable>()")
	g.P("        {")
	for _, filename := range tableFiles {
		value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
		if !found {
			continue
		}
		excelFileRaw := value.(*ExcelFileRaw)
		if shouldGenerateTableWrapper(excelFileRaw) {
			g.P(fmt.Sprintf("            {\"%s\", new Conf%sTable()},", utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
		} else {
			g.P(fmt.Sprintf("            {\"%s\", new ConfigTable<Conf%s>()},", utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
		}
	}
	for _, enumRaw := range g.opts.EnumRaws {
		if !shouldExportEnumInfo(enumRaw) {
			continue
		}
		g.P(fmt.Sprintf("            {\"%sEnumInfo\", new ConfigTable<ConfEnumInfo>()},", utils.SimpleTitle(enumRaw.enumOriginType)))
	}
	g.P("        };")
	g.P("")

	g.P("        public static void Apply()")
	g.P("        {")
	for _, filename := range oneFiles {
		g.P(fmt.Sprintf("            %s = (Conf%s)dicts[\"%s\"];", utils.SimpleTitle(filename), utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
	}
	for _, filename := range tableFiles {
		value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
		if !found {
			continue
		}
		excelFileRaw := value.(*ExcelFileRaw)
		if shouldGenerateTableWrapper(excelFileRaw) {
			g.P(fmt.Sprintf("            %s = (Conf%sTable)tables[\"%s\"];", utils.SimpleTitle(filename), utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
		} else {
			g.P(fmt.Sprintf("            %s = (ConfigTable<Conf%s>)tables[\"%s\"];", utils.SimpleTitle(filename), utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
		}
	}
	for _, enumRaw := range g.opts.EnumRaws {
		if !shouldExportEnumInfo(enumRaw) {
			continue
		}
		g.P(fmt.Sprintf("            %sEnumInfo = (ConfigTable<ConfEnumInfo>)tables[\"%sEnumInfo\"];", utils.SimpleTitle(enumRaw.enumOriginType), utils.SimpleTitle(enumRaw.enumOriginType)))
	}
	g.P("        }")
	g.P("")

	for _, filename := range oneFiles {
		g.P(fmt.Sprintf("        public static Conf%s %s;", utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
	}
	for _, filename := range tableFiles {
		value, found := excelFileRaws.Load(fmt.Sprintf("%s.json", filename))
		if !found {
			continue
		}
		excelFileRaw := value.(*ExcelFileRaw)
		if shouldGenerateTableWrapper(excelFileRaw) {
			g.P(fmt.Sprintf("        public static Conf%sTable %s;", utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
		} else {
			g.P(fmt.Sprintf("        public static ConfigTable<Conf%s> %s;", utils.SimpleTitle(filename), utils.SimpleTitle(filename)))
		}
	}
	for _, enumRaw := range g.opts.EnumRaws {
		if !shouldExportEnumInfo(enumRaw) {
			continue
		}
		g.P(fmt.Sprintf("        public static ConfigTable<ConfEnumInfo> %sEnumInfo;", utils.SimpleTitle(enumRaw.enumOriginType)))
	}

	g.P("    }")
	g.P("}")

	return os.WriteFile(g.opts.FilePath, g.buf.Bytes(), 0666)
}

func (g *CsCodeGenerator) Generate() error {
	if len(g.opts.FilePath) == 0 {
		return errors.New("invalid file path")
	}

	if g.opts.CsData != nil {
		root := buildCsRootObject()
		data, err := json.MarshalIndent(root, "", "    ")
		if err != nil {
			return err
		}
		return os.WriteFile(g.opts.FilePath, data, 0666)
	}

	return g.generateCode()
}

func NewCsCodeGenerator(options ...CsCodeGeneratorOption) *CsCodeGenerator {
	g := &CsCodeGenerator{
		opts: defaultCsOptions(),
	}

	for _, o := range options {
		o(g.opts)
	}

	return g
}

func generateRequiredCsCode(exportPath string, requiredFiles []string) error {
	filename := "Configs.cs"
	names := make([]string, 0, len(requiredFiles))
	for _, filename := range requiredFiles {
		names = append(names, strings.TrimSuffix(filename, ".json"))
	}

	enumRaws := make([]*ExcelEnumRaw, 0, 32)
	excelEnumTypeRaws.Range(func(_, v any) bool {
		enumRaws = append(enumRaws, v.(*ExcelEnumRaw))
		return true
	})

	g := NewCsCodeGenerator(
		CsCodeFilePath(filepath.Join(exportPath, filename)),
		CsCodeFileName(filename),
		CsCodeRequiredFileNames(names),
		CsCodeEnumRaws(enumRaws),
	)

	return g.Generate()
}

func generateEntryCsCode(packageName, exportPath string) error {
	filename := "Data.json"
	g := NewCsCodeGenerator(
		CsCodeDatas(&CsCodeData{}),
		CsCodeFilePath(filepath.Join(exportPath, filename)),
		CsCodeFileName(filename),
	)

	return g.Generate()
}

func generateAllCsCodes(exportPath string, fileNames []string) {
	var exportCsJsonPath string = "json"
	if f := flag.Lookup("exportCsJsonPath"); f != nil {
		exportCsJsonPath = f.Value.String()
	}

	_ = os.MkdirAll(exportCsJsonPath, os.ModePerm)
	err := generateEntryCsCode("config", exportCsJsonPath)
	if !utils.ErrCheck(err, "generateEntryCsCode failed", exportPath) {
		return
	}

	requiredFilenames := make([]string, 0, len(fileNames))
	for _, v := range fileNames {
		name := strings.Replace(v, ".xlsx", ".csv", -1)
		value, found := excelFileRaws.Load(name)
		if !found {
			continue
		}

		excelFileRaw := value.(*ExcelFileRaw)
		if excelFileRaw.Type == ExcelSheetType_Enum || excelFileRaw.Type == ExcelSheetType_Poly {
			continue
		}

		requiredFilenames = append(requiredFilenames, name)
	}

	_ = os.MkdirAll(exportPath, os.ModePerm)
	err = generateRequiredCsCode(exportPath, requiredFilenames)
	if !utils.ErrCheck(err, "generateRequiredCsCode failed", exportPath) {
		return
	}
}
