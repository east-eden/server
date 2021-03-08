package random

import "math/rand"

// Int32n returns, as an int32, a non-negative pseudo-random number in [min,max]
// from the default Source.
// It panics if max - min < 0.
func Int32(min, max int32) int32 {
	if (max - min) < 0 {
		panic("invalid argument to Int32")
	}
	return rand.Int31n(max-min+1) + min
}

func Int(min, max int) int {
	if (max - min) < 0 {
		panic("invalid argument to Int")
	}
	return rand.Intn(max-min+1) + min
}
