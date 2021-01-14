package scene

import log "github.com/rs/zerolog/log"

type MoveCtrl struct {
	owner *SceneUnit // 拥有者
}

func NewMoveCtrl(owner *SceneUnit) *MoveCtrl {
	c := &MoveCtrl{
		owner: owner,
	}

	return c
}

func (c *MoveCtrl) Update() {
	log.Info().Int64("owner_id", c.owner.id).Msg("MoveCtrl update")
}
