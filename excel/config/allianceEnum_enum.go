package config

// 联盟-拳皇/AllianceEnum
type AlliancePositionType int

const (
	AlliancePositionType_BEGIN  AlliancePositionType = 0 //枚举开始
	AlliancePositionType_Normal AlliancePositionType = 0 //成员
	AlliancePositionType_Elite  AlliancePositionType = 1 //精英
	AlliancePositionType_Deputy AlliancePositionType = 2 //副会长
	AlliancePositionType_Leader AlliancePositionType = 3 //会长
	AlliancePositionType_END    AlliancePositionType = 4 //枚举结束
)

// 联盟-拳皇/AllianceEnum
type AllianceTechType int

const (
	AllianceTechType_BEGIN AllianceTechType = 0 //枚举开始
	AllianceTechType_Empty AllianceTechType = 0 //暂无
	AllianceTechType_Chip  AllianceTechType = 1 //碎片
	AllianceTechType_Times AllianceTechType = 2 //次数
	AllianceTechType_Goods AllianceTechType = 3 //道具
	AllianceTechType_END   AllianceTechType = 4 //枚举结束
)

func (e AlliancePositionType) String() string {

	if e == AlliancePositionType_Normal {
		return "Normal"
	}
	if e == AlliancePositionType_Elite {
		return "Elite"
	}
	if e == AlliancePositionType_Deputy {
		return "Deputy"
	}
	if e == AlliancePositionType_Leader {
		return "Leader"
	}

	return ""

}

func GetAlliancePositionType(tp string) (AlliancePositionType, bool) {
	switch tp {
	case "Normal":
		return AlliancePositionType_Normal, true
	case "Elite":
		return AlliancePositionType_Elite, true
	case "Deputy":
		return AlliancePositionType_Deputy, true
	case "Leader":
		return AlliancePositionType_Leader, true

	}
	return -1, false

}

func GetAlliancePositionType_ByAlias(tp string) (AlliancePositionType, bool) {
	switch tp {
	case "成员":
		return AlliancePositionType_Normal, true
	case "精英":
		return AlliancePositionType_Elite, true
	case "副会长":
		return AlliancePositionType_Deputy, true
	case "会长":
		return AlliancePositionType_Leader, true

	}
	return -1, false

}

func (e AllianceTechType) String() string {

	if e == AllianceTechType_Empty {
		return "Empty"
	}
	if e == AllianceTechType_Chip {
		return "Chip"
	}
	if e == AllianceTechType_Times {
		return "Times"
	}
	if e == AllianceTechType_Goods {
		return "Goods"
	}

	return ""

}

func GetAllianceTechType(tp string) (AllianceTechType, bool) {
	switch tp {
	case "Empty":
		return AllianceTechType_Empty, true
	case "Chip":
		return AllianceTechType_Chip, true
	case "Times":
		return AllianceTechType_Times, true
	case "Goods":
		return AllianceTechType_Goods, true

	}
	return -1, false

}

func GetAllianceTechType_ByAlias(tp string) (AllianceTechType, bool) {
	switch tp {
	case "暂无":
		return AllianceTechType_Empty, true
	case "碎片":
		return AllianceTechType_Chip, true
	case "次数":
		return AllianceTechType_Times, true
	case "道具":
		return AllianceTechType_Goods, true

	}
	return -1, false

}
