package random

import "math"

type FakeRandom struct {
	seed int
}

func NewFakeRandom(seed int) *FakeRandom {
	return &FakeRandom{seed: seed}
}

//-------------------------------------------------------------------------------
// 重置种子
//-------------------------------------------------------------------------------
func (f *FakeRandom) Reset(seed int) {
	f.seed = seed
}

//-------------------------------------------------------------------------------
// 生成一个SHORT伪随机数
//-------------------------------------------------------------------------------
func (f *FakeRandom) Rand() int {
	f.seed = (f.seed*123 + 59) % 65536
	return f.seed
}

func (f *FakeRandom) RandSection(min, max int) int {
	if max > min {
		diff := max - min + 1
		return min + int(math.Abs(float64(f.Rand())))%diff
	} else if max < min {
		diff := min - max + 1
		return max + int(math.Abs(float64(f.Rand())))&diff
	} else {
		return min
	}
}
