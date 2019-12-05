package talent

import (
	"encoding/json"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/global"
)

type TalentManager struct {
	OwnerID    int64     `gorm:"type:bigint(20);primary_key;column:owner_id;index:owner_id;default:0;not null"`
	TalentJson string    `gorm:"type:varchar(5120);column:talent_json"`
	Talents    []*Talent `json:"talents"`

	ds *db.Datastore
	sync.RWMutex
}

func NewTalentManager(ownerID int64, ds *db.Datastore) *TalentManager {
	m := &TalentManager{
		OwnerID: ownerID,
		ds:      ds,
		Talents: make([]*Talent, 0),
	}

	// init talents
	//m.initTalents()

	return m
}

func (m *TalentManager) TableName() string {
	return "talent"
}

func Migrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(TalentManager{})
}

func (m *TalentManager) initTalents() {
	//for n := 0; n < define.Talent_End; n++ {
	//m.Talents = append(m.Talents, &Talent{
	//ID:      int32(n),
	//Value:   0,
	//MaxHold: 100000000,
	//entry:   global.GetTalentEntry(int32(n)),
	//})
	//}
}

func (m *TalentManager) LoadFromDB() {
	m.ds.ORM().Find(&m)

	// unmarshal json to talent value
	if len(m.TalentJson) > 0 {
		m.Lock()
		err := json.Unmarshal([]byte(m.TalentJson), &m.Talents)
		m.Unlock()
		if err != nil {
			logger.Error("unmarshal talent json failed:", err)
		}
	}

	// init entry
	m.RLock()
	for _, v := range m.Talents {
		v.entry = global.GetTalentEntry(int32(v.ID))
	}
	m.RUnlock()
}

func (m *TalentManager) Save() error {
	m.RLock()
	data, err := json.Marshal(m.Talents)
	m.RUnlock()
	if err != nil {
		return fmt.Errorf("json marshal failed:", err)
	}

	m.TalentJson = string(data)
	m.ds.ORM().Save(m)
	return nil
}

func (m *TalentManager) AddTalent(id int32) error {
	t := &Talent{ID: id, entry: global.GetTalentEntry(int32(id))}

	if t.entry == nil {
		return fmt.Errorf("add not exist talent entry:%d", id)
	}

	m.Lock()
	defer m.Unlock()

	bFixPrev := false
	bFixMutex := true
	for _, v := range m.Talents {
		if v.ID == t.ID {
			return fmt.Errorf("add existed talent:%d", id)
		}

		// check prev_id
		if t.entry.PrevID == 0 || t.entry.PrevID == v.ID {
			bFixPrev = true
		}

		// check mutex
		if t.entry.MutexID == 0 || t.entry.MutexID == v.entry.MutexID {
			bFixMutex = false
		}
	}

	if bFixPrev && bFixMutex {
		m.Talents = append(m.Talents, t)
	}

	return nil
}

func (m *TalentManager) GetTalent(id int32) *Talent {
	m.RLock()
	defer m.RUnlock()

	for _, v := range m.Talents {
		if v.ID == id {
			return v
		}
	}

	return nil
}

func (m *TalentManager) GetTalentList() []*Talent {
	m.RLock()
	defer m.RUnlock()

	return m.Talents
}
