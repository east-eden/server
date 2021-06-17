package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroEnchantEntries *HeroEnchantEntries //HeroEnchant.csv全局变量

// HeroEnchant.csv属性表
type HeroEnchantEntry struct {
	Id              int32   `json:"Id,omitempty"`              // 主键
	PromoteCostId   []int32 `json:"PromoteCostId,omitempty"`   //突破消耗id
	StarupFragments []int32 `json:"StarupFragments,omitempty"` //升星消耗碎片
}

// HeroEnchant.csv属性表集合
type HeroEnchantEntries struct {
	Rows map[int32]*HeroEnchantEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("HeroEnchant.csv", (*HeroEnchantEntries)(nil))
}

func (e *HeroEnchantEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	heroEnchantEntries = &HeroEnchantEntries{
		Rows: make(map[int32]*HeroEnchantEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &HeroEnchantEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		heroEnchantEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetHeroEnchantEntry(id int32) (*HeroEnchantEntry, bool) {
	entry, ok := heroEnchantEntries.Rows[id]
	return entry, ok
}

func GetHeroEnchantSize() int32 {
	return int32(len(heroEnchantEntries.Rows))
}

func GetHeroEnchantRows() map[int32]*HeroEnchantEntry {
	return heroEnchantEntries.Rows
}
