package quest

import (
	"github.com/east-eden/server/define"
	pbGlobal "github.com/east-eden/server/proto/global"
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

func (obj *QuestObj) GenPB() *pbGlobal.QuestObj {
	msg := &pbGlobal.QuestObj{
		Type:      obj.Type,
		Count:     obj.Count,
		Completed: obj.Completed,
	}

	return msg
}

// 任务
type Quest struct {
	Options `bson:"inline" json:",inline"`
}

func NewQuest(opts ...Option) *Quest {
	q := &Quest{
		Options: DefaultOptions(),
	}

	for _, o := range opts {
		o(&q.Options)
	}

	return q
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
		return define.Event_Type_HeroGain
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

func (q *Quest) CanReward() bool {
	return q.IsComplete() && !q.IsRewarded()
}

func (q *Quest) Refresh() {
	q.State = define.Quest_State_Type_Opened
	for _, obj := range q.Objs {
		obj.Count = 0
		obj.Completed = false
	}
}

func (q *Quest) GenPB() *pbGlobal.Quest {
	pb := &pbGlobal.Quest{
		Id:    q.QuestId,
		State: q.State,
		Objs:  make([]*pbGlobal.QuestObj, 0, len(q.Objs)),
	}

	for _, obj := range q.Objs {
		pbObj := obj.GenPB()
		pb.Objs = append(pb.Objs, pbObj)
	}

	return pb
}
