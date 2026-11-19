// Package shop 实现配置驱动的商店系统。
//
// kof 的 shop_manager 有 1100 多行，包含分层随机商品、售罄重置、平台支付订单等。
// 这里移植其中与配置/账本相关的核心：
//
//   - 商店与商品来自 Shop / ShopGoods 表
//   - 购买校验：等级限制、限购、货币是否足够
//   - 扣费与发货统一走物品账本与 kof 的 Loot 表
//   - 按 cron 重置限购次数（商店级与商品级各自支持）
//
// 平台支付（platformGoodsId）与分层随机商品不在本包内，需要接入支付/随机商品后再补。
package shop

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/utils"
)

var (
	ErrShopNotFound       = errors.New("shop not found")
	ErrGoodsNotFound      = errors.New("shop goods not found")
	ErrGoodsLevelLimit    = errors.New("shop goods level limit")
	ErrGoodsSoldOut       = errors.New("shop goods sold out")
	ErrGoodsCostNotEnough = errors.New("shop goods cost not enough")
	ErrBuyCountInvalid    = errors.New("shop buy count invalid")
)

// GoodsState 商品购买状态
type GoodsState struct {
	Id      int32 `bson:"id" json:"id"`
	Bought  int32 `bson:"bought" json:"bought"`     // 已购买数量
	ResetAt int64 `bson:"reset_at" json:"reset_at"` // 上次重置时间(毫秒)
}

// ShopState 商店状态
type ShopState struct {
	Id      int32                 `bson:"id" json:"id"`
	ResetAt int64                 `bson:"reset_at" json:"reset_at"` // 商店级上次重置时间(毫秒)
	Goods   map[int32]*GoodsState `bson:"goods" json:"goods"`
}

// Owner 商店所需的宿主能力，由 Player 实现
type Owner interface {
	GetLevel() int32
	GetItemCount(entryId int32) int64
	CostItems(costs []*config.IdCount, reason string) error
	GrantItems(items []*config.IdCount, reason string)
	GrantLootById(lootId int32, reason string)
	SaveShops(data map[int32]*ShopState) error
	SendProtoMessage(msg proto.Message)
}

// Manager 玩家商店管理器
type Manager struct {
	owner Owner
	Shops map[int32]*ShopState
}

func NewManager(owner Owner) *Manager {
	return &Manager{
		owner: owner,
		Shops: make(map[int32]*ShopState, 8),
	}
}

// LoadAll 接上读档数据
func (m *Manager) LoadAll(data map[int32]*ShopState) {
	if data != nil {
		m.Shops = data
	}

	if m.Shops == nil {
		m.Shops = make(map[int32]*ShopState, 8)
	}
}

func (m *Manager) ensureShop(shopId int32) *ShopState {
	state, ok := m.Shops[shopId]
	if !ok {
		state = &ShopState{Id: shopId, Goods: make(map[int32]*GoodsState, 8)}
		m.Shops[shopId] = state
	}

	if state.Goods == nil {
		state.Goods = make(map[int32]*GoodsState, 8)
	}

	return state
}

func (m *Manager) ensureGoods(state *ShopState, goodsId int32) *GoodsState {
	goods, ok := state.Goods[goodsId]
	if !ok {
		goods = &GoodsState{Id: goodsId}
		state.Goods[goodsId] = goods
	}

	return goods
}

// resetDue 判断 cron 周期是否已经跨过（上次重置之后是否有新的周期点）
func resetDue(cronSpec string, lastAt int64, now time.Time) bool {
	if cronSpec == "" {
		return false
	}

	cron, err := utils.ParseCron(cronSpec)
	if err != nil {
		utils.ErrPrint(err, "parse shop reset cron failed", cronSpec)
		return false
	}

	prev := cron.Prev(now)
	if prev.IsZero() {
		return false
	}

	return prev.UnixMilli() > lastAt
}

// CheckAll 按配置刷新商店与商品的限购次数
func (m *Manager) CheckAll(now time.Time) {
	for shopId, entry := range config.GetShopRows() {
		state := m.ensureShop(shopId)

		// 商店级重置
		if resetDue(entry.Reset, state.ResetAt, now) {
			state.ResetAt = now.UnixMilli()
			for _, goods := range state.Goods {
				goods.Bought = 0
			}
		}

		// 商品级重置
		for _, goodsId := range entry.Goods {
			goodsEntry, ok := config.GetShopGoodsEntry(goodsId)
			if !ok {
				continue
			}

			goods := m.ensureGoods(state, goodsId)
			if resetDue(goodsEntry.Reset, goods.ResetAt, now) {
				goods.ResetAt = now.UnixMilli()
				goods.Bought = 0
			}
		}
	}
}

// GoodsRemain 返回剩余可购买数量，limit<=0 表示不限购
func (m *Manager) GoodsRemain(shopId, goodsId int32) int32 {
	entry, ok := config.GetShopGoodsEntry(goodsId)
	if !ok || entry.Limit <= 0 {
		return -1
	}

	state := m.ensureShop(shopId)
	goods := m.ensureGoods(state, goodsId)

	return max(entry.Limit-goods.Bought, 0)
}

// checkBuy 校验购买条件
func (m *Manager) checkBuy(shopId, goodsId, count int32) (*config.ShopGoodsEntry, *GoodsState, error) {
	if _, ok := config.GetShopEntry(shopId); !ok {
		return nil, nil, ErrShopNotFound
	}

	entry, ok := config.GetShopGoodsEntry(goodsId)
	if !ok {
		return nil, nil, ErrGoodsNotFound
	}

	if count <= 0 {
		return nil, nil, ErrBuyCountInvalid
	}

	if entry.LevelLimit > 0 && m.owner.GetLevel() < entry.LevelLimit {
		return nil, nil, ErrGoodsLevelLimit
	}

	state := m.ensureShop(shopId)
	goods := m.ensureGoods(state, goodsId)

	if entry.Limit > 0 && goods.Bought+count > entry.Limit {
		return nil, nil, fmt.Errorf("%w: bought %d, limit %d", ErrGoodsSoldOut, goods.Bought, entry.Limit)
	}

	return entry, goods, nil
}

// scaleCost 按购买份数放大消耗
func scaleCost(costs []*config.IdCount, count int32) []*config.IdCount {
	out := make([]*config.IdCount, 0, len(costs))
	for _, c := range costs {
		if c == nil || c.Id <= 0 || c.Count <= 0 {
			continue
		}

		out = append(out, &config.IdCount{Id: c.Id, Count: c.Count * count})
	}

	return out
}

// Buy 购买商品，count 为购买份数
func (m *Manager) Buy(shopId, goodsId, count int32) error {
	entry, goods, err := m.checkBuy(shopId, goodsId, count)
	if err != nil {
		return err
	}

	// 扣费（货币/道具统一走账本）
	if err := m.owner.CostItems(scaleCost(entry.CostItems, count), "shopBuy"); err != nil {
		return fmt.Errorf("%w: %v", ErrGoodsCostNotEnough, err)
	}

	// 发货：固定道具 + 掉落
	if len(entry.GainItems) > 0 {
		m.owner.GrantItems(scaleCost(entry.GainItems, count), "shopBuy")
	}

	for i := int32(0); i < count; i++ {
		if entry.LootId > 0 {
			m.owner.GrantLootById(entry.LootId, "shopBuy")
		}
	}

	goods.Bought += count

	if err := m.owner.SaveShops(m.Shops); err != nil {
		return err
	}

	m.SendUpdate()
	return nil
}

// GenPB 生成下发协议
func (m *Manager) GenPB() []*pbGlobal.ShopInfo {
	list := make([]*pbGlobal.ShopInfo, 0, len(config.GetShopRows()))
	for shopId, entry := range config.GetShopRows() {
		if entry.OpenLv > 0 && m.owner.GetLevel() < entry.OpenLv {
			continue
		}

		state := m.ensureShop(shopId)
		info := &pbGlobal.ShopInfo{Id: shopId}

		for _, goodsId := range entry.Goods {
			goodsEntry, ok := config.GetShopGoodsEntry(goodsId)
			if !ok || goodsEntry.IsHide {
				continue
			}

			goods := m.ensureGoods(state, goodsId)
			remain := int32(-1)
			if goodsEntry.Limit > 0 {
				remain = max(goodsEntry.Limit-goods.Bought, 0)
			}

			info.Goods = append(info.Goods, &pbGlobal.ShopGoodsInfo{
				Id:     goodsId,
				Bought: goods.Bought,
				Remain: remain,
			})
		}

		list = append(list, info)
	}

	return list
}

// SendUpdate 下发商店状态
func (m *Manager) SendUpdate() {
	m.owner.SendProtoMessage(&pbGlobal.S2C_ShopUpdate{Shops: m.GenPB()})
}
