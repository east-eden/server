package token

import (
	"encoding/json"
	"fmt"

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
	OwnerID   int64    `gorm:"type:bigint(20);primary_key;column:owner_id;index:owner_id;default:0;not null"`
	TokenJson string   `gorm:"type:varchar(1024);column:token_json"`
	Tokens    []*Token `json:"tokens"`

	ds *db.Datastore
}

func NewTokenManager(ownerID int64, ds *db.Datastore) *TokenManager {
	m := &TokenManager{
		OwnerID: ownerID,
		ds:      ds,
		Tokens:  make([]*Token, 0),
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

func (m *TokenManager) initTokens() {
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
		err := json.Unmarshal([]byte(m.TokenJson), &m.Tokens)
		if err != nil {
			logger.Error("unmarshal token json failed:", err)
		}
	}
}

func (m *TokenManager) Save() error {
	data, err := json.Marshal(m.Tokens)
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

	m.Tokens[tp].Value += value
	return nil
}

func (m *TokenManager) TokenDec(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type invalid:", tp)
	}

	m.Tokens[tp].Value -= value
	return nil
}

func (m *TokenManager) TokenSet(tp int32, value int64) error {
	if tp < 0 || tp >= define.Token_End {
		return fmt.Errorf("token type invalid:", tp)
	}

	m.Tokens[tp].Value = value
	return nil
}

func (m *TokenManager) GetToken(tp int32) (*Token, error) {
	if tp < 0 || tp >= define.Token_End {
		return nil, fmt.Errorf("token type invalid:", tp)
	}

	return m.Tokens[tp], nil
}
