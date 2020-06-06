package blade

import (
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/talent"
	"github.com/yokaiio/yokai_server/utils"
)

type BladeV1 struct {
	Opts                   *Options              `bson:"inline" redis:"inline"`
	talentManager          *talent.TalentManager `bson:"-" redis:"-"`
	utils.WaitGroupWrapper `bson:"-" redis:"-"`
}

func newPoolBladeV1() interface{} {
	b := &BladeV1{
		Opts: DefaultOptions(),
	}

	b.talentManager = talent.NewTalentManager(b)
	return b
}

func (b *BladeV1) GetType() int32 {
	return define.Plugin_Blade
}

func (b *BladeV1) GetID() int64 {
	return b.Opts.Id
}

func (b *BladeV1) GetLevel() int32 {
	return b.Opts.Level
}

func (b *BladeV1) LoadFromDB() {
	b.Wrap(b.talentManager.LoadFromDB)
	b.Wait()
}

func (b *BladeV1) TalentManager() *talent.TalentManager {
	return b.talentManager
}
