package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/spf13/cast"
	"github.com/valyala/bytebufferpool"
)

var (
	chapterStageUpdateInterval = time.Second * 5 // 每5秒更新一次

	ErrInvalidRequest               = errors.New("invalid request")
	ErrChapterNotFound              = errors.New("chapter not found")
	ErrChapterRewardAlready         = errors.New("chapter reward received already")
	ErrChapterStarsNotEnough        = errors.New("chapter stars not enough")
	ErrStageNotFound                = errors.New("stage not found")
	ErrStagePrevNotPassed           = errors.New("prev stage not passed")
	ErrStageChallengeTimesLimit     = errors.New("stage challenge times limit")
	ErrStageChallengeStrengthLimit  = errors.New("stage challenge strength limit")
	ErrStageChallengeSweepCostLimit = errors.New("stage sweep cost limit")
	ErrStageLevelLimit              = errors.New("stage level limit")
	ErrStageNeedBattleResult        = errors.New("stage battle must be settled by server battle result")
)

func makeChapterKey(chapterId int32, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("chapter_list.")
	_, _ = b.WriteString(cast.ToString(chapterId))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

func makeStageKey(stageId int32, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("stage_list.")
	_, _ = b.WriteString(cast.ToString(stageId))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

// Chapter 章节
type Chapter struct {
	define.ChapterInfo `bson:"inline" json:",inline"`
}

func (c *Chapter) GenChapterPB() *pbGlobal.Chapter {
	pb := &pbGlobal.Chapter{
		Id:      c.Id,
		Stars:   c.Stars,
		Rewards: make([]bool, define.Chapter_Rewards_Num),
	}

	copy(pb.Rewards, c.Rewards[:])
	return pb
}

// Stage 关卡
type Stage struct {
	define.StageInfo `bson:"inline" json:",inline"`
}

func (s *Stage) GenStagePB() *pbGlobal.Stage {
	pb := &pbGlobal.Stage{
		Id:             s.Id,
		Passed:         s.Pass,
		ChallengeTimes: int32(s.ChallengeTimes),
		Achieve:        s.Achieve,
		Objectives:     make([]bool, define.Stage_Objective_Num),
	}

	copy(pb.Objectives, s.Objectives[:])
	return pb
}

// ChapterStageManager 关卡管理器。
//
// 关卡配置改用 kof 的 Level 表（chapter/stage/stageLoot/goalAward），
// 奖励发放统一走 kof 的 Loot 表与物品账本。
type ChapterStageManager struct {
	owner      *Player            `bson:"-" json:"-"`
	nextUpdate int64              `bson:"-" json:"-"`                       // 下次更新时间
	Chapters   map[int32]*Chapter `bson:"chapter_list" json:"chapter_list"` // 章节数据
	Stages     map[int32]*Stage   `bson:"stage_list" json:"stage_list"`     // 关卡数据
}

func NewChapterStageManager(owner *Player) *ChapterStageManager {
	return &ChapterStageManager{
		owner:    owner,
		Chapters: make(map[int32]*Chapter),
		Stages:   make(map[int32]*Stage),
	}
}

func (m *ChapterStageManager) update() {
	if time.Now().Unix() < m.nextUpdate {
		return
	}

	m.nextUpdate = time.Now().Add(chapterStageUpdateInterval).Unix()
}

func (m *ChapterStageManager) onDayChange() {
	for _, stage := range m.Stages {
		if stage.ChallengeTimes == 0 {
			continue
		}

		stage.ChallengeTimes = 0
		fields := map[string]any{
			makeStageKey(stage.Id): stage,
		}

		err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
		utils.ErrPrint(err, "UpdateFields failed when ChapterStageManager.onDayChange", m.owner.ID, fields)
	}
}

// prevLevelId 返回同章节内的上一关 id；chapter 内第一节没有前置关卡
func prevLevelId(level *config.LevelEntry) (int32, bool) {
	if level.Stage <= 1 {
		return 0, false
	}

	return config.GetLevelIdByChapterStage(level.Chapter, level.Stage-1)
}

// StageChallenge 客户端本地战斗的关卡结算。
//
// 一旦配置了外部战斗服（服务器战斗），胜负必须由 C2S_BattleResult 提交，
// 这里会拒绝客户端自报的 win，避免改包刷奖励。
func (m *ChapterStageManager) StageChallenge(stageId int32, win bool, achieve bool, objectives []bool) error {
	if rec := m.owner.BattleManager().Record(); rec != nil && rec.IsServerBattle {
		return ErrStageNeedBattleResult
	}

	return m.stageSettle(stageId, win, achieve, objectives)
}

// StagePassByBattle 服务器战斗结算入口（胜负来自战斗服结果）
func (m *ChapterStageManager) StagePassByBattle(stageId int32, win bool) error {
	return m.stageSettle(stageId, win, false, nil)
}

// stageSettle 关卡结算核心逻辑
func (m *ChapterStageManager) stageSettle(stageId int32, win bool, achieve bool, objectives []bool) error {
	level, ok := config.GetLevelEntry(stageId)
	if !ok {
		return ErrStageNotFound
	}

	// 等级限制
	if level.LevelLimit > 0 && m.owner.GetLevel() < level.LevelLimit {
		return ErrStageLevelLimit
	}

	// 前置关卡限制
	if prevId, found := prevLevelId(level); found {
		if prev, passed := m.Stages[prevId]; !passed || !prev.Pass {
			return ErrStagePrevNotPassed
		}
	}

	stage, stageExist := m.Stages[stageId]
	if !stageExist {
		stage = &Stage{
			StageInfo: define.StageInfo{
				Id:             stageId,
				ChallengeTimes: 0,
			},
		}
	}

	if win {
		stage.Pass = true

		// 通关掉落
		m.owner.GrantLootById(level.StageLoot, "stagePass")

		// 派发通关事件供任务系统判定（ChapterStage 类子条件）
		m.owner.eventManager.DispatchEvent(&event.Event{
			Type:  define.Event_Type_StagePass,
			Miscs: []any{stageId, level.Chapter},
		})
	}

	// 星级目标奖励：达成时一次性发放 GoalAward
	if achieve && !stage.Achieve {
		stage.Achieve = true
		m.owner.GrantAwardList(level.GoalAward, "stageGoal")
	}

	// 更新关卡目标达成状况
	var addStar int32
	for k := range objectives {
		if k >= len(stage.Objectives) || !objectives[k] {
			continue
		}

		if stage.Objectives[k] {
			continue
		}

		stage.Objectives[k] = true
		addStar++
	}

	stage.ChallengeTimes++
	m.Stages[stage.Id] = stage

	// 更新章节星数
	chapter, chapterExist := m.Chapters[level.Chapter]
	if addStar > 0 {
		if chapterExist {
			chapter.Stars += addStar
		} else {
			chapter = &Chapter{
				ChapterInfo: define.ChapterInfo{
					Id:    level.Chapter,
					Stars: addStar,
				},
			}
			m.Chapters[chapter.Id] = chapter
		}

		m.SendChapterUpdate(chapter)
	}

	fields := map[string]any{
		makeStageKey(stage.Id): stage,
	}

	if chapterExist {
		fields[makeChapterKey(chapter.Id)] = chapter
	}

	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when ChapterStageManager.StagePass", m.owner.ID, fields)

	m.SendStageUpdate(stage)
	return nil
}

// StageSweep 扫荡已通关关卡
func (m *ChapterStageManager) StageSweep(stageId int32, times int32) error {
	if times < 0 {
		return ErrInvalidRequest
	}

	level, ok := config.GetLevelEntry(stageId)
	if !ok {
		return ErrStageNotFound
	}

	stage, exist := m.Stages[stageId]
	if !exist || !stage.Pass {
		return ErrStageNotFound
	}

	var swept int32
	for i := int32(0); i < times; i++ {
		// 扫荡按次发放通关掉落
		m.owner.GrantLootById(level.StageLoot, "stageSweep")
		swept++
	}

	if swept <= 0 {
		return ErrStageChallengeSweepCostLimit
	}

	stage.ChallengeTimes += int16(swept)

	fields := map[string]any{
		makeStageKey(stageId, "challenge_times"): stage.ChallengeTimes,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when ChapterStageManager.StageSweep", m.owner.ID, fields)

	m.SendStageUpdate(stage)
	return nil
}

func (m *ChapterStageManager) GenChapterListPB() []*pbGlobal.Chapter {
	chapters := make([]*pbGlobal.Chapter, 0, len(m.Chapters))
	for _, c := range m.Chapters {
		chapters = append(chapters, c.GenChapterPB())
	}

	return chapters
}

func (m *ChapterStageManager) GenStageListPB() []*pbGlobal.Stage {
	stages := make([]*pbGlobal.Stage, 0, len(m.Stages))
	for _, s := range m.Stages {
		stages = append(stages, s.GenStagePB())
	}

	return stages
}

func (m *ChapterStageManager) SendChapterUpdate(c *Chapter) {
	m.owner.SendProtoMessage(&pbGlobal.S2C_ChapterUpdate{Chapter: c.GenChapterPB()})
}

func (m *ChapterStageManager) SendStageUpdate(s *Stage) {
	m.owner.SendProtoMessage(&pbGlobal.S2C_StageUpdate{Stage: s.GenStagePB()})
}

var _ = fmt.Sprintf
