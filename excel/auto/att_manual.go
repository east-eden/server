package auto

import (
	"bitbucket.org/east-eden/server/excel"
)

func init() {
	excel.AddEntryManualLoader("Att.xlsx", (*AttEntries)(nil))
}

// 手动加载处理
func (e *AttEntries) ManualLoad(excelFileRaw *excel.ExcelFileRaw) error {

	return nil
}
