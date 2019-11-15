package global

import (
	"encoding/json"
	"io/ioutil"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/define"
)

type Entrys struct {
	HeroEntrys map[int32]*define.HeroEntry `json:"hero_entry"`
	ItemEntrys map[int32]*define.ItemEntry `json:"item_entry"`
}

var (
	DefaultEntrys *Entrys = newEntrys()
)

func newEntrys() *Entrys {
	m := &Entrys{
		HeroEntry: make(map[int32]*define.HeroEntry, 0),
		ItemENtry: make(map[int32]*define.ItemEntry, 0),
	}

	newEntry("data/entry/hero_entry.json", m.HeroEntry)
	newEntry("data/entry/item_entry.json", m.ItemEntry)

	return m
}

func newEntry(filePath string, m map[interface{}]interface{}) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(data, m)
	if err != nil {
		logger.Fatal(err)
	}
}

func GetEntrys() *Entrys {
	return DefaultEntrys
}
