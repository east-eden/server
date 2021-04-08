package player

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/valyala/bytebufferpool"
)

var (
	chapterStageUpdateInterval  = time.Second * 5 // 每5秒更新一次
	ErrInvalidRequest           = errors.New("invalid request")
	ErrChapterNotFound          = errors.New("chapter not found")
	ErrChapterRewardAlready     = errors.New("chapter reward received already")
	ErrChapterStarsNotEnough    = errors.New("chapter stars not enough")
	ErrStageNotFound            = errors.New("stage not found")
	ErrStagePrevNotPassed       = errors.New("prev stage not passed")
	ErrStageChallengeTimesLimit = errors.New("stage challenge times limit")
)

func makeChapterKey(chapterId int32, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("chapter_list.")
	_, _ = b.WriteString(strconv.Itoa(int(chapterId)))

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
	_, _ = b.WriteString(strconv.Itoa(int(stageId)))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

// 章节
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

// 关卡
type Stage struct {
	define.StageInfo `bson:"inline" json:",inline"`
}

func (s *Stage) GenStagePB() *pbGlobal.Stage {
	pb := &pbGlobal.Stage{
		Id:             s.Id,
		ChallengeTimes: int32(s.ChallengeTimes),
		FirstReward:    s.FirstReward,
		Objectives:     make([]bool, define.Stage_Objective_Num),
	}

	copy(pb.Objectives, s.Objectives[:])
	return pb
}

type ChapterStageManager struct {
	owner         *Player            `bson:"-" json:"-"`
	nextUpdate    int64              `bson:"-" json:"-"`                             // 下次更新时间
	Chapters      map[int32]*Chapter `bson:"chapter_list" json:"chapter_list"`       // 章节数据
	Stages        map[int32]*Stage   `bson:"stage_list" json:"stage_list"`           // 关卡数据
	LastResetTime int32              `bson:"last_reset_time" json:"last_reset_time"` // 上次重置关卡时间
}

func NewChapterStageManager(owner *Player) *ChapterStageManager {
	m := &ChapterStageManager{
		owner:         owner,
		Chapters:      make(map[int32]*Chapter),
		Stages:        make(map[int32]*Stage),
		LastResetTime: int32(time.Now().Unix()),
	}

	return m
}

func (m *ChapterStageManager) update() {
	if time.Now().Unix() < m.nextUpdate {
		return
	}

	// 设置下次更新时间
	m.nextUpdate = time.Now().Add(chapterStageUpdateInterval).Unix()

	// todo 重置章节关卡数据

	// save
	m.LastResetTime = int32(time.Now().Unix())
	fields := map[string]interface{}{
		"last_reset_time": m.LastResetTime,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when ChapterStageManager.update", m.owner.ID, fields)

	// todo sync to client
}

// 关卡通关
func (m *ChapterStageManager) StagePass(stageId int32, objectives []bool) error {
	stageEntry, ok := auto.GetStageEntry(stageId)
	if !ok {
		return ErrStageNotFound
	}

	// 前置关卡限制
	_, ok = m.Stages[stageEntry.PrevStageId]
	if stageEntry.PrevStageId != -1 && !ok {
		return ErrStagePrevNotPassed
	}

	// 条件限制
	if !m.owner.ConditionManager().CheckCondition(stageEntry.ConditionId) {
		return ErrConditionLimit
	}

	stage, stageExist := m.Stages[stageId]

	// 挑战次数限制
	if stageExist && stage.ChallengeTimes >= int16(stageEntry.DailyTimes) {
		return ErrStageChallengeTimesLimit
	}

	// todo 通用限制

	// 通关处理
	if !stageExist {
		stage = &Stage{
			StageInfo: define.StageInfo{
				Id:             stageId,
				ChallengeTimes: 0,
				FirstReward:    false,
			},
		}

		// 首次通关奖励
		err := m.owner.CostLootManager().GainLoot(stageEntry.FirstRewardLootId)
		utils.ErrPrint(err, "GainLoot first reward failed when ChapterStageManager.StagePass", m.owner.ID, stageId)
	}

	// 通关奖励
	err := m.owner.CostLootManager().GainLoot(stageEntry.RewardLootId)
	utils.ErrPrint(err, "GainLoot failed when ChapterStageManager.StagePass", m.owner.ID, stageId)

	// 更新关卡目标达成状况
	var addStar int32
	for k := range objectives {
		if !objectives[k] {
			continue
		}

		if stage.Objectives[k] {
			continue
		}

		stage.Objectives[k] = true
		addStar++
	}

	// 更新关卡数据
	stage.ChallengeTimes++
	m.Stages[stage.Id] = stage

	// 更新章节
	chapter, chapterExist := m.Chapters[stageEntry.ChapterId]
	if addStar > 0 {
		if chapterExist {
			chapter.Stars += addStar
		} else {
			chapter = &Chapter{
				ChapterInfo: define.ChapterInfo{
					Id:    stageEntry.ChapterId,
					Stars: addStar,
				},
			}
			m.Chapters[chapter.Id] = chapter
		}
	}

	fields := map[string]interface{}{
		makeChapterKey(chapter.Id): chapter,
		makeStageKey(stage.Id):     stage,
		"last_reset_time":          m.LastResetTime,
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when ChapterStageManager.StagePass", m.owner.ID, fields)
	return nil
}

// 领取章节奖励
func (m *ChapterStageManager) ReceiveChapterReward(chapterId int32, index int32) error {
	chapterEntry, ok := auto.GetChapterEntry(chapterId)
	if !ok {
		return ErrChapterNotFound
	}

	// 条件限制
	if !m.owner.ConditionManager().CheckCondition(chapterEntry.ConditionId) {
		return ErrConditionLimit
	}

	chapter, exist := m.Chapters[chapterId]
	if !exist {
		return ErrChapterNotFound
	}

	if !utils.BetweenInt32(index, 0, define.Chapter_Rewards_Num) {
		return ErrInvalidRequest
	}

	// 已领过
	if chapter.Rewards[index] {
		return ErrChapterRewardAlready
	}

	// 星数不够
	if chapter.Stars < chapterEntry.RewardStars[index] {
		return ErrChapterStarsNotEnough
	}

	err := m.owner.CostLootManager().GainLoot(chapterEntry.RewardStars[index])
	utils.ErrPrint(err, "GainLoot failed when ChapterStageManager.ReceiveChapterReward", m.owner.ID, chapterId, index)

	chapter.Rewards[index] = true

	fields := map[string]interface{}{
		makeChapterKey(chapterId, fmt.Sprintf("rewards.%d", index)): chapter.Rewards[index],
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when ChapterStageManager.ReceiveChapterReward", m.owner.ID, fields)

	return nil
}

// 关卡扫荡
func (m *ChapterStageManager) StageSweep(stageId int32, times int32) error {
	if times < 0 {
		return ErrInvalidRequest
	}

	globalConfig, _ := auto.GetGlobalConfig()

	stageEntry, ok := auto.GetStageEntry(stageId)
	if !ok {
		return ErrStageNotFound
	}

	// 没通关不能扫荡
	stage, exist := m.Stages[stageId]
	if !exist {
		return ErrStageNotFound
	}

	for n := 0; n < int(times); n++ {
		// 挑战次数限制
		if int32(stage.ChallengeTimes)+1 >= stageEntry.DailyTimes {
			break
		}

		// todo 判断体力

		// 扫荡券
		err := m.owner.ItemManager().CanCost(globalConfig.SweepStageItem, 1)
		if err != nil {
			break
		}

		err = m.owner.ItemManager().DoCost(globalConfig.SweepStageItem, 1)
		utils.ErrPrint(err, "ItemManager.DoCost failed when StageSweep", m.owner.ID)

		err = m.owner.CostLootManager().GainLoot(stageEntry.RewardLootId)
		utils.ErrPrint(err, "GainLoot failed when StageSweep", m.owner.ID, stageEntry.RewardLootId)

		stage.ChallengeTimes++
	}

	// save
	fields := map[string]interface{}{
		makeStageKey(stageId, "challenge_times"): stage.ChallengeTimes,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, m.owner.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when ChapterStageManager.StageSweep", m.owner.ID, fields)
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
	msg := &pbGlobal.S2C_ChapterUpdate{
		Chapter: c.GenChapterPB(),
	}
	m.owner.SendProtoMessage(msg)
}

func (m *ChapterStageManager) SendStageUpdate(s *Stage) {
	msg := &pbGlobal.S2C_StageUpdate{
		Stage: s.GenStagePB(),
	}
	m.owner.SendProtoMessage(msg)
}
