package config

// 礼包表-拳皇/PackageType
type PackagePopType int

const (
	PackagePopType_BEGIN   PackagePopType = 1 //枚举开始
	PackagePopType_Daily   PackagePopType = 1 //每日弹出
	PackagePopType_OnlyOne PackagePopType = 2 //仅一次
	PackagePopType_END     PackagePopType = 3 //枚举结束
)

// 礼包表-拳皇/PackageType
type PackageType int

const (
	PackageType_BEGIN PackageType = 1 //枚举开始
	PackageType_Popup PackageType = 1 //弹出礼包
	PackageType_END   PackageType = 2 //枚举结束
)

// 礼包表-拳皇/PackageType
type PackageIds int

const (
	PackageIds_BEGIN              PackageIds = 1  //枚举开始
	PackageIds_CorpsRecruit       PackageIds = 1  //佣兵招募
	PackageIds_TerritoryResources PackageIds = 2  //领地资源
	PackageIds_RelicsUpgrade      PackageIds = 5  //圣物成长
	PackageIds_CorpLeader         PackageIds = 6  //团长等级成长
	PackageIds_BoxTimeTicket      PackageIds = 7  //加速券
	PackageIds_DiamondNotEnough   PackageIds = 14 //钻石不足
	PackageIds_FirstRechargeNew   PackageIds = 13 //新版首充
	PackageIds_END                PackageIds = 15 //枚举结束
)

// 礼包表-拳皇/PackageType
type WeekCardType int

const (
	WeekCardType_BEGIN            WeekCardType = 0 //枚举开始
	WeekCardType_HeroWeekCard     WeekCardType = 0 //英雄周卡
	WeekCardType_FishWeekCard     WeekCardType = 1 //钓鱼周卡
	WeekCardType_TreasureWeekCard WeekCardType = 2 //宝物周卡
	WeekCardType_END              WeekCardType = 3 //枚举结束
)

// 礼包表-拳皇/PackageType
type PrivilegeCardEnum int

const (
	PrivilegeCardEnum_BEGIN            PrivilegeCardEnum = 1  //枚举开始
	PrivilegeCardEnum_WeekCard         PrivilegeCardEnum = 1  //周卡
	PrivilegeCardEnum_MonthCard        PrivilegeCardEnum = 2  //月卡
	PrivilegeCardEnum_EnjoymentCard    PrivilegeCardEnum = 3  //畅玩卡
	PrivilegeCardEnum_LifeCard         PrivilegeCardEnum = 4  //终身卡
	PrivilegeCardEnum_FishCard         PrivilegeCardEnum = 5  //钓鱼卡
	PrivilegeCardEnum_WeekHeroCard     PrivilegeCardEnum = 6  //招募冲刺卡
	PrivilegeCardEnum_WeekArtifactCard PrivilegeCardEnum = 7  //神器冲刺卡
	PrivilegeCardEnum_WeekAlliTechCard PrivilegeCardEnum = 8  //秘术冲刺卡
	PrivilegeCardEnum_OneYuanCard      PrivilegeCardEnum = 9  //一元权益卡
	PrivilegeCardEnum_SoulCard         PrivilegeCardEnum = 10 //战魂权益卡
	PrivilegeCardEnum_DiamondCard      PrivilegeCardEnum = 11 //钻石权益卡
	PrivilegeCardEnum_END              PrivilegeCardEnum = 12 //枚举结束
)

func (e PackagePopType) String() string {

	if e == PackagePopType_Daily {
		return "Daily"
	}
	if e == PackagePopType_OnlyOne {
		return "OnlyOne"
	}

	return ""

}

func GetPackagePopType(tp string) (PackagePopType, bool) {
	switch tp {
	case "Daily":
		return PackagePopType_Daily, true
	case "OnlyOne":
		return PackagePopType_OnlyOne, true

	}
	return -1, false

}

func GetPackagePopType_ByAlias(tp string) (PackagePopType, bool) {
	switch tp {
	case "每日弹出":
		return PackagePopType_Daily, true
	case "仅一次":
		return PackagePopType_OnlyOne, true

	}
	return -1, false

}

func (e PackageType) String() string {

	if e == PackageType_Popup {
		return "Popup"
	}

	return ""

}

func GetPackageType(tp string) (PackageType, bool) {
	switch tp {
	case "Popup":
		return PackageType_Popup, true

	}
	return -1, false

}

func GetPackageType_ByAlias(tp string) (PackageType, bool) {
	switch tp {
	case "弹出礼包":
		return PackageType_Popup, true

	}
	return -1, false

}

func (e PackageIds) String() string {

	if e == PackageIds_CorpsRecruit {
		return "CorpsRecruit"
	}
	if e == PackageIds_TerritoryResources {
		return "TerritoryResources"
	}
	if e == PackageIds_RelicsUpgrade {
		return "RelicsUpgrade"
	}
	if e == PackageIds_CorpLeader {
		return "CorpLeader"
	}
	if e == PackageIds_BoxTimeTicket {
		return "BoxTimeTicket"
	}
	if e == PackageIds_DiamondNotEnough {
		return "DiamondNotEnough"
	}
	if e == PackageIds_FirstRechargeNew {
		return "FirstRechargeNew"
	}

	return ""

}

func GetPackageIds(tp string) (PackageIds, bool) {
	switch tp {
	case "CorpsRecruit":
		return PackageIds_CorpsRecruit, true
	case "TerritoryResources":
		return PackageIds_TerritoryResources, true
	case "RelicsUpgrade":
		return PackageIds_RelicsUpgrade, true
	case "CorpLeader":
		return PackageIds_CorpLeader, true
	case "BoxTimeTicket":
		return PackageIds_BoxTimeTicket, true
	case "DiamondNotEnough":
		return PackageIds_DiamondNotEnough, true
	case "FirstRechargeNew":
		return PackageIds_FirstRechargeNew, true

	}
	return -1, false

}

func GetPackageIds_ByAlias(tp string) (PackageIds, bool) {
	switch tp {
	case "佣兵招募":
		return PackageIds_CorpsRecruit, true
	case "领地资源":
		return PackageIds_TerritoryResources, true
	case "圣物成长":
		return PackageIds_RelicsUpgrade, true
	case "团长等级成长":
		return PackageIds_CorpLeader, true
	case "加速券":
		return PackageIds_BoxTimeTicket, true
	case "钻石不足":
		return PackageIds_DiamondNotEnough, true
	case "新版首充":
		return PackageIds_FirstRechargeNew, true

	}
	return -1, false

}

func (e WeekCardType) String() string {

	if e == WeekCardType_HeroWeekCard {
		return "HeroWeekCard"
	}
	if e == WeekCardType_FishWeekCard {
		return "FishWeekCard"
	}
	if e == WeekCardType_TreasureWeekCard {
		return "TreasureWeekCard"
	}

	return ""

}

func GetWeekCardType(tp string) (WeekCardType, bool) {
	switch tp {
	case "HeroWeekCard":
		return WeekCardType_HeroWeekCard, true
	case "FishWeekCard":
		return WeekCardType_FishWeekCard, true
	case "TreasureWeekCard":
		return WeekCardType_TreasureWeekCard, true

	}
	return -1, false

}

func GetWeekCardType_ByAlias(tp string) (WeekCardType, bool) {
	switch tp {
	case "英雄周卡":
		return WeekCardType_HeroWeekCard, true
	case "钓鱼周卡":
		return WeekCardType_FishWeekCard, true
	case "宝物周卡":
		return WeekCardType_TreasureWeekCard, true

	}
	return -1, false

}

func (e PrivilegeCardEnum) String() string {

	if e == PrivilegeCardEnum_WeekCard {
		return "WeekCard"
	}
	if e == PrivilegeCardEnum_MonthCard {
		return "MonthCard"
	}
	if e == PrivilegeCardEnum_EnjoymentCard {
		return "EnjoymentCard"
	}
	if e == PrivilegeCardEnum_LifeCard {
		return "LifeCard"
	}
	if e == PrivilegeCardEnum_FishCard {
		return "FishCard"
	}
	if e == PrivilegeCardEnum_WeekHeroCard {
		return "WeekHeroCard"
	}
	if e == PrivilegeCardEnum_WeekArtifactCard {
		return "WeekArtifactCard"
	}
	if e == PrivilegeCardEnum_WeekAlliTechCard {
		return "WeekAlliTechCard"
	}
	if e == PrivilegeCardEnum_OneYuanCard {
		return "OneYuanCard"
	}
	if e == PrivilegeCardEnum_SoulCard {
		return "SoulCard"
	}
	if e == PrivilegeCardEnum_DiamondCard {
		return "DiamondCard"
	}

	return ""

}

func GetPrivilegeCardEnum(tp string) (PrivilegeCardEnum, bool) {
	switch tp {
	case "WeekCard":
		return PrivilegeCardEnum_WeekCard, true
	case "MonthCard":
		return PrivilegeCardEnum_MonthCard, true
	case "EnjoymentCard":
		return PrivilegeCardEnum_EnjoymentCard, true
	case "LifeCard":
		return PrivilegeCardEnum_LifeCard, true
	case "FishCard":
		return PrivilegeCardEnum_FishCard, true
	case "WeekHeroCard":
		return PrivilegeCardEnum_WeekHeroCard, true
	case "WeekArtifactCard":
		return PrivilegeCardEnum_WeekArtifactCard, true
	case "WeekAlliTechCard":
		return PrivilegeCardEnum_WeekAlliTechCard, true
	case "OneYuanCard":
		return PrivilegeCardEnum_OneYuanCard, true
	case "SoulCard":
		return PrivilegeCardEnum_SoulCard, true
	case "DiamondCard":
		return PrivilegeCardEnum_DiamondCard, true

	}
	return -1, false

}

func GetPrivilegeCardEnum_ByAlias(tp string) (PrivilegeCardEnum, bool) {
	switch tp {
	case "周卡":
		return PrivilegeCardEnum_WeekCard, true
	case "月卡":
		return PrivilegeCardEnum_MonthCard, true
	case "畅玩卡":
		return PrivilegeCardEnum_EnjoymentCard, true
	case "终身卡":
		return PrivilegeCardEnum_LifeCard, true
	case "钓鱼卡":
		return PrivilegeCardEnum_FishCard, true
	case "招募冲刺卡":
		return PrivilegeCardEnum_WeekHeroCard, true
	case "神器冲刺卡":
		return PrivilegeCardEnum_WeekArtifactCard, true
	case "秘术冲刺卡":
		return PrivilegeCardEnum_WeekAlliTechCard, true
	case "一元权益卡":
		return PrivilegeCardEnum_OneYuanCard, true
	case "战魂权益卡":
		return PrivilegeCardEnum_SoulCard, true
	case "钻石权益卡":
		return PrivilegeCardEnum_DiamondCard, true

	}
	return -1, false

}
