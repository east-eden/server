package ringbuffer

import (
	"fmt"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	rb := New(10)
	_, err := rb.Write([]byte("abcd"))
	if err != nil {
		t.Fatal(err)
	}

	head, tail := rb.LazyRead(2)
	fmt.Println("head:", head, ", tail:", tail)
	rb.Shift(2)
	head, tail = rb.LazyReadAll()
	fmt.Println("head:", head, ", tail:", tail)
}
