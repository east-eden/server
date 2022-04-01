package container

import (
	"container/list"

	"github.com/east-eden/server/utils"
	"golang.org/x/exp/constraints"
)

// container list
type Container[K comparable, V any] map[K]V
type ContainerArray[T constraints.Integer, K comparable, V any] struct {
	cons []Container[K, V]
}

func New[T constraints.Integer, K constraints.Ordered, V any](size T) *ContainerArray[T, K, V] {
	ca := &ContainerArray[T, K, V]{
		cons: make([]Container[K, V], size),
	}

	list.New()

	for k := range ca.cons {
		ca.cons[k] = make(Container[K, V])
	}

	return ca
}

func (ca *ContainerArray[T, K, V]) Add(idx T, k K, v V) bool {
	if !utils.Between(idx, 0, T(len(ca.cons))) {
		return false
	}

	ca.cons[idx][k] = v
	return true
}

func (ca *ContainerArray[T, K, V]) Get(k K) (V, bool) {
	for idx := range ca.cons {
		if v, ok := ca.cons[idx][k]; ok {
			return v, true
		}
	}

	var v V
	return v, false
}

func (ca *ContainerArray[T, K, V]) GetByIdx(idx T, k K) (V, bool) {
	var v V
	if !utils.Between(idx, 0, T(len(ca.cons))) {
		return v, false
	}

	if v, ok := ca.cons[idx][k]; ok {
		return v, true
	}

	return v, false
}

func (ca *ContainerArray[T, K, V]) Del(k K) bool {
	for idx := range ca.cons {
		if _, ok := ca.cons[idx][k]; ok {
			delete(ca.cons[idx], k)
			return true
		}
	}

	return false
}

func (ca *ContainerArray[T, K, V]) DelByIdx(idx T, k K) bool {
	if !utils.Between(idx, 0, T(len(ca.cons))) {
		return false
	}

	if _, ok := ca.cons[idx][k]; ok {
		delete(ca.cons[idx], k)
		return true
	}

	return false
}

func (ca *ContainerArray[T, K, V]) Size(idx T) T {
	if !utils.Between(idx, 0, T(len(ca.cons))) {
		return 0
	}

	return T(len(ca.cons[idx]))
}

func (ca *ContainerArray[T, K, V]) Range(fn func(v any) bool) {
	for idx := range ca.cons {
		for k := range ca.cons[idx] {
			if !fn(ca.cons[idx][k]) {
				return
			}
		}
	}
}

func (ca *ContainerArray[T, K, V]) RangeByIdx(idx T, fn func(v V) bool) {
	if !utils.Between(idx, 0, T(len(ca.cons))) {
		return
	}

	for k := range ca.cons[idx] {
		if !fn(ca.cons[idx][k]) {
			return
		}
	}
}
