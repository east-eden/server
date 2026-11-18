package config

// 活动表-拳皇/ActivityType
type GiftCardType int

const (
	GiftCardType_BEGIN       GiftCardType = 1 //枚举开始
	GiftCardType_MonthCard   GiftCardType = 1 //月卡
	GiftCardType_ForeverCard GiftCardType = 2 //终身卡
	GiftCardType_END         GiftCardType = 3 //枚举结束
)

// 活动表-拳皇/ActivityType
type ActivityType int

const (
	ActivityType_BEGIN               ActivityType = 1  //枚举开始
	ActivityType_BattlePass          ActivityType = 1  //通行证
	ActivityType_SuperPackage        ActivityType = 2  //超值礼包
	ActivityType_RoyalBoss           ActivityType = 3  //皇家挑战
	ActivityType_SevenDaySign        ActivityType = 4  //七日签到
	ActivityType_SprintTasks         ActivityType = 5  //冲刺任务
	ActivityType_SprintWelfare       ActivityType = 6  //冲刺福利
	ActivityType_AccumulatedRecharge ActivityType = 7  //累计充值
	ActivityType_FirstRecharge       ActivityType = 8  //首充
	ActivityType_SevenDayTask        ActivityType = 9  //七日任务
	ActivityType_WheelIsland         ActivityType = 10 //神秘海岛
	ActivityType_SeasonFestival      ActivityType = 11 //季节庆典
	ActivityType_KofArcade           ActivityType = 12 //拳皇联动
	ActivityType_DailyTopic          ActivityType = 13 //每日话题
	ActivityType_FlipCard            ActivityType = 14 //拳皇翻翻乐
	ActivityType_END                 ActivityType = 15 //枚举结束
)

// 活动表-拳皇/ActivityType
type BattlePassRule int

const (
	BattlePassRule_BEGIN                 BattlePassRule = 1 //枚举开始
	BattlePassRule_HeroRecuritType       BattlePassRule = 1 //英雄招募类型
	BattlePassRule_HeroStarUpCostChip    BattlePassRule = 2 //英雄升星消耗碎片
	BattlePassRule_HeroWakenCostChip     BattlePassRule = 3 //英雄觉醒消耗碎片
	BattlePassRule_ArtifactDrawType      BattlePassRule = 4 //神器抽取类型
	BattlePassRule_AllianceTechDrawTimes BattlePassRule = 5 //公会科技抽取次数
	BattlePassRule_AllianceTechGetType   BattlePassRule = 6 //公会科技抽到类型
	BattlePassRule_ArtifactDrawCost      BattlePassRule = 7 //神器抽取消耗
	BattlePassRule_END                   BattlePassRule = 8 //枚举结束
)

// 活动表-拳皇/ActivityType
type ActPuzzleShowType int

const (
	ActPuzzleShowType_BEGIN  ActPuzzleShowType = 1 //枚举开始
	ActPuzzleShowType_Puzzle ActPuzzleShowType = 1 //拼图
	ActPuzzleShowType_Gift   ActPuzzleShowType = 2 //礼盒
	ActPuzzleShowType_END    ActPuzzleShowType = 3 //枚举结束
)

// 活动表-拳皇/ActivityType
type JiJinSubType int

const (
	JiJinSubType_BEGIN         JiJinSubType = 0  //枚举开始
	JiJinSubType_GrowUp        JiJinSubType = 0  //成长基金
	JiJinSubType_Challenge     JiJinSubType = 1  //挑战基金
	JiJinSubType_Honor         JiJinSubType = 2  //荣誉基金
	JiJinSubType_ArtifactPass  JiJinSubType = 3  //神器基金
	JiJinSubType_CollectPass   JiJinSubType = 4  //收藏基金
	JiJinSubType_NestPass      JiJinSubType = 5  //音巢基金
	JiJinSubType_RecruitPass   JiJinSubType = 6  //招募基金
	JiJinSubType_HeroStarPass  JiJinSubType = 7  //格斗家星级基金
	JiJinSubType_SecreArtsPass JiJinSubType = 8  //秘术基金
	JiJinSubType_HeroSoulPass  JiJinSubType = 9  //八门基金
	JiJinSubType_StagePass     JiJinSubType = 10 //关卡基金
	JiJinSubType_ManorPass     JiJinSubType = 11 //领地基金
	JiJinSubType_PetPass       JiJinSubType = 12 //战魂基金
	JiJinSubType_JiJinCount    JiJinSubType = 13 //基金数量
	JiJinSubType_END           JiJinSubType = 14 //枚举结束
)

func (e GiftCardType) String() string {

	if e == GiftCardType_MonthCard {
		return "MonthCard"
	}
	if e == GiftCardType_ForeverCard {
		return "ForeverCard"
	}

	return ""

}

func GetGiftCardType(tp string) (GiftCardType, bool) {
	switch tp {
	case "MonthCard":
		return GiftCardType_MonthCard, true
	case "ForeverCard":
		return GiftCardType_ForeverCard, true

	}
	return -1, false

}

func GetGiftCardType_ByAlias(tp string) (GiftCardType, bool) {
	switch tp {
	case "月卡":
		return GiftCardType_MonthCard, true
	case "终身卡":
		return GiftCardType_ForeverCard, true

	}
	return -1, false

}

func (e ActivityType) String() string {

	if e == ActivityType_BattlePass {
		return "BattlePass"
	}
	if e == ActivityType_SuperPackage {
		return "SuperPackage"
	}
	if e == ActivityType_RoyalBoss {
		return "RoyalBoss"
	}
	if e == ActivityType_SevenDaySign {
		return "SevenDaySign"
	}
	if e == ActivityType_SprintTasks {
		return "SprintTasks"
	}
	if e == ActivityType_SprintWelfare {
		return "SprintWelfare"
	}
	if e == ActivityType_AccumulatedRecharge {
		return "AccumulatedRecharge"
	}
	if e == ActivityType_FirstRecharge {
		return "FirstRecharge"
	}
	if e == ActivityType_SevenDayTask {
		return "SevenDayTask"
	}
	if e == ActivityType_WheelIsland {
		return "WheelIsland"
	}
	if e == ActivityType_SeasonFestival {
		return "SeasonFestival"
	}
	if e == ActivityType_KofArcade {
		return "KofArcade"
	}
	if e == ActivityType_DailyTopic {
		return "DailyTopic"
	}
	if e == ActivityType_FlipCard {
		return "FlipCard"
	}

	return ""

}

func GetActivityType(tp string) (ActivityType, bool) {
	switch tp {
	case "BattlePass":
		return ActivityType_BattlePass, true
	case "SuperPackage":
		return ActivityType_SuperPackage, true
	case "RoyalBoss":
		return ActivityType_RoyalBoss, true
	case "SevenDaySign":
		return ActivityType_SevenDaySign, true
	case "SprintTasks":
		return ActivityType_SprintTasks, true
	case "SprintWelfare":
		return ActivityType_SprintWelfare, true
	case "AccumulatedRecharge":
		return ActivityType_AccumulatedRecharge, true
	case "FirstRecharge":
		return ActivityType_FirstRecharge, true
	case "SevenDayTask":
		return ActivityType_SevenDayTask, true
	case "WheelIsland":
		return ActivityType_WheelIsland, true
	case "SeasonFestival":
		return ActivityType_SeasonFestival, true
	case "KofArcade":
		return ActivityType_KofArcade, true
	case "DailyTopic":
		return ActivityType_DailyTopic, true
	case "FlipCard":
		return ActivityType_FlipCard, true

	}
	return -1, false

}

func GetActivityType_ByAlias(tp string) (ActivityType, bool) {
	switch tp {
	case "通行证":
		return ActivityType_BattlePass, true
	case "超值礼包":
		return ActivityType_SuperPackage, true
	case "皇家挑战":
		return ActivityType_RoyalBoss, true
	case "七日签到":
		return ActivityType_SevenDaySign, true
	case "冲刺任务":
		return ActivityType_SprintTasks, true
	case "冲刺福利":
		return ActivityType_SprintWelfare, true
	case "累计充值":
		return ActivityType_AccumulatedRecharge, true
	case "首充":
		return ActivityType_FirstRecharge, true
	case "七日任务":
		return ActivityType_SevenDayTask, true
	case "神秘海岛":
		return ActivityType_WheelIsland, true
	case "季节庆典":
		return ActivityType_SeasonFestival, true
	case "拳皇联动":
		return ActivityType_KofArcade, true
	case "每日话题":
		return ActivityType_DailyTopic, true
	case "拳皇翻翻乐":
		return ActivityType_FlipCard, true

	}
	return -1, false

}

func (e BattlePassRule) String() string {

	if e == BattlePassRule_HeroRecuritType {
		return "HeroRecuritType"
	}
	if e == BattlePassRule_HeroStarUpCostChip {
		return "HeroStarUpCostChip"
	}
	if e == BattlePassRule_HeroWakenCostChip {
		return "HeroWakenCostChip"
	}
	if e == BattlePassRule_ArtifactDrawType {
		return "ArtifactDrawType"
	}
	if e == BattlePassRule_AllianceTechDrawTimes {
		return "AllianceTechDrawTimes"
	}
	if e == BattlePassRule_AllianceTechGetType {
		return "AllianceTechGetType"
	}
	if e == BattlePassRule_ArtifactDrawCost {
		return "ArtifactDrawCost"
	}

	return ""

}

func GetBattlePassRule(tp string) (BattlePassRule, bool) {
	switch tp {
	case "HeroRecuritType":
		return BattlePassRule_HeroRecuritType, true
	case "HeroStarUpCostChip":
		return BattlePassRule_HeroStarUpCostChip, true
	case "HeroWakenCostChip":
		return BattlePassRule_HeroWakenCostChip, true
	case "ArtifactDrawType":
		return BattlePassRule_ArtifactDrawType, true
	case "AllianceTechDrawTimes":
		return BattlePassRule_AllianceTechDrawTimes, true
	case "AllianceTechGetType":
		return BattlePassRule_AllianceTechGetType, true
	case "ArtifactDrawCost":
		return BattlePassRule_ArtifactDrawCost, true

	}
	return -1, false

}

func GetBattlePassRule_ByAlias(tp string) (BattlePassRule, bool) {
	switch tp {
	case "英雄招募类型":
		return BattlePassRule_HeroRecuritType, true
	case "英雄升星消耗碎片":
		return BattlePassRule_HeroStarUpCostChip, true
	case "英雄觉醒消耗碎片":
		return BattlePassRule_HeroWakenCostChip, true
	case "神器抽取类型":
		return BattlePassRule_ArtifactDrawType, true
	case "公会科技抽取次数":
		return BattlePassRule_AllianceTechDrawTimes, true
	case "公会科技抽到类型":
		return BattlePassRule_AllianceTechGetType, true
	case "神器抽取消耗":
		return BattlePassRule_ArtifactDrawCost, true

	}
	return -1, false

}

func (e ActPuzzleShowType) String() string {

	if e == ActPuzzleShowType_Puzzle {
		return "Puzzle"
	}
	if e == ActPuzzleShowType_Gift {
		return "Gift"
	}

	return ""

}

func GetActPuzzleShowType(tp string) (ActPuzzleShowType, bool) {
	switch tp {
	case "Puzzle":
		return ActPuzzleShowType_Puzzle, true
	case "Gift":
		return ActPuzzleShowType_Gift, true

	}
	return -1, false

}

func GetActPuzzleShowType_ByAlias(tp string) (ActPuzzleShowType, bool) {
	switch tp {
	case "拼图":
		return ActPuzzleShowType_Puzzle, true
	case "礼盒":
		return ActPuzzleShowType_Gift, true

	}
	return -1, false

}

func (e JiJinSubType) String() string {

	if e == JiJinSubType_GrowUp {
		return "GrowUp"
	}
	if e == JiJinSubType_Challenge {
		return "Challenge"
	}
	if e == JiJinSubType_Honor {
		return "Honor"
	}
	if e == JiJinSubType_ArtifactPass {
		return "ArtifactPass"
	}
	if e == JiJinSubType_CollectPass {
		return "CollectPass"
	}
	if e == JiJinSubType_NestPass {
		return "NestPass"
	}
	if e == JiJinSubType_RecruitPass {
		return "RecruitPass"
	}
	if e == JiJinSubType_HeroStarPass {
		return "HeroStarPass"
	}
	if e == JiJinSubType_SecreArtsPass {
		return "SecreArtsPass"
	}
	if e == JiJinSubType_HeroSoulPass {
		return "HeroSoulPass"
	}
	if e == JiJinSubType_StagePass {
		return "StagePass"
	}
	if e == JiJinSubType_ManorPass {
		return "ManorPass"
	}
	if e == JiJinSubType_PetPass {
		return "PetPass"
	}
	if e == JiJinSubType_JiJinCount {
		return "JiJinCount"
	}

	return ""

}

func GetJiJinSubType(tp string) (JiJinSubType, bool) {
	switch tp {
	case "GrowUp":
		return JiJinSubType_GrowUp, true
	case "Challenge":
		return JiJinSubType_Challenge, true
	case "Honor":
		return JiJinSubType_Honor, true
	case "ArtifactPass":
		return JiJinSubType_ArtifactPass, true
	case "CollectPass":
		return JiJinSubType_CollectPass, true
	case "NestPass":
		return JiJinSubType_NestPass, true
	case "RecruitPass":
		return JiJinSubType_RecruitPass, true
	case "HeroStarPass":
		return JiJinSubType_HeroStarPass, true
	case "SecreArtsPass":
		return JiJinSubType_SecreArtsPass, true
	case "HeroSoulPass":
		return JiJinSubType_HeroSoulPass, true
	case "StagePass":
		return JiJinSubType_StagePass, true
	case "ManorPass":
		return JiJinSubType_ManorPass, true
	case "PetPass":
		return JiJinSubType_PetPass, true
	case "JiJinCount":
		return JiJinSubType_JiJinCount, true

	}
	return -1, false

}

func GetJiJinSubType_ByAlias(tp string) (JiJinSubType, bool) {
	switch tp {
	case "成长基金":
		return JiJinSubType_GrowUp, true
	case "挑战基金":
		return JiJinSubType_Challenge, true
	case "荣誉基金":
		return JiJinSubType_Honor, true
	case "神器基金":
		return JiJinSubType_ArtifactPass, true
	case "收藏基金":
		return JiJinSubType_CollectPass, true
	case "音巢基金":
		return JiJinSubType_NestPass, true
	case "招募基金":
		return JiJinSubType_RecruitPass, true
	case "格斗家星级基金":
		return JiJinSubType_HeroStarPass, true
	case "秘术基金":
		return JiJinSubType_SecreArtsPass, true
	case "八门基金":
		return JiJinSubType_HeroSoulPass, true
	case "关卡基金":
		return JiJinSubType_StagePass, true
	case "领地基金":
		return JiJinSubType_ManorPass, true
	case "战魂基金":
		return JiJinSubType_PetPass, true
	case "基金数量":
		return JiJinSubType_JiJinCount, true

	}
	return -1, false

}
