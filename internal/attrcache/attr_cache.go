// Package attrcache 缓存属性公式的求值结果。
//
// kof 的属性表（BattleProperty.json）用字符串公式描述等级成长，例如
// "580*0.0158+(level-1)*0.0158*15000"；同一公式在同一等级下会被反复求值，
// 直接解析开销很大，因此按“已代入参数的公式串”做一层缓存。
package attrcache

import (
	"strconv"
	"strings"
	"sync"

	"github.com/dengsgo/math-engine/engine"
	"github.com/east-eden/server/utils"
)

var (
	once  sync.Once
	cache sync.Map // formula -> float64
)

// GetAttrCache 解析并缓存公式求值结果（公式需已代入所有变量）
func GetAttrCache(formula string) float64 {
	if v, ok := cache.Load(formula); ok {
		return v.(float64)
	}

	value, err := engine.ParseAndExec(formula)
	if !utils.ErrCheck(err, "ParseAndExec failed", formula) {
		return 0
	}

	cache.Store(formula, value)
	return value
}

// EvalWithVars 把公式里的 level / playerLevel 变量替换成具体值后求值。
//
// 注意替换顺序：先替换 playerLevel，再替换 level，并且用边界匹配避免把
// playerLevel 误替换成 player10。
func EvalWithVars(formula string, level, playerLevel int32) float64 {
	replaced := replaceVar(formula, "playerLevel", playerLevel)
	replaced = replaceVar(replaced, "level", level)

	return GetAttrCache(replaced)
}

func replaceVar(formula, name string, value int32) string {
	var b strings.Builder
	b.Grow(len(formula) + 8)

	s := strconv.Itoa(int(value))
	for i := 0; i < len(formula); {
		if !strings.HasPrefix(formula[i:], name) {
			b.WriteByte(formula[i])
			i++
			continue
		}

		// 变量名两侧必须是边界，避免 level 命中 playerLevel
		prevOK := i == 0 || !isIdent(formula[i-1])
		end := i + len(name)
		nextOK := end >= len(formula) || !isIdent(formula[end])
		if !prevOK || !nextOK {
			b.WriteByte(formula[i])
			i++
			continue
		}

		b.WriteString(s)
		i = end
	}

	return b.String()
}

func isIdent(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
