package player

import (
	"fmt"
	"testing"
)

func BenchmarkBufferPoolOnlyKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = makeHeroKey(int64(i))
	}
}

func BenchmarkBufferPoolKeyAndField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = makeHeroKey(int64(i), "field")
	}
}

func BenchmarkSprintfOnlyKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("hero_map.id_%d", int64(i))
	}
}

func BenchmarkSprintfKeyAndField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("hero_map.id_%d.%s", int64(i), "field")
	}
}
