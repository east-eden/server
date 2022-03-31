package comment

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/zset"
	"github.com/hellodudu/task"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidComment         = errors.New("invalid comment")
	ErrInvalidCommentMetadata = errors.New("invalid comment metadata")
	ErrInvalidCommentStatus   = errors.New("invalid comment status")
	ErrCommentNotExist        = errors.New("comment not exist")
	ErrAddExistComment        = errors.New("add exist comment")

	CommentDataTaskTimeout          = time.Hour       // 评论任务超时
	CommentDataChannelResultTimeout = 5 * time.Second // 评论channel处理超时
)

// 评论话题
type CommentTopicData struct {
	define.CommentTopic `json:"_id" bson:"_id"` // 评论话题
	LastSaveNodeId      int32                   `json:"last_save_node_id" bson:"last_save_node_id"`
	NodeId              int16                   `json:"-" bson:"-"` // 当前节点id
	zsets               *zset.SortedSet         `json:"-" bson:"-"` // 评论按赞排行
	tasker              *task.Tasker            `json:"-" bson:"-"`
	rpcHandler          *RpcHandler             `json:"-" bson:"-"`
}

func NewCommentData() any {
	return &CommentTopicData{}
}

func (c *CommentTopicData) Init(nodeId int16, rpcHandler *RpcHandler) {
	c.LastSaveNodeId = -1
	c.NodeId = nodeId
	c.rpcHandler = rpcHandler
}

func (c *CommentTopicData) InitTask() {
	c.tasker = task.NewTasker()
	c.tasker.Init(
		task.WithStopFns(c.onTaskStop),
		task.WithTimeout(CommentDataTaskTimeout),
	)
}

func (c *CommentTopicData) IsTaskRunning() bool {
	return c.tasker.IsRunning()
}

func (c *CommentTopicData) Load(topic define.CommentTopic) error {
	// 加载话题信息
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Comment, topic, c)

	// 创建新评论数据
	if errors.Is(err, store.ErrNoResult) {
		c.CommentTopic = topic
		c.LastSaveNodeId = int32(c.NodeId)
		errSave := store.GetStore().UpdateOne(context.Background(), define.StoreType_Comment, topic, c, true)
		utils.ErrPrint(errSave, "UpdateOne failed when CommentTopicData.Load", topic)
		return errSave
	}

	if !utils.ErrCheck(err, "FindOne failed when CommentTopicData.Load", topic) {
		return err
	}

	// 加载评论数据
	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Rank, "topic", topic)
	if !utils.ErrCheck(err, "FindAll failed when CommentTopicData.Load", topic) {
		return err
	}

	for _, v := range res {
		vv := v.([]byte)
		metadata := &define.CommentMetadata{}
		err := json.Unmarshal(vv, metadata)
		if !utils.ErrCheck(err, "json.Unmarshal failed when CommentTopicData.Load", vv) {
			continue
		}

		c.zsets.Set(
			float64(metadata.PublisherMetadata.Thumbs),
			metadata.CommentId,
			int64(metadata.PublisherMetadata.Date),
			metadata,
		)
	}

	return nil
}

func (c *CommentTopicData) onTaskStop() {
	log.Info().Caller().Interface("topic", c.CommentTopic).Msg("CommentTopicData task stopped...")
}

func (c *CommentTopicData) TaskRun(ctx context.Context) error {
	return c.tasker.Run(ctx)
}

func (c *CommentTopicData) Stop() {
	c.tasker.Stop()
}

func (c *CommentTopicData) saveLastNode() {
	fields := map[string]any{
		"last_save_node_id": c.NodeId,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Comment, c.CommentTopic, fields, true)
	_ = utils.ErrCheck(err, "UpdateFields failed when CommentTopicData.saveLastNode", c.CommentTopic)
}

func (c *CommentTopicData) AddTask(ctx context.Context, fn task.TaskHandler, p ...any) error {
	return c.tasker.AddWait(ctx, fn, p...)
}

func (c *CommentTopicData) ModThumbs(ctx context.Context, commentId int64, modThumbs int32) error {
	if commentId == -1 {
		return ErrInvalidCommentMetadata
	}

	data, ok := c.zsets.GetData(commentId)
	if !ok {
		return ErrInvalidCommentMetadata
	}

	cm := &define.CommentMetadata{}
	*cm = *data.(*define.CommentMetadata)
	cm.PublisherMetadata.Thumbs *= -1
	cm.PublisherMetadata.Thumbs += modThumbs

	c.zsets.IncrBy(float64(cm.PublisherMetadata.Thumbs), cm.CommentId, int64(cm.PublisherMetadata.Date))

	// save comment metadata
	err := store.GetStore().UpdateOne(ctx, define.StoreType_Comment, cm.CommentId, cm)
	_ = utils.ErrCheck(err, "UpdateOne failed when CommentTopicData.ModThumbs", cm)

	c.saveLastNode()
	return err
}

func (c *CommentTopicData) GetCommentById(ctx context.Context, commentId int64) (rank int64, metadata define.CommentMetadata, err error) {
	zRank, _, data := c.zsets.GetRank(commentId, false)
	rank = zRank
	if data == nil {
		err = ErrCommentNotExist
		return
	}

	metadata = *data.(*define.CommentMetadata)
	metadata.PublisherMetadata.Thumbs *= -1

	return rank, metadata, nil
}

func (c *CommentTopicData) GetCommentByRange(ctx context.Context, start, end int64) (metadatas []*define.CommentMetadata, err error) {
	c.zsets.Range(start, end, func(score float64, key int64, data any) {
		cm := *data.(*define.CommentMetadata)
		cm.PublisherMetadata.Thumbs *= -1
		metadatas = append(metadatas, &cm)
	})
	return
}
