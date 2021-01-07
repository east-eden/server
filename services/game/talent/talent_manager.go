package talent

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/store"
)

type Talent struct {
	Id    int32             `bson:"talent_id" json:"talent_id"`
	entry *auto.TalentEntry `bson:"-" json:"-"`
}

type TalentManager struct {
	store.StoreObjector
	Owner     define.PluginObj `bson:"-" json:"-"`
	OwnerId   int64            `bson:"_id" json:"_id"`
	OwnerType int32            `bson:"owner_type" json:"owner_type"`
	Talents   []*Talent        `bson:"talents" json:"talents"`

	sync.RWMutex `bson:"-" json:"-"`
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

func (m *TalentManager) GetObjID() int64 {
	return m.OwnerId
}

func (m *TalentManager) GetOwnerID() int64 {
	return -1
}

func (m *TalentManager) AfterLoad() error {
	return nil
}

func (m *TalentManager) GetExpire() *time.Timer {
	return nil
}

func (m *TalentManager) LoadFromDB() error {
	err := store.GetStore().LoadObject(define.StoreType_Talent, m.OwnerId, m)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("TalentManager LoadFromDB: %w", err)
	}

	return nil
}

func (m *TalentManager) AddTalent(id int32) error {
	t := &Talent{Id: id}
	t.entry, _ = auto.GetTalentEntry(id)

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
		if t.entry.GroupId != -1 && t.entry.GroupId == v.entry.GroupId {
			return fmt.Errorf("talent group_id conflict:%d", t.entry.GroupId)
		}
	}

	m.Talents = append(m.Talents, t)

	return store.GetStore().SaveObject(define.StoreType_Talent, m)
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
