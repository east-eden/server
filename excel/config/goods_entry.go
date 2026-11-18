package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var goodsEntries *GoodsEntries //Goods.json全局变量

// 属性表 基础配置-拳皇/goods
type GoodsEntry struct {
	Id       int32        `json:"id,omitempty" bson:"id"`             //##道具 id
	Name     string       `json:"name,omitempty" bson:"name"`         //道具名称
	Desc     string       `json:"desc,omitempty" bson:"desc"`         //描述
	Rarity   Rarity       `json:"rarity,omitempty" bson:"rarity"`     //稀有度
	Type     GoodsType    `json:"type,omitempty" bson:"type"`         //道具类型
	SubType  GoodsSubType `json:"subType,omitempty" bson:"subType"`   //子类型
	StackMax int64        `json:"stackMax,omitempty" bson:"stackMax"` //道具数量上限
	Param    string       `json:"param,omitempty" bson:"param"`       //通用参数
}

// Goods.json属性表集合
type GoodsEntries struct {
	Rows   map[int32]*GoodsEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*GoodsEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Goods.json", (*GoodsEntries)(nil))
}

func (e *GoodsEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	goodsEntries = &GoodsEntries{
		Rows:   make(map[int32]*GoodsEntry),
		Arrays: make([]*GoodsEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &GoodsEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return goodsEntries, err
		}

		goodsEntries.Rows[entry.Id] = entry
		goodsEntries.Arrays = append(goodsEntries.Arrays, entry)
	}

	sort.Slice(goodsEntries.Arrays, func(i, j int) bool {
		return goodsEntries.Arrays[i].Id < goodsEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("goods.json", goodsEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return goodsEntries, nil

}

func (e *GoodsEntries) GetArrays() any {
	return goodsEntries.Arrays
}

func GetGoodsEntry(id int32) (*GoodsEntry, bool) {
	entry, ok := goodsEntries.Rows[id]
	return entry, ok
}

func GetGoodsSize() int {
	return int(len(goodsEntries.Rows))
}

func GetGoodsRows() map[int32]*GoodsEntry {
	return goodsEntries.Rows
}

func GetGoodsArrays() []*GoodsEntry {
	return goodsEntries.Arrays
}
