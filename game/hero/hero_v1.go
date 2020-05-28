package hero

import (
	"context"
	"log"
	"time"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/rune"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type HeroV1 struct {
	Opts       *Options        `bson:"inline"`
	equipBar   *item.EquipBar  `gorm:"-" bson:"-"`
	attManager *att.AttManager `gorm:"-" bson:"-"`
	runeBox    *rune.RuneBox   `gorm:"-" bson:"-"`
}

func newPoolHeroV1() interface{} {
	h := &HeroV1{
		Opts: DefaultOptions(),
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = att.NewAttManager(-1)
	h.runeBox = rune.NewRuneBox(h)

	return h
}

func migrateV1(ds *db.Datastore) {
	coll := ds.Database().Collection("hero")

	// check index
	idx := coll.Indexes()

	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := idx.List(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	indexExist := false
	for cursor.Next(context.Background()) {
		var result bson.M
		cursor.Decode(&result)
		if result["name"] == "owner_id" {
			indexExist = true
			break
		}
	}

	// create index
	if !indexExist {
		_, err := coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{"owner_id", bsonx.Int32(1)}},
				Options: options.Index().SetName("owner_id"),
			},
		)

		if err != nil {
			logger.Warn("collection hero create index owner_id failed:", err)
		}
	}

}

func (h *HeroV1) Options() *Options {
	return h.Opts
}

func (h *HeroV1) GetType() int32 {
	return define.Plugin_Hero
}

func (h *HeroV1) GetID() int64 {
	return h.Opts.Id
}

func (h *HeroV1) GetLevel() int32 {
	return h.Opts.Level
}

func (h *HeroV1) GetAttManager() *att.AttManager {
	return h.attManager
}

func (h *HeroV1) GetEquipBar() *item.EquipBar {
	return h.equipBar
}

func (h *HeroV1) GetRuneBox() *rune.RuneBox {
	return h.runeBox
}

func (h *HeroV1) AddExp(exp int64) int64 {
	h.Opts.Exp += exp
	return h.Opts.Exp
}

func (h *HeroV1) AddLevel(level int32) int32 {
	h.Opts.Level += level
	return h.Opts.Level
}

func (h *HeroV1) BeforeDelete() {

}

func (h *HeroV1) CalcAtt() {
	h.attManager.Reset()

	// equip bar
	var n int32
	for n = 0; n < define.Hero_MaxEquip; n++ {
		i := h.equipBar.GetEquipByPos(n)
		if i == nil {
			continue
		}

		i.CalcAtt()
		h.attManager.ModAttManager(i.GetAttManager())
	}

	// rune box
	for n = 0; n < define.Rune_PositionEnd; n++ {
		r := h.runeBox.GetRuneByPos(n)
		if r == nil {
			continue
		}

		r.CalcAtt()
		h.attManager.ModAttManager(r.GetAttManager())
	}

	h.attManager.CalcAtt()
}
