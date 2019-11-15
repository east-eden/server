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
		HeroEntries: make(map[int32]*define.HeroEntry),
		ItemEntries: make(map[int32]*define.ItemEntry),
	}

	var heroEntries define.HeroEntries
	newEntry("../../data/entry/hero_entry.json", &heroEntries)
	for _, v := range heroEntries.Entries {
		if _, ok := m.HeroEntries[v.TypeID]; ok {
			logger.WithFields(logger.Fields{
				"type_id": v.TypeID,
				"file":    "hero_entry.json",
			}).Fatal("adding existed entry")
		}

		m.HeroEntries[v.TypeID] = v
	}

	var itemEntries define.ItemEntries
	newEntry("../../data/entry/item_entry.json", &itemEntries)
	for _, v := range itemEntries.Entries {
		if _, ok := m.ItemEntries[v.TypeID]; ok {
			logger.WithFields(logger.Fields{
				"type_id": v.TypeID,
				"file":    "item_entry.json",
			}).Fatal("adding existed entry")
		}

		m.ItemEntries[v.TypeID] = v
	}

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
