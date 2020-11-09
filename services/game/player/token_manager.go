package player

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/store"
)

type Token struct {
	ID      int32              `bson:"token_id" json:"token_id"`
	Value   int64              `bson:"token_value" json:"token_value"`
	MaxHold int64              `bson:"token_max_hold" json:"token_max_hold"`
	Entry   *define.TokenEntry `bson:"-" json:"-"`
}

type TokenManager struct {
	store.StoreObjector `bson:"-" json:"-"`
	owner               *Player  `bson:"-" json:"-"`
	OwnerType           int32    `bson:"owner_type" json:"owner_type"`
	Tokens              []*Token `bson:"tokens" json:"tokens"`

	sync.RWMutex `bson:"-" json:"-"`
}

func NewTokenManager(owner *Player) *TokenManager {
	m := &TokenManager{
		owner:     owner,
		OwnerType: owner.GetType(),
		Tokens:    make([]*Token, 0),
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
	costNum := int64(num)
	if costNum <= 0 {
		return fmt.Errorf("token manager check token<%d> cost failed, wrong number<%d>", typeMisc, costNum)
	}

	for _, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.Value >= costNum {
				return nil
			}
		}
	}

	return fmt.Errorf("not enough token<%d>, num<%d>", typeMisc, costNum)
}

func (m *TokenManager) DoCost(typeMisc int32, num int32) error {
	costNum := int64(num)
	if costNum <= 0 {
		return fmt.Errorf("token manager cost token<%d> failed, wrong number<%d>", typeMisc, costNum)
	}

	for _, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.Value < costNum {
				log.Warn().
					Int32("cost_type_misc", typeMisc).
					Int64("cost_num", costNum).
					Int64("actual_cost_num", v.Value).
					Msg("token manager cost number error")
			}

			v.Value -= costNum
			if v.Value < 0 {
				v.Value = 0
			}

			break
		}
	}

	m.save()
	return nil
}

func (m *TokenManager) CanGain(typeMisc int32, num int32) error {
	gainNum := int64(num)
	if gainNum <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, gainNum)
	}

	return nil
}

func (m *TokenManager) GainLoot(typeMisc int32, num int32) error {
	gainNum := int64(num)
	if gainNum <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, gainNum)
	}

	for _, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.MaxHold < v.Value+gainNum {
				log.Warn().
					Int32("gain_type_misc", typeMisc).
					Int64("gain_num", gainNum).
					Int64("actual_gain_num", v.MaxHold-v.Value).
					Msg("token manager gain number overflow")
			}

			v.Value += gainNum
			if v.Value > v.MaxHold {
				v.Value = v.MaxHold
			}

			break
		}
	}

	m.save()
	return nil
}

func (m *TokenManager) initTokens() {
	for n := 0; n < define.Token_End; n++ {
		m.Tokens = append(m.Tokens, &Token{
			ID:      int32(n),
			Value:   0,
			MaxHold: 100000000,
			Entry:   entries.GetTokenEntry(int32(n)),
		})
	}
}

func (m *TokenManager) save() error {
	fields := map[string]interface{}{
		"tokens": m.Tokens,
	}
	store.GetStore().SaveFields(define.StoreType_Token, m, fields)

	return nil
}

func (m *TokenManager) LoadAll() error {
	err := store.GetStore().LoadObject(define.StoreType_Token, m.owner.GetID(), m)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("TokenManager LoadAll: %w", err)
	}

	store.GetStore().SaveObject(define.StoreType_Token, m)
	return nil
}

func (m *TokenManager) TokenInc(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	m.Tokens[tp].Value += value
	if m.Tokens[tp].Value > m.Tokens[tp].MaxHold {
		m.Tokens[tp].Value = m.Tokens[tp].MaxHold
	}

	m.save()
	m.SendTokenUpdate(m.Tokens[tp])
	return nil
}

func (m *TokenManager) TokenDec(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	m.Tokens[tp].Value -= value
	if m.Tokens[tp].Value < 0 {
		m.Tokens[tp].Value = 0
	}

	m.save()
	m.SendTokenUpdate(m.Tokens[tp])
	return nil
}

func (m *TokenManager) TokenSet(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	if value < 0 {
		return fmt.Errorf("token<%d> set invalid value<%d>", tp, value)
	}

	m.Tokens[tp].Value = value
	if m.Tokens[tp].Value > m.Tokens[tp].MaxHold {
		m.Tokens[tp].Value = m.Tokens[tp].MaxHold
	}

	m.save()
	m.SendTokenUpdate(m.Tokens[tp])
	return nil
}

func (m *TokenManager) GetToken(tp int32) (*Token, error) {
	if tp < 0 || tp >= define.Token_End {
		return nil, fmt.Errorf("token type<%d> invalid", tp)
	}

	return m.Tokens[tp], nil
}

func (m *TokenManager) SendTokenUpdate(t *Token) {
	msg := &pbGame.M2C_TokenUpdate{
		Info: &pbGame.Token{
			Type:    t.ID,
			Value:   t.Value,
			MaxHold: t.MaxHold,
		},
	}

	m.owner.SendProtoMessage(msg)
}
