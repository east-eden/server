package scene

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/utils"
)

type Scene struct {
	opts *SceneOptions

	id        int64
	result    chan bool
	mapUnits  map[uint64]SceneUnit
	unitIdGen uint64

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func newScene(sceneId int64, opts ...SceneOption) *Scene {
	s := &Scene{
		id:       sceneId,
		mapUnits: make(map[uint64]SceneUnit, define.Scene_MaxUnitPerScene),
		result:   make(chan bool, 1),
		opts:     DefaultSceneOptions(),
	}

	for _, o := range opts {
		o(s.opts)
	}

	// add attack unit list
	for _, unit := range s.opts.AttackUnitList {
		s.addHero(
			WithUnitTypeId(unit.UnitTypeId),
			WithUnitAttList(unit.UnitAttList),
			WithUnitEntry(entries.GetUnitEntry(unit.UnitTypeId)),
		)
	}

	// add defence unit list
	for _, unit := range s.opts.DefenceUnitList {
		s.addHero(
			WithUnitTypeId(unit.UnitTypeId),
			WithUnitAttList(unit.UnitAttList),
			WithUnitEntry(entries.GetUnitEntry(unit.UnitTypeId)),
		)
	}

	// add scene unit list
	if s.opts.Entry.UnitGroupID != -1 {
		if groupEntry := entries.GetUnitGroupEntry(s.opts.Entry.UnitGroupID); groupEntry != nil {
			for k, v := range groupEntry.UnitTypeID {
				s.addCreature(
					WithUnitTypeId(v),
					WithUnitPositionString(groupEntry.Position[k]),
					WithUnitEntry(entries.GetUnitEntry(v)),
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

	h.CombatCtl = NewCombatCtrl(h)
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

	c.CombatCtl = NewCombatCtrl(c)
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

	return nil
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

func (s *Scene) updateUnits() {
	s.RLock()
	defer s.RUnlock()

	for _, unit := range s.mapUnits {
		unit.UpdateSpell()
	}
}
