package container

import (
	"github.com/east-eden/server/utils"
	"golang.org/x/exp/constraints"
)

// Container 是一个带常用操作的 map 封装
type Container[K comparable, V any] map[K]V

func (c Container[K, V]) Add(k K, v V) {
	c[k] = v
}

func (c Container[K, V]) Get(k K) (V, bool) {
	v, found := c[k]
	return v, found
}

func (c Container[K, V]) Del(k K) {
	delete(c, k)
}

func (c Container[K, V]) Size() int {
	return len(c)
}

func (c Container[K, V]) Range(fn func(v V) bool) {
	for k := range c {
		if !fn(c[k]) {
			return
		}
	}
}

// ContainerArray 按索引分组的容器数组，例如按背包类型划分的格子
type ContainerArray[T constraints.Integer, K comparable, V any] struct {
	Containers []Container[K, V] `bson:"containers"`
}

func New[T constraints.Integer, K constraints.Ordered, V any](size T) *ContainerArray[T, K, V] {
	ca := &ContainerArray[T, K, V]{
		Containers: make([]Container[K, V], size),
	}

	for idx := range ca.Containers {
		ca.Containers[idx] = Container[K, V]{}
	}

	return ca
}

func (ca *ContainerArray[T, K, V]) Add(idx T, k K, v V) bool {
	if !utils.Between(idx, 0, T(len(ca.Containers))) {
		return false
	}

	ca.Containers[idx].Add(k, v)
	return true
}

func (ca *ContainerArray[T, K, V]) Get(k K) (V, bool) {
	for idx := range ca.Containers {
		if v, ok := ca.Containers[idx].Get(k); ok {
			return v, true
		}
	}

	var v V
	return v, false
}

func (ca *ContainerArray[T, K, V]) GetByIdx(idx T, k K) (V, bool) {
	var v V
	if !utils.Between(idx, 0, T(len(ca.Containers))) {
		return v, false
	}

	if v, ok := ca.Containers[idx].Get(k); ok {
		return v, true
	}

	return v, false
}

func (ca *ContainerArray[T, K, V]) Del(k K) bool {
	for idx := range ca.Containers {
		if _, ok := ca.Containers[idx].Get(k); ok {
			ca.Containers[idx].Del(k)
			return true
		}
	}

	return false
}

func (ca *ContainerArray[T, K, V]) DelByIdx(idx T, k K) bool {
	if !utils.Between(idx, 0, T(len(ca.Containers))) {
		return false
	}

	if _, ok := ca.Containers[idx].Get(k); ok {
		ca.Containers[idx].Del(k)
		return true
	}

	return false
}

func (ca *ContainerArray[T, K, V]) Size(idx T) T {
	if !utils.Between(idx, 0, T(len(ca.Containers))) {
		return 0
	}

	return T(ca.Containers[idx].Size())
}

func (ca *ContainerArray[T, K, V]) Range(fn func(v V) bool) {
	for idx := range ca.Containers {
		for k := range ca.Containers[idx] {
			if !fn(ca.Containers[idx][k]) {
				return
			}
		}
	}
}

func (ca *ContainerArray[T, K, V]) RangeByIdx(idx T, fn func(v V) bool) {
	if !utils.Between(idx, 0, T(len(ca.Containers))) {
		return
	}

	ca.Containers[idx].Range(fn)
}
