package scene

type Pos struct {
	X int32
	Z int32
}

// 坐标信息
type Position struct {
	Pos
	Rotate int32
}

func IsInDistance(a, b *Position, distance int32) bool {
	x := a.X - b.X
	z := a.Z - b.Z
	return x*x+z*z <= distance*distance
}
