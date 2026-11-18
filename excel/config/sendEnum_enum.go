package config

// 赠送表-拳皇/sendEnum
type SendType int

const (
	SendType_BEGIN          SendType = 1 //枚举开始
	SendType_HeroSoul       SendType = 1 //八门
	SendType_FashionClothes SendType = 2 //时装
	SendType_CommonItem     SendType = 3 //普通物品
	SendType_FashionBox     SendType = 4 //盲盒
	SendType_NewYear        SendType = 5 //五福
	SendType_AniAct         SendType = 6 //周年庆
	SendType_END            SendType = 7 //枚举结束
)

func (e SendType) String() string {

	if e == SendType_HeroSoul {
		return "HeroSoul"
	}
	if e == SendType_FashionClothes {
		return "FashionClothes"
	}
	if e == SendType_CommonItem {
		return "CommonItem"
	}
	if e == SendType_FashionBox {
		return "FashionBox"
	}
	if e == SendType_NewYear {
		return "NewYear"
	}
	if e == SendType_AniAct {
		return "AniAct"
	}

	return ""

}

func GetSendType(tp string) (SendType, bool) {
	switch tp {
	case "HeroSoul":
		return SendType_HeroSoul, true
	case "FashionClothes":
		return SendType_FashionClothes, true
	case "CommonItem":
		return SendType_CommonItem, true
	case "FashionBox":
		return SendType_FashionBox, true
	case "NewYear":
		return SendType_NewYear, true
	case "AniAct":
		return SendType_AniAct, true

	}
	return -1, false

}

func GetSendType_ByAlias(tp string) (SendType, bool) {
	switch tp {
	case "八门":
		return SendType_HeroSoul, true
	case "时装":
		return SendType_FashionClothes, true
	case "普通物品":
		return SendType_CommonItem, true
	case "盲盒":
		return SendType_FashionBox, true
	case "五福":
		return SendType_NewYear, true
	case "周年庆":
		return SendType_AniAct, true

	}
	return -1, false

}
