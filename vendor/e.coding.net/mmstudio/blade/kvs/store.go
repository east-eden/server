package kvs

import (
	"context"
	"errors"
)

type NewStoreFunc func(opts ...StoreOption) (Store, error)

// NewStore creates an instance of store
func NewStore(init NewStoreFunc, opts ...StoreOption) (Store, error) {
	return init(opts...)
}

var (
	// ErrKeyModified is thrown during an atomic operation if the index does not match the one in the store
	ErrKeyModified = errors.New("unable to complete atomic operation, key modified")
	// ErrKeyNotFound is thrown when the key is not found in the store during a Get operation
	ErrKeyNotFound = errors.New("key not found in store")
	// ErrPreviousNotSpecified is thrown when the previous value is not specified for an atomic operation
	ErrPreviousNotSpecified = errors.New("previous K/V pair should be provided for the Atomic operation")
	// ErrKeyExists is thrown when the previous value exists in the case of an AtomicPut
	ErrKeyExists = errors.New("previous K/V pair exists, cannot complete Atomic operation")
)

// Store represents the backend K/V storage
type Store interface {
	// Put a value at the specified key
	Put(key string, value []byte, options *WriteOptions) error

	// Get a value given its key
	Get(key string, options *ReadOptions) (*Pair, error)

	// Delete the value at the specified key
	Delete(key string) error

	// Verify if a Key exists in the store
	Exists(key string, options *ReadOptions) (bool, error)

	// Watch for changes on a key
	Watch(key string, stopCh <-chan struct{}, options *ReadOptions) (<-chan *Pair, error)

	// WatchTree watches for changes on child nodes under
	// a given directory
	WatchTree(directory string, stopCh <-chan struct{}, options *ReadOptions) (<-chan []*Pair, error)

	// NewLock creates a lock for a given key.
	// The returned Locker is not held and must be acquired
	// with `.Lock`. The Value is optional.
	NewLock(key string, options *LockOptions) (Locker, error)

	// List the content of a given prefix
	List(directory string, options *ReadOptions) ([]*Pair, error)

	// DeleteTree deletes a range of keys under a given directory
	DeleteTree(directory string) error

	// Atomic CAS operation on a single value.
	// Pass previous = nil to create a new key.
	AtomicPut(key string, value []byte, previous *Pair, options *WriteOptions) (bool, *Pair, error)

	// Atomic delete of a single value
	AtomicDelete(key string, previous *Pair) (bool, error)

	// Close the store connection
	Close(ctx context.Context)
}

// KVPair represents {Key, Value, Lastindex} tuple
type Pair struct {
	Key       string
	Value     []byte
	LastIndex uint64
}

// Locker provides locking mechanism on top of the store.
// Similar to `sync.Lock` except it may return errors.
type Locker interface {
	Lock(stopChan chan struct{}) (<-chan struct{}, error)
	Unlock() error
}
