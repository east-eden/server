package scene

import (
	log "github.com/rs/zerolog/log"
)

type SceneHero struct {
	id   uint64
	opts *UnitOptions
}

func (h *SceneHero) Guid() uint64 {
	return h.id
}

func (h *SceneHero) CombatCtrl() *CombatCtrl {
	return h.combatCtl
}

func (h *SceneHero) Opts() *UnitOptions {
	return h.opts
}

func (h *SceneHero) UpdateSpell() {
	log.Info().
		Int64("id", h.id).
		Uint32("type_id", h.opts.TypeId).
		Int32("pos", h.opts.Position).
		Msg("hero start UpdateSpell")

	h.CombatCtrl().Update()
}
