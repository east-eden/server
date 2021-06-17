package define

const (
	StoreType_Begin   = iota
	StoreType_Machine = iota - 1
	StoreType_User
	StoreType_Account
	StoreType_Player
	StoreType_Item
	StoreType_Hero
	StoreType_Collection
	StoreType_Token
	StoreType_Fragment
	StoreType_Mail

	StoreType_GlobalMess

	StoreType_End
)
