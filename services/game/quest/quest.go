package quest

import (
	"bitbucket.org/funplus/server/define"
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

// 任务目标监听事件类型
func GetQuestObjListenEvent(objType int32) int32 {
	switch objType {
	case define.QuestObj_Type_StagePass:
		return define.Event_Type_StagePass
	case define.QuestObj_Type_ChapterPass:
		return define.Event_Type_ChapterPass
	case define.QuestObj_Type_PlayerLevel:
		return define.Event_Type_PlayerLevelup
	case define.QuestObj_Type_HeroLevelTimes:
		return define.Event_Type_HeroLevelup
	case define.QuestObj_Type_GainHero, define.QuestObj_Type_GainHeroNum:
		return define.Event_Type_HeroAdd
	}

	return define.Event_Type_Null
}

func (q *Quest) IsComplete() bool {
	return q.State == define.Quest_State_Type_Completed || q.State == define.Quest_State_Type_Rewarded
}

func (q *Quest) Complete() {
	q.State = define.Quest_State_Type_Completed
}

func (q *Quest) CanComplete() bool {
	for _, obj := range q.Objs {
		if !obj.Completed {
			return false
		}
	}

	return true
}

func (q *Quest) IsRewarded() bool {
	return q.State == define.Quest_State_Type_Rewarded
}

func (q *Quest) Rewarded() {
	q.State = define.Quest_State_Type_Rewarded
}
