package hero

import "github.com/yokaiio/yokai_server/define"

type Option func(*Options)

// hero options
type Options struct {
	Id        int64             `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerId   int64             `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	OwnerType int32             `gorm:"type:int(10);column:owner_type;index:owner_type;default:-1;not null" bson:"owner_type"`
	TypeId    int32             `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Exp       int64             `gorm:"type:bigint(20);column:exp;default:0;not null" bson:"exp"`
	Level     int32             `gorm:"type:int(10);column:level;default:1;not null" bson:"level"`
	Entry     *define.HeroEntry `gorm:"-" bson:"-"`
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

func Entry(entry *define.HeroEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}
