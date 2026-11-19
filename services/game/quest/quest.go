package quest

import (
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/subcondition"
)

// SubCondition 任务子条件（共享自 subcondition 包）
type SubCondition = subcondition.SubCondition

// Options 任务数据
type Options struct {
	Id             int32                  `bson:"id" json:"id"`                         // 任务id
	SubConditionId int32                  `bson:"subConditionId" json:"subConditionId"` // 主条件id
	Entry          *config.MainQuestEntry `bson:"-" json:"-"`                           // 任务配置
	SubConditions  []*SubCondition        `bson:"subConditions" json:"subConditions"`   // 子条件数据
	State          config.QuestState      `bson:"state" json:"state"`                   // 任务状态
}

func DefaultOptions() Options {
	return Options{
		Id:             -1,
		SubConditionId: -1,
		SubConditions:  make([]*SubCondition, 0, 2),
		State:          config.QuestState_Open,
	}
}

type Option func(*Options)

func WithId(id int32) Option {
	return func(o *Options) {
		o.Id = id
	}
}

func WithConfig(entry *config.MainQuestEntry) Option {
	return func(o *Options) {
		if entry == nil {
			return
		}

		o.Id = entry.Id
		o.SubConditionId = entry.SubConditionId
		o.Entry = entry
		o.SubConditions = buildSubConditions(entry.SubConditionId)
	}
}

func WithState(state config.QuestState) Option {
	return func(o *Options) {
		o.State = state
	}
}

// buildSubConditions 按主条件展开子条件（共享实现）
func buildSubConditions(subConditionId int32) []*SubCondition {
	return subcondition.BuildFromConditionId(subConditionId)
}

// Quest 任务
type Quest struct {
	Options `bson:"inline" json:",inline"`
}

func NewQuest(opts ...Option) *Quest {
	q := &Quest{Options: DefaultOptions()}
	for _, o := range opts {
		o(&q.Options)
	}

	return q
}

func (q *Quest) IsValid() bool {
	return q.Id != -1
}

func (q *Quest) IsComplete() bool {
	return q.State == config.QuestState_Completed
}

func (q *Quest) Complete() {
	q.State = config.QuestState_Completed
}

func (q *Quest) IsRewarded() bool {
	return q.State == config.QuestState_Rewarded
}

func (q *Quest) Reward() {
	q.State = config.QuestState_Rewarded
}

func (q *Quest) IsOpen() bool {
	return q.State == config.QuestState_Open
}

// CanComplete 子条件是否全部达成
func (q *Quest) CanComplete() bool {
	if q.IsRewarded() || q.IsComplete() {
		return q.IsComplete()
	}

	if len(q.SubConditions) == 0 {
		return false
	}

	relation := config.SubConditionRelation_And
	if cond, ok := config.GetConditionEntry(q.SubConditionId); ok {
		relation = cond.SubCondRelation
	}

	if relation == config.SubConditionRelation_Or {
		for _, sc := range q.SubConditions {
			if sc.Completed {
				return true
			}
		}

		return false
	}

	for _, sc := range q.SubConditions {
		if !sc.Completed {
			return false
		}
	}

	return true
}

// CanReward 已完成且未领奖
func (q *Quest) CanReward() bool {
	return q.IsComplete() && !q.IsRewarded()
}

// GenPB 生成下发协议
func (q *Quest) GenPB() *pbGlobal.Quest {
	msg := &pbGlobal.Quest{
		Id:    q.Id,
		State: int32(q.State),
		Objs:  make([]*pbGlobal.QuestObj, 0, len(q.SubConditions)),
	}

	for _, sc := range q.SubConditions {
		msg.Objs = append(msg.Objs, sc.GenPB())
	}

	return msg
}
