package player

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/rune"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/valyala/bytebufferpool"
)

func MakeRuneKey(runeId int64, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	b.B = append(b.B, "rune_map.id_"...)
	b.B = append(b.B, strconv.Itoa(int(runeId))...)

	for _, f := range fields {
		b.B = append(b.B, "."...)
		b.B = append(b.B, f...)
	}

	return b.String()
}

type RuneManager struct {
	owner   *Player              `bson:"-" json:"-"`
	RuneMap map[int64]*rune.Rune `bson:"rune_map" json:"rune_map"`
}

func NewRuneManager(owner *Player) *RuneManager {
	m := &RuneManager{
		owner:   owner,
		RuneMap: make(map[int64]*rune.Rune),
	}

	return m
}

func (m *RuneManager) createRune(typeId int32) (*rune.Rune, error) {
	runeEntry, ok := auto.GetRuneEntry(typeId)
	if !ok {
		return nil, fmt.Errorf("GetRuneEntry<%d> failed", typeId)
	}

	r, err := m.createEntryRune(runeEntry)
	if err != nil {
		return nil, err
	}

	m.RuneMap[r.GetOptions().Id] = r

	fields := map[string]interface{}{
		MakeRuneKey(r.Id): r,
	}
	err = store.GetStore().SaveFields(define.StoreType_Rune, m.owner.ID, fields)

	return r, err
}

func (m *RuneManager) delRune(id int64) error {
	r, ok := m.RuneMap[id]
	if !ok {
		return fmt.Errorf("invalid rune id<%d>", id)
	}

	r.GetOptions().EquipObj = -1
	delete(m.RuneMap, id)

	fieldsName := []string{MakeRuneKey(id)}
	err := store.GetStore().DeleteFields(define.StoreType_Rune, m.owner.ID, fieldsName)
	rune.GetRunePool().Put(r)
	return err
}

func (m *RuneManager) createRuneAtt(r *rune.Rune) {

	switch r.GetOptions().Entry.Pos {

	//1号位    主属性   攻击
	case define.Rune_Position1:
		attMain := &rune.RuneAtt{AttType: define.Att_Atk, AttValue: 100}
		r.SetAtt(0, attMain)

	//2号位    主属性   体%、攻%、防%、速度（随机）
	case define.Rune_Position2:
		tp := []int32{
			define.Att_Armor,
			define.Att_Atk,
			define.Att_Crit,
			define.Att_AtbSpeed,
		}
		attMain := &rune.RuneAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)

	//3号位    主属性   速度
	case define.Rune_Position3:
		attMain := &rune.RuneAtt{AttType: define.Att_AtbSpeed, AttValue: 100}
		r.SetAtt(0, attMain)

	//4号位    主属性   体%、攻%、防%（随机）
	case define.Rune_Position4:
		tp := []int32{
			define.Att_Armor,
			define.Att_Atk,
			define.Att_AtbSpeed,
		}
		attMain := &rune.RuneAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)

	//5号位    主属性   体力
	case define.Rune_Position5:
		attMain := &rune.RuneAtt{AttType: define.Att_Crit, AttValue: 100}
		r.SetAtt(0, attMain)

	//6号位    主属性   体%、攻%、防%、暴击%、暴伤%（随机）
	case define.Rune_Position6:
		tp := []int32{
			define.Att_AtbSpeed,
			define.Att_Atk,
			define.Att_Armor,
			define.Att_Crit,
			define.Att_CritInc,
		}
		attMain := &rune.RuneAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)
	}
}

func (m *RuneManager) createEntryRune(entry *auto.RuneEntry) (*rune.Rune, error) {
	if entry == nil {
		return nil, errors.New("invalid RuneEntry")
	}

	id, err := utils.NextID(define.SnowFlake_Rune)
	if err != nil {
		return nil, err
	}

	r := rune.NewRune(
		rune.Id(id),
		rune.OwnerId(m.owner.GetID()),
		rune.TypeId(entry.Id),
		rune.Entry(entry),
	)

	m.createRuneAtt(r)
	m.RuneMap[r.GetOptions().Id] = r

	fields := map[string]interface{}{
		MakeRuneKey(id): r,
	}
	err = store.GetStore().SaveFields(define.StoreType_Rune, m.owner.ID, fields)

	r.CalcAtt()

	return r, err
}

// interface of cost_loot
func (m *RuneManager) GetCostLootType() int32 {
	return define.CostLoot_Rune
}

func (m *RuneManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager check item<%d> cost failed, wrong number<%d>", typeMisc, num)
	}

	var fixNum int32
	for _, v := range m.RuneMap {
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
	loadRunes := struct {
		RuneMap map[string]*rune.Rune `bson:"rune_map" json:"rune_map"`
	}{
		RuneMap: make(map[string]*rune.Rune),
	}

	err := store.GetStore().LoadObject(define.StoreType_Rune, m.owner.ID, &loadRunes)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("RuneManager LoadAll: %w", err)
	}

	for _, v := range loadRunes.RuneMap {
		r := rune.NewRune()
		r.Options = v.Options
		err := m.initLoadedRune(r)
		if err != nil {
			return fmt.Errorf("RuneManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *RuneManager) initLoadedRune(r *rune.Rune) error {
	entry, ok := auto.GetRuneEntry(r.GetOptions().TypeId)
	if !ok {
		return fmt.Errorf("rune<%d> entry invalid", r.GetOptions().TypeId)
	}

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

	m.RuneMap[r.GetOptions().Id] = r
	r.CalcAtt()
	return nil
}

func (m *RuneManager) Save(id int64) error {
	r := m.GetRune(id)
	if r == nil {
		return fmt.Errorf("invalid rune id<%d>", id)
	}

	fields := map[string]interface{}{
		MakeRuneKey(id): r,
	}
	return store.GetStore().SaveFields(define.StoreType_Rune, m.owner.ID, fields)
}

func (m *RuneManager) GetRune(id int64) *rune.Rune {
	return m.RuneMap[id]
}

func (m *RuneManager) GetRuneNums() int {
	return len(m.RuneMap)
}

func (m *RuneManager) GetRuneList() []*rune.Rune {
	list := make([]*rune.Rune, 0)

	for _, v := range m.RuneMap {
		list = append(list, v)
	}

	return list
}

func (m *RuneManager) AddRuneByTypeID(typeId int32) error {
	r, err := m.createRune(typeId)
	if err != nil {
		return err
	}

	m.SendRuneAdd(r)
	return nil
}

func (m *RuneManager) DeleteRune(id int64) error {
	if r := m.GetRune(id); r == nil {
		return fmt.Errorf("cannot find rune<%d> while DeleteRune", id)
	}

	err := m.delRune(id)
	m.SendRuneDelete(id)

	return err
}

func (m *RuneManager) CostRuneByTypeID(typeId int32, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec rune error, invalid number:%d", num)
	}

	var err error
	decNum := num
	for _, v := range m.RuneMap {
		if decNum <= 0 {
			break
		}

		if v.GetOptions().Entry.Id == typeId && v.GetEquipObj() == -1 {
			decNum--
			delId := v.GetOptions().Id
			if errDel := m.delRune(delId); errDel != nil {
				err = errDel
				utils.ErrPrint(errDel, "delRune failed when CostRuneByTypeID", typeId, num, m.owner.ID)
				continue
			}

			m.SendRuneDelete(delId)
		}
	}

	if decNum > 0 {
		log.Warn().
			Int32("need_dec", num).
			Int32("actual_dec", num-decNum).
			Msg("cost rune not enough")
	}

	return err
}

func (m *RuneManager) CostRuneByID(id int64) error {
	r := m.GetRune(id)
	if r == nil {
		return fmt.Errorf("cannot find rune by id:%d", id)
	}

	err := m.delRune(id)
	m.SendRuneDelete(id)

	return err
}

func (m *RuneManager) SetRuneEquiped(id int64, objId int64) error {
	r, ok := m.RuneMap[id]
	if !ok {
		return fmt.Errorf("invalid rune id<%d>", id)
	}

	r.GetOptions().EquipObj = objId

	fields := map[string]interface{}{
		MakeRuneKey(id, "equip_obj"): r.GetOptions().EquipObj,
	}
	err := store.GetStore().SaveFields(define.StoreType_Rune, m.owner.ID, fields)

	m.SendRuneUpdate(r)
	return err
}

func (m *RuneManager) SetRuneUnEquiped(id int64) error {
	r, ok := m.RuneMap[id]
	if !ok {
		return fmt.Errorf("invalid rune id<%d>", id)
	}

	r.GetOptions().EquipObj = -1

	fields := map[string]interface{}{
		MakeRuneKey(id, "equip_obj"): r.GetOptions().EquipObj,
	}
	err := store.GetStore().SaveFields(define.StoreType_Rune, m.owner.ID, fields)

	m.SendRuneUpdate(r)
	return err
}

func (m *RuneManager) SendRuneAdd(r *rune.Rune) {
	msg := &pbGlobal.S2C_RuneAdd{
		Rune: &pbGlobal.Rune{
			Id:     r.GetOptions().Id,
			TypeId: int32(r.GetOptions().TypeId),
		},
	}

	m.owner.SendProtoMessage(msg)
}

func (m *RuneManager) SendRuneDelete(id int64) {
	msg := &pbGlobal.S2C_DelRune{
		RuneId: id,
	}

	m.owner.SendProtoMessage(msg)
}

func (m *RuneManager) SendRuneUpdate(r *rune.Rune) {
	msg := &pbGlobal.S2C_RuneUpdate{
		Rune: &pbGlobal.Rune{
			Id:     r.GetOptions().Id,
			TypeId: int32(r.GetOptions().TypeId),
		},
	}

	m.owner.SendProtoMessage(msg)
}
