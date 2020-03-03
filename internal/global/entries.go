package global

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Entries struct {
	HeroEntries     map[int32]*define.HeroEntry
	ItemEntries     map[int32]*define.ItemEntry
	TokenEntries    map[int32]*define.TokenEntry
	TalentEntries   map[int32]*define.TalentEntry
	BladeEntries    map[int32]*define.BladeEntry
	CostLootEntries map[int32]*define.CostLootEntry

	PlayerLevelupEntries map[int32]*define.PlayerLevelupEntry
}

var (
	DefaultEntries *Entries = newEntries()
)

func GetHeroEntry(id int32) *define.HeroEntry {
	return DefaultEntries.HeroEntries[id]
}

func GetItemEntry(id int32) *define.ItemEntry {
	return DefaultEntries.ItemEntries[id]
}

func GetTokenEntry(id int32) *define.TokenEntry {
	return DefaultEntries.TokenEntries[id]
}

func GetTalentEntry(id int32) *define.TalentEntry {
	return DefaultEntries.TalentEntries[id]
}

func GetBladeEntry(id int32) *define.BladeEntry {
	return DefaultEntries.BladeEntries[id]
}

func GetCostLootEntry(id int32) *define.CostLootEntry {
	return DefaultEntries.CostLootEntries[id]
}

func GetPlayerLevelupEntry(id int32) *define.PlayerLevelupEntry {
	return DefaultEntries.PlayerLevelupEntries[id]
}

func newEntries() *Entries {
	var wg utils.WaitGroupWrapper

	m := &Entries{
		HeroEntries:     make(map[int32]*define.HeroEntry),
		ItemEntries:     make(map[int32]*define.ItemEntry),
		TokenEntries:    make(map[int32]*define.TokenEntry),
		TalentEntries:   make(map[int32]*define.TalentEntry),
		BladeEntries:    make(map[int32]*define.BladeEntry),
		CostLootEntries: make(map[int32]*define.CostLootEntry),

		PlayerLevelupEntries: make(map[int32]*define.PlayerLevelupEntry),
	}

	// hero_entry.json
	wg.Wrap(func() {
		entry := make([]*define.HeroEntry, 0)
		readEntry("hero_entry.json", &entry, m.HeroEntries)
	})

	// ItemConfig.json
	wg.Wrap(func() {
		entry := make([]*define.ItemEntry, 0)
		readEntry("ItemConfig.json", &entry, m.ItemEntries)
	})

	// token_entry.json
	wg.Wrap(func() {
		entry := make([]*define.TokenEntry, 0)
		readEntry("token_entry.json", &entry, m.TokenEntries)
	})

	// talent_entry.json
	wg.Wrap(func() {
		entry := make([]*define.TalentEntry, 0)
		readEntry("talent_entry.json", &entry, m.TalentEntries)
	})

	// blade_entry.json
	wg.Wrap(func() {
		entry := make([]*define.BladeEntry, 0)
		readEntry("blade_entry.json", &entry, m.BladeEntries)
	})

	// cost_loot_entry.json
	wg.Wrap(func() {
		entry := make([]*define.CostLootEntry, 0)
		readEntry("cost_loot_entry.json", &entry, m.CostLootEntries)
	})

	// player_levelup_entry.json
	wg.Wrap(func() {
		entry := make([]*define.PlayerLevelupEntry, 0)
		readEntry("player_levelup_entry.json", &entry, m.PlayerLevelupEntries)
	})

	wg.Wait()
	return m
}

// read entries(v) to map(m)
func readEntry(filePath string, v interface{}, m interface{}) {
	absPath := strings.Join([]string{"../../config/entry/", filePath}, "")
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		logger.Fatal(err)
	}

	tp := reflect.TypeOf(v)
	if tp.Kind() == reflect.Ptr || tp.Kind() == reflect.Struct {
		entryField := reflect.ValueOf(v).Elem()
		mapValue := reflect.ValueOf(m)

		for n := 0; n < entryField.Len(); n++ {
			elem := entryField.Index(n)
			key := elem.Elem().FieldByName("ID")

			// if key existed
			keyExist := mapValue.MapIndex(key)
			if keyExist.IsValid() {
				logger.WithFields(logger.Fields{
					"file": absPath,
					"id":   key.Int(),
				}).Fatal("error loading entry")
			}

			mapValue.SetMapIndex(key, elem)
		}
	} else {
		logger.WithFields(logger.Fields{
			"file": absPath,
		}).Fatal("skip reading entry")
	}
}
