package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var shopGoodsEntries *ShopGoodsEntries //ShopGoods.json全局变量

// 属性表 商店商品表-拳皇/ShopGoods
type ShopGoodsEntry struct {
	Id                    int32      `json:"id,omitempty" bson:"id"`                                       //##商品id
	Name                  string     `json:"name,omitempty" bson:"name"`                                   //商品名
	QuickBuy              bool       `json:"quickBuy,omitempty" bson:"quickBuy"`                           //是否快速购买
	Unit                  int32      `json:"unit,omitempty" bson:"unit"`                                   //购买单位数量
	Rarity                Rarity     `json:"rarity,omitempty" bson:"rarity"`                               //品质
	CostItems             []*IdCount `json:"costItems,omitempty" bson:"costItems"`                         //购买价格
	Limit                 int32      `json:"limit,omitempty" bson:"limit"`                                 //限量
	LevelLimit            int32      `json:"levelLimit,omitempty" bson:"levelLimit"`                       //等级限制
	ShowConditionId       int32      `json:"showConditionId,omitempty" bson:"showConditionId"`             //显示商品条件id
	IsHide                bool       `json:"isHide,omitempty" bson:"isHide"`                               //隐藏商品
	GainItems             []*IdCount `json:"gainItems,omitempty" bson:"gainItems"`                         //固定获得道具【Ver4.29 必须给功能解锁道具本体】
	LootId                int32      `json:"lootId,omitempty" bson:"lootId"`                               //奖励掉落id
	ShowAwardItems_       []*IdCount `json:"showAwardItems_,omitempty" bson:"showAwardItems_"`             //展示奖励道具列表【必给的道具本体不展示则不填】
	IsVideo               bool       `json:"isVideo,omitempty" bson:"isVideo"`                             //是否播放视频
	Reset                 string     `json:"reset,omitempty" bson:"reset"`                                 //重置(cron)
	DoubleCharge          bool       `json:"doubleCharge,omitempty" bson:"doubleCharge"`                   //是否首充双倍
	PlatformGoodsId       string     `json:"platformGoodsId,omitempty" bson:"platformGoodsId"`             //平台商品ID
	OriCostItems          []*IdCount `json:"oriCostItems,omitempty" bson:"oriCostItems"`                   //展示原价消耗道具列表
	Discount              float32    `json:"discount,omitempty" bson:"discount"`                           //显示折扣
	UnlockAllianceLevel   int32      `json:"unlockAllianceLevel,omitempty" bson:"unlockAllianceLevel"`     //商品解锁联盟等级
	UnlockVipLevel        int32      `json:"unlockVipLevel,omitempty" bson:"unlockVipLevel"`               //商品解锁vip等级
	UnlockBattlePassPoint *StringVal `json:"unlockBattlePassPoint,omitempty" bson:"unlockBattlePassPoint"` //解锁对应的battlepass积分
	PreGoodsId            int32      `json:"preGoodsId,omitempty" bson:"preGoodsId"`                       //前置商品id
}

// ShopGoods.json属性表集合
type ShopGoodsEntries struct {
	Rows   map[int32]*ShopGoodsEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ShopGoodsEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ShopGoods.json", (*ShopGoodsEntries)(nil))
}

func (e *ShopGoodsEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	shopGoodsEntries = &ShopGoodsEntries{
		Rows:   make(map[int32]*ShopGoodsEntry),
		Arrays: make([]*ShopGoodsEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ShopGoodsEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return shopGoodsEntries, err
		}

		shopGoodsEntries.Rows[entry.Id] = entry
		shopGoodsEntries.Arrays = append(shopGoodsEntries.Arrays, entry)
	}

	sort.Slice(shopGoodsEntries.Arrays, func(i, j int) bool {
		return shopGoodsEntries.Arrays[i].Id < shopGoodsEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("shopGoods.json", shopGoodsEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return shopGoodsEntries, nil

}

func (e *ShopGoodsEntries) GetArrays() any {
	return shopGoodsEntries.Arrays
}

func GetShopGoodsEntry(id int32) (*ShopGoodsEntry, bool) {
	entry, ok := shopGoodsEntries.Rows[id]
	return entry, ok
}

func GetShopGoodsSize() int {
	return int(len(shopGoodsEntries.Rows))
}

func GetShopGoodsRows() map[int32]*ShopGoodsEntry {
	return shopGoodsEntries.Rows
}

func GetShopGoodsArrays() []*ShopGoodsEntry {
	return shopGoodsEntries.Arrays
}
