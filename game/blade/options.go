package blade

import (
	"github.com/yokaiio/yokai_server/define"
)

type Option func(*Options)

// blade options
type Options struct {
	Id        int64              `bson:"_id" redis:"_id"`
	OwnerId   int64              `bson:"owner_id" redis:"owner_id"`
	OwnerType int32              `bson:"owner_type" redis:"owner_type"`
	TypeId    int32              `bson:"type_id" redis:"type_id"`
	Exp       int64              `bson:"exp" redis:"exp"`
	Level     int32              `bson:"level" redis:"level"`
	Entry     *define.BladeEntry `bson:"-" redis:"-"`
}

func DefaultOptions() *Options {
	return &Options{
		Id:        -1,
		OwnerId:   -1,
		OwnerType: -1,
		TypeId:    -1,
		Exp:       0,
		Level:     1,
		Entry:     nil,
	}
}

func Id(id int64) Option {
	return func(o *Options) {
		o.Id = id
	}
}

func OwnerId(id int64) Option {
	return func(o *Options) {
		o.OwnerId = id
	}
}

func OwnerType(tp int32) Option {
	return func(o *Options) {
		o.OwnerType = tp
	}
}

func TypeId(id int32) Option {
	return func(o *Options) {
		o.TypeId = id
	}
}

func Exp(exp int64) Option {
	return func(o *Options) {
		o.Exp = exp
	}
}

func Level(level int32) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func Entry(entry *define.BladeEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}
