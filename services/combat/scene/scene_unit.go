package scene

import "github.com/yokaiio/yokai_server/define"

type SceneUnit interface {
	Guid() uint64
	GetLevel() uint32
	GetScene() *Scene
	UpdateSpell()
	CombatCtrl() *CombatCtrl
	Opts() *UnitOptions

	HasState(define.EHeroState) bool
	HasImmunityAny(define.EImmunityType, uint32)
	GetState64() uint64
	BeatBack(SceneUnit)
}
