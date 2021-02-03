package auto

import "bitbucket.org/east-eden/server/excel"

func init() {
	excel.AddEntryManualLoader("buff.xlsx", (*BuffEntries)(nil))
}

// 特殊类型字段处理
func (e *BuffEntries) ManualLoad(excelFileRaw *excel.ExcelFileRaw) error {
	// rows := GetBuffRows()
	// for _, row := range rows {
	// row.TestField = []int32{1, 2, 3, 4, 5}
	// }
	return nil
}
