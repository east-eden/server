package rank

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"e.coding.net/mmstudio/blade/server/utils/zset"
	"github.com/hellodudu/task"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidRank       = errors.New("invalid rank")
	ErrInvalidRankRaw    = errors.New("invalid rank raw")
	ErrInvalidRankStatus = errors.New("invalid rank status")
	ErrRankNotExist      = errors.New("rank not exist")
	ErrAddExistRank      = errors.New("add exist rank")

	RankDataTaskTimeout          = time.Hour       // 邮箱任务超时
	RankDataChannelResultTimeout = 5 * time.Second // 邮箱channel处理超时
)

// 排行榜数据
type RankData struct {
	RankId         int32           `json:"_id" bson:"_id"`
	LastSaveNodeId int32           `json:"last_save_node_id" bson:"last_save_node_id"`
	NodeId         int16           `json:"-" bson:"-"` // 当前节点id
	zsets          *zset.SortedSet `json:"-" bson:"-"` // 排行zset
	tasker         *task.Tasker    `json:"-" bson:"-"`
	rpcHandler     *RpcHandler     `json:"-" bson:"-"`
	entry          *auto.RankEntry `json:"-" bson:"-"`
}

func NewRankData() interface{} {
	return &RankData{}
}

func (r *RankData) Init(nodeId int16, rpcHandler *RpcHandler) {
	r.RankId = -1
	r.LastSaveNodeId = -1
	r.NodeId = nodeId
	r.zsets = zset.New()
	r.rpcHandler = rpcHandler
}

func (r *RankData) InitTask() {
	r.tasker = task.NewTasker()
	r.tasker.Init(
		task.WithStopFns(r.onTaskStop),
		// task.WithUpdateFn(r.onTaskUpdate),
		task.WithTimeout(RankDataTaskTimeout),
		task.WithSleep(time.Millisecond),
	)
}

func (r *RankData) IsTaskRunning() bool {
	return r.tasker.IsRunning()
}

func (r *RankData) Load(rankId int32) error {
	var ok bool
	r.entry, ok = auto.GetRankEntry(rankId)
	if !ok {
		return ErrInvalidRank
	}

	// 加载排行榜信息
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Rank, rankId, r)

	// 创建新排行榜数据
	if errors.Is(err, store.ErrNoResult) {
		r.RankId = rankId
		r.LastSaveNodeId = int32(r.NodeId)
		errSave := store.GetStore().UpdateOne(context.Background(), define.StoreType_Rank, rankId, r, true)
		utils.ErrPrint(errSave, "UpdateOne failed when RankData.Load", rankId)
		return errSave
	}

	if !utils.ErrCheck(err, "FindOne failed when RankData.Load", rankId) {
		return err
	}

	// 加载排行榜数据
	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Rank, "_id.rank_id", rankId)
	if !utils.ErrCheck(err, "FindAll failed when RankData.Load", rankId) {
		return err
	}

	for _, v := range res {
		vv := v.([]byte)
		raw := &define.RankRaw{}
		err := json.Unmarshal(vv, raw)
		if !utils.ErrCheck(err, "json.Unmarshal failed when RankData.Load", vv) {
			continue
		}

		r.zsets.Set(raw.Score, raw.ObjId, raw.Date, raw)
	}

	return nil
}

func (r *RankData) onTaskStop() {
	log.Info().Caller().Int32("rank_id", r.RankId).Msg("RankData task stopped...")
}

// func (r *RankData) onTaskUpdate() {
// }

func (r *RankData) TaskRun(ctx context.Context) error {
	return r.tasker.Run(ctx)
}

func (r *RankData) Stop() {
	r.tasker.Stop()
}

func (r *RankData) saveLastNode() {
	fields := map[string]interface{}{
		"last_save_node_id": r.NodeId,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Rank, r.RankId, fields, true)
	_ = utils.ErrCheck(err, "UpdateFields failed when RankData.saveLastNode", r.RankId)
}

func (r *RankData) AddTask(ctx context.Context, fn task.TaskHandler, p ...interface{}) error {
	return r.tasker.AddWait(ctx, fn, p...)
}

func (r *RankData) SetScore(ctx context.Context, rankRaw *define.RankRaw) error {
	if rankRaw == nil {
		return ErrInvalidRankRaw
	}

	rr := &define.RankRaw{}
	*rr = *rankRaw

	if r.entry.Desc {
		rr.Score *= -1
	}

	r.zsets.Set(rr.Score, rr.ObjId, rr.Date, rr)

	// save rank raw
	err := store.GetStore().UpdateOne(ctx, define.StoreType_Rank, rr.RankKey, rr)
	_ = utils.ErrCheck(err, "UpdateOne failed when RankData.SetScore", rr)

	r.saveLastNode()
	return err
}

func (r *RankData) GetRankByObjId(ctx context.Context, objId int64) (rank int64, raw define.RankRaw, err error) {
	zRank, _, data := r.zsets.GetRank(objId, false)
	rank = zRank
	if data == nil {
		err = ErrRankNotExist
		return
	}

	raw = *data.(*define.RankRaw)
	if r.entry.Desc {
		raw.Score *= -1
	}

	return rank, raw, nil
}

func (r *RankData) GetRankByRange(ctx context.Context, start, end int64) (raws []define.RankRaw, err error) {
	r.zsets.Range(start, end, func(score float64, key int64, data interface{}) {
		rr := *data.(*define.RankRaw)
		if r.entry.Desc {
			rr.Score *= -1
		}
		raws = append(raws, rr)
	})
	return
}
