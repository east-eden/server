package config

// 剧情对话表-拳皇/talkEnum
type TalkNPCPos int

const (
	TalkNPCPos_BEGIN  TalkNPCPos = 1 //枚举开始
	TalkNPCPos_Player TalkNPCPos = 1 //主角
	TalkNPCPos_NPC1   TalkNPCPos = 2 //NPC1
	TalkNPCPos_NPC2   TalkNPCPos = 3 //NPC2
	TalkNPCPos_NPC3   TalkNPCPos = 4 //NPC3
	TalkNPCPos_NPC4   TalkNPCPos = 5 //NPC4
	TalkNPCPos_NPC5   TalkNPCPos = 6 //NPC5
	TalkNPCPos_END    TalkNPCPos = 7 //枚举结束
)

// 剧情对话表-拳皇/talkEnum
type LevelTalkActorType int

const (
	LevelTalkActorType_BEGIN   LevelTalkActorType = 1 //枚举开始
	LevelTalkActorType_Player  LevelTalkActorType = 1 //英雄
	LevelTalkActorType_Monster LevelTalkActorType = 2 //怪物
	LevelTalkActorType_Boss    LevelTalkActorType = 3 //Boss
	LevelTalkActorType_END     LevelTalkActorType = 4 //枚举结束
)

// 剧情对话表-拳皇/talkEnum
type LevelLoadingTipType int

const (
	LevelLoadingTipType_BEGIN LevelLoadingTipType = 1 //枚举开始
	LevelLoadingTipType_Pve   LevelLoadingTipType = 1 //普通战斗
	LevelLoadingTipType_END   LevelLoadingTipType = 2 //枚举结束
)

func (e TalkNPCPos) String() string {

	if e == TalkNPCPos_Player {
		return "Player"
	}
	if e == TalkNPCPos_NPC1 {
		return "NPC1"
	}
	if e == TalkNPCPos_NPC2 {
		return "NPC2"
	}
	if e == TalkNPCPos_NPC3 {
		return "NPC3"
	}
	if e == TalkNPCPos_NPC4 {
		return "NPC4"
	}
	if e == TalkNPCPos_NPC5 {
		return "NPC5"
	}

	return ""

}

func GetTalkNPCPos(tp string) (TalkNPCPos, bool) {
	switch tp {
	case "Player":
		return TalkNPCPos_Player, true
	case "NPC1":
		return TalkNPCPos_NPC1, true
	case "NPC2":
		return TalkNPCPos_NPC2, true
	case "NPC3":
		return TalkNPCPos_NPC3, true
	case "NPC4":
		return TalkNPCPos_NPC4, true
	case "NPC5":
		return TalkNPCPos_NPC5, true

	}
	return -1, false

}

func GetTalkNPCPos_ByAlias(tp string) (TalkNPCPos, bool) {
	switch tp {
	case "主角":
		return TalkNPCPos_Player, true
	case "NPC1":
		return TalkNPCPos_NPC1, true
	case "NPC2":
		return TalkNPCPos_NPC2, true
	case "NPC3":
		return TalkNPCPos_NPC3, true
	case "NPC4":
		return TalkNPCPos_NPC4, true
	case "NPC5":
		return TalkNPCPos_NPC5, true

	}
	return -1, false

}

func (e LevelTalkActorType) String() string {

	if e == LevelTalkActorType_Player {
		return "Player"
	}
	if e == LevelTalkActorType_Monster {
		return "Monster"
	}
	if e == LevelTalkActorType_Boss {
		return "Boss"
	}

	return ""

}

func GetLevelTalkActorType(tp string) (LevelTalkActorType, bool) {
	switch tp {
	case "Player":
		return LevelTalkActorType_Player, true
	case "Monster":
		return LevelTalkActorType_Monster, true
	case "Boss":
		return LevelTalkActorType_Boss, true

	}
	return -1, false

}

func GetLevelTalkActorType_ByAlias(tp string) (LevelTalkActorType, bool) {
	switch tp {
	case "英雄":
		return LevelTalkActorType_Player, true
	case "怪物":
		return LevelTalkActorType_Monster, true
	case "Boss":
		return LevelTalkActorType_Boss, true

	}
	return -1, false

}

func (e LevelLoadingTipType) String() string {

	if e == LevelLoadingTipType_Pve {
		return "Pve"
	}

	return ""

}

func GetLevelLoadingTipType(tp string) (LevelLoadingTipType, bool) {
	switch tp {
	case "Pve":
		return LevelLoadingTipType_Pve, true

	}
	return -1, false

}

func GetLevelLoadingTipType_ByAlias(tp string) (LevelLoadingTipType, bool) {
	switch tp {
	case "普通战斗":
		return LevelLoadingTipType_Pve, true

	}
	return -1, false

}
