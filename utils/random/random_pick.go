package random

import (
	"errors"
	"math/rand"
	"time"

	"github.com/east-eden/server/define"
)

var (
	ErrNoResult = errors.New("no result")
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Item[T define.Integer] interface {
	GetId() T
	GetWeight() T
}

// 按权重随机接口
type RandomPicker[T define.Integer] interface {
	GetItemList() []Item[T]
}

// 限制器
type Limiter[T define.Integer] func(Item[T]) bool

// 按权重随机一个
func PickOne[T define.Integer](rp RandomPicker[T], limiter Limiter[T]) (Item[T], error) {
	itemList := rp.GetItemList()

	if len(itemList) == 0 {
		return nil, errors.New("RandomPicker PickOne: not enough item")
	}

	// 总权重
	totalWeight := func() T {
		var total T
		for _, item := range itemList {
			if limiter == nil || limiter(item) {
				total += item.GetWeight()
			}
		}
		return total
	}()

	if totalWeight <= 0 {
		return nil, ErrNoResult
	}

	rd := Int(1, int(totalWeight))
	for _, item := range itemList {
		if limiter == nil || limiter(item) {
			rd -= int(item.GetWeight())
			if rd <= 0 {
				return item, nil
			}
		}
	}

	return nil, ErrNoResult
}

// 按权重随机n个不重复的结果
func PickUnrepeated[T define.Integer](rp RandomPicker[T], num int, limiter Limiter[T]) ([]Item[T], error) {
	itemList := rp.GetItemList()

	if num < 0 {
		return nil, errors.New("RandomPicker PickUnrepeated: invalid num")
	}

	if len(itemList) < num {
		return nil, errors.New("RandomPicker PickUnrepeated: not enough item")
	}

	// 总权重
	totalWeight := func() int {
		var total int
		for _, item := range itemList {
			if limiter == nil || limiter(item) {
				total += int(item.GetWeight())
			}
		}
		return total
	}()

	if totalWeight <= 0 {
		return nil, ErrNoResult
	}

	result := make([]Item[T], 0, num)
	var n int
	for n = 0; n < num; n++ {
		rd := Int(1, totalWeight)
		for k, item := range itemList {
			if limiter == nil || limiter(item) {
				rd -= int(item.GetWeight())
				if rd <= 0 {
					result = append(result, item)
					totalWeight -= int(item.GetWeight())
					itemList = append(itemList[:k], itemList[k+1:]...)
					break
				}
			}
		}
	}

	if len(result) != num {
		return nil, ErrNoResult
	}

	return result, nil
}
