package scene

import (
	log "github.com/rs/zerolog/log"
)

type SceneCreature struct {
	id   uint64
	opts *UnitOptions
}

func (c *SceneCreature) Guid() uint64 {
	return c.id
}

func (c *SceneCreature) CombatCtrl() *CombatCtrl {
	return c.opts.CombatCtl
}

func (c *SceneCreature) Opts() *UnitOptions {
	return c.opts
}

func (c *SceneCreature) UpdateSpell() {
	log.Info().
		Int64("id", c.id).
		Uint32("type_id", c.opts.TypeId).
		Int32("pos", c.opts.Position).
		Msg("creature start UpdateSpell")

	c.CombatCtrl().Update()
}
