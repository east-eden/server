package player

import (
	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/talent"
	"github.com/yokaiio/yokai_server/game/token"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type DefaultPlayer struct {
	ds *db.Datastore
	wg utils.WaitGroupWrapper

	itemManager   *item.ItemManager
	heroManager   *hero.HeroManager
	tokenManager  *token.TokenManager
	talentManager *talent.TalentManager

	ID       int64  `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null"`
	ClientID int64  `gorm:"type:bigint(20);column:client_id;default:-1;not null"`
	Name     string `gorm:"type:varchar(32);column:name;not null"`
	Exp      int64  `gorm:"type:bigint(20);column:exp;default:0;not null"`
	Level    int32  `gorm:"type:int(10);column:level;default:1;not null"`
}

func newDefaultPlayer(id int64, name string, ds *db.Datastore) Player {
	return &DefaultPlayer{
		ds:            ds,
		ID:            id,
		ClientID:      -1,
		Name:          name,
		Exp:           0,
		Level:         1,
		itemManager:   item.NewItemManager(id, ds),
		heroManager:   hero.NewHeroManager(id, ds),
		tokenManager:  token.NewTokenManager(id, ds),
		talentManager: talent.NewTalentManager(id, ds),
	}
}

func defaultMigrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(DefaultPlayer{})
	item.Migrate(ds)
	hero.Migrate(ds)
	token.Migrate(ds)
	talent.Migrate(ds)
}

func (p *DefaultPlayer) TableName() string {
	return "player"
}

func (p *DefaultPlayer) GetID() int64 {
	return p.ID
}

func (p *DefaultPlayer) GetClientID() int64 {
	return p.ClientID
}

func (p *DefaultPlayer) GetName() string {
	return p.Name
}

func (p *DefaultPlayer) GetExp() int64 {
	return p.Exp
}

func (p *DefaultPlayer) GetLevel() int32 {
	return p.Level
}

func (p *DefaultPlayer) SetClientID(id int64) {
	p.ClientID = id
}

func (p *DefaultPlayer) SetName(name string) {
	p.Name = name
}

func (p *DefaultPlayer) SetExp(exp int64) {
	p.Exp = exp
}

func (p *DefaultPlayer) SetLevel(level int32) {
	p.Level = level
}

func (p *DefaultPlayer) HeroManager() *hero.HeroManager {
	return p.heroManager
}

func (p *DefaultPlayer) ItemManager() *item.ItemManager {
	return p.itemManager
}

func (p *DefaultPlayer) TokenManager() *token.TokenManager {
	return p.tokenManager
}

func (p *DefaultPlayer) TalentManager() *talent.TalentManager {
	return p.talentManager
}

func (p *DefaultPlayer) LoadFromDB() {
	p.wg.Wrap(p.heroManager.LoadFromDB)
	p.wg.Wrap(p.itemManager.LoadFromDB)
	p.wg.Wrap(p.tokenManager.LoadFromDB)
	p.wg.Wrap(p.talentManager.LoadFromDB)
	p.wg.Wait()
}

func (p *DefaultPlayer) AfterLoad() {
	items := p.itemManager.GetItemList()
	for _, v := range items {
		if v.GetEquipObj() == -1 {
			continue
		}

		if err := p.HeroManager().PutonEquip(v.GetEquipObj(), v.GetID(), v.Entry().EquipPos); err != nil {
			logger.Warn("Hero puton equip error when loading from db:", err)
		}
	}
}

func (p *DefaultPlayer) Save() {
	p.ds.ORM().Save(p)
}

func (p *DefaultPlayer) ChangeExp(add int64) {
	p.Exp += add

	p.ds.ORM().Model(p).Updates(DefaultPlayer{
		Exp: p.Exp,
	})
}

func (p *DefaultPlayer) ChangeLevel(add int32) {
	p.Level += add

	p.ds.ORM().Model(p).Updates(DefaultPlayer{
		Level: p.Level,
	})
}
