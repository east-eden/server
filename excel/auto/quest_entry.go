package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var questEntries *QuestEntries //Quest.csv全局变量

// Quest.csv属性表
type QuestEntry struct {
	Id           int32   `json:"Id,omitempty"`           // 主键
	Type         int32   `json:"Type,omitempty"`         //任务类型
	RefreshType  int32   `json:"RefreshType,omitempty"`  //刷新方式
	ObjTypes     []int32 `json:"ObjTypes,omitempty"`     //任务目标类型列表
	ObjParams1   []int32 `json:"ObjParams1,omitempty"`   //任务目标参数1
	ObjParams2   []int32 `json:"ObjParams2,omitempty"`   //任务目标参数2
	ObjCount     []int32 `json:"ObjCount,omitempty"`     //任务目标计数
	RewardLootId int32   `json:"RewardLootId,omitempty"` //任务奖励id
}

// Quest.csv属性表集合
type QuestEntries struct {
	Rows map[int32]*QuestEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Quest.csv", (*QuestEntries)(nil))
}

func (e *QuestEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	questEntries = &QuestEntries{
		Rows: make(map[int32]*QuestEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &QuestEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		questEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetQuestEntry(id int32) (*QuestEntry, bool) {
	entry, ok := questEntries.Rows[id]
	return entry, ok
}

func GetQuestSize() int32 {
	return int32(len(questEntries.Rows))
}

func GetQuestRows() map[int32]*QuestEntry {
	return questEntries.Rows
}
