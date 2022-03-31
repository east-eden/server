package global

import (
	"context"
	"errors"
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
)

// 塔最佳记录
type TowerBestInfo struct {
	PlayerId    int64   `bson:"player_id" json:"player_id"`
	PlayerName  string  `bson:"player_name" json:"player_name"`
	Seconds     int32   `bson:"seconds" json:"seconds"`           // 通关时间
	RecordId    int64   `bson:"record_id" json:"record_id"`       // 录像id
	BattleArray []int64 `bson:"battle_array" json:"battle_array"` // 阵容
}

// 获取塔通关时间
func (g *GlobalController) GetTowerBestSeconds(towerType int32, floor int32) (int32, error) {
	if !utils.Between(towerType, define.Tower_Type_Begin, define.Tower_Type_End) {
		return 0, errors.New("invalid tower type")
	}

	if !utils.Between(floor, 0, define.TowerMaxFloor) {
		return 0, errors.New("invalid tower floor")
	}

	g.RLock()
	defer g.RUnlock()

	record := g.TowerBestRecord[towerType][floor]
	if record == nil {
		return -1, nil
	}

	return record.Seconds, nil
}

// 塔通关
func (g *GlobalController) onEventTowerPass(e *event.Event) error {
	if e.Type != define.Event_Type_TowerPass {
		return errors.New("invalid event type")
	}

	towerType := e.Miscs[0].(int32)
	floor := e.Miscs[1].(int32)
	towerInfo := e.Miscs[2].(*TowerBestInfo)

	// 记录没超越
	if record := g.TowerBestRecord[towerType][floor]; record != nil && record.Seconds <= towerInfo.Seconds {
		return nil
	}

	g.TowerBestRecord[towerType][floor] = towerInfo

	fields := map[string]any{
		fmt.Sprintf("tower_best_record.%d.%d", towerType, floor): g.TowerBestRecord[towerType][floor],
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_GlobalMess, g.GameId, fields)
	utils.ErrPrint(err, "UpdateFields failed when GlobalMess.onEventTowerPass", g.GameId, fields)
	return nil
}
