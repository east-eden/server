package collection

import (
	"sync"

	"bitbucket.org/funplus/server/define"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/quest"
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
}

func newPoolCollection() interface{} {
	return &Collection{}
}

func (c *Collection) Init(opts ...Option) {
	c.Options = DefaultOptions()

	for _, o := range opts {
		o(c.GetOptions())
	}

	c.QuestManager = quest.NewQuestManager()
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
	}

	return pb
}
