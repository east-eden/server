package config

// 英雄表-拳皇/heroEnum
type HeroEquipPos int

const (
	HeroEquipPos_BEGIN  HeroEquipPos = 0 //枚举开始
	HeroEquipPos_Weapon HeroEquipPos = 0 //武器
	HeroEquipPos_Helmet HeroEquipPos = 1 //头盔
	HeroEquipPos_Chest  HeroEquipPos = 2 //胸甲
	HeroEquipPos_Boots  HeroEquipPos = 3 //靴子
	HeroEquipPos_END    HeroEquipPos = 4 //枚举结束
)

// 英雄表-拳皇/heroEnum
type BonusTargetType int

const (
	BonusTargetType_BEGIN           BonusTargetType = 0  //枚举开始
	BonusTargetType_All             BonusTargetType = 0  //全部英雄
	BonusTargetType_Career          BonusTargetType = 1  //职业
	BonusTargetType_Personal        BonusTargetType = 2  //自身
	BonusTargetType_FetterTeam      BonusTargetType = 3  //羁绊队伍
	BonusTargetType_TargetHero      BonusTargetType = 4  //目标英雄
	BonusTargetType_CareerAttack    BonusTargetType = 11 //职业攻
	BonusTargetType_CareerDefense   BonusTargetType = 12 //职业防
	BonusTargetType_CareerSkill     BonusTargetType = 13 //职业技
	BonusTargetType_CareerAuxiliary BonusTargetType = 14 //职业辅
	BonusTargetType_END             BonusTargetType = 15 //枚举结束
)

// 英雄表-拳皇/heroEnum
type BattleRoleType int

const (
	BattleRoleType_BEGIN  BattleRoleType = 1 //枚举开始
	BattleRoleType_Melee  BattleRoleType = 1 //近战
	BattleRoleType_Remote BattleRoleType = 2 //远程
	BattleRoleType_END    BattleRoleType = 3 //枚举结束
)

// 英雄表-拳皇/heroEnum
type Career int

const (
	Career_BEGIN     Career = 0 //枚举开始
	Career_Empty     Career = 0 //空职业
	Career_Attack    Career = 1 //攻
	Career_Defense   Career = 2 //防
	Career_Skill     Career = 3 //技
	Career_Auxiliary Career = 4 //辅
	Career_END       Career = 5 //枚举结束
)

// 英雄表-拳皇/heroEnum
type HeroSkillType int

const (
	HeroSkillType_BEGIN    HeroSkillType = 1 //枚举开始
	HeroSkillType_Base     HeroSkillType = 1 //普攻
	HeroSkillType_Passive  HeroSkillType = 2 //被动
	HeroSkillType_Super    HeroSkillType = 3 //必杀
	HeroSkillType_Normal   HeroSkillType = 4 //一般
	HeroSkillType_Artifact HeroSkillType = 5 //神器
	HeroSkillType_Active   HeroSkillType = 6 //主动
	HeroSkillType_END      HeroSkillType = 7 //枚举结束
)

// 英雄表-拳皇/heroEnum
type HeroRightsType int

const (
	HeroRightsType_BEGIN     HeroRightsType = 1 //枚举开始
	HeroRightsType_Skill     HeroRightsType = 1 //技能
	HeroRightsType_Fromation HeroRightsType = 2 //布阵
	HeroRightsType_Show      HeroRightsType = 3 //表演
	HeroRightsType_Finish    HeroRightsType = 4 //终结
	HeroRightsType_Painting  HeroRightsType = 5 //立绘
	HeroRightsType_END       HeroRightsType = 6 //枚举结束
)

// 英雄表-拳皇/heroEnum
type HeroCampType int

const (
	HeroCampType_BEGIN           HeroCampType = 1 //枚举开始
	HeroCampType_SacredTreasures HeroCampType = 1 //三神器
	HeroCampType_Orochi          HeroCampType = 2 //大蛇
	HeroCampType_Neutral         HeroCampType = 3 //中立
	HeroCampType_END             HeroCampType = 4 //枚举结束
)

// 英雄表-拳皇/heroEnum
type HeroMarkType int

const (
	HeroMarkType_BEGIN  HeroMarkType = 1 //枚举开始
	HeroMarkType_Baoji  HeroMarkType = 1 //暴击流
	HeroMarkType_Lianji HeroMarkType = 2 //连击流
	HeroMarkType_Gedang HeroMarkType = 3 //格挡流
	HeroMarkType_Jijiu  HeroMarkType = 4 //急救流
	HeroMarkType_Fanji  HeroMarkType = 5 //反击流
	HeroMarkType_Xixue  HeroMarkType = 6 //吸血流
	HeroMarkType_Suoyou HeroMarkType = 7 //所有
	HeroMarkType_END    HeroMarkType = 8 //枚举结束
)

func (e HeroEquipPos) String() string {

	if e == HeroEquipPos_Weapon {
		return "Weapon"
	}
	if e == HeroEquipPos_Helmet {
		return "Helmet"
	}
	if e == HeroEquipPos_Chest {
		return "Chest"
	}
	if e == HeroEquipPos_Boots {
		return "Boots"
	}

	return ""

}

func GetHeroEquipPos(tp string) (HeroEquipPos, bool) {
	switch tp {
	case "Weapon":
		return HeroEquipPos_Weapon, true
	case "Helmet":
		return HeroEquipPos_Helmet, true
	case "Chest":
		return HeroEquipPos_Chest, true
	case "Boots":
		return HeroEquipPos_Boots, true

	}
	return -1, false

}

func GetHeroEquipPos_ByAlias(tp string) (HeroEquipPos, bool) {
	switch tp {
	case "武器":
		return HeroEquipPos_Weapon, true
	case "头盔":
		return HeroEquipPos_Helmet, true
	case "胸甲":
		return HeroEquipPos_Chest, true
	case "靴子":
		return HeroEquipPos_Boots, true

	}
	return -1, false

}

func (e BonusTargetType) String() string {

	if e == BonusTargetType_All {
		return "All"
	}
	if e == BonusTargetType_Career {
		return "Career"
	}
	if e == BonusTargetType_Personal {
		return "Personal"
	}
	if e == BonusTargetType_FetterTeam {
		return "FetterTeam"
	}
	if e == BonusTargetType_TargetHero {
		return "TargetHero"
	}
	if e == BonusTargetType_CareerAttack {
		return "CareerAttack"
	}
	if e == BonusTargetType_CareerDefense {
		return "CareerDefense"
	}
	if e == BonusTargetType_CareerSkill {
		return "CareerSkill"
	}
	if e == BonusTargetType_CareerAuxiliary {
		return "CareerAuxiliary"
	}

	return ""

}

func GetBonusTargetType(tp string) (BonusTargetType, bool) {
	switch tp {
	case "All":
		return BonusTargetType_All, true
	case "Career":
		return BonusTargetType_Career, true
	case "Personal":
		return BonusTargetType_Personal, true
	case "FetterTeam":
		return BonusTargetType_FetterTeam, true
	case "TargetHero":
		return BonusTargetType_TargetHero, true
	case "CareerAttack":
		return BonusTargetType_CareerAttack, true
	case "CareerDefense":
		return BonusTargetType_CareerDefense, true
	case "CareerSkill":
		return BonusTargetType_CareerSkill, true
	case "CareerAuxiliary":
		return BonusTargetType_CareerAuxiliary, true

	}
	return -1, false

}

func GetBonusTargetType_ByAlias(tp string) (BonusTargetType, bool) {
	switch tp {
	case "全部英雄":
		return BonusTargetType_All, true
	case "职业":
		return BonusTargetType_Career, true
	case "自身":
		return BonusTargetType_Personal, true
	case "羁绊队伍":
		return BonusTargetType_FetterTeam, true
	case "目标英雄":
		return BonusTargetType_TargetHero, true
	case "职业攻":
		return BonusTargetType_CareerAttack, true
	case "职业防":
		return BonusTargetType_CareerDefense, true
	case "职业技":
		return BonusTargetType_CareerSkill, true
	case "职业辅":
		return BonusTargetType_CareerAuxiliary, true

	}
	return -1, false

}

func (e BattleRoleType) String() string {

	if e == BattleRoleType_Melee {
		return "Melee"
	}
	if e == BattleRoleType_Remote {
		return "Remote"
	}

	return ""

}

func GetBattleRoleType(tp string) (BattleRoleType, bool) {
	switch tp {
	case "Melee":
		return BattleRoleType_Melee, true
	case "Remote":
		return BattleRoleType_Remote, true

	}
	return -1, false

}

func GetBattleRoleType_ByAlias(tp string) (BattleRoleType, bool) {
	switch tp {
	case "近战":
		return BattleRoleType_Melee, true
	case "远程":
		return BattleRoleType_Remote, true

	}
	return -1, false

}

func (e Career) String() string {

	if e == Career_Empty {
		return "Empty"
	}
	if e == Career_Attack {
		return "Attack"
	}
	if e == Career_Defense {
		return "Defense"
	}
	if e == Career_Skill {
		return "Skill"
	}
	if e == Career_Auxiliary {
		return "Auxiliary"
	}

	return ""

}

func GetCareer(tp string) (Career, bool) {
	switch tp {
	case "Empty":
		return Career_Empty, true
	case "Attack":
		return Career_Attack, true
	case "Defense":
		return Career_Defense, true
	case "Skill":
		return Career_Skill, true
	case "Auxiliary":
		return Career_Auxiliary, true

	}
	return -1, false

}

func GetCareer_ByAlias(tp string) (Career, bool) {
	switch tp {
	case "空职业":
		return Career_Empty, true
	case "攻":
		return Career_Attack, true
	case "防":
		return Career_Defense, true
	case "技":
		return Career_Skill, true
	case "辅":
		return Career_Auxiliary, true

	}
	return -1, false

}

func (e HeroSkillType) String() string {

	if e == HeroSkillType_Base {
		return "Base"
	}
	if e == HeroSkillType_Passive {
		return "Passive"
	}
	if e == HeroSkillType_Super {
		return "Super"
	}
	if e == HeroSkillType_Normal {
		return "Normal"
	}
	if e == HeroSkillType_Artifact {
		return "Artifact"
	}
	if e == HeroSkillType_Active {
		return "Active"
	}

	return ""

}

func GetHeroSkillType(tp string) (HeroSkillType, bool) {
	switch tp {
	case "Base":
		return HeroSkillType_Base, true
	case "Passive":
		return HeroSkillType_Passive, true
	case "Super":
		return HeroSkillType_Super, true
	case "Normal":
		return HeroSkillType_Normal, true
	case "Artifact":
		return HeroSkillType_Artifact, true
	case "Active":
		return HeroSkillType_Active, true

	}
	return -1, false

}

func GetHeroSkillType_ByAlias(tp string) (HeroSkillType, bool) {
	switch tp {
	case "普攻":
		return HeroSkillType_Base, true
	case "被动":
		return HeroSkillType_Passive, true
	case "必杀":
		return HeroSkillType_Super, true
	case "一般":
		return HeroSkillType_Normal, true
	case "神器":
		return HeroSkillType_Artifact, true
	case "主动":
		return HeroSkillType_Active, true

	}
	return -1, false

}

func (e HeroRightsType) String() string {

	if e == HeroRightsType_Skill {
		return "Skill"
	}
	if e == HeroRightsType_Fromation {
		return "Fromation"
	}
	if e == HeroRightsType_Show {
		return "Show"
	}
	if e == HeroRightsType_Finish {
		return "Finish"
	}
	if e == HeroRightsType_Painting {
		return "Painting"
	}

	return ""

}

func GetHeroRightsType(tp string) (HeroRightsType, bool) {
	switch tp {
	case "Skill":
		return HeroRightsType_Skill, true
	case "Fromation":
		return HeroRightsType_Fromation, true
	case "Show":
		return HeroRightsType_Show, true
	case "Finish":
		return HeroRightsType_Finish, true
	case "Painting":
		return HeroRightsType_Painting, true

	}
	return -1, false

}

func GetHeroRightsType_ByAlias(tp string) (HeroRightsType, bool) {
	switch tp {
	case "技能":
		return HeroRightsType_Skill, true
	case "布阵":
		return HeroRightsType_Fromation, true
	case "表演":
		return HeroRightsType_Show, true
	case "终结":
		return HeroRightsType_Finish, true
	case "立绘":
		return HeroRightsType_Painting, true

	}
	return -1, false

}

func (e HeroCampType) String() string {

	if e == HeroCampType_SacredTreasures {
		return "SacredTreasures"
	}
	if e == HeroCampType_Orochi {
		return "Orochi"
	}
	if e == HeroCampType_Neutral {
		return "Neutral"
	}

	return ""

}

func GetHeroCampType(tp string) (HeroCampType, bool) {
	switch tp {
	case "SacredTreasures":
		return HeroCampType_SacredTreasures, true
	case "Orochi":
		return HeroCampType_Orochi, true
	case "Neutral":
		return HeroCampType_Neutral, true

	}
	return -1, false

}

func GetHeroCampType_ByAlias(tp string) (HeroCampType, bool) {
	switch tp {
	case "三神器":
		return HeroCampType_SacredTreasures, true
	case "大蛇":
		return HeroCampType_Orochi, true
	case "中立":
		return HeroCampType_Neutral, true

	}
	return -1, false

}

func (e HeroMarkType) String() string {

	if e == HeroMarkType_Baoji {
		return "Baoji"
	}
	if e == HeroMarkType_Lianji {
		return "Lianji"
	}
	if e == HeroMarkType_Gedang {
		return "Gedang"
	}
	if e == HeroMarkType_Jijiu {
		return "Jijiu"
	}
	if e == HeroMarkType_Fanji {
		return "Fanji"
	}
	if e == HeroMarkType_Xixue {
		return "Xixue"
	}
	if e == HeroMarkType_Suoyou {
		return "Suoyou"
	}

	return ""

}

func GetHeroMarkType(tp string) (HeroMarkType, bool) {
	switch tp {
	case "Baoji":
		return HeroMarkType_Baoji, true
	case "Lianji":
		return HeroMarkType_Lianji, true
	case "Gedang":
		return HeroMarkType_Gedang, true
	case "Jijiu":
		return HeroMarkType_Jijiu, true
	case "Fanji":
		return HeroMarkType_Fanji, true
	case "Xixue":
		return HeroMarkType_Xixue, true
	case "Suoyou":
		return HeroMarkType_Suoyou, true

	}
	return -1, false

}

func GetHeroMarkType_ByAlias(tp string) (HeroMarkType, bool) {
	switch tp {
	case "暴击流":
		return HeroMarkType_Baoji, true
	case "连击流":
		return HeroMarkType_Lianji, true
	case "格挡流":
		return HeroMarkType_Gedang, true
	case "急救流":
		return HeroMarkType_Jijiu, true
	case "反击流":
		return HeroMarkType_Fanji, true
	case "吸血流":
		return HeroMarkType_Xixue, true
	case "所有":
		return HeroMarkType_Suoyou, true

	}
	return -1, false

}
