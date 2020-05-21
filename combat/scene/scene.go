package scene

import (
	"context"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	pbCombat "github.com/yokaiio/yokai_server/proto/combat"
	"github.com/yokaiio/yokai_server/utils"
)

type Scene struct {
	id              int64
	attackId        int64
	defenceId       int64
	attackUnitList  []*pbCombat.UnitAtt
	defenceUnitList []*pbCombat.UnitAtt
	result          chan bool

	entry    *define.SceneEntry
	mapUnits map[int64]Unit

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func newScene(sceneId int64, entry *define.SceneEntry, attackId, defenceId int64, attackUnitList, defenceUnitList []*pbCombat.UnitAtt) *Scene {
	return &Scene{
		id:              sceneId,
		attackId:        attackId,
		defenceId:       defenceId,
		attackUnitList:  attackUnitList,
		defenceUnitList: defenceUnitList,
		entry:           entry,
		mapUnits:        make(map[int64]Unit, define.Scene_MaxUnitPerScene),
		result:          make(chan bool, 1),
	}
}

func (s *Scene) Main(ctx context.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				logger.WithFields(logger.Fields{
					"scene_id":   s.id,
					"scene_type": s.entry.Type,
					"error":      err,
				}).Error("scene main() return error")
			}
			exitCh <- err
		})
	}

	s.wg.Wrap(func() {
		exitFunc(s.Run(ctx))
	})

	s.wg.Wrap(func() {
		logger.Info("scene begin count down 10 seconds")
		tm := time.NewTimer(time.Second * 10)
		<-tm.C
		logger.Info("scene complete count down 10 seconds")
	})

	return <-exitCh
}

func (s *Scene) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.WithFields(logger.Fields{
				"scene_id": s.id,
			}).Info("scene context done")
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
