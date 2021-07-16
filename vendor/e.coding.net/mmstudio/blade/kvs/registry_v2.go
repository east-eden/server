package kvs

import (
	"context"
	"fmt"
	"path"
	"time"

	"e.coding.net/mmstudio/blade/kvs/sync2"
)

// The Registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type RegistryV2 interface {
	// Watcher must be provided by every backend.
	WatcherV2
	Register(string, *Entry, ...WriteOption) error
	Deregister(string, *Entry) error
	GetAll(string, ...ReadOption) (Entries, error)
	Close(ctx context.Context)
	Store() Store
}

// Watcher provides watching over a cluster for nodes joining and leaving.
type WatcherV2 interface {
	// Watch the discovery for entry changes.
	// Returns a channel that will receive changes or an error.
	// Providing a non-nil stopCh can be used to stop watching
	Watch(string, ...WatchOption) (<-chan Entries, <-chan error)
	WatchWithNotifier(string, WatchNotifier, ...WatchOption)
}

////////// base on store
func NewRegistryV2(store Store, opts ...RegistryOption) RegistryV2 {
	r := &storeRegistryV2{
		store:    store,
		cc:       NewRegistryOptions(opts...),
		stopChan: make(chan struct{}),
	}
	return r
}

type storeRegistryV2 struct {
	store     Store
	cc        *RegistryOptions
	stopChan  chan struct{}
	closeFlag sync2.AtomicInt32
}

//Store return the underline store current using
func (r *storeRegistryV2) Store() Store { return r.store }

// Watch the store until either there's a store error or we receive a stop request.
// Returns false if we shouldn't attempt watching the store anymore (stop request received).
func (r *storeRegistryV2) watchOnce(root string, stopCh <-chan struct{}, watchCh <-chan []*Pair, notifier WatchNotifier) bool {
	for {
		select {
		case pairs := <-watchCh:
			if pairs == nil {
				return true
			}
			var errGot error
			entries := Entries{}
			for _, pair := range pairs {
				if entry, err := NewEntryFromJson(pair.Value); err == nil {
					entries = append(entries, entry)
				} else {
					errGot = err
				}
			}
			if errGot != nil {
				notifier.NotifyError(root, errGot)
			} else {
				notifier.NotifyUpdate(root, entries)
			}
		case <-stopCh:
			// We were requested to stop watching.
			return false
		}
	}
}

func (r *storeRegistryV2) WatchWithNotifier(root string, notifier WatchNotifier, opts ...WatchOption) {
	go func(rootIn string, notifierIn WatchNotifier, watchOptions *WatchOptions) {
		defer func() {
			if c, ok := notifierIn.(interface {
				Close()
			}); ok {
				c.Close()
			}
		}()
		// Forever: Create a store watch, watch until we get an error and then try again.
		// Will only stop if we receive a stopCh request.
		for {
			// Create the path to watch if it does not exist yet
			exists, err := r.store.Exists(rootIn, watchOptions.KVReadOptions)
			if err != nil {
				notifierIn.NotifyError(rootIn, err)
			}
			if !exists {
				if err := r.store.Put(rootIn, []byte(""), nil); err != nil {
					notifierIn.NotifyError(rootIn, err)
				}
			}

			// Set up a watch.
			watchCh, err := r.store.WatchTree(rootIn, r.stopChan, watchOptions.KVReadOptions)
			if err != nil {
				notifierIn.NotifyError(rootIn, err)
			} else {
				if !r.watchOnce(rootIn, r.stopChan, watchCh, notifierIn) {
					return
				}
			}

			needRetry := true
			select {
			case <-r.stopChan:
				needRetry = false
			default:
			}

			if !needRetry {
				return
			}

			// If we get here it means the store watch channel was closed. This
			// is unexpected so let'r retry later.
			notifierIn.NotifyError(rootIn, fmt.Errorf("unexpected watch error, retry after:%s", watchOptions.Heartbeat))
			time.Sleep(watchOptions.Heartbeat)
		}
	}(root, notifier, NewWatchOptions(opts...))
}

// Watch is exported
func (r *storeRegistryV2) Watch(root string, opts ...WatchOption) (<-chan Entries, <-chan error) {
	on := &watchNotifier{
		entriesChan: make(chan Entries),
		errChan:     make(chan error),
	}
	r.WatchWithNotifier(root, on, opts...)
	return on.entriesChan, on.errChan
}

type watchNotifier struct {
	entriesChan chan Entries
	errChan     chan error
}

func (w *watchNotifier) NotifyUpdate(root string, e Entries) {
	w.entriesChan <- e
}
func (w *watchNotifier) NotifyError(root string, err error) {
	w.errChan <- err
}

func (w *watchNotifier) Close() {
	close(w.entriesChan)
	close(w.errChan)
}

// Register is exported
func (r *storeRegistryV2) Register(root string, entry *Entry, opts ...WriteOption) error {
	wo := *r.cc.BaseStoreWriteOptionsInner
	for _, opt := range opts {
		wo.SetOption(opt)
	}
	return r.store.Put(path.Join(root, entry.Identifier), entry.Marshal(), &wo)
}

// Deregister the given Entry
func (r *storeRegistryV2) Deregister(root string, entry *Entry) error {
	return r.store.Delete(path.Join(root, entry.Identifier))
}

// GetAll get all entries under given dir
func (r *storeRegistryV2) GetAll(root string, opts ...ReadOption) (Entries, error) {
	so := *r.cc.BaseStoreReadOptionsInner
	for _, opt := range opts {
		so.SetOption(opt)
	}
	ret, err := r.store.List(root, &so)
	if err != nil {
		return nil, fmt.Errorf("GetAll with error:%w", err)
	}
	addrSlice := make([][]byte, len(ret))
	for _, pair := range ret {
		addrSlice = append(addrSlice, pair.Value)
	}
	return CreateEntriesFromJson(addrSlice)
}

// Close close the under store and all watchers
func (r *storeRegistryV2) Close(ctx context.Context) {
	if r.closeFlag.CompareAndSwap(0, 1) {
		close(r.stopChan)
		if r.store == nil {
			return
		}
		r.store.Close(ctx)
	}
}
