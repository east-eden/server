package scene

import "sync"

var (
	scenePool  sync.Pool // 场景池
	skillPool  sync.Pool // 技能池
	buffPool   sync.Pool // aura池
	actionPool sync.Pool // 行动池
)

func init() {
	scenePool.New = func() interface{} {
		return &Scene{}
	}

	skillPool.New = func() interface{} {
		return &Skill{}
	}

	buffPool.New = func() interface{} {
		return &Buff{}
	}

	actionPool.New = func() interface{} {
		return &Action{}
	}
}

func NewScene() *Scene {
	return scenePool.Get().(*Scene)
}

func ReleaseScene(s *Scene) {
	scenePool.Put(s)
}

func NewSkill() *Skill {
	return skillPool.Get().(*Skill)
}

func ReleaseSkill(s *Skill) {
	skillPool.Put(s)
}

func NewBuff() *Buff {
	return buffPool.Get().(*Buff)
}

func ReleaseBuff(b *Buff) {
	buffPool.Put(b)
}

func NewAction() *Action {
	return actionPool.Get().(*Action)
}

func ReleaseAction(a *Action) {
	actionPool.Put(a)
}
