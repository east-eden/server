package utils

import (
	"fmt"
	"sync"
)

type IDGenerator struct {
	id int64
	mu sync.RWMutex
}

var idgs []*IDGenerator

func NewIDGenerators(size int) []*IDGenerator {
	if idgs == nil {
		idgs = make([]*IDGenerator, 0, size)

		for n := 0; n < size; n++ {
			idgs = append(idgs, &IDGenerator{id: 0})
		}
	}

	return idgs
}

func GeneralIDGet(tp int) (int64, error) {
	if tp >= len(idgs) {
		return -1, fmt.Errorf("wrong id generator type:%d", tp)
	}

	idgs[tp].mu.RLock()
	defer idgs[tp].mu.RUnlock()
	return idgs[tp].id, nil
}

func GeneralIDSet(tp int, id int64) error {
	if tp >= len(idgs) {
		return fmt.Errorf("wrong id generator type:%d", tp)
	}

	idgs[tp].mu.Lock()
	idgs[tp].id = id
	defer idgs[tp].mu.Unlock()
	return nil
}

func GeneralIDGen(tp int) (int64, error) {
	if tp >= len(idgs) {
		return -1, fmt.Errorf("wrong id generator type:%d", tp)
	}

	idgs[tp].mu.Lock()
	defer idgs[tp].mu.Unlock()
	idgs[tp].id++
	return idgs[tp].id, nil
}
