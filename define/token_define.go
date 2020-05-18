package define

const (
	Token_Gold = iota
	Token_Diamond
	Token_Honour

	Token_End
)

// token entry
type TokenEntry struct {
	ID      int32  `json:"_id"`
	Name    string `json:"name"`
	MaxHold int64  `json:"max_hold"`
}
