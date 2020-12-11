package blade

import (
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/services/game/talent"
	"github.com/east-eden/server/utils"
)

type BladeV1 struct {
	Options                `bson:"inline" json:",inline"`
	TalentManager          *talent.TalentManager `bson:"-" json:"-"`
	AttManager             *att.AttManager       `bson:"-" json:"-"`
	utils.WaitGroupWrapper `bson:"-" json:"-"`
}

func newPoolBladeV1() interface{} {
	b := &BladeV1{
		Options: DefaultOptions(),
	}

	b.AttManager = att.NewAttManager(-1)

	return b
}

func (b *BladeV1) GetOptions() *Options {
	return &b.Options
}

func (b *BladeV1) GetObjID() int64 {
	return b.Options.Id
}

func (b *BladeV1) GetStoreIndex() int64 {
	return b.Options.OwnerId
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
	if b.TalentManager == nil {
		return nil
	}

	// load blade's talent
	var errLoad error = nil
	b.Wrap(func() {
		if err := b.TalentManager.LoadFromDB(); err != nil {
			errLoad = err
		}
	})

	b.Wait()

	if errLoad != nil {
		return fmt.Errorf("BladeV1 LoadFromDb: %w", errLoad)
	}

	return nil
}

func (b *BladeV1) AddExp(exp int64) int64 {
	b.Options.Exp += exp
	return b.Options.Exp
}

func (b *BladeV1) AddLevel(level int32) int32 {
	b.Options.Level += level
	return b.Options.Level
}

func (b *BladeV1) CalcAtt() {

}

func (b *BladeV1) AfterLoad() error {
	return nil
}
