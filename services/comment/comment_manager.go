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

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"e.coding.net/mmstudio/blade/server/utils/cache"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	commentCleanupInterval = 1 * time.Minute // cache cleanup interval
	commentCacheExpire     = 1 * time.Hour   // cache缓存1小时
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

func (m *CommentManager) KickAllCommentData() {
	m.cacheCommentDatas.DeleteAll()
}

// 踢掉排行缓存
func (m *CommentManager) KickCommentData(commentId int32, commentNodeId int32) error {
	if commentId == -1 {
		return nil
	}

	// 踢掉本服CommentData
	if commentNodeId == int32(m.r.ID) {
		cd, ok := m.cacheCommentDatas.Get(commentId)
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
		rs, err := m.r.rpcHandler.CallKickCommentTopicData(commentId, commentNodeId)
		if !utils.ErrCheck(err, "kick comment topic data failed", commentId, commentNodeId, rs) {
			return err
		}

		// rpc调用成功
		if rs.GetCommentId() == commentId {
			return nil
		}

		return errors.New("kick comment topic data invalid error")
	}
}

// 获取comment数据
func (m *CommentManager) getCommentData(commentId int32) (*CommentTopicData, error) {
	if commentId == -1 {
		return nil, ErrInvalidComment
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	cache, ok := m.cacheCommentDatas.Get(commentId)

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
		err := cd.Load(commentId)
		if !utils.ErrCheck(err, "CommentData Load failed when CommentManager.getCommentData", commentId) {
			m.commentPool.Put(cache)
			return nil, err
		}

		// 踢掉上一个节点的缓存
		if cd.LastSaveNodeId != -1 && cd.LastSaveNodeId != int32(m.r.ID) {
			err := m.KickCommentData(cd.CommentId, cd.LastSaveNodeId)
			if !utils.ErrCheck(err, "kick CommentData failed", cd.CommentId, cd.LastSaveNodeId, m.r.ID) {
				return nil, err
			}
		}

		m.cacheCommentDatas.Set(commentId, cache, commentCacheExpire)
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
			m.cacheCommentDatas.Delete(cache.(*CommentTopicData).CommentId)
		}()

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
		err := cache.(*CommentTopicData).TaskRun(ctx)
		utils.ErrPrint(err, "CommentData run failed", cache.(*CommentTopicData).CommentId)
	})

	return cd, nil
}

func (m *CommentManager) AddTask(ctx context.Context, commentId int32, fn task.TaskHandler) error {
	cd, err := m.getCommentData(commentId)

	if err != nil {
		return err
	}

	return cd.AddTask(ctx, fn, cd)
}
