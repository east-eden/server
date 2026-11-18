package config

// 拳皇大明星-拳皇/kofStarType
type KofStarCareer int

const (
	KofStarCareer_BEGIN     KofStarCareer = 0 //枚举开始
	KofStarCareer_Empty     KofStarCareer = 0 //空
	KofStarCareer_Attack    KofStarCareer = 1 //攻
	KofStarCareer_Defense   KofStarCareer = 2 //防
	KofStarCareer_Skill     KofStarCareer = 3 //技
	KofStarCareer_Auxiliary KofStarCareer = 4 //辅
	KofStarCareer_Special   KofStarCareer = 5 //特
	KofStarCareer_END       KofStarCareer = 6 //枚举结束
)

// 拳皇大明星-拳皇/kofStarType
type KofStarRandomType int

const (
	KofStarRandomType_BEGIN KofStarRandomType = 1 //枚举开始
	KofStarRandomType_A     KofStarRandomType = 1 //音巢科技
	KofStarRandomType_B     KofStarRandomType = 2 //武道秘法
	KofStarRandomType_END   KofStarRandomType = 3 //枚举结束
)

// 拳皇大明星-拳皇/kofStarType
type KofStarSkillOwnerType int

const (
	KofStarSkillOwnerType_BEGIN  KofStarSkillOwnerType = 1 //枚举开始
	KofStarSkillOwnerType_Hero   KofStarSkillOwnerType = 1 //英雄
	KofStarSkillOwnerType_Career KofStarSkillOwnerType = 2 //职业
	KofStarSkillOwnerType_Global KofStarSkillOwnerType = 3 //全局
	KofStarSkillOwnerType_END    KofStarSkillOwnerType = 4 //枚举结束
)

// 拳皇大明星-拳皇/kofStarType
type KofStarSkillType int

const (
	KofStarSkillType_BEGIN          KofStarSkillType = 1 //枚举开始
	KofStarSkillType_CareeBase      KofStarSkillType = 1 //职业基础
	KofStarSkillType_CareerCore     KofStarSkillType = 2 //职业核心
	KofStarSkillType_CareerNegative KofStarSkillType = 3 //职业被动
	KofStarSkillType_HeroSpecial    KofStarSkillType = 4 //英雄专属
	KofStarSkillType_END            KofStarSkillType = 5 //枚举结束
)

// 拳皇大明星-拳皇/kofStarType
type KofRogueCardType int

const (
	KofRogueCardType_BEGIN KofRogueCardType = 1 //枚举开始
	KofRogueCardType_Hero  KofRogueCardType = 1 //英雄
	KofRogueCardType_Skill KofRogueCardType = 2 //技能
	KofRogueCardType_END   KofRogueCardType = 3 //枚举结束
)

func (e KofStarCareer) String() string {

	if e == KofStarCareer_Empty {
		return "Empty"
	}
	if e == KofStarCareer_Attack {
		return "Attack"
	}
	if e == KofStarCareer_Defense {
		return "Defense"
	}
	if e == KofStarCareer_Skill {
		return "Skill"
	}
	if e == KofStarCareer_Auxiliary {
		return "Auxiliary"
	}
	if e == KofStarCareer_Special {
		return "Special"
	}

	return ""

}

func GetKofStarCareer(tp string) (KofStarCareer, bool) {
	switch tp {
	case "Empty":
		return KofStarCareer_Empty, true
	case "Attack":
		return KofStarCareer_Attack, true
	case "Defense":
		return KofStarCareer_Defense, true
	case "Skill":
		return KofStarCareer_Skill, true
	case "Auxiliary":
		return KofStarCareer_Auxiliary, true
	case "Special":
		return KofStarCareer_Special, true

	}
	return -1, false

}

func GetKofStarCareer_ByAlias(tp string) (KofStarCareer, bool) {
	switch tp {
	case "空":
		return KofStarCareer_Empty, true
	case "攻":
		return KofStarCareer_Attack, true
	case "防":
		return KofStarCareer_Defense, true
	case "技":
		return KofStarCareer_Skill, true
	case "辅":
		return KofStarCareer_Auxiliary, true
	case "特":
		return KofStarCareer_Special, true

	}
	return -1, false

}

func (e KofStarRandomType) String() string {

	if e == KofStarRandomType_A {
		return "A"
	}
	if e == KofStarRandomType_B {
		return "B"
	}

	return ""

}

func GetKofStarRandomType(tp string) (KofStarRandomType, bool) {
	switch tp {
	case "A":
		return KofStarRandomType_A, true
	case "B":
		return KofStarRandomType_B, true

	}
	return -1, false

}

func GetKofStarRandomType_ByAlias(tp string) (KofStarRandomType, bool) {
	switch tp {
	case "音巢科技":
		return KofStarRandomType_A, true
	case "武道秘法":
		return KofStarRandomType_B, true

	}
	return -1, false

}

func (e KofStarSkillOwnerType) String() string {

	if e == KofStarSkillOwnerType_Hero {
		return "Hero"
	}
	if e == KofStarSkillOwnerType_Career {
		return "Career"
	}
	if e == KofStarSkillOwnerType_Global {
		return "Global"
	}

	return ""

}

func GetKofStarSkillOwnerType(tp string) (KofStarSkillOwnerType, bool) {
	switch tp {
	case "Hero":
		return KofStarSkillOwnerType_Hero, true
	case "Career":
		return KofStarSkillOwnerType_Career, true
	case "Global":
		return KofStarSkillOwnerType_Global, true

	}
	return -1, false

}

func GetKofStarSkillOwnerType_ByAlias(tp string) (KofStarSkillOwnerType, bool) {
	switch tp {
	case "英雄":
		return KofStarSkillOwnerType_Hero, true
	case "职业":
		return KofStarSkillOwnerType_Career, true
	case "全局":
		return KofStarSkillOwnerType_Global, true

	}
	return -1, false

}

func (e KofStarSkillType) String() string {

	if e == KofStarSkillType_CareeBase {
		return "CareeBase"
	}
	if e == KofStarSkillType_CareerCore {
		return "CareerCore"
	}
	if e == KofStarSkillType_CareerNegative {
		return "CareerNegative"
	}
	if e == KofStarSkillType_HeroSpecial {
		return "HeroSpecial"
	}

	return ""

}

func GetKofStarSkillType(tp string) (KofStarSkillType, bool) {
	switch tp {
	case "CareeBase":
		return KofStarSkillType_CareeBase, true
	case "CareerCore":
		return KofStarSkillType_CareerCore, true
	case "CareerNegative":
		return KofStarSkillType_CareerNegative, true
	case "HeroSpecial":
		return KofStarSkillType_HeroSpecial, true

	}
	return -1, false

}

func GetKofStarSkillType_ByAlias(tp string) (KofStarSkillType, bool) {
	switch tp {
	case "职业基础":
		return KofStarSkillType_CareeBase, true
	case "职业核心":
		return KofStarSkillType_CareerCore, true
	case "职业被动":
		return KofStarSkillType_CareerNegative, true
	case "英雄专属":
		return KofStarSkillType_HeroSpecial, true

	}
	return -1, false

}

func (e KofRogueCardType) String() string {

	if e == KofRogueCardType_Hero {
		return "Hero"
	}
	if e == KofRogueCardType_Skill {
		return "Skill"
	}

	return ""

}

func GetKofRogueCardType(tp string) (KofRogueCardType, bool) {
	switch tp {
	case "Hero":
		return KofRogueCardType_Hero, true
	case "Skill":
		return KofRogueCardType_Skill, true

	}
	return -1, false

}

func GetKofRogueCardType_ByAlias(tp string) (KofRogueCardType, bool) {
	switch tp {
	case "英雄":
		return KofRogueCardType_Hero, true
	case "技能":
		return KofRogueCardType_Skill, true

	}
	return -1, false

}
