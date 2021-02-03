package define

const (
	StoreType_Begin = iota
	StoreType_User  = iota - 1
	StoreType_Account
	StoreType_PlayerInfo
	StoreType_Player
	StoreType_Item
	StoreType_Hero
	StoreType_Blade
	StoreType_Token
	StoreType_Rune
	StoreType_Talent

	StoreType_End
)
