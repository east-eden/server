package quest

import (
	"bitbucket.org/funplus/server/excel/auto"
)

var (
	Quest_Max_Obj = 5 // 任务最多目标数
)

// 任务目标
type QuestObj struct {
	Type      int32 `bson:"type" json:"type"`           // 任务目标类型
	Count     int32 `bson:"count" json:"count"`         // 任务目标计数
	Completed bool  `bson:"completed" json:"completed"` // 任务目标是否达成
}

// 任务
type Quest struct {
	Options `bson:"inline" json:",inline"`
	Entry   *auto.QuestEntry `bson:"-" json:"-"`
}

func NewQuest(opts ...Option) *Quest {
	return &Quest{
		Options: DefaultOptions(),
		Entry:   nil,
	}
}
