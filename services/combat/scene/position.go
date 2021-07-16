package scene

import "github.com/shopspring/decimal"

type Pos struct {
	X decimal.Decimal
	Z decimal.Decimal
}

// 坐标信息
type Position struct {
	Pos
	Rotate decimal.Decimal
}

func IsInDistance(a, b *Position, distance int32) bool {
	x2 := a.X.Sub(b.X).Mul(a.X)
	z2 := a.Z.Sub(b.Z).Mul(a.Z)
	d := decimal.NewFromInt32(distance)
	return x2.Add(z2).LessThanOrEqual(d.Mul(d))
}
