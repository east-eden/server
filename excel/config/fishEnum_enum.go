package config

// 钓鱼-拳皇/FishEnum
type FishAttrType int

const (
	FishAttrType_BEGIN              FishAttrType = 1   //枚举开始
	FishAttrType_BigFish            FishAttrType = 601 //大鱼概率
	FishAttrType_FishingSpeed       FishAttrType = 602 //垂钓速度
	FishAttrType_FishLucky          FishAttrType = 603 //幸运
	FishAttrType_Eye                FishAttrType = 604 //慧眼
	FishAttrType_Focus              FishAttrType = 605 //专注
	FishAttrType_Skillful           FishAttrType = 606 //精通
	FishAttrType_Fisher             FishAttrType = 607 //渔夫
	FishAttrType_SmallFisher        FishAttrType = 608 //小型鱼猎手
	FishAttrType_MiddleFisher       FishAttrType = 609 //中型鱼猎手
	FishAttrType_BigFisher          FishAttrType = 610 //大型鱼猎手
	FishAttrType_Biologist          FishAttrType = 611 //生物学家
	FishAttrType_Efficiently        FishAttrType = 612 //高效
	FishAttrType_Rubbish            FishAttrType = 613 //拾荒
	FishAttrType_LowVision          FishAttrType = 614 //小型鱼克星
	FishAttrType_BadVision          FishAttrType = 615 //中型鱼克星
	FishAttrType_NoVision           FishAttrType = 616 //大型鱼克星
	FishAttrType_Strength           FishAttrType = 617 //强度
	FishAttrType_Sweeper            FishAttrType = 618 //清道夫
	FishAttrType_Feeder             FishAttrType = 619 //饲养员
	FishAttrType_MoneyFisher        FishAttrType = 620 //赏金钓客
	FishAttrType_Vitality           FishAttrType = 621 //活力
	FishAttrType_SmallScore         FishAttrType = 623 //小型鱼积分
	FishAttrType_MiddleScore        FishAttrType = 624 //中型鱼积分
	FishAttrType_BigScore           FishAttrType = 625 //大型鱼积分
	FishAttrType_BiologyScore       FishAttrType = 626 //水中生物积分
	FishAttrType_SmallSeaman        FishAttrType = 627 //专一的水手
	FishAttrType_MiddleSeaman       FishAttrType = 628 //专一的水手Ⅱ
	FishAttrType_BigSeaman          FishAttrType = 629 //专一的水手Ⅲ
	FishAttrType_BiologySeaman      FishAttrType = 630 //专一的水手Ⅳ
	FishAttrType_RubbishSeaman      FishAttrType = 631 //专一的水手Ⅴ
	FishAttrType_RubbishScoreUp     FishAttrType = 632 //变废为宝
	FishAttrType_BigLowButFast      FishAttrType = 633 //风驰电掣
	FishAttrType_NoFishScore        FishAttrType = 634 //不幸中的万幸
	FishAttrType_UnlimitedBait      FishAttrType = 635 //无限饵料
	FishAttrType_FishFriend         FishAttrType = 636 //鱼类之友
	FishAttrType_SuperBait          FishAttrType = 637 //超级饵料
	FishAttrType_FishMutation       FishAttrType = 638 //突然变异
	FishAttrType_StrengthMore       FishAttrType = 639 //金线
	FishAttrType_DiscoverCollection FishAttrType = 640 //考古
	FishAttrType_BackBait           FishAttrType = 641 //及时止损
	FishAttrType_PearlBonus         FishAttrType = 642 //意外发现
	FishAttrType_DiscoverUnlock     FishAttrType = 643 //见多识广
	FishAttrType_FeederTwo          FishAttrType = 644 //饲养员II
	FishAttrType_CollectExpert      FishAttrType = 701 //收藏专家
	FishAttrType_END                FishAttrType = 702 //枚举结束
)

func (e FishAttrType) String() string {

	if e == FishAttrType_BigFish {
		return "BigFish"
	}
	if e == FishAttrType_FishingSpeed {
		return "FishingSpeed"
	}
	if e == FishAttrType_FishLucky {
		return "FishLucky"
	}
	if e == FishAttrType_Eye {
		return "Eye"
	}
	if e == FishAttrType_Focus {
		return "Focus"
	}
	if e == FishAttrType_Skillful {
		return "Skillful"
	}
	if e == FishAttrType_Fisher {
		return "Fisher"
	}
	if e == FishAttrType_SmallFisher {
		return "SmallFisher"
	}
	if e == FishAttrType_MiddleFisher {
		return "MiddleFisher"
	}
	if e == FishAttrType_BigFisher {
		return "BigFisher"
	}
	if e == FishAttrType_Biologist {
		return "Biologist"
	}
	if e == FishAttrType_Efficiently {
		return "Efficiently"
	}
	if e == FishAttrType_Rubbish {
		return "Rubbish"
	}
	if e == FishAttrType_LowVision {
		return "LowVision"
	}
	if e == FishAttrType_BadVision {
		return "BadVision"
	}
	if e == FishAttrType_NoVision {
		return "NoVision"
	}
	if e == FishAttrType_Strength {
		return "Strength"
	}
	if e == FishAttrType_Sweeper {
		return "Sweeper"
	}
	if e == FishAttrType_Feeder {
		return "Feeder"
	}
	if e == FishAttrType_MoneyFisher {
		return "MoneyFisher"
	}
	if e == FishAttrType_Vitality {
		return "Vitality"
	}
	if e == FishAttrType_SmallScore {
		return "SmallScore"
	}
	if e == FishAttrType_MiddleScore {
		return "MiddleScore"
	}
	if e == FishAttrType_BigScore {
		return "BigScore"
	}
	if e == FishAttrType_BiologyScore {
		return "BiologyScore"
	}
	if e == FishAttrType_SmallSeaman {
		return "SmallSeaman"
	}
	if e == FishAttrType_MiddleSeaman {
		return "MiddleSeaman"
	}
	if e == FishAttrType_BigSeaman {
		return "BigSeaman"
	}
	if e == FishAttrType_BiologySeaman {
		return "BiologySeaman"
	}
	if e == FishAttrType_RubbishSeaman {
		return "RubbishSeaman"
	}
	if e == FishAttrType_RubbishScoreUp {
		return "RubbishScoreUp"
	}
	if e == FishAttrType_BigLowButFast {
		return "BigLowButFast"
	}
	if e == FishAttrType_NoFishScore {
		return "NoFishScore"
	}
	if e == FishAttrType_UnlimitedBait {
		return "UnlimitedBait"
	}
	if e == FishAttrType_FishFriend {
		return "FishFriend"
	}
	if e == FishAttrType_SuperBait {
		return "SuperBait"
	}
	if e == FishAttrType_FishMutation {
		return "FishMutation"
	}
	if e == FishAttrType_StrengthMore {
		return "StrengthMore"
	}
	if e == FishAttrType_DiscoverCollection {
		return "DiscoverCollection"
	}
	if e == FishAttrType_BackBait {
		return "BackBait"
	}
	if e == FishAttrType_PearlBonus {
		return "PearlBonus"
	}
	if e == FishAttrType_DiscoverUnlock {
		return "DiscoverUnlock"
	}
	if e == FishAttrType_FeederTwo {
		return "FeederTwo"
	}
	if e == FishAttrType_CollectExpert {
		return "CollectExpert"
	}

	return ""

}

func GetFishAttrType(tp string) (FishAttrType, bool) {
	switch tp {
	case "BigFish":
		return FishAttrType_BigFish, true
	case "FishingSpeed":
		return FishAttrType_FishingSpeed, true
	case "FishLucky":
		return FishAttrType_FishLucky, true
	case "Eye":
		return FishAttrType_Eye, true
	case "Focus":
		return FishAttrType_Focus, true
	case "Skillful":
		return FishAttrType_Skillful, true
	case "Fisher":
		return FishAttrType_Fisher, true
	case "SmallFisher":
		return FishAttrType_SmallFisher, true
	case "MiddleFisher":
		return FishAttrType_MiddleFisher, true
	case "BigFisher":
		return FishAttrType_BigFisher, true
	case "Biologist":
		return FishAttrType_Biologist, true
	case "Efficiently":
		return FishAttrType_Efficiently, true
	case "Rubbish":
		return FishAttrType_Rubbish, true
	case "LowVision":
		return FishAttrType_LowVision, true
	case "BadVision":
		return FishAttrType_BadVision, true
	case "NoVision":
		return FishAttrType_NoVision, true
	case "Strength":
		return FishAttrType_Strength, true
	case "Sweeper":
		return FishAttrType_Sweeper, true
	case "Feeder":
		return FishAttrType_Feeder, true
	case "MoneyFisher":
		return FishAttrType_MoneyFisher, true
	case "Vitality":
		return FishAttrType_Vitality, true
	case "SmallScore":
		return FishAttrType_SmallScore, true
	case "MiddleScore":
		return FishAttrType_MiddleScore, true
	case "BigScore":
		return FishAttrType_BigScore, true
	case "BiologyScore":
		return FishAttrType_BiologyScore, true
	case "SmallSeaman":
		return FishAttrType_SmallSeaman, true
	case "MiddleSeaman":
		return FishAttrType_MiddleSeaman, true
	case "BigSeaman":
		return FishAttrType_BigSeaman, true
	case "BiologySeaman":
		return FishAttrType_BiologySeaman, true
	case "RubbishSeaman":
		return FishAttrType_RubbishSeaman, true
	case "RubbishScoreUp":
		return FishAttrType_RubbishScoreUp, true
	case "BigLowButFast":
		return FishAttrType_BigLowButFast, true
	case "NoFishScore":
		return FishAttrType_NoFishScore, true
	case "UnlimitedBait":
		return FishAttrType_UnlimitedBait, true
	case "FishFriend":
		return FishAttrType_FishFriend, true
	case "SuperBait":
		return FishAttrType_SuperBait, true
	case "FishMutation":
		return FishAttrType_FishMutation, true
	case "StrengthMore":
		return FishAttrType_StrengthMore, true
	case "DiscoverCollection":
		return FishAttrType_DiscoverCollection, true
	case "BackBait":
		return FishAttrType_BackBait, true
	case "PearlBonus":
		return FishAttrType_PearlBonus, true
	case "DiscoverUnlock":
		return FishAttrType_DiscoverUnlock, true
	case "FeederTwo":
		return FishAttrType_FeederTwo, true
	case "CollectExpert":
		return FishAttrType_CollectExpert, true

	}
	return -1, false

}

func GetFishAttrType_ByAlias(tp string) (FishAttrType, bool) {
	switch tp {
	case "大鱼概率":
		return FishAttrType_BigFish, true
	case "垂钓速度":
		return FishAttrType_FishingSpeed, true
	case "幸运":
		return FishAttrType_FishLucky, true
	case "慧眼":
		return FishAttrType_Eye, true
	case "专注":
		return FishAttrType_Focus, true
	case "精通":
		return FishAttrType_Skillful, true
	case "渔夫":
		return FishAttrType_Fisher, true
	case "小型鱼猎手":
		return FishAttrType_SmallFisher, true
	case "中型鱼猎手":
		return FishAttrType_MiddleFisher, true
	case "大型鱼猎手":
		return FishAttrType_BigFisher, true
	case "生物学家":
		return FishAttrType_Biologist, true
	case "高效":
		return FishAttrType_Efficiently, true
	case "拾荒":
		return FishAttrType_Rubbish, true
	case "小型鱼克星":
		return FishAttrType_LowVision, true
	case "中型鱼克星":
		return FishAttrType_BadVision, true
	case "大型鱼克星":
		return FishAttrType_NoVision, true
	case "强度":
		return FishAttrType_Strength, true
	case "清道夫":
		return FishAttrType_Sweeper, true
	case "饲养员":
		return FishAttrType_Feeder, true
	case "赏金钓客":
		return FishAttrType_MoneyFisher, true
	case "活力":
		return FishAttrType_Vitality, true
	case "小型鱼积分":
		return FishAttrType_SmallScore, true
	case "中型鱼积分":
		return FishAttrType_MiddleScore, true
	case "大型鱼积分":
		return FishAttrType_BigScore, true
	case "水中生物积分":
		return FishAttrType_BiologyScore, true
	case "专一的水手":
		return FishAttrType_SmallSeaman, true
	case "专一的水手Ⅱ":
		return FishAttrType_MiddleSeaman, true
	case "专一的水手Ⅲ":
		return FishAttrType_BigSeaman, true
	case "专一的水手Ⅳ":
		return FishAttrType_BiologySeaman, true
	case "专一的水手Ⅴ":
		return FishAttrType_RubbishSeaman, true
	case "变废为宝":
		return FishAttrType_RubbishScoreUp, true
	case "风驰电掣":
		return FishAttrType_BigLowButFast, true
	case "不幸中的万幸":
		return FishAttrType_NoFishScore, true
	case "无限饵料":
		return FishAttrType_UnlimitedBait, true
	case "鱼类之友":
		return FishAttrType_FishFriend, true
	case "超级饵料":
		return FishAttrType_SuperBait, true
	case "突然变异":
		return FishAttrType_FishMutation, true
	case "金线":
		return FishAttrType_StrengthMore, true
	case "考古":
		return FishAttrType_DiscoverCollection, true
	case "及时止损":
		return FishAttrType_BackBait, true
	case "意外发现":
		return FishAttrType_PearlBonus, true
	case "见多识广":
		return FishAttrType_DiscoverUnlock, true
	case "饲养员II":
		return FishAttrType_FeederTwo, true
	case "收藏专家":
		return FishAttrType_CollectExpert, true

	}
	return -1, false

}
