package comment

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/cache"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	commentCleanupInterval       = 1 * time.Minute // cache cleanup interval
	commentCacheExpire           = 1 * time.Hour   // cache缓存1小时
	commentDefaultLoad     int64 = 10              // 默认加载前10条评论
)

type CommentManager struct {
	r                 *Comment
	cacheCommentDatas *cache.Cache
	commentPool       sync.Pool
	wg                utils.WaitGroupWrapper
	mu                sync.Mutex
}

func NewCommentManager(ctx *cli.Context, r *Comment) *CommentManager {
	manager := &CommentManager{
		r:                 r,
		cacheCommentDatas: cache.New(commentCacheExpire, commentCleanupInterval),
	}

	// 排行榜池
	manager.commentPool.New = NewCommentData

	// 排行缓存删除时处理
	manager.cacheCommentDatas.OnEvicted(func(k, v interface{}) {
		v.(*CommentTopicData).Stop()
		manager.commentPool.Put(v)
	})

	// 初始化db
	store.GetStore().AddStoreInfo(define.StoreType_Comment, "comment", "_id")
	if err := store.GetStore().MigrateDbTable("comment"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection comment failed")
	}

	log.Info().Msg("CommentManager init ok ...")
	return manager
}

func (m *CommentManager) Run(ctx *cli.Context) error {
	<-ctx.Done()
	log.Info().Msg("CommentManager context done...")
	return nil
}

func (m *CommentManager) Exit(ctx *cli.Context) {
	m.wg.Wait()
	log.Info().Msg("CommentManager exit...")
}

func (m *CommentManager) KickAllCommentTopicData() {
	m.cacheCommentDatas.DeleteAll()
}

// 踢掉评论topic缓存
func (m *CommentManager) KickCommentTopicData(topic define.CommentTopic, commentNodeId int32) error {
	if !topic.Valid() {
		return nil
	}

	// 踢掉本服CommentTopicData
	if commentNodeId == int32(m.r.ID) {
		topicId := utils.PackId(topic.Type, topic.TypeId)
		cd, ok := m.cacheCommentDatas.Get(topicId)
		if !ok {
			return nil
		}

		cd.(*CommentTopicData).Stop()
		store.GetStore().Flush()
		return nil

	} else {
		// comment节点不存在的话不用发送rpc
		nodeId := fmt.Sprintf("comment-%d", commentNodeId)
		srvs, err := m.r.mi.srv.Options().Registry.GetService("comment")
		if err != nil {
			return nil
		}

		hit := false
		for _, srv := range srvs {
			for _, node := range srv.Nodes {
				if node.Id == nodeId {
					hit = true
					break
				}
			}
		}

		if !hit {
			return nil
		}

		// 发送rpc踢掉其他服CommentTopicData
		rs, err := m.r.rpcHandler.CallKickCommentTopicData(topic, commentNodeId)
		if !utils.ErrCheck(err, "kick comment topic data failed", topic, commentNodeId, rs) {
			return err
		}

		// rpc调用成功
		if rs.GetTopic().GetTopicType() == topic.Type && rs.GetTopic().GetTopicTypeId() == topic.TypeId {
			return nil
		}

		return errors.New("kick comment topic data invalid error")
	}
}

// 获取comment数据
func (m *CommentManager) getCommentData(topic define.CommentTopic) (*CommentTopicData, error) {
	if !topic.Valid() {
		return nil, ErrInvalidComment
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	topicId := utils.PackId(topic.Type, topic.TypeId)
	cache, ok := m.cacheCommentDatas.Get(topicId)

	if ok {
		rd := cache.(*CommentTopicData)
		if rd.IsTaskRunning() {
			return rd, nil
		}

	} else {

		// 缓存没有，从db加载
		cache = m.commentPool.Get()
		cd := cache.(*CommentTopicData)
		cd.Init(m.r.ID, m.r.rpcHandler)
		err := cd.Load(topic)
		if !utils.ErrCheck(err, "CommentTopicData Load failed when CommentManager.getCommentData", topic) {
			m.commentPool.Put(cache)
			return nil, err
		}

		// 踢掉上一个节点的缓存
		if cd.LastSaveNodeId != -1 && cd.LastSaveNodeId != int32(m.r.ID) {
			err := m.KickCommentTopicData(topic, cd.LastSaveNodeId)
			if !utils.ErrCheck(err, "kick CommentTopicData failed", topic, cd.LastSaveNodeId, m.r.ID) {
				return nil, err
			}
		}

		m.cacheCommentDatas.Set(topicId, cache, commentCacheExpire)
	}

	cd := cache.(*CommentTopicData)
	cd.InitTask()
	m.wg.Wrap(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			// 立即删除缓存
			topicId := utils.PackId(cache.(*CommentTopicData).Type, cache.(*CommentTopicData).TypeId)
			m.cacheCommentDatas.Delete(topicId)
		}()

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
		err := cache.(*CommentTopicData).TaskRun(ctx)
		utils.ErrPrint(err, "CommentData run failed", cache.(*CommentTopicData).CommentTopic)
	})

	return cd, nil
}

func (m *CommentManager) AddTask(ctx context.Context, topic define.CommentTopic, fn task.TaskHandler) error {
	cd, err := m.getCommentData(topic)

	if err != nil {
		return err
	}

	return cd.AddTask(ctx, fn, cd)
}

func (m *CommentManager) QueryCommentTopic(ctx context.Context, topic define.CommentTopic) (metadatas []*define.CommentMetadata, err error) {
	err = m.AddTask(
		ctx,
		topic,
		func(c context.Context, p ...interface{}) error {
			var e error
			ctd := p[0].(*CommentTopicData)
			metadatas, e = ctd.GetCommentByRange(c, 0, commentDefaultLoad)
			return e
		},
	)

	_ = utils.ErrCheck(err, "AddTask failed when CommentManager.QueryCommentTopic", topic)
	return
}

func (m *CommentManager) QueryCommentTopicRange(ctx context.Context, topic define.CommentTopic, start, end int64) (metadatas []*define.CommentMetadata, err error) {
	err = m.AddTask(
		ctx,
		topic,
		func(c context.Context, p ...interface{}) error {
			var e error
			ctd := p[0].(*CommentTopicData)
			metadatas, e = ctd.GetCommentByRange(c, start, end)
			return e
		},
	)

	_ = utils.ErrCheck(err, "AddTask failed when CommentManager.QueryCommentTopic", topic, start, end)
	return
}

func (m *CommentManager) ModCommentThumbs(ctx context.Context, topic define.CommentTopic, commentId int64, modThumbs int32) error {
	err := m.AddTask(
		ctx,
		topic,
		func(c context.Context, p ...interface{}) error {
			ctd := p[0].(*CommentTopicData)
			return ctd.ModThumbs(ctx, commentId, modThumbs)
		},
	)

	_ = utils.ErrCheck(err, "AddTask failed when CommentManager.ModCommentThumbs", topic, commentId, modThumbs)
	return err
}
