package container

import "bitbucket.org/funplus/server/utils"

// container list
type Container map[interface{}]interface{}
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

func (ca *ContainerArray) Add(idx int, k, v interface{}) bool {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return false
	}

	ca.cons[idx][k] = v
	return true
}

func (ca *ContainerArray) Get(k interface{}) (interface{}, bool) {
	for idx := range ca.cons {
		if v, ok := ca.cons[idx][k]; ok {
			return v, true
		}
	}

	return nil, false
}

func (ca *ContainerArray) GetByIdx(idx int, k interface{}) (interface{}, bool) {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return nil, false
	}

	if v, ok := ca.cons[idx][k]; ok {
		return v, true
	}

	return nil, false
}

func (ca *ContainerArray) Del(k interface{}) bool {
	for idx := range ca.cons {
		if _, ok := ca.cons[idx][k]; ok {
			delete(ca.cons[idx], k)
			return true
		}
	}

	return false
}

func (ca *ContainerArray) DelByIdx(idx int, k interface{}) bool {
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

func (ca *ContainerArray) Range(fn func(v interface{}) bool) {
	for idx := range ca.cons {
		for k := range ca.cons[idx] {
			ok := fn(ca.cons[idx][k])
			if !ok {
				return
			}
		}
	}
}

func (ca *ContainerArray) RangeByIdx(idx int, fn func(v interface{}) bool) {
	if !utils.Between(idx, 0, len(ca.cons)) {
		return
	}

	for k := range ca.cons[idx] {
		ok := fn(ca.cons[idx][k])
		if !ok {
			return
		}
	}
}
