package cross_server

import "github.com/east-eden/server/excel/config"

type ICrossServerDivider interface {
	SetOrder(order int32)
}

type BaseCrossServerDivider struct {
	ObjectId int64
	Score    int64
	Order    int32
	ICrossServerDivider
}

func (d *BaseCrossServerDivider) SetOrder(order int32) {
	d.Order = order
}

func (d *BaseCrossServerDivider) GetOrder() int32 {
	return d.Order
}

func (d *BaseCrossServerDivider) GetObjectId() int64 {
	return d.ObjectId
}

// 分组
func Divide(arrayDividers []ICrossServerDivider, objectsPerGroup int32, divideType config.GameplayCrossDivideType) {
	groupNum := int32(len(arrayDividers)) / objectsPerGroup

	switch divideType {
	case config.GameplayCrossDivideType_Snake:
		// 前整数分组按蛇形排序
		var inc bool
		for n := int32(0); n < int32(groupNum)*objectsPerGroup; n++ {
			// 每遍历groupNum一次，顺序翻转
			if n%int32(groupNum) == 0 {
				inc = !inc
			}

			if inc {
				arrayDividers[n].SetOrder(n%int32(groupNum) + 1)
			} else {
				arrayDividers[n].SetOrder(int32(groupNum) - 1 - (n % int32(groupNum)) + 1)
			}
		}

		// 剩余的全部分到新一组中
		for n := int32(groupNum) * objectsPerGroup; n < int32(len(arrayDividers)); n++ {
			arrayDividers[n].SetOrder(int32(groupNum) + 1)
		}

	case config.GameplayCrossDivideType_Order:
		// 按顺序分组
		for n := int32(0); n < int32(len(arrayDividers)); n++ {
			arrayDividers[n].SetOrder(n/objectsPerGroup + 1)
		}
	}
}
