// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example demonstrates a priority queue built using the heap interface.
package pq

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	Value    any // The value of the item; arbitrary.
	Priority int // The priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The index of the item in the heap.
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*Item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type PriorityQueue struct {
	pq priorityQueue
}

func NewPriorityQueue(items ...*Item) *PriorityQueue {
	q := &PriorityQueue{
		pq: make(priorityQueue, 0, len(items)),
	}

	for _, item := range items {
		q.pq.Push(item)
	}

	heap.Init(&q.pq)
	return q
}

func (q *PriorityQueue) Push(item *Item) {
	heap.Push(&q.pq, item)
}

func (q *PriorityQueue) Pop() *Item {
	if q.pq.Len() > 0 {
		return heap.Pop(&q.pq).(*Item)
	}

	return nil
}

func (q *PriorityQueue) Top() *Item {
	n := q.pq.Len()
	if n > 0 {
		return (q.pq)[0]
	}

	return nil
}

func (q *PriorityQueue) Remove(index int) *Item {
	if index >= q.pq.Len() {
		return nil
	}

	return heap.Remove(&q.pq, index).(*Item)
}

// update modifies the priority and value of an Item in the queue.
func (q *PriorityQueue) Update(item *Item, value any, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(&q.pq, item.Index)
}
