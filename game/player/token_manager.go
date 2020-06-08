package player

import (
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/store"
)

type Token struct {
	ID      int32              `json:"token_id" bson:"token_id" redis:"token_id"`
	Value   int64              `json:"token_value" bson:"token_value" redis:"token_value"`
	MaxHold int64              `json:"token_max_hold" bson:"token_max_hold" redis:"token_max_hold"`
	entry   *define.TokenEntry `json:"-" bson:"-" redis:"-"`
}

type TokenManager struct {
	store.StoreObjector `bson:"-" redis:"-"`
	owner               *Player  `bson:"-" redis:"-"`
	OwnerId             int64    `bson:"_id" redis:"_id"`
	OwnerType           int32    `bson:"owner_type" redis:"owner_type"`
	Tokens              []*Token `json:"tokens" bson:"tokens" redis:"tokens"`

	sync.RWMutex `bson:"-" redis:"-"`
}

func NewTokenManager(owner *Player) *TokenManager {
	m := &TokenManager{
		owner:     owner,
		OwnerId:   owner.GetID(),
		OwnerType: owner.GetType(),
		Tokens:    make([]*Token, 0),
	}

	// init tokens
	m.initTokens()

	return m
}

func (m *TokenManager) TableName() string {
	return "token"
}

func (m *TokenManager) AfterLoad() {

}

func (m *TokenManager) GetObjID() interface{} {
	return m.owner.GetID()
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

	for tp, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.Value < costNum {
				logger.WithFields(logger.Fields{
					"cost_type_misc":  typeMisc,
					"cost_num":        costNum,
					"actual_cost_num": v.Value,
				}).Warn("token manager cost number error")
			}

			v.Value -= costNum
			if v.Value < 0 {
				v.Value = 0
			}

			m.save(int32(tp))
		}
	}

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

	for tp, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.MaxHold < v.Value+gainNum {
				logger.WithFields(logger.Fields{
					"gain_type_misc":  typeMisc,
					"gain_num":        gainNum,
					"actual_gain_num": v.MaxHold - v.Value,
				}).Warn("token manager gain number overflow")
			}

			v.Value += gainNum
			if v.Value > v.MaxHold {
				v.Value = v.MaxHold
			}

			m.save(int32(tp))
		}
	}

	return nil
}

func (m *TokenManager) initTokens() {
	for n := 0; n < define.Token_End; n++ {
		m.Tokens = append(m.Tokens, &Token{
			ID:      int32(n),
			Value:   0,
			MaxHold: 100000000,
			entry:   entries.GetTokenEntry(int32(n)),
		})
	}
}

func (m *TokenManager) save(tp int32) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid when save", tp)
	}

	fields := map[string]interface{}{
		"tokens": m.Tokens,
	}
	m.owner.store.SaveFieldsToCacheAndDB(store.StoreType_Token, m, fields)

	return nil
}

func (m *TokenManager) LoadAll() {
	err := m.owner.store.LoadObjectFromCacheAndDB(store.StoreType_Token, "_id", m.owner.GetID(), m)
	if err != nil {
		m.owner.store.SaveObjectToCacheAndDB(store.StoreType_Token, m)
	}
}

func (m *TokenManager) TokenInc(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	m.Tokens[tp].Value += value
	if m.Tokens[tp].Value > m.Tokens[tp].MaxHold {
		m.Tokens[tp].Value = m.Tokens[tp].MaxHold
	}

	m.save(tp)
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

	m.save(tp)
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

	m.save(tp)
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
