package att

import (
	"os"
	"testing"

	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
)

func TestAttManager(t *testing.T) {
	os.Chdir("../../")
	global.InitEntries()

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
