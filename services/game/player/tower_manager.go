package player

import (
	"errors"
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
	log "github.com/rs/zerolog/log"
)

var (
	ErrTowerMaxFloor          = errors.New("tower reach max floor")
	ErrTowerInvalidFloor      = errors.New("invalid tower floor")
	ErrTowerPhaseAwardPending = errors.New("tower phase award pending")
	ErrTowerTicketNotEnough   = errors.New("tower ticket not enough")
	ErrTowerLevelMissing      = errors.New("tower level config missing")
)

// TowerManager 爬塔管理器。
//
// 采用 kof 的单线推进模型：PassedFloor 记录已通过层数，
// 每次挑战下一层，消耗爬塔挑战券，按 Tower 表发放层奖励与阶段奖励。
type TowerManager struct {
	owner *Player `bson:"-" json:"-"`

	PassedFloor  int32 `bson:"passed_floor" json:"passed_floor"`   // 已通过层数
	PhaseAwarded bool  `bson:"phase_awarded" json:"phase_awarded"` // 当前阶段奖励是否已领
}

func NewTowerManager(owner *Player) *TowerManager {
	return &TowerManager{owner: owner}
}

func (m *TowerManager) start() {}

func (m *TowerManager) update() {}

func (m *TowerManager) OnHourChange(curHour int) {}

func towerConfig() *config.TowerConfigEntry {
	cfg := config.GetTowerConfigEntry()
	if cfg == nil {
		return &config.TowerConfigEntry{FloorPerPhase: 10}
	}

	return cfg
}

// isPhaseFloor 判断当前层是否为阶段奖励层
func (m *TowerManager) isPhaseFloor(floor int32) bool {
	per := towerConfig().FloorPerPhase
	if per <= 0 {
		return false
	}

	return floor > 0 && floor%per == 0
}

// Challenge 挑战下一层。
//
// towerType/towerFloor 仅用于校验客户端请求，实际推进以服务端 PassedFloor 为准。
func (m *TowerManager) Challenge(towerType int32, towerFloor int32, battleArray []int64) error {
	nextFloor := m.PassedFloor + 1

	if towerFloor > 0 && towerFloor != nextFloor {
		return fmt.Errorf("%w: request floor %d, expect %d", ErrTowerInvalidFloor, towerFloor, nextFloor)
	}

	entry, ok := config.GetTowerEntry(nextFloor)
	if !ok {
		return ErrTowerMaxFloor
	}

	// 上一个是阶段层且未发奖时不允许继续挑战（与 kof 的宝箱领取前置一致）
	if m.isPhaseFloor(m.PassedFloor) && m.PhaseAwarded == false && m.PassedFloor > 0 {
		return ErrTowerPhaseAwardPending
	}

	if entry.LevelId <= 0 {
		return ErrTowerLevelMissing
	}

	// 挑战券
	ticketId := int32(config.GoodsConfigId_TowerTicket)
	itemMgr := m.owner.ItemManager()
	if err := itemMgr.OnCheckEnough(ticketId, 1); err != nil {
		return ErrTowerTicketNotEnough
	}

	if err := itemMgr.OnItemCost(ticketId, 1, "towerChallenge"); err != nil {
		return err
	}

	// 服务器战斗（配置了战斗服时）
	if _, _, err := m.owner.BattleManager().CreateCommonBattle(
		pbGlobal.EMBattleType_BattleType_Tower,
		entry.LevelId,
		nil, nil, &pbGlobal.BattleExtraData{BattleIndex: nextFloor},
	); err != nil {
		log.Warn().Err(err).Int32("floor", nextFloor).Msg("tower battle request failed, fallback to client battle")
	}

	m.PassedFloor = nextFloor

	// 层奖励
	m.owner.GrantAwardList(entry.FloorAward, "towerFloor")

	// 阶段奖励
	m.PhaseAwarded = false
	if m.isPhaseFloor(m.PassedFloor) {
		m.owner.GrantAwardList(entry.PhaseAward, "towerPhase")
		m.PhaseAwarded = true
	}

	// 派发爬塔事件供任务判定（TowerFloorPass 类子条件）
	m.owner.eventManager.DispatchEvent(&event.Event{
		Type:  define.Event_Type_TowerPass,
		Miscs: []any{m.PassedFloor},
	})

	m.SendTowerUpdate(towerType)
	return nil
}

// TowerBattleResult 战斗结果回填（由战斗系统调用）
func (m *TowerManager) TowerBattleResult(win bool) {
	if !win {
		return
	}

	entry, ok := config.GetTowerEntry(m.PassedFloor)
	if !ok {
		return
	}

	m.owner.GrantAwardList(entry.FloorAward, "towerBattleWin")
}

func (m *TowerManager) GmFloorPass(towerType int32, floor int32) error {
	if floor < 0 {
		return ErrTowerInvalidFloor
	}

	m.PassedFloor = floor
	m.SendTowerUpdate(towerType)
	return nil
}

func (m *TowerManager) GmSettleReward(days int) {}

func (m *TowerManager) GenTowerInfoPB() []*pbGlobal.Tower {
	return []*pbGlobal.Tower{
		{
			Type:  0,
			Floor: m.PassedFloor,
		},
	}
}

func (m *TowerManager) SendTowerUpdate(tp int32) {
	for _, t := range m.GenTowerInfoPB() {
		m.owner.SendProtoMessage(&pbGlobal.S2C_TowerUpdate{Tower: t})
	}
}
