package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var sceneEntries *SceneEntries //Scene.xlsx全局变量

// Scene.xlsx属性表
type SceneEntry struct {
	Id     int32           `json:"Id,omitempty"`     // 主键
	Type   int32           `json:"Type,omitempty"`   //场景类型
	Radius decimal.Decimal `json:"Radius,omitempty"` //战斗区域半径
}

// Scene.xlsx属性表集合
type SceneEntries struct {
	Rows map[int32]*SceneEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Scene.xlsx", (*SceneEntries)(nil))
}

func (e *SceneEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	sceneEntries = &SceneEntries{
		Rows: make(map[int32]*SceneEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SceneEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		sceneEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetSceneEntry(id int32) (*SceneEntry, bool) {
	entry, ok := sceneEntries.Rows[id]
	return entry, ok
}

func GetSceneSize() int32 {
	return int32(len(sceneEntries.Rows))
}

func GetSceneRows() map[int32]*SceneEntry {
	return sceneEntries.Rows
}
