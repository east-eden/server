package attrcache

import (
	"math"
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
)

func TestEvalWithVars(t *testing.T) {
	// 等级 1 时 (level-1) 项为 0
	got := EvalWithVars("580*0.0158+(level-1)*0.0158*15000", 1, 1)
	want := 580 * 0.0158
	if math.Abs(got-want) > 1e-6 {
		t.Fatalf("got %v, want %v", got, want)
	}

	// 等级 2 时第二项生效
	got2 := EvalWithVars("580*0.0158+(level-1)*0.0158*15000", 2, 1)
	want2 := 580*0.0158 + 0.0158*15000
	if math.Abs(got2-want2) > 1e-6 {
		t.Fatalf("got %v, want %v", got2, want2)
	}
}

func TestReplaceVarBoundary(t *testing.T) {
	// playerLevel 不能被 level 的替换污染
	got := replaceVar(replaceVar("level+playerLevel", "playerLevel", 7), "level", 3)
	if got != "3+7" {
		t.Fatalf("got %q, want %q", got, "3+7")
	}
}

// TestRealBattlePropertyFormula 用真实属性表公式验证求值链路
func TestRealBattlePropertyFormula(t *testing.T) {
	excel.ReadAllJsonEntries("../../config/json")

	rows := config.GetBattlePropertyRows()
	if len(rows) == 0 {
		t.Skip("battle property table not loaded")
	}

	var formula string
	for _, e := range rows {
		if e.Level == 1 && len(e.MAXHP) > 0 {
			formula = e.MAXHP
			break
		}
	}

	if formula == "" {
		t.Skip("no level 1 formula found")
	}

	if v := EvalWithVars(formula, 1, 1); v <= 0 {
		t.Fatalf("formula %q evaluated to %v", formula, v)
	}
}
