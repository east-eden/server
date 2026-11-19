package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var shopEntries *ShopEntries //Shop.json全局变量

// 属性表 商店商品表-拳皇/Shop
type ShopEntry struct {
	Id                   int32     `json:"id,omitempty" bson:"id"`                                     //##商店id
	Name                 string    `json:"name,omitempty" bson:"name"`                                 //商店名
	OpenLv               int32     `json:"openLv,omitempty" bson:"openLv"`                             //商城开启等级
	Goods                []int32   `json:"goods,omitempty" bson:"goods"`                               //商品id列表
	GoodsLayer           [][]int32 `json:"goodsLayer,omitempty" bson:"goodsLayer"`                     //商品分层
	SlotGoodsRandomLoots []int32   `json:"slotGoodsRandomLoots,omitempty" bson:"slotGoodsRandomLoots"` //slot对应随机商品掉落组
	Reset                string    `json:"reset,omitempty" bson:"reset"`                               //商品重置(cron)
	SoldOutReset         int32     `json:"soldOutReset,omitempty" bson:"soldOutReset"`                 //售尽后自动重置次数
	RefreshLimit         int32     `json:"refreshLimit,omitempty" bson:"refreshLimit"`                 //刷新次数上限
	RefGoodsCost         *IdCount  `json:"refGoodsCost,omitempty" bson:"refGoodsCost"`                 //刷新商品消耗
	OnePackageGoodsId    int32     `json:"onePackageGoodsId,omitempty" bson:"onePackageGoodsId"`       //一键打包商品ID
	ShowCostGoodsIds     []int32   `json:"showCostGoodsIds,omitempty" bson:"showCostGoodsIds"`         //展示消耗道具列表
}

// Shop.json属性表集合
type ShopEntries struct {
	Rows   map[int32]*ShopEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ShopEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Shop.json", (*ShopEntries)(nil))
}

func (e *ShopEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	shopEntries = &ShopEntries{
		Rows:   make(map[int32]*ShopEntry),
		Arrays: make([]*ShopEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ShopEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return shopEntries, err
		}

		shopEntries.Rows[entry.Id] = entry
		shopEntries.Arrays = append(shopEntries.Arrays, entry)
	}

	sort.Slice(shopEntries.Arrays, func(i, j int) bool {
		return shopEntries.Arrays[i].Id < shopEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("shop.json", shopEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return shopEntries, nil

}

func (e *ShopEntries) GetArrays() any {
	return shopEntries.Arrays
}

func GetShopEntry(id int32) (*ShopEntry, bool) {
	entry, ok := shopEntries.Rows[id]
	return entry, ok
}

func GetShopSize() int {
	return int(len(shopEntries.Rows))
}

func GetShopRows() map[int32]*ShopEntry {
	return shopEntries.Rows
}

func GetShopArrays() []*ShopEntry {
	return shopEntries.Arrays
}
