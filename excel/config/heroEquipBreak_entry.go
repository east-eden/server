package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var heroEquipBreakEntries *HeroEquipBreakEntries //HeroEquipBreak.json全局变量

// 属性表 英雄表-拳皇/heroEquipBreak
type HeroEquipBreakEntry struct {
	Id              int32      `json:"id,omitempty" bson:"id"`                           //装备突破属性id
	Level           int32      `json:"level,omitempty" bson:"level"`                     //装备等级
	WeaponAttrValue *AttrValue `json:"weaponAttrValue,omitempty" bson:"weaponAttrValue"` //武器对应突破属性
	HelmetAttrValue *AttrValue `json:"helmetAttrValue,omitempty" bson:"helmetAttrValue"` //头盔对应突破属性
	ChestAttrValue  *AttrValue `json:"chestAttrValue,omitempty" bson:"chestAttrValue"`   //胸甲对应突破属性
	BootsAttrValue  *AttrValue `json:"bootsAttrValue,omitempty" bson:"bootsAttrValue"`   //靴子对应突破属性
}

// HeroEquipBreak.json属性表集合
type HeroEquipBreakEntries struct {
	Rows   map[string]*HeroEquipBreakEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroEquipBreakEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroEquipBreak.json", (*HeroEquipBreakEntries)(nil))
}

func (e *HeroEquipBreakEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroEquipBreakEntries = &HeroEquipBreakEntries{
		Rows:   make(map[string]*HeroEquipBreakEntry),
		Arrays: make([]*HeroEquipBreakEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroEquipBreakEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroEquipBreakEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		heroEquipBreakEntries.Rows[key] = entry
		heroEquipBreakEntries.Arrays = append(heroEquipBreakEntries.Arrays, entry)
	}

	sort.Slice(heroEquipBreakEntries.Arrays, func(i, j int) bool {
		return heroEquipBreakEntries.Arrays[i].Id < heroEquipBreakEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroEquipBreak.json", heroEquipBreakEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroEquipBreakEntries, nil

}

func (e *HeroEquipBreakEntries) GetArrays() any {
	return heroEquipBreakEntries.Arrays
}

func GetHeroEquipBreakEntry(keys ...string) (*HeroEquipBreakEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroEquipBreakEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroEquipBreakSize() int {
	return int(len(heroEquipBreakEntries.Rows))
}

func GetHeroEquipBreakRows() map[string]*HeroEquipBreakEntry {
	return heroEquipBreakEntries.Rows
}

func GetHeroEquipBreakArrays() []*HeroEquipBreakEntry {
	return heroEquipBreakEntries.Arrays
}
