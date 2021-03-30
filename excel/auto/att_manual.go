package auto

import (
	"github.com/east-eden/server/excel"
)

// number = excel读取数值 * 10000
type number int32

func init() {
	excel.AddEntryManualLoader("Att.xlsx", (*AttEntries)(nil))
}

// 手动加载处理
func (e *AttEntries) ManualLoad(excelFileRaw *excel.ExcelFileRaw) error {
	return nil
}
