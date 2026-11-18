package config

// LD技能表-拳皇/LDSkillEnum
type LDSkillShowUIType int

const (
	LDSkillShowUIType_BEGIN LDSkillShowUIType = 0 //枚举开始
	LDSkillShowUIType_None  LDSkillShowUIType = 0 //不显示
	LDSkillShowUIType_Troop LDSkillShowUIType = 1 //仅阵容显示
	LDSkillShowUIType_Hero  LDSkillShowUIType = 2 //仅培养显示
	LDSkillShowUIType_All   LDSkillShowUIType = 3 //显示
	LDSkillShowUIType_END   LDSkillShowUIType = 4 //枚举结束
)

// LD技能表-拳皇/LDSkillEnum
type SkillDescAttrType int

const (
	SkillDescAttrType_BEGIN   SkillDescAttrType = 1 //枚举开始
	SkillDescAttrType_Physics SkillDescAttrType = 1 //物理伤害
	SkillDescAttrType_Magic   SkillDescAttrType = 2 //魔法伤害
	SkillDescAttrType_Stun    SkillDescAttrType = 3 //眩晕
	SkillDescAttrType_END     SkillDescAttrType = 4 //枚举结束
)

// LD技能表-拳皇/LDSkillEnum
type LDSkillIconType int

const (
	LDSkillIconType_BEGIN     LDSkillIconType = 1 //枚举开始
	LDSkillIconType_Proactive LDSkillIconType = 1 //技
	LDSkillIconType_Passively LDSkillIconType = 2 //被
	LDSkillIconType_Ultimate  LDSkillIconType = 3 //大
	LDSkillIconType_Normal    LDSkillIconType = 4 //普
	LDSkillIconType_Positive  LDSkillIconType = 5 //主
	LDSkillIconType_Fetters   LDSkillIconType = 8 //羁
	LDSkillIconType_END       LDSkillIconType = 9 //枚举结束
)

func (e LDSkillShowUIType) String() string {

	if e == LDSkillShowUIType_None {
		return "None"
	}
	if e == LDSkillShowUIType_Troop {
		return "Troop"
	}
	if e == LDSkillShowUIType_Hero {
		return "Hero"
	}
	if e == LDSkillShowUIType_All {
		return "All"
	}

	return ""

}

func GetLDSkillShowUIType(tp string) (LDSkillShowUIType, bool) {
	switch tp {
	case "None":
		return LDSkillShowUIType_None, true
	case "Troop":
		return LDSkillShowUIType_Troop, true
	case "Hero":
		return LDSkillShowUIType_Hero, true
	case "All":
		return LDSkillShowUIType_All, true

	}
	return -1, false

}

func GetLDSkillShowUIType_ByAlias(tp string) (LDSkillShowUIType, bool) {
	switch tp {
	case "不显示":
		return LDSkillShowUIType_None, true
	case "仅阵容显示":
		return LDSkillShowUIType_Troop, true
	case "仅培养显示":
		return LDSkillShowUIType_Hero, true
	case "显示":
		return LDSkillShowUIType_All, true

	}
	return -1, false

}

func (e SkillDescAttrType) String() string {

	if e == SkillDescAttrType_Physics {
		return "Physics"
	}
	if e == SkillDescAttrType_Magic {
		return "Magic"
	}
	if e == SkillDescAttrType_Stun {
		return "Stun"
	}

	return ""

}

func GetSkillDescAttrType(tp string) (SkillDescAttrType, bool) {
	switch tp {
	case "Physics":
		return SkillDescAttrType_Physics, true
	case "Magic":
		return SkillDescAttrType_Magic, true
	case "Stun":
		return SkillDescAttrType_Stun, true

	}
	return -1, false

}

func GetSkillDescAttrType_ByAlias(tp string) (SkillDescAttrType, bool) {
	switch tp {
	case "物理伤害":
		return SkillDescAttrType_Physics, true
	case "魔法伤害":
		return SkillDescAttrType_Magic, true
	case "眩晕":
		return SkillDescAttrType_Stun, true

	}
	return -1, false

}

func (e LDSkillIconType) String() string {

	if e == LDSkillIconType_Proactive {
		return "Proactive"
	}
	if e == LDSkillIconType_Passively {
		return "Passively"
	}
	if e == LDSkillIconType_Ultimate {
		return "Ultimate"
	}
	if e == LDSkillIconType_Normal {
		return "Normal"
	}
	if e == LDSkillIconType_Positive {
		return "Positive"
	}
	if e == LDSkillIconType_Fetters {
		return "Fetters"
	}

	return ""

}

func GetLDSkillIconType(tp string) (LDSkillIconType, bool) {
	switch tp {
	case "Proactive":
		return LDSkillIconType_Proactive, true
	case "Passively":
		return LDSkillIconType_Passively, true
	case "Ultimate":
		return LDSkillIconType_Ultimate, true
	case "Normal":
		return LDSkillIconType_Normal, true
	case "Positive":
		return LDSkillIconType_Positive, true
	case "Fetters":
		return LDSkillIconType_Fetters, true

	}
	return -1, false

}

func GetLDSkillIconType_ByAlias(tp string) (LDSkillIconType, bool) {
	switch tp {
	case "技":
		return LDSkillIconType_Proactive, true
	case "被":
		return LDSkillIconType_Passively, true
	case "大":
		return LDSkillIconType_Ultimate, true
	case "普":
		return LDSkillIconType_Normal, true
	case "主":
		return LDSkillIconType_Positive, true
	case "羁":
		return LDSkillIconType_Fetters, true

	}
	return -1, false

}
