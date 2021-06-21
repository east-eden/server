package kvs

import (
	"context"
)

// The Registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	// Watcher must be provided by every backend.
	Watcher
	Register(*Entry) error
	Deregister(*Entry) error
	GetAll() (Entries, error)
	Close(ctx context.Context)
	Path() string
	Store() Store
}

// Watcher provides watching over a cluster for nodes joining and leaving.
type Watcher interface {
	// Watch the discovery for entry changes.
	// Returns a channel that will receive changes or an error.
	// Providing a non-nil stopCh can be used to stop watching
	Watch(...WatchOption) (<-chan Entries, <-chan error)
	WatchWithNotifier(WatchNotifier, ...WatchOption)
}

// WatchNotifier
type WatchNotifier interface {
	NotifyUpdate(string, Entries)
	NotifyError(string, error)
}

////////// base on store
func NewRegistry(path string, store Store, opts ...RegistryOption) Registry {
	r := &storeRegistry{
		path:            path,
		storeRegistryV2: NewRegistryV2(store, opts...),
	}
	return r
}

// storeRegistry is exported
type storeRegistry struct {
	path            string
	storeRegistryV2 RegistryV2
}

func (r *storeRegistry) Path() string              { return r.path }
func (r *storeRegistry) Register(e *Entry) error   { return r.storeRegistryV2.Register(r.path, e) }
func (r *storeRegistry) Deregister(e *Entry) error { return r.storeRegistryV2.Deregister(r.path, e) }
func (r *storeRegistry) GetAll() (Entries, error)  { return r.storeRegistryV2.GetAll(r.path) }
func (r *storeRegistry) Close(ctx context.Context) { r.storeRegistryV2.Close(ctx) }
func (r *storeRegistry) Store() Store              { return r.storeRegistryV2.Store() }
func (r *storeRegistry) Watch(opts ...WatchOption) (<-chan Entries, <-chan error) {
	return r.storeRegistryV2.Watch(r.path, opts...)
}
func (r *storeRegistry) WatchWithNotifier(wn WatchNotifier, opts ...WatchOption) {
	r.storeRegistryV2.WatchWithNotifier(r.path, wn, opts...)
}
