package config

// 拳皇之路-拳皇/KofRoadStarType
type KofRoadStarType int

const (
	KofRoadStarType_BEGIN    KofRoadStarType = 1 //枚举开始
	KofRoadStarType_Pass     KofRoadStarType = 1 //通关
	KofRoadStarType_Combo    KofRoadStarType = 2 //连击
	KofRoadStarType_RemainHp KofRoadStarType = 3 //剩余血量
	KofRoadStarType_END      KofRoadStarType = 4 //枚举结束
)

func (e KofRoadStarType) String() string {

	if e == KofRoadStarType_Pass {
		return "Pass"
	}
	if e == KofRoadStarType_Combo {
		return "Combo"
	}
	if e == KofRoadStarType_RemainHp {
		return "RemainHp"
	}

	return ""

}

func GetKofRoadStarType(tp string) (KofRoadStarType, bool) {
	switch tp {
	case "Pass":
		return KofRoadStarType_Pass, true
	case "Combo":
		return KofRoadStarType_Combo, true
	case "RemainHp":
		return KofRoadStarType_RemainHp, true

	}
	return -1, false

}

func GetKofRoadStarType_ByAlias(tp string) (KofRoadStarType, bool) {
	switch tp {
	case "通关":
		return KofRoadStarType_Pass, true
	case "连击":
		return KofRoadStarType_Combo, true
	case "剩余血量":
		return KofRoadStarType_RemainHp, true

	}
	return -1, false

}
