package selector

import (
	"e.coding.net/mmstudio/blade/kvs"
	"errors"
)

type Selector interface {
	// Select returns a node
	Select(service string, opts ...SelectOption) (*kvs.Entry, error)
	// Mark sets the success/error against a node
	Mark(service string, node *kvs.Entry, err error)
	// Reset returns state back to zero for a service
	Reset(service string)
	// Close renders the selector unusable
	Close() error
	// Name of the selector
	String() string
}

// Filter is used to filter a service during the selection process
type Filter func(kvs.Entries) kvs.Entries

// Strategy is a selection strategy e.g random, round robin
type Strategy func(kvs.Entries) (*kvs.Entry, error)

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)
