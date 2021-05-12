package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var collectionBoxEntries *CollectionBoxEntries //CollectionBox.xlsx全局变量

// CollectionBox.xlsx属性表
type CollectionBoxEntry struct {
	Id      int32   `json:"Id,omitempty"`      // 主键
	MaxSlot int32   `json:"MaxSlot,omitempty"` //可放置收集品数
	Scores  []int32 `json:"Scores,omitempty"`  //激活评分分段
	Effects []int32 `json:"Effects,omitempty"` //激活评分效果
}

// CollectionBox.xlsx属性表集合
type CollectionBoxEntries struct {
	Rows map[int32]*CollectionBoxEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("CollectionBox.xlsx", (*CollectionBoxEntries)(nil))
}

func (e *CollectionBoxEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	collectionBoxEntries = &CollectionBoxEntries{
		Rows: make(map[int32]*CollectionBoxEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CollectionBoxEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		collectionBoxEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCollectionBoxEntry(id int32) (*CollectionBoxEntry, bool) {
	entry, ok := collectionBoxEntries.Rows[id]
	return entry, ok
}

func GetCollectionBoxSize() int32 {
	return int32(len(collectionBoxEntries.Rows))
}

func GetCollectionBoxRows() map[int32]*CollectionBoxEntry {
	return collectionBoxEntries.Rows
}
