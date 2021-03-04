package item

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

// 晶石
type Crystal struct {
	Item           `bson:"inline" json:",inline"`
	CrystalOptions `bson:"inline" json:",inline"`
	MainAtt        CrystalMainAtt     `bson:"main_att"`
	ViceAtts       []CrystalViceAtt   `bson:"vice_atts" json:"vice_atts"`
	attManager     *CrystalAttManager `bson:"-" json:"-"`
}

func (c *Crystal) InitCrystal(opts ...CrystalOption) {
	for _, o := range opts {
		o(&c.CrystalOptions)
	}
}

func (c *Crystal) OnDelete() {
	c.CrystalObj = -1
	c.Item.OnDelete()
}

func (c *Crystal) GetAttManager() *CrystalAttManager {
	return c.attManager
}
