// Code generated by gotemplate. DO NOT EDIT.

package xlistener

import (
	"sync"
)

// A thread safe map.
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
var (
	SHARD_COUNTConcurrentMapInt64FdInfo = 32
)

// template type ConcurrentMap(KType,VType,KeyHash)

type ConcurrentMapInt64FdInfo []*shardedConcurrentMapInt64FdInfo

type shardedConcurrentMapInt64FdInfo struct {
	items map[int]*fdInfo
	sync.RWMutex
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type TupleConcurrentMapInt64FdInfo struct {
	Key int
	Val *fdInfo
}

func NewConcurrentMapInt64FdInfo() ConcurrentMapInt64FdInfo {
	this := make(ConcurrentMapInt64FdInfo, SHARD_COUNTConcurrentMapInt64FdInfo)
	for i := 0; i < SHARD_COUNTConcurrentMapInt64FdInfo; i++ {
		this[i] = &shardedConcurrentMapInt64FdInfo{items: make(map[int]*fdInfo)}
	}
	return this
}

// Returns shard under given key.
func (m ConcurrentMapInt64FdInfo) GetShard(key int) *shardedConcurrentMapInt64FdInfo {
	return m[uint64(func(key int) int {
		return key
	}(key))%uint64(SHARD_COUNTConcurrentMapInt64FdInfo)]
}

func (m ConcurrentMapInt64FdInfo) MSet(data map[int]*fdInfo) {
	for key, value := range data {
		shard := m.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

// IsEmpty checks if map is empty.
func (m ConcurrentMapInt64FdInfo) IsEmpty() bool {
	return m.Count() == 0
}

func (m *ConcurrentMapInt64FdInfo) Set(key int, value *fdInfo) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// like redis SETNX
// return true if the key was set
// return false if the key was not set
func (m *ConcurrentMapInt64FdInfo) SetNX(key int, value *fdInfo) bool {
	shard := m.GetShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return true
}

func (m ConcurrentMapInt64FdInfo) Get(key int) (*fdInfo, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m ConcurrentMapInt64FdInfo) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNTConcurrentMapInt64FdInfo; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (m *ConcurrentMapInt64FdInfo) Has(key int) bool {
	shard := m.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

func (m *ConcurrentMapInt64FdInfo) Remove(key int) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

func (m *ConcurrentMapInt64FdInfo) GetAndRemove(key int) (*fdInfo, bool) {
	shard := m.GetShard(key)
	shard.Lock()
	val, ok := shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return val, ok
}

// Returns an iterator which could be used in a for range loop.
func (m ConcurrentMapInt64FdInfo) Iter() <-chan TupleConcurrentMapInt64FdInfo {
	ch := make(chan TupleConcurrentMapInt64FdInfo)
	go func() {
		for _, shard := range m {
			shard.RLock()
			for key, val := range shard.items {
				ch <- TupleConcurrentMapInt64FdInfo{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// Returns a buffered iterator which could be used in a for range loop.
func (m ConcurrentMapInt64FdInfo) IterBuffered() <-chan TupleConcurrentMapInt64FdInfo {
	ch := make(chan TupleConcurrentMapInt64FdInfo, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- TupleConcurrentMapInt64FdInfo{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}