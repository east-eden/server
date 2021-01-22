package att

import (
	"testing"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
)

func TestAttManager(t *testing.T) {
	// reload to project root path
	if err := utils.RelocatePath(); err != nil {
		t.Fatalf("TestAttManager failed: %s", err.Error())
	}

	excel.ReadAllEntries("config/excel")

	attManager := NewAttManager(1)
	attManager.ModBaseAtt(define.Att_Atk, 100)

	attManager2 := NewAttManager(2)
	attManager.ModAttManager(attManager2)
	attManager.CalcAtt()
	_ = attManager.GetAttValue(define.Att_Atk)
}
