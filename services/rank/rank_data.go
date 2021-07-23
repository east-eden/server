package rank

import (
	"context"
	"errors"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/hellodudu/task"
	"github.com/liyiheng/zset"
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
	Id             int32                     `json:"_id" bson:"_id"`
	LastSaveNodeId int32                     `json:"last_save_node_id" bson:"last_save_node_id"`
	Raws           map[int64]*define.RankRaw `json:"raws" bson:"raws"`
	NodeId         int16                     `json:"-" bson:"-"` // 当前节点id
	ZSets          *zset.SortedSet           `json:"-" bson:"-"` // 排行zset
	tasker         *task.Tasker              `json:"-" bson:"-"`
	rpcHandler     *RpcHandler               `json:"-" bson:"-"`
	entry          *auto.RankEntry           `json:"-" bson:"-"`
}

func NewRankData() interface{} {
	return &RankData{}
}

func (r *RankData) Init(nodeId int16, rpcHandler *RpcHandler) {
	r.Id = -1
	r.LastSaveNodeId = -1
	r.Raws = make(map[int64]*define.RankRaw)
	r.NodeId = nodeId
	r.ZSets = zset.New()
	r.rpcHandler = rpcHandler
}

func (r *RankData) InitTask() {
	r.tasker = task.NewTasker()
	r.tasker.Init(
		task.WithStopFns(r.onTaskStop),
		task.WithUpdateFn(r.onTaskUpdate),
		task.WithTimeout(RankDataTaskTimeout),
		task.WithSleep(time.Second),
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

	var storeType int
	if r.entry.Local {
		storeType = define.StoreType_LocalRank
	} else {
		storeType = define.StoreType_GlobalRank
	}

	// 加载排行榜信息
	err := store.GetStore().FindOne(context.Background(), storeType, rankId, r)

	// 创建新排行榜数据
	if errors.Is(err, store.ErrNoResult) {
		r.Id = rankId
		r.LastSaveNodeId = int32(r.NodeId)
		errSave := store.GetStore().UpdateOne(context.Background(), storeType, rankId, r, true)
		utils.ErrPrint(errSave, "UpdateOne failed when RankData.Load", rankId)
		return errSave
	}

	if !utils.ErrCheck(err, "FindOne failed when RankData.Load", rankId) {
		return err
	}

	// 数据排序
	for key, value := range r.Raws {
		r.ZSets.Set(value.Score, key, value)
	}

	return nil
}

func (r *RankData) onTaskStop() {
	log.Info().Caller().Int32("rank_id", r.Id).Msg("RankData task stopped...")
}

func (r *RankData) onTaskUpdate() {
}

func (r *RankData) TaskRun(ctx context.Context) error {
	return r.tasker.Run(ctx)
}

func (r *RankData) Stop() {
	r.tasker.Stop()
}

func (r *RankData) AddTask(ctx context.Context, fn task.TaskHandler, p ...interface{}) error {
	return r.tasker.AddWait(ctx, fn, p...)
}

func (r *RankData) SetScore(ctx context.Context, rankRaw *define.RankRaw) error {
	if rankRaw == nil {
		return ErrInvalidRankRaw
	}

	r.ZSets.Set(rankRaw.Score, rankRaw.ObjId, rankRaw)

	// todo save to mongodb
	return nil
}

func (r *RankData) GetRankByKey(ctx context.Context, key int64) (*define.RankRaw, error) {
	data, ok := r.ZSets.GetData(key)
	if !ok {
		return nil, ErrRankNotExist
	}

	return data.(*define.RankRaw), nil
}

func (r *RankData) GetRankByIndex(ctx context.Context, start, end int64) ([]*define.RankRaw, error) {
	res := make([]*define.RankRaw, 0, 64)
	if r.entry.Desc {
		r.ZSets.RevRange(start, end, func(score float64, key int64, data interface{}) {
			res = append(res, data.(*define.RankRaw))
		})
	} else {
		r.ZSets.Range(start, end, func(score float64, key int64, data interface{}) {
			res = append(res, data.(*define.RankRaw))
		})
	}
	return res, nil
}
