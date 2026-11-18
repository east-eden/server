package config

// 联盟-拳皇/AlliancePlayerPower
type AlliancePowerType int

const (
	AlliancePowerType_BEGIN               AlliancePowerType = 0  //枚举开始
	AlliancePowerType_None                AlliancePowerType = 0  //无
	AlliancePowerType_ModifyName          AlliancePowerType = 1  //修改名字
	AlliancePowerType_ModifyNotice        AlliancePowerType = 2  //修改公告
	AlliancePowerType_Setting             AlliancePowerType = 3  //设置
	AlliancePowerType_UploadQrcode        AlliancePowerType = 4  //上传二维码
	AlliancePowerType_HandleRequest       AlliancePowerType = 5  //处理公会申请消息
	AlliancePowerType_Dismiss             AlliancePowerType = 6  //解散公会
	AlliancePowerType_UpgradeStatue       AlliancePowerType = 7  //升级雕像
	AlliancePowerType_ModifyJoinCondition AlliancePowerType = 8  //修改加入条件和开启审核
	AlliancePowerType_KickMember          AlliancePowerType = 9  //移除成员
	AlliancePowerType_AppointDeputy       AlliancePowerType = 10 //任命副会长
	AlliancePowerType_AppointElite        AlliancePowerType = 11 //任命精英
	AlliancePowerType_END                 AlliancePowerType = 12 //枚举结束
)

func (e AlliancePowerType) String() string {

	if e == AlliancePowerType_None {
		return "None"
	}
	if e == AlliancePowerType_ModifyName {
		return "ModifyName"
	}
	if e == AlliancePowerType_ModifyNotice {
		return "ModifyNotice"
	}
	if e == AlliancePowerType_Setting {
		return "Setting"
	}
	if e == AlliancePowerType_UploadQrcode {
		return "UploadQrcode"
	}
	if e == AlliancePowerType_HandleRequest {
		return "HandleRequest"
	}
	if e == AlliancePowerType_Dismiss {
		return "Dismiss"
	}
	if e == AlliancePowerType_UpgradeStatue {
		return "UpgradeStatue"
	}
	if e == AlliancePowerType_ModifyJoinCondition {
		return "ModifyJoinCondition"
	}
	if e == AlliancePowerType_KickMember {
		return "KickMember"
	}
	if e == AlliancePowerType_AppointDeputy {
		return "AppointDeputy"
	}
	if e == AlliancePowerType_AppointElite {
		return "AppointElite"
	}

	return ""

}

func GetAlliancePowerType(tp string) (AlliancePowerType, bool) {
	switch tp {
	case "None":
		return AlliancePowerType_None, true
	case "ModifyName":
		return AlliancePowerType_ModifyName, true
	case "ModifyNotice":
		return AlliancePowerType_ModifyNotice, true
	case "Setting":
		return AlliancePowerType_Setting, true
	case "UploadQrcode":
		return AlliancePowerType_UploadQrcode, true
	case "HandleRequest":
		return AlliancePowerType_HandleRequest, true
	case "Dismiss":
		return AlliancePowerType_Dismiss, true
	case "UpgradeStatue":
		return AlliancePowerType_UpgradeStatue, true
	case "ModifyJoinCondition":
		return AlliancePowerType_ModifyJoinCondition, true
	case "KickMember":
		return AlliancePowerType_KickMember, true
	case "AppointDeputy":
		return AlliancePowerType_AppointDeputy, true
	case "AppointElite":
		return AlliancePowerType_AppointElite, true

	}
	return -1, false

}

func GetAlliancePowerType_ByAlias(tp string) (AlliancePowerType, bool) {
	switch tp {
	case "无":
		return AlliancePowerType_None, true
	case "修改名字":
		return AlliancePowerType_ModifyName, true
	case "修改公告":
		return AlliancePowerType_ModifyNotice, true
	case "设置":
		return AlliancePowerType_Setting, true
	case "上传二维码":
		return AlliancePowerType_UploadQrcode, true
	case "处理公会申请消息":
		return AlliancePowerType_HandleRequest, true
	case "解散公会":
		return AlliancePowerType_Dismiss, true
	case "升级雕像":
		return AlliancePowerType_UpgradeStatue, true
	case "修改加入条件和开启审核":
		return AlliancePowerType_ModifyJoinCondition, true
	case "移除成员":
		return AlliancePowerType_KickMember, true
	case "任命副会长":
		return AlliancePowerType_AppointDeputy, true
	case "任命精英":
		return AlliancePowerType_AppointElite, true

	}
	return -1, false

}
