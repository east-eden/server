package define

import "time"

const (
	Player_MaxPlayer = 10 // how many players can be created per client
	Player_MaxLevel  = 60
	Player_MemExpire = 2 * time.Hour // memory expire time
)

// player level up entry
type PlayerLevelupEntry struct {
	ID  int32 `json:"id"`
	Exp int64 `json:"exp"`
}
