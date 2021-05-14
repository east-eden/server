package player

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/event"
	"bitbucket.org/funplus/server/services/game/global"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
)

var (
	ErrTowerInvalidType        = errors.New("invalid tower type")
	ErrTowerInvalidFloor       = errors.New("invalid tower floor")
	ErrTowerInvalidEntry       = errors.New("invalid tower entry")
	ErrTowerLevelLimit         = errors.New("level limit")
	ErrTowerGeneralFloorLimit  = errors.New("general tower floor limit")
	ErrTowerInvalidBattleArray = errors.New("invalid tower battle array")
)

type TowerManager struct {
	owner    *Player                      `bson:"-" json:"-"`
	CurFloor [define.Tower_Type_End]int32 `bson:"cur_floor" json:"cur_floor"`
}

func NewTowerManager(owner *Player) *TowerManager {
	m := &TowerManager{
		owner: owner,
	}

	return m
}

// 刷新记录处理
func (m *TowerManager) refreshRecord(towerType int32, floor int32, battleArray []int64) {
	seconds, err := global.GetGlobalMess().GetTowerBestSeconds(towerType, floor)
	if err != nil {
		return
	}

	// todo check record seconds
	if seconds != -1 {
		return
	}

	e := &event.Event{
		Type: define.Event_Type_TowerPass,
		Miscs: []interface{}{
			towerType,
			floor,
			&global.TowerBestInfo{
				PlayerId:    m.owner.ID,
				PlayerName:  m.owner.Name,
				Seconds:     10,
				RecordId:    1001,
				BattleArray: make([]int64, len(battleArray)),
			},
		},
	}

	copy(e.Miscs[2].(*global.TowerBestInfo).BattleArray[:], battleArray[:])
	global.GetGlobalMess().AddEvent(e)
}

func (m *TowerManager) Challenge(towerType int32, floor int32, battleArray []int64) error {
	if !utils.BetweenInt32(towerType, define.Tower_Type_Begin, define.Tower_Type_End) {
		return ErrTowerInvalidType
	}

	if floor != m.CurFloor[towerType]+1 {
		return ErrTowerInvalidFloor
	}

	// battle hero id repeated
	battleMap := make(map[int64]struct{}, 8)
	for _, heroId := range battleArray {
		battleMap[heroId] = struct{}{}
	}
	if len(battleMap) != len(battleArray) {
		return ErrTowerInvalidBattleArray
	}

	// limit
	towerEntry, ok := auto.GetTowerEntry(towerType, floor)
	if !ok {
		return ErrTowerInvalidEntry
	}

	if m.owner.Level < towerEntry.LevelLimit {
		return ErrTowerLevelLimit
	}

	// 种族塔需要综合试炼达到一定层数才开启
	if towerType != define.Tower_Type_General && floor == 1 {
		if m.CurFloor[define.Tower_Type_General] < 10 {
			return ErrTowerGeneralFloorLimit
		}
	}

	// 检查阵容种族是否符合
	switch towerType {
	// 综合塔没有种族限制
	case define.Tower_Type_General:
		break
	default:
		for _, heroId := range battleArray {
			if heroId == -1 {
				continue
			}

			h := m.owner.HeroManager().GetHero(heroId)
			if h == nil {
				continue
			}

			if h.Entry.Race != towerType {
				return ErrTowerInvalidBattleArray
			}
		}
	}

	// pass
	m.CurFloor[towerType]++

	// first reward
	err := m.owner.CostLootManager().GainLoot(towerEntry.FirstRewardId)
	utils.ErrPrint(err, "GainLoot failed when TowerManager.FloorPass", m.owner.ID, towerType, floor)

	// save
	fields := map[string]interface{}{
		fmt.Sprintf("cur_floor.%d", towerType): m.CurFloor,
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when TowerManager.FloorPass", m.owner.ID, fields)

	m.SendTowerUpdate(towerType)

	m.refreshRecord(towerType, floor, battleArray)

	return err
}

func (m *TowerManager) GmFloorPass(towerType int32, floor int32) error {
	if !utils.BetweenInt32(towerType, define.Tower_Type_Begin, define.Tower_Type_End) {
		return ErrTowerInvalidType
	}

	if floor <= 0 {
		return ErrTowerInvalidFloor
	}

	towerEntry, ok := auto.GetTowerEntry(towerType, floor)
	if !ok {
		return ErrTowerInvalidEntry
	}

	m.CurFloor[towerType] = floor

	// first reward
	err := m.owner.CostLootManager().GainLoot(towerEntry.FirstRewardId)
	utils.ErrPrint(err, "GainLoot failed when TowerManager.FloorPass", m.owner.ID, towerType, floor)

	// save
	fields := map[string]interface{}{
		fmt.Sprintf("cur_floor.%d", towerType): m.CurFloor[towerType],
	}

	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when TowerManager.GmFloorPass", m.owner.ID, fields)

	m.refreshRecord(towerType, floor, []int64{1, 2, 3})
	return err
}

func (m *TowerManager) GenTowerInfoPB() []*pbGlobal.Tower {
	pb := make([]*pbGlobal.Tower, 0, define.Tower_Type_End)
	for tp, floor := range m.CurFloor {
		pb = append(pb, &pbGlobal.Tower{
			Type:  int32(tp),
			Floor: floor,
		})
	}

	return pb
}

func (m *TowerManager) SendTowerUpdate(tp int32) {
	if !utils.BetweenInt32(tp, define.Tower_Type_Begin, define.Tower_Type_End) {
		return
	}

	msg := &pbGlobal.S2C_TowerUpdate{
		Tower: &pbGlobal.Tower{
			Type:  tp,
			Floor: m.CurFloor[tp],
		},
	}

	m.owner.SendProtoMessage(msg)
}
