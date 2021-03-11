package scene

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
)

type SceneManager struct {
	mapScenes map[int64]*Scene
	scenePool sync.Pool

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func NewSceneManager() *SceneManager {
	m := &SceneManager{
		mapScenes: make(map[int64]*Scene, define.Scene_MaxNumPerCombat),
	}

	m.scenePool.New = func() interface{} {
		return NewScene()
	}

	return m
}

func (m *SceneManager) createEntryScene(sceneId int64, opts ...SceneOption) (*Scene, error) {
	s := m.scenePool.Get().(*Scene)
	s.Init(sceneId, opts...)
	return s, nil
}

func (m *SceneManager) CreateScene(ctx context.Context, sceneId int64, sceneType int32, opts ...SceneOption) (*Scene, error) {
	if sceneType < define.Scene_TypeBegin || sceneType >= define.Scene_TypeEnd {
		return nil, fmt.Errorf("unknown scene type<%d>", sceneType)
	}

	if len(m.mapScenes) >= define.Scene_MaxNumPerCombat {
		return nil, errors.New("full of scene instance")
	}

	// compile comment
	// var entry *auto.SceneEntry
	// var ok bool
	// if sceneType == define.Scene_TypeStage {
	// 	if entry, ok = auto.GetSceneEntry(1); ok {
	// 		return nil, fmt.Errorf("invalid scene entry by id<%d>", 1)
	// 	}
	// }

	// newOpts := append(opts, WithSceneEntry(entry))
	s, err := m.createEntryScene(sceneId, opts...)
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
		m.DestroyScene(s)
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

func (m *SceneManager) DestroyScene(s *Scene) {
	m.Lock()
	defer m.Unlock()

	delete(m.mapScenes, s.id)
	m.scenePool.Put(s)
}
