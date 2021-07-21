package rank

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
	rankCleanupInterval = 1 * time.Minute  // cache cleanup interval
	rankCacheExpire     = 10 * time.Minute // cache缓存10分钟
	ErrInvalidOwner     = errors.New("invalid owner")
)

type RankManager struct {
	r          *Rank
	cacheRanks *cache.Cache
	rankPool   sync.Pool
	wg         utils.WaitGroupWrapper
}

func NewRankManager(ctx *cli.Context, r *Rank) *RankManager {
	manager := &RankManager{
		r:          r,
		cacheRanks: cache.New(rankCacheExpire, rankCleanupInterval),
	}

	// 排行榜池
	manager.rankPool.New = NewRankData

	// 排行缓存删除时处理
	manager.cacheRanks.OnEvicted(func(k, v interface{}) {
		v.(*RankData).Stop()
		manager.rankPool.Put(v)
	})

	// 初始化db
	store.GetStore().AddStoreInfo(define.StoreType_LocalRank, "rank", "_id")
	if err := store.GetStore().MigrateDbTable("rank", "owner_id", "rank_list._id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection rank failed")
	}

	log.Info().Msg("RankManager init ok ...")
	return manager
}

func (m *RankManager) Run(ctx *cli.Context) error {
	<-ctx.Done()
	log.Info().Msg("RankManager context done...")
	return nil
}

func (m *RankManager) Exit(ctx *cli.Context) {
	m.wg.Wait()
	log.Info().Msg("RankManager exit...")
}

func (m *RankManager) KickAllRankData() {
	m.cacheRanks.DeleteAll()
}

// 踢掉排行缓存
func (m *RankManager) KickRankData(rankId int64, rankNodeId int32) error {
	if rankId == -1 {
		return nil
	}

	// 踢掉本服rankdata
	if rankNodeId == int32(m.r.ID) {
		mb, ok := m.cacheRanks.Get(rankId)
		if !ok {
			return nil
		}

		mb.(*RankData).Stop()
		store.GetStore().Flush()
		return nil

	} else {
		// rank节点不存在的话不用发送rpc
		nodeId := fmt.Sprintf("rank-%d", rankNodeId)
		srvs, err := m.r.mi.srv.Server().Options().Registry.GetService("rank")
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

		// 发送rpc踢掉其他服rankdata
		rs, err := m.r.rpcHandler.CallKickRankData(rankId, rankNodeId)
		if !utils.ErrCheck(err, "kick rank data failed", rankId, rankNodeId, rs) {
			return err
		}

		// rpc调用成功
		if rs.GetRankId() == rankId {
			return nil
		}

		return errors.New("kick rank data invalid error")
	}
}

// 获取rank数据
func (m *RankManager) getRankData(rankId int64) (*RankData, error) {
	if rankId == -1 {
		return nil, ErrInvalidOwner
	}

	cache, ok := m.cacheRanks.Get(rankId)

	if ok {
		mb := cache.(*RankData)
		if mb.IsTaskRunning() {
			return mb, nil
		}

	} else {

		// 缓存没有，从db加载
		cache = m.rankPool.Get()
		rd := cache.(*RankData)
		rd.Init(m.r.ID, m.r.rpcHandler)
		err := rd.Load(rankId)
		if !utils.ErrCheck(err, "rank data Load failed when RankManager.getRankData", rankId) {
			m.rankPool.Put(cache)
			return nil, err
		}

		// 踢掉上一个节点的缓存
		if rd.LastSaveNodeId != -1 && rd.LastSaveNodeId != int32(m.r.ID) {
			err := m.KickRankData(rd.Id, rd.LastSaveNodeId)
			if !utils.ErrCheck(err, "kick rank data failed", rd.Id, rd.LastSaveNodeId, m.r.ID) {
				return nil, err
			}
		}

		m.cacheRanks.Set(rankId, cache, rankCacheExpire)
	}

	rd := cache.(*RankData)
	rd.InitTask()
	m.wg.Wrap(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			// 立即删除缓存
			m.cacheRanks.Delete(cache.(*RankData).Id)
		}()

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
		err := cache.(*RankData).TaskRun(ctx)
		utils.ErrPrint(err, "rank data run failed", cache.(*RankData).Id)
	})

	return rd, nil
}

func (m *RankManager) AddTask(ctx context.Context, rankId int64, fn task.TaskHandler) error {
	mb, err := m.getRankData(rankId)
	if err != nil {
		return err
	}

	return mb.AddTask(ctx, fn, mb)
}
