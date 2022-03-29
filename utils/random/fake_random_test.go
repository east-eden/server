package random

import (
	"testing"
)

func FuzzReverse(f *testing.F) {
	testcases := []any{1, 2, 7000}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig int) {
		fr := NewFakeRandom(orig)
		for n := 0; n < 10; n++ {
			t.Logf("fuzzing generate random seed<%d>'s random number<%d>: %d\n", orig, n, fr.Rand())
		}
	})
}
