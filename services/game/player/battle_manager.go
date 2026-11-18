package player

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/battle"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/hero"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
)

var (
	ErrBattleNotStarted     = errors.New("battle not started")
	ErrBattleTokenInvalid   = errors.New("battle verify token invalid")
	ErrBattleResultMismatch = errors.New("client battle result mismatch with server")
)

// PlayerBattleRecord 玩家最近一次战斗记录
type PlayerBattleRecord struct {
	BattleType     int32                            `bson:"battle_type" json:"battle_type"`
	BattleId       string                           `bson:"battle_id" json:"battle_id"`
	BattleIndex    int32                            `bson:"battle_index" json:"battle_index"`
	IsServerBattle bool                             `bson:"is_server_battle" json:"is_server_battle"`
	LevelId        int32                            `bson:"level_id" json:"level_id"`
	Dmg            float64                          `bson:"dmg" json:"dmg"`
	Win            bool                             `bson:"win" json:"win"`
	Winner         int64                            `bson:"winner" json:"winner"`
	BattleTime     float32                          `bson:"battle_time" json:"battle_time"`
	BossSurplusHp  float32                          `bson:"boss_surplus_hp" json:"boss_surplus_hp"`
	BossTotalHp    float32                          `bson:"boss_total_hp" json:"boss_total_hp"`
	ServerResult   *pbGlobal.RespServerBattleResult `bson:"-" json:"-"` // 战斗服原始结果，仅内存保留
}

// BattleManager 战斗管理器。
//
// 与 kof 一致：开战时生成一次性校验 token，把战斗请求发给外部战斗服；
// 客户端提交战斗结果时校验 token，并与服务器战斗结果比对，防止改包。
type BattleManager struct {
	owner        *Player
	BattleRecord *PlayerBattleRecord `bson:"battle_record" json:"battle_record"`

	battleServerURL string
	tokens          map[int32]string // 战斗类型 -> 校验 token
	lock            sync.Mutex
}

func NewBattleManager(p *Player, battleServerURL string) *BattleManager {
	return &BattleManager{
		owner:           p,
		battleServerURL: battleServerURL,
		tokens:          make(map[int32]string, 4),
	}
}

// Migrate 在玩家读档后把存档中的战斗记录接上
func (m *BattleManager) Migrate() {
	if m.BattleRecord == nil {
		m.BattleRecord = &PlayerBattleRecord{}
	}
}

// GenBattleVerifyToken 生成一次性战斗校验 token
func (m *BattleManager) GenBattleVerifyToken(battleType pbGlobal.EMBattleType) string {
	id, err := utils.NextID(define.SnowFlake_Battle)
	if err != nil {
		log.Error().Err(err).Msg("gen battle verify token failed")
		return ""
	}

	token := strconv.FormatInt(id, 10)

	m.lock.Lock()
	m.tokens[int32(battleType)] = token
	m.lock.Unlock()

	return token
}

// GetBattleVerifyToken 读取当前 token（不消费）
func (m *BattleManager) GetBattleVerifyToken(battleType int32) (string, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	token, ok := m.tokens[battleType]
	return token, ok
}

// ConsumeBattleVerifyToken 作废 token
func (m *BattleManager) ConsumeBattleVerifyToken(battleType int32) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.tokens, battleType)
}

// BuildBattlePlayer 用上阵英雄构造战斗单位
func (m *BattleManager) BuildBattlePlayer() *pbGlobal.BattlePlayer {
	p := m.owner

	bp := &pbGlobal.BattlePlayer{
		Id:        p.ID,
		Name:      p.Name,
		Level:     p.Level,
		Formation: make([]int64, 0, len(p.BattleArray)),
		Heros:     make([]*pbGlobal.BattleHeroCard, 0, len(p.BattleArray)),
	}

	for _, heroId := range p.BattleArray {
		if heroId == -1 {
			continue
		}

		card := p.HeroManager().GetHero(heroId)
		if card == nil {
			continue
		}

		bp.Formation = append(bp.Formation, int64(card.ConfigId))
		bp.Heros = append(bp.Heros, battleHeroCard(card))
	}

	// 未布阵时按配置 id 取前几个英雄参战
	if len(bp.Heros) == 0 {
		cards := p.HeroManager().GetHeroList()
		sort.Slice(cards, func(i, j int) bool { return cards[i].ConfigId < cards[j].ConfigId })

		for i, card := range cards {
			if i >= PlayerBattleArrayMaxHero {
				break
			}

			bp.Formation = append(bp.Formation, int64(card.ConfigId))
			bp.Heros = append(bp.Heros, battleHeroCard(card))
		}
	}

	return bp
}

func battleHeroCard(card *hero.Card) *pbGlobal.BattleHeroCard {
	return &pbGlobal.BattleHeroCard{
		ConfigId:    card.ConfigId,
		Level:       card.Level,
		Star:        card.Star,
		PromotionLv: card.PromotionLv,
		Rarity:      card.Rarity,
		Attrs:       card.Attrs,
		BattlePower: card.BattlePower,
	}
}

// CreateCommonBattle 发起一场服务器战斗。
//
// 返回发给战斗服的请求；战斗服不可用时返回请求本身，由调用方决定降级策略。
func (m *BattleManager) CreateCommonBattle(
	battleType pbGlobal.EMBattleType,
	levelId int32,
	atkPlayer, defPlayer *pbGlobal.BattlePlayer,
	extra *pbGlobal.BattleExtraData,
) (*pbGlobal.ReqServerBattleStart, *pbGlobal.RespServerBattleResult, error) {
	if atkPlayer == nil {
		atkPlayer = m.BuildBattlePlayer()
	}

	req := battle.BuildServerBattleRequest(battleType, levelId, atkPlayer, defPlayer, extra)

	m.BattleRecord = &PlayerBattleRecord{
		BattleType:     int32(battleType),
		BattleId:       req.BattleId,
		BattleIndex:    extra.GetBattleIndex(),
		IsServerBattle: m.battleServerURL != "",
		LevelId:        levelId,
	}

	var result *pbGlobal.RespServerBattleResult
	if m.battleServerURL != "" {
		var err error
		result, err = battle.StartServerBattle(m.battleServerURL, req)
		if err != nil {
			return req, nil, fmt.Errorf("start server battle failed: %w", err)
		}

		m.BattleRecord.Win = result.Win
		m.BattleRecord.Winner = result.Winner
		m.BattleRecord.Dmg = float64(result.TotalDmg)
		m.BattleRecord.BossSurplusHp = result.BossSurplusHp
		m.BattleRecord.BossTotalHp = result.BossTotalHp
		m.BattleRecord.ServerResult = result
	}

	return req, result, nil
}

// MergeBattleRecord 合并客户端上报的战斗结果。
//
// 校验顺序：token -> 服务器战斗结果一致性 -> 写入战斗记录。
// 没有服务器战斗结果时以客户端结果为准（PVE 本地战斗）。
func (m *BattleManager) MergeBattleRecord(sign *pbGlobal.ChallengeSignInfo, battleType int32) error {
	if sign == nil {
		return ErrBattleNotStarted
	}

	// 先校验再消费：token 不匹配时不能让合法的 token 失效
	expect, ok := m.GetBattleVerifyToken(battleType)
	if !ok || expect == "" || sign.Token != expect {
		return ErrBattleTokenInvalid
	}

	m.ConsumeBattleVerifyToken(battleType)

	if m.BattleRecord == nil {
		return ErrBattleNotStarted
	}

	// 服务器战斗：结果必须与战斗服一致
	if m.BattleRecord.IsServerBattle {
		if m.BattleRecord.ServerResult == nil {
			return ErrBattleNotStarted
		}

		if sign.Win != m.BattleRecord.ServerResult.Win {
			return ErrBattleResultMismatch
		}
	}

	m.BattleRecord.Win = sign.Win
	m.BattleRecord.Dmg = sign.Dmg

	return nil
}

// Record 返回最近一次战斗记录
func (m *BattleManager) Record() *PlayerBattleRecord {
	return m.BattleRecord
}

// Reset 清空战斗记录（战斗结算完成后调用）
func (m *BattleManager) Reset() {
	m.BattleRecord = &PlayerBattleRecord{}
}
