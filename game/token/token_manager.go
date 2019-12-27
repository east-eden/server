package token

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Token struct {
	ID      int32              `json:"id" bson:"token_id"`
	Value   int32              `json:"value" bson:"token_value"`
	MaxHold int32              `json:"max_hold" bson:"token_maxhold"`
	entry   *define.TokenEntry `json:"-" bson:"-"`
}

type TokenManager struct {
	OwnerID   int64    `gorm:"type:bigint(20);primary_key;column:owner_id;index:owner_id;default:-1;not null" bson:"_id"`
	OwnerType int32    `gorm:"type:int(10);column:owner_type;index:owner_type;default:-1;not null" bson:"owner_type"`
	Tokens    []*Token `json:"tokens" bson:"tokens"`

	sync.RWMutex `bson:"-"`
	ds           *db.Datastore     `bson:"-"`
	coll         *mongo.Collection `bson:"-"`
}

func NewTokenManager(owner define.PluginObj, ds *db.Datastore) *TokenManager {
	m := &TokenManager{
		OwnerID:   owner.GetID(),
		OwnerType: owner.GetType(),
		ds:        ds,
		Tokens:    make([]*Token, 0),
	}

	m.coll = ds.Database().Collection(m.TableName())

	// init tokens
	m.initTokens()

	return m
}

func (m *TokenManager) TableName() string {
	return "token"
}

func Migrate(ds *db.Datastore) {
	//ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(TokenManager{})
}

// interface of cost_loot
func (m *TokenManager) GetCostLootType() int32 {
	return define.CostLoot_Token
}

func (m *TokenManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("token manager check token<%d> cost failed, wrong number<%d>", typeMisc, num)
	}

	for _, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.Value >= num {
				return nil
			}
		}
	}

	return fmt.Errorf("not enough token<%d>, num<%d>", typeMisc, num)
}

func (m *TokenManager) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("token manager cost token<%d> failed, wrong number<%d>", typeMisc, num)
	}

	for _, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.Value < num {
				logger.WithFields(logger.Fields{
					"cost_type_misc":  typeMisc,
					"cost_num":        num,
					"actual_cost_num": v.Value,
				}).Warn("token manager cost number error")
			}

			v.Value -= num
			if v.Value < 0 {
				v.Value = 0
			}

			m.Save()
		}
	}

	return nil
}

func (m *TokenManager) CanGain(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return nil
}

func (m *TokenManager) GainLoot(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("token manager check gain token<%d> failed, wrong number<%d>", typeMisc, num)
	}

	for _, v := range m.Tokens {
		if v.ID == typeMisc {
			if v.MaxHold < v.Value+num {
				logger.WithFields(logger.Fields{
					"gain_type_misc":  typeMisc,
					"gain_num":        num,
					"actual_gain_num": v.MaxHold - v.Value,
				}).Warn("token manager gain number overflow")
			}

			v.Value += num
			if v.Value > v.MaxHold {
				v.Value = v.MaxHold
			}

			m.Save()
		}
	}

	return nil
}

func (m *TokenManager) initTokens() {
	m.Lock()
	defer m.Unlock()
	for n := 0; n < define.Token_End; n++ {
		m.Tokens = append(m.Tokens, &Token{
			ID:      int32(n),
			Value:   0,
			MaxHold: 100000000,
			entry:   global.GetTokenEntry(int32(n)),
		})
	}
}

func (m *TokenManager) LoadFromDB() {
	res := m.coll.FindOne(context.Background(), bson.D{{"_id": m.OwnerID}})
	if res.Err() == mongo.ErrNoDocuments {
		m.coll.InsertOne(context.Background(), m)
	} else {
		res.Decode(m)
	}
}

func (m *TokenManager) Save() error {
	m.RLock()
	data, err := json.Marshal(m.Tokens)
	m.RUnlock()
	if err != nil {
		return fmt.Errorf("json marshal failed:%s", err.Error())
	}

	filter := bson.D{{"_id", m.OwnerID}}
	update := bson.D{{"$set", m}}
	op := options.Update().SetUpsert(true)
	m.coll.UpdateOne(context.Background(), filter, update, op)
	return nil
}

func (m *TokenManager) TokenInc(tp int32, value int32) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	m.Tokens[tp].Value += value
	if m.Tokens[tp].Value > m.Tokens[tp].MaxHold {
		m.Tokens[tp].Value = m.Tokens[tp].MaxHold
	}

	return nil
}

func (m *TokenManager) TokenDec(tp int32, value int32) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type<%d> invalid", tp)
	}

	m.Tokens[tp].Value -= value
	if m.Tokens[tp].Value < 0 {
		m.Tokens[tp].Value = 0
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

	m.Tokens[tp].Value = value
	if m.Tokens[tp].Value > m.Tokens[tp].MaxHold {
		m.Tokens[tp].Value = m.Tokens[tp].MaxHold
	}

	return nil
}

func (m *TokenManager) GetToken(tp int32) (*Token, error) {
	if tp < 0 || tp >= define.Token_End {
		return nil, fmt.Errorf("token type<%d> invalid", tp)
	}

	return m.Tokens[tp], nil
}
