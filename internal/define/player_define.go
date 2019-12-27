package define

import "time"

const (
	Player_MaxPlayer     = 10 // how many players can be created per client
	Player_MaxLevel      = 60
	Player_MemExpire     = 60 * time.Second // memory expire time
	Player_ExpireChanNum = 1000             // player expire channel num
)

// player level up entry
type PlayerLevelupEntry struct {
	ID  int32 `json:"id"`
	Exp int64 `json:"exp"`
}
