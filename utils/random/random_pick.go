package random

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Item interface {
	GetId() int
	GetWeight() int
}

// 按权重随机接口
type RandomPicker interface {
	GetItemList() []Item
}

// 限制器
type Limiter func(Item) bool

// 按权重随机一个
func PickOne(rp RandomPicker, limiter Limiter) (Item, error) {
	itemList := rp.GetItemList()

	if len(itemList) == 0 {
		return nil, errors.New("RandomPicker PickOne: not enough item")
	}

	// 总权重
	totalWeight := func() int {
		var total int
		for _, item := range itemList {
			if limiter(item) {
				total += item.GetWeight()
			}
		}
		return total
	}()

	if totalWeight <= 0 {
		return nil, errors.New("RandomPicker PickOne: total weight is zero")
	}

	rd := rand.Intn(totalWeight + 1)
	for _, item := range itemList {
		if limiter(item) {
			rd -= item.GetWeight()
			if rd <= 0 {
				return item, nil
			}
		}
	}

	return nil, errors.New("RandomPicker PickOne: pick no rand items")
}

// 按权重随机n个不重复的结果
func PickUnrepeated(rp RandomPicker, num int, limiter Limiter) ([]Item, error) {
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
			if limiter(item) {
				total += item.GetWeight()
			}
		}
		return total
	}()

	if totalWeight <= 0 {
		return nil, errors.New("RandomPicker PickUnrepeated: total weight is zero")
	}

	result := make([]Item, 0, num)
	var n int
	for n = 0; n < num; n++ {
		rd := rand.Intn(totalWeight + 1)
		for k, item := range itemList {
			if limiter(item) {
				rd -= item.GetWeight()
				if rd <= 0 {
					result = append(result, item)
					totalWeight -= item.GetWeight()
					itemList = append(itemList[:k], itemList[k+1:]...)
					continue
				}
			}
		}
	}

	if len(result) != num {
		return nil, errors.New("RandomPicker PickUnrepeated: pick no rand items")
	}

	return result, nil
}
