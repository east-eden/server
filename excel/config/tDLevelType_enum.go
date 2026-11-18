package config

// 副玩法关卡表-拳皇/TDLevelType
type TDLevelGeneratorType int

const (
	TDLevelGeneratorType_BEGIN   TDLevelGeneratorType = 1 //枚举开始
	TDLevelGeneratorType_Random  TDLevelGeneratorType = 1 //随机
	TDLevelGeneratorType_Uniform TDLevelGeneratorType = 2 //均匀
	TDLevelGeneratorType_END     TDLevelGeneratorType = 3 //枚举结束
)

// 副玩法关卡表-拳皇/TDLevelType
type RougeSkillTargetType int

const (
	RougeSkillTargetType_BEGIN  RougeSkillTargetType = 1 //枚举开始
	RougeSkillTargetType_Wall   RougeSkillTargetType = 1 //城墙
	RougeSkillTargetType_Player RougeSkillTargetType = 2 //我方阵营
	RougeSkillTargetType_Enemy  RougeSkillTargetType = 3 //敌方阵营
	RougeSkillTargetType_END    RougeSkillTargetType = 4 //枚举结束
)

func (e TDLevelGeneratorType) String() string {

	if e == TDLevelGeneratorType_Random {
		return "Random"
	}
	if e == TDLevelGeneratorType_Uniform {
		return "Uniform"
	}

	return ""

}

func GetTDLevelGeneratorType(tp string) (TDLevelGeneratorType, bool) {
	switch tp {
	case "Random":
		return TDLevelGeneratorType_Random, true
	case "Uniform":
		return TDLevelGeneratorType_Uniform, true

	}
	return -1, false

}

func GetTDLevelGeneratorType_ByAlias(tp string) (TDLevelGeneratorType, bool) {
	switch tp {
	case "随机":
		return TDLevelGeneratorType_Random, true
	case "均匀":
		return TDLevelGeneratorType_Uniform, true

	}
	return -1, false

}

func (e RougeSkillTargetType) String() string {

	if e == RougeSkillTargetType_Wall {
		return "Wall"
	}
	if e == RougeSkillTargetType_Player {
		return "Player"
	}
	if e == RougeSkillTargetType_Enemy {
		return "Enemy"
	}

	return ""

}

func GetRougeSkillTargetType(tp string) (RougeSkillTargetType, bool) {
	switch tp {
	case "Wall":
		return RougeSkillTargetType_Wall, true
	case "Player":
		return RougeSkillTargetType_Player, true
	case "Enemy":
		return RougeSkillTargetType_Enemy, true

	}
	return -1, false

}

func GetRougeSkillTargetType_ByAlias(tp string) (RougeSkillTargetType, bool) {
	switch tp {
	case "城墙":
		return RougeSkillTargetType_Wall, true
	case "我方阵营":
		return RougeSkillTargetType_Player, true
	case "敌方阵营":
		return RougeSkillTargetType_Enemy, true

	}
	return -1, false

}
