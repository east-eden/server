package entries

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/utils"
)

type Entries struct {
	HeroEntries         map[int32]*define.HeroEntry
	ItemEntries         map[int32]*define.ItemEntry
	EquipEnchantEntries map[int32]*define.EquipEnchantEntry
	TokenEntries        map[int32]*define.TokenEntry
	TalentEntries       map[int32]*define.TalentEntry
	BladeEntries        map[int32]*define.BladeEntry
	RuneEntries         map[int32]*define.RuneEntry
	RuneSuitEntries     map[int32]*define.RuneSuitEntry
	CostLootEntries     map[int32]*define.CostLootEntry
	AttEntries          map[int32]*define.AttEntry
	SceneEntries        map[int32]*define.SceneEntry
	UnitGroupEntries    map[int32]*define.UnitGroupEntry
	UnitEntries         map[int32]*define.UnitEntry

	PlayerLevelupEntries map[int32]*define.PlayerLevelupEntry
}

var (
	DefaultEntries *Entries
)

func InitEntries() {
	DefaultEntries = newEntries()
}

func GetHeroEntry(id int32) *define.HeroEntry {
	return DefaultEntries.HeroEntries[id]
}

func GetItemEntry(id int32) *define.ItemEntry {
	return DefaultEntries.ItemEntries[id]
}

func GetEquipEnchantEntry(id int32) *define.EquipEnchantEntry {
	return DefaultEntries.EquipEnchantEntries[id]
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

func GetRuneEntry(id int32) *define.RuneEntry {
	return DefaultEntries.RuneEntries[id]
}

func GetRuneSuitEntry(id int32) *define.RuneSuitEntry {
	return DefaultEntries.RuneSuitEntries[id]
}

func GetCostLootEntry(id int32) *define.CostLootEntry {
	return DefaultEntries.CostLootEntries[id]
}

func GetAttEntry(id int32) *define.AttEntry {
	return DefaultEntries.AttEntries[id]
}

func GetSceneEntry(id int32) *define.SceneEntry {
	return DefaultEntries.SceneEntries[id]
}

func GetUnitGroupEntry(id int32) *define.UnitGroupEntry {
	return DefaultEntries.UnitGroupEntries[id]
}

func GetUnitEntry(id int32) *define.UnitEntry {
	return DefaultEntries.UnitEntries[id]
}

func GetPlayerLevelupEntry(id int32) *define.PlayerLevelupEntry {
	return DefaultEntries.PlayerLevelupEntries[id]
}

func newEntries() *Entries {
	var wg utils.WaitGroupWrapper

	m := &Entries{
		HeroEntries:         make(map[int32]*define.HeroEntry),
		ItemEntries:         make(map[int32]*define.ItemEntry),
		EquipEnchantEntries: make(map[int32]*define.EquipEnchantEntry),
		TokenEntries:        make(map[int32]*define.TokenEntry),
		TalentEntries:       make(map[int32]*define.TalentEntry),
		BladeEntries:        make(map[int32]*define.BladeEntry),
		RuneEntries:         make(map[int32]*define.RuneEntry),
		RuneSuitEntries:     make(map[int32]*define.RuneSuitEntry),
		CostLootEntries:     make(map[int32]*define.CostLootEntry),
		AttEntries:          make(map[int32]*define.AttEntry),
		SceneEntries:        make(map[int32]*define.SceneEntry),
		UnitGroupEntries:    make(map[int32]*define.UnitGroupEntry),
		UnitEntries:         make(map[int32]*define.UnitEntry),

		PlayerLevelupEntries: make(map[int32]*define.PlayerLevelupEntry),
	}

	// HeroConfig.json
	wg.Wrap(func() {
		entry := make([]*define.HeroEntry, 0)
		if err := readEntry("HeroConfig.json", &entry, m.HeroEntries); err != nil {
			log.Fatal(err)
		}
	})

	// ItemConfig.json
	wg.Wrap(func() {
		entry := make([]*define.ItemEntry, 0)
		if err := readEntry("ItemConfig.json", &entry, m.ItemEntries); err != nil {
			log.Fatal(err)
		}
	})

	// EquipEnchantConfig.json
	wg.Wrap(func() {
		entry := make([]*define.EquipEnchantEntry, 0)
		if err := readEntry("EquipEnchantConfig.json", &entry, m.EquipEnchantEntries); err != nil {
			log.Fatal(err)
		}
	})

	// TokenConfig.json
	wg.Wrap(func() {
		entry := make([]*define.TokenEntry, 0)
		if err := readEntry("TokenConfig.json", &entry, m.TokenEntries); err != nil {
			log.Fatal(err)
		}
	})

	// talent_entry.json
	wg.Wrap(func() {
		entry := make([]*define.TalentEntry, 0)
		if err := readEntry("talent_entry.json", &entry, m.TalentEntries); err != nil {
			log.Fatal(err)
		}
	})

	// blade_entry.json
	wg.Wrap(func() {
		entry := make([]*define.BladeEntry, 0)
		if err := readEntry("blade_entry.json", &entry, m.BladeEntries); err != nil {
			log.Fatal(err)
		}
	})

	// RuneConfig.json
	wg.Wrap(func() {
		entry := make([]*define.RuneEntry, 0)
		if err := readEntry("RuneConfig.json", &entry, m.RuneEntries); err != nil {
			log.Fatal(err)
		}
	})

	// RuneSuitConfig.json
	wg.Wrap(func() {
		entry := make([]*define.RuneSuitEntry, 0)
		if err := readEntry("RuneSuitConfig.json", &entry, m.RuneSuitEntries); err != nil {
			log.Fatal(err)
		}
	})

	// cost_loot_entry.json
	wg.Wrap(func() {
		entry := make([]*define.CostLootEntry, 0)
		if err := readEntry("CostLootConfig.json", &entry, m.CostLootEntries); err != nil {
			log.Fatal(err)
		}
	})

	// AttConfig.json
	wg.Wrap(func() {
		entry := make([]*define.AttEntry, 0)
		if err := readEntry("AttConfig.json", &entry, m.AttEntries); err != nil {
			log.Fatal(err)
		}
	})

	// SceneConfig.json
	wg.Wrap(func() {
		entry := make([]*define.SceneEntry, 0)
		if err := readEntry("SceneConfig.json", &entry, m.SceneEntries); err != nil {
			log.Fatal(err)
		}
	})

	// UnitGroupConfig.json
	wg.Wrap(func() {
		entry := make([]*define.UnitGroupEntry, 0)
		if err := readEntry("UnitGroupConfig.json", &entry, m.UnitGroupEntries); err != nil {
			log.Fatal(err)
		}
	})

	// UnitConfig.json
	wg.Wrap(func() {
		entry := make([]*define.UnitEntry, 0)
		if err := readEntry("UnitConfig.json", &entry, m.UnitEntries); err != nil {
			log.Fatal(err)
		}
	})

	// player_levelup_entry.json
	wg.Wrap(func() {
		entry := make([]*define.PlayerLevelupEntry, 0)
		if err := readEntry("PlayerLevelupConfig.json", &entry, m.PlayerLevelupEntries); err != nil {
			log.Fatal(err)
		}
	})

	wg.Wait()
	return m
}

// read entries(v) to map(m)
func readEntry(filePath string, v interface{}, m interface{}) error {
	absPath := strings.Join([]string{"config/entry/", filePath}, "")
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("readEntry failed: ioutil.ReadFile<%s> error:%w", absPath, err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("readEntry failed: json.Unmarshal<%s> error:%w", absPath, err)
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
				return fmt.Errorf("readEntry failed: key dunplicate<%d>, file<%s>", key.Int(), absPath)
			}

			mapValue.SetMapIndex(key, elem)
		}
	} else {
		return fmt.Errorf("readEntry failed: skip reading entry<%s>", absPath)
	}

	return nil
}
