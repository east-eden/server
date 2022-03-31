package container

import (
	"github.com/east-eden/server/utils"
)

// container list
type Container map[any]any
type ContainerArray struct {
	cons []Container
}

func New(size int) *ContainerArray {
	ca := &ContainerArray{
		cons: make([]Container, size),
	}

	for k := range ca.cons {
		ca.cons[k] = make(Container)
	}

	return ca
}

func (ca *ContainerArray) Add(idx int, k, v any) bool {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return false
	}

	ca.cons[idx][k] = v
	return true
}

func (ca *ContainerArray) Get(k any) (any, bool) {
	for idx := range ca.cons {
		if v, ok := ca.cons[idx][k]; ok {
			return v, true
		}
	}

	return nil, false
}

func (ca *ContainerArray) GetByIdx(idx int, k any) (any, bool) {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return nil, false
	}

	if v, ok := ca.cons[idx][k]; ok {
		return v, true
	}

	return nil, false
}

func (ca *ContainerArray) Del(k any) bool {
	for idx := range ca.cons {
		if _, ok := ca.cons[idx][k]; ok {
			delete(ca.cons[idx], k)
			return true
		}
	}

	return false
}

func (ca *ContainerArray) DelByIdx(idx int, k any) bool {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return false
	}

	if _, ok := ca.cons[idx][k]; ok {
		delete(ca.cons[idx], k)
		return true
	}

	return false
}

func (ca *ContainerArray) Size(idx int) int {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return 0
	}

	return len(ca.cons[idx])
}

func (ca *ContainerArray) Range(fn func(v any) bool) {
	for idx := range ca.cons {
		for k := range ca.cons[idx] {
			if !fn(ca.cons[idx][k]) {
				return
			}
		}
	}
}

func (ca *ContainerArray) RangeByIdx(idx int, fn func(v any) bool) {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return
	}

	for k := range ca.cons[idx] {
		if !fn(ca.cons[idx][k]) {
			return
		}
	}
}
