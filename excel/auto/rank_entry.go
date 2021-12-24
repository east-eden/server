package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var rankEntries *RankEntries //Rank.csv全局变量

// Rank.csv属性表
type RankEntry struct {
	Id          int32 `json:"Id,omitempty"`          // 主键
	RefreshType int32 `json:"RefreshType,omitempty"` //刷新方式
	Local       bool  `json:"Local,omitempty"`       //是否仅为本服排行榜
	Desc        bool  `json:"Desc,omitempty"`        //是否降序排序
}

// Rank.csv属性表集合
type RankEntries struct {
	Rows map[int32]*RankEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Rank.csv", (*RankEntries)(nil))
}

func (e *RankEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	rankEntries = &RankEntries{
		Rows: make(map[int32]*RankEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &RankEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		rankEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetRankEntry(id int32) (*RankEntry, bool) {
	entry, ok := rankEntries.Rows[id]
	return entry, ok
}

func GetRankSize() int32 {
	return int32(len(rankEntries.Rows))
}

func GetRankRows() map[int32]*RankEntry {
	return rankEntries.Rows
}
