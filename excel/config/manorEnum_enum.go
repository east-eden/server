package config

// 领地表-拳皇/ManorEnum
type ManorBuildingId int

const (
	ManorBuildingId_BEGIN    ManorBuildingId = 1  //枚举开始
	ManorBuildingId_Hall     ManorBuildingId = 1  //大厅
	ManorBuildingId_Tavern   ManorBuildingId = 2  //酒馆
	ManorBuildingId_Altar    ManorBuildingId = 4  //祭坛
	ManorBuildingId_Training ManorBuildingId = 5  //训练营
	ManorBuildingId_Corps1   ManorBuildingId = 21 //兵团1
	ManorBuildingId_Corps2   ManorBuildingId = 22 //兵团2
	ManorBuildingId_Corps3   ManorBuildingId = 23 //兵团3
	ManorBuildingId_Corps4   ManorBuildingId = 24 //兵团4
	ManorBuildingId_Corps5   ManorBuildingId = 25 //兵团5
	ManorBuildingId_END      ManorBuildingId = 26 //枚举结束
)

// 领地表-拳皇/ManorEnum
type ManorResourceType int

const (
	ManorResourceType_BEGIN             ManorResourceType = 1  //枚举开始
	ManorResourceType_Gold              ManorResourceType = 1  //金币
	ManorResourceType_Diamond           ManorResourceType = 2  //钻石
	ManorResourceType_Box               ManorResourceType = 3  //格斗家收藏盒
	ManorResourceType_SpeedTicket       ManorResourceType = 4  //加速券
	ManorResourceType_Stone             ManorResourceType = 5  //石料
	ManorResourceType_Wood              ManorResourceType = 6  //木材
	ManorResourceType_BoxExp            ManorResourceType = 7  //收藏经验
	ManorResourceType_HeroSoul          ManorResourceType = 8  //随机本源
	ManorResourceType_HeroSoulSynthesis ManorResourceType = 9  //本源凝液
	ManorResourceType_Meteor            ManorResourceType = 10 //陨石箱
	ManorResourceType_ItemBag           ManorResourceType = 11 //道具卡包
	ManorResourceType_BindClothesBox    ManorResourceType = 12 //绑定盲盒
	ManorResourceType_ClothesBox        ManorResourceType = 13 //盲盒
	ManorResourceType_END               ManorResourceType = 14 //枚举结束
)

// 领地表-拳皇/ManorEnum
type ManorTechType int

const (
	ManorTechType_BEGIN         ManorTechType = 0 //枚举开始
	ManorTechType_HorseStable   ManorTechType = 0 //马场
	ManorTechType_WoodenPost    ManorTechType = 1 //木桩
	ManorTechType_ShootingRange ManorTechType = 2 //靶场
	ManorTechType_Forge         ManorTechType = 3 //铁匠铺
	ManorTechType_END           ManorTechType = 4 //枚举结束
)

func (e ManorBuildingId) String() string {

	if e == ManorBuildingId_Hall {
		return "Hall"
	}
	if e == ManorBuildingId_Tavern {
		return "Tavern"
	}
	if e == ManorBuildingId_Altar {
		return "Altar"
	}
	if e == ManorBuildingId_Training {
		return "Training"
	}
	if e == ManorBuildingId_Corps1 {
		return "Corps1"
	}
	if e == ManorBuildingId_Corps2 {
		return "Corps2"
	}
	if e == ManorBuildingId_Corps3 {
		return "Corps3"
	}
	if e == ManorBuildingId_Corps4 {
		return "Corps4"
	}
	if e == ManorBuildingId_Corps5 {
		return "Corps5"
	}

	return ""

}

func GetManorBuildingId(tp string) (ManorBuildingId, bool) {
	switch tp {
	case "Hall":
		return ManorBuildingId_Hall, true
	case "Tavern":
		return ManorBuildingId_Tavern, true
	case "Altar":
		return ManorBuildingId_Altar, true
	case "Training":
		return ManorBuildingId_Training, true
	case "Corps1":
		return ManorBuildingId_Corps1, true
	case "Corps2":
		return ManorBuildingId_Corps2, true
	case "Corps3":
		return ManorBuildingId_Corps3, true
	case "Corps4":
		return ManorBuildingId_Corps4, true
	case "Corps5":
		return ManorBuildingId_Corps5, true

	}
	return -1, false

}

func GetManorBuildingId_ByAlias(tp string) (ManorBuildingId, bool) {
	switch tp {
	case "大厅":
		return ManorBuildingId_Hall, true
	case "酒馆":
		return ManorBuildingId_Tavern, true
	case "祭坛":
		return ManorBuildingId_Altar, true
	case "训练营":
		return ManorBuildingId_Training, true
	case "兵团1":
		return ManorBuildingId_Corps1, true
	case "兵团2":
		return ManorBuildingId_Corps2, true
	case "兵团3":
		return ManorBuildingId_Corps3, true
	case "兵团4":
		return ManorBuildingId_Corps4, true
	case "兵团5":
		return ManorBuildingId_Corps5, true

	}
	return -1, false

}

func (e ManorResourceType) String() string {

	if e == ManorResourceType_Gold {
		return "Gold"
	}
	if e == ManorResourceType_Diamond {
		return "Diamond"
	}
	if e == ManorResourceType_Box {
		return "Box"
	}
	if e == ManorResourceType_SpeedTicket {
		return "SpeedTicket"
	}
	if e == ManorResourceType_Stone {
		return "Stone"
	}
	if e == ManorResourceType_Wood {
		return "Wood"
	}
	if e == ManorResourceType_BoxExp {
		return "BoxExp"
	}
	if e == ManorResourceType_HeroSoul {
		return "HeroSoul"
	}
	if e == ManorResourceType_HeroSoulSynthesis {
		return "HeroSoulSynthesis"
	}
	if e == ManorResourceType_Meteor {
		return "Meteor"
	}
	if e == ManorResourceType_ItemBag {
		return "ItemBag"
	}
	if e == ManorResourceType_BindClothesBox {
		return "BindClothesBox"
	}
	if e == ManorResourceType_ClothesBox {
		return "ClothesBox"
	}

	return ""

}

func GetManorResourceType(tp string) (ManorResourceType, bool) {
	switch tp {
	case "Gold":
		return ManorResourceType_Gold, true
	case "Diamond":
		return ManorResourceType_Diamond, true
	case "Box":
		return ManorResourceType_Box, true
	case "SpeedTicket":
		return ManorResourceType_SpeedTicket, true
	case "Stone":
		return ManorResourceType_Stone, true
	case "Wood":
		return ManorResourceType_Wood, true
	case "BoxExp":
		return ManorResourceType_BoxExp, true
	case "HeroSoul":
		return ManorResourceType_HeroSoul, true
	case "HeroSoulSynthesis":
		return ManorResourceType_HeroSoulSynthesis, true
	case "Meteor":
		return ManorResourceType_Meteor, true
	case "ItemBag":
		return ManorResourceType_ItemBag, true
	case "BindClothesBox":
		return ManorResourceType_BindClothesBox, true
	case "ClothesBox":
		return ManorResourceType_ClothesBox, true

	}
	return -1, false

}

func GetManorResourceType_ByAlias(tp string) (ManorResourceType, bool) {
	switch tp {
	case "金币":
		return ManorResourceType_Gold, true
	case "钻石":
		return ManorResourceType_Diamond, true
	case "格斗家收藏盒":
		return ManorResourceType_Box, true
	case "加速券":
		return ManorResourceType_SpeedTicket, true
	case "石料":
		return ManorResourceType_Stone, true
	case "木材":
		return ManorResourceType_Wood, true
	case "收藏经验":
		return ManorResourceType_BoxExp, true
	case "随机本源":
		return ManorResourceType_HeroSoul, true
	case "本源凝液":
		return ManorResourceType_HeroSoulSynthesis, true
	case "陨石箱":
		return ManorResourceType_Meteor, true
	case "道具卡包":
		return ManorResourceType_ItemBag, true
	case "绑定盲盒":
		return ManorResourceType_BindClothesBox, true
	case "盲盒":
		return ManorResourceType_ClothesBox, true

	}
	return -1, false

}

func (e ManorTechType) String() string {

	if e == ManorTechType_HorseStable {
		return "HorseStable"
	}
	if e == ManorTechType_WoodenPost {
		return "WoodenPost"
	}
	if e == ManorTechType_ShootingRange {
		return "ShootingRange"
	}
	if e == ManorTechType_Forge {
		return "Forge"
	}

	return ""

}

func GetManorTechType(tp string) (ManorTechType, bool) {
	switch tp {
	case "HorseStable":
		return ManorTechType_HorseStable, true
	case "WoodenPost":
		return ManorTechType_WoodenPost, true
	case "ShootingRange":
		return ManorTechType_ShootingRange, true
	case "Forge":
		return ManorTechType_Forge, true

	}
	return -1, false

}

func GetManorTechType_ByAlias(tp string) (ManorTechType, bool) {
	switch tp {
	case "马场":
		return ManorTechType_HorseStable, true
	case "木桩":
		return ManorTechType_WoodenPost, true
	case "靶场":
		return ManorTechType_ShootingRange, true
	case "铁匠铺":
		return ManorTechType_Forge, true

	}
	return -1, false

}
