package config

// 基础配置-拳皇/BackpackItemType
type BackpackTab int

const (
	BackpackTab_BEGIN    BackpackTab = 0   //枚举开始
	BackpackTab_Prop     BackpackTab = 0   //道具
	BackpackTab_Material BackpackTab = 1   //碎片
	BackpackTab_HeroSoul BackpackTab = 2   //奥义石
	BackpackTab_Hide     BackpackTab = 99  //不显示
	BackpackTab_END      BackpackTab = 100 //枚举结束
)

// 基础配置-拳皇/BackpackItemType
type BackpackUseType int

const (
	BackpackUseType_BEGIN     BackpackUseType = 0 //枚举开始
	BackpackUseType_NonUse    BackpackUseType = 0 //不可使用
	BackpackUseType_Optional  BackpackUseType = 1 //自选
	BackpackUseType_Normal    BackpackUseType = 3 //直接使用
	BackpackUseType_Compose   BackpackUseType = 4 //合成
	BackpackUseType_RedPacket BackpackUseType = 5 //红包
	BackpackUseType_END       BackpackUseType = 6 //枚举结束
)

func (e BackpackTab) String() string {

	if e == BackpackTab_Prop {
		return "Prop"
	}
	if e == BackpackTab_Material {
		return "Material"
	}
	if e == BackpackTab_HeroSoul {
		return "HeroSoul"
	}
	if e == BackpackTab_Hide {
		return "Hide"
	}

	return ""

}

func GetBackpackTab(tp string) (BackpackTab, bool) {
	switch tp {
	case "Prop":
		return BackpackTab_Prop, true
	case "Material":
		return BackpackTab_Material, true
	case "HeroSoul":
		return BackpackTab_HeroSoul, true
	case "Hide":
		return BackpackTab_Hide, true

	}
	return -1, false

}

func GetBackpackTab_ByAlias(tp string) (BackpackTab, bool) {
	switch tp {
	case "道具":
		return BackpackTab_Prop, true
	case "碎片":
		return BackpackTab_Material, true
	case "奥义石":
		return BackpackTab_HeroSoul, true
	case "不显示":
		return BackpackTab_Hide, true

	}
	return -1, false

}

func (e BackpackUseType) String() string {

	if e == BackpackUseType_NonUse {
		return "NonUse"
	}
	if e == BackpackUseType_Optional {
		return "Optional"
	}
	if e == BackpackUseType_Normal {
		return "Normal"
	}
	if e == BackpackUseType_Compose {
		return "Compose"
	}
	if e == BackpackUseType_RedPacket {
		return "RedPacket"
	}

	return ""

}

func GetBackpackUseType(tp string) (BackpackUseType, bool) {
	switch tp {
	case "NonUse":
		return BackpackUseType_NonUse, true
	case "Optional":
		return BackpackUseType_Optional, true
	case "Normal":
		return BackpackUseType_Normal, true
	case "Compose":
		return BackpackUseType_Compose, true
	case "RedPacket":
		return BackpackUseType_RedPacket, true

	}
	return -1, false

}

func GetBackpackUseType_ByAlias(tp string) (BackpackUseType, bool) {
	switch tp {
	case "不可使用":
		return BackpackUseType_NonUse, true
	case "自选":
		return BackpackUseType_Optional, true
	case "直接使用":
		return BackpackUseType_Normal, true
	case "合成":
		return BackpackUseType_Compose, true
	case "红包":
		return BackpackUseType_RedPacket, true

	}
	return -1, false

}
