// Package settle 记录玩家的结算数据。
//
// 与 kof 的 SettleManager 语义一致：按玩家 + 结算类型保存一份带过期时间的
// 序列化结算数据，供跨服/异步结算流程读写。kof 用的是 Redis hash + 30 天过期，
// 这里落在 Mongo 的 player_settle 表并显式记录过期时间，读取方按 ExpireAt 判断有效性。
package settle

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
)

var (
	once         sync.Once
	ErrNilPlayer = errors.New("nil player")
)

// SettleRecord 单条结算记录
type SettleRecord struct {
	SettleData []byte `bson:"settle_data" json:"settle_data"`
	ExpireAt   int64  `bson:"expire_at" json:"expire_at"`
}

// playerSettleDoc 一个玩家的全部结算记录
type playerSettleDoc struct {
	Id      int64                    `bson:"_id" json:"_id"`
	Records map[string]*SettleRecord `bson:"records" json:"records"`
}

// InitStore 注册结算表，服务启动时调用
func InitStore() {
	once.Do(func() {
		store.GetStore().AddStoreInfo(define.StoreType_Settle, "player_settle", "_id")
		utils.ErrPrint(store.GetStore().MigrateDbTable("player_settle"), "migrate player_settle failed")
	})
}

// RecordPlayerSettle 记录玩家结算数据
func RecordPlayerSettle(playerId int64, tp config.SettleType, settleData []byte, expireTime time.Duration) error {
	if playerId <= 0 {
		return ErrNilPlayer
	}

	InitStore()

	ctx := context.Background()
	doc := &playerSettleDoc{}

	err := store.GetStore().FindOne(ctx, define.StoreType_Settle, playerId, doc)
	if err != nil && !errors.Is(err, store.ErrNoResult) {
		return err
	}

	if doc.Records == nil {
		doc.Records = make(map[string]*SettleRecord, 1)
	}

	doc.Id = playerId
	doc.Records[tp.String()] = &SettleRecord{
		SettleData: settleData,
		ExpireAt:   utils.NowInLocation().Add(expireTime).UnixMilli(),
	}

	log.Info().
		Int64("player_id", playerId).
		Str("settle_type", tp.String()).
		Int64("expire_at", doc.Records[tp.String()].ExpireAt).
		Msg("record player settle")

	return store.GetStore().UpdateOne(ctx, define.StoreType_Settle, playerId, doc)
}

// GetPlayerSettle 读取玩家结算数据，过期视为不存在
func GetPlayerSettle(playerId int64, tp config.SettleType) ([]byte, bool, error) {
	if playerId <= 0 {
		return nil, false, ErrNilPlayer
	}

	InitStore()

	doc := &playerSettleDoc{}
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Settle, playerId, doc)
	if errors.Is(err, store.ErrNoResult) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	record, ok := doc.Records[tp.String()]
	if !ok {
		return nil, false, nil
	}

	if record.ExpireAt > 0 && record.ExpireAt < utils.NowInLocation().UnixMilli() {
		return nil, false, nil
	}

	return record.SettleData, true, nil
}

// DeletePlayerSettle 删除指定类型的结算数据
func DeletePlayerSettle(playerId int64, tp config.SettleType) error {
	if playerId <= 0 {
		return ErrNilPlayer
	}

	InitStore()

	return store.GetStore().DeleteFields(
		context.Background(),
		define.StoreType_Settle,
		playerId,
		[]string{"records." + tp.String()},
	)
}
