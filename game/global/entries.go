package global

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Entries struct {
	HeroEntries map[int32]*define.HeroEntry
	ItemEntries map[int32]*define.ItemEntry
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

func newEntries() *Entries {
	var wg utils.WaitGroupWrapper

	m := &Entries{
		HeroEntries: make(map[int32]*define.HeroEntry),
		ItemEntries: make(map[int32]*define.ItemEntry),
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
