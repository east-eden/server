package define

const (
	Event_Type_Begin         int32 = iota
	Event_Type_Sign          int32 = iota - 1 // 0 签到
	Event_Type_PlayerLevelup                  // 1 角色升级
	Event_Type_HeroLevelup                    // 2 英雄升级
	Event_Type_End
)

// 事件属性
type Event struct {
	Type  int32         `bson:"type" json:"type"`   // 事件类型
	Miscs []interface{} `bson:"miscs" json:"miscs"` // 事件参数
}
