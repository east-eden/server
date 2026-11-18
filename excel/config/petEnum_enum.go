package config

// 宠物表-拳皇/petEnum
type EMPetType int

const (
	EMPetType_BEGIN    EMPetType = 0 //枚举开始
	EMPetType_Normal   EMPetType = 0 //普通宠物
	EMPetType_Slime    EMPetType = 1 //史莱姆
	EMPetType_Special  EMPetType = 2 //谜团
	EMPetType_Optional EMPetType = 3 //自选宠物
	EMPetType_END      EMPetType = 4 //枚举结束
)

// 宠物表-拳皇/petEnum
type PetRarity int

const (
	PetRarity_BEGIN PetRarity = 0 //枚举开始
	PetRarity_Gray  PetRarity = 0 //灰
	PetRarity_END   PetRarity = 1 //枚举结束
)

// 宠物表-拳皇/petEnum
type EMPetTargetType int

const (
	EMPetTargetType_BEGIN         EMPetTargetType = 0 //枚举开始
	EMPetTargetType_Empty         EMPetTargetType = 0 //空
	EMPetTargetType_Crit          EMPetTargetType = 1 //暴击
	EMPetTargetType_Counterattack EMPetTargetType = 2 //反击
	EMPetTargetType_Block         EMPetTargetType = 3 //格挡
	EMPetTargetType_FirstAid      EMPetTargetType = 4 //急救
	EMPetTargetType_Combo         EMPetTargetType = 5 //连击
	EMPetTargetType_BloodSucking  EMPetTargetType = 6 //吸血
	EMPetTargetType_General       EMPetTargetType = 7 //通用
	EMPetTargetType_END           EMPetTargetType = 8 //枚举结束
)

func (e EMPetType) String() string {

	if e == EMPetType_Normal {
		return "Normal"
	}
	if e == EMPetType_Slime {
		return "Slime"
	}
	if e == EMPetType_Special {
		return "Special"
	}
	if e == EMPetType_Optional {
		return "Optional"
	}

	return ""

}

func GetEMPetType(tp string) (EMPetType, bool) {
	switch tp {
	case "Normal":
		return EMPetType_Normal, true
	case "Slime":
		return EMPetType_Slime, true
	case "Special":
		return EMPetType_Special, true
	case "Optional":
		return EMPetType_Optional, true

	}
	return -1, false

}

func GetEMPetType_ByAlias(tp string) (EMPetType, bool) {
	switch tp {
	case "普通宠物":
		return EMPetType_Normal, true
	case "史莱姆":
		return EMPetType_Slime, true
	case "谜团":
		return EMPetType_Special, true
	case "自选宠物":
		return EMPetType_Optional, true

	}
	return -1, false

}

func (e PetRarity) String() string {

	if e == PetRarity_Gray {
		return "Gray"
	}

	return ""

}

func GetPetRarity(tp string) (PetRarity, bool) {
	switch tp {
	case "Gray":
		return PetRarity_Gray, true

	}
	return -1, false

}

func GetPetRarity_ByAlias(tp string) (PetRarity, bool) {
	switch tp {
	case "灰":
		return PetRarity_Gray, true

	}
	return -1, false

}

func (e EMPetTargetType) String() string {

	if e == EMPetTargetType_Empty {
		return "Empty"
	}
	if e == EMPetTargetType_Crit {
		return "Crit"
	}
	if e == EMPetTargetType_Counterattack {
		return "Counterattack"
	}
	if e == EMPetTargetType_Block {
		return "Block"
	}
	if e == EMPetTargetType_FirstAid {
		return "FirstAid"
	}
	if e == EMPetTargetType_Combo {
		return "Combo"
	}
	if e == EMPetTargetType_BloodSucking {
		return "BloodSucking"
	}
	if e == EMPetTargetType_General {
		return "General"
	}

	return ""

}

func GetEMPetTargetType(tp string) (EMPetTargetType, bool) {
	switch tp {
	case "Empty":
		return EMPetTargetType_Empty, true
	case "Crit":
		return EMPetTargetType_Crit, true
	case "Counterattack":
		return EMPetTargetType_Counterattack, true
	case "Block":
		return EMPetTargetType_Block, true
	case "FirstAid":
		return EMPetTargetType_FirstAid, true
	case "Combo":
		return EMPetTargetType_Combo, true
	case "BloodSucking":
		return EMPetTargetType_BloodSucking, true
	case "General":
		return EMPetTargetType_General, true

	}
	return -1, false

}

func GetEMPetTargetType_ByAlias(tp string) (EMPetTargetType, bool) {
	switch tp {
	case "空":
		return EMPetTargetType_Empty, true
	case "暴击":
		return EMPetTargetType_Crit, true
	case "反击":
		return EMPetTargetType_Counterattack, true
	case "格挡":
		return EMPetTargetType_Block, true
	case "急救":
		return EMPetTargetType_FirstAid, true
	case "连击":
		return EMPetTargetType_Combo, true
	case "吸血":
		return EMPetTargetType_BloodSucking, true
	case "通用":
		return EMPetTargetType_General, true

	}
	return -1, false

}
