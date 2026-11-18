package config

// 聊天表-拳皇/ChatHelpType
type ChatHelpType int

const (
	ChatHelpType_BEGIN           ChatHelpType = 1 //枚举开始
	ChatHelpType_BoxGift         ChatHelpType = 1 //箱子礼物
	ChatHelpType_Fish            ChatHelpType = 2 //鱼缸
	ChatHelpType_RecruitLevel    ChatHelpType = 3 //招募等级
	ChatHelpType_NestsTrap       ChatHelpType = 4 //音巢陷阱
	ChatHelpType_NestsTrapDebuff ChatHelpType = 5 //音巢debuff陷阱
	ChatHelpType_END             ChatHelpType = 6 //枚举结束
)

// 聊天表-拳皇/ChatHelpType
type ChatEmojiType int

const (
	ChatEmojiType_BEGIN ChatEmojiType = 1 //枚举开始
	ChatEmojiType_Image ChatEmojiType = 1 //图片
	ChatEmojiType_Spine ChatEmojiType = 2 //Spine
	ChatEmojiType_END   ChatEmojiType = 3 //枚举结束
)

// 聊天表-拳皇/ChatHelpType
type RedEnvelopType int

const (
	RedEnvelopType_BEGIN   RedEnvelopType = 1 //枚举开始
	RedEnvelopType_Equally RedEnvelopType = 1 //平分
	RedEnvelopType_Random  RedEnvelopType = 2 //随机
	RedEnvelopType_END     RedEnvelopType = 3 //枚举结束
)

func (e ChatHelpType) String() string {

	if e == ChatHelpType_BoxGift {
		return "BoxGift"
	}
	if e == ChatHelpType_Fish {
		return "Fish"
	}
	if e == ChatHelpType_RecruitLevel {
		return "RecruitLevel"
	}
	if e == ChatHelpType_NestsTrap {
		return "NestsTrap"
	}
	if e == ChatHelpType_NestsTrapDebuff {
		return "NestsTrapDebuff"
	}

	return ""

}

func GetChatHelpType(tp string) (ChatHelpType, bool) {
	switch tp {
	case "BoxGift":
		return ChatHelpType_BoxGift, true
	case "Fish":
		return ChatHelpType_Fish, true
	case "RecruitLevel":
		return ChatHelpType_RecruitLevel, true
	case "NestsTrap":
		return ChatHelpType_NestsTrap, true
	case "NestsTrapDebuff":
		return ChatHelpType_NestsTrapDebuff, true

	}
	return -1, false

}

func GetChatHelpType_ByAlias(tp string) (ChatHelpType, bool) {
	switch tp {
	case "箱子礼物":
		return ChatHelpType_BoxGift, true
	case "鱼缸":
		return ChatHelpType_Fish, true
	case "招募等级":
		return ChatHelpType_RecruitLevel, true
	case "音巢陷阱":
		return ChatHelpType_NestsTrap, true
	case "音巢debuff陷阱":
		return ChatHelpType_NestsTrapDebuff, true

	}
	return -1, false

}

func (e ChatEmojiType) String() string {

	if e == ChatEmojiType_Image {
		return "Image"
	}
	if e == ChatEmojiType_Spine {
		return "Spine"
	}

	return ""

}

func GetChatEmojiType(tp string) (ChatEmojiType, bool) {
	switch tp {
	case "Image":
		return ChatEmojiType_Image, true
	case "Spine":
		return ChatEmojiType_Spine, true

	}
	return -1, false

}

func GetChatEmojiType_ByAlias(tp string) (ChatEmojiType, bool) {
	switch tp {
	case "图片":
		return ChatEmojiType_Image, true
	case "Spine":
		return ChatEmojiType_Spine, true

	}
	return -1, false

}

func (e RedEnvelopType) String() string {

	if e == RedEnvelopType_Equally {
		return "Equally"
	}
	if e == RedEnvelopType_Random {
		return "Random"
	}

	return ""

}

func GetRedEnvelopType(tp string) (RedEnvelopType, bool) {
	switch tp {
	case "Equally":
		return RedEnvelopType_Equally, true
	case "Random":
		return RedEnvelopType_Random, true

	}
	return -1, false

}

func GetRedEnvelopType_ByAlias(tp string) (RedEnvelopType, bool) {
	switch tp {
	case "平分":
		return RedEnvelopType_Equally, true
	case "随机":
		return RedEnvelopType_Random, true

	}
	return -1, false

}
