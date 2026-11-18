package config

// SpTD关卡表-拳皇/SpTDEnum
type AbilityCheckType int

const (
	AbilityCheckType_BEGIN   AbilityCheckType = 1 //枚举开始
	AbilityCheckType_Carrer  AbilityCheckType = 1 //职业
	AbilityCheckType_Team    AbilityCheckType = 2 //队伍
	AbilityCheckType_TeamCnt AbilityCheckType = 3 //队伍人数
	AbilityCheckType_HeroCnt AbilityCheckType = 4 //英雄在场
	AbilityCheckType_END     AbilityCheckType = 5 //枚举结束
)

// SpTD关卡表-拳皇/SpTDEnum
type SpLevelGeneratorType int

const (
	SpLevelGeneratorType_BEGIN  SpLevelGeneratorType = 1 //枚举开始
	SpLevelGeneratorType_Normal SpLevelGeneratorType = 1 //主线
	SpLevelGeneratorType_Boss   SpLevelGeneratorType = 2 //Boss
	SpLevelGeneratorType_END    SpLevelGeneratorType = 3 //枚举结束
)

// SpTD关卡表-拳皇/SpTDEnum
type EntranceActionType int

const (
	EntranceActionType_BEGIN           EntranceActionType = -1 //枚举开始
	EntranceActionType_None            EntranceActionType = -1 //无
	EntranceActionType_Avoid           EntranceActionType = 1  //入场迫使躲避
	EntranceActionType_LeaveAttackBack EntranceActionType = 2  //离场击退
	EntranceActionType_LeaveAttackFly  EntranceActionType = 3  //离场击飞
	EntranceActionType_END             EntranceActionType = 4  //枚举结束
)

// SpTD关卡表-拳皇/SpTDEnum
type EnemyType int

const (
	EnemyType_BEGIN   EnemyType = 0 //枚举开始
	EnemyType_Sheild  EnemyType = 0 //护盾流
	EnemyType_HitStun EnemyType = 1 //推推流
	EnemyType_Counter EnemyType = 2 //反击流
	EnemyType_Sniper  EnemyType = 3 //后排流
	EnemyType_Combo   EnemyType = 4 //连击流
	EnemyType_END     EnemyType = 5 //枚举结束
)

func (e AbilityCheckType) String() string {

	if e == AbilityCheckType_Carrer {
		return "Carrer"
	}
	if e == AbilityCheckType_Team {
		return "Team"
	}
	if e == AbilityCheckType_TeamCnt {
		return "TeamCnt"
	}
	if e == AbilityCheckType_HeroCnt {
		return "HeroCnt"
	}

	return ""

}

func GetAbilityCheckType(tp string) (AbilityCheckType, bool) {
	switch tp {
	case "Carrer":
		return AbilityCheckType_Carrer, true
	case "Team":
		return AbilityCheckType_Team, true
	case "TeamCnt":
		return AbilityCheckType_TeamCnt, true
	case "HeroCnt":
		return AbilityCheckType_HeroCnt, true

	}
	return -1, false

}

func GetAbilityCheckType_ByAlias(tp string) (AbilityCheckType, bool) {
	switch tp {
	case "职业":
		return AbilityCheckType_Carrer, true
	case "队伍":
		return AbilityCheckType_Team, true
	case "队伍人数":
		return AbilityCheckType_TeamCnt, true
	case "英雄在场":
		return AbilityCheckType_HeroCnt, true

	}
	return -1, false

}

func (e SpLevelGeneratorType) String() string {

	if e == SpLevelGeneratorType_Normal {
		return "Normal"
	}
	if e == SpLevelGeneratorType_Boss {
		return "Boss"
	}

	return ""

}

func GetSpLevelGeneratorType(tp string) (SpLevelGeneratorType, bool) {
	switch tp {
	case "Normal":
		return SpLevelGeneratorType_Normal, true
	case "Boss":
		return SpLevelGeneratorType_Boss, true

	}
	return -1, false

}

func GetSpLevelGeneratorType_ByAlias(tp string) (SpLevelGeneratorType, bool) {
	switch tp {
	case "主线":
		return SpLevelGeneratorType_Normal, true
	case "Boss":
		return SpLevelGeneratorType_Boss, true

	}
	return -1, false

}

func (e EntranceActionType) String() string {

	if e == EntranceActionType_None {
		return "None"
	}
	if e == EntranceActionType_Avoid {
		return "Avoid"
	}
	if e == EntranceActionType_LeaveAttackBack {
		return "LeaveAttackBack"
	}
	if e == EntranceActionType_LeaveAttackFly {
		return "LeaveAttackFly"
	}

	return ""

}

func GetEntranceActionType(tp string) (EntranceActionType, bool) {
	switch tp {
	case "None":
		return EntranceActionType_None, true
	case "Avoid":
		return EntranceActionType_Avoid, true
	case "LeaveAttackBack":
		return EntranceActionType_LeaveAttackBack, true
	case "LeaveAttackFly":
		return EntranceActionType_LeaveAttackFly, true

	}
	return -1, false

}

func GetEntranceActionType_ByAlias(tp string) (EntranceActionType, bool) {
	switch tp {
	case "无":
		return EntranceActionType_None, true
	case "入场迫使躲避":
		return EntranceActionType_Avoid, true
	case "离场击退":
		return EntranceActionType_LeaveAttackBack, true
	case "离场击飞":
		return EntranceActionType_LeaveAttackFly, true

	}
	return -1, false

}

func (e EnemyType) String() string {

	if e == EnemyType_Sheild {
		return "Sheild"
	}
	if e == EnemyType_HitStun {
		return "HitStun"
	}
	if e == EnemyType_Counter {
		return "Counter"
	}
	if e == EnemyType_Sniper {
		return "Sniper"
	}
	if e == EnemyType_Combo {
		return "Combo"
	}

	return ""

}

func GetEnemyType(tp string) (EnemyType, bool) {
	switch tp {
	case "Sheild":
		return EnemyType_Sheild, true
	case "HitStun":
		return EnemyType_HitStun, true
	case "Counter":
		return EnemyType_Counter, true
	case "Sniper":
		return EnemyType_Sniper, true
	case "Combo":
		return EnemyType_Combo, true

	}
	return -1, false

}

func GetEnemyType_ByAlias(tp string) (EnemyType, bool) {
	switch tp {
	case "护盾流":
		return EnemyType_Sheild, true
	case "推推流":
		return EnemyType_HitStun, true
	case "反击流":
		return EnemyType_Counter, true
	case "后排流":
		return EnemyType_Sniper, true
	case "连击流":
		return EnemyType_Combo, true

	}
	return -1, false

}
