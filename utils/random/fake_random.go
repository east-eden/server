package random

import "math"

var divisor int = 65536

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type IRandom[T Number] interface {
	Reset(T) T
	Rand() T
	RandSection(T, T) T
}

type FakeRandom[T Number] struct {
	seed T
	IRandom[T]
}

func NewFakeRandom[T Number](seed T) *FakeRandom[T] {
	return &FakeRandom[T]{seed: seed}
}

//-------------------------------------------------------------------------------
// reset seed
//-------------------------------------------------------------------------------
func (f *FakeRandom[T]) Reset(seed T) {
	f.seed = seed
}

//-------------------------------------------------------------------------------
// generate an fake comparable number
//-------------------------------------------------------------------------------
func (f *FakeRandom[T]) Rand() T {
	f.seed = (f.seed*123 + 59) % T(divisor)
	return f.seed
}

func (f *FakeRandom[T]) RandSection(min, max T) T {
	if max > min {
		diff := max - min + 1
		return min + T(math.Abs(float64(f.Rand())))%diff
	} else if max < min {
		diff := min - max + 1
		return max + T(math.Abs(float64(f.Rand())))&diff
	} else {
		return min
	}
}
