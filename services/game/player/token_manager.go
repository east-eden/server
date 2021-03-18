package player

import (
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
)

type TokenManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	owner  *Player `bson:"-" json:"-"`
	Tokens []int32 `bson:"tokens" json:"tokens"`
}

func NewTokenManager(owner *Player) *TokenManager {
	m := &TokenManager{
		owner:  owner,
		Tokens: make([]int32, define.Token_End),
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
	costNum := num
	if costNum <= 0 {
		return fmt.Errorf("token manager cost token<%d> failed, wrong number<%d>", typeMisc, costNum)
	}

	return m.TokenDec(typeMisc, num)
}

func (m *TokenManager) CanGain(typeMisc int32, num int32) error {
	gainNum := int64(num)
	if gainNum <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, gainNum)
	}

	return nil
}

func (m *TokenManager) GainLoot(typeMisc int32, num int32) error {
	gainNum := num
	if gainNum <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, gainNum)
	}

	return m.TokenInc(typeMisc, num)
}

func (m *TokenManager) initTokens() {
	var n int32
	for n = 0; n < define.Token_End; n++ {
		m.Tokens[n] = 0
	}
}

func (m *TokenManager) save(tp int32) error {
	fields := map[string]interface{}{}
	fields[fmt.Sprintf("tokens[%d]", tp)] = m.Tokens[tp]
	return store.GetStore().SaveFields(define.StoreType_Token, m.owner.ID, fields)
}

func (m *TokenManager) LoadAll() error {
	err := store.GetStore().LoadObject(define.StoreType_Token, m.owner.GetID(), m)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("TokenManager LoadAll: %w", err)
	}

	return nil
}

func (m *TokenManager) TokenInc(tp int32, value int32) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	if m.Tokens[tp]+value < 0 {
		return fmt.Errorf("token<%d> overflow when TokenInc", tp)
	}

	entry, ok := auto.GetTokenEntry(tp)
	if !ok {
		return fmt.Errorf("GetTokenEntry<%d> failed when TokenInc", tp)
	}

	m.Tokens[tp] += value
	if m.Tokens[tp] > entry.MaxHold {
		m.Tokens[tp] = entry.MaxHold
	}

	err := m.save(tp)
	m.SendTokenUpdate(tp, m.Tokens[tp])
	return err
}

func (m *TokenManager) TokenDec(tp int32, value int32) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	m.Tokens[tp] -= value
	if m.Tokens[tp] < 0 {
		m.Tokens[tp] = 0
	}

	err := m.save(tp)
	m.SendTokenUpdate(tp, m.Tokens[tp])
	return err
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
