package att

import (
	"testing"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
)

func TestAttManager(t *testing.T) {
	entries.InitEntries()

	attManager := NewAttManager(1)
	attManager.ModBaseAtt(define.Att_Str, 100)

	attManager2 := NewAttManager(2)
	attManager.ModAttManager(attManager2)
	attManager.CalcAtt()
	attStr := attManager.GetAttValue(define.Att_Str)
	if attStr == 0 {
		t.Errorf("att manager calc failed")
	}
}
