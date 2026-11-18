package config

// 联盟攻城战-拳皇/AllianceWarEnum
type AllianceWarProgress int

const (
	AllianceWarProgress_BEGIN   AllianceWarProgress = 1  //枚举开始
	AllianceWarProgress_Signup  AllianceWarProgress = 1  //报名期
	AllianceWarProgress_Prepare AllianceWarProgress = 2  //备战期
	AllianceWarProgress_Fight1  AllianceWarProgress = 3  //进行中·第1轮
	AllianceWarProgress_Rest1   AllianceWarProgress = 4  //休战期·第1轮后
	AllianceWarProgress_Fight2  AllianceWarProgress = 5  //进行中·第2轮
	AllianceWarProgress_Rest2   AllianceWarProgress = 6  //休战期·第2轮后(竞猜窗口)
	AllianceWarProgress_Fight3  AllianceWarProgress = 7  //进行中·第3轮
	AllianceWarProgress_Reward  AllianceWarProgress = 8  //领奖期
	AllianceWarProgress_Closed  AllianceWarProgress = 9  //关闭状态
	AllianceWarProgress_END     AllianceWarProgress = 10 //枚举结束
)

// 联盟攻城战-拳皇/AllianceWarEnum
type AllianceWarBattlefield int

const (
	AllianceWarBattlefield_BEGIN    AllianceWarBattlefield = 1 //枚举开始
	AllianceWarBattlefield_Beginner AllianceWarBattlefield = 1 //初级战场
	AllianceWarBattlefield_Middle   AllianceWarBattlefield = 2 //中级战场
	AllianceWarBattlefield_Senior   AllianceWarBattlefield = 3 //高级战场
	AllianceWarBattlefield_END      AllianceWarBattlefield = 4 //枚举结束
)

// 联盟攻城战-拳皇/AllianceWarEnum
type AllianceWarDonateType int

const (
	AllianceWarDonateType_BEGIN AllianceWarDonateType = 1 //枚举开始
	AllianceWarDonateType_Gem   AllianceWarDonateType = 1 //钻石捐献
	AllianceWarDonateType_Item  AllianceWarDonateType = 2 //道具捐献
	AllianceWarDonateType_END   AllianceWarDonateType = 3 //枚举结束
)

// 联盟攻城战-拳皇/AllianceWarEnum
type AllianceWarGuessHit int

const (
	AllianceWarGuessHit_BEGIN       AllianceWarGuessHit = 1 //枚举开始
	AllianceWarGuessHit_All3        AllianceWarGuessHit = 1 //成功猜对冠亚季军
	AllianceWarGuessHit_ChampRunner AllianceWarGuessHit = 2 //成功猜对冠亚军
	AllianceWarGuessHit_Champ       AllianceWarGuessHit = 3 //成功猜对冠军
	AllianceWarGuessHit_RunnerUp    AllianceWarGuessHit = 4 //成功猜对亚军
	AllianceWarGuessHit_Third       AllianceWarGuessHit = 5 //成功猜对季军
	AllianceWarGuessHit_END         AllianceWarGuessHit = 6 //枚举结束
)

func (e AllianceWarProgress) String() string {

	if e == AllianceWarProgress_Signup {
		return "Signup"
	}
	if e == AllianceWarProgress_Prepare {
		return "Prepare"
	}
	if e == AllianceWarProgress_Fight1 {
		return "Fight1"
	}
	if e == AllianceWarProgress_Rest1 {
		return "Rest1"
	}
	if e == AllianceWarProgress_Fight2 {
		return "Fight2"
	}
	if e == AllianceWarProgress_Rest2 {
		return "Rest2"
	}
	if e == AllianceWarProgress_Fight3 {
		return "Fight3"
	}
	if e == AllianceWarProgress_Reward {
		return "Reward"
	}
	if e == AllianceWarProgress_Closed {
		return "Closed"
	}

	return ""

}

func GetAllianceWarProgress(tp string) (AllianceWarProgress, bool) {
	switch tp {
	case "Signup":
		return AllianceWarProgress_Signup, true
	case "Prepare":
		return AllianceWarProgress_Prepare, true
	case "Fight1":
		return AllianceWarProgress_Fight1, true
	case "Rest1":
		return AllianceWarProgress_Rest1, true
	case "Fight2":
		return AllianceWarProgress_Fight2, true
	case "Rest2":
		return AllianceWarProgress_Rest2, true
	case "Fight3":
		return AllianceWarProgress_Fight3, true
	case "Reward":
		return AllianceWarProgress_Reward, true
	case "Closed":
		return AllianceWarProgress_Closed, true

	}
	return -1, false

}

func GetAllianceWarProgress_ByAlias(tp string) (AllianceWarProgress, bool) {
	switch tp {
	case "报名期":
		return AllianceWarProgress_Signup, true
	case "备战期":
		return AllianceWarProgress_Prepare, true
	case "进行中·第1轮":
		return AllianceWarProgress_Fight1, true
	case "休战期·第1轮后":
		return AllianceWarProgress_Rest1, true
	case "进行中·第2轮":
		return AllianceWarProgress_Fight2, true
	case "休战期·第2轮后(竞猜窗口)":
		return AllianceWarProgress_Rest2, true
	case "进行中·第3轮":
		return AllianceWarProgress_Fight3, true
	case "领奖期":
		return AllianceWarProgress_Reward, true
	case "关闭状态":
		return AllianceWarProgress_Closed, true

	}
	return -1, false

}

func (e AllianceWarBattlefield) String() string {

	if e == AllianceWarBattlefield_Beginner {
		return "Beginner"
	}
	if e == AllianceWarBattlefield_Middle {
		return "Middle"
	}
	if e == AllianceWarBattlefield_Senior {
		return "Senior"
	}

	return ""

}

func GetAllianceWarBattlefield(tp string) (AllianceWarBattlefield, bool) {
	switch tp {
	case "Beginner":
		return AllianceWarBattlefield_Beginner, true
	case "Middle":
		return AllianceWarBattlefield_Middle, true
	case "Senior":
		return AllianceWarBattlefield_Senior, true

	}
	return -1, false

}

func GetAllianceWarBattlefield_ByAlias(tp string) (AllianceWarBattlefield, bool) {
	switch tp {
	case "初级战场":
		return AllianceWarBattlefield_Beginner, true
	case "中级战场":
		return AllianceWarBattlefield_Middle, true
	case "高级战场":
		return AllianceWarBattlefield_Senior, true

	}
	return -1, false

}

func (e AllianceWarDonateType) String() string {

	if e == AllianceWarDonateType_Gem {
		return "Gem"
	}
	if e == AllianceWarDonateType_Item {
		return "Item"
	}

	return ""

}

func GetAllianceWarDonateType(tp string) (AllianceWarDonateType, bool) {
	switch tp {
	case "Gem":
		return AllianceWarDonateType_Gem, true
	case "Item":
		return AllianceWarDonateType_Item, true

	}
	return -1, false

}

func GetAllianceWarDonateType_ByAlias(tp string) (AllianceWarDonateType, bool) {
	switch tp {
	case "钻石捐献":
		return AllianceWarDonateType_Gem, true
	case "道具捐献":
		return AllianceWarDonateType_Item, true

	}
	return -1, false

}

func (e AllianceWarGuessHit) String() string {

	if e == AllianceWarGuessHit_All3 {
		return "All3"
	}
	if e == AllianceWarGuessHit_ChampRunner {
		return "ChampRunner"
	}
	if e == AllianceWarGuessHit_Champ {
		return "Champ"
	}
	if e == AllianceWarGuessHit_RunnerUp {
		return "RunnerUp"
	}
	if e == AllianceWarGuessHit_Third {
		return "Third"
	}

	return ""

}

func GetAllianceWarGuessHit(tp string) (AllianceWarGuessHit, bool) {
	switch tp {
	case "All3":
		return AllianceWarGuessHit_All3, true
	case "ChampRunner":
		return AllianceWarGuessHit_ChampRunner, true
	case "Champ":
		return AllianceWarGuessHit_Champ, true
	case "RunnerUp":
		return AllianceWarGuessHit_RunnerUp, true
	case "Third":
		return AllianceWarGuessHit_Third, true

	}
	return -1, false

}

func GetAllianceWarGuessHit_ByAlias(tp string) (AllianceWarGuessHit, bool) {
	switch tp {
	case "成功猜对冠亚季军":
		return AllianceWarGuessHit_All3, true
	case "成功猜对冠亚军":
		return AllianceWarGuessHit_ChampRunner, true
	case "成功猜对冠军":
		return AllianceWarGuessHit_Champ, true
	case "成功猜对亚军":
		return AllianceWarGuessHit_RunnerUp, true
	case "成功猜对季军":
		return AllianceWarGuessHit_Third, true

	}
	return -1, false

}
