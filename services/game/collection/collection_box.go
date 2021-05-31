package collection

import (
	"errors"

	"github.com/east-eden/server/excel/auto"
)

var (
	ErrCollectionAlreadyInBox = errors.New("collection already in box")
	ErrCollectionBoxSlotLimit = errors.New("collection box slot limit")
	ErrCollectionNotPutinBox  = errors.New("collection not put in box")
)

// 收集品放置管理
type CollectionBox struct {
	tp             int32
	collectionList map[int32]*Collection
	Entry          *auto.CollectionBoxEntry
}

func NewCollectionBox(tp int32) *CollectionBox {
	m := &CollectionBox{
		tp:             tp,
		collectionList: make(map[int32]*Collection, 8),
	}

	m.Entry, _ = auto.GetCollectionBoxEntry(tp)
	return m
}

func (cb *CollectionBox) PutonCollection(c *Collection) error {
	if c.BoxId != -1 {
		return ErrCollectionAlreadyInBox
	}

	if _, ok := cb.collectionList[c.TypeId]; ok {
		return ErrCollectionAlreadyInBox
	}

	if len(cb.collectionList) >= int(cb.Entry.MaxSlot) {
		return ErrCollectionBoxSlotLimit
	}

	cb.collectionList[c.TypeId] = c
	c.BoxId = cb.tp
	return nil
}

func (cb *CollectionBox) TakeoffCollection(c *Collection) error {
	if c.BoxId == -1 {
		return ErrCollectionNotPutinBox
	}

	if _, ok := cb.collectionList[c.TypeId]; !ok {
		return ErrCollectionNotPutinBox
	}

	delete(cb.collectionList, c.TypeId)
	c.BoxId = -1
	return nil
}

func (cb *CollectionBox) GetActiveEffect() int32 {
	var totalScore int32
	for _, c := range cb.collectionList {
		totalScore += c.score
	}

	var curIdx int = -1
	for k, score := range cb.Entry.Scores {
		if totalScore < score {
			break
		}

		curIdx = k
	}

	if curIdx == -1 {
		return -1
	}

	return cb.Entry.Effects[curIdx]
}
