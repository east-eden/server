package costloot

import (
	"testing"

	"github.com/yokaiio/yokai_server/internal/define"
)

func init() {

}

type PlayerObj struct {
}

func (o *PlayerObj) GetType() int32 {
	return define.Plugin_Player
}

func (o *PlayerObj) GetID() int64 {
	return 1
}

func (o *PlayerObj) GetLevel() int32 {
	return 1
}

func TestCostLoot(t *testing.T) {

}
