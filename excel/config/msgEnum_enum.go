package config

// 文本配置表-拳皇/MsgEnum
type SubscribeType int

const (
	SubscribeType_BEGIN             SubscribeType = 1  //枚举开始
	SubscribeType_Afk               SubscribeType = 1  //挂机奖励
	SubscribeType_Gift              SubscribeType = 2  //收藏升级
	SubscribeType_Manor             SubscribeType = 3  //领地拉车
	SubscribeType_Fish              SubscribeType = 4  //钓鱼体力
	SubscribeType_Bounty            SubscribeType = 5  //赏金战报名
	SubscribeType_CrossArena        SubscribeType = 6  //跨服32强资格赛
	SubscribeType_ManorRefresh      SubscribeType = 7  //领地刷新
	SubscribeType_ManorPlunder      SubscribeType = 8  //领地掠夺
	SubscribeType_BountyRaid        SubscribeType = 9  //赏金掠夺开启
	SubscribeType_CrossArenaCheer   SubscribeType = 10 //跨服32强公会助威
	SubscribeType_CrossArenaRank    SubscribeType = 11 //跨服32强排位赛
	SubscribeType_DailyAward        SubscribeType = 12 //每日活动奖励
	SubscribeType_BountyChallenge   SubscribeType = 13 //赏金挑战开启
	SubscribeType_CrossArenaElimate SubscribeType = 14 //跨服32强淘汰赛
	SubscribeType_END               SubscribeType = 15 //枚举结束
)

// 文本配置表-拳皇/MsgEnum
type DouyinFeedType int

const (
	DouyinFeedType_BEGIN DouyinFeedType = 1 //枚举开始
	DouyinFeedType_Afk   DouyinFeedType = 1 //赛事奖励
	DouyinFeedType_Fish  DouyinFeedType = 2 //钓鱼体力
	DouyinFeedType_Daily DouyinFeedType = 3 //每日玩法
	DouyinFeedType_END   DouyinFeedType = 4 //枚举结束
)

// 文本配置表-拳皇/MsgEnum
type DouyinNewBeeFeedType int

const (
	DouyinNewBeeFeedType_BEGIN       DouyinNewBeeFeedType = 1 //枚举开始
	DouyinNewBeeFeedType_Newbee      DouyinNewBeeFeedType = 1 //新手战斗
	DouyinNewBeeFeedType_SuperClark  DouyinNewBeeFeedType = 2 //克拉克抛接
	DouyinNewBeeFeedType_OnePunch    DouyinNewBeeFeedType = 3 //一拳超人
	DouyinNewBeeFeedType_KofRoad     DouyinNewBeeFeedType = 4 //拳皇之路
	DouyinNewBeeFeedType_SimpleLogin DouyinNewBeeFeedType = 5 //登录界面
	DouyinNewBeeFeedType_END         DouyinNewBeeFeedType = 6 //枚举结束
)

// 文本配置表-拳皇/MsgEnum
type LoadingTipType int

const (
	LoadingTipType_BEGIN     LoadingTipType = 1 //枚举开始
	LoadingTipType_LDLoading LoadingTipType = 1 //LD战斗
	LoadingTipType_Normal    LoadingTipType = 2 //通用结算
	LoadingTipType_END       LoadingTipType = 3 //枚举结束
)

// 文本配置表-拳皇/MsgEnum
type InteractVideoAction int

const (
	InteractVideoAction_BEGIN     InteractVideoAction = 1  //枚举开始
	InteractVideoAction_Up        InteractVideoAction = 1  //上
	InteractVideoAction_Down      InteractVideoAction = 2  //下
	InteractVideoAction_LeftUp    InteractVideoAction = 3  //左上
	InteractVideoAction_LeftDown  InteractVideoAction = 4  //左下
	InteractVideoAction_Left      InteractVideoAction = 5  //左
	InteractVideoAction_Right     InteractVideoAction = 6  //右
	InteractVideoAction_RightUp   InteractVideoAction = 7  //右上
	InteractVideoAction_RightDown InteractVideoAction = 8  //右下
	InteractVideoAction_A         InteractVideoAction = 9  //A
	InteractVideoAction_B         InteractVideoAction = 10 //B
	InteractVideoAction_C         InteractVideoAction = 11 //C
	InteractVideoAction_D         InteractVideoAction = 12 //D
	InteractVideoAction_END       InteractVideoAction = 13 //枚举结束
)

// 文本配置表-拳皇/MsgEnum
type InterActVideoDifficulty int

const (
	InterActVideoDifficulty_BEGIN     InterActVideoDifficulty = 1 //枚举开始
	InterActVideoDifficulty_Easy      InterActVideoDifficulty = 1 //简单
	InterActVideoDifficulty_Difficult InterActVideoDifficulty = 2 //困难
	InterActVideoDifficulty_END       InterActVideoDifficulty = 3 //枚举结束
)

func (e SubscribeType) String() string {

	if e == SubscribeType_Afk {
		return "Afk"
	}
	if e == SubscribeType_Gift {
		return "Gift"
	}
	if e == SubscribeType_Manor {
		return "Manor"
	}
	if e == SubscribeType_Fish {
		return "Fish"
	}
	if e == SubscribeType_Bounty {
		return "Bounty"
	}
	if e == SubscribeType_CrossArena {
		return "CrossArena"
	}
	if e == SubscribeType_ManorRefresh {
		return "ManorRefresh"
	}
	if e == SubscribeType_ManorPlunder {
		return "ManorPlunder"
	}
	if e == SubscribeType_BountyRaid {
		return "BountyRaid"
	}
	if e == SubscribeType_CrossArenaCheer {
		return "CrossArenaCheer"
	}
	if e == SubscribeType_CrossArenaRank {
		return "CrossArenaRank"
	}
	if e == SubscribeType_DailyAward {
		return "DailyAward"
	}
	if e == SubscribeType_BountyChallenge {
		return "BountyChallenge"
	}
	if e == SubscribeType_CrossArenaElimate {
		return "CrossArenaElimate"
	}

	return ""

}

func GetSubscribeType(tp string) (SubscribeType, bool) {
	switch tp {
	case "Afk":
		return SubscribeType_Afk, true
	case "Gift":
		return SubscribeType_Gift, true
	case "Manor":
		return SubscribeType_Manor, true
	case "Fish":
		return SubscribeType_Fish, true
	case "Bounty":
		return SubscribeType_Bounty, true
	case "CrossArena":
		return SubscribeType_CrossArena, true
	case "ManorRefresh":
		return SubscribeType_ManorRefresh, true
	case "ManorPlunder":
		return SubscribeType_ManorPlunder, true
	case "BountyRaid":
		return SubscribeType_BountyRaid, true
	case "CrossArenaCheer":
		return SubscribeType_CrossArenaCheer, true
	case "CrossArenaRank":
		return SubscribeType_CrossArenaRank, true
	case "DailyAward":
		return SubscribeType_DailyAward, true
	case "BountyChallenge":
		return SubscribeType_BountyChallenge, true
	case "CrossArenaElimate":
		return SubscribeType_CrossArenaElimate, true

	}
	return -1, false

}

func GetSubscribeType_ByAlias(tp string) (SubscribeType, bool) {
	switch tp {
	case "挂机奖励":
		return SubscribeType_Afk, true
	case "收藏升级":
		return SubscribeType_Gift, true
	case "领地拉车":
		return SubscribeType_Manor, true
	case "钓鱼体力":
		return SubscribeType_Fish, true
	case "赏金战报名":
		return SubscribeType_Bounty, true
	case "跨服32强资格赛":
		return SubscribeType_CrossArena, true
	case "领地刷新":
		return SubscribeType_ManorRefresh, true
	case "领地掠夺":
		return SubscribeType_ManorPlunder, true
	case "赏金掠夺开启":
		return SubscribeType_BountyRaid, true
	case "跨服32强公会助威":
		return SubscribeType_CrossArenaCheer, true
	case "跨服32强排位赛":
		return SubscribeType_CrossArenaRank, true
	case "每日活动奖励":
		return SubscribeType_DailyAward, true
	case "赏金挑战开启":
		return SubscribeType_BountyChallenge, true
	case "跨服32强淘汰赛":
		return SubscribeType_CrossArenaElimate, true

	}
	return -1, false

}

func (e DouyinFeedType) String() string {

	if e == DouyinFeedType_Afk {
		return "Afk"
	}
	if e == DouyinFeedType_Fish {
		return "Fish"
	}
	if e == DouyinFeedType_Daily {
		return "Daily"
	}

	return ""

}

func GetDouyinFeedType(tp string) (DouyinFeedType, bool) {
	switch tp {
	case "Afk":
		return DouyinFeedType_Afk, true
	case "Fish":
		return DouyinFeedType_Fish, true
	case "Daily":
		return DouyinFeedType_Daily, true

	}
	return -1, false

}

func GetDouyinFeedType_ByAlias(tp string) (DouyinFeedType, bool) {
	switch tp {
	case "赛事奖励":
		return DouyinFeedType_Afk, true
	case "钓鱼体力":
		return DouyinFeedType_Fish, true
	case "每日玩法":
		return DouyinFeedType_Daily, true

	}
	return -1, false

}

func (e DouyinNewBeeFeedType) String() string {

	if e == DouyinNewBeeFeedType_Newbee {
		return "Newbee"
	}
	if e == DouyinNewBeeFeedType_SuperClark {
		return "SuperClark"
	}
	if e == DouyinNewBeeFeedType_OnePunch {
		return "OnePunch"
	}
	if e == DouyinNewBeeFeedType_KofRoad {
		return "KofRoad"
	}
	if e == DouyinNewBeeFeedType_SimpleLogin {
		return "SimpleLogin"
	}

	return ""

}

func GetDouyinNewBeeFeedType(tp string) (DouyinNewBeeFeedType, bool) {
	switch tp {
	case "Newbee":
		return DouyinNewBeeFeedType_Newbee, true
	case "SuperClark":
		return DouyinNewBeeFeedType_SuperClark, true
	case "OnePunch":
		return DouyinNewBeeFeedType_OnePunch, true
	case "KofRoad":
		return DouyinNewBeeFeedType_KofRoad, true
	case "SimpleLogin":
		return DouyinNewBeeFeedType_SimpleLogin, true

	}
	return -1, false

}

func GetDouyinNewBeeFeedType_ByAlias(tp string) (DouyinNewBeeFeedType, bool) {
	switch tp {
	case "新手战斗":
		return DouyinNewBeeFeedType_Newbee, true
	case "克拉克抛接":
		return DouyinNewBeeFeedType_SuperClark, true
	case "一拳超人":
		return DouyinNewBeeFeedType_OnePunch, true
	case "拳皇之路":
		return DouyinNewBeeFeedType_KofRoad, true
	case "登录界面":
		return DouyinNewBeeFeedType_SimpleLogin, true

	}
	return -1, false

}

func (e LoadingTipType) String() string {

	if e == LoadingTipType_LDLoading {
		return "LDLoading"
	}
	if e == LoadingTipType_Normal {
		return "Normal"
	}

	return ""

}

func GetLoadingTipType(tp string) (LoadingTipType, bool) {
	switch tp {
	case "LDLoading":
		return LoadingTipType_LDLoading, true
	case "Normal":
		return LoadingTipType_Normal, true

	}
	return -1, false

}

func GetLoadingTipType_ByAlias(tp string) (LoadingTipType, bool) {
	switch tp {
	case "LD战斗":
		return LoadingTipType_LDLoading, true
	case "通用结算":
		return LoadingTipType_Normal, true

	}
	return -1, false

}

func (e InteractVideoAction) String() string {

	if e == InteractVideoAction_Up {
		return "Up"
	}
	if e == InteractVideoAction_Down {
		return "Down"
	}
	if e == InteractVideoAction_LeftUp {
		return "LeftUp"
	}
	if e == InteractVideoAction_LeftDown {
		return "LeftDown"
	}
	if e == InteractVideoAction_Left {
		return "Left"
	}
	if e == InteractVideoAction_Right {
		return "Right"
	}
	if e == InteractVideoAction_RightUp {
		return "RightUp"
	}
	if e == InteractVideoAction_RightDown {
		return "RightDown"
	}
	if e == InteractVideoAction_A {
		return "A"
	}
	if e == InteractVideoAction_B {
		return "B"
	}
	if e == InteractVideoAction_C {
		return "C"
	}
	if e == InteractVideoAction_D {
		return "D"
	}

	return ""

}

func GetInteractVideoAction(tp string) (InteractVideoAction, bool) {
	switch tp {
	case "Up":
		return InteractVideoAction_Up, true
	case "Down":
		return InteractVideoAction_Down, true
	case "LeftUp":
		return InteractVideoAction_LeftUp, true
	case "LeftDown":
		return InteractVideoAction_LeftDown, true
	case "Left":
		return InteractVideoAction_Left, true
	case "Right":
		return InteractVideoAction_Right, true
	case "RightUp":
		return InteractVideoAction_RightUp, true
	case "RightDown":
		return InteractVideoAction_RightDown, true
	case "A":
		return InteractVideoAction_A, true
	case "B":
		return InteractVideoAction_B, true
	case "C":
		return InteractVideoAction_C, true
	case "D":
		return InteractVideoAction_D, true

	}
	return -1, false

}

func GetInteractVideoAction_ByAlias(tp string) (InteractVideoAction, bool) {
	switch tp {
	case "上":
		return InteractVideoAction_Up, true
	case "下":
		return InteractVideoAction_Down, true
	case "左上":
		return InteractVideoAction_LeftUp, true
	case "左下":
		return InteractVideoAction_LeftDown, true
	case "左":
		return InteractVideoAction_Left, true
	case "右":
		return InteractVideoAction_Right, true
	case "右上":
		return InteractVideoAction_RightUp, true
	case "右下":
		return InteractVideoAction_RightDown, true
	case "A":
		return InteractVideoAction_A, true
	case "B":
		return InteractVideoAction_B, true
	case "C":
		return InteractVideoAction_C, true
	case "D":
		return InteractVideoAction_D, true

	}
	return -1, false

}

func (e InterActVideoDifficulty) String() string {

	if e == InterActVideoDifficulty_Easy {
		return "Easy"
	}
	if e == InterActVideoDifficulty_Difficult {
		return "Difficult"
	}

	return ""

}

func GetInterActVideoDifficulty(tp string) (InterActVideoDifficulty, bool) {
	switch tp {
	case "Easy":
		return InterActVideoDifficulty_Easy, true
	case "Difficult":
		return InterActVideoDifficulty_Difficult, true

	}
	return -1, false

}

func GetInterActVideoDifficulty_ByAlias(tp string) (InterActVideoDifficulty, bool) {
	switch tp {
	case "简单":
		return InterActVideoDifficulty_Easy, true
	case "困难":
		return InterActVideoDifficulty_Difficult, true

	}
	return -1, false

}
