package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/spf13/cast"
	"github.com/valyala/bytebufferpool"
)

var (
	strengthRegenInterval = time.Minute * 5 // 体力每5分钟更新一次
	ErrInvalidTokenType   = errors.New("invalid token type")
)

func makeTokenKey(tp int32) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("tokens.")
	_, _ = b.WriteString(cast.ToString(tp))

	return b.String()
}

type TokenManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`
	NextStrengthRegenTime int32 `bson:"next_strength_regen_time" json:"next_strength_regen_time"` // 下次体力恢复时间

	owner  *Player                 `bson:"-" json:"-"`
	Tokens [define.Token_End]int32 `bson:"tokens" json:"tokens"`
}

func NewTokenManager(owner *Player) *TokenManager {
	m := &TokenManager{
		owner:                 owner,
		NextStrengthRegenTime: int32(time.Now().Unix()),
	}

	// init tokens
	m.initTokens()

	return m
}

// interface of cost_loot
func (m *TokenManager) GetCostLootType() int32 {
	return define.CostLoot_Token
}

func (m *TokenManager) CanCost(typeMisc int32, num int32) error {
	costNum := num
	if costNum <= 0 {
		return fmt.Errorf("token manager check token<%d> cost failed, wrong number<%d>", typeMisc, costNum)
	}

	if !utils.BetweenInt32(typeMisc, define.Token_Begin, define.Token_End) {
		return errors.New("invalid token type")
	}

	if m.Tokens[typeMisc] < costNum {
		return errors.New("not enough token")
	}

	return nil
}

func (m *TokenManager) DoCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.DoCost(typeMisc, num)
	if err != nil {
		return err
	}

	if !utils.BetweenInt32(typeMisc, 0, define.Token_End) {
		return ErrInvalidTokenType
	}

	m.Tokens[typeMisc] -= num
	if m.Tokens[typeMisc] < 0 {
		m.Tokens[typeMisc] = 0
	}

	err = m.save(typeMisc)
	m.SendTokenUpdate(typeMisc, m.Tokens[typeMisc])
	return err
}

func (m *TokenManager) CanGain(typeMisc int32, num int32) error {
	gainNum := int64(num)
	if gainNum <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, gainNum)
	}

	return nil
}

func (m *TokenManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	if !utils.BetweenInt32(typeMisc, 0, define.Token_End) {
		return ErrInvalidTokenType
	}

	if m.Tokens[typeMisc]+num < 0 {
		return fmt.Errorf("token<%d> overflow when TokenInc", typeMisc)
	}

	entry, ok := auto.GetTokenEntry(typeMisc)
	if !ok {
		return fmt.Errorf("GetTokenEntry<%d> failed when TokenInc", typeMisc)
	}

	m.Tokens[typeMisc] += num
	if m.Tokens[typeMisc] > entry.MaxHold {
		m.tokenOverflow(typeMisc, entry.MaxHold)
	}

	err = m.save(typeMisc)
	m.SendTokenUpdate(typeMisc, m.Tokens[typeMisc])
	return err
}

func (m *TokenManager) initTokens() {
	var n int32
	for n = 0; n < define.Token_End; n++ {
		m.Tokens[n] = 0
	}
}

func (m *TokenManager) save(tp int32) error {
	fields := map[string]interface{}{
		makeTokenKey(tp): m.Tokens[tp],
	}
	return store.GetStore().UpdateFields(context.Background(), define.StoreType_Token, m.owner.ID, fields)
}

func (m *TokenManager) update() {
	// 体力恢复
	if m.NextStrengthRegenTime > int32(time.Now().Unix()) {
		return
	}

	tm := time.Unix(int64(m.NextStrengthRegenTime), 0)
	d := time.Since(tm)
	times := d / strengthRegenInterval

	// 设置下次更新时间
	m.NextStrengthRegenTime = int32(time.Now().Add(strengthRegenInterval).Unix())
	fields := map[string]interface{}{
		"next_strength_regen_time": m.NextStrengthRegenTime,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Token, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when TokenMananger.update", m.owner.ID, fields)

	// 恢复体力
	_ = m.GainLoot(define.Token_Strength, int32(1+times))
}

func (m *TokenManager) tokenOverflow(tp int32, maxHold int32) {
	switch tp {
	case define.Token_Strength:
		_ = m.GainLoot(define.Token_StrengthStore, m.Tokens[tp]-maxHold)
		m.Tokens[tp] = maxHold
	}
}

func (m *TokenManager) LoadAll() error {
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Token, m.owner.GetId(), m)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("TokenManager LoadAll: %w", err)
	}

	return nil
}

func (m *TokenManager) TokenSet(tp int32, value int32) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	if value < 0 {
		return fmt.Errorf("token<%d> set invalid value<%d>", tp, value)
	}

	entry, ok := auto.GetTokenEntry(tp)
	if !ok {
		return fmt.Errorf("GetTokenEntry<%d> failed when TokenInc", tp)
	}

	m.Tokens[tp] = value
	if m.Tokens[tp] > entry.MaxHold {
		m.Tokens[tp] = entry.MaxHold
	}

	err := m.save(tp)
	m.SendTokenUpdate(tp, m.Tokens[tp])
	return err
}

func (m *TokenManager) GetToken(tp int32) (int32, error) {
	if tp < 0 || tp >= define.Token_End {
		return 0, fmt.Errorf("token type<%d> invalid", tp)
	}

	return m.Tokens[tp], nil
}

func (m *TokenManager) SendTokenUpdate(tp, value int32) {
	msg := &pbGlobal.S2C_TokenUpdate{
		Type:  tp,
		Value: value,
	}

	m.owner.SendProtoMessage(msg)
}
