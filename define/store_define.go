package define

// StoreType 存储类型。
// 与 kof 对齐为具名类型，便于在接口签名与持久化配置里做类型约束。
type StoreType int

const (
	StoreType_Begin   StoreType = iota
	StoreType_Machine StoreType = iota - 1
	StoreType_User
	StoreType_Account
	StoreType_Player
	StoreType_Item
	StoreType_Hero
	StoreType_Collection
	StoreType_Fragment
	StoreType_Mail
	StoreType_Rank
	StoreType_Comment
	StoreType_GlobalMess

	StoreType_End
)

// String 返回存储类型对应的表名
func (s StoreType) String() string {
	switch s {
	case StoreType_Machine:
		return "machine"
	case StoreType_User:
		return "user"
	case StoreType_Account:
		return "account"
	case StoreType_Player:
		return "player"
	case StoreType_Item:
		return "player_item"
	case StoreType_Hero:
		return "player_hero"
	case StoreType_Collection:
		return "player_collection"
	case StoreType_Fragment:
		return "player_fragment"
	case StoreType_Mail:
		return "mail"
	case StoreType_Rank:
		return "rank"
	case StoreType_Comment:
		return "comment"
	case StoreType_GlobalMess:
		return "global_mess"
	default:
		return ""
	}
}
