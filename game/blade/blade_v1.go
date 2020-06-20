package blade

import (
	"fmt"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/talent"
	"github.com/yokaiio/yokai_server/utils"
)

type BladeV1 struct {
	Options                `bson:"inline" json:",inline"`
	talentManager          *talent.TalentManager `bson:"-" json:"-"`
	utils.WaitGroupWrapper `bson:"-" json:"-"`
}

func newPoolBladeV1() interface{} {
	b := &BladeV1{
		Options: DefaultOptions(),
	}

	return b
}

func (b *BladeV1) GetOptions() *Options {
	return &b.Options
}

func (b *BladeV1) GetType() int32 {
	return define.Plugin_Blade
}

func (b *BladeV1) GetID() int64 {
	return b.Options.Id
}

func (b *BladeV1) GetLevel() int32 {
	return b.Options.Level
}

func (b *BladeV1) LoadFromDB() error {
	if b.talentManager == nil {
		return nil
	}

	// load blade's talent
	var errLoad error = nil
	b.Wrap(func() {
		if err := b.talentManager.LoadFromDB(); err != nil {
			errLoad = err
		}
	})

	b.Wait()

	if errLoad != nil {
		return fmt.Errorf("BladeV1 LoadFromDb: %w", errLoad)
	}

	return nil
}

func (b *BladeV1) TalentManager() *talent.TalentManager {
	return b.talentManager
}

func (b *BladeV1) SetTalentManager(m *talent.TalentManager) {
	b.talentManager = m
}
