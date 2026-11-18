package config

// 交流群表-拳皇/ChatGroupEnum
type ChatGroupPlatform int

const (
	ChatGroupPlatform_BEGIN  ChatGroupPlatform = 1 //枚举开始
	ChatGroupPlatform_Weixin ChatGroupPlatform = 1 //微信
	ChatGroupPlatform_Douyin ChatGroupPlatform = 2 //抖音
	ChatGroupPlatform_Alipay ChatGroupPlatform = 3 //支付宝
	ChatGroupPlatform_END    ChatGroupPlatform = 4 //枚举结束
)

func (e ChatGroupPlatform) String() string {

	if e == ChatGroupPlatform_Weixin {
		return "Weixin"
	}
	if e == ChatGroupPlatform_Douyin {
		return "Douyin"
	}
	if e == ChatGroupPlatform_Alipay {
		return "Alipay"
	}

	return ""

}

func GetChatGroupPlatform(tp string) (ChatGroupPlatform, bool) {
	switch tp {
	case "Weixin":
		return ChatGroupPlatform_Weixin, true
	case "Douyin":
		return ChatGroupPlatform_Douyin, true
	case "Alipay":
		return ChatGroupPlatform_Alipay, true

	}
	return -1, false

}

func GetChatGroupPlatform_ByAlias(tp string) (ChatGroupPlatform, bool) {
	switch tp {
	case "微信":
		return ChatGroupPlatform_Weixin, true
	case "抖音":
		return ChatGroupPlatform_Douyin, true
	case "支付宝":
		return ChatGroupPlatform_Alipay, true

	}
	return -1, false

}
