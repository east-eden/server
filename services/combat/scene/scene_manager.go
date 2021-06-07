package scene

import (
	"context"
	"errors"
	"sync"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	ErrSceneNumLimit = errors.New("scene num limit")
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

func (m *SceneManager) createEntryScene(opts ...SceneOption) (*Scene, error) {
	s := NewScene()

	sceneId, err := utils.NextID(define.SnowFlake_Scene)
	if err != nil {
		return nil, err
	}

	s.Init(sceneId, opts...)
	return s, nil
}

func (m *SceneManager) CreateScene(ctx context.Context, opts ...SceneOption) (*Scene, error) {
	m.RLock()
	sceneNum := len(m.mapScenes)
	m.RUnlock()
	if sceneNum >= define.Scene_MaxNumPerCombat {
		return nil, ErrSceneNumLimit
	}

	s, err := m.createEntryScene(opts...)
	if err != nil {
		return nil, err
	}

	m.Lock()
	m.mapScenes[s.GetId()] = s
	m.Unlock()

	// make scene run
	m.wg.Wrap(func() {
		defer utils.CaptureException()
		err := s.Run(ctx)
		_ = utils.ErrCheck(err, "scene.Rune failed", s.GetId())
		s.Exit(ctx)
		m.DestroyScene(s)
	})

	log.Info().
		Int64("scene_id", s.GetId()).
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
		defer utils.CaptureException()
		exitFunc(m.Run(ctx))
	})

	return <-exitCh
}

func (m *SceneManager) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("all scenes are closed graceful...")
	return nil
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
	ReleaseScene(s)
}
