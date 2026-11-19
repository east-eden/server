package alliance

import (
	"context"
	"errors"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/store"
)

// MongoStore 基于 store 层实现公会持久化。
//
// 公会本身用 alliance 表存（_id = 公会id），
// 名字索引单独存 alliance_name 表（_id = 名字），用于重名校验。
type MongoStore struct {
	table     string
	nameTable string
}

func NewMongoStore() *MongoStore {
	return &MongoStore{table: "alliance", nameTable: "alliance_name"}
}

// Init 注册表并建索引，服务启动时调用一次
func (s *MongoStore) Init() {
	store.GetStore().AddStoreInfo(define.StoreType_Alliance, s.table, "_id")
	store.GetStore().AddStoreInfo(define.StoreType_AllianceName, s.nameTable, "_id")

	_ = store.GetStore().MigrateDbTable(s.table)
	_ = store.GetStore().MigrateDbTable(s.nameTable)
}

func (s *MongoStore) Load(ctx context.Context, id int64) (*Guild, error) {
	g := &Guild{}
	err := store.GetStore().FindOne(ctx, define.StoreType_Alliance, id, g)
	if errors.Is(err, store.ErrNoResult) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if g.Members == nil {
		g.Members = make(map[int64]int32, 8)
	}

	return g, nil
}

func (s *MongoStore) Save(ctx context.Context, g *Guild) error {
	return store.GetStore().UpdateOne(ctx, define.StoreType_Alliance, g.Id, g)
}

func (s *MongoStore) Delete(ctx context.Context, id int64) error {
	return store.GetStore().DeleteOne(ctx, define.StoreType_Alliance, id)
}

type nameDoc struct {
	Id  string `bson:"_id"`
	Gid int64  `bson:"gid"`
}

func (s *MongoStore) FindByName(ctx context.Context, name string) (int64, bool, error) {
	doc := &nameDoc{}
	err := store.GetStore().FindOne(ctx, define.StoreType_AllianceName, name, doc)
	if errors.Is(err, store.ErrNoResult) {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, err
	}

	return doc.Gid, true, nil
}

func (s *MongoStore) SaveName(ctx context.Context, name string, id int64) error {
	return store.GetStore().UpdateOne(ctx, define.StoreType_AllianceName, name, &nameDoc{Id: name, Gid: id})
}
