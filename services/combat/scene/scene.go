package scene

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/east-eden/server/auto"
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

type Scene struct {
	opts *SceneOptions

	id        int64
	result    chan bool
	mapUnits  map[uint64]SceneUnit
	unitIdGen uint64
	rand      *utils.FakeRandom

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func newScene(sceneId int64, opts ...SceneOption) *Scene {
	s := &Scene{
		id:       sceneId,
		mapUnits: make(map[uint64]SceneUnit, define.Scene_MaxUnitPerScene),
		result:   make(chan bool, 1),
		opts:     DefaultSceneOptions(),
		rand:     utils.NewFakeRandom(int(time.Now().Unix())),
	}

	for _, o := range opts {
		o(s.opts)
	}

	// add attack unit list
	for _, unit := range s.opts.AttackUnitList {
		entry, _ := auto.GetUnitEntry(unit.UnitTypeId)
		s.addHero(
			WithUnitTypeId(unit.UnitTypeId),
			WithUnitAttList(unit.UnitAttList),
			WithUnitEntry(entry),
		)
	}

	// add defence unit list
	for _, unit := range s.opts.DefenceUnitList {
		entry, _ := auto.GetUnitEntry(unit.UnitTypeId)
		s.addHero(
			WithUnitTypeId(unit.UnitTypeId),
			WithUnitAttList(unit.UnitAttList),
			WithUnitEntry(entry),
		)
	}

	// add scene unit list
	if s.opts.Entry.UnitGroupID != -1 {
		if groupEntry, ok := auto.GetUnitGroupEntry(s.opts.Entry.UnitGroupID); ok {
			for k, v := range groupEntry.UnitTypeID {
				entry, _ := auto.GetUnitEntry(v)
				s.addCreature(
					WithUnitTypeId(v),
					WithUnitPositionString(groupEntry.Position[k]),
					WithUnitEntry(entry),
				)
			}
		}
	}

	return s
}

func (s *Scene) addHero(opts ...UnitOption) error {
	id := atomic.AddUint64(&s.unitIdGen, 1)
	h := &SceneHero{
		id:   id,
		opts: DefaultUnitOptions(),
	}

	for _, o := range opts {
		o(h.opts)
	}

	h.opts.CombatCtrl = NewCombatCtrl(h)
	s.mapUnits[id] = h
	return nil
}

func (s *Scene) addCreature(opts ...UnitOption) error {
	id := atomic.AddUint64(&s.unitIdGen, 1)
	c := &SceneCreature{
		id:   id,
		opts: DefaultUnitOptions(),
	}

	for _, o := range opts {
		o(c.opts)
	}

	c.opts.CombatCtrl = NewCombatCtrl(c)
	s.mapUnits[id] = c
	return nil
}

func (s *Scene) Main(ctx context.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Error().
					Int64("scene_id", s.id).
					Int32("scene_type", s.opts.Entry.Type).
					Err(err).
					Msg("scene main return error")
			}
			exitCh <- err
		})
	}

	s.wg.Wrap(func() {
		exitFunc(s.Run(ctx))
	})

	s.wg.Wrap(func() {
		log.Info().Msg("scene begin count down 10 seconds")
		time.AfterFunc(time.Second*10, func() {
			log.Info().Msg("scene complete count down 10 seconds")
		})
	})

	return <-exitCh
}

func (s *Scene) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("scene_id", s.id).Msg("scene context done...")
			return nil
		default:
			t := time.Now()

			// update
			s.updateUnits()

			d := time.Since(t)
			time.Sleep(time.Millisecond*200 - d)
		}
	}
}

func (s *Scene) Exit(ctx context.Context) {
	s.wg.Wait()
}

func (s *Scene) GetID() int64 {
	return s.id
}

func (s *Scene) GetResult() bool {
	return <-s.result
}

// todo
func (s *Scene) IsOnlyRecord() bool {
	return false
}

func (s *Scene) Rand(min, max int) int {
	return s.rand.RandSection(min, max)
}

func (s *Scene) updateUnits() {
	s.RLock()
	defer s.RUnlock()

	for _, unit := range s.mapUnits {
		unit.UpdateSpell()
	}
}

func (s *Scene) SendDamage(dmgInfo *CalcDamageInfo) {
	// todo
}
