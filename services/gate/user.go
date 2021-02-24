package gate

type UserInfo struct {
	UserID      int64  `bson:"_id" json:"_id"`
	AccountID   int64  `bson:"account_id" json:"account_id"`
	PlayerID    int64  `bson:"player_id" json:"player_id"`
	PlayerName  string `bson:"player_name" json:"player_name"`
	PlayerLevel int32  `bson:"player_level" json:"player_level"`
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
