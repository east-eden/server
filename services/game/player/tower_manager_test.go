package player

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type SettleTimeCase struct {
	last   time.Time
	now    time.Time
	expect int
	index  int
}

var settleCases = []*SettleTimeCase{
	{
		last:   time.Date(2021, time.May, 18, 2, 0, 0, 0, time.UTC),
		now:    time.Date(2021, time.May, 18, 3, 0, 0, 0, time.UTC),
		expect: 0,
		index:  1,
	},

	{
		last:   time.Date(2021, time.May, 18, 2, 0, 0, 0, time.UTC),
		now:    time.Date(2021, time.May, 18, 5, 0, 0, 0, time.UTC),
		expect: 1,
		index:  2,
	},

	{
		last:   time.Date(2021, time.May, 18, 2, 0, 0, 0, time.UTC),
		now:    time.Date(2021, time.May, 19, 1, 0, 0, 0, time.UTC),
		expect: 1,
		index:  3,
	},

	{
		last:   time.Date(2021, time.May, 18, 6, 0, 0, 0, time.UTC),
		now:    time.Date(2021, time.May, 18, 7, 0, 0, 0, time.UTC),
		expect: 0,
		index:  4,
	},

	{
		last:   time.Date(2021, time.May, 18, 6, 0, 0, 0, time.UTC),
		now:    time.Date(2021, time.May, 19, 1, 0, 0, 0, time.UTC),
		expect: 0,
		index:  5,
	},

	{
		last:   time.Date(2021, time.May, 18, 6, 0, 0, 0, time.UTC),
		now:    time.Date(2021, time.May, 19, 5, 1, 0, 0, time.UTC),
		expect: 1,
		index:  6,
	},

	{
		last:   time.Date(2021, time.May, 18, 6, 0, 0, 0, time.UTC),
		now:    time.Date(2021, time.May, 20, 1, 0, 0, 0, time.UTC),
		expect: 1,
		index:  7,
	},
}

func TestSettleDays(t *testing.T) {
	for _, c := range settleCases {
		days := CalcSettleDays(c.now, c.last)
		diff := cmp.Diff(days, c.expect)
		if diff != "" {
			t.Fatalf("result not expected, index = %d, diff = %s", c.index, diff)
		}
	}
}
