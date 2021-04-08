package player

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/valyala/bytebufferpool"
)

func makeHeroKey(heroId int64, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("hero_list.")
	_, _ = b.WriteString(strconv.Itoa(int(heroId)))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

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
