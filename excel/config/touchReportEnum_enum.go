package config

// 客户端触摸打点-拳皇/TouchReportEnum
type TouchReportKey int

const (
	TouchReportKey_BEGIN      TouchReportKey = 1 //枚举开始
	TouchReportKey_ManorApply TouchReportKey = 1 //领地上车
	TouchReportKey_END        TouchReportKey = 2 //枚举结束
)

func (e TouchReportKey) String() string {

	if e == TouchReportKey_ManorApply {
		return "ManorApply"
	}

	return ""

}

func GetTouchReportKey(tp string) (TouchReportKey, bool) {
	switch tp {
	case "ManorApply":
		return TouchReportKey_ManorApply, true

	}
	return -1, false

}

func GetTouchReportKey_ByAlias(tp string) (TouchReportKey, bool) {
	switch tp {
	case "领地上车":
		return TouchReportKey_ManorApply, true

	}
	return -1, false

}
