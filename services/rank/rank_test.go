package rank

import (
	"fmt"
	"testing"
	"time"

	"github.com/liyiheng/zset"
)

func TestRankZSet(t *testing.T) {
	now := time.Now()
	tm1 := int32(now.Unix())
	tm2 := int32(now.Add(2 * time.Second).Unix())
	tm3 := int32(now.Add(4 * time.Second).Unix())

	s := zset.New()
	// add data
	s.Set(66, 1001, tm1)
	s.Set(66, 1002, tm2)
	s.Set(88, 1003, tm3)

	// get rank by id
	rank, score, extra := s.GetRank(1002, false)
	fmt.Printf("1002's rank = %d, score = %f, extra = %v\n", rank, score, extra)

	// get data by rank
	id, score, extra := s.GetDataByRank(0, true)
	fmt.Printf("rank0's key = %d, score = %f, extra = %v\n", id, score, extra)

	// get data by id
	data, ok := s.GetData(1001)
	if ok {
		fmt.Printf("1001's data = %v\n", data)
	}

	// delete data by id
	s.Delete(1001)

	// Increase score
	s.IncrBy(5.0, 1001)

	// ZRANGE, ASC
	five := make([]int64, 0, 5)
	s.Range(0, 5, func(score float64, k int64, _ interface{}) {
		five = append(five, k)
	})

	// ZREVRANGE, DESC
	all := make([]int64, 0)
	s.RevRange(0, -1, func(score float64, k int64, _ interface{}) {
		all = append(all, k)
	})
}
