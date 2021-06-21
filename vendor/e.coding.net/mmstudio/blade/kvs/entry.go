package kvs

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"net"
	"time"
)

const MetaKeyCreated = "kvs-created"

var EntryEqualOperator = func(e1 *Entry, e2 *Entry) bool {
	return e1.Host == e2.Host && e1.Port == e2.Port && e1.Identifier == e2.Identifier
}

type Entry struct {
	Host       string
	Port       string
	Identifier string
	Metadata   map[string]string
	Raw        []byte
}

func NewEntry(host string, port string, opts ...EntryOption) *Entry {
	e := &Entry{
		Host: host,
		Port: port,
	}
	for _, opt := range opts {
		_ = opt(e)
	}
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[MetaKeyCreated] = time.Now().String()
	if e.Identifier == "" {
		e.Identifier = xid.New().String()
	}
	return e
}

func NewEntryFromJson(jsonb []byte) (*Entry, error) {
	e := &Entry{}
	err := json.Unmarshal(jsonb, e)
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	return e, err
}

// Marshal the entry into []byte as value in the registry
func (e *Entry) Marshal() []byte {
	b, _ := json.Marshal(e)
	return b
}

// Equals returns true if cmp contains the same data.
func (e *Entry) Equals(cmp *Entry) bool { return EntryEqualOperator(e, cmp) }

// Address returns the identifier an entry.
func (e *Entry) Address() string { return net.JoinHostPort(e.Host, e.Port) }

// String returns the string form of an entry.
func (e *Entry) String() string { return fmt.Sprintf("%s_%s", e.Identifier, e.Address()) }

func (e *Entry) SetMetadata(key, value string) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
}

func (e *Entry) GetMetadata(key string) string {
	if e.Metadata == nil {
		return ""
	}
	return e.Metadata[key]
}

func (e *Entry) SetOption(opt EntryOption) {
	_ = opt(e)
}

func (e *Entry) ApplyOption(opts ...EntryOption) {
	for _, opt := range opts {
		_ = opt(e)
	}
}

func (e *Entry) GetSetOption(opt EntryOption) EntryOption {
	return opt(e)
}

type EntryOption func(cc *Entry) EntryOption

func WithEntryIdentifier(v string) EntryOption {
	return func(cc *Entry) EntryOption {
		previous := cc.Identifier
		cc.Identifier = v
		return WithEntryIdentifier(previous)
	}
}
func WithEntryRaw(v []byte) EntryOption {
	return func(cc *Entry) EntryOption {
		previous := cc.Raw
		cc.Raw = v
		return WithEntryRaw(previous)
	}
}
func WithEntryMetadata(v map[string]string) EntryOption {
	return func(cc *Entry) EntryOption {
		previous := cc.Metadata
		cc.Metadata = v
		return WithEntryMetadata(previous)
	}
}
