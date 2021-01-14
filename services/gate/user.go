package gate

import (
	"e.coding.net/mmstudio/blade/server/store"
)

type UserInfo struct {
	store.StoreObjector `bson:"-" json:"-"`

	UserID      int64  `bson:"_id" json:"_id"`
	AccountID   int64  `bson:"account_id" json:"account_id"`
	PlayerID    int64  `bson:"player_id" json:"player_id"`
	PlayerName  string `bson:"player_name" json:"player_name"`
	PlayerLevel int32  `bson:"player_level" json:"player_level"`
}

func (u *UserInfo) GetObjID() int64 {
	return u.UserID
}

func (u *UserInfo) GetStoreIndex() int64 {
	return -1
}

func (u *UserInfo) AfterLoad() error {
	return nil
}

func NewUserInfo() interface{} {
	return &UserInfo{
		UserID:      -1,
		AccountID:   -1,
		PlayerID:    -1,
		PlayerName:  "",
		PlayerLevel: 1,
	}
}
