package collection

import (
	"sync"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
)

// collection create pool
var collectionPool = &sync.Pool{New: newPoolCollection}

func GetCollectionPool() *sync.Pool {
	return collectionPool
}

func NewCollection() *Collection {
	return collectionPool.Get().(*Collection)
}

type Collection struct {
	Options `bson:"inline" json:",inline"`
}

func newPoolCollection() interface{} {
	return &Collection{}
}

func (c *Collection) Init(opts ...Option) {
	c.Options = DefaultOptions()

	for _, o := range opts {
		o(c.GetOptions())
	}
}

func (c *Collection) GetOptions() *Options {
	return &c.Options
}

func (c *Collection) GenCollectionPB() *pbGlobal.Collection {
	pb := &pbGlobal.Collection{
		TypeId: c.TypeId,
		Active: c.Active,
		Star:   int32(c.Star),
	}

	return pb
}
