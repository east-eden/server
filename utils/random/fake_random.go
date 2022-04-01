package random

import (
	"math"

	"golang.org/x/exp/constraints"
)

var (
	divisor8  uint8  = math.MaxUint8
	divisor16 uint16 = math.MaxUint16
	divisor32 uint32 = math.MaxUint32
	divisor64 uint64 = math.MaxUint64
)

type IRandom[T constraints.Integer] interface {
	Reset(T) T
	Rand() T
	RandSection(T, T) T
}

type FakeRandom[T constraints.Integer] struct {
	seed T
	IRandom[T]
}

func NewFakeRandom[T constraints.Integer](seed T) *FakeRandom[T] {
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
	switch any(f.seed).(type) {
	case int8, uint8:
		f.seed = (f.seed*123 + 59) % T(divisor8)
	case int16, uint16:
		f.seed = (f.seed*123 + 59) % T(divisor16)
	case int32, uint32:
		f.seed = (f.seed*123 + 59) % T(divisor32)
	case int64, uint64:
		f.seed = (f.seed*123 + 59) % T(divisor64)
	}
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
