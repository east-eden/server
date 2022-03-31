package collection

import (
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/quest"
	"github.com/shopspring/decimal"
)

// collection create pool
var collectionPool = &sync.Pool{New: newPoolCollection}

func GetCollectionPool() *sync.Pool {
	return collectionPool
}

func NewCollection() *Collection {
	return collectionPool.Get().(*Collection)
}

type Collection struct {
	Options             `bson:"inline" json:",inline"`
	*quest.QuestManager `bson:"inline" json:",inline"`
	score               int32 `bson:"-" json:"-"`
}

func newPoolCollection() any {
	return &Collection{}
}

func (c *Collection) Init(opts ...Option) {
	c.Options = DefaultOptions()

	for _, o := range opts {
		o(c.GetOptions())
	}

	c.QuestManager = quest.NewQuestManager()
	c.CalcScore()
}

func (c *Collection) InitQuestManager() {
	questList := make([]int32, 0, 1)
	if c.Entry.QuestId != -1 {
		questList = append(questList, c.Entry.QuestId)
	}

	c.QuestManager.Init(
		quest.WithManagerOwnerId(c.Id),
		quest.WithManagerOwnerType(define.QuestOwner_Type_Collection),
		quest.WithManagerStoreType(define.StoreType_Collection),
		quest.WithManagerAdditionalQuestId(questList...),
		quest.WithManagerEventManager(c.eventManager),
		quest.WithManagerQuestChangedCb(func(q *quest.Quest) {
			c.questUpdateCb(q)
		}),
	)
}

func (c *Collection) GetOptions() *Options {
	return &c.Options
}

func (c *Collection) GenCollectionPB() *pbGlobal.Collection {
	pb := &pbGlobal.Collection{
		TypeId: c.TypeId,
		Active: c.Active,
		Star:   int32(c.Star),
		BoxId:  c.BoxId,
		Score:  c.score,
	}

	return pb
}

func (c *Collection) CalcScore() {
	globalConfig, _ := auto.GetGlobalConfig()

	if !c.Active {
		c.score = 0
		return
	}

	// 基础分 * 品质系数 * 星级系数 * 觉醒系数
	baseScore := decimal.NewFromInt32(c.Entry.BaseScore)
	qualityFactor := globalConfig.CollectionQualityScoreFactor[c.Entry.Quality]
	starFactor := globalConfig.CollectionStarScoreFactor[c.Star]
	wakeupFactor := decimal.NewFromInt32(1)
	if c.Wakeup {
		wakeupFactor = globalConfig.CollectionWakeupScoreFactor
	}

	c.score = int32(qualityFactor.Mul(starFactor).Mul(wakeupFactor).Mul(baseScore).Round(0).IntPart())
}
