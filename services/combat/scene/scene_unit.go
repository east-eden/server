package scene

type SceneUnit interface {
	Guid() uint64
	UpdateSpell()
	CombatCtrl() *CombatCtrl
	Opts() *UnitOptions
}
