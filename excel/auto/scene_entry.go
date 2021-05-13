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
	Id              int32             `json:"Id,omitempty"`              // 主键
	Resource        string            `json:"Resource,omitempty"`        //场景资源
	Type            int32             `json:"Type,omitempty"`            //场景类型
	Radius          decimal.Decimal   `json:"Radius,omitempty"`          //战斗区域半径
	UsPosition1     []decimal.Decimal `json:"UsPosition1,omitempty"`     //单位1生成点
	UsPosition2     []decimal.Decimal `json:"UsPosition2,omitempty"`     //单位2生成点
	UsPosition3     []decimal.Decimal `json:"UsPosition3,omitempty"`     //单位3生成点
	UsPosition4     []decimal.Decimal `json:"UsPosition4,omitempty"`     //单位4生成点
	UsRotate1       decimal.Decimal   `json:"UsRotate1,omitempty"`       //单位1朝向
	UsRotate2       decimal.Decimal   `json:"UsRotate2,omitempty"`       //单位2朝向
	UsRotate3       decimal.Decimal   `json:"UsRotate3,omitempty"`       //单位3朝向
	UsRotate4       decimal.Decimal   `json:"UsRotate4,omitempty"`       //单位4朝向
	UsInitalCom1    decimal.Decimal   `json:"UsInitalCom1,omitempty"`    //初始COM_1
	UsInitalCom2    decimal.Decimal   `json:"UsInitalCom2,omitempty"`    //初始COM_2
	UsInitalCom3    decimal.Decimal   `json:"UsInitalCom3,omitempty"`    //初始COM_3
	UsInitalCom4    decimal.Decimal   `json:"UsInitalCom4,omitempty"`    //初始COM_4
	EnemyPosition1  []decimal.Decimal `json:"EnemyPosition1,omitempty"`  //单位1生成点
	EnemyPosition2  []decimal.Decimal `json:"EnemyPosition2,omitempty"`  //单位2生成点
	EnemyPosition3  []decimal.Decimal `json:"EnemyPosition3,omitempty"`  //单位3生成点
	EnemyPosition4  []decimal.Decimal `json:"EnemyPosition4,omitempty"`  //单位4生成点
	EnemyPosition5  []decimal.Decimal `json:"EnemyPosition5,omitempty"`  //单位5生成点
	EnemyPosition6  []decimal.Decimal `json:"EnemyPosition6,omitempty"`  //单位6生成点
	EnemyRotate1    decimal.Decimal   `json:"EnemyRotate1,omitempty"`    //单位1朝向
	EnemyRotate2    decimal.Decimal   `json:"EnemyRotate2,omitempty"`    //单位2朝向
	EnemyRotate3    decimal.Decimal   `json:"EnemyRotate3,omitempty"`    //单位3朝向
	EnemyRotate4    decimal.Decimal   `json:"EnemyRotate4,omitempty"`    //单位4朝向
	EnemyRotate5    decimal.Decimal   `json:"EnemyRotate5,omitempty"`    //单位5朝向
	EnemyRotate6    decimal.Decimal   `json:"EnemyRotate6,omitempty"`    //单位6朝向
	EnemyInitalCom1 decimal.Decimal   `json:"EnemyInitalCom1,omitempty"` //初始COM_1
	EnemyInitalCom2 decimal.Decimal   `json:"EnemyInitalCom2,omitempty"` //初始COM_2
	EnemyInitalCom3 decimal.Decimal   `json:"EnemyInitalCom3,omitempty"` //初始COM_3
	EnemyInitalCom4 decimal.Decimal   `json:"EnemyInitalCom4,omitempty"` //初始COM_4
	EnemyInitalCom5 decimal.Decimal   `json:"EnemyInitalCom5,omitempty"` //初始COM_5
	EnemyInitalCom6 decimal.Decimal   `json:"EnemyInitalCom6,omitempty"` //初始COM_6
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
