package config

// 格斗家礼物-拳皇/BoxGiftEnum
type BoxGiftBlessBreakType int

const (
	BoxGiftBlessBreakType_BEGIN       BoxGiftBlessBreakType = 0 //枚举开始
	BoxGiftBlessBreakType_None        BoxGiftBlessBreakType = 0 //无
	BoxGiftBlessBreakType_BlessBreak  BoxGiftBlessBreakType = 1 //祝福突破
	BoxGiftBlessBreakType_RarityBreak BoxGiftBlessBreakType = 2 //品质突破
	BoxGiftBlessBreakType_END         BoxGiftBlessBreakType = 3 //枚举结束
)

func (e BoxGiftBlessBreakType) String() string {

	if e == BoxGiftBlessBreakType_None {
		return "None"
	}
	if e == BoxGiftBlessBreakType_BlessBreak {
		return "BlessBreak"
	}
	if e == BoxGiftBlessBreakType_RarityBreak {
		return "RarityBreak"
	}

	return ""

}

func GetBoxGiftBlessBreakType(tp string) (BoxGiftBlessBreakType, bool) {
	switch tp {
	case "None":
		return BoxGiftBlessBreakType_None, true
	case "BlessBreak":
		return BoxGiftBlessBreakType_BlessBreak, true
	case "RarityBreak":
		return BoxGiftBlessBreakType_RarityBreak, true

	}
	return -1, false

}

func GetBoxGiftBlessBreakType_ByAlias(tp string) (BoxGiftBlessBreakType, bool) {
	switch tp {
	case "无":
		return BoxGiftBlessBreakType_None, true
	case "祝福突破":
		return BoxGiftBlessBreakType_BlessBreak, true
	case "品质突破":
		return BoxGiftBlessBreakType_RarityBreak, true

	}
	return -1, false

}
