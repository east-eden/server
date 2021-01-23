package hero

import (
	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
)

type Option func(*Options)

// hero options
type Options struct {
	define.HeroInfo
	Entry *auto.HeroEntry `bson:"-" json:"-"`
}

func DefaultOptions() Options {
	return Options{
		HeroInfo: define.HeroInfo{
			Id:        -1,
			OwnerId:   -1,
			OwnerType: -1,
			TypeId:    -1,
			Exp:       0,
			Level:     1,
		},
		Entry: nil,
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

func Entry(entry *auto.HeroEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}
