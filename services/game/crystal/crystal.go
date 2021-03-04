package crystal

import (
	"sync"

	"bitbucket.org/funplus/server/define"
)

// crystal create pool
var crystalPool = &sync.Pool{New: func() interface{} { return &Crystal{} }}

func NewPoolCrystal() *Crystal {
	return crystalPool.Get().(*Crystal)
}

func GetCrystalPool() *sync.Pool {
	return crystalPool
}

func NewCrystal(opts ...Option) *Crystal {
	c := NewPoolCrystal()
	c.Options = DefaultOptions()
	c.attManager = NewCrystalAttManager(c)

	for _, o := range opts {
		o(c.GetOptions())
	}

	return c
}

// 晶石主属性
type CrystalMainAtt struct {
	AttRepoId    int32 `bson:"att_repo_id" json:"att_repo_id"`       // 主属性库id
	AttRandRatio int32 `bson:"att_rand_ratio" json:"att_rand_ratio"` // 主属性随机区间系数
}

// 晶石副属性
type CrystalViceAtt struct {
	AttId        int32 `bson:"att_id" json:"att_id"`                 // 副属性id
	AttRandRatio int32 `bson:"att_rand_ratio" json:"att_rand_ratio"` // 副属性随机区间
}

type Crystal struct {
	Options    `bson:"inline" json:",inline"`
	MainAtt    CrystalMainAtt     `bson:"main_att"`
	ViceAtts   []CrystalViceAtt   `bson:"vice_atts" json:"vice_atts"`
	attManager *CrystalAttManager `bson:"-" json:"-"`
}

func (c *Crystal) GetOptions() *Options {
	return &c.Options
}

func (c *Crystal) GetType() int32 {
	return define.Plugin_Crystal
}

func (c *Crystal) GetID() int64 {
	return c.Id
}

func (c *Crystal) GetOwnerID() int64 {
	return c.OwnerId
}

func (c *Crystal) GetTypeID() int32 {
	return c.TypeId
}

func (c *Crystal) GetEquipObj() int64 {
	return c.EquipObj
}

func (c *Crystal) GetAttManager() *CrystalAttManager {
	return c.attManager
}
