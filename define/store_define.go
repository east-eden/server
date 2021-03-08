package define

const (
	StoreType_Begin = iota
	StoreType_User  = iota - 1
	StoreType_Account
	StoreType_PlayerInfo
	StoreType_Player
	StoreType_Item
	StoreType_Hero
	StoreType_Token
	StoreType_Fragment

	StoreType_End
)
