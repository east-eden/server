package scene

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	pbCombat "github.com/yokaiio/yokai_server/proto/combat"
	"github.com/yokaiio/yokai_server/utils"
)

type SceneManager struct {
	mapScenes map[int64]*Scene

	wg     utils.WaitGroupWrapper
	ctx    context.Context
	cancel context.CancelFunc
	sync.RWMutex
}

func NewSceneManager(ctx context.Context) *SceneManager {
	m := &SceneManager{
		mapScenes: make(map[int64]*Scene, define.Scene_MaxNumPerCombat),
	}

	m.ctx, m.cancel = context.WithCancel(ctx)

	return m
}

func (m *SceneManager) createEntryScene(sceneId int64, entry *define.SceneEntry, attackId, defenceId int64, attackUnitList, defenceUnitList []*pbCombat.UnitAtt) (*Scene, error) {
	s := newScene(sceneId, entry, attackId, defenceId, attackUnitList, defenceUnitList)
	s.ctx, s.cancel = context.WithCancel(m.ctx)

	return s, nil
}

func (m *SceneManager) CreateScene(sceneId int64, sceneType int32, attackId, defenceId int64, attackUnitList, defenceUnitList []*pbCombat.UnitAtt) (*Scene, error) {
	if sceneType < define.Scene_TypeBegin || sceneType >= define.Scene_TypeEnd {
		return nil, fmt.Errorf("unknown scene type<%d>", sceneType)
	}

	if len(m.mapScenes) >= define.Scene_MaxNumPerCombat {
		return nil, errors.New("full of scene instance")
	}

	var entry *define.SceneEntry
	if sceneType == define.Scene_TypeStage {
		if entry = entries.GetSceneEntry(1); entry == nil {
			return nil, fmt.Errorf("invalid scene entry by id<%d>", 1)
		}
	}

	s, err := m.createEntryScene(sceneId, entry, attackId, defenceId, attackUnitList, defenceUnitList)
	if err != nil {
		return nil, err
	}

	m.Lock()
	m.mapScenes[s.GetID()] = s
	m.Unlock()

	// make scene run
	m.wg.Wrap(func() {
		s.Main()
	})

	return s, nil
}

func (m *SceneManager) Main() error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("SceneManager Main() error:", err)
			}
			exitCh <- err
		})
	}

	m.wg.Wrap(func() {
		exitFunc(m.Run())
	})

	return <-exitCh
}

func (m *SceneManager) Run() error {
	for {
		select {
		case <-m.ctx.Done():
			m.wg.Wait()
			logger.Info("all scenes are closed graceful...")
			return nil
		}
	}

	return nil
}

func (m *SceneManager) GetScene(sceneId int64) *Scene {
	m.RLock()
	defer m.RUnlock()

	return m.mapScenes[sceneId]
}
