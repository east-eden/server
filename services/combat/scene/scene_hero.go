package scene

import (
	log "github.com/rs/zerolog/log"
	"github.com/willf/bitset"
	"github.com/yokaiio/yokai_server/define"
)

type SceneHero struct {
	id    uint64
	level uint32
	opts  *UnitOptions
}

func (h *SceneHero) Guid() uint64 {
	return h.id
}

func (h *SceneHero) GetLevel() uint32 {
	return h.level
}

func (h *SceneHero) GetScene() *Scene {
	return h.opts.Scene
}

func (h *SceneHero) GetCamp() int32 {
	return 0
}

func (h *SceneHero) CombatCtrl() *CombatCtrl {
	return h.opts.CombatCtrl
}

func (h *SceneHero) Opts() *UnitOptions {
	return h.opts
}

func (h *SceneHero) UpdateSpell() {
	log.Info().
		Uint64("id", h.id).
		Int32("type_id", h.opts.TypeId).
		Floats32("pos", h.opts.Position[:]).
		Msg("hero start UpdateSpell")

	h.CombatCtrl().Update()
}

func (h *SceneHero) HasState(e define.EHeroState) bool {
	return h.opts.State.Test(uint(e))
}

func (h *SceneHero) HasStateAny(flag uint32) bool {
	compare := bitset.From([]uint64{uint64(flag)})
	return h.opts.State.Intersection(compare).Any()
}

func (h *SceneHero) GetState64() uint64 {
	return h.opts.State.Bytes()[0]
}

func (h *SceneHero) HasImmunityAny(tp define.EImmunityType, flag uint32) bool {
	compare := bitset.From([]uint64{uint64(flag)})
	return h.opts.Immunity[tp].Intersection(compare).Any()
}

func (h *SceneHero) BeatBack(target SceneUnit) {

}

func (h *SceneHero) DoneDamage(caster SceneUnit, dmgInfo *CalcDamageInfo) {

}

func (h *SceneHero) OnDamage(target SceneUnit, dmgInfo *CalcDamageInfo) {

}

func (h *SceneHero) OnBeDamaged(caster SceneUnit, dmgInfo *CalcDamageInfo) {

}
