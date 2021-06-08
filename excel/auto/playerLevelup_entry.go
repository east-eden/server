package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var playerLevelupEntries *PlayerLevelupEntries //PlayerLevelup.csv全局变量

// PlayerLevelup.csv属性表
type PlayerLevelupEntry struct {
	Id     int32 `json:"Id,omitempty"`     // 主键
	Exp    int32 `json:"Exp,omitempty"`    //账号经验
	LootId int32 `json:"LootId,omitempty"` //升级奖励
}

// PlayerLevelup.csv属性表集合
type PlayerLevelupEntries struct {
	Rows map[int32]*PlayerLevelupEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("PlayerLevelup.csv", (*PlayerLevelupEntries)(nil))
}

func (e *PlayerLevelupEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	playerLevelupEntries = &PlayerLevelupEntries{
		Rows: make(map[int32]*PlayerLevelupEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &PlayerLevelupEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		playerLevelupEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetPlayerLevelupEntry(id int32) (*PlayerLevelupEntry, bool) {
	entry, ok := playerLevelupEntries.Rows[id]
	return entry, ok
}

func GetPlayerLevelupSize() int32 {
	return int32(len(playerLevelupEntries.Rows))
}

func GetPlayerLevelupRows() map[int32]*PlayerLevelupEntry {
	return playerLevelupEntries.Rows
}
