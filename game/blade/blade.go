package blade

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/talent"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Blade struct {
	ID        int64              `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null"`
	OwnerID   int64              `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null"`
	OwnerType int32              `gorm:"type:int(10);column:owner_type;index:owner_type;default:-1;not null"`
	TypeID    int32              `gorm:"type:int(10);column:type_id;default:-1;not null"`
	Exp       int64              `gorm:"type:bigint(20);column:exp;default:0;not null"`
	Level     int32              `gorm:"type:int(10);column:level;default:1;not null"`
	Entry     *define.BladeEntry `gorm:"-"`

	talentManager *talent.TalentManager
	wg            utils.WaitGroupWrapper
}

func newBlade(id int64, owner define.PluginObj, ds *db.Datastore) *Blade {
	b := &Blade{
		ID:        id,
		OwnerID:   owner.GetID(),
		OwnerType: owner.GetType(),
		TypeID:    -1,
		Exp:       0,
		Level:     1,
	}

	b.talentManager = talent.NewTalentManager(b, ds)
	return b
}

func (b *Blade) GetType() int32 {
	return define.Plugin_Blade
}

func (b *Blade) GetID() int64 {
	return b.ID
}

func (b *Blade) GetLevel() int32 {
	return b.Level
}

func (b *Blade) LoadFromDB() {
	b.wg.Wrap(b.talentManager.LoadFromDB)
	b.wg.Wait()
}

func (b *Blade) TalentManager() *talent.TalentManager {
	return b.talentManager
}

func (b *Blade) TableName() string {
	return "blade"
}
