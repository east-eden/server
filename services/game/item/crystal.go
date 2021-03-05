package item

// 晶石属性
type CrystalAtt struct {
	AttRepoId    int32 `bson:"att_repo_id" json:"att_repo_id"`       // 属性库id
	AttRandRatio int32 `bson:"att_rand_ratio" json:"att_rand_ratio"` // 属性随机区间系数
}

// 晶石
type Crystal struct {
	Item           `bson:"inline" json:",inline"`
	CrystalOptions `bson:"inline" json:",inline"`
	MainAtt        CrystalAtt         `bson:"main_att" json:"main_att"`
	ViceAtts       []CrystalAtt       `bson:"vice_atts" json:"vice_atts"`
	attManager     *CrystalAttManager `bson:"-" json:"-"`
}

func (c *Crystal) InitCrystal(opts ...CrystalOption) {
	// 主属性初始化
	c.MainAtt.AttRepoId = -1
	c.MainAtt.AttRandRatio = 0

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
