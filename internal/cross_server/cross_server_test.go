package cross_server

import "testing"

type testDivider struct {
	BaseCrossServerDivider
}

func newDividers(n int) []ICrossServerDivider {
	list := make([]ICrossServerDivider, 0, n)
	for i := 0; i < n; i++ {
		d := &testDivider{}
		d.ObjectId = int64(i)
		list = append(list, d)
	}

	return list
}

// TestDivideOrder 顺序分组：每 objectsPerGroup 个一组
func TestDivideOrder(t *testing.T) {
	dividers := newDividers(10)
	Divide(dividers, 3, 1) // GameplayCrossDivideType_Order

	want := []int32{1, 1, 1, 2, 2, 2, 3, 3, 3, 4}
	for i, d := range dividers {
		if got := d.(*testDivider).GetOrder(); got != want[i] {
			t.Fatalf("index %d order = %d, want %d", i, got, want[i])
		}
	}
}

// TestDivideSnake 蛇形分组：组内顺序来回翻转
func TestDivideSnake(t *testing.T) {
	dividers := newDividers(6)
	Divide(dividers, 3, 0) // GameplayCrossDivideType_Snake

	// groupNum = 6/3 = 2，每 3 个为一轮做蛇形翻转：正向 1,2 再反向 2,1
	want := []int32{1, 2, 2, 1, 1, 2}
	for i, d := range dividers {
		if got := d.(*testDivider).GetOrder(); got != want[i] {
			t.Fatalf("index %d order = %d, want %d", i, got, want[i])
		}
	}
}

// TestDivideSnakeRemainder 多余对象分到新一组
func TestDivideSnakeRemainder(t *testing.T) {
	dividers := newDividers(8)
	Divide(dividers, 3, 0)

	// groupNum = 2，前 6 个蛇形，剩下 2 个分到第 3 组
	if got := dividers[6].(*testDivider).GetOrder(); got != 3 {
		t.Fatalf("remainder order = %d, want 3", got)
	}

	if got := dividers[7].(*testDivider).GetOrder(); got != 3 {
		t.Fatalf("remainder order = %d, want 3", got)
	}
}
