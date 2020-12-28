package define

const (
	Hero_MaxEquip = 4 // how many equips can put on per hero
)

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
