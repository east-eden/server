package config

// 跨服竞技32强-拳皇/CrossArenaEnum
type CrossArenaProgress int

const (
	CrossArenaProgress_BEGIN           CrossArenaProgress = 1  //枚举开始
	CrossArenaProgress_WaitRegister    CrossArenaProgress = 1  //等待名单
	CrossArenaProgress_WaitOpen        CrossArenaProgress = 2  //等待开启
	CrossArenaProgress_Cheering        CrossArenaProgress = 3  //助威中
	CrossArenaProgress_Ranking         CrossArenaProgress = 4  //排位赛
	CrossArenaProgress_WaitElimination CrossArenaProgress = 5  //等待淘汰赛
	CrossArenaProgress_WaitTop16       CrossArenaProgress = 6  //等待16强
	CrossArenaProgress_Top16           CrossArenaProgress = 7  //角逐16强
	CrossArenaProgress_WaitTop8        CrossArenaProgress = 8  //等待8强战
	CrossArenaProgress_Top8            CrossArenaProgress = 9  //角逐8强
	CrossArenaProgress_WaitTop4        CrossArenaProgress = 10 //等待4强战
	CrossArenaProgress_Top4            CrossArenaProgress = 11 //角逐4强
	CrossArenaProgress_WaitTop2        CrossArenaProgress = 12 //等待2强战
	CrossArenaProgress_Top2            CrossArenaProgress = 13 //角逐2强
	CrossArenaProgress_WaitTop1        CrossArenaProgress = 14 //等待冠军战
	CrossArenaProgress_Top1            CrossArenaProgress = 15 //角逐冠军
	CrossArenaProgress_Exhibition      CrossArenaProgress = 16 //冠军展示
	CrossArenaProgress_Closed          CrossArenaProgress = 17 //未开启
	CrossArenaProgress_END             CrossArenaProgress = 18 //枚举结束
)

// 跨服竞技32强-拳皇/CrossArenaEnum
type CrossArenaGroup int

const (
	CrossArenaGroup_BEGIN  CrossArenaGroup = 1 //枚举开始
	CrossArenaGroup_Rookie CrossArenaGroup = 1 //新手组
	CrossArenaGroup_Elite  CrossArenaGroup = 2 //精英组
	CrossArenaGroup_Honor  CrossArenaGroup = 3 //荣耀组
	CrossArenaGroup_King   CrossArenaGroup = 4 //王者组
	CrossArenaGroup_END    CrossArenaGroup = 5 //枚举结束
)

func (e CrossArenaProgress) String() string {

	if e == CrossArenaProgress_WaitRegister {
		return "WaitRegister"
	}
	if e == CrossArenaProgress_WaitOpen {
		return "WaitOpen"
	}
	if e == CrossArenaProgress_Cheering {
		return "Cheering"
	}
	if e == CrossArenaProgress_Ranking {
		return "Ranking"
	}
	if e == CrossArenaProgress_WaitElimination {
		return "WaitElimination"
	}
	if e == CrossArenaProgress_WaitTop16 {
		return "WaitTop16"
	}
	if e == CrossArenaProgress_Top16 {
		return "Top16"
	}
	if e == CrossArenaProgress_WaitTop8 {
		return "WaitTop8"
	}
	if e == CrossArenaProgress_Top8 {
		return "Top8"
	}
	if e == CrossArenaProgress_WaitTop4 {
		return "WaitTop4"
	}
	if e == CrossArenaProgress_Top4 {
		return "Top4"
	}
	if e == CrossArenaProgress_WaitTop2 {
		return "WaitTop2"
	}
	if e == CrossArenaProgress_Top2 {
		return "Top2"
	}
	if e == CrossArenaProgress_WaitTop1 {
		return "WaitTop1"
	}
	if e == CrossArenaProgress_Top1 {
		return "Top1"
	}
	if e == CrossArenaProgress_Exhibition {
		return "Exhibition"
	}
	if e == CrossArenaProgress_Closed {
		return "Closed"
	}

	return ""

}

func GetCrossArenaProgress(tp string) (CrossArenaProgress, bool) {
	switch tp {
	case "WaitRegister":
		return CrossArenaProgress_WaitRegister, true
	case "WaitOpen":
		return CrossArenaProgress_WaitOpen, true
	case "Cheering":
		return CrossArenaProgress_Cheering, true
	case "Ranking":
		return CrossArenaProgress_Ranking, true
	case "WaitElimination":
		return CrossArenaProgress_WaitElimination, true
	case "WaitTop16":
		return CrossArenaProgress_WaitTop16, true
	case "Top16":
		return CrossArenaProgress_Top16, true
	case "WaitTop8":
		return CrossArenaProgress_WaitTop8, true
	case "Top8":
		return CrossArenaProgress_Top8, true
	case "WaitTop4":
		return CrossArenaProgress_WaitTop4, true
	case "Top4":
		return CrossArenaProgress_Top4, true
	case "WaitTop2":
		return CrossArenaProgress_WaitTop2, true
	case "Top2":
		return CrossArenaProgress_Top2, true
	case "WaitTop1":
		return CrossArenaProgress_WaitTop1, true
	case "Top1":
		return CrossArenaProgress_Top1, true
	case "Exhibition":
		return CrossArenaProgress_Exhibition, true
	case "Closed":
		return CrossArenaProgress_Closed, true

	}
	return -1, false

}

func GetCrossArenaProgress_ByAlias(tp string) (CrossArenaProgress, bool) {
	switch tp {
	case "等待名单":
		return CrossArenaProgress_WaitRegister, true
	case "等待开启":
		return CrossArenaProgress_WaitOpen, true
	case "助威中":
		return CrossArenaProgress_Cheering, true
	case "排位赛":
		return CrossArenaProgress_Ranking, true
	case "等待淘汰赛":
		return CrossArenaProgress_WaitElimination, true
	case "等待16强":
		return CrossArenaProgress_WaitTop16, true
	case "角逐16强":
		return CrossArenaProgress_Top16, true
	case "等待8强战":
		return CrossArenaProgress_WaitTop8, true
	case "角逐8强":
		return CrossArenaProgress_Top8, true
	case "等待4强战":
		return CrossArenaProgress_WaitTop4, true
	case "角逐4强":
		return CrossArenaProgress_Top4, true
	case "等待2强战":
		return CrossArenaProgress_WaitTop2, true
	case "角逐2强":
		return CrossArenaProgress_Top2, true
	case "等待冠军战":
		return CrossArenaProgress_WaitTop1, true
	case "角逐冠军":
		return CrossArenaProgress_Top1, true
	case "冠军展示":
		return CrossArenaProgress_Exhibition, true
	case "未开启":
		return CrossArenaProgress_Closed, true

	}
	return -1, false

}

func (e CrossArenaGroup) String() string {

	if e == CrossArenaGroup_Rookie {
		return "Rookie"
	}
	if e == CrossArenaGroup_Elite {
		return "Elite"
	}
	if e == CrossArenaGroup_Honor {
		return "Honor"
	}
	if e == CrossArenaGroup_King {
		return "King"
	}

	return ""

}

func GetCrossArenaGroup(tp string) (CrossArenaGroup, bool) {
	switch tp {
	case "Rookie":
		return CrossArenaGroup_Rookie, true
	case "Elite":
		return CrossArenaGroup_Elite, true
	case "Honor":
		return CrossArenaGroup_Honor, true
	case "King":
		return CrossArenaGroup_King, true

	}
	return -1, false

}

func GetCrossArenaGroup_ByAlias(tp string) (CrossArenaGroup, bool) {
	switch tp {
	case "新手组":
		return CrossArenaGroup_Rookie, true
	case "精英组":
		return CrossArenaGroup_Elite, true
	case "荣耀组":
		return CrossArenaGroup_Honor, true
	case "王者组":
		return CrossArenaGroup_King, true

	}
	return -1, false

}
