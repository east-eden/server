package hero

import (
	"sync"

	"bitbucket.org/funplus/server/define"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/item"
	"bitbucket.org/funplus/server/services/game/talent"
)

// hero create pool
var heroPool = &sync.Pool{New: newPoolHero}

func GetHeroPool() *sync.Pool {
	return heroPool
}

func NewHero() *Hero {
	return heroPool.Get().(*Hero)
}

type Hero struct {
	Options    `bson:"inline" json:",inline"`
	equipBar   *item.EquipBar    `bson:"-" json:"-"`
	attManager *HeroAttManager   `bson:"-" json:"-"`
	crystalBox *item.CrystalBox  `bson:"-" json:"-"`
	TalentBox  *talent.TalentBox `bson:"inline" json:",inline"`
}

func newPoolHero() interface{} {
	return &Hero{}
}

func (h *Hero) Init(opts ...Option) {
	h.Options = DefaultOptions()

	for _, o := range opts {
		o(h.GetOptions())
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = NewHeroAttManager(h)
	h.crystalBox = item.NewCrystalBox(h)
	h.TalentBox = talent.NewTalentBox(h, nil, define.Talent_Type_Hero)
}

func (h *Hero) GetOptions() *Options {
	return &h.Options
}

func (h *Hero) GetType() int32 {
	return define.Plugin_Hero
}

func (h *Hero) GetTypeId() int32 {
	return h.Entry.Id
}

func (h *Hero) GetStoreType() int {
	return define.StoreType_Hero
}

func (h *Hero) GetId() int64 {
	return h.Options.Id
}

func (h *Hero) GetLevel() int32 {
	return int32(h.Options.Level)
}

func (h *Hero) GetAttManager() *HeroAttManager {
	return h.attManager
}

func (h *Hero) GetEquipBar() *item.EquipBar {
	return h.equipBar
}

func (h *Hero) GetCrystalBox() *item.CrystalBox {
	return h.crystalBox
}

func (h *Hero) GetTalentBox() *talent.TalentBox {
	return h.TalentBox
}

func (h *Hero) AddExp(exp int32) int32 {
	h.Exp += exp
	return h.Exp
}

func (h *Hero) AddLevel(level int16) int16 {
	h.Level += level
	return h.Level
}

func (h *Hero) GenHeroPB() *pbGlobal.Hero {
	pb := &pbGlobal.Hero{
		Id:            h.Id,
		TypeId:        h.TypeId,
		Exp:           h.Exp,
		Level:         int32(h.Level),
		PromoteLevel:  int32(h.PromoteLevel),
		Star:          int32(h.Star),
		Friendship:    h.Friendship,
		FashionId:     h.FashionId,
		CrystalSkills: h.crystalBox.GetSkills(),
		TalentList:    h.GetTalentBox().GenTalentList(),
		AttValues:     h.attManager.ExportInt32(),
	}

	return pb
}

func (h *Hero) GenEntityInfoPB() *pbGlobal.EntityInfo {
	h.attManager.CalcAtt()

	pb := &pbGlobal.EntityInfo{
		HeroTypeId:    h.TypeId,
		CrystalSkills: h.crystalBox.GetSkills(),
		AttValue:      h.attManager.ExportInt32(),
	}

	return pb
}
