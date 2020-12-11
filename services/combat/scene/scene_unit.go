package scene

import "github.com/east-eden/server/define"

type SceneUnit interface {
	Guid() uint64
	GetLevel() uint32
	GetScene() *Scene
	GetCamp() int32
	UpdateSpell()
	CombatCtrl() *CombatCtrl
	Opts() *UnitOptions

	HasState(define.EHeroState) bool
	HasStateAny(uint32) bool
	HasImmunityAny(define.EImmunityType, uint32) bool
	GetState64() uint64
	BeatBack(SceneUnit)
	DoneDamage(SceneUnit, *CalcDamageInfo)
	OnDamage(SceneUnit, *CalcDamageInfo)
	OnBeDamaged(SceneUnit, *CalcDamageInfo)
}
