package config

// 广告分享表-拳皇/iaaEnum
type IAAAction int

const (
	IAAAction_BEGIN          IAAAction = 0  //枚举开始
	IAAAction_None           IAAAction = 0  //无
	IAAAction_Share          IAAAction = 1  //分享
	IAAAction_ShareVideo     IAAAction = 2  //录屏
	IAAAction_Video          IAAAction = 3  //视频
	IAAAction_Banner         IAAAction = 4  //Banner
	IAAAction_ChaPing        IAAAction = 5  //插屏
	IAAAction_ShareOnNoVideo IAAAction = 90 //视频保底分享
	IAAAction_END            IAAAction = 91 //枚举结束
)

// 广告分享表-拳皇/iaaEnum
type IAACompatible int

const (
	IAACompatible_BEGIN     IAACompatible = 0 //枚举开始
	IAACompatible_All       IAACompatible = 0 //不限制
	IAACompatible_OnlyVideo IAACompatible = 1 //仅视频
	IAACompatible_OnlyShare IAACompatible = 2 //仅分享
	IAACompatible_END       IAACompatible = 3 //枚举结束
)

// 广告分享表-拳皇/iaaEnum
type AdType int

const (
	AdType_BEGIN        AdType = 1 //枚举开始
	AdType_AllianceTech AdType = 1 //公会秘术
	AdType_Astrology    AdType = 2 //祈福占星
	AdType_END          AdType = 3 //枚举结束
)

func (e IAAAction) String() string {

	if e == IAAAction_None {
		return "None"
	}
	if e == IAAAction_Share {
		return "Share"
	}
	if e == IAAAction_ShareVideo {
		return "ShareVideo"
	}
	if e == IAAAction_Video {
		return "Video"
	}
	if e == IAAAction_Banner {
		return "Banner"
	}
	if e == IAAAction_ChaPing {
		return "ChaPing"
	}
	if e == IAAAction_ShareOnNoVideo {
		return "ShareOnNoVideo"
	}

	return ""

}

func GetIAAAction(tp string) (IAAAction, bool) {
	switch tp {
	case "None":
		return IAAAction_None, true
	case "Share":
		return IAAAction_Share, true
	case "ShareVideo":
		return IAAAction_ShareVideo, true
	case "Video":
		return IAAAction_Video, true
	case "Banner":
		return IAAAction_Banner, true
	case "ChaPing":
		return IAAAction_ChaPing, true
	case "ShareOnNoVideo":
		return IAAAction_ShareOnNoVideo, true

	}
	return -1, false

}

func GetIAAAction_ByAlias(tp string) (IAAAction, bool) {
	switch tp {
	case "无":
		return IAAAction_None, true
	case "分享":
		return IAAAction_Share, true
	case "录屏":
		return IAAAction_ShareVideo, true
	case "视频":
		return IAAAction_Video, true
	case "Banner":
		return IAAAction_Banner, true
	case "插屏":
		return IAAAction_ChaPing, true
	case "视频保底分享":
		return IAAAction_ShareOnNoVideo, true

	}
	return -1, false

}

func (e IAACompatible) String() string {

	if e == IAACompatible_All {
		return "All"
	}
	if e == IAACompatible_OnlyVideo {
		return "OnlyVideo"
	}
	if e == IAACompatible_OnlyShare {
		return "OnlyShare"
	}

	return ""

}

func GetIAACompatible(tp string) (IAACompatible, bool) {
	switch tp {
	case "All":
		return IAACompatible_All, true
	case "OnlyVideo":
		return IAACompatible_OnlyVideo, true
	case "OnlyShare":
		return IAACompatible_OnlyShare, true

	}
	return -1, false

}

func GetIAACompatible_ByAlias(tp string) (IAACompatible, bool) {
	switch tp {
	case "不限制":
		return IAACompatible_All, true
	case "仅视频":
		return IAACompatible_OnlyVideo, true
	case "仅分享":
		return IAACompatible_OnlyShare, true

	}
	return -1, false

}

func (e AdType) String() string {

	if e == AdType_AllianceTech {
		return "AllianceTech"
	}
	if e == AdType_Astrology {
		return "Astrology"
	}

	return ""

}

func GetAdType(tp string) (AdType, bool) {
	switch tp {
	case "AllianceTech":
		return AdType_AllianceTech, true
	case "Astrology":
		return AdType_Astrology, true

	}
	return -1, false

}

func GetAdType_ByAlias(tp string) (AdType, bool) {
	switch tp {
	case "公会秘术":
		return AdType_AllianceTech, true
	case "祈福占星":
		return AdType_Astrology, true

	}
	return -1, false

}
