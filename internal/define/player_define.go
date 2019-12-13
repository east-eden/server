package define

const (
	Player_MaxPlayer = 10 // how many players can be created per client
	Player_MaxLevel  = 60
)

// player level up entry
type PlayerLevelupEntry struct {
	ID  int32 `json:"id"`
	Exp int64 `json:"exp"`
}
