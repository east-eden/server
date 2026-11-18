package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var chapterEntries *ChapterEntries //Chapter.json全局变量

// 属性表 关卡表-拳皇/Chapter
type ChapterEntry struct {
	Id int32 `json:"id,omitempty" bson:"id"` //##id
}

// Chapter.json属性表集合
type ChapterEntries struct {
	Rows   map[int32]*ChapterEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ChapterEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Chapter.json", (*ChapterEntries)(nil))
}

func (e *ChapterEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	chapterEntries = &ChapterEntries{
		Rows:   make(map[int32]*ChapterEntry),
		Arrays: make([]*ChapterEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ChapterEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return chapterEntries, err
		}

		chapterEntries.Rows[entry.Id] = entry
		chapterEntries.Arrays = append(chapterEntries.Arrays, entry)
	}

	sort.Slice(chapterEntries.Arrays, func(i, j int) bool {
		return chapterEntries.Arrays[i].Id < chapterEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("chapter.json", chapterEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return chapterEntries, nil

}

func (e *ChapterEntries) GetArrays() any {
	return chapterEntries.Arrays
}

func GetChapterEntry(id int32) (*ChapterEntry, bool) {
	entry, ok := chapterEntries.Rows[id]
	return entry, ok
}

func GetChapterSize() int {
	return int(len(chapterEntries.Rows))
}

func GetChapterRows() map[int32]*ChapterEntry {
	return chapterEntries.Rows
}

func GetChapterArrays() []*ChapterEntry {
	return chapterEntries.Arrays
}
