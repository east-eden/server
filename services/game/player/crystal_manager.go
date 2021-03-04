package player

import (
	"errors"
	"fmt"
	"strconv"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/att"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/crystal"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"bitbucket.org/funplus/server/utils/random"
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
	itemEntry, ok := auto.GetItemEntry(typeId)
	if !ok {
		return nil, fmt.Errorf("GetItemEntry<%d> failed", typeId)
	}

	crystalEntry, ok := auto.GetCrystalEntry(typeId)
	if !ok {
		return nil, fmt.Errorf("GetCrystalEntry<%d> failed", typeId)
	}

	c, err := m.createEntryCrystal(itemEntry, crystalEntry)
	if err != nil {
		return nil, err
	}

	m.CrystalMap[c.Id] = c

	fields := map[string]interface{}{
		MakeCrystalKey(c.Id): c,
	}
	err = store.GetStore().SaveFields(define.StoreType_Crystal, m.owner.ID, fields)

	return c, err
}

func (m *CrystalManager) delCrystal(id int64) error {
	c, ok := m.CrystalMap[id]
	if !ok {
		return fmt.Errorf("invalid crystal id<%d>", id)
	}

	c.EquipObj = -1
	delete(m.CrystalMap, id)

	fieldsName := []string{MakeCrystalKey(id)}
	err := store.GetStore().DeleteFields(define.StoreType_Crystal, m.owner.ID, fieldsName)
	crystal.GetCrystalPool().Put(c)
	return err
}

func (m *CrystalManager) initCrystalAtt(c *crystal.Crystal) {
	globalConfig, _ := auto.GetGlobalConfig()

	// 初始主属性
	mainAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, c.ItemEntry.Quality, define.Crystal_AttTypeMain)
	mainAttItem, err := random.PickOne(mainAttRepoList)
	if err != nil {
		log.Error().Err(err).Int64("crystal_id", c.Id).Msg("pick crystal main att failed")
		return
	}

	// 记录主属性库id
	mainAttRepoEntry := mainAttItem.(*auto.CrystalAttRepoEntry)
	c.MainAtt.AttRepoId = mainAttRepoEntry.Id
	c.MainAtt.AttRandRatio = random.Int32(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1])

	// 随机几条副属性
	viceAttNum := auto.GetCrystalInitViceAttNum(c.ItemEntry.Quality)

	// 初始副属性 todo new att manager
	viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, c.ItemEntry.Quality, define.Crystal_AttTypeVice)
	viceAttItems, err := random.PickUnrepeated(viceAttRepoList, viceAttNum)
	for k, v := range viceAttItems {
		viceAttRepoEntry := v.(*auto.CrystalAttRepoEntry)
		c.ViceAtts[k] = crystal.CrystalViceAtt{
			AttId:        viceAttRepoEntry.AttId,
			AttRandRatio: random.Int32(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1]),
		}

		am := att.NewAttManager()
		am.SetBaseAttId(viceAttRepoEntry.AttId)
	}
}

func (m *CrystalManager) createEntryCrystal(itemEntry *auto.ItemEntry, crystalEntry *auto.CrystalEntry) (*crystal.Crystal, error) {
	if itemEntry == nil {
		return nil, errors.New("invalid ItemEntry")
	}

	if crystalEntry == nil {
		return nil, errors.New("invalid CrystalEntry")
	}

	id, err := utils.NextID(define.SnowFlake_Crystal)
	if err != nil {
		return nil, err
	}

	c := crystal.NewCrystal(
		crystal.Id(id),
		crystal.OwnerId(m.owner.GetID()),
		crystal.TypeId(crystalEntry.Id),
		crystal.ItemEntry(itemEntry),
		crystal.CrystalEntry(crystalEntry),
	)

	// 生成初始属性
	m.initCrystalAtt(c)

	m.CrystalMap[c.GetOptions().Id] = c
	fields := map[string]interface{}{
		MakeCrystalKey(id): c,
	}
	err = store.GetStore().SaveFields(define.StoreType_Crystal, m.owner.ID, fields)

	return c, err
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
		if err := m.AddCrystalByTypeId(typeMisc); err != nil {
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
		c := crystal.NewCrystal()
		c.Options = v.Options
		c.MainAtt = v.MainAtt
		c.ViceAtts = v.ViceAtts

		itemEntry, ok := auto.GetItemEntry(c.TypeId)
		if !ok {
			log.Error().Int32("typeid", c.TypeId).Msg("cannot find item entry")
			continue
		}

		crystalEntry, ok := auto.GetCrystalEntry(c.TypeId)
		if !ok {
			log.Error().Int32("typeid", c.TypeId).Msg("cannot find crystal entry")
			continue
		}

		c.ItemEntry = itemEntry
		c.CrystalEntry = crystalEntry
		m.CrystalMap[c.Id] = c
	}

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

func (m *CrystalManager) AddCrystalByTypeId(typeId int32) error {
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

		if v.CrystalEntry.Id == typeId && v.GetEquipObj() == -1 {
			decNum--
			delId := v.Id
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
