package scene

import log "github.com/rs/zerolog/log"

type AtbCtrl struct {
	*AtbOption
	owner *SceneEntity // 拥有者
}

type AtbOption func(*AtbOptions)
type AtbOptions struct {
}

func DefaultAtbOptions() *AtbOptions {
	o := &AtbOptions{}
	return o
}

func NewAtbCtrl(owner *SceneEntity, opts ...AtbOption) *AtbCtrl {
	c := &AtbCtrl{
		owner: owner,
	}

	return c
}

func (c *AtbCtrl) Update() {
	log.Info().Int64("owner_id", c.owner.id).Msg("AtbCtrl update")
}
