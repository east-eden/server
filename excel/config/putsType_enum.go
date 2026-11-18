package config

// LD创意工坊-拳皇 client_skip/PutsType
type WsPutsType int

const (
	WsPutsType_BEGIN       WsPutsType = 1 //枚举开始
	WsPutsType_Grid        WsPutsType = 1 //格子
	WsPutsType_MonsterBorn WsPutsType = 2 //刷怪点
	WsPutsType_Bg          WsPutsType = 3 //背景
	WsPutsType_END         WsPutsType = 4 //枚举结束
)

// LD创意工坊-拳皇 client_skip/PutsType
type WsLoopType int

const (
	WsLoopType_BEGIN    WsLoopType = 0 //枚举开始
	WsLoopType_Not      WsLoopType = 0 //不循环
	WsLoopType_Loop     WsLoopType = 1 //循环
	WsLoopType_PingPong WsLoopType = 2 //往返
	WsLoopType_END      WsLoopType = 3 //枚举结束
)

// LD创意工坊-拳皇 client_skip/PutsType
type WsLevelStatus int

const (
	WsLevelStatus_BEGIN     WsLevelStatus = 1 //枚举开始
	WsLevelStatus_Graft     WsLevelStatus = 1 //草稿
	WsLevelStatus_Published WsLevelStatus = 2 //已发布
	WsLevelStatus_Trash     WsLevelStatus = 3 //回收站
	WsLevelStatus_END       WsLevelStatus = 4 //枚举结束
)

// LD创意工坊-拳皇 client_skip/PutsType
type WSPlayCondition int

const (
	WSPlayCondition_BEGIN WSPlayCondition = 0 //枚举开始
	WSPlayCondition_Free  WSPlayCondition = 0 //免费
	WSPlayCondition_Ad    WSPlayCondition = 1 //广告
	WSPlayCondition_Props WSPlayCondition = 2 //消耗道具
	WSPlayCondition_END   WSPlayCondition = 3 //枚举结束
)

// LD创意工坊-拳皇 client_skip/PutsType
type WsHeroCamp int

const (
	WsHeroCamp_BEGIN    WsHeroCamp = 1 //枚举开始
	WsHeroCamp_Daily    WsHeroCamp = 1 //每日阵容
	WsHeroCamp_Grow     WsHeroCamp = 2 //培养阵容
	WsHeroCamp_Specific WsHeroCamp = 3 //指定神话阵容
	WsHeroCamp_END      WsHeroCamp = 4 //枚举结束
)

// LD创意工坊-拳皇 client_skip/PutsType
type WsRandomWay int

const (
	WsRandomWay_BEGIN         WsRandomWay = 1 //枚举开始
	WsRandomWay_No            WsRandomWay = 1 //不随机
	WsRandomWay_Once          WsRandomWay = 2 //一次性
	WsRandomWay_EveryInterval WsRandomWay = 3 //每次
	WsRandomWay_END           WsRandomWay = 4 //枚举结束
)

func (e WsPutsType) String() string {

	if e == WsPutsType_Grid {
		return "Grid"
	}
	if e == WsPutsType_MonsterBorn {
		return "MonsterBorn"
	}
	if e == WsPutsType_Bg {
		return "Bg"
	}

	return ""

}

func GetWsPutsType(tp string) (WsPutsType, bool) {
	switch tp {
	case "Grid":
		return WsPutsType_Grid, true
	case "MonsterBorn":
		return WsPutsType_MonsterBorn, true
	case "Bg":
		return WsPutsType_Bg, true

	}
	return -1, false

}

func GetWsPutsType_ByAlias(tp string) (WsPutsType, bool) {
	switch tp {
	case "格子":
		return WsPutsType_Grid, true
	case "刷怪点":
		return WsPutsType_MonsterBorn, true
	case "背景":
		return WsPutsType_Bg, true

	}
	return -1, false

}

func (e WsLoopType) String() string {

	if e == WsLoopType_Not {
		return "Not"
	}
	if e == WsLoopType_Loop {
		return "Loop"
	}
	if e == WsLoopType_PingPong {
		return "PingPong"
	}

	return ""

}

func GetWsLoopType(tp string) (WsLoopType, bool) {
	switch tp {
	case "Not":
		return WsLoopType_Not, true
	case "Loop":
		return WsLoopType_Loop, true
	case "PingPong":
		return WsLoopType_PingPong, true

	}
	return -1, false

}

func GetWsLoopType_ByAlias(tp string) (WsLoopType, bool) {
	switch tp {
	case "不循环":
		return WsLoopType_Not, true
	case "循环":
		return WsLoopType_Loop, true
	case "往返":
		return WsLoopType_PingPong, true

	}
	return -1, false

}

func (e WsLevelStatus) String() string {

	if e == WsLevelStatus_Graft {
		return "Graft"
	}
	if e == WsLevelStatus_Published {
		return "Published"
	}
	if e == WsLevelStatus_Trash {
		return "Trash"
	}

	return ""

}

func GetWsLevelStatus(tp string) (WsLevelStatus, bool) {
	switch tp {
	case "Graft":
		return WsLevelStatus_Graft, true
	case "Published":
		return WsLevelStatus_Published, true
	case "Trash":
		return WsLevelStatus_Trash, true

	}
	return -1, false

}

func GetWsLevelStatus_ByAlias(tp string) (WsLevelStatus, bool) {
	switch tp {
	case "草稿":
		return WsLevelStatus_Graft, true
	case "已发布":
		return WsLevelStatus_Published, true
	case "回收站":
		return WsLevelStatus_Trash, true

	}
	return -1, false

}

func (e WSPlayCondition) String() string {

	if e == WSPlayCondition_Free {
		return "Free"
	}
	if e == WSPlayCondition_Ad {
		return "Ad"
	}
	if e == WSPlayCondition_Props {
		return "Props"
	}

	return ""

}

func GetWSPlayCondition(tp string) (WSPlayCondition, bool) {
	switch tp {
	case "Free":
		return WSPlayCondition_Free, true
	case "Ad":
		return WSPlayCondition_Ad, true
	case "Props":
		return WSPlayCondition_Props, true

	}
	return -1, false

}

func GetWSPlayCondition_ByAlias(tp string) (WSPlayCondition, bool) {
	switch tp {
	case "免费":
		return WSPlayCondition_Free, true
	case "广告":
		return WSPlayCondition_Ad, true
	case "消耗道具":
		return WSPlayCondition_Props, true

	}
	return -1, false

}

func (e WsHeroCamp) String() string {

	if e == WsHeroCamp_Daily {
		return "Daily"
	}
	if e == WsHeroCamp_Grow {
		return "Grow"
	}
	if e == WsHeroCamp_Specific {
		return "Specific"
	}

	return ""

}

func GetWsHeroCamp(tp string) (WsHeroCamp, bool) {
	switch tp {
	case "Daily":
		return WsHeroCamp_Daily, true
	case "Grow":
		return WsHeroCamp_Grow, true
	case "Specific":
		return WsHeroCamp_Specific, true

	}
	return -1, false

}

func GetWsHeroCamp_ByAlias(tp string) (WsHeroCamp, bool) {
	switch tp {
	case "每日阵容":
		return WsHeroCamp_Daily, true
	case "培养阵容":
		return WsHeroCamp_Grow, true
	case "指定神话阵容":
		return WsHeroCamp_Specific, true

	}
	return -1, false

}

func (e WsRandomWay) String() string {

	if e == WsRandomWay_No {
		return "No"
	}
	if e == WsRandomWay_Once {
		return "Once"
	}
	if e == WsRandomWay_EveryInterval {
		return "EveryInterval"
	}

	return ""

}

func GetWsRandomWay(tp string) (WsRandomWay, bool) {
	switch tp {
	case "No":
		return WsRandomWay_No, true
	case "Once":
		return WsRandomWay_Once, true
	case "EveryInterval":
		return WsRandomWay_EveryInterval, true

	}
	return -1, false

}

func GetWsRandomWay_ByAlias(tp string) (WsRandomWay, bool) {
	switch tp {
	case "不随机":
		return WsRandomWay_No, true
	case "一次性":
		return WsRandomWay_Once, true
	case "每次":
		return WsRandomWay_EveryInterval, true

	}
	return -1, false

}
