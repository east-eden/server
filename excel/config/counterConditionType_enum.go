package config

// 计数任务表-拳皇/CounterConditionType
type CounterConditionType int

const (
	CounterConditionType_BEGIN                 CounterConditionType = 1  //枚举开始
	CounterConditionType_HeroSoulCompose       CounterConditionType = 1  //八门合成
	CounterConditionType_PlayerLevel           CounterConditionType = 2  //玩家等级
	CounterConditionType_TrialChallenge        CounterConditionType = 3  //试炼挑战
	CounterConditionType_ArenaWin              CounterConditionType = 4  //竞技场获胜
	CounterConditionType_MainStagePass         CounterConditionType = 5  //主线通关
	CounterConditionType_AllianceBossChallenge CounterConditionType = 6  //公会boss挑战
	CounterConditionType_FishRedTask           CounterConditionType = 7  //钓鱼品质任务
	CounterConditionType_ManorResourse         CounterConditionType = 8  //领地资源
	CounterConditionType_ManorPlunder          CounterConditionType = 9  //领地掠夺
	CounterConditionType_AllianceBuild         CounterConditionType = 10 //公会捐献
	CounterConditionType_ShopBuy               CounterConditionType = 11 //商店购买
	CounterConditionType_AfkQuick              CounterConditionType = 12 //快速赛事奖励
	CounterConditionType_CostItem              CounterConditionType = 13 //消耗X道具Y个
	CounterConditionType_ElfTowerChallenge     CounterConditionType = 14 //南镇挑战次数
	CounterConditionType_CostItemBySubType     CounterConditionType = 15 //按类型消耗道具
	CounterConditionType_AllianceBargain       CounterConditionType = 16 //公会砍价
	CounterConditionType_AllianceBargainBuy    CounterConditionType = 17 //公会砍价购买
	CounterConditionType_ManorEvictEnemy       CounterConditionType = 18 //阻击对方掠夺
	CounterConditionType_WatchAd               CounterConditionType = 19 //看广告
	CounterConditionType_Recharge              CounterConditionType = 20 //充值金额
	CounterConditionType_HeroReplace           CounterConditionType = 21 //无损替换
	CounterConditionType_TDMatch               CounterConditionType = 22 //TD匹配次数
	CounterConditionType_WorldBossAttack       CounterConditionType = 23 //封印之地次数
	CounterConditionType_AllianceHelp          CounterConditionType = 24 //帮助公会成员加速
	CounterConditionType_AllianceAidSuccess    CounterConditionType = 25 //援助公会成员
	CounterConditionType_ArtifactAstrology     CounterConditionType = 26 //祈福次数
	CounterConditionType_GainItem              CounterConditionType = 27 //获得X道具Y个
	CounterConditionType_GainItemBySubType     CounterConditionType = 28 //按类型获得道具
	CounterConditionType_HeroUpStarById        CounterConditionType = 29 //指定ID卡牌升星
	CounterConditionType_HeroUpStarByRarity    CounterConditionType = 30 //按所见品质卡牌升星次数
	CounterConditionType_ArtifactUpStar        CounterConditionType = 31 //神器升星次数
	CounterConditionType_FetterLvUp            CounterConditionType = 32 //羁绊升级次数
	CounterConditionType_BoxGiftByRarity       CounterConditionType = 33 //获得指定品质手办次数
	CounterConditionType_RecruitCnt            CounterConditionType = 34 //招募次数
	CounterConditionType_WorldBossDamage       CounterConditionType = 35 //世界boss伤害
	CounterConditionType_LDPass                CounterConditionType = 36 //td波数通关
	CounterConditionType_ElfTowerKill          CounterConditionType = 37 //南镇击杀
	CounterConditionType_FIshPhysicalStrength  CounterConditionType = 38 //钓鱼体力
	CounterConditionType_TowerPassFloor        CounterConditionType = 39 //吉斯塔挑战层
	CounterConditionType_StageLevelPass        CounterConditionType = 40 //主线挑战关卡
	CounterConditionType_FightPowerUp          CounterConditionType = 41 //战力提升
	CounterConditionType_AllianceTechGacha     CounterConditionType = 42 //公会科技修炼
	CounterConditionType_END                   CounterConditionType = 43 //枚举结束
)

func (e CounterConditionType) String() string {

	if e == CounterConditionType_HeroSoulCompose {
		return "HeroSoulCompose"
	}
	if e == CounterConditionType_PlayerLevel {
		return "PlayerLevel"
	}
	if e == CounterConditionType_TrialChallenge {
		return "TrialChallenge"
	}
	if e == CounterConditionType_ArenaWin {
		return "ArenaWin"
	}
	if e == CounterConditionType_MainStagePass {
		return "MainStagePass"
	}
	if e == CounterConditionType_AllianceBossChallenge {
		return "AllianceBossChallenge"
	}
	if e == CounterConditionType_FishRedTask {
		return "FishRedTask"
	}
	if e == CounterConditionType_ManorResourse {
		return "ManorResourse"
	}
	if e == CounterConditionType_ManorPlunder {
		return "ManorPlunder"
	}
	if e == CounterConditionType_AllianceBuild {
		return "AllianceBuild"
	}
	if e == CounterConditionType_ShopBuy {
		return "ShopBuy"
	}
	if e == CounterConditionType_AfkQuick {
		return "AfkQuick"
	}
	if e == CounterConditionType_CostItem {
		return "CostItem"
	}
	if e == CounterConditionType_ElfTowerChallenge {
		return "ElfTowerChallenge"
	}
	if e == CounterConditionType_CostItemBySubType {
		return "CostItemBySubType"
	}
	if e == CounterConditionType_AllianceBargain {
		return "AllianceBargain"
	}
	if e == CounterConditionType_AllianceBargainBuy {
		return "AllianceBargainBuy"
	}
	if e == CounterConditionType_ManorEvictEnemy {
		return "ManorEvictEnemy"
	}
	if e == CounterConditionType_WatchAd {
		return "WatchAd"
	}
	if e == CounterConditionType_Recharge {
		return "Recharge"
	}
	if e == CounterConditionType_HeroReplace {
		return "HeroReplace"
	}
	if e == CounterConditionType_TDMatch {
		return "TDMatch"
	}
	if e == CounterConditionType_WorldBossAttack {
		return "WorldBossAttack"
	}
	if e == CounterConditionType_AllianceHelp {
		return "AllianceHelp"
	}
	if e == CounterConditionType_AllianceAidSuccess {
		return "AllianceAidSuccess"
	}
	if e == CounterConditionType_ArtifactAstrology {
		return "ArtifactAstrology"
	}
	if e == CounterConditionType_GainItem {
		return "GainItem"
	}
	if e == CounterConditionType_GainItemBySubType {
		return "GainItemBySubType"
	}
	if e == CounterConditionType_HeroUpStarById {
		return "HeroUpStarById"
	}
	if e == CounterConditionType_HeroUpStarByRarity {
		return "HeroUpStarByRarity"
	}
	if e == CounterConditionType_ArtifactUpStar {
		return "ArtifactUpStar"
	}
	if e == CounterConditionType_FetterLvUp {
		return "FetterLvUp"
	}
	if e == CounterConditionType_BoxGiftByRarity {
		return "BoxGiftByRarity"
	}
	if e == CounterConditionType_RecruitCnt {
		return "RecruitCnt"
	}
	if e == CounterConditionType_WorldBossDamage {
		return "WorldBossDamage"
	}
	if e == CounterConditionType_LDPass {
		return "LDPass"
	}
	if e == CounterConditionType_ElfTowerKill {
		return "ElfTowerKill"
	}
	if e == CounterConditionType_FIshPhysicalStrength {
		return "FIshPhysicalStrength"
	}
	if e == CounterConditionType_TowerPassFloor {
		return "TowerPassFloor"
	}
	if e == CounterConditionType_StageLevelPass {
		return "StageLevelPass"
	}
	if e == CounterConditionType_FightPowerUp {
		return "FightPowerUp"
	}
	if e == CounterConditionType_AllianceTechGacha {
		return "AllianceTechGacha"
	}

	return ""

}

func GetCounterConditionType(tp string) (CounterConditionType, bool) {
	switch tp {
	case "HeroSoulCompose":
		return CounterConditionType_HeroSoulCompose, true
	case "PlayerLevel":
		return CounterConditionType_PlayerLevel, true
	case "TrialChallenge":
		return CounterConditionType_TrialChallenge, true
	case "ArenaWin":
		return CounterConditionType_ArenaWin, true
	case "MainStagePass":
		return CounterConditionType_MainStagePass, true
	case "AllianceBossChallenge":
		return CounterConditionType_AllianceBossChallenge, true
	case "FishRedTask":
		return CounterConditionType_FishRedTask, true
	case "ManorResourse":
		return CounterConditionType_ManorResourse, true
	case "ManorPlunder":
		return CounterConditionType_ManorPlunder, true
	case "AllianceBuild":
		return CounterConditionType_AllianceBuild, true
	case "ShopBuy":
		return CounterConditionType_ShopBuy, true
	case "AfkQuick":
		return CounterConditionType_AfkQuick, true
	case "CostItem":
		return CounterConditionType_CostItem, true
	case "ElfTowerChallenge":
		return CounterConditionType_ElfTowerChallenge, true
	case "CostItemBySubType":
		return CounterConditionType_CostItemBySubType, true
	case "AllianceBargain":
		return CounterConditionType_AllianceBargain, true
	case "AllianceBargainBuy":
		return CounterConditionType_AllianceBargainBuy, true
	case "ManorEvictEnemy":
		return CounterConditionType_ManorEvictEnemy, true
	case "WatchAd":
		return CounterConditionType_WatchAd, true
	case "Recharge":
		return CounterConditionType_Recharge, true
	case "HeroReplace":
		return CounterConditionType_HeroReplace, true
	case "TDMatch":
		return CounterConditionType_TDMatch, true
	case "WorldBossAttack":
		return CounterConditionType_WorldBossAttack, true
	case "AllianceHelp":
		return CounterConditionType_AllianceHelp, true
	case "AllianceAidSuccess":
		return CounterConditionType_AllianceAidSuccess, true
	case "ArtifactAstrology":
		return CounterConditionType_ArtifactAstrology, true
	case "GainItem":
		return CounterConditionType_GainItem, true
	case "GainItemBySubType":
		return CounterConditionType_GainItemBySubType, true
	case "HeroUpStarById":
		return CounterConditionType_HeroUpStarById, true
	case "HeroUpStarByRarity":
		return CounterConditionType_HeroUpStarByRarity, true
	case "ArtifactUpStar":
		return CounterConditionType_ArtifactUpStar, true
	case "FetterLvUp":
		return CounterConditionType_FetterLvUp, true
	case "BoxGiftByRarity":
		return CounterConditionType_BoxGiftByRarity, true
	case "RecruitCnt":
		return CounterConditionType_RecruitCnt, true
	case "WorldBossDamage":
		return CounterConditionType_WorldBossDamage, true
	case "LDPass":
		return CounterConditionType_LDPass, true
	case "ElfTowerKill":
		return CounterConditionType_ElfTowerKill, true
	case "FIshPhysicalStrength":
		return CounterConditionType_FIshPhysicalStrength, true
	case "TowerPassFloor":
		return CounterConditionType_TowerPassFloor, true
	case "StageLevelPass":
		return CounterConditionType_StageLevelPass, true
	case "FightPowerUp":
		return CounterConditionType_FightPowerUp, true
	case "AllianceTechGacha":
		return CounterConditionType_AllianceTechGacha, true

	}
	return -1, false

}

func GetCounterConditionType_ByAlias(tp string) (CounterConditionType, bool) {
	switch tp {
	case "八门合成":
		return CounterConditionType_HeroSoulCompose, true
	case "玩家等级":
		return CounterConditionType_PlayerLevel, true
	case "试炼挑战":
		return CounterConditionType_TrialChallenge, true
	case "竞技场获胜":
		return CounterConditionType_ArenaWin, true
	case "主线通关":
		return CounterConditionType_MainStagePass, true
	case "公会boss挑战":
		return CounterConditionType_AllianceBossChallenge, true
	case "钓鱼品质任务":
		return CounterConditionType_FishRedTask, true
	case "领地资源":
		return CounterConditionType_ManorResourse, true
	case "领地掠夺":
		return CounterConditionType_ManorPlunder, true
	case "公会捐献":
		return CounterConditionType_AllianceBuild, true
	case "商店购买":
		return CounterConditionType_ShopBuy, true
	case "快速赛事奖励":
		return CounterConditionType_AfkQuick, true
	case "消耗X道具Y个":
		return CounterConditionType_CostItem, true
	case "南镇挑战次数":
		return CounterConditionType_ElfTowerChallenge, true
	case "按类型消耗道具":
		return CounterConditionType_CostItemBySubType, true
	case "公会砍价":
		return CounterConditionType_AllianceBargain, true
	case "公会砍价购买":
		return CounterConditionType_AllianceBargainBuy, true
	case "阻击对方掠夺":
		return CounterConditionType_ManorEvictEnemy, true
	case "看广告":
		return CounterConditionType_WatchAd, true
	case "充值金额":
		return CounterConditionType_Recharge, true
	case "无损替换":
		return CounterConditionType_HeroReplace, true
	case "TD匹配次数":
		return CounterConditionType_TDMatch, true
	case "封印之地次数":
		return CounterConditionType_WorldBossAttack, true
	case "帮助公会成员加速":
		return CounterConditionType_AllianceHelp, true
	case "援助公会成员":
		return CounterConditionType_AllianceAidSuccess, true
	case "祈福次数":
		return CounterConditionType_ArtifactAstrology, true
	case "获得X道具Y个":
		return CounterConditionType_GainItem, true
	case "按类型获得道具":
		return CounterConditionType_GainItemBySubType, true
	case "指定ID卡牌升星":
		return CounterConditionType_HeroUpStarById, true
	case "按所见品质卡牌升星次数":
		return CounterConditionType_HeroUpStarByRarity, true
	case "神器升星次数":
		return CounterConditionType_ArtifactUpStar, true
	case "羁绊升级次数":
		return CounterConditionType_FetterLvUp, true
	case "获得指定品质手办次数":
		return CounterConditionType_BoxGiftByRarity, true
	case "招募次数":
		return CounterConditionType_RecruitCnt, true
	case "世界boss伤害":
		return CounterConditionType_WorldBossDamage, true
	case "td波数通关":
		return CounterConditionType_LDPass, true
	case "南镇击杀":
		return CounterConditionType_ElfTowerKill, true
	case "钓鱼体力":
		return CounterConditionType_FIshPhysicalStrength, true
	case "吉斯塔挑战层":
		return CounterConditionType_TowerPassFloor, true
	case "主线挑战关卡":
		return CounterConditionType_StageLevelPass, true
	case "战力提升":
		return CounterConditionType_FightPowerUp, true
	case "公会科技修炼":
		return CounterConditionType_AllianceTechGacha, true

	}
	return -1, false

}
