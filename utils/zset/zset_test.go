package zset

import (
	"fmt"
	"testing"
)

type zsetData struct {
	id   int64
	name string
}

func TestRankZSet(t *testing.T) {
	s := New()

	// add data
	s.Set(77, 1, 1000, &zsetData{id: 1, name: "player1"})
	s.Set(80, 2, 1001, &zsetData{id: 2, name: "player2"})
	s.Set(80, 3, 1002, &zsetData{id: 3, name: "player3"})

	rank, score, data := s.GetRank(1, true)
	fmt.Println("player1:", rank, score, data)

	rank, score, data = s.GetRank(2, true)
	fmt.Println("player2:", rank, score, data)

	rank, score, data = s.GetRank(3, true)
	fmt.Println("player3:", rank, score, data)
}
