package quest

// 任务目标
type QuestObj struct {
	Type  int32 `bson:"type" json:"type"`   // 任务目标类型
	Count int32 `bson:"count" json:"count"` // 任务目标计数
	State int32 `bson:"state" json:"state"` // 任务目标状态
}

// 任务
type Quest struct {
	Id      int32       `bson:"_id" json:"_id"`           // 任务id
	OwnerId int64       `bson:"owner_id" json:"owner_id"` // 玩家id
	Objs    []*QuestObj `bson:"objs" json:"objs"`         // 任务目标数据
	State   int32       `bson:"state" json:"state"`       // 任务状态

	Entry *auto.QuestEntry `bson:"-" json:"-"`
}
