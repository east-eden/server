// Package activity 实现配置驱动的活动系统。
//
// kof 的活动管理器把 68 种活动类型（战斗通行证、七日任务、签到、礼包……）
// 写在一个 6000 多行的文件里，并与三十多个子系统耦合。这里抽出其中可复用的部分：
//
//   - 活动生命周期：按 Activity 表的 startMode/endMode 计算开启与结束时间
//   - 任务型活动：ActivityTask 表（子条件 + 奖励）驱动进度与领奖
//   - 签到型活动：ActivitySignIn 表（天数 + 奖励）驱动每日签到
//
// 具体活动类型的专属逻辑（通行证轮次、商店、抽奖等）不在这里，需要时按同样的
// 状态模型接入。
package activity

import (
	"errors"
	"fmt"
	"time"

	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/internal/activity"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/services/game/subcondition"
	"github.com/east-eden/server/utils"
	"google.golang.org/protobuf/proto"
)

var (
	ErrActivityNotFound   = errors.New("activity not found")
	ErrActivityClosed     = errors.New("activity not open")
	ErrTaskNotFound       = errors.New("activity task not found")
	ErrTaskNotCompleted   = errors.New("activity task not completed")
	ErrTaskAlreadyClaimed = errors.New("activity task already claimed")
	ErrSignInvalidDay     = errors.New("invalid sign in day")
	ErrSignAlreadyDone    = errors.New("sign in already done today")
	ErrBattlePassLocked   = errors.New("battle pass premium not unlocked")
	ErrBattlePassNoReward = errors.New("battle pass has no claimable reward")
)

// TaskState 活动任务的运行时状态
type TaskState struct {
	Id            int32                        `bson:"id" json:"id"`
	SubConditions []*subcondition.SubCondition `bson:"sub_conditions" json:"sub_conditions"`
	Claimed       bool                         `bson:"claimed" json:"claimed"`
}

// ActivityState 单个活动的运行时状态
type ActivityState struct {
	Id         string               `bson:"id" json:"id"`
	StartAt    int64                `bson:"start_at" json:"start_at"`
	EndAt      int64                `bson:"end_at" json:"end_at"`
	Open       bool                 `bson:"open" json:"open"`
	SignDays   int32                `bson:"sign_days" json:"sign_days"`     // 已签到天数
	Tasks      map[int32]*TaskState `bson:"tasks" json:"tasks"`             // 任务id -> 状态
	BattlePass *BattlePassState     `bson:"battle_pass" json:"battle_pass"` // 通行证进度
}

// BattlePassState 战斗通行证进度
type BattlePassState struct {
	Score     int32 `bson:"score" json:"score"`           // 当前积分
	Unlocked  bool  `bson:"unlocked" json:"unlocked"`     // 是否解锁高级奖励
	FreeIndex int32 `bson:"free_index" json:"free_index"` // 已领取的免费阶段数
	HighIndex int32 `bson:"high_index" json:"high_index"` // 已领取的高级阶段数
}

// IsOpen 判断活动在指定时间是否开启
func (s *ActivityState) IsOpen(now time.Time) bool {
	ts := now.UnixMilli()
	return s.Open && (s.StartAt == 0 || ts >= s.StartAt) && (s.EndAt == 0 || ts <= s.EndAt)
}

// Owner 活动管理器需要宿主提供的能力，由 Player 实现
type Owner interface {
	// GetLevel 玩家等级，用于按等级开启的活动
	GetLevel() int32
	// GrantAwardList 发放固定奖励
	GrantAwardList(awards []*config.IdCount, reason string)
	// SaveActivities 持久化活动数据
	SaveActivities(data map[string]*ActivityState) error
	// SendProtoMessage 下发协议
	SendProtoMessage(msg proto.Message)
}

// Manager 玩家活动管理器
type Manager struct {
	owner  Owner
	Data   map[string]*ActivityState
	lastAt int64
}

func NewManager(owner Owner) *Manager {
	return &Manager{
		owner: owner,
		Data:  make(map[string]*ActivityState, 16),
	}
}

// LoadAll 读档后接上已有数据
func (m *Manager) LoadAll(data map[string]*ActivityState) {
	if data != nil {
		m.Data = data
	}

	if m.Data == nil {
		m.Data = make(map[string]*ActivityState, 16)
	}
}

// CheckAll 按配置刷新所有活动状态：开启/关闭与结束时间。
//
// createdAt 为角色创建时间，供 createtime 开启模式使用。
func (m *Manager) CheckAll(now time.Time, createdAt time.Time) {
	playerLevel := m.owner.GetLevel()

	for id, entry := range config.GetActivityRows() {
		conf := &activity.ActivityConf{
			StartMode:              entry.StartMode,
			Start:                  entry.Start,
			EndMode:                entry.EndMode,
			End:                    entry.End,
			CreatetimeSpecialParam: entry.CreatetimeSpecialParam,
		}

		start, end := activity.GetActivityTime(conf, now, createdAt, playerLevel)

		state, ok := m.Data[id]
		if !ok {
			state = &ActivityState{
				Id:    id,
				Tasks: make(map[int32]*TaskState, 4),
			}
			m.Data[id] = state
		}

		state.StartAt = start.UnixMilli()
		state.EndAt = end.UnixMilli()
		state.Open = !now.Before(start) && !now.After(end)

		if !state.Open {
			continue
		}

		// 活动开启时按配置补建任务状态，后续事件才能推进进度
		for _, taskEntry := range config.GetActivityTasksByActivityId(id) {
			m.ensureTask(state, taskEntry)
		}
	}
}

// OnEvent 把玩家事件推进到所有开启中的活动任务上
func (m *Manager) OnEvent(e *event.Event) {
	now := utils.NowInLocation()
	playerLevel := m.owner.GetLevel()

	for _, state := range m.Data {
		if !state.IsOpen(now) {
			continue
		}

		for _, task := range state.Tasks {
			m.progressTask(task, e, playerLevel)
		}
	}
}

func (m *Manager) progressTask(task *TaskState, e *event.Event, playerLevel int32) {
	for _, sc := range task.SubConditions {
		if sc.Completed {
			continue
		}

		if subcondition.SupportedEventTypes[sc.Type] != e.Type {
			continue
		}

		ok, err := subcondition.Eval(sc, e, playerLevel)
		if err != nil {
			utils.ErrPrint(err, "eval activity sub condition failed", task.Id, sc.Id)
			continue
		}

		if !ok {
			continue
		}

		sc.Count = sc.MaxCount
		if sc.Count <= 0 {
			sc.Count = 1
		}
		sc.Completed = true
	}
}

// TaskCompleted 判断任务是否全部子条件达成
func (m *Manager) TaskCompleted(task *TaskState) bool {
	if len(task.SubConditions) == 0 {
		return false
	}

	for _, sc := range task.SubConditions {
		if !sc.Completed {
			return false
		}
	}

	return true
}

// ensureTask 按配置补建任务状态（新开启的活动或新增任务时）
func (m *Manager) ensureTask(state *ActivityState, entry *config.ActivityTaskEntry) *TaskState {
	if task, ok := state.Tasks[entry.Id]; ok {
		return task
	}

	task := &TaskState{
		Id:            entry.Id,
		SubConditions: subcondition.BuildFromConditionId(entry.SubConditionId),
	}

	state.Tasks[entry.Id] = task
	return task
}

// ClaimTask 领取活动任务奖励
func (m *Manager) ClaimTask(activityId string, taskId int32) error {
	state, ok := m.Data[activityId]
	if !ok {
		return ErrActivityNotFound
	}

	if !state.IsOpen(utils.NowInLocation()) {
		return ErrActivityClosed
	}

	entry, ok := config.GetActivityTaskEntry(taskId)
	if !ok || entry.ActivityId != activityId {
		return ErrTaskNotFound
	}

	task := m.ensureTask(state, entry)

	if !m.TaskCompleted(task) {
		return ErrTaskNotCompleted
	}

	if task.Claimed {
		return ErrTaskAlreadyClaimed
	}

	task.Claimed = true
	m.owner.GrantAwardList(entry.Awards, "activityTask")

	if err := m.owner.SaveActivities(m.Data); err != nil {
		return err
	}

	m.SendUpdate()
	return nil
}

// SignIn 活动签到，day 必须等于当前已签到天数 + 1
func (m *Manager) SignIn(activityId string, day int32) error {
	state, ok := m.Data[activityId]
	if !ok {
		return ErrActivityNotFound
	}

	if !state.IsOpen(utils.NowInLocation()) {
		return ErrActivityClosed
	}

	if day != state.SignDays+1 {
		if day <= state.SignDays {
			return ErrSignAlreadyDone
		}

		return fmt.Errorf("%w: day %d, expect %d", ErrSignInvalidDay, day, state.SignDays+1)
	}

	entry, ok := config.GetActivitySignInEntry(activityId, fmt.Sprintf("%d", day))
	if !ok {
		return ErrSignInvalidDay
	}

	state.SignDays = day
	m.owner.GrantAwardList(entry.Awards, "activitySignIn")

	if err := m.owner.SaveActivities(m.Data); err != nil {
		return err
	}

	m.SendUpdate()
	return nil
}

// GenPB 生成下发协议
func (m *Manager) GenPB() []*pbGlobal.ActivityInfo {
	now := utils.NowInLocation()

	list := make([]*pbGlobal.ActivityInfo, 0, len(m.Data))
	for _, state := range m.Data {
		if !state.IsOpen(now) {
			continue
		}

		info := &pbGlobal.ActivityInfo{
			Id:       state.Id,
			StartAt:  state.StartAt,
			EndAt:    state.EndAt,
			Open:     true,
			SignDays: state.SignDays,
		}

		if bp := state.BattlePass; bp != nil {
			info.BattlePass = &pbGlobal.BattlePassInfo{
				Score:     bp.Score,
				Unlocked:  bp.Unlocked,
				FreeIndex: bp.FreeIndex,
				HighIndex: bp.HighIndex,
			}
		}

		for id, task := range state.Tasks {
			ti := &pbGlobal.ActivityTaskInfo{Id: id, Claimed: task.Claimed}
			for _, sc := range task.SubConditions {
				ti.Objs = append(ti.Objs, sc.GenPB())
			}

			info.Tasks = append(info.Tasks, ti)
		}

		list = append(list, info)
	}

	return list
}

// SendUpdate 下发活动状态
func (m *Manager) SendUpdate() {
	m.owner.SendProtoMessage(&pbGlobal.S2C_ActivityUpdate{Activities: m.GenPB()})
}

////////////////////////////////////////////////////////
// 战斗通行证

// battlePass 返回（必要时创建）通行证状态
func (s *ActivityState) battlePass() *BattlePassState {
	if s.BattlePass == nil {
		s.BattlePass = &BattlePassState{}
	}

	return s.BattlePass
}

// AddBattlePassPoint 增加通行证积分
func (m *Manager) AddBattlePassPoint(activityId string, point int32) error {
	if point <= 0 {
		return nil
	}

	state, ok := m.Data[activityId]
	if !ok {
		return ErrActivityNotFound
	}

	if !state.IsOpen(utils.NowInLocation()) {
		return ErrActivityClosed
	}

	if len(config.GetActivityBattlePassLevels(activityId)) == 0 {
		return ErrActivityNotFound
	}

	bp := state.battlePass()
	bp.Score += point

	return m.owner.SaveActivities(m.Data)
}

// UnlockBattlePass 解锁高级奖励（支付校验由调用方完成）
func (m *Manager) UnlockBattlePass(activityId string) error {
	state, ok := m.Data[activityId]
	if !ok {
		return ErrActivityNotFound
	}

	if !state.IsOpen(utils.NowInLocation()) {
		return ErrActivityClosed
	}

	bp := state.battlePass()
	if bp.Unlocked {
		return nil
	}

	bp.Unlocked = true

	if err := m.owner.SaveActivities(m.Data); err != nil {
		return err
	}

	m.SendUpdate()
	return nil
}

// ClaimBattlePass 领取所有已达成的通行证阶段奖励。
//
// 与 kof 一致：一次领完从上次位置到当前积分之间的所有阶段，避免逐级点选。
func (m *Manager) ClaimBattlePass(activityId string, free bool) error {
	state, ok := m.Data[activityId]
	if !ok {
		return ErrActivityNotFound
	}

	if !state.IsOpen(utils.NowInLocation()) {
		return ErrActivityClosed
	}

	levels := config.GetActivityBattlePassLevels(activityId)
	if len(levels) == 0 {
		return ErrActivityNotFound
	}

	bp := state.battlePass()
	if !free && !bp.Unlocked {
		return ErrBattlePassLocked
	}

	index := bp.FreeIndex
	if !free {
		index = bp.HighIndex
	}

	var awards []*config.IdCount

	for int(index) < len(levels) {
		level := levels[index]
		if level.Score > bp.Score {
			break
		}

		if free {
			awards = append(awards, level.AwardItems...)
		} else {
			awards = append(awards, level.HightAwardItems...)
		}

		index++
	}

	if index == bp.FreeIndex && free {
		return ErrBattlePassNoReward
	}

	if index == bp.HighIndex && !free {
		return ErrBattlePassNoReward
	}

	if free {
		bp.FreeIndex = index
	} else {
		bp.HighIndex = index
	}

	m.owner.GrantAwardList(awards, "battlePass")

	if err := m.owner.SaveActivities(m.Data); err != nil {
		return err
	}

	m.SendUpdate()
	return nil
}
