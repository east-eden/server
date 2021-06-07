package random

import (
	"math/rand"

	"e.coding.net/mmstudio/blade/server/define"
	"github.com/shopspring/decimal"
)

// Int32n returns, as an int32, a non-negative pseudo-random number in [min,max]
// from the default Source.
// It panics if max - min < 0.
func Int32(min, max int32) int32 {
	if min == max {
		return min
	}

	if (max - min) < 0 {
		panic("invalid argument to Int32")
	}

	return rand.Int31n(max-min+1) + min
}

func Int(min, max int) int {
	if min == max {
		return min
	}

	if (max - min) < 0 {
		panic("invalid argument to Int")
	}

	return rand.Intn(max-min+1) + min
}

func Decimal(min, max decimal.Decimal) decimal.Decimal {
	if min.Equal(max) {
		return min
	}

	imin := min.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()
	imax := max.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()
	if (imax - imin) < 0 {
		panic("invalid argument to Decimal")
	}

	rd := rand.Int63n(imax-imin+1) + imin
	ret := decimal.NewFromInt(rd)
	return ret.Div(decimal.NewFromInt(define.PercentBase))
}

func DecimalFake(min, max decimal.Decimal, fake *FakeRandom) decimal.Decimal {
	if min.Equal(max) {
		return min
	}

	imin := min.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()
	imax := max.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()
	if (imax - imin) < 0 {
		panic("invalid argument to DecimalFake")
	}

	rd := fake.RandSection(int(imin), int(imax))
	ret := decimal.NewFromInt(int64(rd))
	return ret.Div(decimal.NewFromInt(define.PercentBase))
}

func Interface(min, max interface{}) interface{} {
	switch min.(type) {
	case int:
		return Int(min.(int), max.(int))
	case int32:
		return Int32(min.(int32), max.(int32))
	case decimal.Decimal:
		return Decimal(min.(decimal.Decimal), max.(decimal.Decimal))
	default:
		return int(0)
	}
}
