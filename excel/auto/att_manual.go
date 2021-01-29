package auto

import (
	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel"
)

func init() {
	excel.AddEntryManualLoader("att.xlsx", (*AttEntries)(nil))
}

// 手动加载处理
func (e *AttEntries) ManualLoad(excelFileRaw *excel.ExcelFileRaw) error {
	rows := GetAttRows()
	for _, row := range rows {
		row.DmgOfType = make([]int32, define.Att_DmgTypeEnd-define.Att_DmgTypeBegin)
		row.ResOfType = make([]int32, define.Att_ResTypeEnd-define.Att_ResTypeBegin)

		row.DmgOfType[define.Att_DmgPhysics-define.Att_DmgTypeBegin] = row.DmgPhysics
		row.DmgOfType[define.Att_DmgEarth-define.Att_DmgTypeBegin] = row.DmgEarth
		row.DmgOfType[define.Att_DmgWater-define.Att_DmgTypeBegin] = row.DmgWater
		row.DmgOfType[define.Att_DmgFire-define.Att_DmgTypeBegin] = row.DmgFire
		row.DmgOfType[define.Att_DmgWind-define.Att_DmgTypeBegin] = row.DmgWind
		row.DmgOfType[define.Att_DmgTime-define.Att_DmgTypeBegin] = row.DmgTime
		row.DmgOfType[define.Att_DmgSpace-define.Att_DmgTypeBegin] = row.DmgSpace
		row.DmgOfType[define.Att_DmgMirage-define.Att_DmgTypeBegin] = row.DmgMirage
		row.DmgOfType[define.Att_DmgLight-define.Att_DmgTypeBegin] = row.DmgLight
		row.DmgOfType[define.Att_DmgDark-define.Att_DmgTypeBegin] = row.DmgDark

		row.ResOfType[define.Att_ResPhysics-define.Att_ResTypeBegin] = row.ResPhysics
		row.ResOfType[define.Att_ResEarth-define.Att_ResTypeBegin] = row.ResEarth
		row.ResOfType[define.Att_ResWater-define.Att_ResTypeBegin] = row.ResWater
		row.ResOfType[define.Att_ResFire-define.Att_ResTypeBegin] = row.ResFire
		row.ResOfType[define.Att_ResWind-define.Att_ResTypeBegin] = row.ResWind
		row.ResOfType[define.Att_ResTime-define.Att_ResTypeBegin] = row.ResTime
		row.ResOfType[define.Att_ResSpace-define.Att_ResTypeBegin] = row.ResSpace
		row.ResOfType[define.Att_ResMirage-define.Att_ResTypeBegin] = row.ResMirage
		row.ResOfType[define.Att_ResLight-define.Att_ResTypeBegin] = row.ResLight
		row.ResOfType[define.Att_ResDark-define.Att_ResTypeBegin] = row.ResDark
	}
	return nil
}
