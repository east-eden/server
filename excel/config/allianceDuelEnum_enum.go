package config

// 公会对决-拳皇/AllianceDuelEnum
type AllianceDuelProgress int

const (
	AllianceDuelProgress_BEGIN       AllianceDuelProgress = 1  //枚举开始
	AllianceDuelProgress_Matching    AllianceDuelProgress = 1  //匹配对手
	AllianceDuelProgress_TaskPeriod1 AllianceDuelProgress = 2  //任务阶段1
	AllianceDuelProgress_TaskPeriod2 AllianceDuelProgress = 3  //任务阶段2
	AllianceDuelProgress_TaskPeriod3 AllianceDuelProgress = 4  //任务阶段3
	AllianceDuelProgress_TaskPeriod4 AllianceDuelProgress = 5  //任务阶段4
	AllianceDuelProgress_TaskPeriod5 AllianceDuelProgress = 6  //任务阶段5
	AllianceDuelProgress_Fighting    AllianceDuelProgress = 7  //战斗阶段
	AllianceDuelProgress_Settlement  AllianceDuelProgress = 8  //领奖阶段
	AllianceDuelProgress_Closed      AllianceDuelProgress = 9  //关闭状态
	AllianceDuelProgress_END         AllianceDuelProgress = 10 //枚举结束
)

func (e AllianceDuelProgress) String() string {

	if e == AllianceDuelProgress_Matching {
		return "Matching"
	}
	if e == AllianceDuelProgress_TaskPeriod1 {
		return "TaskPeriod1"
	}
	if e == AllianceDuelProgress_TaskPeriod2 {
		return "TaskPeriod2"
	}
	if e == AllianceDuelProgress_TaskPeriod3 {
		return "TaskPeriod3"
	}
	if e == AllianceDuelProgress_TaskPeriod4 {
		return "TaskPeriod4"
	}
	if e == AllianceDuelProgress_TaskPeriod5 {
		return "TaskPeriod5"
	}
	if e == AllianceDuelProgress_Fighting {
		return "Fighting"
	}
	if e == AllianceDuelProgress_Settlement {
		return "Settlement"
	}
	if e == AllianceDuelProgress_Closed {
		return "Closed"
	}

	return ""

}

func GetAllianceDuelProgress(tp string) (AllianceDuelProgress, bool) {
	switch tp {
	case "Matching":
		return AllianceDuelProgress_Matching, true
	case "TaskPeriod1":
		return AllianceDuelProgress_TaskPeriod1, true
	case "TaskPeriod2":
		return AllianceDuelProgress_TaskPeriod2, true
	case "TaskPeriod3":
		return AllianceDuelProgress_TaskPeriod3, true
	case "TaskPeriod4":
		return AllianceDuelProgress_TaskPeriod4, true
	case "TaskPeriod5":
		return AllianceDuelProgress_TaskPeriod5, true
	case "Fighting":
		return AllianceDuelProgress_Fighting, true
	case "Settlement":
		return AllianceDuelProgress_Settlement, true
	case "Closed":
		return AllianceDuelProgress_Closed, true

	}
	return -1, false

}

func GetAllianceDuelProgress_ByAlias(tp string) (AllianceDuelProgress, bool) {
	switch tp {
	case "匹配对手":
		return AllianceDuelProgress_Matching, true
	case "任务阶段1":
		return AllianceDuelProgress_TaskPeriod1, true
	case "任务阶段2":
		return AllianceDuelProgress_TaskPeriod2, true
	case "任务阶段3":
		return AllianceDuelProgress_TaskPeriod3, true
	case "任务阶段4":
		return AllianceDuelProgress_TaskPeriod4, true
	case "任务阶段5":
		return AllianceDuelProgress_TaskPeriod5, true
	case "战斗阶段":
		return AllianceDuelProgress_Fighting, true
	case "领奖阶段":
		return AllianceDuelProgress_Settlement, true
	case "关闭状态":
		return AllianceDuelProgress_Closed, true

	}
	return -1, false

}
