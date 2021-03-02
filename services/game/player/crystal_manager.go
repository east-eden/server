package player

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/crystal"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/valyala/bytebufferpool"
)

func MakeCrystalKey(crystalId int64, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	b.B = append(b.B, "crystal_map.id_"...)
	b.B = append(b.B, strconv.Itoa(int(crystalId))...)

	for _, f := range fields {
		b.B = append(b.B, "."...)
		b.B = append(b.B, f...)
	}

	return b.String()
}

type CrystalManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	owner      *Player                    `bson:"-" json:"-"`
	CrystalMap map[int64]*crystal.Crystal `bson:"crystal_map" json:"crystal_map"`
}

func NewCrystalManager(owner *Player) *CrystalManager {
	m := &CrystalManager{
		owner:      owner,
		CrystalMap: make(map[int64]*crystal.Crystal),
	}

	return m
}

func (m *CrystalManager) Destroy() {
	for _, r := range m.CrystalMap {
		crystal.GetCrystalPool().Put(r)
	}
}

func (m *CrystalManager) createCrystal(typeId int32) (*crystal.Crystal, error) {
	crystalEntry, ok := auto.GetCrystalEntry(typeId)
	if !ok {
		return nil, fmt.Errorf("GetCrystalEntry<%d> failed", typeId)
	}

	r, err := m.createEntryCrystal(crystalEntry)
	if err != nil {
		return nil, err
	}

	m.CrystalMap[r.GetOptions().Id] = r

	fields := map[string]interface{}{
		MakeCrystalKey(r.Id): r,
	}
	err = store.GetStore().SaveFields(define.StoreType_Crystal, m.owner.ID, fields)

	return r, err
}

func (m *CrystalManager) delCrystal(id int64) error {
	r, ok := m.CrystalMap[id]
	if !ok {
		return fmt.Errorf("invalid crystal id<%d>", id)
	}

	r.GetOptions().EquipObj = -1
	delete(m.CrystalMap, id)

	fieldsName := []string{MakeCrystalKey(id)}
	err := store.GetStore().DeleteFields(define.StoreType_Crystal, m.owner.ID, fieldsName)
	crystal.GetCrystalPool().Put(r)
	return err
}

func (m *CrystalManager) createCrystalAtt(r *crystal.Crystal) {

	switch r.GetOptions().Entry.Pos {

	//1号位    主属性   攻击
	case define.Crystal_Pos1:
		attMain := &crystal.CrystalAtt{AttType: define.Att_Atk, AttValue: 100}
		r.SetAtt(0, attMain)

	//2号位    主属性   体%、攻%、防%、速度（随机）
	case define.Crystal_Pos2:
		tp := []int32{
			define.Att_Armor,
			define.Att_Atk,
			define.Att_Crit,
			define.Att_AtbSpeed,
		}
		attMain := &crystal.CrystalAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)

	//3号位    主属性   速度
	case define.Crystal_Pos3:
		attMain := &crystal.CrystalAtt{AttType: define.Att_AtbSpeed, AttValue: 100}
		r.SetAtt(0, attMain)

	//4号位    主属性   体%、攻%、防%（随机）
	case define.Crystal_Pos4:
		tp := []int32{
			define.Att_Armor,
			define.Att_Atk,
			define.Att_AtbSpeed,
		}
		attMain := &crystal.CrystalAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)

	//5号位    主属性   体力
	case define.Crystal_Pos5:
		attMain := &crystal.CrystalAtt{AttType: define.Att_Crit, AttValue: 100}
		r.SetAtt(0, attMain)

	//6号位    主属性   体%、攻%、防%、暴击%、暴伤%（随机）
	case define.Crystal_Pos6:
		tp := []int32{
			define.Att_AtbSpeed,
			define.Att_Atk,
			define.Att_Armor,
			define.Att_Crit,
			define.Att_CritInc,
		}
		attMain := &crystal.CrystalAtt{AttType: tp[rand.Intn(len(tp))], AttValue: 100}
		r.SetAtt(0, attMain)
	}
}

func (m *CrystalManager) createEntryCrystal(entry *auto.CrystalEntry) (*crystal.Crystal, error) {
	if entry == nil {
		return nil, errors.New("invalid CrystalEntry")
	}

	id, err := utils.NextID(define.SnowFlake_Crystal)
	if err != nil {
		return nil, err
	}

	r := crystal.NewCrystal(
		crystal.Id(id),
		crystal.OwnerId(m.owner.GetID()),
		crystal.TypeId(entry.Id),
		crystal.Entry(entry),
	)

	m.createCrystalAtt(r)
	m.CrystalMap[r.GetOptions().Id] = r

	fields := map[string]interface{}{
		MakeCrystalKey(id): r,
	}
	err = store.GetStore().SaveFields(define.StoreType_Crystal, m.owner.ID, fields)

	r.CalcAtt()

	return r, err
}

// interface of cost_loot
func (m *CrystalManager) GetCostLootType() int32 {
	return define.CostLoot_Crystal
}

func (m *CrystalManager) CanCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.CanCost(typeMisc, num)
	if err != nil {
		return err
	}

	var fixNum int32
	for _, v := range m.CrystalMap {
		if v.GetOptions().TypeId == typeMisc && v.GetEquipObj() == -1 {
			fixNum += 1
		}
	}

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough crystal<%d>, num<%d>", typeMisc, num)
}

func (m *CrystalManager) DoCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.DoCost(typeMisc, num)
	if err != nil {
		return err
	}

	return m.CostCrystalByTypeID(typeMisc, num)
}

func (m *CrystalManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	var n int32
	for n = 0; n < num; n++ {
		if err := m.AddCrystalByTypeID(typeMisc); err != nil {
			return err
		}
	}

	return nil
}

func (m *CrystalManager) LoadAll() error {
	loadCrystals := struct {
		CrystalMap map[string]*crystal.Crystal `bson:"crystal_map" json:"crystal_map"`
	}{
		CrystalMap: make(map[string]*crystal.Crystal),
	}

	err := store.GetStore().LoadObject(define.StoreType_Crystal, m.owner.ID, &loadCrystals)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("CrystalManager LoadAll: %w", err)
	}

	for _, v := range loadCrystals.CrystalMap {
		r := crystal.NewCrystal()
		r.Options = v.Options
		err := m.initLoadedCrystal(r)
		if err != nil {
			return fmt.Errorf("CrystalManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *CrystalManager) initLoadedCrystal(r *crystal.Crystal) error {
	entry, ok := auto.GetCrystalEntry(r.GetOptions().TypeId)
	if !ok {
		return fmt.Errorf("crystal<%d> entry invalid", r.GetOptions().TypeId)
	}

	if r.GetOptions().Entry == nil {
		return fmt.Errorf("crystal<%d> entry invalid", r.GetOptions().TypeId)
	}

	r.GetOptions().Entry = entry

	var n int32
	for n = 0; n < define.Crystal_AttNum; n++ {
		if oldAtt := r.GetAtt(int32(n)); oldAtt != nil {
			att := &crystal.CrystalAtt{AttType: oldAtt.AttType, AttValue: oldAtt.AttValue}
			r.SetAtt(n, att)
		}
	}

	m.CrystalMap[r.GetOptions().Id] = r
	r.CalcAtt()
	return nil
}

func (m *CrystalManager) SaveCrystalEquiped(id int64, equipObj int64) error {
	c := m.GetCrystal(id)
	if c == nil {
		return fmt.Errorf("invalid crystal id<%d>", id)
	}

	fields := map[string]interface{}{
		MakeCrystalKey(id, "equip_obj"): c.GetEquipObj(),
	}
	return store.GetStore().SaveFields(define.StoreType_Crystal, m.owner.ID, fields)
}

func (m *CrystalManager) GetCrystal(id int64) *crystal.Crystal {
	return m.CrystalMap[id]
}

func (m *CrystalManager) GetCrystalNums() int {
	return len(m.CrystalMap)
}

func (m *CrystalManager) GetCrystalList() []*crystal.Crystal {
	list := make([]*crystal.Crystal, 0)

	for _, v := range m.CrystalMap {
		list = append(list, v)
	}

	return list
}

func (m *CrystalManager) AddCrystalByTypeID(typeId int32) error {
	r, err := m.createCrystal(typeId)
	if err != nil {
		return err
	}

	m.SendCrystalAdd(r)
	return nil
}

func (m *CrystalManager) DeleteCrystal(id int64) error {
	if r := m.GetCrystal(id); r == nil {
		return fmt.Errorf("cannot find crystal<%d> while DeleteCrystal", id)
	}

	err := m.delCrystal(id)
	m.SendCrystalDelete(id)

	return err
}

func (m *CrystalManager) CostCrystalByTypeID(typeId int32, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec crystal error, invalid number:%d", num)
	}

	var err error
	decNum := num
	for _, v := range m.CrystalMap {
		if decNum <= 0 {
			break
		}

		if v.GetOptions().Entry.Id == typeId && v.GetEquipObj() == -1 {
			decNum--
			delId := v.GetOptions().Id
			if errDel := m.delCrystal(delId); errDel != nil {
				err = errDel
				utils.ErrPrint(errDel, "delCrystal failed when CostCrystalByTypeID", typeId, num, m.owner.ID)
				continue
			}

			m.SendCrystalDelete(delId)
		}
	}

	if decNum > 0 {
		log.Warn().
			Int32("need_dec", num).
			Int32("actual_dec", num-decNum).
			Msg("cost crystal not enough")
	}

	return err
}

func (m *CrystalManager) CostCrystalByID(id int64) error {
	r := m.GetCrystal(id)
	if r == nil {
		return fmt.Errorf("cannot find crystal by id:%d", id)
	}

	err := m.delCrystal(id)
	m.SendCrystalDelete(id)

	return err
}

func (m *CrystalManager) SetCrystalEquiped(id int64, objId int64) error {
	r, ok := m.CrystalMap[id]
	if !ok {
		return fmt.Errorf("invalid crystal id<%d>", id)
	}

	r.GetOptions().EquipObj = objId

	fields := map[string]interface{}{
		MakeCrystalKey(id, "equip_obj"): r.GetOptions().EquipObj,
	}
	err := store.GetStore().SaveFields(define.StoreType_Crystal, m.owner.ID, fields)

	m.SendCrystalUpdate(r)
	return err
}

func (m *CrystalManager) SetCrystalUnEquiped(id int64) error {
	r, ok := m.CrystalMap[id]
	if !ok {
		return fmt.Errorf("invalid crystal id<%d>", id)
	}

	r.GetOptions().EquipObj = -1

	fields := map[string]interface{}{
		MakeCrystalKey(id, "equip_obj"): r.GetOptions().EquipObj,
	}
	err := store.GetStore().SaveFields(define.StoreType_Crystal, m.owner.ID, fields)

	m.SendCrystalUpdate(r)
	return err
}

func (m *CrystalManager) SendCrystalAdd(r *crystal.Crystal) {
	msg := &pbGlobal.S2C_CrystalAdd{
		Crystal: &pbGlobal.Crystal{
			Id:     r.GetOptions().Id,
			TypeId: int32(r.GetOptions().TypeId),
		},
	}

	m.owner.SendProtoMessage(msg)
}

func (m *CrystalManager) SendCrystalDelete(id int64) {
	msg := &pbGlobal.S2C_DelCrystal{
		CrystalId: id,
	}

	m.owner.SendProtoMessage(msg)
}

func (m *CrystalManager) SendCrystalUpdate(r *crystal.Crystal) {
	msg := &pbGlobal.S2C_CrystalUpdate{
		Crystal: &pbGlobal.Crystal{
			Id:     r.GetOptions().Id,
			TypeId: int32(r.GetOptions().TypeId),
		},
	}

	m.owner.SendProtoMessage(msg)
}
