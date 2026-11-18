package config

// 基础配置-拳皇/reportType
type ReportType int

const (
	ReportType_BEGIN       ReportType = 1 //枚举开始
	ReportType_LimitWords  ReportType = 1 //不当字眼
	ReportType_Plug        ReportType = 2 //使用恶意软件
	ReportType_Uncooperate ReportType = 3 //消极游戏
	ReportType_END         ReportType = 4 //枚举结束
)

func (e ReportType) String() string {

	if e == ReportType_LimitWords {
		return "LimitWords"
	}
	if e == ReportType_Plug {
		return "Plug"
	}
	if e == ReportType_Uncooperate {
		return "Uncooperate"
	}

	return ""

}

func GetReportType(tp string) (ReportType, bool) {
	switch tp {
	case "LimitWords":
		return ReportType_LimitWords, true
	case "Plug":
		return ReportType_Plug, true
	case "Uncooperate":
		return ReportType_Uncooperate, true

	}
	return -1, false

}

func GetReportType_ByAlias(tp string) (ReportType, bool) {
	switch tp {
	case "不当字眼":
		return ReportType_LimitWords, true
	case "使用恶意软件":
		return ReportType_Plug, true
	case "消极游戏":
		return ReportType_Uncooperate, true

	}
	return -1, false

}
