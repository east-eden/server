package game

import (
	"encoding"
)

type User struct {
	UserID                     int64  `bson:"_id" json:"_id"`
	AccountID                  int64  `bson:"account_id" json:"account_id"`
	PlayerID                   int64  `bson:"player_id" json:"player_id"`
	PlayerName                 string `bson:"player_name" json:"player_name"`
	PlayerLevel                int32  `bson:"player_level" json:"player_level"`
	encoding.BinaryMarshaler   `bson:"-" json:"-"`
	encoding.BinaryUnmarshaler `bson:"-" json:"-"`
}

func NewUser() interface{} {
	return &User{}
}

func (u *User) Init() {
	u.UserID = -1
	u.AccountID = -1
	u.PlayerID = 1
	u.PlayerName = ""
	u.PlayerLevel = 1
}
