package scene

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
)

type SceneUnit interface {
	Entry() *define.UnitEntry
	UpdateSpell()
}

type SceneHero struct {
	id   int64
	opts *UnitOptions
}

func (u *SceneHero) Entry() *define.UnitEntry {
	return u.opts.Entry
}

func (u *SceneHero) UpdateSpell() {
	logger.WithFields(logger.Fields{
		"id":      u.id,
		"type_id": u.opts.TypeId,
		"pos":     u.opts.Position,
	}).Info("scene hero update spell")
}

type SceneCreature struct {
	id   int64
	opts *UnitOptions
}

func (u *SceneCreature) Entry() *define.UnitEntry {
	return u.opts.Entry
}

func (u *SceneCreature) UpdateSpell() {
	logger.WithFields(logger.Fields{
		"id":      u.id,
		"type_id": u.opts.TypeId,
		"pos":     u.opts.Position,
	}).Info("scene creature update spell")
}
