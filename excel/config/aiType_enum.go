package config

// LD ai表-拳皇/aiType
type AIActionType int

const (
	AIActionType_BEGIN            AIActionType = 0  //枚举开始
	AIActionType_SummonHero       AIActionType = 0  //随机召唤
	AIActionType_Gamble           AIActionType = 1  //普通赌博
	AIActionType_RandomGamble     AIActionType = 2  //随机赌博
	AIActionType_MergeFormulaHero AIActionType = 3  //合成指定神话英雄
	AIActionType_GambleSpecHero   AIActionType = 4  //赌博出指定英雄
	AIActionType_ShopEnchant      AIActionType = 5  //商店强化
	AIActionType_HeroEnchant      AIActionType = 6  //英雄强化
	AIActionType_AddItem          AIActionType = 7  //增加道具
	AIActionType_SellHeroForGrid  AIActionType = 8  //出售到格子数
	AIActionType_SummonSpecHero   AIActionType = 9  //召唤出指定英雄
	AIActionType_END              AIActionType = 10 //枚举结束
)

// LD ai表-拳皇/aiType
type AIActionTriggerType int

const (
	AIActionTriggerType_BEGIN        AIActionTriggerType = 0 //枚举开始
	AIActionTriggerType_Empty        AIActionTriggerType = 0 //无条件
	AIActionTriggerType_Coin         AIActionTriggerType = 1 //硬币数
	AIActionTriggerType_LuckStone    AIActionTriggerType = 2 //幸运石数
	AIActionTriggerType_Wave         AIActionTriggerType = 3 //回合数
	AIActionTriggerType_MonsterNum   AIActionTriggerType = 4 //怪物数达到
	AIActionTriggerType_HasHeros     AIActionTriggerType = 5 //拥有英雄
	AIActionTriggerType_PlayerSummon AIActionTriggerType = 6 //玩家召唤
	AIActionTriggerType_END          AIActionTriggerType = 7 //枚举结束
)

// LD ai表-拳皇/aiType
type AIActionExitType int

const (
	AIActionExitType_BEGIN         AIActionExitType = 0 //枚举开始
	AIActionExitType_Empty         AIActionExitType = 0 //无条件
	AIActionExitType_Times         AIActionExitType = 1 //执行次数
	AIActionExitType_EpicHeroesNum AIActionExitType = 2 //有多少个神话英雄
	AIActionExitType_END           AIActionExitType = 3 //枚举结束
)

func (e AIActionType) String() string {

	if e == AIActionType_SummonHero {
		return "SummonHero"
	}
	if e == AIActionType_Gamble {
		return "Gamble"
	}
	if e == AIActionType_RandomGamble {
		return "RandomGamble"
	}
	if e == AIActionType_MergeFormulaHero {
		return "MergeFormulaHero"
	}
	if e == AIActionType_GambleSpecHero {
		return "GambleSpecHero"
	}
	if e == AIActionType_ShopEnchant {
		return "ShopEnchant"
	}
	if e == AIActionType_HeroEnchant {
		return "HeroEnchant"
	}
	if e == AIActionType_AddItem {
		return "AddItem"
	}
	if e == AIActionType_SellHeroForGrid {
		return "SellHeroForGrid"
	}
	if e == AIActionType_SummonSpecHero {
		return "SummonSpecHero"
	}

	return ""

}

func GetAIActionType(tp string) (AIActionType, bool) {
	switch tp {
	case "SummonHero":
		return AIActionType_SummonHero, true
	case "Gamble":
		return AIActionType_Gamble, true
	case "RandomGamble":
		return AIActionType_RandomGamble, true
	case "MergeFormulaHero":
		return AIActionType_MergeFormulaHero, true
	case "GambleSpecHero":
		return AIActionType_GambleSpecHero, true
	case "ShopEnchant":
		return AIActionType_ShopEnchant, true
	case "HeroEnchant":
		return AIActionType_HeroEnchant, true
	case "AddItem":
		return AIActionType_AddItem, true
	case "SellHeroForGrid":
		return AIActionType_SellHeroForGrid, true
	case "SummonSpecHero":
		return AIActionType_SummonSpecHero, true

	}
	return -1, false

}

func GetAIActionType_ByAlias(tp string) (AIActionType, bool) {
	switch tp {
	case "随机召唤":
		return AIActionType_SummonHero, true
	case "普通赌博":
		return AIActionType_Gamble, true
	case "随机赌博":
		return AIActionType_RandomGamble, true
	case "合成指定神话英雄":
		return AIActionType_MergeFormulaHero, true
	case "赌博出指定英雄":
		return AIActionType_GambleSpecHero, true
	case "商店强化":
		return AIActionType_ShopEnchant, true
	case "英雄强化":
		return AIActionType_HeroEnchant, true
	case "增加道具":
		return AIActionType_AddItem, true
	case "出售到格子数":
		return AIActionType_SellHeroForGrid, true
	case "召唤出指定英雄":
		return AIActionType_SummonSpecHero, true

	}
	return -1, false

}

func (e AIActionTriggerType) String() string {

	if e == AIActionTriggerType_Empty {
		return "Empty"
	}
	if e == AIActionTriggerType_Coin {
		return "Coin"
	}
	if e == AIActionTriggerType_LuckStone {
		return "LuckStone"
	}
	if e == AIActionTriggerType_Wave {
		return "Wave"
	}
	if e == AIActionTriggerType_MonsterNum {
		return "MonsterNum"
	}
	if e == AIActionTriggerType_HasHeros {
		return "HasHeros"
	}
	if e == AIActionTriggerType_PlayerSummon {
		return "PlayerSummon"
	}

	return ""

}

func GetAIActionTriggerType(tp string) (AIActionTriggerType, bool) {
	switch tp {
	case "Empty":
		return AIActionTriggerType_Empty, true
	case "Coin":
		return AIActionTriggerType_Coin, true
	case "LuckStone":
		return AIActionTriggerType_LuckStone, true
	case "Wave":
		return AIActionTriggerType_Wave, true
	case "MonsterNum":
		return AIActionTriggerType_MonsterNum, true
	case "HasHeros":
		return AIActionTriggerType_HasHeros, true
	case "PlayerSummon":
		return AIActionTriggerType_PlayerSummon, true

	}
	return -1, false

}

func GetAIActionTriggerType_ByAlias(tp string) (AIActionTriggerType, bool) {
	switch tp {
	case "无条件":
		return AIActionTriggerType_Empty, true
	case "硬币数":
		return AIActionTriggerType_Coin, true
	case "幸运石数":
		return AIActionTriggerType_LuckStone, true
	case "回合数":
		return AIActionTriggerType_Wave, true
	case "怪物数达到":
		return AIActionTriggerType_MonsterNum, true
	case "拥有英雄":
		return AIActionTriggerType_HasHeros, true
	case "玩家召唤":
		return AIActionTriggerType_PlayerSummon, true

	}
	return -1, false

}

func (e AIActionExitType) String() string {

	if e == AIActionExitType_Empty {
		return "Empty"
	}
	if e == AIActionExitType_Times {
		return "Times"
	}
	if e == AIActionExitType_EpicHeroesNum {
		return "EpicHeroesNum"
	}

	return ""

}

func GetAIActionExitType(tp string) (AIActionExitType, bool) {
	switch tp {
	case "Empty":
		return AIActionExitType_Empty, true
	case "Times":
		return AIActionExitType_Times, true
	case "EpicHeroesNum":
		return AIActionExitType_EpicHeroesNum, true

	}
	return -1, false

}

func GetAIActionExitType_ByAlias(tp string) (AIActionExitType, bool) {
	switch tp {
	case "无条件":
		return AIActionExitType_Empty, true
	case "执行次数":
		return AIActionExitType_Times, true
	case "有多少个神话英雄":
		return AIActionExitType_EpicHeroesNum, true

	}
	return -1, false

}
