package player

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/spf13/cast"
	"github.com/valyala/bytebufferpool"
)

//////////////////////////////////////////////////////////////
// BaseFragmentManager
type BaseFragmentManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	costLootType int32           `bson:"-" json:"-"`
	owner        *Player         `bson:"-" json:"-"`
	FragmentList map[int32]int32 `bson:"fragment_list" json:"fragment_list"` // 碎片包 map[英雄typeid]数量
}

func (m *BaseFragmentManager) makeFragmentKey(fragmentId int32, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	var prefix string
	switch m.costLootType {
	case define.CostLoot_HeroFragment:
		prefix = "hero_fragment."
	case define.CostLoot_CollectionFragment:
		prefix = "collection_fragment."
	}

	_, _ = b.WriteString(prefix)
	_, _ = b.WriteString("fragment_list.")
	_, _ = b.WriteString(cast.ToString(fragmentId))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

func (m *BaseFragmentManager) CanCost(typeMisc int32, num int32) error {
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

func (m *BaseFragmentManager) DoCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.DoCost(typeMisc, num)
	if err != nil {
		return err
	}

	m.FragmentList[typeMisc] -= num
	if m.FragmentList[typeMisc] < 0 {
		m.FragmentList[typeMisc] = 0
	}

	fields := map[string]interface{}{
		m.makeFragmentKey(typeMisc): m.FragmentList[typeMisc],
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.DoCost", typeMisc, num)
	return err
}

func (m *BaseFragmentManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	m.FragmentList[typeMisc] += num
	if m.FragmentList[typeMisc] < 0 {
		m.FragmentList[typeMisc] = 0
	}

	fields := map[string]interface{}{
		m.makeFragmentKey(typeMisc): m.FragmentList[typeMisc],
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.GainLoot", typeMisc, num)
	return err
}

func (m *BaseFragmentManager) GetFragmentList() []*pbGlobal.Fragment {
	reply := make([]*pbGlobal.Fragment, 0, len(m.FragmentList))
	for k, v := range m.FragmentList {
		reply = append(reply, &pbGlobal.Fragment{
			Id:  k,
			Num: v,
		})
	}

	return reply
}

func (m *BaseFragmentManager) GenFragmentListPB() []*pbGlobal.Fragment {
	frags := make([]*pbGlobal.Fragment, 0, len(m.FragmentList))
	for typeId, num := range m.FragmentList {
		frags = append(frags, &pbGlobal.Fragment{
			Id:  typeId,
			Num: num,
		})
	}

	return frags
}

//////////////////////////////////////////////////////////
// HeroFragmentManager
type HeroFragmentManager struct {
	*BaseFragmentManager `bson:"hero_fragment" json:"hero_fragment"`
}

func (m *HeroFragmentManager) GetCostLootType() int32 {
	return define.CostLoot_HeroFragment
}

func (m *HeroFragmentManager) Inc(id, num int32) {
	_ = m.GainLoot(id, num)
	m.SendFragmentsUpdate(id)
}

func (m *HeroFragmentManager) SendFragmentsUpdate(ids ...int32) {
	reply := &pbGlobal.S2C_HeroFragmentsUpdate{
		Frags: make([]*pbGlobal.Fragment, len(ids)),
	}

	for _, id := range ids {
		reply.Frags = append(reply.Frags, &pbGlobal.Fragment{
			Id:  id,
			Num: m.FragmentList[id],
		})
	}

	m.owner.SendProtoMessage(reply)
}

//////////////////////////////////////////////////////////
// CollectionFragmentManager
type CollectionFragmentManager struct {
	*BaseFragmentManager `bson:"collection_fragment" json:"collection_fragment"`
}

func (m *CollectionFragmentManager) GetCostLootType() int32 {
	return define.CostLoot_CollectionFragment
}

func (m *CollectionFragmentManager) Inc(id, num int32) {
	_ = m.GainLoot(id, num)
	m.SendFragmentsUpdate(id)
}

func (m *CollectionFragmentManager) SendFragmentsUpdate(ids ...int32) {
	reply := &pbGlobal.S2C_HeroFragmentsUpdate{
		Frags: make([]*pbGlobal.Fragment, len(ids)),
	}

	for _, id := range ids {
		reply.Frags = append(reply.Frags, &pbGlobal.Fragment{
			Id:  id,
			Num: m.FragmentList[id],
		})
	}

	m.owner.SendProtoMessage(reply)
}

// 所有碎片管理
type FragmentManager struct {
	owner                      *Player `bson:"-" json:"-"`
	*HeroFragmentManager       `bson:"inline" json:",inline"`
	*CollectionFragmentManager `bson:"inline" json:",inline"`
}

func NewFragmentManager(owner *Player) *FragmentManager {
	m := &FragmentManager{
		owner: owner,
		HeroFragmentManager: &HeroFragmentManager{
			BaseFragmentManager: &BaseFragmentManager{
				costLootType: define.CostLoot_HeroFragment,
				owner:        owner,
				FragmentList: make(map[int32]int32),
			},
		},

		CollectionFragmentManager: &CollectionFragmentManager{
			BaseFragmentManager: &BaseFragmentManager{
				costLootType: define.CostLoot_CollectionFragment,
				owner:        owner,
				FragmentList: make(map[int32]int32),
			},
		},
	}

	return m
}

func (m *FragmentManager) LoadAll() error {
	loadFragments := struct {
		HeroFragmentList       map[int32]int32 `bson:"hero_fragment.fragment_list" json:"hero_fragment.fragment_list"`
		CollectionFragmentList map[int32]int32 `bson:"collection_fragment.fragment_list" json:"collection_fragment.fragment_list"`
	}{
		HeroFragmentList:       make(map[int32]int32),
		CollectionFragmentList: make(map[int32]int32),
	}

	err := store.GetStore().FindOne(context.Background(), define.StoreType_Fragment, m.owner.ID, &loadFragments)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("FragmentManager LoadAll: %w", err)
	}

	for fragmentId, num := range loadFragments.HeroFragmentList {
		m.HeroFragmentManager.FragmentList[fragmentId] = num
	}

	for fragmentId, num := range loadFragments.CollectionFragmentList {
		m.CollectionFragmentManager.FragmentList[fragmentId] = num
	}

	return nil
}

func (m *FragmentManager) HeroCompose(id int32) error {
	heroEntry, ok := auto.GetHeroEntry(id)
	if !ok {
		return fmt.Errorf("cannot find hero entry by id<%d>", id)
	}

	if heroEntry.FragmentCompose <= 0 {
		return fmt.Errorf("invalid hero entry<%d> fragmentCompose<%d>", id, heroEntry.FragmentCompose)
	}

	err := m.HeroFragmentManager.CanCost(id, heroEntry.FragmentCompose)
	if err != nil {
		return err
	}

	err = m.HeroFragmentManager.DoCost(id, heroEntry.FragmentCompose)
	utils.ErrPrint(err, "HeroFragmentManager.DoCost failed when HeroCompose", m.owner.ID, id)

	fields := map[string]interface{}{
		m.HeroFragmentManager.makeFragmentKey(id): m.HeroFragmentManager.FragmentList[id],
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.HeroCompose", m.owner.ID, fields)
	return err
}

func (m *FragmentManager) CollectionCompose(id int32) error {
	collectionEntry, ok := auto.GetCollectionEntry(id)
	if !ok {
		return fmt.Errorf("cannot find collection entry by id<%d>", id)
	}

	if collectionEntry.FragmentCompose <= 0 {
		return fmt.Errorf("invalid collection entry<%d> fragmentCompose<%d>", id, collectionEntry.FragmentCompose)
	}

	err := m.CollectionFragmentManager.CanCost(id, collectionEntry.FragmentCompose)
	if err != nil {
		return err
	}

	err = m.CollectionFragmentManager.DoCost(id, collectionEntry.FragmentCompose)
	utils.ErrPrint(err, "CollectionFragmentManager.DoCost failed when CollectionCompose", m.owner.ID, id)

	fields := map[string]interface{}{
		m.CollectionFragmentManager.makeFragmentKey(id): m.CollectionFragmentManager.FragmentList[id],
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Fragment, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when FragmentManager.CollectionCompose", m.owner.ID, fields)
	return err
}
