package token

import (
	"encoding/json"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
)

type Token struct {
	ID      int32 `json:"id"`
	Value   int64 `json:"value"`
	MaxHold int64 `json:"max_hold"`
	entry   *define.TokenEntry
}

type TokenManager struct {
	OwnerID   int64    `gorm:"type:bigint(20);primary_key;column:owner_id;index:owner_id;default:-1;not null"`
	OwnerType int32    `gorm:"type:int(10);primary_key;column:owner_type;index:owner_type;default:-1;not null"`
	TokenJson string   `gorm:"type:varchar(1024);column:token_json"`
	Tokens    []*Token `json:"tokens"`

	sync.RWMutex
	ds *db.Datastore
}

func NewTokenManager(owner define.PluginObj, ds *db.Datastore) *TokenManager {
	m := &TokenManager{
		OwnerID:   owner.GetID(),
		OwnerType: owner.GetType(),
		ds:        ds,
		Tokens:    make([]*Token, 0),
	}

	// init tokens
	m.initTokens()

	return m
}

func (m *TokenManager) TableName() string {
	return "token"
}

func Migrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(TokenManager{})
}

// interface of cost_loot
func (m *TokenManager) GetCostLootType() int32 {
	return define.CostLoot_Token
}

func (m *TokenManager) CanCost(typeMisc int32, num int32) error {
	return nil
}

func (m *TokenManager) DoCost(typeMisc int32, num int32) error {
	return nil
}

func (m *TokenManager) CanGain(typeMisc int32, num int32) error {
	return nil
}

func (m *TokenManager) GainLoot(typeMisc int32, num int32) error {
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
	m.ds.ORM().Find(&m)

	// unmarshal json to token value
	if len(m.TokenJson) > 0 {
		m.Lock()
		err := json.Unmarshal([]byte(m.TokenJson), &m.Tokens)
		if err != nil {
			logger.Error("unmarshal token json failed:", err)
		}
		m.Unlock()
	}
}

func (m *TokenManager) Save() error {
	m.RLock()
	data, err := json.Marshal(m.Tokens)
	m.RUnlock()
	if err != nil {
		return fmt.Errorf("json marshal failed:", err)
	}

	m.TokenJson = string(data)
	m.ds.ORM().Save(m)
	return nil
}

func (m *TokenManager) TokenInc(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type invalid:", tp)
	}

	m.Lock()
	m.Tokens[tp].Value += value
	m.Unlock()
	return nil
}

func (m *TokenManager) TokenDec(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type invalid:", tp)
	}

	m.Lock()
	m.Tokens[tp].Value -= value
	m.Unlock()
	return nil
}

func (m *TokenManager) TokenSet(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type invalid:", tp)
	}

	m.Lock()
	m.Tokens[tp].Value = value
	m.Unlock()
	return nil
}

func (m *TokenManager) GetToken(tp int32) (*Token, error) {
	if tp < 0 || tp >= define.Token_End {
		return nil, fmt.Errorf("token type invalid:", tp)
	}

	m.RLock()
	defer m.RUnlock()
	return m.Tokens[tp], nil
}
