package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	sceneEntries	*SceneEntries	//scene.xlsx全局变量

// scene.xlsx属性表
type SceneEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//场景id      
	Desc      	string    	`json:"Desc,omitempty"`	//场景描述      
	Type      	string    	`json:"Type,omitempty"`	//场景类型      
	UnitGroupId	int       	`json:"UnitGroupId,omitempty"`	//场景怪物组id   
}

// scene.xlsx属性表集合
type SceneEntries struct {
	Rows      	map[int]*SceneEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("scene.xlsx", sceneEntries)
}

func (e *SceneEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	sceneEntries = &SceneEntries{
		Rows: make(map[int]*SceneEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SceneEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	sceneEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetSceneEntry(id int) (*SceneEntry, bool) {
	entry, ok := sceneEntries.Rows[id]
	return entry, ok
}

func  GetSceneSize() int {
	return len(sceneEntries.Rows)
}

