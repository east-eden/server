package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var collectionEntries *CollectionEntries //Collection.xlsx全局变量

// Collection.xlsx属性表
type CollectionEntry struct {
	Id                int32   `json:"Id,omitempty"`                // 主键
	Quality           int32   `json:"Quality,omitempty"`           //品质
	QuestId           int32   `json:"QuestId,omitempty"`           //激活任务id
	BaseScore         int32   `json:"BaseScore,omitempty"`         //基础评分
	FragmentCompose   int32   `json:"FragmentCompose,omitempty"`   //合成所需碎片数量
	FragmentTransform int32   `json:"FragmentTransform,omitempty"` //重复获得返还碎片数
	StarAttIds        []int32 `json:"StarAttIds,omitempty"`        //星级对应的属性id
	WakeupAttId       int32   `json:"WakeupAttId,omitempty"`       //觉醒属性id
	WakeupCostId      int32   `json:"WakeupCostId,omitempty"`      //觉醒消耗id
	ActiveCostId      int32   `json:"ActiveCostId,omitempty"`      //激活消耗id
	ExhibitionAttId   int32   `json:"ExhibitionAttId,omitempty"`   //放置属性id
	StarCostFragments []int32 `json:"StarCostFragments,omitempty"` //升星消耗碎片数量
}

// Collection.xlsx属性表集合
type CollectionEntries struct {
	Rows map[int32]*CollectionEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Collection.xlsx", (*CollectionEntries)(nil))
}

func (e *CollectionEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	collectionEntries = &CollectionEntries{
		Rows: make(map[int32]*CollectionEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CollectionEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		collectionEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCollectionEntry(id int32) (*CollectionEntry, bool) {
	entry, ok := collectionEntries.Rows[id]
	return entry, ok
}

func GetCollectionSize() int32 {
	return int32(len(collectionEntries.Rows))
}

func GetCollectionRows() map[int32]*CollectionEntry {
	return collectionEntries.Rows
}
