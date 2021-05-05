package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var stageEntries *StageEntries //Stage.xlsx全局变量

// Stage.xlsx属性表
type StageEntry struct {
	Id                int32             `json:"Id,omitempty"`                // 主键
	PrevStageId       int32             `json:"PrevStageId,omitempty"`       //前置关卡id
	ConditionId       int32             `json:"ConditionId,omitempty"`       //解锁条件id
	CostStrength      int32             `json:"CostStrength,omitempty"`      //消耗体力
	ChapterId         int32             `json:"ChapterId,omitempty"`         //所属章节id
	Type              int32             `json:"Type,omitempty"`              //关卡类型
	WinCondition      int32             `json:"WinCondition,omitempty"`      //胜利条件
	LostCondition     int32             `json:"LostCondition,omitempty"`     //失败条件
	StarConditionIds  []int32           `json:"StarConditionIds,omitempty"`  //星级挑战条件列表
	FirstRewardLootId int32             `json:"FirstRewardLootId,omitempty"` //首次通关掉落id
	RewardLootId      int32             `json:"RewardLootId,omitempty"`      //通关掉落id
	DailyTimes        int32             `json:"DailyTimes,omitempty"`        //每日挑战次数
	SceneId           int32             `json:"SceneId,omitempty"`           //场景id
	UsPosition1       []decimal.Decimal `json:"UsPosition1,omitempty"`       //单位1生成点
	UsPosition2       []decimal.Decimal `json:"UsPosition2,omitempty"`       //单位2生成点
	UsPosition3       []decimal.Decimal `json:"UsPosition3,omitempty"`       //单位3生成点
	UsPosition4       []decimal.Decimal `json:"UsPosition4,omitempty"`       //单位4生成点
	Rotate1           decimal.Decimal   `json:"Rotate1,omitempty"`           //单位1朝向
	Rotate2           decimal.Decimal   `json:"Rotate2,omitempty"`           //单位2朝向
	Rotate3           decimal.Decimal   `json:"Rotate3,omitempty"`           //单位3朝向
	Rotate4           decimal.Decimal   `json:"Rotate4,omitempty"`           //单位4朝向
	UsInitalCom1      decimal.Decimal   `json:"UsInitalCom1,omitempty"`      //初始COM_1
	UsInitalCom2      decimal.Decimal   `json:"UsInitalCom2,omitempty"`      //初始COM_2
	UsInitalCom3      decimal.Decimal   `json:"UsInitalCom3,omitempty"`      //初始COM_3
	UsInitalCom4      decimal.Decimal   `json:"UsInitalCom4,omitempty"`      //初始COM_4
	UnitGroupId       int32             `json:"UnitGroupId,omitempty"`       //怪物组id
}

// Stage.xlsx属性表集合
type StageEntries struct {
	Rows map[int32]*StageEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Stage.xlsx", (*StageEntries)(nil))
}

func (e *StageEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	stageEntries = &StageEntries{
		Rows: make(map[int32]*StageEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &StageEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		stageEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetStageEntry(id int32) (*StageEntry, bool) {
	entry, ok := stageEntries.Rows[id]
	return entry, ok
}

func GetStageSize() int32 {
	return int32(len(stageEntries.Rows))
}

func GetStageRows() map[int32]*StageEntry {
	return stageEntries.Rows
}
