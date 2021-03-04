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

// 按权重随机一个
func PickOne(rp RandomPicker) (Item, error) {
	itemList := rp.GetItemList()

	if len(itemList) == 0 {
		return nil, errors.New("RandomPicker PickOne: not enough item")
	}

	// 总权重
	totalWeight := func() int {
		var total int
		for _, value := range itemList {
			total += value.GetWeight()
		}
		return total
	}()

	if totalWeight <= 0 {
		return nil, errors.New("RandomPicker PickOne: total weight is zero")
	}

	rd := rand.Intn(totalWeight + 1)
	for _, it := range itemList {
		rd -= it.GetWeight()
		if rd <= 0 {
			return it, nil
		}
	}

	return nil, errors.New("RandomPicker PickOne: pick no rand items")
}

// 按权重随机n个不重复的结果
func PickUnrepeated(rp RandomPicker, num int) ([]Item, error) {
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
		for _, value := range itemList {
			total += value.GetWeight()
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
		for k, it := range itemList {
			rd -= it.GetWeight()
			if rd <= 0 {
				result = append(result, it)
				totalWeight -= it.GetWeight()
				itemList = append(itemList[:k], itemList[k+1:]...)
				continue
			}
		}
	}

	if len(result) != num {
		return nil, errors.New("RandomPicker PickUnrepeated: pick no rand items")
	}

	return result, nil
}
