package quest

import (
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
)

// SubCondition 任务子条件（kof SubCondition.json 的运行时形态）
type SubCondition struct {
	Id        int32                              `bson:"id" json:"id"`               // 子条件id
	Type      config.SubConditionType            `bson:"type" json:"type"`           // 子条件类型
	Params    []int32                            `bson:"params" json:"params"`       // 子条件参数
	OpTypes   []config.SubConditionOperationType `bson:"opTypes" json:"opTypes"`     // 比较方式
	Count     int32                              `bson:"count" json:"count"`         // 当前计数
	MaxCount  int32                              `bson:"maxCount" json:"maxCount"`   // 目标计数
	Completed bool                               `bson:"completed" json:"completed"` // 是否达成
}

// GenPB 生成任务目标协议结构
func (sc *SubCondition) GenPB() *pbGlobal.QuestObj {
	return &pbGlobal.QuestObj{
		Type:      int32(sc.Type),
		Count:     sc.Count,
		Completed: sc.Completed,
	}
}

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

// buildSubConditions 按主条件展开子条件。
// 主条件里复合条件（Condition）与直接子条件（SubCondition）两种形态都存在，
// 这里统一展开成子条件列表，任一缺失就返回空列表由调用方判定为“未配置”。
func buildSubConditions(subConditionId int32) []*SubCondition {
	subs := make([]*SubCondition, 0, 2)

	// 先按复合条件展开
	if cond, ok := config.GetConditionEntry(subConditionId); ok {
		for _, sid := range cond.SubCondIds {
			if sc, ok := newSubCondition(sid); ok {
				subs = append(subs, sc)
			}
		}

		return subs
	}

	// 再按直接子条件处理
	if sc, ok := newSubCondition(subConditionId); ok {
		subs = append(subs, sc)
	}

	return subs
}

func newSubCondition(id int32) (*SubCondition, bool) {
	entry, ok := config.GetSubConditionEntry(id)
	if !ok {
		return nil, false
	}

	return &SubCondition{
		Id:      entry.Id,
		Type:    entry.SubCondType,
		Params:  entry.SubCondParams,
		OpTypes: entry.OperationType,
	}, true
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
