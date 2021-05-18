package utils

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type timeCase struct {
	t1     time.Time
	t2     time.Time
	hour   int
	expect bool
	index  int
}

var timeCases = []*timeCase{
	{
		t1:     time.Date(2021, time.May, 18, 0, 0, 0, 0, time.UTC),
		t2:     time.Date(2021, time.May, 18, 1, 0, 0, 0, time.UTC),
		hour:   0,
		expect: true,
	},

	{
		t1:     time.Date(2021, time.May, 18, 0, 0, 0, 0, time.UTC),
		t2:     time.Date(2021, time.May, 18, 1, 0, 0, 0, time.UTC),
		hour:   5,
		expect: true,
	},

	{
		t1:     time.Date(2021, time.May, 17, 6, 0, 0, 0, time.UTC),
		t2:     time.Date(2021, time.May, 18, 4, 0, 0, 0, time.UTC),
		hour:   5,
		expect: true,
	},

	{
		t1:     time.Date(2021, time.May, 18, 6, 0, 0, 0, time.UTC),
		t2:     time.Date(2021, time.May, 18, 23, 0, 0, 0, time.UTC),
		hour:   5,
		expect: true,
	},

	{
		t1:     time.Date(2021, time.May, 18, 4, 0, 0, 0, time.UTC),
		t2:     time.Date(2021, time.May, 18, 5, 0, 0, 0, time.UTC),
		hour:   5,
		expect: false,
	},

	{
		t1:     time.Date(2021, time.May, 18, 6, 0, 0, 0, time.UTC),
		t2:     time.Date(2021, time.May, 19, 4, 0, 0, 0, time.UTC),
		hour:   5,
		expect: true,
	},

	{
		t1:     time.Date(2021, time.May, 18, 6, 0, 0, 0, time.UTC),
		t2:     time.Date(2021, time.May, 19, 5, 0, 0, 0, time.UTC),
		hour:   5,
		expect: false,
	},
}

func TestInSameDay(t *testing.T) {
	for _, c := range timeCases {
		b := IsInSameDay(c.t1, c.t2, c.hour)
		diff := cmp.Diff(b, c.expect)
		if diff != "" {
			t.Fatalf("result not expected, index = %d, diff = %s", c.index, diff)
		}
	}
}
