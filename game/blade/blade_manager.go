package blade

import (
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type BladeManager struct {
	Owner    define.PluginObj
	mapBlade map[int64]*Blade

	ds *db.Datastore
	sync.RWMutex
	wg utils.WaitGroupWrapper
}

func NewBladeManager(obj define.PluginObj, ds *db.Datastore) *BladeManager {
	m := &BladeManager{
		Owner:    obj,
		ds:       ds,
		mapBlade: make(map[int64]*Blade, 0),
	}

	return m
}

func Migrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(Blade{})
}

// interface of cost_loot
func (m *BladeManager) GetCostLootType() int32 {
	return define.CostLoot_Blade
}

func (m *BladeManager) CanCost(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) DoCost(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) CanGain(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) GainLoot(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) LoadFromDB() {
	list := make([]*Blade, 0)
	m.ds.ORM().Where("owner_id = ?", m.Owner.GetID()).Find(&list)

	for _, v := range list {
		m.newDBBlade(v)
	}

	m.wg.Wait()
}

func (m *BladeManager) newEntryBlade(entry *define.BladeEntry) *Blade {
	if entry == nil {
		logger.Error("newEntryBlade with nil BladeEntry")
		return nil
	}

	id, err := utils.NextID(define.Plugin_Blade)
	if err != nil {
		logger.Error(err)
		return nil
	}

	blade := newBlade(id, m.Owner, m.ds)
	blade.OwnerID = m.Owner.GetID()
	blade.TypeID = entry.ID
	blade.Entry = entry

	m.mapBlade[blade.GetID()] = blade

	return blade
}

func (m *BladeManager) newDBBlade(b *Blade) *Blade {
	blade := newBlade(b.GetID(), m.Owner, m.ds)
	blade.OwnerID = m.Owner.GetID()
	blade.TypeID = b.TypeID
	blade.Entry = global.GetBladeEntry(b.TypeID)

	m.mapBlade[blade.GetID()] = blade

	// load from db
	m.wg.Wrap(blade.LoadFromDB)

	return blade
}

func (m *BladeManager) GetBlade(id int64) *Blade {
	return m.mapBlade[id]
}

func (m *BladeManager) GetBladeNums() int {
	return len(m.mapBlade)
}

func (m *BladeManager) GetBladeList() []*Blade {
	list := make([]*Blade, 0)

	m.RLock()
	for _, v := range m.mapBlade {
		list = append(list, v)
	}
	m.RUnlock()

	return list
}

func (m *BladeManager) AddBlade(typeID int32) *Blade {
	bladeEntry := global.GetBladeEntry(typeID)
	blade := m.newEntryBlade(bladeEntry)
	if blade == nil {
		return nil
	}

	m.ds.ORM().Save(blade)
	return blade
}

func (m *BladeManager) DelBlade(id int64) {
	h, ok := m.mapBlade[id]
	if !ok {
		return
	}

	delete(m.mapBlade, id)
	m.ds.ORM().Delete(h)
}

func (m *BladeManager) BladeAddExp(id int64, exp int64) {
	blade, ok := m.mapBlade[id]

	if ok {
		blade.Exp += exp
		m.ds.ORM().Save(blade)
	}
}

func (m *BladeManager) BladeAddLevel(id int64, level int32) {
	blade, ok := m.mapBlade[id]

	if ok {
		blade.Level += level
		m.ds.ORM().Save(blade)
	}
}

func (m *BladeManager) PutonEquip(bladeID int64, equipID int64) error {
	/*blade, ok := m.mapBlade[bladeID]*/
	//if !ok {
	//return fmt.Errorf("invalid bladeid")
	/*}*/

	return nil
}

func (m *BladeManager) TakeoffEquip(bladeID int64) error {
	/*blade, ok := m.mapBlade[bladeID]*/
	//if !ok {
	//return fmt.Errorf("invalid blade_id")
	/*}*/

	return nil
}
