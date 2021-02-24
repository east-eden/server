package utils

import (
	"fmt"

	"github.com/willf/bitset"
)

type CountableBitset struct {
	state *bitset.BitSet
	count []int16
}

func NewCountableBitset(length uint) *CountableBitset {
	return &CountableBitset{
		state: bitset.New(length),
		count: make([]int16, length),
	}
}

func FromCountableBitset(buf []uint64, count []int16) *CountableBitset {
	return &CountableBitset{
		state: bitset.From(buf),
		count: count,
	}
}

func (b *CountableBitset) Bytes() []uint64 {
	return b.state.Bytes()
}

func (b *CountableBitset) Test(i uint) bool {
	if i > b.state.Len() {
		panic(fmt.Sprintf("CountableBitset: index<%d> out of range", i))
	}

	return b.state.Test(i)
}

func (b *CountableBitset) Any() bool {
	return b.state.Any()
}

func (b *CountableBitset) Intersection(compare *CountableBitset) (result *CountableBitset) {
	result = NewCountableBitset(b.state.Len())
	result.state = b.state.Intersection(compare.state)
	return
}

func (b *CountableBitset) Set(i uint, count int16) *CountableBitset {
	if i > b.state.Len() {
		panic(fmt.Sprintf("CountableBitset: index<%d> out of range", i))
	}

	b.state.Set(i)
	b.count[i] += count
	return b
}

func (b *CountableBitset) Clear(i uint, count int16) *CountableBitset {
	if i > b.state.Len() {
		panic(fmt.Sprintf("CountableBitset: index<%d> out of range", i))
	}

	b.state.Clear(i)
	b.count[i] -= count
	return b
}

func (b *CountableBitset) ClearAll() *CountableBitset {
	b.state.ClearAll()
	for k := range b.count {
		b.count[k] = 0
	}

	return b
}
