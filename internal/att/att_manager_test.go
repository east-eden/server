package att

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dengsgo/math-engine/engine"
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel"
	"github.com/shopspring/decimal"
)

func TestAttManager(t *testing.T) {
	// 从包目录回到仓库根目录的配置目录：att 表仍来自 CSV，同时加载 JSON 配置
	excel.LoadAll(filepath.Join("..", "..", "config", "csv"), filepath.Join("..", "..", "config", "json"))

	attManager := NewAttManager()
	attManager.SetBaseAttId(1)
	// attManager.ModBaseAtt(define.Att_Atk, 100)

	attManager2 := NewAttManager()
	attManager2.SetBaseAttId(2)
	attManager.ModAttManager(attManager2)
	attManager.CalcAtt()
	_ = attManager.GetFinalAttValue(define.Att_AtkBase)

	d1, _ := decimal.NewFromString("101.57")
	d2, _ := decimal.NewFromString("-382.4")
	x := d1.Round(0).IntPart()
	fmt.Println(x)
	y := d2.IntPart()
	fmt.Println(y)
	round := int32(d1.Round(0).IntPart())
	fmt.Println(round)
	d3 := d1.Mul(d2)
	r := d3.Floor().BigInt().Int64()
	fmt.Println(r)
}

func BenchmarkFormula(b *testing.B) {
	// for n := 0; n < b.N; n++ {
	formula := "(roleLvl^2+2*roleLvl+108)*monsterLvl"
	formula = strings.ReplaceAll(formula, "roleLvl", "100")
	formula = strings.ReplaceAll(formula, "monsterLvl", "200")
	_, err := engine.ParseAndExec(formula)
	if err != nil {
		b.Fatal(err)
	}
	// fmt.Println(result, err)
	// }
}
