package auto

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var skillTrackEntries *SkillTrackEntries //SkillTrack.xlsx全局变量

// SkillTrack.xlsx属性表
type SkillTrackEntry struct {
	Id           int32  `json:"Id,omitempty"`           // 主键
	TrackID      int32  `json:"TrackID,omitempty"`      // 多主键之一
	TrackType    int32  `json:"TrackType,omitempty"`    //track类型
	StartTime    number `json:"StartTime,omitempty"`    //开始时间
	DurationTime number `json:"DurationTime,omitempty"` //持续时间
	EffectIndex  int32  `json:"EffectIndex,omitempty"`  //第几个effect
	AnimName     string `json:"AnimName,omitempty"`     //动作名称
	FxName       string `json:"FxName,omitempty"`       //特效名称
	HitAnimName  string `json:"HitAnimName,omitempty"`  //受击动作
	HitFxName    string `json:"HitFxName,omitempty"`    //受击特效
	HitFxSlot    string `json:"HitFxSlot,omitempty"`    //受击特效插槽
	HitStopTime  number `json:"HitStopTime,omitempty"`  //受击动作
}

// SkillTrack.xlsx属性表集合
type SkillTrackEntries struct {
	Rows map[string]*SkillTrackEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("SkillTrack.xlsx", (*SkillTrackEntries)(nil))
}

func (e *SkillTrackEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	skillTrackEntries = &SkillTrackEntries{
		Rows: make(map[string]*SkillTrackEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillTrackEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		key := fmt.Sprintf("%d+%d", entry.Id, entry.TrackID)
		skillTrackEntries.Rows[key] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetSkillTrackEntry(keys ...int32) (*SkillTrackEntry, bool) {
	keyName := make([]string, 0, len(keys))
	for _, key := range keys {
		keyName = append(keyName, strconv.Itoa(int(key)))
	}

	finalKey := strings.Join(keyName, "+")
	entry, ok := skillTrackEntries.Rows[finalKey]
	return entry, ok
}

func GetSkillTrackSize() int32 {
	return int32(len(skillTrackEntries.Rows))
}

func GetSkillTrackRows() map[string]*SkillTrackEntry {
	return skillTrackEntries.Rows
}
