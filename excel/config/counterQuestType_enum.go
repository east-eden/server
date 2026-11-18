package config

// 计数任务表-拳皇/CounterQuestType
type QuestCounterType int

const (
	QuestCounterType_BEGIN                QuestCounterType = 1 //枚举开始
	QuestCounterType_ActivityCounterQuest QuestCounterType = 1 //活动计数任务类型
	QuestCounterType_END                  QuestCounterType = 2 //枚举结束
)

func (e QuestCounterType) String() string {

	if e == QuestCounterType_ActivityCounterQuest {
		return "ActivityCounterQuest"
	}

	return ""

}

func GetQuestCounterType(tp string) (QuestCounterType, bool) {
	switch tp {
	case "ActivityCounterQuest":
		return QuestCounterType_ActivityCounterQuest, true

	}
	return -1, false

}

func GetQuestCounterType_ByAlias(tp string) (QuestCounterType, bool) {
	switch tp {
	case "活动计数任务类型":
		return QuestCounterType_ActivityCounterQuest, true

	}
	return -1, false

}
