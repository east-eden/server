package define

//------------------------------------------------------------------------------------------
// 英雄种族
//------------------------------------------------------------------------------------------
type EHeroRaceType int32

const (
	HeroRace_Begin  EHeroRaceType = iota
	HeroRace_Kindom EHeroRaceType = iota - 1 // 王国
	HeroRace_forest                          // 森林
	HeroRace_End
)

//-------------------------------------------------------------------------------
// 状态类型
//-------------------------------------------------------------------------------
type EHeroState int32

const (
	HeroState_Begin            EHeroState = iota
	HeroState_Dead             EHeroState = iota - 1 // 0 死亡
	HeroState_Solid                                  // 1 石化
	HeroState_Freeze                                 // 2 冻结
	HeroState_Stun                                   // 3 眩晕
	HeroState_Fire                                   // 4 灼烧
	HeroState_Seal                                   // 5 封印
	HeroState_UnBeat                                 // 6 无敌
	HeroState_UnDead                                 // 7 不死
	HeroState_Anger                                  // 8 风怒
	HeroState_DoubleAttack                           // 9 连击
	HeroState_Stealth                                // 10 隐匿
	HeroState_Injury                                 // 11 重伤
	HeroState_Poison                                 // 12 中毒
	HeroState_Chaos                                  // 13 混乱
	HeroState_AntiHidden                             // 14 鹰眼
	HeroState_ImmunityGroupDmg                       // 15 免疫群体伤害
	HeroState_Paralyzed                              // 16 麻痹
	HeroState_Taunt                                  // 17 嘲讽
	HeroState_End
)

const (
	Hero_Quality_Begin  int32 = iota
	Hero_Quality_White  int32 = iota - 1 // 白
	Hero_Quality_Green                   // 绿
	Hero_Quality_Blue                    // 蓝
	Hero_Quality_Purple                  // 紫
	Hero_Quality_Orange                  // 橙
	Hero_Quality_Red                     // 红
	Hero_Quality_End
)

const (
	Hero_Max_Promote_Times = 6 // 突破次数上限
)

// 英雄信息
type HeroInfo struct {
	Id           int64 `bson:"_id" json:"_id"`
	OwnerId      int64 `bson:"owner_id" json:"owner_id"`
	OwnerType    int32 `bson:"owner_type" json:"owner_type"`
	TypeId       int32 `bson:"type_id" json:"type_id"`
	Exp          int32 `bson:"exp" json:"exp"`
	Level        int16 `bson:"level" json:"level"`
	PromoteLevel int8  `bson:"promote_level" json:"promote_level"` // 突破等级
	Star         int8  `bson:"star" json:"star"`                   // 星级
	Friendship   int32 `bson:"friendship" json:"friendship"`       // 友好度
	FashionId    int32 `bson:"fashion_id" json:"fashion_id"`       // 时装id
}
