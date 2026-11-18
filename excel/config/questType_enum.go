package config

// 任务表-拳皇/QuestType
type MainTaskType int

const (
	MainTaskType_BEGIN       MainTaskType = 1 //枚举开始
	MainTaskType_PlayerLevel MainTaskType = 1 //提升等级至
	MainTaskType_END         MainTaskType = 2 //枚举结束
)

// 任务表-拳皇/QuestType
type QuestType int

const (
	QuestType_BEGIN              QuestType = 1  //枚举开始
	QuestType_MainQuest          QuestType = 1  //主线任务
	QuestType_ActivityQuest      QuestType = 2  //活动任务
	QuestType_AllianceQuest      QuestType = 3  //联盟任务
	QuestType_PopupPackageQuest  QuestType = 4  //弹出礼包任务
	QuestType_PromotionQuest     QuestType = 5  //晋升任务
	QuestType_RegularQuest       QuestType = 6  //日常周常任务
	QuestType_FishOrder          QuestType = 7  //钓鱼订单任务
	QuestType_FishSceneUnlock    QuestType = 8  //钓鱼场景解锁任务
	QuestType_FishRodBreak       QuestType = 9  //钓鱼钓竿突破任务
	QuestType_FishHookBreak      QuestType = 10 //钓鱼鱼钩突破任务
	QuestType_FishWireBreak      QuestType = 11 //钓鱼鱼线突破任务
	QuestType_FishFloatBreak     QuestType = 12 //钓鱼鱼漂突破任务
	QuestType_RoomEvent          QuestType = 13 //房间事件
	QuestType_Avatar             QuestType = 14 //头像解锁
	QuestType_AvatarFrame        QuestType = 15 //头像框解锁
	QuestType_Title              QuestType = 16 //称号解锁
	QuestType_ChatBubble         QuestType = 17 //聊天气泡解锁
	QuestType_ShopGoods          QuestType = 18 //商店商品任务
	QuestType_ActivityPointQuest QuestType = 19 //活动积分任务
	QuestType_CityWarQuest       QuestType = 20 //城市争夺成就任务
	QuestType_FameForgeUnlock    QuestType = 21 //名人堂锻造解锁
	QuestType_DestinyQuest       QuestType = 22 //宿命任务
	QuestType_DestinyFetterQuest QuestType = 23 //宿命羁绊任务
	QuestType_Leaderboard        QuestType = 24 //排行聚合
	QuestType_Counter            QuestType = 25 //计数任务索引
	QuestType_UnlimitedDraw      QuestType = 26 //随心抽任务
	QuestType_TeamLeague         QuestType = 27 //公会团赛
	QuestType_Team               QuestType = 28 //组队
	QuestType_TeamChapter        QuestType = 29 //组队章节
	QuestType_Trade              QuestType = 30 //贸易
	QuestType_KingWar            QuestType = 31 //皇城争锋
	QuestType_END                QuestType = 32 //枚举结束
)

// 任务表-拳皇/QuestType
type QuestState int

const (
	QuestState_BEGIN     QuestState = 0 //枚举开始
	QuestState_Open      QuestState = 0 //开启
	QuestState_Completed QuestState = 1 //完成
	QuestState_Rewarded  QuestState = 2 //已领奖
	QuestState_END       QuestState = 3 //枚举结束
)

func (e MainTaskType) String() string {

	if e == MainTaskType_PlayerLevel {
		return "PlayerLevel"
	}

	return ""

}

func GetMainTaskType(tp string) (MainTaskType, bool) {
	switch tp {
	case "PlayerLevel":
		return MainTaskType_PlayerLevel, true

	}
	return -1, false

}

func GetMainTaskType_ByAlias(tp string) (MainTaskType, bool) {
	switch tp {
	case "提升等级至":
		return MainTaskType_PlayerLevel, true

	}
	return -1, false

}

func (e QuestType) String() string {

	if e == QuestType_MainQuest {
		return "MainQuest"
	}
	if e == QuestType_ActivityQuest {
		return "ActivityQuest"
	}
	if e == QuestType_AllianceQuest {
		return "AllianceQuest"
	}
	if e == QuestType_PopupPackageQuest {
		return "PopupPackageQuest"
	}
	if e == QuestType_PromotionQuest {
		return "PromotionQuest"
	}
	if e == QuestType_RegularQuest {
		return "RegularQuest"
	}
	if e == QuestType_FishOrder {
		return "FishOrder"
	}
	if e == QuestType_FishSceneUnlock {
		return "FishSceneUnlock"
	}
	if e == QuestType_FishRodBreak {
		return "FishRodBreak"
	}
	if e == QuestType_FishHookBreak {
		return "FishHookBreak"
	}
	if e == QuestType_FishWireBreak {
		return "FishWireBreak"
	}
	if e == QuestType_FishFloatBreak {
		return "FishFloatBreak"
	}
	if e == QuestType_RoomEvent {
		return "RoomEvent"
	}
	if e == QuestType_Avatar {
		return "Avatar"
	}
	if e == QuestType_AvatarFrame {
		return "AvatarFrame"
	}
	if e == QuestType_Title {
		return "Title"
	}
	if e == QuestType_ChatBubble {
		return "ChatBubble"
	}
	if e == QuestType_ShopGoods {
		return "ShopGoods"
	}
	if e == QuestType_ActivityPointQuest {
		return "ActivityPointQuest"
	}
	if e == QuestType_CityWarQuest {
		return "CityWarQuest"
	}
	if e == QuestType_FameForgeUnlock {
		return "FameForgeUnlock"
	}
	if e == QuestType_DestinyQuest {
		return "DestinyQuest"
	}
	if e == QuestType_DestinyFetterQuest {
		return "DestinyFetterQuest"
	}
	if e == QuestType_Leaderboard {
		return "Leaderboard"
	}
	if e == QuestType_Counter {
		return "Counter"
	}
	if e == QuestType_UnlimitedDraw {
		return "UnlimitedDraw"
	}
	if e == QuestType_TeamLeague {
		return "TeamLeague"
	}
	if e == QuestType_Team {
		return "Team"
	}
	if e == QuestType_TeamChapter {
		return "TeamChapter"
	}
	if e == QuestType_Trade {
		return "Trade"
	}
	if e == QuestType_KingWar {
		return "KingWar"
	}

	return ""

}

func GetQuestType(tp string) (QuestType, bool) {
	switch tp {
	case "MainQuest":
		return QuestType_MainQuest, true
	case "ActivityQuest":
		return QuestType_ActivityQuest, true
	case "AllianceQuest":
		return QuestType_AllianceQuest, true
	case "PopupPackageQuest":
		return QuestType_PopupPackageQuest, true
	case "PromotionQuest":
		return QuestType_PromotionQuest, true
	case "RegularQuest":
		return QuestType_RegularQuest, true
	case "FishOrder":
		return QuestType_FishOrder, true
	case "FishSceneUnlock":
		return QuestType_FishSceneUnlock, true
	case "FishRodBreak":
		return QuestType_FishRodBreak, true
	case "FishHookBreak":
		return QuestType_FishHookBreak, true
	case "FishWireBreak":
		return QuestType_FishWireBreak, true
	case "FishFloatBreak":
		return QuestType_FishFloatBreak, true
	case "RoomEvent":
		return QuestType_RoomEvent, true
	case "Avatar":
		return QuestType_Avatar, true
	case "AvatarFrame":
		return QuestType_AvatarFrame, true
	case "Title":
		return QuestType_Title, true
	case "ChatBubble":
		return QuestType_ChatBubble, true
	case "ShopGoods":
		return QuestType_ShopGoods, true
	case "ActivityPointQuest":
		return QuestType_ActivityPointQuest, true
	case "CityWarQuest":
		return QuestType_CityWarQuest, true
	case "FameForgeUnlock":
		return QuestType_FameForgeUnlock, true
	case "DestinyQuest":
		return QuestType_DestinyQuest, true
	case "DestinyFetterQuest":
		return QuestType_DestinyFetterQuest, true
	case "Leaderboard":
		return QuestType_Leaderboard, true
	case "Counter":
		return QuestType_Counter, true
	case "UnlimitedDraw":
		return QuestType_UnlimitedDraw, true
	case "TeamLeague":
		return QuestType_TeamLeague, true
	case "Team":
		return QuestType_Team, true
	case "TeamChapter":
		return QuestType_TeamChapter, true
	case "Trade":
		return QuestType_Trade, true
	case "KingWar":
		return QuestType_KingWar, true

	}
	return -1, false

}

func GetQuestType_ByAlias(tp string) (QuestType, bool) {
	switch tp {
	case "主线任务":
		return QuestType_MainQuest, true
	case "活动任务":
		return QuestType_ActivityQuest, true
	case "联盟任务":
		return QuestType_AllianceQuest, true
	case "弹出礼包任务":
		return QuestType_PopupPackageQuest, true
	case "晋升任务":
		return QuestType_PromotionQuest, true
	case "日常周常任务":
		return QuestType_RegularQuest, true
	case "钓鱼订单任务":
		return QuestType_FishOrder, true
	case "钓鱼场景解锁任务":
		return QuestType_FishSceneUnlock, true
	case "钓鱼钓竿突破任务":
		return QuestType_FishRodBreak, true
	case "钓鱼鱼钩突破任务":
		return QuestType_FishHookBreak, true
	case "钓鱼鱼线突破任务":
		return QuestType_FishWireBreak, true
	case "钓鱼鱼漂突破任务":
		return QuestType_FishFloatBreak, true
	case "房间事件":
		return QuestType_RoomEvent, true
	case "头像解锁":
		return QuestType_Avatar, true
	case "头像框解锁":
		return QuestType_AvatarFrame, true
	case "称号解锁":
		return QuestType_Title, true
	case "聊天气泡解锁":
		return QuestType_ChatBubble, true
	case "商店商品任务":
		return QuestType_ShopGoods, true
	case "活动积分任务":
		return QuestType_ActivityPointQuest, true
	case "城市争夺成就任务":
		return QuestType_CityWarQuest, true
	case "名人堂锻造解锁":
		return QuestType_FameForgeUnlock, true
	case "宿命任务":
		return QuestType_DestinyQuest, true
	case "宿命羁绊任务":
		return QuestType_DestinyFetterQuest, true
	case "排行聚合":
		return QuestType_Leaderboard, true
	case "计数任务索引":
		return QuestType_Counter, true
	case "随心抽任务":
		return QuestType_UnlimitedDraw, true
	case "公会团赛":
		return QuestType_TeamLeague, true
	case "组队":
		return QuestType_Team, true
	case "组队章节":
		return QuestType_TeamChapter, true
	case "贸易":
		return QuestType_Trade, true
	case "皇城争锋":
		return QuestType_KingWar, true

	}
	return -1, false

}

func (e QuestState) String() string {

	if e == QuestState_Open {
		return "Open"
	}
	if e == QuestState_Completed {
		return "Completed"
	}
	if e == QuestState_Rewarded {
		return "Rewarded"
	}

	return ""

}

func GetQuestState(tp string) (QuestState, bool) {
	switch tp {
	case "Open":
		return QuestState_Open, true
	case "Completed":
		return QuestState_Completed, true
	case "Rewarded":
		return QuestState_Rewarded, true

	}
	return -1, false

}

func GetQuestState_ByAlias(tp string) (QuestState, bool) {
	switch tp {
	case "开启":
		return QuestState_Open, true
	case "完成":
		return QuestState_Completed, true
	case "已领奖":
		return QuestState_Rewarded, true

	}
	return -1, false

}
