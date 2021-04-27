package player

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/spf13/cast"
	"github.com/valyala/bytebufferpool"
)

func makeFragmentKey(fragmentId int32, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("fragment_list.")
	_, _ = b.WriteString(cast.ToString(fragmentId))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

type FragmentManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	owner        *Player         `bson:"-" json:"-"`
	FragmentList map[int32]int32 `bson:"fragment_list" json:"fragment_list"` // 碎片包 map[英雄typeid]数量
}

func NewFragmentManager(owner *Player) *FragmentManager {
	m := &FragmentManager{
		owner:        owner,
		FragmentList: make(map[int32]int32),
	}

	return m
}

func (m *FragmentManager) LoadAll() error {
	loadFragments := struct {
		FragmentList map[int32]int32 `bson:"fragment_list" json:"fragment_list"`
	}{
		FragmentList: make(map[int32]int32),
	}

	err := store.GetStore().FindOne(context.Background(), define.StoreType_Fragment, m.owner.ID, &loadFragments)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("FragmentManager LoadAll: %w", err)
	}

	for fragmentId, num := range loadFragments.FragmentList {
		m.FragmentList[fragmentId] = num
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

	for k, v := range m.FragmentList {
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

	m.FragmentList[typeMisc] -= num
	if m.FragmentList[typeMisc] < 0 {
		m.FragmentList[typeMisc] = 0
	}

	fields := map[string]interface{}{
		makeFragmentKey(typeMisc): m.FragmentList[typeMisc],
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.DoCost", typeMisc, num)
	return err
}

func (m *FragmentManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	m.FragmentList[typeMisc] += num
	if m.FragmentList[typeMisc] < 0 {
		m.FragmentList[typeMisc] = 0
	}

	fields := map[string]interface{}{
		makeFragmentKey(typeMisc): m.FragmentList[typeMisc],
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.GainLoot", typeMisc, num)
	return err
}

func (m *FragmentManager) GetFragmentList() []*pbCommon.Fragment {
	reply := make([]*pbCommon.Fragment, 0, len(m.FragmentList))
	for k, v := range m.FragmentList {
		reply = append(reply, &pbCommon.Fragment{
			Id:  k,
			Num: v,
		})
	}

	return reply
}

func (m *FragmentManager) Inc(id, num int32) {
	m.FragmentList[id] += num
	fields := map[string]interface{}{
		makeFragmentKey(id): m.FragmentList[id],
	}

	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.Inc", m.owner.ID, fields)

	m.SendFragmentsUpdate(id)
}

func (m *FragmentManager) Compose(id int32) error {
	heroEntry, ok := auto.GetHeroEntry(id)
	if !ok {
		return fmt.Errorf("cannot find hero entry by id<%d>", id)
	}

	if heroEntry.FragmentCompose <= 0 {
		return fmt.Errorf("invalid hero entry<%d> fragmentCompose<%d>", id, heroEntry.FragmentCompose)
	}

	curNum := m.FragmentList[id]
	if curNum < heroEntry.FragmentCompose {
		return fmt.Errorf("not enough fragment<%d> num<%d>", id, curNum)
	}

	_ = m.owner.HeroManager().AddHeroByTypeId(id)
	m.FragmentList[id] -= heroEntry.FragmentCompose

	fields := map[string]interface{}{
		makeFragmentKey(id): curNum - heroEntry.FragmentCompose,
	}

	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.Compose", m.owner.ID, fields)
	return err
}

func (m *FragmentManager) GenFragmentListPB() []*pbCommon.Fragment {
	frags := make([]*pbCommon.Fragment, 0, len(m.FragmentList))
	for typeId, num := range m.FragmentList {
		frags = append(frags, &pbCommon.Fragment{
			Id:  typeId,
			Num: num,
		})
	}

	return frags
}

func (m *FragmentManager) SendFragmentsUpdate(ids ...int32) {
	reply := &pbCommon.S2C_FragmentsUpdate{
		Frags: make([]*pbCommon.Fragment, len(ids)),
	}

	for _, id := range ids {
		reply.Frags = append(reply.Frags, &pbCommon.Fragment{
			Id:  id,
			Num: m.FragmentList[id],
		})
	}

	m.owner.SendProtoMessage(reply)
}
