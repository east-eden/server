package collection

import (
	"bitbucket.org/funplus/server/excel/auto"
)

type Option func(*Options)

// collection options
type Options struct {
	TypeId  int32                 `bson:"type_id" json:"type_id"`
	OwnerId int64                 `bson:"owner_id" json:"owner_id"`
	Active  bool                  `bson:"active" json:"active"`
	Star    int8                  `bson:"star" json:"star"`
	Entry   *auto.CollectionEntry `bson:"-" json:"-"`
}

func DefaultOptions() Options {
	return Options{
		TypeId:  -1,
		OwnerId: -1,
		Active:  false,
		Star:    0,
		Entry:   nil,
	}
}

func TypeId(id int32) Option {
	return func(o *Options) {
		o.TypeId = id
	}
}

func OwnerId(id int64) Option {
	return func(o *Options) {
		o.OwnerId = id
	}
}

func Star(star int8) Option {
	return func(o *Options) {
		o.Star = star
	}
}

func Entry(entry *auto.CollectionEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}
