package config

// Mega招募-拳皇/MegaRecruitEnum
type MegaRecruitType int

const (
	MegaRecruitType_BEGIN          MegaRecruitType = 1 //枚举开始
	MegaRecruitType_Rugal          MegaRecruitType = 1 //暴走卢卡尔
	MegaRecruitType_DaShe          MegaRecruitType = 2 //大蛇
	MegaRecruitType_DaSheBack      MegaRecruitType = 3 //大蛇返场
	MegaRecruitType_Crazy8         MegaRecruitType = 4 //疯八
	MegaRecruitType_RugalBack      MegaRecruitType = 5 //暴走卢卡尔返场
	MegaRecruitType_Crazy8Back     MegaRecruitType = 6 //疯八返场
	MegaRecruitType_CrazyLeona     MegaRecruitType = 7 //狂暴莉安娜
	MegaRecruitType_CrazyLeonaBack MegaRecruitType = 8 //狂暴莉安娜返场
	MegaRecruitType_END            MegaRecruitType = 9 //枚举结束
)

// Mega招募-拳皇/MegaRecruitEnum
type MegaActSubType int

const (
	MegaActSubType_BEGIN      MegaActSubType = 1 //枚举开始
	MegaActSubType_Recruit    MegaActSubType = 1 //招募
	MegaActSubType_Nurture    MegaActSubType = 2 //养成
	MegaActSubType_Sign       MegaActSubType = 3 //签到
	MegaActSubType_Rank       MegaActSubType = 4 //排行
	MegaActSubType_TimePortal MegaActSubType = 5 //时空之门
	MegaActSubType_END        MegaActSubType = 6 //枚举结束
)

func (e MegaRecruitType) String() string {

	if e == MegaRecruitType_Rugal {
		return "Rugal"
	}
	if e == MegaRecruitType_DaShe {
		return "DaShe"
	}
	if e == MegaRecruitType_DaSheBack {
		return "DaSheBack"
	}
	if e == MegaRecruitType_Crazy8 {
		return "Crazy8"
	}
	if e == MegaRecruitType_RugalBack {
		return "RugalBack"
	}
	if e == MegaRecruitType_Crazy8Back {
		return "Crazy8Back"
	}
	if e == MegaRecruitType_CrazyLeona {
		return "CrazyLeona"
	}
	if e == MegaRecruitType_CrazyLeonaBack {
		return "CrazyLeonaBack"
	}

	return ""

}

func GetMegaRecruitType(tp string) (MegaRecruitType, bool) {
	switch tp {
	case "Rugal":
		return MegaRecruitType_Rugal, true
	case "DaShe":
		return MegaRecruitType_DaShe, true
	case "DaSheBack":
		return MegaRecruitType_DaSheBack, true
	case "Crazy8":
		return MegaRecruitType_Crazy8, true
	case "RugalBack":
		return MegaRecruitType_RugalBack, true
	case "Crazy8Back":
		return MegaRecruitType_Crazy8Back, true
	case "CrazyLeona":
		return MegaRecruitType_CrazyLeona, true
	case "CrazyLeonaBack":
		return MegaRecruitType_CrazyLeonaBack, true

	}
	return -1, false

}

func GetMegaRecruitType_ByAlias(tp string) (MegaRecruitType, bool) {
	switch tp {
	case "暴走卢卡尔":
		return MegaRecruitType_Rugal, true
	case "大蛇":
		return MegaRecruitType_DaShe, true
	case "大蛇返场":
		return MegaRecruitType_DaSheBack, true
	case "疯八":
		return MegaRecruitType_Crazy8, true
	case "暴走卢卡尔返场":
		return MegaRecruitType_RugalBack, true
	case "疯八返场":
		return MegaRecruitType_Crazy8Back, true
	case "狂暴莉安娜":
		return MegaRecruitType_CrazyLeona, true
	case "狂暴莉安娜返场":
		return MegaRecruitType_CrazyLeonaBack, true

	}
	return -1, false

}

func (e MegaActSubType) String() string {

	if e == MegaActSubType_Recruit {
		return "Recruit"
	}
	if e == MegaActSubType_Nurture {
		return "Nurture"
	}
	if e == MegaActSubType_Sign {
		return "Sign"
	}
	if e == MegaActSubType_Rank {
		return "Rank"
	}
	if e == MegaActSubType_TimePortal {
		return "TimePortal"
	}

	return ""

}

func GetMegaActSubType(tp string) (MegaActSubType, bool) {
	switch tp {
	case "Recruit":
		return MegaActSubType_Recruit, true
	case "Nurture":
		return MegaActSubType_Nurture, true
	case "Sign":
		return MegaActSubType_Sign, true
	case "Rank":
		return MegaActSubType_Rank, true
	case "TimePortal":
		return MegaActSubType_TimePortal, true

	}
	return -1, false

}

func GetMegaActSubType_ByAlias(tp string) (MegaActSubType, bool) {
	switch tp {
	case "招募":
		return MegaActSubType_Recruit, true
	case "养成":
		return MegaActSubType_Nurture, true
	case "签到":
		return MegaActSubType_Sign, true
	case "排行":
		return MegaActSubType_Rank, true
	case "时空之门":
		return MegaActSubType_TimePortal, true

	}
	return -1, false

}
