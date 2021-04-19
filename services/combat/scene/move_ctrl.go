package scene

import log "github.com/rs/zerolog/log"

type MoveCtrl struct {
	owner *SceneEntity // 拥有者
}

func NewMoveCtrl(owner *SceneEntity) *MoveCtrl {
	c := &MoveCtrl{
		owner: owner,
	}

	return c
}

func (c *MoveCtrl) Update() {
	log.Info().Int64("owner_id", c.owner.id).Msg("MoveCtrl update")
}
