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
	HeroEntries   map[int32]*define.HeroEntry
	ItemEntries   map[int32]*define.ItemEntry
	TokenEntries  map[int32]*define.TokenEntry
	TalentEntries map[int32]*define.TalentEntry
	BladeEntries  map[int32]*define.BladeEntry

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

func GetPlayerLevelupEntry(id int32) *define.PlayerLevelupEntry {
	return DefaultEntries.PlayerLevelupEntries[id]
}

func newEntries() *Entries {
	var wg utils.WaitGroupWrapper

	m := &Entries{
		HeroEntries:   make(map[int32]*define.HeroEntry),
		ItemEntries:   make(map[int32]*define.ItemEntry),
		TokenEntries:  make(map[int32]*define.TokenEntry),
		TalentEntries: make(map[int32]*define.TalentEntry),
		BladeEntries:  make(map[int32]*define.BladeEntry),

		PlayerLevelupEntries: make(map[int32]*define.PlayerLevelupEntry),
	}

	// hero_entry.json
	wg.Wrap(func() {
		var heroEntries struct {
			Entries []*define.HeroEntry `json:"hero_entry"`
		}
		readEntry("hero_entry.json", &heroEntries, m.HeroEntries)
	})

	// item_entry.json
	wg.Wrap(func() {
		var itemEntries struct {
			Entries []*define.ItemEntry `json:"item_entry"`
		}
		readEntry("item_entry.json", &itemEntries, m.ItemEntries)
	})

	// token_entry.json
	wg.Wrap(func() {
		var tokenEntries struct {
			Entries []*define.TokenEntry `json:"token_entry"`
		}
		readEntry("token_entry.json", &tokenEntries, m.TokenEntries)
	})

	// talent_entry.json
	wg.Wrap(func() {
		var talentEntries struct {
			Entries []*define.TalentEntry `json:"talent_entry"`
		}
		readEntry("talent_entry.json", &talentEntries, m.TalentEntries)
	})

	// blade_entry.json
	wg.Wrap(func() {
		var bladeEntries struct {
			Entries []*define.BladeEntry `json:"blade_entry"`
		}
		readEntry("blade_entry.json", &bladeEntries, m.BladeEntries)
	})

	// player_levelup_entry.json
	wg.Wrap(func() {
		var playerLevelupEntries struct {
			Entries []*define.PlayerLevelupEntry `json:"player_levelup_entry"`
		}
		readEntry("player_levelup_entry.json", &playerLevelupEntries, m.PlayerLevelupEntries)
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
		entryField := reflect.ValueOf(v).Elem().FieldByName("Entries")
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
