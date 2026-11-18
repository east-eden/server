package config

// 街机厅-拳皇/ArcadeEnum
type ArcadeProgress int

const (
	ArcadeProgress_BEGIN  ArcadeProgress = 1 //枚举开始
	ArcadeProgress_Open   ArcadeProgress = 1 //开始
	ArcadeProgress_Closed ArcadeProgress = 2 //未开启
	ArcadeProgress_END    ArcadeProgress = 3 //枚举结束
)

// 街机厅-拳皇/ArcadeEnum
type ArcadeGameType int

const (
	ArcadeGameType_BEGIN   ArcadeGameType = 0 //枚举开始
	ArcadeGameType_TypeAll ArcadeGameType = 0 //全部
	ArcadeGameType_TypeA   ArcadeGameType = 1 //类型一
	ArcadeGameType_TypeB   ArcadeGameType = 2 //类型二
	ArcadeGameType_TypeC   ArcadeGameType = 3 //类型三
	ArcadeGameType_TypeD   ArcadeGameType = 4 //类型四
	ArcadeGameType_END     ArcadeGameType = 5 //枚举结束
)

func (e ArcadeProgress) String() string {

	if e == ArcadeProgress_Open {
		return "Open"
	}
	if e == ArcadeProgress_Closed {
		return "Closed"
	}

	return ""

}

func GetArcadeProgress(tp string) (ArcadeProgress, bool) {
	switch tp {
	case "Open":
		return ArcadeProgress_Open, true
	case "Closed":
		return ArcadeProgress_Closed, true

	}
	return -1, false

}

func GetArcadeProgress_ByAlias(tp string) (ArcadeProgress, bool) {
	switch tp {
	case "开始":
		return ArcadeProgress_Open, true
	case "未开启":
		return ArcadeProgress_Closed, true

	}
	return -1, false

}

func (e ArcadeGameType) String() string {

	if e == ArcadeGameType_TypeAll {
		return "TypeAll"
	}
	if e == ArcadeGameType_TypeA {
		return "TypeA"
	}
	if e == ArcadeGameType_TypeB {
		return "TypeB"
	}
	if e == ArcadeGameType_TypeC {
		return "TypeC"
	}
	if e == ArcadeGameType_TypeD {
		return "TypeD"
	}

	return ""

}

func GetArcadeGameType(tp string) (ArcadeGameType, bool) {
	switch tp {
	case "TypeAll":
		return ArcadeGameType_TypeAll, true
	case "TypeA":
		return ArcadeGameType_TypeA, true
	case "TypeB":
		return ArcadeGameType_TypeB, true
	case "TypeC":
		return ArcadeGameType_TypeC, true
	case "TypeD":
		return ArcadeGameType_TypeD, true

	}
	return -1, false

}

func GetArcadeGameType_ByAlias(tp string) (ArcadeGameType, bool) {
	switch tp {
	case "全部":
		return ArcadeGameType_TypeAll, true
	case "类型一":
		return ArcadeGameType_TypeA, true
	case "类型二":
		return ArcadeGameType_TypeB, true
	case "类型三":
		return ArcadeGameType_TypeC, true
	case "类型四":
		return ArcadeGameType_TypeD, true

	}
	return -1, false

}
