package talent

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/logger"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/store"
)

type Talent struct {
	Id    int32               `json:"talent_id" bson:"talent_id" redis:"talent_id"`
	entry *define.TalentEntry `bson:"-" redis:"-"`
}

type TalentManager struct {
	store.StoreObjector
	Owner     define.PluginObj `bson:"-" redis:"-"`
	OwnerId   int64            `bson:"_id" redis:"_id"`
	OwnerType int32            `bson:"owner_type" redis:"owner_type"`
	Talents   []*Talent        `json:"talents" bson:"talents" redis:"talents"`

	sync.RWMutex `bson:"-" redis:"-"`
}

func NewTalentManager(owner define.PluginObj) *TalentManager {
	m := &TalentManager{
		Owner:     owner,
		OwnerId:   owner.GetID(),
		OwnerType: owner.GetType(),
		Talents:   make([]*Talent, 0),
	}

	return m
}

func (m *TalentManager) TableName() string {
	return "talent"
}

func (m *TalentManager) GetObjID() interface{} {
	return m.OwnerId
}

func (m *TalentManager) AfterLoad() {

}

func (m *TalentManager) GetExpire() *time.Timer {
	return nil
}

func (m *TalentManager) LoadFromDB() {
	err := store.GetStore().LoadObjectFromCacheAndDB(store.StoreType_Talent, "_id", m.OwnerId, m)
	if err != nil {
		logger.Error("talent manager load failed:", err)
	}
}

func (m *TalentManager) AddTalent(id int32) error {
	t := &Talent{Id: id, entry: entries.GetTalentEntry(int32(id))}

	if t.entry == nil {
		return fmt.Errorf("add not exist talent entry:%d", id)
	}

	if m.Owner.GetLevel() < t.entry.LevelLimit {
		return fmt.Errorf("level limit:%d", t.entry.LevelLimit)
	}

	// check group_id
	for _, v := range m.Talents {
		if v.Id == t.Id {
			return fmt.Errorf("add existed talent:%d", id)
		}

		// check group_id
		if t.entry.GroupID != -1 && t.entry.GroupID == v.entry.GroupID {
			return fmt.Errorf("talent group_id conflict:%d", t.entry.GroupID)
		}
	}

	m.Talents = append(m.Talents, t)

	return store.GetStore().SaveObjectToCacheAndDB(store.StoreType_Talent, m)
}

func (m *TalentManager) GetTalent(id int32) *Talent {

	for _, v := range m.Talents {
		if v.Id == id {
			return v
		}
	}

	return nil
}

func (m *TalentManager) GetTalentList() []*Talent {

	return m.Talents
}
