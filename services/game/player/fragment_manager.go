package player

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/utils"
	"github.com/valyala/bytebufferpool"
)

func MakeFragmentKey(fragmentId int32, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	b.B = append(b.B, "fragment_map.id_"...)
	b.B = append(b.B, strconv.Itoa(int(fragmentId))...)

	for _, f := range fields {
		b.B = append(b.B, "."...)
		b.B = append(b.B, f...)
	}

	return b.String()
}

type FragmentManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	owner       *Player         `bson:"-" json:"-"`
	FragmentMap map[int32]int32 `bson:"fragment_map" json:"fragment_map"` // 碎片包
}

func NewFragmentManager(owner *Player) *FragmentManager {
	m := &FragmentManager{
		owner:       owner,
		FragmentMap: make(map[int32]int32),
	}

	return m
}

func (m *FragmentManager) LoadAll() error {
	loadFragments := struct {
		FragmentMap map[string]int32 `bson:"fragment_map" json:"fragment_map"`
	}{
		FragmentMap: make(map[string]int32),
	}

	err := store.GetStore().LoadObject(define.StoreType_Fragment, m.owner.ID, &loadFragments)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("FragmentManager LoadAll: %w", err)
	}

	for k, v := range loadFragments.FragmentMap {
		ids := strings.Split(k, "id_")
		fragmentId, err := strconv.ParseInt(ids[len(ids)-1], 10, 32)
		if pass := utils.ErrCheck(err, "fragment id invalid", ids[len(ids)-1]); !pass {
			return err
		}

		m.FragmentMap[int32(fragmentId)] = v
	}

	return nil
}

// interface of cost_loot
func (m *FragmentManager) GetCostLootType() int32 {
	return define.CostLoot_Fragment
}

func (m *FragmentManager) CanCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.CanCost(typeMisc, num)
	if err != nil {
		return err
	}

	for k, v := range m.FragmentMap {
		if k != typeMisc {
			continue
		}

		if v >= num {
			return nil
		}
	}

	return fmt.Errorf("not enough fragment<%d>, num<%d>", typeMisc, num)
}

func (m *FragmentManager) DoCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.DoCost(typeMisc, num)
	if err != nil {
		return err
	}

	m.FragmentMap[typeMisc] -= num
	if m.FragmentMap[typeMisc] < 0 {
		m.FragmentMap[typeMisc] = 0
	}

	fields := map[string]interface{}{
		MakeFragmentKey(typeMisc): m.FragmentMap[typeMisc],
	}

	err = store.GetStore().SaveFields(define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "FragmentManager cost failed", typeMisc, num)
	return err
}

func (m *FragmentManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	m.FragmentMap[typeMisc] += num
	if m.FragmentMap[typeMisc] < 0 {
		m.FragmentMap[typeMisc] = 0
	}

	fields := map[string]interface{}{
		MakeFragmentKey(typeMisc): m.FragmentMap[typeMisc],
	}

	err = store.GetStore().SaveFields(define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "FragmentManager cost failed", typeMisc, num)
	return err
}

func (m *FragmentManager) GetFragmentList() []*pbGlobal.Fragment {
	reply := make([]*pbGlobal.Fragment, 0, len(m.FragmentMap))
	for k, v := range m.FragmentMap {
		reply = append(reply, &pbGlobal.Fragment{
			Id:  k,
			Num: v,
		})
	}

	return reply
}

func (m *FragmentManager) Inc(id, num int32) {
	m.FragmentMap[id] += num
	fields := map[string]interface{}{
		MakeFragmentKey(id): m.FragmentMap[id],
	}

	err := store.GetStore().SaveFields(define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "store SaveFields failed when FragmentManager Inc", m.owner.ID, fields)
}

func (m *FragmentManager) Compose(id int32) error {
	heroEntry, ok := auto.GetHeroEntry(id)
	if !ok {
		return fmt.Errorf("cannot find hero entry by id<%d>", id)
	}

	if heroEntry.FragmentCompose <= 0 {
		return fmt.Errorf("invalid hero entry<%d> fragmentCompose<%d>", id, heroEntry.FragmentCompose)
	}

	curNum := m.FragmentMap[id]
	if curNum < heroEntry.FragmentCompose {
		return fmt.Errorf("not enough fragment<%d> num<%d>", id, curNum)
	}

	_ = m.owner.HeroManager().AddHeroByTypeID(id)
	m.FragmentMap[id] -= heroEntry.FragmentCompose

	fields := map[string]interface{}{
		MakeFragmentKey(id): curNum - heroEntry.FragmentCompose,
	}

	err := store.GetStore().SaveFields(define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "store SaveFields failed when FragmentManager Compose", m.owner.ID, fields)
	return err
}

func (m *FragmentManager) SendFragmentsUpdate(ids ...int32) {
	reply := &pbGlobal.S2C_FragmentsUpdate{
		Frags: make([]*pbGlobal.Fragment, len(ids)),
	}

	for _, id := range ids {
		reply.Frags = append(reply.Frags, &pbGlobal.Fragment{
			Id:  id,
			Num: m.FragmentMap[id],
		})
	}

	m.owner.SendProtoMessage(reply)
}
