package blade

import (
	"fmt"
	"sync"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/internal/att"
	"bitbucket.org/east-eden/server/services/game/talent"
	"bitbucket.org/east-eden/server/utils"
)

// blade create pool
var bladePool = &sync.Pool{New: newPoolBlade}

func NewPoolBlade() *Blade {
	return bladePool.Get().(*Blade)
}

func GetBladePool() *sync.Pool {
	return bladePool
}

func ReleasePoolBlade(x interface{}) {
	bladePool.Put(x)
}

func NewBlade(opts ...Option) *Blade {
	b := NewPoolBlade()

	for _, o := range opts {
		o(b.GetOptions())
	}

	return b
}

type Blade struct {
	Options                `bson:"inline" json:",inline"`
	TalentManager          *talent.TalentManager `bson:"-" json:"-"`
	AttManager             *att.AttManager       `bson:"-" json:"-"`
	utils.WaitGroupWrapper `bson:"-" json:"-"`
}

func newPoolBlade() interface{} {
	b := &Blade{
		Options: DefaultOptions(),
	}

	b.AttManager = att.NewAttManager(-1)

	return b
}

func (b *Blade) GetOptions() *Options {
	return &b.Options
}

func (b *Blade) GetObjID() int64 {
	return b.Options.Id
}

func (b *Blade) GetStoreIndex() int64 {
	return b.Options.OwnerId
}

func (b *Blade) GetType() int32 {
	return define.Plugin_Blade
}

func (b *Blade) GetID() int64 {
	return b.Options.Id
}

func (b *Blade) GetLevel() int32 {
	return b.Options.Level
}

func (b *Blade) GetAttManager() *att.AttManager {
	return b.AttManager
}

func (b *Blade) SetTalentManager(t *talent.TalentManager) {
	b.TalentManager = t
}

func (b *Blade) GetTalentManager() *talent.TalentManager {
	return b.TalentManager
}

func (b *Blade) LoadFromDB() error {
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

func (b *Blade) AddExp(exp int64) int64 {
	b.Options.Exp += exp
	return b.Options.Exp
}

func (b *Blade) AddLevel(level int32) int32 {
	b.Options.Level += level
	return b.Options.Level
}

func (b *Blade) CalcAtt() {

}

func (b *Blade) AfterLoad() error {
	return nil
}
