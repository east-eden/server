package config

// LD关卡表-拳皇/ServerEffectEnum
type ServerEffectType int

const (
	ServerEffectType_BEGIN                     ServerEffectType = 0    //枚举开始
	ServerEffectType_Empty                     ServerEffectType = 0    //空
	ServerEffectType_MonsterMaxLimit           ServerEffectType = 1    //怪物计数上限
	ServerEffectType_WaveAddCoin               ServerEffectType = 2    //每回合发放硬币
	ServerEffectType_WaveAddPercentCoin        ServerEffectType = 3    //每回合发放持有硬币百分比
	ServerEffectType_StartAddCoin              ServerEffectType = 4    //对局开始时发放金币
	ServerEffectType_AddLuckyStonePer10Wave    ServerEffectType = 5    //每10回合发放幸运石
	ServerEffectType_GamblePurpleProbAdd       ServerEffectType = 6    //赌博时英雄兵种（紫色）召唤概率增加
	ServerEffectType_GambleLegendProbAdd       ServerEffectType = 7    //赌博时神话兵种（金色）召唤概率增加
	ServerEffectType_SummonRareProbAdd         ServerEffectType = 8    //增加抽取稀有品质英雄概率
	ServerEffectType_SellHeroGetLuckyStoneProb ServerEffectType = 9    //出售兵种时有概率额外获得1个幸运石
	ServerEffectType_FinishQuestGetLuckyStone  ServerEffectType = 10   //完成关卡内任务时获得幸运石
	ServerEffectType_GambleFailGuarantee       ServerEffectType = 11   //赌博失败时有概率发放本次赌博低一阶段的兵种
	ServerEffectType_GambleFailReturn          ServerEffectType = 12   //赌博失败时有概率返还赌博费用
	ServerEffectType_HeroMaxLimitAdd           ServerEffectType = 13   //可持有的兵种上限数量增加
	ServerEffectType_SummonHero                ServerEffectType = 14   //召唤英雄
	ServerEffectType_KillAddCoin               ServerEffectType = 15   //击杀时硬币增加
	ServerEffectType_CriticalHit               ServerEffectType = 16   //触发暴击
	ServerEffectType_GambleBlueProbAdd         ServerEffectType = 17   //赌博时英雄兵种（蓝色）召唤概率增加
	ServerEffectType_Heal                      ServerEffectType = 18   //治疗
	ServerEffectType_ClientAddHeroAtk          ServerEffectType = 9901 //增加英雄攻击
	ServerEffectType_ClientMoneyAddAtk         ServerEffectType = 9902 //金币增加攻击
	ServerEffectType_ClientMagicCrit           ServerEffectType = 9903 //增加魔法暴击
	ServerEffectType_ClientRetunSp             ServerEffectType = 9904 //返还SP
	ServerEffectType_END                       ServerEffectType = 9905 //枚举结束
)

func (e ServerEffectType) String() string {

	if e == ServerEffectType_Empty {
		return "Empty"
	}
	if e == ServerEffectType_MonsterMaxLimit {
		return "MonsterMaxLimit"
	}
	if e == ServerEffectType_WaveAddCoin {
		return "WaveAddCoin"
	}
	if e == ServerEffectType_WaveAddPercentCoin {
		return "WaveAddPercentCoin"
	}
	if e == ServerEffectType_StartAddCoin {
		return "StartAddCoin"
	}
	if e == ServerEffectType_AddLuckyStonePer10Wave {
		return "AddLuckyStonePer10Wave"
	}
	if e == ServerEffectType_GamblePurpleProbAdd {
		return "GamblePurpleProbAdd"
	}
	if e == ServerEffectType_GambleLegendProbAdd {
		return "GambleLegendProbAdd"
	}
	if e == ServerEffectType_SummonRareProbAdd {
		return "SummonRareProbAdd"
	}
	if e == ServerEffectType_SellHeroGetLuckyStoneProb {
		return "SellHeroGetLuckyStoneProb"
	}
	if e == ServerEffectType_FinishQuestGetLuckyStone {
		return "FinishQuestGetLuckyStone"
	}
	if e == ServerEffectType_GambleFailGuarantee {
		return "GambleFailGuarantee"
	}
	if e == ServerEffectType_GambleFailReturn {
		return "GambleFailReturn"
	}
	if e == ServerEffectType_HeroMaxLimitAdd {
		return "HeroMaxLimitAdd"
	}
	if e == ServerEffectType_SummonHero {
		return "SummonHero"
	}
	if e == ServerEffectType_KillAddCoin {
		return "KillAddCoin"
	}
	if e == ServerEffectType_CriticalHit {
		return "CriticalHit"
	}
	if e == ServerEffectType_GambleBlueProbAdd {
		return "GambleBlueProbAdd"
	}
	if e == ServerEffectType_Heal {
		return "Heal"
	}
	if e == ServerEffectType_ClientAddHeroAtk {
		return "ClientAddHeroAtk"
	}
	if e == ServerEffectType_ClientMoneyAddAtk {
		return "ClientMoneyAddAtk"
	}
	if e == ServerEffectType_ClientMagicCrit {
		return "ClientMagicCrit"
	}
	if e == ServerEffectType_ClientRetunSp {
		return "ClientRetunSp"
	}

	return ""

}

func GetServerEffectType(tp string) (ServerEffectType, bool) {
	switch tp {
	case "Empty":
		return ServerEffectType_Empty, true
	case "MonsterMaxLimit":
		return ServerEffectType_MonsterMaxLimit, true
	case "WaveAddCoin":
		return ServerEffectType_WaveAddCoin, true
	case "WaveAddPercentCoin":
		return ServerEffectType_WaveAddPercentCoin, true
	case "StartAddCoin":
		return ServerEffectType_StartAddCoin, true
	case "AddLuckyStonePer10Wave":
		return ServerEffectType_AddLuckyStonePer10Wave, true
	case "GamblePurpleProbAdd":
		return ServerEffectType_GamblePurpleProbAdd, true
	case "GambleLegendProbAdd":
		return ServerEffectType_GambleLegendProbAdd, true
	case "SummonRareProbAdd":
		return ServerEffectType_SummonRareProbAdd, true
	case "SellHeroGetLuckyStoneProb":
		return ServerEffectType_SellHeroGetLuckyStoneProb, true
	case "FinishQuestGetLuckyStone":
		return ServerEffectType_FinishQuestGetLuckyStone, true
	case "GambleFailGuarantee":
		return ServerEffectType_GambleFailGuarantee, true
	case "GambleFailReturn":
		return ServerEffectType_GambleFailReturn, true
	case "HeroMaxLimitAdd":
		return ServerEffectType_HeroMaxLimitAdd, true
	case "SummonHero":
		return ServerEffectType_SummonHero, true
	case "KillAddCoin":
		return ServerEffectType_KillAddCoin, true
	case "CriticalHit":
		return ServerEffectType_CriticalHit, true
	case "GambleBlueProbAdd":
		return ServerEffectType_GambleBlueProbAdd, true
	case "Heal":
		return ServerEffectType_Heal, true
	case "ClientAddHeroAtk":
		return ServerEffectType_ClientAddHeroAtk, true
	case "ClientMoneyAddAtk":
		return ServerEffectType_ClientMoneyAddAtk, true
	case "ClientMagicCrit":
		return ServerEffectType_ClientMagicCrit, true
	case "ClientRetunSp":
		return ServerEffectType_ClientRetunSp, true

	}
	return -1, false

}

func GetServerEffectType_ByAlias(tp string) (ServerEffectType, bool) {
	switch tp {
	case "空":
		return ServerEffectType_Empty, true
	case "怪物计数上限":
		return ServerEffectType_MonsterMaxLimit, true
	case "每回合发放硬币":
		return ServerEffectType_WaveAddCoin, true
	case "每回合发放持有硬币百分比":
		return ServerEffectType_WaveAddPercentCoin, true
	case "对局开始时发放金币":
		return ServerEffectType_StartAddCoin, true
	case "每10回合发放幸运石":
		return ServerEffectType_AddLuckyStonePer10Wave, true
	case "赌博时英雄兵种（紫色）召唤概率增加":
		return ServerEffectType_GamblePurpleProbAdd, true
	case "赌博时神话兵种（金色）召唤概率增加":
		return ServerEffectType_GambleLegendProbAdd, true
	case "增加抽取稀有品质英雄概率":
		return ServerEffectType_SummonRareProbAdd, true
	case "出售兵种时有概率额外获得1个幸运石":
		return ServerEffectType_SellHeroGetLuckyStoneProb, true
	case "完成关卡内任务时获得幸运石":
		return ServerEffectType_FinishQuestGetLuckyStone, true
	case "赌博失败时有概率发放本次赌博低一阶段的兵种":
		return ServerEffectType_GambleFailGuarantee, true
	case "赌博失败时有概率返还赌博费用":
		return ServerEffectType_GambleFailReturn, true
	case "可持有的兵种上限数量增加":
		return ServerEffectType_HeroMaxLimitAdd, true
	case "召唤英雄":
		return ServerEffectType_SummonHero, true
	case "击杀时硬币增加":
		return ServerEffectType_KillAddCoin, true
	case "触发暴击":
		return ServerEffectType_CriticalHit, true
	case "赌博时英雄兵种（蓝色）召唤概率增加":
		return ServerEffectType_GambleBlueProbAdd, true
	case "治疗":
		return ServerEffectType_Heal, true
	case "增加英雄攻击":
		return ServerEffectType_ClientAddHeroAtk, true
	case "金币增加攻击":
		return ServerEffectType_ClientMoneyAddAtk, true
	case "增加魔法暴击":
		return ServerEffectType_ClientMagicCrit, true
	case "返还SP":
		return ServerEffectType_ClientRetunSp, true

	}
	return -1, false

}
