package scene

import (
	"fmt"
	"sync/atomic"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
)

type CombatCtrl struct {
	mapSpells map[uint64]*Spell // 技能列表
	mapAuras  map[uint64]*Aura  // aura列表
	owner     SceneUnit         // 拥有者
	idGen     uint64            // id generator
}

func NewCombatCtrl(owner SceneUnit) *CombatCtrl {
	c := &CombatCtrl{
		mapSpells: make(map[uint64]*Spell, define.Combat_MaxSpell),
		mapAuras:  make(map[uint64]*Spell, define.Combat_MaxAura),
	}

	c.owner = owner
	idGen = 0
	return c
}

func (c *CombatCtrl) CastSpell(spellId uint32, caster, target SceneUnit, triggered bool) error {
	if len(c.mapSpells) >= define.Combat_MaxSpell {
		err := fmt.Errorf("spell list length >= <%d>", len(c.mapSpells))
		log.Warn().Err(err).Uint32("spell_id", spellId).Send()
		return err
	}

	entry := entries.GetSpellEntry(spellId)
	if entry == nil {
		err := fmt.Errorf("get spell entry failed")
		log.Warn().Err(err).Uint32("spell_id", spellId).Send()
		return err
	}

	s := NewSpell(spellId,
		WithSpellEntry(entry),
		WithSpellCaster(caster),
		WithSpellTarget(target),
		WithSpellTriggered(triggered),
	)

	s.Cast()

	c.mapSpells[atomic.AddUint64(&c.idGen, 1)] = s

	return nil
}

func (c *CombatCtrl) Update() {
	for spellGuid, s := range c.mapSpells {
		log.Info().Uint32("spell_id", s.GetOptions().Id).Msg("spell update...")
		delete(c.mapSpells, spellGuid)
	}

}
