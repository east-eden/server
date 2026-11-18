package config

// 最强108-拳皇/Strongest108Enum
type Strongest108Tier int

const (
	Strongest108Tier_BEGIN        Strongest108Tier = 1  //枚举开始
	Strongest108Tier_Iron         Strongest108Tier = 1  //黑铁
	Strongest108Tier_Bronze       Strongest108Tier = 2  //青铜
	Strongest108Tier_Silver       Strongest108Tier = 3  //白银
	Strongest108Tier_Gold         Strongest108Tier = 4  //黄金
	Strongest108Tier_Platinum     Strongest108Tier = 5  //铂金
	Strongest108Tier_Emerald      Strongest108Tier = 6  //翡翠
	Strongest108Tier_Diamond      Strongest108Tier = 7  //钻石
	Strongest108Tier_Master       Strongest108Tier = 8  //大师
	Strongest108Tier_Grandmaster  Strongest108Tier = 9  //宗师
	Strongest108Tier_Challenger   Strongest108Tier = 10 //王者
	Strongest108Tier_Strongest108 Strongest108Tier = 11 //最强108
	Strongest108Tier_END          Strongest108Tier = 12 //枚举结束
)

// 最强108-拳皇/Strongest108Enum
type Strongest108Progress int

const (
	Strongest108Progress_BEGIN      Strongest108Progress = 1 //枚举开始
	Strongest108Progress_Register   Strongest108Progress = 1 //报名
	Strongest108Progress_Matching   Strongest108Progress = 2 //匹配
	Strongest108Progress_Ranking    Strongest108Progress = 3 //排位
	Strongest108Progress_Settlement Strongest108Progress = 4 //结算领奖
	Strongest108Progress_Closed     Strongest108Progress = 5 //未开启
	Strongest108Progress_END        Strongest108Progress = 6 //枚举结束
)

// 最强108-拳皇/Strongest108Enum
type Strongest108AwardShowType int

const (
	Strongest108AwardShowType_BEGIN     Strongest108AwardShowType = 1 //枚举开始
	Strongest108AwardShowType_One       Strongest108AwardShowType = 1 //第一名
	Strongest108AwardShowType_Two       Strongest108AwardShowType = 2 //第二名
	Strongest108AwardShowType_Three     Strongest108AwardShowType = 3 //第三名
	Strongest108AwardShowType_Promotion Strongest108AwardShowType = 4 //晋级
	Strongest108AwardShowType_Normal    Strongest108AwardShowType = 5 //平级
	Strongest108AwardShowType_Downgrade Strongest108AwardShowType = 6 //降级
	Strongest108AwardShowType_END       Strongest108AwardShowType = 7 //枚举结束
)

func (e Strongest108Tier) String() string {

	if e == Strongest108Tier_Iron {
		return "Iron"
	}
	if e == Strongest108Tier_Bronze {
		return "Bronze"
	}
	if e == Strongest108Tier_Silver {
		return "Silver"
	}
	if e == Strongest108Tier_Gold {
		return "Gold"
	}
	if e == Strongest108Tier_Platinum {
		return "Platinum"
	}
	if e == Strongest108Tier_Emerald {
		return "Emerald"
	}
	if e == Strongest108Tier_Diamond {
		return "Diamond"
	}
	if e == Strongest108Tier_Master {
		return "Master"
	}
	if e == Strongest108Tier_Grandmaster {
		return "Grandmaster"
	}
	if e == Strongest108Tier_Challenger {
		return "Challenger"
	}
	if e == Strongest108Tier_Strongest108 {
		return "Strongest108"
	}

	return ""

}

func GetStrongest108Tier(tp string) (Strongest108Tier, bool) {
	switch tp {
	case "Iron":
		return Strongest108Tier_Iron, true
	case "Bronze":
		return Strongest108Tier_Bronze, true
	case "Silver":
		return Strongest108Tier_Silver, true
	case "Gold":
		return Strongest108Tier_Gold, true
	case "Platinum":
		return Strongest108Tier_Platinum, true
	case "Emerald":
		return Strongest108Tier_Emerald, true
	case "Diamond":
		return Strongest108Tier_Diamond, true
	case "Master":
		return Strongest108Tier_Master, true
	case "Grandmaster":
		return Strongest108Tier_Grandmaster, true
	case "Challenger":
		return Strongest108Tier_Challenger, true
	case "Strongest108":
		return Strongest108Tier_Strongest108, true

	}
	return -1, false

}

func GetStrongest108Tier_ByAlias(tp string) (Strongest108Tier, bool) {
	switch tp {
	case "黑铁":
		return Strongest108Tier_Iron, true
	case "青铜":
		return Strongest108Tier_Bronze, true
	case "白银":
		return Strongest108Tier_Silver, true
	case "黄金":
		return Strongest108Tier_Gold, true
	case "铂金":
		return Strongest108Tier_Platinum, true
	case "翡翠":
		return Strongest108Tier_Emerald, true
	case "钻石":
		return Strongest108Tier_Diamond, true
	case "大师":
		return Strongest108Tier_Master, true
	case "宗师":
		return Strongest108Tier_Grandmaster, true
	case "王者":
		return Strongest108Tier_Challenger, true
	case "最强108":
		return Strongest108Tier_Strongest108, true

	}
	return -1, false

}

func (e Strongest108Progress) String() string {

	if e == Strongest108Progress_Register {
		return "Register"
	}
	if e == Strongest108Progress_Matching {
		return "Matching"
	}
	if e == Strongest108Progress_Ranking {
		return "Ranking"
	}
	if e == Strongest108Progress_Settlement {
		return "Settlement"
	}
	if e == Strongest108Progress_Closed {
		return "Closed"
	}

	return ""

}

func GetStrongest108Progress(tp string) (Strongest108Progress, bool) {
	switch tp {
	case "Register":
		return Strongest108Progress_Register, true
	case "Matching":
		return Strongest108Progress_Matching, true
	case "Ranking":
		return Strongest108Progress_Ranking, true
	case "Settlement":
		return Strongest108Progress_Settlement, true
	case "Closed":
		return Strongest108Progress_Closed, true

	}
	return -1, false

}

func GetStrongest108Progress_ByAlias(tp string) (Strongest108Progress, bool) {
	switch tp {
	case "报名":
		return Strongest108Progress_Register, true
	case "匹配":
		return Strongest108Progress_Matching, true
	case "排位":
		return Strongest108Progress_Ranking, true
	case "结算领奖":
		return Strongest108Progress_Settlement, true
	case "未开启":
		return Strongest108Progress_Closed, true

	}
	return -1, false

}

func (e Strongest108AwardShowType) String() string {

	if e == Strongest108AwardShowType_One {
		return "One"
	}
	if e == Strongest108AwardShowType_Two {
		return "Two"
	}
	if e == Strongest108AwardShowType_Three {
		return "Three"
	}
	if e == Strongest108AwardShowType_Promotion {
		return "Promotion"
	}
	if e == Strongest108AwardShowType_Normal {
		return "Normal"
	}
	if e == Strongest108AwardShowType_Downgrade {
		return "Downgrade"
	}

	return ""

}

func GetStrongest108AwardShowType(tp string) (Strongest108AwardShowType, bool) {
	switch tp {
	case "One":
		return Strongest108AwardShowType_One, true
	case "Two":
		return Strongest108AwardShowType_Two, true
	case "Three":
		return Strongest108AwardShowType_Three, true
	case "Promotion":
		return Strongest108AwardShowType_Promotion, true
	case "Normal":
		return Strongest108AwardShowType_Normal, true
	case "Downgrade":
		return Strongest108AwardShowType_Downgrade, true

	}
	return -1, false

}

func GetStrongest108AwardShowType_ByAlias(tp string) (Strongest108AwardShowType, bool) {
	switch tp {
	case "第一名":
		return Strongest108AwardShowType_One, true
	case "第二名":
		return Strongest108AwardShowType_Two, true
	case "第三名":
		return Strongest108AwardShowType_Three, true
	case "晋级":
		return Strongest108AwardShowType_Promotion, true
	case "平级":
		return Strongest108AwardShowType_Normal, true
	case "降级":
		return Strongest108AwardShowType_Downgrade, true

	}
	return -1, false

}
