package gate

import (
	"time"

	"github.com/yokaiio/yokai_server/store"
)

var userExpireTime time.Duration = 30 * time.Minute

type UserInfo struct {
	store.StoreObjector `bson:"-" json:"-"`

	UserID      int64  `bson:"_id" json:"_id"`
	AccountID   int64  `bson:"account_id" json:"account_id"`
	GameID      int16  `bson:"game_id" json:"game_id"`
	PlayerID    int64  `bson:"player_id" json:"player_id"`
	PlayerName  string `bson:"player_name" json:"player_name"`
	PlayerLevel int32  `bson:"player_level" json:"player_level"`
}

func (u *UserInfo) TableName() string {
	return "user"
}

func (u *UserInfo) GetObjID() interface{} {
	return u.UserID
}

func (u *UserInfo) AfterLoad() error {
	return nil
}

func NewUserInfo() interface{} {
	return &UserInfo{
		UserID:      -1,
		AccountID:   -1,
		GameID:      int16(-1),
		PlayerID:    -1,
		PlayerName:  "",
		PlayerLevel: 1,
	}
}
