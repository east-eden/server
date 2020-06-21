package player

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/rune"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/store/db"
	"github.com/yokaiio/yokai_server/utils"
)

type RuneManager struct {
	owner   *Player
	mapRune map[int64]rune.Rune

	sync.RWMutex
}

func NewRuneManager(owner *Player) *RuneManager {
	m := &RuneManager{
		owner:   owner,
		mapRune: make(map[int64]rune.Rune, 0),
	}

	return m
}

func (m *RuneManager) TableName() string {
	return "rune"
}

func (m *RuneManager) createRune(typeID int32) rune.Rune {
	runeEntry := entries.GetRuneEntry(typeID)
	r := m.createEntryRune(runeEntry)
	if r == nil {
		logger.Warning("new rune failed when createRune:", typeID)
		return nil
	}

	m.mapRune[r.GetOptions().Id] = r
	store.GetStore().SaveObject(store.StoreType_Rune, r)

	return r
}

func (m *RuneManager) delRune(id int64) {
	r, ok := m.mapRune[id]
	if !ok {
		return
	}

	r.GetOptions().EquipObj = -1
	delete(m.mapRune, id)
	store.GetStore().DeleteObject(store.StoreType_Rune, r)
	rune.ReleasePoolRune(r)
}

func (m *RuneManager) createRuneAtt(r rune.Rune) {

	switch r.GetOptions().Entry.Pos {

	//1号位    主属性   攻击
	case define.Rune_Position1:
		attMain := &rune.RuneAtt{AttType: define.Att_Atk, AttValue: 100}
		r.SetAtt(0, attMain)

	//2号位    主属性   体%、攻%、防%、速度（随机）
	case define.Rune_Position2:
		tp := []int32{
			define.Att_ConPercent,
			define.Att_AtkPercent,
			define.Att_DefPercent,
			define.Att_AtkSpeed,
		}
		attMain := &rune.RuneAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)

	//3号位    主属性   速度
	case define.Rune_Position3:
		attMain := &rune.RuneAtt{AttType: define.Att_AtkSpeed, AttValue: 100}
		r.SetAtt(0, attMain)

	//4号位    主属性   体%、攻%、防%（随机）
	case define.Rune_Position4:
		tp := []int32{
			define.Att_ConPercent,
			define.Att_AtkPercent,
			define.Att_DefPercent,
		}
		attMain := &rune.RuneAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)

	//5号位    主属性   体力
	case define.Rune_Position5:
		attMain := &rune.RuneAtt{AttType: define.Att_Con, AttValue: 100}
		r.SetAtt(0, attMain)

	//6号位    主属性   体%、攻%、防%、暴击%、暴伤%（随机）
	case define.Rune_Position6:
		tp := []int32{
			define.Att_ConPercent,
			define.Att_AtkPercent,
			define.Att_DefPercent,
			define.Att_CriProb,
			define.Att_CriDmg,
		}
		attMain := &rune.RuneAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)
	}
}

func (m *RuneManager) createEntryRune(entry *define.RuneEntry) rune.Rune {
	if entry == nil {
		logger.Error("createEntryRune with nil RuneEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Rune)
	if err != nil {
		logger.Error(err)
		return nil
	}

	r := rune.NewRune(
		rune.Id(id),
		rune.OwnerId(m.owner.GetID()),
		rune.TypeId(entry.ID),
		rune.Entry(entry),
	)

	m.createRuneAtt(r)
	m.mapRune[r.GetOptions().Id] = r
	store.GetStore().SaveObject(store.StoreType_Rune, r)

	r.CalcAtt()

	return r
}

// interface of cost_loot
func (m *RuneManager) GetCostLootType() int32 {
	return define.CostLoot_Rune
}

func (m *RuneManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager check item<%d> cost failed, wrong number<%d>", typeMisc, num)
	}

	var fixNum int32 = 0
	for _, v := range m.mapRune {
		if v.GetOptions().TypeId == typeMisc && v.GetEquipObj() == -1 {
			fixNum += 1
		}
	}

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough rune<%d>, num<%d>", typeMisc, num)
}

func (m *RuneManager) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager cost item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return m.CostRuneByTypeID(typeMisc, num)
}

func (m *RuneManager) CanGain(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager check gain item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	// todo bag max item

	return nil
}

func (m *RuneManager) GainLoot(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager gain rune<%d> failed, wrong number<%d>", typeMisc, num)
	}

	var n int32
	for n = 0; n < num; n++ {
		if err := m.AddRuneByTypeID(typeMisc); err != nil {
			return err
		}
	}

	return nil
}

func (m *RuneManager) LoadAll() error {
	runeList, err := store.GetStore().LoadArray(store.StoreType_Rune, "owner_id", m.owner.GetID(), rune.GetRunePool())
	if errors.Is(err, db.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("RuneManager LoadAll: %w", err)
	}

	for _, r := range runeList {
		err := m.initLoadedRune(r.(rune.Rune))
		if err != nil {
			return fmt.Errorf("RuneManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *RuneManager) initLoadedRune(r rune.Rune) error {
	entry := entries.GetRuneEntry(r.GetOptions().TypeId)

	if r.GetOptions().Entry == nil {
		return fmt.Errorf("rune<%d> entry invalid", r.GetOptions().TypeId)
	}

	r.GetOptions().Entry = entry

	var n int32
	for n = 0; n < define.Rune_AttNum; n++ {
		if oldAtt := r.GetAtt(int32(n)); oldAtt != nil {
			att := &rune.RuneAtt{AttType: oldAtt.AttType, AttValue: oldAtt.AttValue}
			r.SetAtt(n, att)
		}
	}

	m.mapRune[r.GetOptions().Id] = r
	store.GetStore().SaveObject(store.StoreType_Rune, r)

	r.CalcAtt()
	return nil
}

func (m *RuneManager) Save(id int64) {
	if r := m.GetRune(id); r != nil {
		store.GetStore().SaveObject(store.StoreType_Rune, r)
	}
}

func (m *RuneManager) GetRune(id int64) rune.Rune {
	return m.mapRune[id]
}

func (m *RuneManager) GetRuneNums() int {
	return len(m.mapRune)
}

func (m *RuneManager) GetRuneList() []rune.Rune {
	list := make([]rune.Rune, 0)

	for _, v := range m.mapRune {
		list = append(list, v)
	}

	return list
}

func (m *RuneManager) AddRuneByTypeID(typeID int32) error {
	r := m.createRune(typeID)
	if r == nil {
		return fmt.Errorf("AddRuneByTypeID failed: type_id = %d", typeID)
	}

	m.SendRuneAdd(r)
	return nil
}

func (m *RuneManager) DeleteRune(id int64) error {
	if r := m.GetRune(id); r == nil {
		return fmt.Errorf("cannot find rune<%d> while DeleteRune", id)
	}

	m.delRune(id)
	m.SendRuneDelete(id)

	return nil
}

func (m *RuneManager) CostRuneByTypeID(typeID int32, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec rune error, invalid number:%d", num)
	}

	decNum := num
	for _, v := range m.mapRune {
		if decNum <= 0 {
			break
		}

		if v.GetOptions().Entry.ID == typeID && v.GetEquipObj() == -1 {
			decNum--
			delId := v.GetOptions().Id
			m.delRune(delId)
			m.SendRuneDelete(delId)
		}
	}

	if decNum > 0 {
		logger.WithFields(logger.Fields{
			"need_dec":   num,
			"actual_dec": num - decNum,
		}).Warning("CostRuneByTypeID warning")
	}

	return nil
}

func (m *RuneManager) CostRuneByID(id int64) error {
	r := m.GetRune(id)
	if r == nil {
		return fmt.Errorf("cannot find rune by id:%d", id)
	}

	m.delRune(id)
	m.SendRuneDelete(id)

	return nil
}

func (m *RuneManager) SetRuneEquiped(id int64, objId int64) {
	r, ok := m.mapRune[id]
	if !ok {
		return
	}

	r.GetOptions().EquipObj = objId
	store.GetStore().SaveObject(store.StoreType_Rune, r)
	m.SendRuneUpdate(r)
}

func (m *RuneManager) SetRuneUnEquiped(id int64) {
	r, ok := m.mapRune[id]
	if !ok {
		return
	}

	r.GetOptions().EquipObj = -1
	store.GetStore().SaveObject(store.StoreType_Rune, r)
	m.SendRuneUpdate(r)
}

func (m *RuneManager) SendRuneAdd(r rune.Rune) {
	msg := &pbGame.M2C_RuneAdd{
		Rune: &pbGame.Rune{
			Id:     r.GetOptions().Id,
			TypeId: r.GetOptions().TypeId,
		},
	}

	m.owner.SendProtoMessage(msg)
}

func (m *RuneManager) SendRuneDelete(id int64) {
	msg := &pbGame.M2C_DelRune{
		RuneId: id,
	}

	m.owner.SendProtoMessage(msg)
}

func (m *RuneManager) SendRuneUpdate(r rune.Rune) {
	msg := &pbGame.M2C_RuneUpdate{
		Rune: &pbGame.Rune{
			Id:     r.GetOptions().Id,
			TypeId: r.GetOptions().TypeId,
		},
	}

	m.owner.SendProtoMessage(msg)
}
