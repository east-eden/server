package crystal

import (
	"sync"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/internal/att"
)

// crystal create pool
var crystalPool = &sync.Pool{New: newPoolCrystal}

func NewPoolCrystal() *Crystal {
	return crystalPool.Get().(*Crystal)
}

func GetCrystalPool() *sync.Pool {
	return crystalPool
}

func NewCrystal(opts ...Option) *Crystal {
	r := NewPoolCrystal()

	for _, o := range opts {
		o(r.GetOptions())
	}

	return r
}

type CrystalAtt struct {
	AttType  int32 `bson:"att_type" json:"att_type"`
	AttValue int32 `bson:"att_value" json:"att_value"`
}

type Crystal struct {
	Options    `bson:"inline" json:",inline"`
	Atts       [define.Crystal_AttNum]*CrystalAtt `bson:"atts" json:"atts"`
	attManager *att.AttManager                    `bson:"-" json:"-"`
}

func newPoolCrystal() interface{} {
	r := &Crystal{
		Options: DefaultOptions(),
	}

	r.attManager = att.NewAttManager()

	return r
}

func (r *Crystal) GetExpire() *time.Timer {
	return nil
}

func (r *Crystal) GetStoreIndex() int64 {
	return r.Options.OwnerId
}

func (r *Crystal) GetOptions() *Options {
	return &r.Options
}

func (r *Crystal) GetType() int32 {
	return define.Plugin_Crystal
}

func (r *Crystal) GetID() int64 {
	return r.Options.Id
}

func (r *Crystal) GetOwnerID() int64 {
	return r.Options.OwnerId
}

func (r *Crystal) GetTypeID() int32 {
	return r.Options.TypeId
}

func (r *Crystal) GetEquipObj() int64 {
	return r.Options.EquipObj
}

func (r *Crystal) GetAttManager() *att.AttManager {
	return r.attManager
}

func (r *Crystal) GetAtt(idx int32) *CrystalAtt {
	if idx < 0 || idx >= define.Crystal_AttNum {
		return nil
	}

	return r.Atts[idx]
}

func (r *Crystal) SetAtt(idx int32, att *CrystalAtt) {
	if idx < 0 || idx >= define.Crystal_AttNum {
		return
	}

	r.Atts[idx] = att
}

func (r *Crystal) CalcAtt() {
	r.attManager.Reset()

	var n int32
	for n = 0; n < define.Crystal_AttNum; n++ {
		att := r.Atts[n]
		if att == nil {
			continue
		}

		r.attManager.ModBaseAtt(int(att.AttType), att.AttValue)
	}

	r.attManager.CalcAtt()
}
