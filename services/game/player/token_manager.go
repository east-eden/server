package player

import (
	"errors"
	"fmt"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/store"
	log "github.com/rs/zerolog/log"
)

type TokenManager struct {
	store.StoreObjector `bson:"-" json:"-"`
	owner               *Player `bson:"-" json:"-"`
	OwnerType           int32   `bson:"owner_type" json:"owner_type"`
	Tokens              []int32 `bson:"tokens" json:"tokens"`
}

func NewTokenManager(owner *Player) *TokenManager {
	m := &TokenManager{
		owner:     owner,
		OwnerType: owner.GetType(),
		Tokens:    make([]int32, 0, define.Token_End),
	}

	// init tokens
	m.initTokens()

	return m
}

func (m *TokenManager) AfterLoad() error {
	return nil
}

func (m *TokenManager) GetObjID() int64 {
	return m.owner.GetID()
}

func (m *TokenManager) GetStoreIndex() int64 {
	return -1
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

	for k, v := range m.Tokens {
		if int32(k) == typeMisc {
			if v >= costNum {
				return nil
			}
		}
	}

	return fmt.Errorf("not enough token<%d>, num<%d>", typeMisc, costNum)
}

func (m *TokenManager) DoCost(typeMisc int32, num int32) error {
	costNum := num
	if costNum <= 0 {
		return fmt.Errorf("token manager cost token<%d> failed, wrong number<%d>", typeMisc, costNum)
	}

	for k := range m.Tokens {
		if int32(k) == typeMisc {
			if m.Tokens[k] < costNum {
				log.Warn().
					Int32("cost_type_misc", typeMisc).
					Int32("cost_num", costNum).
					Int32("actual_cost_num", m.Tokens[k]).
					Msg("token manager cost number error")
			}

			m.Tokens[k] -= costNum
			if m.Tokens[k] < 0 {
				m.Tokens[k] = 0
			}

			break
		}
	}

	err := m.save()
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
	gainNum := num
	if gainNum <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, gainNum)
	}

	for k := range m.Tokens {
		if int32(k) == typeMisc {
			entry, ok := auto.GetTokenEntry(int32(k))
			if !ok {
				return fmt.Errorf("GetTokenEntry<%d> failed when GainLoot", k)
			}

			if m.Tokens[k]+gainNum < 0 {
				return fmt.Errorf("token overflow when GainLoot")
			}

			m.Tokens[k] += gainNum
			if m.Tokens[k] > entry.MaxHold {
				m.Tokens[k] = entry.MaxHold
			}

			break
		}
	}

	err := m.save()
	return err
}

func (m *TokenManager) initTokens() {
	var n int32
	for n = 0; n < define.Token_End; n++ {
		m.Tokens[n] = 0
	}
}

func (m *TokenManager) save() error {
	fields := map[string]interface{}{
		"tokens": m.Tokens,
	}
	return store.GetStore().SaveFields(define.StoreType_Token, m, fields)
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

	err := m.save()
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

	err := m.save()
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

	err := m.save()
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
