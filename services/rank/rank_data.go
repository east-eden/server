package rank

import (
	"context"
	"errors"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/hellodudu/task"
	"github.com/liyiheng/zset"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidRank       = errors.New("invalid rank")
	ErrInvalidRankStatus = errors.New("invalid rank status")
	ErrAddExistRank      = errors.New("add exist rank")

	RankDataTaskTimeout          = time.Hour       // 邮箱任务超时
	RankDataChannelResultTimeout = 5 * time.Second // 邮箱channel处理超时
)

// 排行榜数据
type RankData struct {
	Id             int64           `json:"_id" bson:"_id"`
	LastSaveNodeId int32           `json:"last_save_node_id" bson:"last_save_node_id"`
	NodeId         int16           `json:"-" bson:"-"` // 当前节点id
	ZSets          *zset.SortedSet `json:"-" bson:"-"` // 排行zset
	tasker         *task.Tasker    `json:"-" bson:"-"`
	rpcHandler     *RpcHandler     `json:"-" bson:"-"`
}

func NewRankData() interface{} {
	return &RankData{}
}

func (r *RankData) Init(nodeId int16, rpcHandler *RpcHandler) {
	r.Id = -1
	r.LastSaveNodeId = -1
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

func (r *RankData) ResetTaskTimeout() {
	r.tasker.ResetTimer()
}

func (r *RankData) IsTaskRunning() bool {
	return r.tasker.IsRunning()
}

func (r *RankData) Load(ownerId int64) error {
	// 加载排行榜信息
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Mail, ownerId, r)

	// 创建新排行榜数据
	if errors.Is(err, store.ErrNoResult) {
		r.Id = ownerId
		r.LastSaveNodeId = int32(r.NodeId)
		errSave := store.GetStore().UpdateOne(context.Background(), define.StoreType_Mail, ownerId, r, true)
		utils.ErrPrint(errSave, "UpdateOne failed when MailBox.Load", ownerId)
		return errSave
	}

	if !utils.ErrCheck(err, "FindOne failed when MailBox.Load", ownerId) {
		return err
	}

	// 加载所有邮件
	_, errMails := store.GetStore().FindAll(context.Background(), define.StoreType_Mail, "owner_id", ownerId)
	if !utils.ErrCheck(errMails, "FindAll failed when MailBox.Load", ownerId) {
		return errMails
	}

	// for _, v := range res {
	// 	vv := v.([]byte)
	// 	mail := &define.Mail{}
	// 	err := json.Unmarshal(vv, mail)
	// 	if !utils.ErrCheck(err, "json.Unmarshal failed when MailBox.Load", ownerId) {
	// 		continue
	// 	}

	// 	r.Mails[mail.Id] = mail
	// }

	return nil
}

func (r *RankData) onTaskStop() {
	log.Info().Caller().Int64("owner_id", r.Id).Msg("mailbox task stopped...")
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
