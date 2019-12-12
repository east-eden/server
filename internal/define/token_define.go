package define

const (
	Token_Gold = iota
	Token_Diamond
	Token_Honour

	Token_End
)

// token entry
type TokenEntry struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	MaxHold int32  `json:"max_hold"`
}
