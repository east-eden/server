package att

import (
	"testing"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
)

func TestAttManager(t *testing.T) {
	// reload to project root path
	if err := utils.RelocatePath("/server"); err != nil {
		t.Fatalf("TestAttManager failed: %s", err.Error())
	}

	excel.ReadAllEntries("config/excel/")

	attManager := NewAttManager()
	attManager.SetBaseAttId(1)
	attManager.ModBaseAtt(define.Att_Atk, 100)

	attManager2 := NewAttManager()
	attManager2.SetBaseAttId(2)
	attManager.ModAttManager(attManager2)
	attManager.CalcAtt()
	_ = attManager.GetAttValue(define.Att_Atk)
}
