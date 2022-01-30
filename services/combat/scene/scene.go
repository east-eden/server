package scene

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/random"
	"github.com/emirpasic/gods/maps/treemap"
	god_utils "github.com/emirpasic/gods/utils"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
)

var (
	ErrSceneModelNotFound = errors.New("model not found")
)

type Scene struct {
	opts   *SceneOptions
	tasker *task.Tasker

	id          int64
	entityIdGen int64
	entityMap   *treemap.Map // 战斗unit列表
	curRound    int32
	maxRound    int32
	result      chan bool
	rand        *random.FakeRandom
	camps       [define.Scene_Camp_End]*SceneCamp

	comFinishList *list.List // com条结束的entity列表
	spellList     *list.List // 场景内技能列表

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func (s *Scene) Init(sceneId int64, opts ...SceneOption) *Scene {
	s.id = sceneId
	s.entityMap = treemap.NewWith(god_utils.Int64Comparator)
	s.comFinishList = list.New()
	s.spellList = list.New()
	s.result = make(chan bool, 1)
	s.opts = DefaultSceneOptions()
	s.rand = random.NewFakeRandom(int(time.Now().Unix()))
	s.tasker = task.NewTasker()

	for n := define.Scene_Camp_Begin; n < define.Scene_Camp_End; n++ {
		s.camps[n] = NewSceneCamp(s, n)
	}

	for _, o := range opts {
		o(s.opts)
	}

	// add attack unit list
	for _, unitInfo := range s.opts.AttackEntityList {
		err := s.AddEntityByPB(s.camps[define.Scene_Camp_Attack], unitInfo)
		utils.ErrPrint(err, "AddEntityByPB failed when Scene.Init", sceneId, s.opts.SceneEntry.Id, unitInfo.HeroTypeId)
	}

	// add defence unit list
	for _, unitInfo := range s.opts.DefenceEntityList {
		err := s.AddEntityByPB(s.camps[define.Scene_Camp_Defence], unitInfo)
		_ = utils.ErrCheck(err, "AddEntityByPB failed when Scene.Init", sceneId, s.opts.SceneEntry.Id, unitInfo.HeroTypeId)
	}

	// 目前只有一波
	battleWaveEntry := s.opts.BattleWaveEntries[0]
	if battleWaveEntry != nil {
		for idx := range battleWaveEntry.MonsterID {
			// hero id invalid
			if battleWaveEntry.MonsterID[idx] == -1 {
				continue
			}

			monsterEntry, ok := auto.GetMonsterEntry(battleWaveEntry.MonsterID[idx])
			if !ok {
				continue
			}

			err := s.AddEntityByOptions(
				s.camps[define.Scene_Camp_Defence],
				WithEntityMonsterId(battleWaveEntry.MonsterID[idx]),
				WithEntityMonsterEntry(monsterEntry),
				WithEntityPosition(battleWaveEntry.PositionX[idx], battleWaveEntry.PositionZ[idx], battleWaveEntry.Rotation[idx]),
				WithEntityInitAtbValue(battleWaveEntry.InitalCom[idx]),
			)

			_ = utils.ErrCheck(err, "AddEntityByOptions failed when Scene.Init", battleWaveEntry.MonsterID[idx])
		}
	}

	// tasker init
	s.tasker.Init(
		task.WithStartFns(func() {
			s.onTaskStart()
		}),

		task.WithStopFns(func() {
			s.onTaskStop()
		}),

		task.WithUpdateFn(func() {
			s.onTaskUpdate()
		}),
	)

	return s
}

func (s *Scene) onTaskStart() {
	it := s.entityMap.Iterator()
	for it.Next() {
		it.Value().(*SceneEntity).OnSceneStart()
	}
}

func (s *Scene) onTaskStop() {
	log.Info().
		Int32("scene_type_id", s.opts.SceneEntry.Id).
		Int64("scene_id", s.GetId()).
		Msg("scene context done...")
}

func (s *Scene) onTaskUpdate() {
	s.updateEntities()
	s.updateCamps()
}

func (s *Scene) TaskRun(ctx context.Context) error {
	return s.tasker.Run(ctx)
}

func (s *Scene) TaskStop() {
	s.tasker.Stop()
}

func (s *Scene) Exit(ctx context.Context) {
	s.wg.Wait()
}

func (s *Scene) GetId() int64 {
	return s.id
}

func (s *Scene) GetCamp(camp int32) *SceneCamp {
	return s.camps[camp]
}

func (s *Scene) GetEntity(id int64) (*SceneEntity, bool) {
	val, ok := s.entityMap.Get(id)
	if ok {
		return val.(*SceneEntity), ok
	}

	return nil, ok
}

func (s *Scene) GetEntityMap() *treemap.Map {
	return s.entityMap
}

// 寻找单位
func (s *Scene) findEnemyEntityByHead(camp int32) (*SceneEntity, bool) {
	if s.entityMap.Size() == 0 {
		return nil, false
	}

	it := s.entityMap.Iterator()
	for it.Next() {
		e := it.Value().(*SceneEntity)
		if e.GetCamp().camp != camp {
			return e, true
		}
	}

	return nil, false
}

func (s *Scene) GetResult() bool {
	return <-s.result
}

func (s *Scene) GetSceneCamp(camp int32) (*SceneCamp, bool) {
	if !utils.BetweenInt32(camp, define.Scene_Camp_Begin, define.Scene_Camp_End) {
		return nil, false
	}

	return s.camps[camp], true
}

func (s *Scene) Rand(min, max int) int {
	return s.rand.RandSection(min, max)
}

func (s *Scene) GetRand() *random.FakeRandom {
	return s.rand
}

func (s *Scene) updateCamps() {

	// 是否攻击方先手
	bAttackFirst := true

	for ; s.curRound+1 <= s.maxRound; s.curRound++ {
		bEnterNextRound := false

		nActionRount := 0
		for !bEnterNextRound {
			nActionRount++

			if bAttackFirst {
				s.camps[int(define.Scene_Camp_Attack)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Attack)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Attack)].Attack(s.camps[int(define.Scene_Camp_Defence)])
				}

				s.camps[int(define.Scene_Camp_Defence)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Defence)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Defence)].Attack(s.camps[int(define.Scene_Camp_Attack)])
				}
			} else {
				s.camps[int(define.Scene_Camp_Defence)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Defence)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Defence)].Attack(s.camps[int(define.Scene_Camp_Defence)])
				}

				s.camps[int(define.Scene_Camp_Attack)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Attack)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Attack)].Attack(s.camps[int(define.Scene_Camp_Defence)])
				}
			}

			if s.camps[int(define.Scene_Camp_Attack)].IsLoopEnd() &&
				s.camps[int(define.Scene_Camp_Defence)].IsLoopEnd() {
				for i := nActionRount; i < Camp_Max_Unit; i++ {
					s.camps[int(define.Scene_Camp_Defence)].Update()
					s.camps[int(define.Scene_Camp_Attack)].Update()
				}

				nActionRount = 0
				bEnterNextRound = true
			}

			// 释放符文技能
			// compile comment
			// s.UpdateRuneSpell(bAttackFirst);
		}

		// 重置攻击顺序
		s.camps[int(define.Scene_Camp_Attack)].ResetLoopIndex()
		s.camps[int(define.Scene_Camp_Defence)].ResetLoopIndex()

		// 战斗结束
		if !s.camps[int(define.Scene_Camp_Attack)].IsValid() ||
			!s.camps[int(define.Scene_Camp_Defence)].IsValid() {
			break
		}
	}
}

// 更新场景内技能
func (s *Scene) updateSpells() {
	var next *list.Element
	for e := s.spellList.Front(); e != nil; e = next {
		next = e.Next()

		skill := e.Value.(*Skill)
		skill.Update()

		// 删除已作用玩的技能
		if skill.IsCompleted() {
			s.spellList.Remove(e)
		}
	}
}

func (s *Scene) updateEntities() {
	it := s.entityMap.Iterator()
	for it.Next() {
		it.Value().(*SceneEntity).Update()
	}
}

func (s *Scene) AddEntityByPB(camp *SceneCamp, unitInfo *pbGlobal.EntityInfo) error {
	entry, ok := auto.GetHeroEntry(unitInfo.HeroTypeId)
	if !ok {
		return fmt.Errorf("GetUnitEntry failed: type_id<%d>", unitInfo.HeroTypeId)
	}

	modelEntry, ok := auto.GetModelEntry(entry.ModelID)
	if !ok {
		return fmt.Errorf("err:<%w>, model_id:<%d>", ErrSceneModelNotFound, entry.ModelID)
	}

	id := atomic.AddInt64(&s.entityIdGen, 1)
	e, err := NewSceneEntity(
		s,
		id,
		WithEntityHeroId(unitInfo.HeroTypeId),
		WithEntityAttList(unitInfo.AttValue),
		WithEntityHeroEntry(entry),
		WithEntityModelEntry(modelEntry),
	)

	if err != nil {
		return err
	}

	s.entityMap.Put(id, e)

	return nil
}

func (s *Scene) AddEntityByOptions(camp *SceneCamp, opts ...EntityOption) error {
	id := atomic.AddInt64(&s.entityIdGen, 1)
	opts = append(opts, WithEntitySceneCamp(camp))
	e, err := NewSceneEntity(s, id, opts...)
	if err != nil {
		return err
	}

	s.entityMap.Put(id, e)
	return nil
}

func (s *Scene) ClearEntities() {
	s.entityMap.Clear()
}
