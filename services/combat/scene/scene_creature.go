package scene

import (
	log "github.com/rs/zerolog/log"
	"github.com/willf/bitset"
	"github.com/east-eden/server/define"
)

type SceneCreature struct {
	id    uint64
	level uint32
	opts  *UnitOptions
}

func (c *SceneCreature) Guid() uint64 {
	return c.id
}

func (c *SceneCreature) GetLevel() uint32 {
	return c.level
}

func (c *SceneCreature) GetScene() *Scene {
	return c.opts.Scene
}

func (c *SceneCreature) GetCamp() int32 {
	return 0
}

func (c *SceneCreature) CombatCtrl() *CombatCtrl {
	return c.opts.CombatCtrl
}

func (c *SceneCreature) Opts() *UnitOptions {
	return c.opts
}

func (c *SceneCreature) UpdateSpell() {
	log.Info().
		Uint64("id", c.id).
		Int32("type_id", c.opts.TypeId).
		Floats32("pos", c.opts.Position[:]).
		Msg("creature start UpdateSpell")

	c.CombatCtrl().Update()
}

func (c *SceneCreature) HasState(e define.EHeroState) bool {
	return c.opts.State.Test(uint(e))
}

func (h *SceneCreature) HasStateAny(flag uint32) bool {
	compare := bitset.From([]uint64{uint64(flag)})
	return h.opts.State.Intersection(compare).Any()
}

func (h *SceneCreature) GetState64() uint64 {
	return h.opts.State.Bytes()[0]
}

func (h *SceneCreature) HasImmunityAny(tp define.EImmunityType, flag uint32) bool {
	compare := bitset.From([]uint64{uint64(flag)})
	return h.opts.Immunity[tp].Intersection(compare).Any()
}

func (h *SceneCreature) BeatBack(target SceneUnit) {

}

func (h *SceneCreature) DoneDamage(caster SceneUnit, dmgInfo *CalcDamageInfo) {

}

func (h *SceneCreature) OnDamage(target SceneUnit, dmgInfo *CalcDamageInfo) {

}

func (h *SceneCreature) OnBeDamaged(caster SceneUnit, dmgInfo *CalcDamageInfo) {

}
