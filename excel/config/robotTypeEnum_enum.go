package config

// 竞技场-拳皇/RobotTypeEnum
type RobotType int

const (
	RobotType_BEGIN    RobotType = 1 //枚举开始
	RobotType_Arena    RobotType = 1 //竞技场
	RobotType_Fighting RobotType = 2 //格斗圣殿
	RobotType_Arcade   RobotType = 3 //街机厅
	RobotType_END      RobotType = 4 //枚举结束
)

func (e RobotType) String() string {

	if e == RobotType_Arena {
		return "Arena"
	}
	if e == RobotType_Fighting {
		return "Fighting"
	}
	if e == RobotType_Arcade {
		return "Arcade"
	}

	return ""

}

func GetRobotType(tp string) (RobotType, bool) {
	switch tp {
	case "Arena":
		return RobotType_Arena, true
	case "Fighting":
		return RobotType_Fighting, true
	case "Arcade":
		return RobotType_Arcade, true

	}
	return -1, false

}

func GetRobotType_ByAlias(tp string) (RobotType, bool) {
	switch tp {
	case "竞技场":
		return RobotType_Arena, true
	case "格斗圣殿":
		return RobotType_Fighting, true
	case "街机厅":
		return RobotType_Arcade, true

	}
	return -1, false

}
