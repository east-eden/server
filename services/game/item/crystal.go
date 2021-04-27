package item

import (
	"bitbucket.org/funplus/server/define"
	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	"github.com/shopspring/decimal"
)

// 晶石属性
type CrystalAtt struct {
	AttRepoId    int32           `bson:"att_repo_id" json:"att_repo_id"`       // 属性库id
	AttRandRatio decimal.Decimal `bson:"att_rand_ratio" json:"att_rand_ratio"` // 属性随机区间系数
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

	for _, o := range opts {
		o(&c.CrystalOptions)
	}
}

func (c *Crystal) GetAttManager() *CrystalAttManager {
	return c.attManager
}

func (c *Crystal) GenCrystalDataPB() *pbCommon.CrystalData {
	pb := &pbCommon.CrystalData{
		Level:      int32(c.Level),
		Exp:        int32(c.Exp),
		CrystalObj: c.CrystalObj,
		MainAtt: &pbCommon.CrystalAtt{
			AttRepoId:    c.MainAtt.AttRepoId,
			AttRandRatio: int32(c.MainAtt.AttRandRatio.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()),
		},
		ViceAtts: make([]*pbCommon.CrystalAtt, 0, len(c.ViceAtts)),
	}

	for _, att := range c.ViceAtts {
		pb.ViceAtts = append(pb.ViceAtts, &pbCommon.CrystalAtt{
			AttRepoId:    att.AttRepoId,
			AttRandRatio: int32(att.AttRandRatio.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()),
		})
	}

	return pb
}

func (c *Crystal) GenCrystalPB() *pbCommon.Crystal {
	pb := &pbCommon.Crystal{
		Item:        c.GenItemPB(),
		CrystalData: c.GenCrystalDataPB(),
	}

	return pb
}
