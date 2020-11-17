package scene

import (
	"context"
	"errors"
	"fmt"
	"sync"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/utils"
)

type SceneManager struct {
	mapScenes map[int64]*Scene

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func NewSceneManager() *SceneManager {
	m := &SceneManager{
		mapScenes: make(map[int64]*Scene, define.Scene_MaxNumPerCombat),
	}

	return m
}

func (m *SceneManager) createEntryScene(sceneId int64, opts ...SceneOption) (*Scene, error) {
	s := newScene(sceneId, opts...)

	return s, nil
}

func (m *SceneManager) CreateScene(ctx context.Context, sceneId int64, sceneType int32, opts ...SceneOption) (*Scene, error) {
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

	newOpts := append(opts, WithSceneEntry(entry))
	s, err := m.createEntryScene(sceneId, newOpts...)
	if err != nil {
		return nil, err
	}

	m.Lock()
	m.mapScenes[s.GetID()] = s
	m.Unlock()

	// make scene run
	m.wg.Wrap(func() {
		s.Main(ctx)
		s.Exit(ctx)
	})

	log.Info().
		Int64("scene_id", sceneId).
		Int32("scene_type", sceneType).
		Int64("attack_id", s.opts.AttackId).
		Int64("defence_id", s.opts.DefenceId).
		Msg("create a new scene")

	return s, nil
}

func (m *SceneManager) Main(ctx context.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal().Err(err).Msg("SceneManager Main() error")
			}
			exitCh <- err
		})
	}

	m.wg.Wrap(func() {
		exitFunc(m.Run(ctx))
	})

	// test create scene
	m.CreateScene(ctx, 12345, 0)

	return <-exitCh
}

func (m *SceneManager) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("all scenes are closed graceful...")
			return nil
		}
	}
}

func (m *SceneManager) Exit() {
	m.wg.Wait()
}

func (m *SceneManager) GetScene(sceneId int64) *Scene {
	m.RLock()
	defer m.RUnlock()

	return m.mapScenes[sceneId]
}
