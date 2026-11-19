package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var vipLevelEntries *VipLevelEntries //VipLevel.json全局变量

// 属性表 vip表-拳皇/VipLevel
type VipLevelEntry struct {
	Id                    int32      `json:"id,omitempty" bson:"id"`                                       //##vip等级
	NextExp               int32      `json:"nextExp,omitempty" bson:"nextExp"`                             //升到下一级所需经验值
	DailyFreeShopGoods    int32      `json:"dailyFreeShopGoods,omitempty" bson:"dailyFreeShopGoods"`       //每日免费特权礼包商品id
	DailyPremiumShopGoods int32      `json:"dailyPremiumShopGoods,omitempty" bson:"dailyPremiumShopGoods"` //每日尊享特权礼包商品id
	PreminumCostItems     []*IdCount `json:"preminumCostItems,omitempty" bson:"preminumCostItems"`         //每日尊享特权礼包折扣价
	RightsList            []int32    `json:"rightsList,omitempty" bson:"rightsList"`                       //vip专属权益
	AddVal                []int32    `json:"addVal,omitempty" bson:"addVal"`                               //权益等级
	DailyAwards           []*IdCount `json:"dailyAwards,omitempty" bson:"dailyAwards"`                     //vip每日奖励列表
}

// VipLevel.json属性表集合
type VipLevelEntries struct {
	Rows   map[int32]*VipLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*VipLevelEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("VipLevel.json", (*VipLevelEntries)(nil))
}

func (e *VipLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	vipLevelEntries = &VipLevelEntries{
		Rows:   make(map[int32]*VipLevelEntry),
		Arrays: make([]*VipLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &VipLevelEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return vipLevelEntries, err
		}

		vipLevelEntries.Rows[entry.Id] = entry
		vipLevelEntries.Arrays = append(vipLevelEntries.Arrays, entry)
	}

	sort.Slice(vipLevelEntries.Arrays, func(i, j int) bool {
		return vipLevelEntries.Arrays[i].Id < vipLevelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("vipLevel.json", vipLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return vipLevelEntries, nil

}

func (e *VipLevelEntries) GetArrays() any {
	return vipLevelEntries.Arrays
}

func GetVipLevelEntry(id int32) (*VipLevelEntry, bool) {
	entry, ok := vipLevelEntries.Rows[id]
	return entry, ok
}

func GetVipLevelSize() int {
	return int(len(vipLevelEntries.Rows))
}

func GetVipLevelRows() map[int32]*VipLevelEntry {
	return vipLevelEntries.Rows
}

func GetVipLevelArrays() []*VipLevelEntry {
	return vipLevelEntries.Arrays
}
