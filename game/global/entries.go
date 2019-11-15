package global

import (
	"encoding/json"
	"io/ioutil"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/define"
)

type Entries struct {
	HeroEntries map[int32]*define.HeroEntry `json:"hero_entry"`
	ItemEntries map[int32]*define.ItemEntry `json:"item_entry"`
}

var (
	DefaultEntries *Entries = newEntries()
)

func newEntries() *Entries {
	m := &Entries{
		HeroEntries: make(map[int32]*define.HeroEntry, 0),
		ItemEntries: make(map[int32]*define.ItemEntry, 0),
	}

	newEntry("data/entry/hero_entry.json", m.HeroEntries)
	newEntry("data/entry/item_entry.json", m.ItemEntries)

	return m
}

func newEntry(filePath string, v interface{}) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		logger.Fatal(err)
	}
}

func GetHeroEntry(typeID int32) *define.HeroEntry {
	return DefaultEntries.HeroEntries[typeID]
}

func GetItemEntry(typeID int32) *define.ItemEntry {
	return DefaultEntries.ItemEntries[typeID]
}
