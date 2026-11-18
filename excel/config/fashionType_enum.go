package config

// 时装表-拳皇/fashionType
type FashionType int

const (
	FashionType_BEGIN              FashionType = 1  //枚举开始
	FashionType_Avatar             FashionType = 1  //头像
	FashionType_AvatarFrame        FashionType = 2  //头像框
	FashionType_ChatBubble         FashionType = 3  //聊天气泡
	FashionType_Cutscene           FashionType = 4  //战斗场景
	FashionType_DamageEffects      FashionType = 5  //伤害特效
	FashionType_DynamicBackground  FashionType = 6  //动态背景
	FashionType_Wing               FashionType = 7  //背饰
	FashionType_KillSpecialEffects FashionType = 8  //必杀特效
	FashionType_FinalEffects       FashionType = 9  //终结特效
	FashionType_EntryAction        FashionType = 10 //入场动作
	FashionType_StandbyAction      FashionType = 11 //待机动作
	FashionType_VerticalPainting   FashionType = 12 //立绘
	FashionType_StandinPuppet      FashionType = 13 //替身人偶
	FashionType_FashionableDress   FashionType = 14 //时装
	FashionType_ManorSkin          FashionType = 15 //领地皮肤
	FashionType_HomeSkin           FashionType = 16 //主页皮肤
	FashionType_AoYi               FashionType = 17 //奥义技
	FashionType_AidSkill           FashionType = 18 //援助技
	FashionType_HeroPart1          FashionType = 19 //变身I
	FashionType_HeroPart2          FashionType = 20 //变身II
	FashionType_HeroPart3          FashionType = 21 //变身III
	FashionType_HeroPart4          FashionType = 22 //变身IV
	FashionType_HeroAction         FashionType = 24 //动作
	FashionType_END                FashionType = 25 //枚举结束
)

// 时装表-拳皇/fashionType
type FashionClothesType int

const (
	FashionClothesType_BEGIN    FashionClothesType = 1 //枚举开始
	FashionClothesType_Subject  FashionClothesType = 1 //主体
	FashionClothesType_Material FashionClothesType = 2 //素材
	FashionClothesType_END      FashionClothesType = 3 //枚举结束
)

// 时装表-拳皇/fashionType
type FashionEffectType int

const (
	FashionEffectType_BEGIN            FashionEffectType = 1 //枚举开始
	FashionEffectType_Before           FashionEffectType = 1 //强化前
	FashionEffectType_Success          FashionEffectType = 2 //成功
	FashionEffectType_Failed           FashionEffectType = 3 //失败
	FashionEffectType_ExtraSuccessRate FashionEffectType = 4 //额外成功率
	FashionEffectType_SpecialWeight    FashionEffectType = 5 //特性权重
	FashionEffectType_DownLv           FashionEffectType = 6 //降低等级
	FashionEffectType_DownLvBefore     FashionEffectType = 7 //降低等级前
	FashionEffectType_END              FashionEffectType = 8 //枚举结束
)

// 时装表-拳皇/fashionType
type FashionViewPart int

const (
	FashionViewPart_BEGIN    FashionViewPart = 1 //枚举开始
	FashionViewPart_Head     FashionViewPart = 1 //头
	FashionViewPart_Arm      FashionViewPart = 2 //臂
	FashionViewPart_UpBody   FashionViewPart = 3 //上身
	FashionViewPart_DownBody FashionViewPart = 4 //下身
	FashionViewPart_END      FashionViewPart = 5 //枚举结束
)

// 时装表-拳皇/fashionType
type AoYiType int

const (
	AoYiType_BEGIN AoYiType = 1 //枚举开始
	AoYiType_A     AoYiType = 1 //奥义1
	AoYiType_B     AoYiType = 2 //奥义2
	AoYiType_C     AoYiType = 3 //奥义3
	AoYiType_D     AoYiType = 4 //奥义4
	AoYiType_E     AoYiType = 5 //奥义5
	AoYiType_F     AoYiType = 6 //奥义6
	AoYiType_END   AoYiType = 7 //枚举结束
)

// 时装表-拳皇/fashionType
type PaintingOffsetType int

const (
	PaintingOffsetType_BEGIN          PaintingOffsetType = 1 //枚举开始
	PaintingOffsetType_BattlePos      PaintingOffsetType = 1 //战斗偏移
	PaintingOffsetType_SupportPos     PaintingOffsetType = 2 //援护偏移
	PaintingOffsetType_HeroMainPosArr PaintingOffsetType = 3 //英雄主界面偏移
	PaintingOffsetType_ResultPosArr   PaintingOffsetType = 4 //战斗结果偏移
	PaintingOffsetType_FamePosArr     PaintingOffsetType = 5 //名人堂偏移
	PaintingOffsetType_END            PaintingOffsetType = 6 //枚举结束
)

func (e FashionType) String() string {

	if e == FashionType_Avatar {
		return "Avatar"
	}
	if e == FashionType_AvatarFrame {
		return "AvatarFrame"
	}
	if e == FashionType_ChatBubble {
		return "ChatBubble"
	}
	if e == FashionType_Cutscene {
		return "Cutscene"
	}
	if e == FashionType_DamageEffects {
		return "DamageEffects"
	}
	if e == FashionType_DynamicBackground {
		return "DynamicBackground"
	}
	if e == FashionType_Wing {
		return "Wing"
	}
	if e == FashionType_KillSpecialEffects {
		return "KillSpecialEffects"
	}
	if e == FashionType_FinalEffects {
		return "FinalEffects"
	}
	if e == FashionType_EntryAction {
		return "EntryAction"
	}
	if e == FashionType_StandbyAction {
		return "StandbyAction"
	}
	if e == FashionType_VerticalPainting {
		return "VerticalPainting"
	}
	if e == FashionType_StandinPuppet {
		return "StandinPuppet"
	}
	if e == FashionType_FashionableDress {
		return "FashionableDress"
	}
	if e == FashionType_ManorSkin {
		return "ManorSkin"
	}
	if e == FashionType_HomeSkin {
		return "HomeSkin"
	}
	if e == FashionType_AoYi {
		return "AoYi"
	}
	if e == FashionType_AidSkill {
		return "AidSkill"
	}
	if e == FashionType_HeroPart1 {
		return "HeroPart1"
	}
	if e == FashionType_HeroPart2 {
		return "HeroPart2"
	}
	if e == FashionType_HeroPart3 {
		return "HeroPart3"
	}
	if e == FashionType_HeroPart4 {
		return "HeroPart4"
	}
	if e == FashionType_HeroAction {
		return "HeroAction"
	}

	return ""

}

func GetFashionType(tp string) (FashionType, bool) {
	switch tp {
	case "Avatar":
		return FashionType_Avatar, true
	case "AvatarFrame":
		return FashionType_AvatarFrame, true
	case "ChatBubble":
		return FashionType_ChatBubble, true
	case "Cutscene":
		return FashionType_Cutscene, true
	case "DamageEffects":
		return FashionType_DamageEffects, true
	case "DynamicBackground":
		return FashionType_DynamicBackground, true
	case "Wing":
		return FashionType_Wing, true
	case "KillSpecialEffects":
		return FashionType_KillSpecialEffects, true
	case "FinalEffects":
		return FashionType_FinalEffects, true
	case "EntryAction":
		return FashionType_EntryAction, true
	case "StandbyAction":
		return FashionType_StandbyAction, true
	case "VerticalPainting":
		return FashionType_VerticalPainting, true
	case "StandinPuppet":
		return FashionType_StandinPuppet, true
	case "FashionableDress":
		return FashionType_FashionableDress, true
	case "ManorSkin":
		return FashionType_ManorSkin, true
	case "HomeSkin":
		return FashionType_HomeSkin, true
	case "AoYi":
		return FashionType_AoYi, true
	case "AidSkill":
		return FashionType_AidSkill, true
	case "HeroPart1":
		return FashionType_HeroPart1, true
	case "HeroPart2":
		return FashionType_HeroPart2, true
	case "HeroPart3":
		return FashionType_HeroPart3, true
	case "HeroPart4":
		return FashionType_HeroPart4, true
	case "HeroAction":
		return FashionType_HeroAction, true

	}
	return -1, false

}

func GetFashionType_ByAlias(tp string) (FashionType, bool) {
	switch tp {
	case "头像":
		return FashionType_Avatar, true
	case "头像框":
		return FashionType_AvatarFrame, true
	case "聊天气泡":
		return FashionType_ChatBubble, true
	case "战斗场景":
		return FashionType_Cutscene, true
	case "伤害特效":
		return FashionType_DamageEffects, true
	case "动态背景":
		return FashionType_DynamicBackground, true
	case "背饰":
		return FashionType_Wing, true
	case "必杀特效":
		return FashionType_KillSpecialEffects, true
	case "终结特效":
		return FashionType_FinalEffects, true
	case "入场动作":
		return FashionType_EntryAction, true
	case "待机动作":
		return FashionType_StandbyAction, true
	case "立绘":
		return FashionType_VerticalPainting, true
	case "替身人偶":
		return FashionType_StandinPuppet, true
	case "时装":
		return FashionType_FashionableDress, true
	case "领地皮肤":
		return FashionType_ManorSkin, true
	case "主页皮肤":
		return FashionType_HomeSkin, true
	case "奥义技":
		return FashionType_AoYi, true
	case "援助技":
		return FashionType_AidSkill, true
	case "变身I":
		return FashionType_HeroPart1, true
	case "变身II":
		return FashionType_HeroPart2, true
	case "变身III":
		return FashionType_HeroPart3, true
	case "变身IV":
		return FashionType_HeroPart4, true
	case "动作":
		return FashionType_HeroAction, true

	}
	return -1, false

}

func (e FashionClothesType) String() string {

	if e == FashionClothesType_Subject {
		return "Subject"
	}
	if e == FashionClothesType_Material {
		return "Material"
	}

	return ""

}

func GetFashionClothesType(tp string) (FashionClothesType, bool) {
	switch tp {
	case "Subject":
		return FashionClothesType_Subject, true
	case "Material":
		return FashionClothesType_Material, true

	}
	return -1, false

}

func GetFashionClothesType_ByAlias(tp string) (FashionClothesType, bool) {
	switch tp {
	case "主体":
		return FashionClothesType_Subject, true
	case "素材":
		return FashionClothesType_Material, true

	}
	return -1, false

}

func (e FashionEffectType) String() string {

	if e == FashionEffectType_Before {
		return "Before"
	}
	if e == FashionEffectType_Success {
		return "Success"
	}
	if e == FashionEffectType_Failed {
		return "Failed"
	}
	if e == FashionEffectType_ExtraSuccessRate {
		return "ExtraSuccessRate"
	}
	if e == FashionEffectType_SpecialWeight {
		return "SpecialWeight"
	}
	if e == FashionEffectType_DownLv {
		return "DownLv"
	}
	if e == FashionEffectType_DownLvBefore {
		return "DownLvBefore"
	}

	return ""

}

func GetFashionEffectType(tp string) (FashionEffectType, bool) {
	switch tp {
	case "Before":
		return FashionEffectType_Before, true
	case "Success":
		return FashionEffectType_Success, true
	case "Failed":
		return FashionEffectType_Failed, true
	case "ExtraSuccessRate":
		return FashionEffectType_ExtraSuccessRate, true
	case "SpecialWeight":
		return FashionEffectType_SpecialWeight, true
	case "DownLv":
		return FashionEffectType_DownLv, true
	case "DownLvBefore":
		return FashionEffectType_DownLvBefore, true

	}
	return -1, false

}

func GetFashionEffectType_ByAlias(tp string) (FashionEffectType, bool) {
	switch tp {
	case "强化前":
		return FashionEffectType_Before, true
	case "成功":
		return FashionEffectType_Success, true
	case "失败":
		return FashionEffectType_Failed, true
	case "额外成功率":
		return FashionEffectType_ExtraSuccessRate, true
	case "特性权重":
		return FashionEffectType_SpecialWeight, true
	case "降低等级":
		return FashionEffectType_DownLv, true
	case "降低等级前":
		return FashionEffectType_DownLvBefore, true

	}
	return -1, false

}

func (e FashionViewPart) String() string {

	if e == FashionViewPart_Head {
		return "Head"
	}
	if e == FashionViewPart_Arm {
		return "Arm"
	}
	if e == FashionViewPart_UpBody {
		return "UpBody"
	}
	if e == FashionViewPart_DownBody {
		return "DownBody"
	}

	return ""

}

func GetFashionViewPart(tp string) (FashionViewPart, bool) {
	switch tp {
	case "Head":
		return FashionViewPart_Head, true
	case "Arm":
		return FashionViewPart_Arm, true
	case "UpBody":
		return FashionViewPart_UpBody, true
	case "DownBody":
		return FashionViewPart_DownBody, true

	}
	return -1, false

}

func GetFashionViewPart_ByAlias(tp string) (FashionViewPart, bool) {
	switch tp {
	case "头":
		return FashionViewPart_Head, true
	case "臂":
		return FashionViewPart_Arm, true
	case "上身":
		return FashionViewPart_UpBody, true
	case "下身":
		return FashionViewPart_DownBody, true

	}
	return -1, false

}

func (e AoYiType) String() string {

	if e == AoYiType_A {
		return "A"
	}
	if e == AoYiType_B {
		return "B"
	}
	if e == AoYiType_C {
		return "C"
	}
	if e == AoYiType_D {
		return "D"
	}
	if e == AoYiType_E {
		return "E"
	}
	if e == AoYiType_F {
		return "F"
	}

	return ""

}

func GetAoYiType(tp string) (AoYiType, bool) {
	switch tp {
	case "A":
		return AoYiType_A, true
	case "B":
		return AoYiType_B, true
	case "C":
		return AoYiType_C, true
	case "D":
		return AoYiType_D, true
	case "E":
		return AoYiType_E, true
	case "F":
		return AoYiType_F, true

	}
	return -1, false

}

func GetAoYiType_ByAlias(tp string) (AoYiType, bool) {
	switch tp {
	case "奥义1":
		return AoYiType_A, true
	case "奥义2":
		return AoYiType_B, true
	case "奥义3":
		return AoYiType_C, true
	case "奥义4":
		return AoYiType_D, true
	case "奥义5":
		return AoYiType_E, true
	case "奥义6":
		return AoYiType_F, true

	}
	return -1, false

}

func (e PaintingOffsetType) String() string {

	if e == PaintingOffsetType_BattlePos {
		return "BattlePos"
	}
	if e == PaintingOffsetType_SupportPos {
		return "SupportPos"
	}
	if e == PaintingOffsetType_HeroMainPosArr {
		return "HeroMainPosArr"
	}
	if e == PaintingOffsetType_ResultPosArr {
		return "ResultPosArr"
	}
	if e == PaintingOffsetType_FamePosArr {
		return "FamePosArr"
	}

	return ""

}

func GetPaintingOffsetType(tp string) (PaintingOffsetType, bool) {
	switch tp {
	case "BattlePos":
		return PaintingOffsetType_BattlePos, true
	case "SupportPos":
		return PaintingOffsetType_SupportPos, true
	case "HeroMainPosArr":
		return PaintingOffsetType_HeroMainPosArr, true
	case "ResultPosArr":
		return PaintingOffsetType_ResultPosArr, true
	case "FamePosArr":
		return PaintingOffsetType_FamePosArr, true

	}
	return -1, false

}

func GetPaintingOffsetType_ByAlias(tp string) (PaintingOffsetType, bool) {
	switch tp {
	case "战斗偏移":
		return PaintingOffsetType_BattlePos, true
	case "援护偏移":
		return PaintingOffsetType_SupportPos, true
	case "英雄主界面偏移":
		return PaintingOffsetType_HeroMainPosArr, true
	case "战斗结果偏移":
		return PaintingOffsetType_ResultPosArr, true
	case "名人堂偏移":
		return PaintingOffsetType_FamePosArr, true

	}
	return -1, false

}
