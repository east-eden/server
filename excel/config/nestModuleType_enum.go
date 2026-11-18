package config

// 音巢科技表-拳皇/NestModuleType
type NestModuleType int

const (
	NestModuleType_BEGIN NestModuleType = 1 //枚举开始
	NestModuleType_Hand  NestModuleType = 1 //手
	NestModuleType_Chest NestModuleType = 2 //胸
	NestModuleType_Head  NestModuleType = 3 //头
	NestModuleType_Foot  NestModuleType = 4 //脚
	NestModuleType_END   NestModuleType = 5 //枚举结束
)

func (e NestModuleType) String() string {

	if e == NestModuleType_Hand {
		return "Hand"
	}
	if e == NestModuleType_Chest {
		return "Chest"
	}
	if e == NestModuleType_Head {
		return "Head"
	}
	if e == NestModuleType_Foot {
		return "Foot"
	}

	return ""

}

func GetNestModuleType(tp string) (NestModuleType, bool) {
	switch tp {
	case "Hand":
		return NestModuleType_Hand, true
	case "Chest":
		return NestModuleType_Chest, true
	case "Head":
		return NestModuleType_Head, true
	case "Foot":
		return NestModuleType_Foot, true

	}
	return -1, false

}

func GetNestModuleType_ByAlias(tp string) (NestModuleType, bool) {
	switch tp {
	case "手":
		return NestModuleType_Hand, true
	case "胸":
		return NestModuleType_Chest, true
	case "头":
		return NestModuleType_Head, true
	case "脚":
		return NestModuleType_Foot, true

	}
	return -1, false

}
