package config

// 技能表-拳皇/SkillEnum
type TargetType int

const (
	TargetType_BEGIN              TargetType = 1   //枚举开始
	TargetType_EnemyAll           TargetType = 101 //敌方全体
	TargetType_EnemyFront         TargetType = 102 //敌方前排
	TargetType_EnemyBack          TargetType = 103 //敌方后排
	TargetType_EnemyMonster       TargetType = 104 //敌方怪物
	TargetType_EnemyRomte         TargetType = 105 //敌方远程
	TargetType_EnemyMelee         TargetType = 106 //敌方近战
	TargetType_Self               TargetType = 200 //自身
	TargetType_SelfAll            TargetType = 201 //己方全体
	TargetType_SelfFont           TargetType = 202 //己方前排
	TargetType_SelfBack           TargetType = 203 //己方后排
	TargetType_SelfMonster        TargetType = 204 //己方怪物
	TargetType_SelfCampExcudeSelf TargetType = 205 //己方除自己
	TargetType_SelfRomte          TargetType = 206 //己方远程
	TargetType_SelfMelee          TargetType = 207 //己方近战
	TargetType_Wall               TargetType = 8   //墙体
	TargetType_ResPoint           TargetType = 9   //采集点
	TargetType_All                TargetType = 11  //场上全部
	TargetType_Support            TargetType = 12  //援护对象
	TargetType_PlayerFront        TargetType = 13  //玩家前排
	TargetType_PlayerBack         TargetType = 14  //玩家后排
	TargetType_PlayerRemote       TargetType = 15  //玩家远程
	TargetType_PlayerMelee        TargetType = 16  //玩家近战
	TargetType_NpcExcudeSelf      TargetType = 17  //Npc除自身
	TargetType_PlayerAll          TargetType = 18  //玩家全体
	TargetType_GlobalMan          TargetType = 19  //机器人
	TargetType_SelfGlobalMan      TargetType = 20  //己方机器人
	TargetType_EnemyGlobalMan     TargetType = 21  //敌方机器人
	TargetType_AllExcludeSelf     TargetType = 22  //全部除自己
	TargetType_END                TargetType = 208 //枚举结束
)

// 技能表-拳皇/SkillEnum
type LockType int

const (
	LockType_BEGIN     LockType = 1 //枚举开始
	LockType_Target    LockType = 1 //锁定目标
	LockType_TargetPos LockType = 2 //锁定目标位置
	LockType_END       LockType = 3 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillPerfromType int

const (
	SkillPerfromType_BEGIN         SkillPerfromType = 1 //枚举开始
	SkillPerfromType_Instantaneous SkillPerfromType = 1 //瞬发
	SkillPerfromType_Sing          SkillPerfromType = 2 //吟唱
	SkillPerfromType_Guide         SkillPerfromType = 3 //引导
	SkillPerfromType_END           SkillPerfromType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SequenceTriggerType int

const (
	SequenceTriggerType_BEGIN            SequenceTriggerType = 1  //枚举开始
	SequenceTriggerType_EnterLevel       SequenceTriggerType = 1  //进入关卡
	SequenceTriggerType_EnterBattle      SequenceTriggerType = 2  //进入战斗
	SequenceTriggerType_HpLessThan       SequenceTriggerType = 3  //生命值低于
	SequenceTriggerType_DamageOver       SequenceTriggerType = 4  //伤害超过
	SequenceTriggerType_EnemyDead        SequenceTriggerType = 5  //敌人死亡
	SequenceTriggerType_EnterStage       SequenceTriggerType = 6  //进入阶段
	SequenceTriggerType_StageTimeOver    SequenceTriggerType = 7  //阶段时长超过
	SequenceTriggerType_DeadRelease      SequenceTriggerType = 8  //死亡释放
	SequenceTriggerType_BattleTimeOver   SequenceTriggerType = 9  //战斗时间超过
	SequenceTriggerType_OnNormalAttack   SequenceTriggerType = 10 //普攻发生
	SequenceTriggerType_BeHitMoreThan    SequenceTriggerType = 11 //受击超过
	SequenceTriggerType_SpMax            SequenceTriggerType = 12 //能量值满
	SequenceTriggerType_BeAimed          SequenceTriggerType = 13 //被技能锁定
	SequenceTriggerType_CoolDown         SequenceTriggerType = 14 //冷却完成
	SequenceTriggerType_CDAndEnterBattle SequenceTriggerType = 16 //进入战斗且不在CD中
	SequenceTriggerType_BeHit            SequenceTriggerType = 17 //受击触发
	SequenceTriggerType_END              SequenceTriggerType = 18 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillType int

const (
	SkillType_BEGIN    SkillType = 1 //枚举开始
	SkillType_Positive SkillType = 1 //主动
	SkillType_Negative SkillType = 2 //被动
	SkillType_Attr     SkillType = 3 //属性
	SkillType_END      SkillType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillSequenceExcuteAndCDType int

const (
	SkillSequenceExcuteAndCDType_BEGIN       SkillSequenceExcuteAndCDType = 1 //枚举开始
	SkillSequenceExcuteAndCDType_Timing      SkillSequenceExcuteAndCDType = 1 //定时
	SkillSequenceExcuteAndCDType_Probability SkillSequenceExcuteAndCDType = 2 //概率
	SkillSequenceExcuteAndCDType_PerDamage   SkillSequenceExcuteAndCDType = 3 //生命每降低
	SkillSequenceExcuteAndCDType_TimingRange SkillSequenceExcuteAndCDType = 4 //范围随机
	SkillSequenceExcuteAndCDType_END         SkillSequenceExcuteAndCDType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillSiftTargetType int

const (
	SkillSiftTargetType_BEGIN                  SkillSiftTargetType = 1  //枚举开始
	SkillSiftTargetType_Nearest                SkillSiftTargetType = 1  //最近
	SkillSiftTargetType_Farest                 SkillSiftTargetType = 2  //最远
	SkillSiftTargetType_Random                 SkillSiftTargetType = 3  //随机
	SkillSiftTargetType_UnSummon               SkillSiftTargetType = 4  //非召唤
	SkillSiftTargetType_Summon                 SkillSiftTargetType = 5  //召唤
	SkillSiftTargetType_MaxDamage              SkillSiftTargetType = 8  //伤害最高
	SkillSiftTargetType_ExculdeSelf            SkillSiftTargetType = 11 //除自己
	SkillSiftTargetType_MinHp                  SkillSiftTargetType = 12 //血量最少
	SkillSiftTargetType_WithBuff               SkillSiftTargetType = 13 //处于Buff状态
	SkillSiftTargetType_MinSp                  SkillSiftTargetType = 14 //能量最少
	SkillSiftTargetType_MaxHp                  SkillSiftTargetType = 15 //血量最多
	SkillSiftTargetType_NotBoss                SkillSiftTargetType = 16 //非Boss
	SkillSiftTargetType_UseSp                  SkillSiftTargetType = 17 //使用能量者
	SkillSiftTargetType_Wall                   SkillSiftTargetType = 18 //墙体
	SkillSiftTargetType_PosRange               SkillSiftTargetType = 19 //位置区间
	SkillSiftTargetType_PosList                SkillSiftTargetType = 20 //位置列表
	SkillSiftTargetType_MaxDefence             SkillSiftTargetType = 21 //防御最高
	SkillSiftTargetType_MinDefence             SkillSiftTargetType = 22 //防御最低
	SkillSiftTargetType_HeadFirst              SkillSiftTargetType = 23 //前排优先
	SkillSiftTargetType_TailFrist              SkillSiftTargetType = 24 //后排优先
	SkillSiftTargetType_MaxSp                  SkillSiftTargetType = 25 //能量最多
	SkillSiftTargetType_SameTeam               SkillSiftTargetType = 26 //同队伍
	SkillSiftTargetType_Carrer                 SkillSiftTargetType = 27 //职业
	SkillSiftTargetType_SameTeamCarrer         SkillSiftTargetType = 28 //同队伍职业筛选
	SkillSiftTargetType_CertainHero            SkillSiftTargetType = 29 //指定英雄
	SkillSiftTargetType_FixedHeroPriority      SkillSiftTargetType = 30 //固定英雄优先
	SkillSiftTargetType_MinHpPrecent           SkillSiftTargetType = 31 //血量百分比最少
	SkillSiftTargetType_MaxPower               SkillSiftTargetType = 32 //最高战力
	SkillSiftTargetType_SameTeamMaxPower       SkillSiftTargetType = 33 //同队最高战力
	SkillSiftTargetType_SameTeamHeadFirst      SkillSiftTargetType = 34 //同队站位靠前
	SkillSiftTargetType_SameTeamTailFirst      SkillSiftTargetType = 35 //同队站位靠后
	SkillSiftTargetType_KofStarHeadFirst       SkillSiftTargetType = 36 //大明星前排优先
	SkillSiftTargetType_NotBossAndMinHpPrecent SkillSiftTargetType = 37 //非Boss且血量最少
	SkillSiftTargetType_ForntOfMe              SkillSiftTargetType = 38 //我及我前方
	SkillSiftTargetType_BehindOfMe             SkillSiftTargetType = 39 //我及我后方
	SkillSiftTargetType_END                    SkillSiftTargetType = 40 //枚举结束
)

// 技能表-拳皇/SkillEnum
type DamageType int

const (
	DamageType_BEGIN    DamageType = 0 //枚举开始
	DamageType_All      DamageType = 0 //全种类
	DamageType_Physical DamageType = 1 //物理伤害
	DamageType_Magic    DamageType = 2 //法术伤害
	DamageType_TrueDmg  DamageType = 3 //真实伤害
	DamageType_Fire     DamageType = 4 //火系伤害
	DamageType_END      DamageType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type DamageCalculateType int

const (
	DamageCalculateType_BEGIN       DamageCalculateType = 1 //枚举开始
	DamageCalculateType_HpPercent   DamageCalculateType = 1 //血量百分比
	DamageCalculateType_MaxPercent  DamageCalculateType = 2 //血量上限百分比
	DamageCalculateType_AttrPercent DamageCalculateType = 3 //属性百分比
	DamageCalculateType_END         DamageCalculateType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SubDamageType int

const (
	SubDamageType_BEGIN SubDamageType = 1 //枚举开始
	SubDamageType_Level SubDamageType = 1 //等级
	SubDamageType_Fixed SubDamageType = 2 //固定
	SubDamageType_END   SubDamageType = 3 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillLockType int

const (
	SkillLockType_BEGIN   SkillLockType = 1 //枚举开始
	SkillLockType_Target  SkillLockType = 1 //锁目标
	SkillLockType_Postion SkillLockType = 2 //锁位置
	SkillLockType_END     SkillLockType = 3 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillQuaility int

const (
	SkillQuaility_BEGIN      SkillQuaility = -1 //枚举开始
	SkillQuaility_Collect    SkillQuaility = -1 //采集
	SkillQuaility_Normal     SkillQuaility = 1  //普通
	SkillQuaility_Excellence SkillQuaility = 2  //卓越
	SkillQuaility_Rare       SkillQuaility = 3  //稀有
	SkillQuaility_Epic       SkillQuaility = 4  //史诗
	SkillQuaility_Legend     SkillQuaility = 5  //传说
	SkillQuaility_END        SkillQuaility = 6  //枚举结束
)

// 技能表-拳皇/SkillEnum
type BuffType int

const (
	BuffType_BEGIN                       BuffType = -2 //枚举开始
	BuffType_Stun                        BuffType = 1  //眩晕
	BuffType_Freeze                      BuffType = 2  //冰冻
	BuffType_Silent                      BuffType = 3  //沉默
	BuffType_Blood                       BuffType = 4  //流血
	BuffType_Poison                      BuffType = 5  //中毒
	BuffType_Burn                        BuffType = 6  //灼烧
	BuffType_Positive                    BuffType = 7  //正面
	BuffType_Negative                    BuffType = 8  //负面
	BuffType_Control                     BuffType = 9  //控制
	BuffType_Fixation                    BuffType = 10 //定身
	BuffType_Decelerate                  BuffType = 11 //减速
	BuffType_Harm                        BuffType = 12 //伤害
	BuffType_Heal                        BuffType = 13 //治愈
	BuffType_Special                     BuffType = 14 //特殊
	BuffType_Frozen                      BuffType = 15 //冻结
	BuffType_Lianna                      BuffType = 16 //莉安娜吸血
	BuffType_Shichuizhenwu               BuffType = 17 //矢吹真吾灼烧
	BuffType_TimePasue                   BuffType = 18 //时停
	BuffType_Banqiliangxiaozhaozhuangtai BuffType = 19 //坂崎良小招状态
	BuffType_Oumijia1                    BuffType = 20 //欧米茄充能
	BuffType_Oumijia2                    BuffType = 21 //欧米茄充能协同
	BuffType_AllDizzy                    BuffType = -1 //全部
	BuffType_None                        BuffType = -2 //无
	BuffType_END                         BuffType = 22 //枚举结束
)

// 技能表-拳皇/SkillEnum
type BuffGroupType int

const (
	BuffGroupType_BEGIN        BuffGroupType = 1 //枚举开始
	BuffGroupType_Unique       BuffGroupType = 1 //唯一
	BuffGroupType_Scattered    BuffGroupType = 2 //分散
	BuffGroupType_Concentrated BuffGroupType = 3 //集中
	BuffGroupType_Share        BuffGroupType = 4 //共享
	BuffGroupType_END          BuffGroupType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type BuffFlagType int

const (
	BuffFlagType_BEGIN             BuffFlagType = 1  //枚举开始
	BuffFlagType_Dizzy             BuffFlagType = 1  //控制
	BuffFlagType_Endure            BuffFlagType = 2  //霸体
	BuffFlagType_Invincible        BuffFlagType = 3  //免伤
	BuffFlagType_BadlyHurt         BuffFlagType = 4  //重伤
	BuffFlagType_Silent            BuffFlagType = 5  //沉默
	BuffFlagType_Freeze            BuffFlagType = 6  //冰冻
	BuffFlagType_Snare             BuffFlagType = 7  //定身
	BuffFlagType_Disarm            BuffFlagType = 8  //缴械
	BuffFlagType_Taunt             BuffFlagType = 9  //嘲讽
	BuffFlagType_Charm             BuffFlagType = 10 //魅惑
	BuffFlagType_Bokeh             BuffFlagType = 11 //虚化
	BuffFlagType_Support           BuffFlagType = 12 //援护
	BuffFlagType_ResistDeath       BuffFlagType = 13 //抗死
	BuffFlagType_HunterSign        BuffFlagType = 14 //标记
	BuffFlagType_Stealth           BuffFlagType = 15 //隐身
	BuffFlagType_Imprison          BuffFlagType = 16 //禁锢
	BuffFlagType_LimtSp            BuffFlagType = 17 //SP限制
	BuffFlagType_Frozen            BuffFlagType = 18 //冻结
	BuffFlagType_Transfiguration   BuffFlagType = 19 //变身
	BuffFlagType_DeTrans           BuffFlagType = 20 //抗变身
	BuffFlagType_ChangeNormalSkill BuffFlagType = 21 //变更普攻
	BuffFlagType_DeChange          BuffFlagType = 22 //禁普攻改变
	BuffFlagType_ScaleChange       BuffFlagType = 23 //放缩
	BuffFlagType_ForbidenRotate    BuffFlagType = 24 //禁转向
	BuffFlagType_ForbidenSpeedDown BuffFlagType = 25 //禁减速
	BuffFlagType_BreakEndrue       BuffFlagType = 26 //破霸体
	BuffFlagType_ChargeImmunity    BuffFlagType = 27 //冲锋免疫
	BuffFlagType_ForbidenHeal      BuffFlagType = 28 //禁疗
	BuffFlagType_KyoBreak          BuffFlagType = 29 //草薙流识破
	BuffFlagType_NoDOT             BuffFlagType = 30 //DOT免疫
	BuffFlagType_NoDedcut          BuffFlagType = 31 //Dedcut免疫
	BuffFlagType_AllTrueDamage     BuffFlagType = 32 //绝对真伤
	BuffFlagType_IgnoreArmor       BuffFlagType = 33 //护盾失效
	BuffFlagType_END               BuffFlagType = 34 //枚举结束
)

// 技能表-拳皇/SkillEnum
type ShockType int

const (
	ShockType_BEGIN   ShockType = 1 //枚举开始
	ShockType_Rect    ShockType = 1 //矩形
	ShockType_Sector  ShockType = 2 //扇形
	ShockType_Annular ShockType = 3 //环形
	ShockType_END     ShockType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type ShieldBrokenType int

const (
	ShieldBrokenType_BEGIN       ShieldBrokenType = 1 //枚举开始
	ShieldBrokenType_EffectTimes ShieldBrokenType = 1 //次数
	ShieldBrokenType_Time        ShieldBrokenType = 2 //时间
	ShieldBrokenType_SheildValue ShieldBrokenType = 3 //吸收值
	ShieldBrokenType_Puppet      ShieldBrokenType = 4 //娃娃盾
	ShieldBrokenType_END         ShieldBrokenType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectPlayType int

const (
	EffectPlayType_BEGIN    EffectPlayType = 1 //枚举开始
	EffectPlayType_None     EffectPlayType = 1 //无
	EffectPlayType_Sprite   EffectPlayType = 2 //图片
	EffectPlayType_FramSeq  EffectPlayType = 3 //序列帧
	EffectPlayType_Skeleton EffectPlayType = 4 //动画
	EffectPlayType_END      EffectPlayType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectOwnerType int

const (
	EffectOwnerType_BEGIN EffectOwnerType = 1 //枚举开始
	EffectOwnerType_Skill EffectOwnerType = 1 //技能
	EffectOwnerType_Buff  EffectOwnerType = 2 //状态
	EffectOwnerType_END   EffectOwnerType = 3 //枚举结束
)

// 技能表-拳皇/SkillEnum
type BulletMoveType int

const (
	BulletMoveType_BEGIN           BulletMoveType = 1 //枚举开始
	BulletMoveType_Line            BulletMoveType = 1 //直线
	BulletMoveType_Parabola        BulletMoveType = 2 //抛物线
	BulletMoveType_ReverseParbola  BulletMoveType = 3 //反抛物线
	BulletMoveType_Ballistics      BulletMoveType = 4 //弹道
	BulletMoveType_FaceBallistics  BulletMoveType = 5 //面向弹道
	BulletMoveType_GroudBallistics BulletMoveType = 6 //贴地弹道
	BulletMoveType_END             BulletMoveType = 7 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillMoveTargetType int

const (
	SkillMoveTargetType_BEGIN         SkillMoveTargetType = 1 //枚举开始
	SkillMoveTargetType_TargetPos     SkillMoveTargetType = 1 //目标位置
	SkillMoveTargetType_SelfPos       SkillMoveTargetType = 2 //自身位置
	SkillMoveTargetType_Escape        SkillMoveTargetType = 3 //逃逸位置
	SkillMoveTargetType_AbsoluteAngle SkillMoveTargetType = 4 //固定角度
	SkillMoveTargetType_END           SkillMoveTargetType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillMoveDirectionType int

const (
	SkillMoveDirectionType_BEGIN         SkillMoveDirectionType = 1 //枚举开始
	SkillMoveDirectionType_TargetForward SkillMoveDirectionType = 1 //目标正向
	SkillMoveDirectionType_SelfForward   SkillMoveDirectionType = 2 //自身正向
	SkillMoveDirectionType_END           SkillMoveDirectionType = 3 //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectRoleType int

const (
	EffectRoleType_BEGIN                  EffectRoleType = 1  //枚举开始
	EffectRoleType_SkillTarget            EffectRoleType = 1  //技能目标
	EffectRoleType_SkillPreform           EffectRoleType = 2  //安装者
	EffectRoleType_EnemyAll               EffectRoleType = 3  //敌方全体
	EffectRoleType_EnemyFront             EffectRoleType = 4  //敌方前排
	EffectRoleType_EnemyBack              EffectRoleType = 5  //敌方后排
	EffectRoleType_SelfAll                EffectRoleType = 6  //己方全体
	EffectRoleType_SelfFront              EffectRoleType = 7  //己方前排
	EffectRoleType_SelfBack               EffectRoleType = 8  //己方后排
	EffectRoleType_All                    EffectRoleType = 9  //场上全部
	EffectRoleType_Self                   EffectRoleType = 10 //自身
	EffectRoleType_Wall                   EffectRoleType = 11 //墙体
	EffectRoleType_BuffOwner              EffectRoleType = 12 //拥有者
	EffectRoleType_CheckTarget            EffectRoleType = 13 //检测对象
	EffectRoleType_CurAttack              EffectRoleType = 14 //当前目标
	EffectRoleType_MySummon               EffectRoleType = 15 //自身召唤物
	EffectRoleType_Master                 EffectRoleType = 16 //召唤者
	EffectRoleType_Attacker               EffectRoleType = 17 //攻击者
	EffectRoleType_Support                EffectRoleType = 18 //援护对象
	EffectRoleType_AllPlayers             EffectRoleType = 19 //玩家全体
	EffectRoleType_SelfExcudeSelf         EffectRoleType = 20 //己方除自己
	EffectRoleType_AllExcludeSelf         EffectRoleType = 21 //全部除自己
	EffectRoleType_EnemyMonster           EffectRoleType = 22 //敌方怪物
	EffectRoleType_SelfMonster            EffectRoleType = 23 //己方怪物
	EffectRoleType_EnemyAllWithWall       EffectRoleType = 24 //敌方全体以及城墙
	EffectRoleType_SelfWall               EffectRoleType = 25 //己方墙体
	EffectRoleType_SelfMonsterWithoutHunt EffectRoleType = 26 //己方怪物除狩猎
	EffectRoleType_GlobalMan              EffectRoleType = 27 //机器人
	EffectRoleType_SelfGlobalMan          EffectRoleType = 28 //己方机器人
	EffectRoleType_EnemyGlobalMan         EffectRoleType = 29 //敌方机器人
	EffectRoleType_LockTarget             EffectRoleType = 30 //当前锁定目标
	EffectRoleType_END                    EffectRoleType = 31 //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectDirType int

const (
	EffectDirType_BEGIN       EffectDirType = 1 //枚举开始
	EffectDirType_None        EffectDirType = 1 //无朝向
	EffectDirType_Target      EffectDirType = 2 //目标朝向
	EffectDirType_Self        EffectDirType = 3 //自身朝向
	EffectDirType_SelfAndFilp EffectDirType = 4 //自身朝向与翻转
	EffectDirType_LogicDir    EffectDirType = 5 //逻辑朝向
	EffectDirType_END         EffectDirType = 6 //枚举结束
)

// 技能表-拳皇/SkillEnum
type BuffFixedDamageType int

const (
	BuffFixedDamageType_BEGIN   BuffFixedDamageType = 1 //枚举开始
	BuffFixedDamageType_Sub     BuffFixedDamageType = 1 //额外附加
	BuffFixedDamageType_Replace BuffFixedDamageType = 2 //变为
	BuffFixedDamageType_END     BuffFixedDamageType = 3 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SummonCampType int

const (
	SummonCampType_BEGIN           SummonCampType = 1 //枚举开始
	SummonCampType_SameWithMaster  SummonCampType = 1 //同主人
	SummonCampType_EnemyWithMaster SummonCampType = 2 //与主人敌对
	SummonCampType_ThirdCamp       SummonCampType = 3 //第三方阵营
	SummonCampType_None            SummonCampType = 4 //无阵营
	SummonCampType_END             SummonCampType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SummonMoveType int

const (
	SummonMoveType_BEGIN  SummonMoveType = 1 //枚举开始
	SummonMoveType_Follow SummonMoveType = 1 //跟随
	SummonMoveType_Auto   SummonMoveType = 2 //自动
	SummonMoveType_None   SummonMoveType = 3 //不动
	SummonMoveType_END    SummonMoveType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type ReplaceTargetType int

const (
	ReplaceTargetType_BEGIN  ReplaceTargetType = -1 //枚举开始
	ReplaceTargetType_Circle ReplaceTargetType = 1  //圆形范围
	ReplaceTargetType_Sector ReplaceTargetType = 2  //扇形范围
	ReplaceTargetType_None   ReplaceTargetType = -1 //不替换
	ReplaceTargetType_END    ReplaceTargetType = 3  //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectShapeType int

const (
	EffectShapeType_BEGIN  EffectShapeType = 1 //枚举开始
	EffectShapeType_Circle EffectShapeType = 1 //圆形
	EffectShapeType_Rect   EffectShapeType = 2 //矩形
	EffectShapeType_Sector EffectShapeType = 3 //扇形
	EffectShapeType_END    EffectShapeType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SummonShapCopyType int

const (
	SummonShapCopyType_BEGIN  SummonShapCopyType = -1 //枚举开始
	SummonShapCopyType_Self   SummonShapCopyType = 1  //自己
	SummonShapCopyType_Target SummonShapCopyType = 2  //目标
	SummonShapCopyType_None   SummonShapCopyType = -1 //不复制
	SummonShapCopyType_END    SummonShapCopyType = 3  //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectLayerType int

const (
	EffectLayerType_BEGIN EffectLayerType = 1 //枚举开始
	EffectLayerType_UP    EffectLayerType = 1 //上特效层
	EffectLayerType_Down  EffectLayerType = 2 //下特效层
	EffectLayerType_UI    EffectLayerType = 3 //UI层
	EffectLayerType_Black EffectLayerType = 4 //黑遮罩
	EffectLayerType_END   EffectLayerType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectPosType int

const (
	EffectPosType_BEGIN     EffectPosType = 1 //枚举开始
	EffectPosType_Role      EffectPosType = 1 //角色
	EffectPosType_HitPoint  EffectPosType = 2 //受击点
	EffectPosType_HpPos     EffectPosType = 3 //血条
	EffectPosType_FirePoint EffectPosType = 4 //开火点
	EffectPosType_END       EffectPosType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type BattleRoleCamp int

const (
	BattleRoleCamp_BEGIN       BattleRoleCamp = 1  //枚举开始
	BattleRoleCamp_Player      BattleRoleCamp = 1  //玩家
	BattleRoleCamp_Enemy       BattleRoleCamp = 2  //敌方
	BattleRoleCamp_Third       BattleRoleCamp = 3  //第三方
	BattleRoleCamp_NPC         BattleRoleCamp = 4  //第三方Npc
	BattleRoleCamp_FriendNpc   BattleRoleCamp = 5  //友方NPC
	BattleRoleCamp_Friend      BattleRoleCamp = 6  //友方
	BattleRoleCamp_FriendEnemy BattleRoleCamp = 7  //友方敌人
	BattleRoleCamp_SelfHero    BattleRoleCamp = 8  //我方英雄
	BattleRoleCamp_EnemyHero   BattleRoleCamp = 9  //敌方英雄
	BattleRoleCamp_Max         BattleRoleCamp = 10 //
	BattleRoleCamp_END         BattleRoleCamp = 11 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillIconType int

const (
	SkillIconType_BEGIN    SkillIconType = 1 //枚举开始
	SkillIconType_Healing  SkillIconType = 1 //治
	SkillIconType_Negative SkillIconType = 2 //被
	SkillIconType_Ultimate SkillIconType = 3 //奥
	SkillIconType_Normal   SkillIconType = 4 //普
	SkillIconType_Positive SkillIconType = 5 //主
	SkillIconType_Fetters  SkillIconType = 6 //羁
	SkillIconType_END      SkillIconType = 7 //枚举结束
)

// 技能表-拳皇/SkillEnum
type ConnectionType int

const (
	ConnectionType_BEGIN     ConnectionType = 1 //枚举开始
	ConnectionType_Flash     ConnectionType = 1 //闪现连携
	ConnectionType_Move      ConnectionType = 2 //移动连携
	ConnectionType_Immediate ConnectionType = 3 //瞬时连携
	ConnectionType_END       ConnectionType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type HitStunOverlay int

const (
	HitStunOverlay_BEGIN       HitStunOverlay = 1 //枚举开始
	HitStunOverlay_Cannot      HitStunOverlay = 1 //不可操作
	HitStunOverlay_Replace     HitStunOverlay = 2 //替换
	HitStunOverlay_Overlay     HitStunOverlay = 3 //可叠加
	HitStunOverlay_SameOverlay HitStunOverlay = 4 //同种可叠加
	HitStunOverlay_SameReplace HitStunOverlay = 5 //同种可替换
	HitStunOverlay_END         HitStunOverlay = 6 //枚举结束
)

// 技能表-拳皇/SkillEnum
type HitStunType int

const (
	HitStunType_BEGIN        HitStunType = 1 //枚举开始
	HitStunType_HitBack      HitStunType = 1 //击退
	HitStunType_LockedUp     HitStunType = 2 //击飞
	HitStunType_HitDown      HitStunType = 3 //击倒
	HitStunType_SkillHitBack HitStunType = 4 //技能特殊击退
	HitStunType_END          HitStunType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SpecialStun int

const (
	SpecialStun_BEGIN SpecialStun = 1 //枚举开始
	SpecialStun_Guard SpecialStun = 6 //防御硬直
	SpecialStun_Hurt  SpecialStun = 7 //受击硬直
	SpecialStun_END   SpecialStun = 8 //枚举结束
)

// 技能表-拳皇/SkillEnum
type JumpBackType int

const (
	JumpBackType_BEGIN      JumpBackType = -1 //枚举开始
	JumpBackType_None       JumpBackType = -1 //无
	JumpBackType_Immediate  JumpBackType = 1  //立即
	JumpBackType_AfterCombo JumpBackType = 2  //连击后
	JumpBackType_END        JumpBackType = 3  //枚举结束
)

// 技能表-拳皇/SkillEnum
type ConnectCameraType int

const (
	ConnectCameraType_BEGIN     ConnectCameraType = -1 //枚举开始
	ConnectCameraType_None      ConnectCameraType = -1 //无
	ConnectCameraType_Lerp      ConnectCameraType = 1  //缓动
	ConnectCameraType_Immeidate ConnectCameraType = 2  //立即
	ConnectCameraType_END       ConnectCameraType = 3  //枚举结束
)

// 技能表-拳皇/SkillEnum
type GrabPointType int

const (
	GrabPointType_BEGIN GrabPointType = -1 //枚举开始
	GrabPointType_Check GrabPointType = -1 //胸
	GrabPointType_Head  GrabPointType = 1  //头
	GrabPointType_END   GrabPointType = 2  //枚举结束
)

// 技能表-拳皇/SkillEnum
type ImportSkillID int

const (
	ImportSkillID_BEGIN             ImportSkillID = 1    //枚举开始
	ImportSkillID_CommonEndrue      ImportSkillID = 3    //通用受击霸体
	ImportSkillID_CommoRemoteNormal ImportSkillID = 5    //远程普攻检测
	ImportSkillID_CounterAttack     ImportSkillID = 11   //反击技能
	ImportSkillID_ComboAttack       ImportSkillID = 12   //连击技能
	ImportSkillID_FirstMoveGain     ImportSkillID = 1701 //先手增益
	ImportSkillID_CNT               ImportSkillID = 14   //结束
	ImportSkillID_END               ImportSkillID = 1702 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SkillTagType int

const (
	SkillTagType_BEGIN            SkillTagType = -1 //枚举开始
	SkillTagType_NormalAttack     SkillTagType = 1  //普攻
	SkillTagType_SkillA           SkillTagType = 2  //小招A
	SkillTagType_SkillB           SkillTagType = 3  //小招B
	SkillTagType_Ultimate         SkillTagType = 4  //大招
	SkillTagType_UltimateNext     SkillTagType = 5  //大招后续
	SkillTagType_Negative         SkillTagType = 6  //被动
	SkillTagType_Connect          SkillTagType = 7  //连携
	SkillTagType_SkillC           SkillTagType = 8  //小招C
	SkillTagType_NoCDReplace      SkillTagType = 9  //无冷却普攻替换
	SkillTagType_AtrSkillA        SkillTagType = 10 //奥义A
	SkillTagType_AtrSkillB        SkillTagType = 11 //奥义B
	SkillTagType_AssistSkillA     SkillTagType = 12 //援护A
	SkillTagType_AssistSkillB     SkillTagType = 13 //援护B
	SkillTagType_AssistSkillC     SkillTagType = 14 //援护C
	SkillTagType_PetSkillA        SkillTagType = 15 //宠物技能A
	SkillTagType_PetSkillB        SkillTagType = 16 //宠物技能B
	SkillTagType_PetSkillManifest SkillTagType = 17 //宠物技能幻化
	SkillTagType_None             SkillTagType = -1 //无
	SkillTagType_END              SkillTagType = 18 //枚举结束
)

// 技能表-拳皇/SkillEnum
type ConnectFocusType int

const (
	ConnectFocusType_BEGIN  ConnectFocusType = -1 //枚举开始
	ConnectFocusType_Form   ConnectFocusType = -1 //发起者
	ConnectFocusType_Target ConnectFocusType = 1  //目标
	ConnectFocusType_Center ConnectFocusType = 2  //中点
	ConnectFocusType_END    ConnectFocusType = 3  //枚举结束
)

// 技能表-拳皇/SkillEnum
type ComboFntShowType int

const (
	ComboFntShowType_BEGIN        ComboFntShowType = 1 //枚举开始
	ComboFntShowType_AmplifyScale ComboFntShowType = 1 //放缩
	ComboFntShowType_Flash        ComboFntShowType = 2 //位移
	ComboFntShowType_Appear       ComboFntShowType = 3 //渐显
	ComboFntShowType_END          ComboFntShowType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type ComboFntMoveInDir int

const (
	ComboFntMoveInDir_BEGIN     ComboFntMoveInDir = 1 //枚举开始
	ComboFntMoveInDir_Up        ComboFntMoveInDir = 1 //上
	ComboFntMoveInDir_Down      ComboFntMoveInDir = 2 //下
	ComboFntMoveInDir_Left      ComboFntMoveInDir = 3 //左
	ComboFntMoveInDir_Right     ComboFntMoveInDir = 4 //右
	ComboFntMoveInDir_LeftUp    ComboFntMoveInDir = 5 //左上
	ComboFntMoveInDir_RightUp   ComboFntMoveInDir = 6 //右上
	ComboFntMoveInDir_LeftDown  ComboFntMoveInDir = 7 //左下
	ComboFntMoveInDir_RightDown ComboFntMoveInDir = 8 //右下
	ComboFntMoveInDir_END       ComboFntMoveInDir = 9 //枚举结束
)

// 技能表-拳皇/SkillEnum
type BattleAreaType int

const (
	BattleAreaType_BEGIN   BattleAreaType = 1 //枚举开始
	BattleAreaType_Left    BattleAreaType = 1 //左区
	BattleAreaType_Midddle BattleAreaType = 2 //中区
	BattleAreaType_Right   BattleAreaType = 3 //右区
	BattleAreaType_Self    BattleAreaType = 4 //己方区
	BattleAreaType_Enemy   BattleAreaType = 5 //敌方区
	BattleAreaType_END     BattleAreaType = 6 //枚举结束
)

// 技能表-拳皇/SkillEnum
type CameraLockType int

const (
	CameraLockType_BEGIN  CameraLockType = 1 //枚举开始
	CameraLockType_From   CameraLockType = 1 //发起者
	CameraLockType_Target CameraLockType = 2 //目标
	CameraLockType_Center CameraLockType = 3 //中心
	CameraLockType_END    CameraLockType = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type EffectAddParam int

const (
	EffectAddParam_BEGIN    EffectAddParam = 1 //枚举开始
	EffectAddParam_SkillDmg EffectAddParam = 1 //技伤
	EffectAddParam_Heal     EffectAddParam = 2 //技疗
	EffectAddParam_CNT      EffectAddParam = 3 //结束
	EffectAddParam_END      EffectAddParam = 4 //枚举结束
)

// 技能表-拳皇/SkillEnum
type DmgCalculateRateType int

const (
	DmgCalculateRateType_BEGIN      DmgCalculateRateType = 1 //枚举开始
	DmgCalculateRateType_Crital     DmgCalculateRateType = 1 //暴击
	DmgCalculateRateType_HealCrital DmgCalculateRateType = 2 //治疗暴击
	DmgCalculateRateType_Dodge      DmgCalculateRateType = 3 //闪避
	DmgCalculateRateType_Block      DmgCalculateRateType = 4 //格挡
	DmgCalculateRateType_END        DmgCalculateRateType = 5 //枚举结束
)

// 技能表-拳皇/SkillEnum
type SummonAppearType int

const (
	SummonAppearType_BEGIN    SummonAppearType = -1 //枚举开始
	SummonAppearType_None     SummonAppearType = -1 //无
	SummonAppearType_Straight SummonAppearType = 1  //直线冲刺
	SummonAppearType_Parabola SummonAppearType = 2  //曲线进场
	SummonAppearType_Fall     SummonAppearType = 3  //上空掉落
	SummonAppearType_END      SummonAppearType = 4  //枚举结束
)

func (e TargetType) String() string {

	if e == TargetType_EnemyAll {
		return "EnemyAll"
	}
	if e == TargetType_EnemyFront {
		return "EnemyFront"
	}
	if e == TargetType_EnemyBack {
		return "EnemyBack"
	}
	if e == TargetType_EnemyMonster {
		return "EnemyMonster"
	}
	if e == TargetType_EnemyRomte {
		return "EnemyRomte"
	}
	if e == TargetType_EnemyMelee {
		return "EnemyMelee"
	}
	if e == TargetType_Self {
		return "Self"
	}
	if e == TargetType_SelfAll {
		return "SelfAll"
	}
	if e == TargetType_SelfFont {
		return "SelfFont"
	}
	if e == TargetType_SelfBack {
		return "SelfBack"
	}
	if e == TargetType_SelfMonster {
		return "SelfMonster"
	}
	if e == TargetType_SelfCampExcudeSelf {
		return "SelfCampExcudeSelf"
	}
	if e == TargetType_SelfRomte {
		return "SelfRomte"
	}
	if e == TargetType_SelfMelee {
		return "SelfMelee"
	}
	if e == TargetType_Wall {
		return "Wall"
	}
	if e == TargetType_ResPoint {
		return "ResPoint"
	}
	if e == TargetType_All {
		return "All"
	}
	if e == TargetType_Support {
		return "Support"
	}
	if e == TargetType_PlayerFront {
		return "PlayerFront"
	}
	if e == TargetType_PlayerBack {
		return "PlayerBack"
	}
	if e == TargetType_PlayerRemote {
		return "PlayerRemote"
	}
	if e == TargetType_PlayerMelee {
		return "PlayerMelee"
	}
	if e == TargetType_NpcExcudeSelf {
		return "NpcExcudeSelf"
	}
	if e == TargetType_PlayerAll {
		return "PlayerAll"
	}
	if e == TargetType_GlobalMan {
		return "GlobalMan"
	}
	if e == TargetType_SelfGlobalMan {
		return "SelfGlobalMan"
	}
	if e == TargetType_EnemyGlobalMan {
		return "EnemyGlobalMan"
	}
	if e == TargetType_AllExcludeSelf {
		return "AllExcludeSelf"
	}

	return ""

}

func GetTargetType(tp string) (TargetType, bool) {
	switch tp {
	case "EnemyAll":
		return TargetType_EnemyAll, true
	case "EnemyFront":
		return TargetType_EnemyFront, true
	case "EnemyBack":
		return TargetType_EnemyBack, true
	case "EnemyMonster":
		return TargetType_EnemyMonster, true
	case "EnemyRomte":
		return TargetType_EnemyRomte, true
	case "EnemyMelee":
		return TargetType_EnemyMelee, true
	case "Self":
		return TargetType_Self, true
	case "SelfAll":
		return TargetType_SelfAll, true
	case "SelfFont":
		return TargetType_SelfFont, true
	case "SelfBack":
		return TargetType_SelfBack, true
	case "SelfMonster":
		return TargetType_SelfMonster, true
	case "SelfCampExcudeSelf":
		return TargetType_SelfCampExcudeSelf, true
	case "SelfRomte":
		return TargetType_SelfRomte, true
	case "SelfMelee":
		return TargetType_SelfMelee, true
	case "Wall":
		return TargetType_Wall, true
	case "ResPoint":
		return TargetType_ResPoint, true
	case "All":
		return TargetType_All, true
	case "Support":
		return TargetType_Support, true
	case "PlayerFront":
		return TargetType_PlayerFront, true
	case "PlayerBack":
		return TargetType_PlayerBack, true
	case "PlayerRemote":
		return TargetType_PlayerRemote, true
	case "PlayerMelee":
		return TargetType_PlayerMelee, true
	case "NpcExcudeSelf":
		return TargetType_NpcExcudeSelf, true
	case "PlayerAll":
		return TargetType_PlayerAll, true
	case "GlobalMan":
		return TargetType_GlobalMan, true
	case "SelfGlobalMan":
		return TargetType_SelfGlobalMan, true
	case "EnemyGlobalMan":
		return TargetType_EnemyGlobalMan, true
	case "AllExcludeSelf":
		return TargetType_AllExcludeSelf, true

	}
	return -1, false

}

func GetTargetType_ByAlias(tp string) (TargetType, bool) {
	switch tp {
	case "敌方全体":
		return TargetType_EnemyAll, true
	case "敌方前排":
		return TargetType_EnemyFront, true
	case "敌方后排":
		return TargetType_EnemyBack, true
	case "敌方怪物":
		return TargetType_EnemyMonster, true
	case "敌方远程":
		return TargetType_EnemyRomte, true
	case "敌方近战":
		return TargetType_EnemyMelee, true
	case "自身":
		return TargetType_Self, true
	case "己方全体":
		return TargetType_SelfAll, true
	case "己方前排":
		return TargetType_SelfFont, true
	case "己方后排":
		return TargetType_SelfBack, true
	case "己方怪物":
		return TargetType_SelfMonster, true
	case "己方除自己":
		return TargetType_SelfCampExcudeSelf, true
	case "己方远程":
		return TargetType_SelfRomte, true
	case "己方近战":
		return TargetType_SelfMelee, true
	case "墙体":
		return TargetType_Wall, true
	case "采集点":
		return TargetType_ResPoint, true
	case "场上全部":
		return TargetType_All, true
	case "援护对象":
		return TargetType_Support, true
	case "玩家前排":
		return TargetType_PlayerFront, true
	case "玩家后排":
		return TargetType_PlayerBack, true
	case "玩家远程":
		return TargetType_PlayerRemote, true
	case "玩家近战":
		return TargetType_PlayerMelee, true
	case "Npc除自身":
		return TargetType_NpcExcudeSelf, true
	case "玩家全体":
		return TargetType_PlayerAll, true
	case "机器人":
		return TargetType_GlobalMan, true
	case "己方机器人":
		return TargetType_SelfGlobalMan, true
	case "敌方机器人":
		return TargetType_EnemyGlobalMan, true
	case "全部除自己":
		return TargetType_AllExcludeSelf, true

	}
	return -1, false

}

func (e LockType) String() string {

	if e == LockType_Target {
		return "Target"
	}
	if e == LockType_TargetPos {
		return "TargetPos"
	}

	return ""

}

func GetLockType(tp string) (LockType, bool) {
	switch tp {
	case "Target":
		return LockType_Target, true
	case "TargetPos":
		return LockType_TargetPos, true

	}
	return -1, false

}

func GetLockType_ByAlias(tp string) (LockType, bool) {
	switch tp {
	case "锁定目标":
		return LockType_Target, true
	case "锁定目标位置":
		return LockType_TargetPos, true

	}
	return -1, false

}

func (e SkillPerfromType) String() string {

	if e == SkillPerfromType_Instantaneous {
		return "Instantaneous"
	}
	if e == SkillPerfromType_Sing {
		return "Sing"
	}
	if e == SkillPerfromType_Guide {
		return "Guide"
	}

	return ""

}

func GetSkillPerfromType(tp string) (SkillPerfromType, bool) {
	switch tp {
	case "Instantaneous":
		return SkillPerfromType_Instantaneous, true
	case "Sing":
		return SkillPerfromType_Sing, true
	case "Guide":
		return SkillPerfromType_Guide, true

	}
	return -1, false

}

func GetSkillPerfromType_ByAlias(tp string) (SkillPerfromType, bool) {
	switch tp {
	case "瞬发":
		return SkillPerfromType_Instantaneous, true
	case "吟唱":
		return SkillPerfromType_Sing, true
	case "引导":
		return SkillPerfromType_Guide, true

	}
	return -1, false

}

func (e SequenceTriggerType) String() string {

	if e == SequenceTriggerType_EnterLevel {
		return "EnterLevel"
	}
	if e == SequenceTriggerType_EnterBattle {
		return "EnterBattle"
	}
	if e == SequenceTriggerType_HpLessThan {
		return "HpLessThan"
	}
	if e == SequenceTriggerType_DamageOver {
		return "DamageOver"
	}
	if e == SequenceTriggerType_EnemyDead {
		return "EnemyDead"
	}
	if e == SequenceTriggerType_EnterStage {
		return "EnterStage"
	}
	if e == SequenceTriggerType_StageTimeOver {
		return "StageTimeOver"
	}
	if e == SequenceTriggerType_DeadRelease {
		return "DeadRelease"
	}
	if e == SequenceTriggerType_BattleTimeOver {
		return "BattleTimeOver"
	}
	if e == SequenceTriggerType_OnNormalAttack {
		return "OnNormalAttack"
	}
	if e == SequenceTriggerType_BeHitMoreThan {
		return "BeHitMoreThan"
	}
	if e == SequenceTriggerType_SpMax {
		return "SpMax"
	}
	if e == SequenceTriggerType_BeAimed {
		return "BeAimed"
	}
	if e == SequenceTriggerType_CoolDown {
		return "CoolDown"
	}
	if e == SequenceTriggerType_CDAndEnterBattle {
		return "CDAndEnterBattle"
	}
	if e == SequenceTriggerType_BeHit {
		return "BeHit"
	}

	return ""

}

func GetSequenceTriggerType(tp string) (SequenceTriggerType, bool) {
	switch tp {
	case "EnterLevel":
		return SequenceTriggerType_EnterLevel, true
	case "EnterBattle":
		return SequenceTriggerType_EnterBattle, true
	case "HpLessThan":
		return SequenceTriggerType_HpLessThan, true
	case "DamageOver":
		return SequenceTriggerType_DamageOver, true
	case "EnemyDead":
		return SequenceTriggerType_EnemyDead, true
	case "EnterStage":
		return SequenceTriggerType_EnterStage, true
	case "StageTimeOver":
		return SequenceTriggerType_StageTimeOver, true
	case "DeadRelease":
		return SequenceTriggerType_DeadRelease, true
	case "BattleTimeOver":
		return SequenceTriggerType_BattleTimeOver, true
	case "OnNormalAttack":
		return SequenceTriggerType_OnNormalAttack, true
	case "BeHitMoreThan":
		return SequenceTriggerType_BeHitMoreThan, true
	case "SpMax":
		return SequenceTriggerType_SpMax, true
	case "BeAimed":
		return SequenceTriggerType_BeAimed, true
	case "CoolDown":
		return SequenceTriggerType_CoolDown, true
	case "CDAndEnterBattle":
		return SequenceTriggerType_CDAndEnterBattle, true
	case "BeHit":
		return SequenceTriggerType_BeHit, true

	}
	return -1, false

}

func GetSequenceTriggerType_ByAlias(tp string) (SequenceTriggerType, bool) {
	switch tp {
	case "进入关卡":
		return SequenceTriggerType_EnterLevel, true
	case "进入战斗":
		return SequenceTriggerType_EnterBattle, true
	case "生命值低于":
		return SequenceTriggerType_HpLessThan, true
	case "伤害超过":
		return SequenceTriggerType_DamageOver, true
	case "敌人死亡":
		return SequenceTriggerType_EnemyDead, true
	case "进入阶段":
		return SequenceTriggerType_EnterStage, true
	case "阶段时长超过":
		return SequenceTriggerType_StageTimeOver, true
	case "死亡释放":
		return SequenceTriggerType_DeadRelease, true
	case "战斗时间超过":
		return SequenceTriggerType_BattleTimeOver, true
	case "普攻发生":
		return SequenceTriggerType_OnNormalAttack, true
	case "受击超过":
		return SequenceTriggerType_BeHitMoreThan, true
	case "能量值满":
		return SequenceTriggerType_SpMax, true
	case "被技能锁定":
		return SequenceTriggerType_BeAimed, true
	case "冷却完成":
		return SequenceTriggerType_CoolDown, true
	case "进入战斗且不在CD中":
		return SequenceTriggerType_CDAndEnterBattle, true
	case "受击触发":
		return SequenceTriggerType_BeHit, true

	}
	return -1, false

}

func (e SkillType) String() string {

	if e == SkillType_Positive {
		return "Positive"
	}
	if e == SkillType_Negative {
		return "Negative"
	}
	if e == SkillType_Attr {
		return "Attr"
	}

	return ""

}

func GetSkillType(tp string) (SkillType, bool) {
	switch tp {
	case "Positive":
		return SkillType_Positive, true
	case "Negative":
		return SkillType_Negative, true
	case "Attr":
		return SkillType_Attr, true

	}
	return -1, false

}

func GetSkillType_ByAlias(tp string) (SkillType, bool) {
	switch tp {
	case "主动":
		return SkillType_Positive, true
	case "被动":
		return SkillType_Negative, true
	case "属性":
		return SkillType_Attr, true

	}
	return -1, false

}

func (e SkillSequenceExcuteAndCDType) String() string {

	if e == SkillSequenceExcuteAndCDType_Timing {
		return "Timing"
	}
	if e == SkillSequenceExcuteAndCDType_Probability {
		return "Probability"
	}
	if e == SkillSequenceExcuteAndCDType_PerDamage {
		return "PerDamage"
	}
	if e == SkillSequenceExcuteAndCDType_TimingRange {
		return "TimingRange"
	}

	return ""

}

func GetSkillSequenceExcuteAndCDType(tp string) (SkillSequenceExcuteAndCDType, bool) {
	switch tp {
	case "Timing":
		return SkillSequenceExcuteAndCDType_Timing, true
	case "Probability":
		return SkillSequenceExcuteAndCDType_Probability, true
	case "PerDamage":
		return SkillSequenceExcuteAndCDType_PerDamage, true
	case "TimingRange":
		return SkillSequenceExcuteAndCDType_TimingRange, true

	}
	return -1, false

}

func GetSkillSequenceExcuteAndCDType_ByAlias(tp string) (SkillSequenceExcuteAndCDType, bool) {
	switch tp {
	case "定时":
		return SkillSequenceExcuteAndCDType_Timing, true
	case "概率":
		return SkillSequenceExcuteAndCDType_Probability, true
	case "生命每降低":
		return SkillSequenceExcuteAndCDType_PerDamage, true
	case "范围随机":
		return SkillSequenceExcuteAndCDType_TimingRange, true

	}
	return -1, false

}

func (e SkillSiftTargetType) String() string {

	if e == SkillSiftTargetType_Nearest {
		return "Nearest"
	}
	if e == SkillSiftTargetType_Farest {
		return "Farest"
	}
	if e == SkillSiftTargetType_Random {
		return "Random"
	}
	if e == SkillSiftTargetType_UnSummon {
		return "UnSummon"
	}
	if e == SkillSiftTargetType_Summon {
		return "Summon"
	}
	if e == SkillSiftTargetType_MaxDamage {
		return "MaxDamage"
	}
	if e == SkillSiftTargetType_ExculdeSelf {
		return "ExculdeSelf"
	}
	if e == SkillSiftTargetType_MinHp {
		return "MinHp"
	}
	if e == SkillSiftTargetType_WithBuff {
		return "WithBuff"
	}
	if e == SkillSiftTargetType_MinSp {
		return "MinSp"
	}
	if e == SkillSiftTargetType_MaxHp {
		return "MaxHp"
	}
	if e == SkillSiftTargetType_NotBoss {
		return "NotBoss"
	}
	if e == SkillSiftTargetType_UseSp {
		return "UseSp"
	}
	if e == SkillSiftTargetType_Wall {
		return "Wall"
	}
	if e == SkillSiftTargetType_PosRange {
		return "PosRange"
	}
	if e == SkillSiftTargetType_PosList {
		return "PosList"
	}
	if e == SkillSiftTargetType_MaxDefence {
		return "MaxDefence"
	}
	if e == SkillSiftTargetType_MinDefence {
		return "MinDefence"
	}
	if e == SkillSiftTargetType_HeadFirst {
		return "HeadFirst"
	}
	if e == SkillSiftTargetType_TailFrist {
		return "TailFrist"
	}
	if e == SkillSiftTargetType_MaxSp {
		return "MaxSp"
	}
	if e == SkillSiftTargetType_SameTeam {
		return "SameTeam"
	}
	if e == SkillSiftTargetType_Carrer {
		return "Carrer"
	}
	if e == SkillSiftTargetType_SameTeamCarrer {
		return "SameTeamCarrer"
	}
	if e == SkillSiftTargetType_CertainHero {
		return "CertainHero"
	}
	if e == SkillSiftTargetType_FixedHeroPriority {
		return "FixedHeroPriority"
	}
	if e == SkillSiftTargetType_MinHpPrecent {
		return "MinHpPrecent"
	}
	if e == SkillSiftTargetType_MaxPower {
		return "MaxPower"
	}
	if e == SkillSiftTargetType_SameTeamMaxPower {
		return "SameTeamMaxPower"
	}
	if e == SkillSiftTargetType_SameTeamHeadFirst {
		return "SameTeamHeadFirst"
	}
	if e == SkillSiftTargetType_SameTeamTailFirst {
		return "SameTeamTailFirst"
	}
	if e == SkillSiftTargetType_KofStarHeadFirst {
		return "KofStarHeadFirst"
	}
	if e == SkillSiftTargetType_NotBossAndMinHpPrecent {
		return "NotBossAndMinHpPrecent"
	}
	if e == SkillSiftTargetType_ForntOfMe {
		return "ForntOfMe"
	}
	if e == SkillSiftTargetType_BehindOfMe {
		return "BehindOfMe"
	}

	return ""

}

func GetSkillSiftTargetType(tp string) (SkillSiftTargetType, bool) {
	switch tp {
	case "Nearest":
		return SkillSiftTargetType_Nearest, true
	case "Farest":
		return SkillSiftTargetType_Farest, true
	case "Random":
		return SkillSiftTargetType_Random, true
	case "UnSummon":
		return SkillSiftTargetType_UnSummon, true
	case "Summon":
		return SkillSiftTargetType_Summon, true
	case "MaxDamage":
		return SkillSiftTargetType_MaxDamage, true
	case "ExculdeSelf":
		return SkillSiftTargetType_ExculdeSelf, true
	case "MinHp":
		return SkillSiftTargetType_MinHp, true
	case "WithBuff":
		return SkillSiftTargetType_WithBuff, true
	case "MinSp":
		return SkillSiftTargetType_MinSp, true
	case "MaxHp":
		return SkillSiftTargetType_MaxHp, true
	case "NotBoss":
		return SkillSiftTargetType_NotBoss, true
	case "UseSp":
		return SkillSiftTargetType_UseSp, true
	case "Wall":
		return SkillSiftTargetType_Wall, true
	case "PosRange":
		return SkillSiftTargetType_PosRange, true
	case "PosList":
		return SkillSiftTargetType_PosList, true
	case "MaxDefence":
		return SkillSiftTargetType_MaxDefence, true
	case "MinDefence":
		return SkillSiftTargetType_MinDefence, true
	case "HeadFirst":
		return SkillSiftTargetType_HeadFirst, true
	case "TailFrist":
		return SkillSiftTargetType_TailFrist, true
	case "MaxSp":
		return SkillSiftTargetType_MaxSp, true
	case "SameTeam":
		return SkillSiftTargetType_SameTeam, true
	case "Carrer":
		return SkillSiftTargetType_Carrer, true
	case "SameTeamCarrer":
		return SkillSiftTargetType_SameTeamCarrer, true
	case "CertainHero":
		return SkillSiftTargetType_CertainHero, true
	case "FixedHeroPriority":
		return SkillSiftTargetType_FixedHeroPriority, true
	case "MinHpPrecent":
		return SkillSiftTargetType_MinHpPrecent, true
	case "MaxPower":
		return SkillSiftTargetType_MaxPower, true
	case "SameTeamMaxPower":
		return SkillSiftTargetType_SameTeamMaxPower, true
	case "SameTeamHeadFirst":
		return SkillSiftTargetType_SameTeamHeadFirst, true
	case "SameTeamTailFirst":
		return SkillSiftTargetType_SameTeamTailFirst, true
	case "KofStarHeadFirst":
		return SkillSiftTargetType_KofStarHeadFirst, true
	case "NotBossAndMinHpPrecent":
		return SkillSiftTargetType_NotBossAndMinHpPrecent, true
	case "ForntOfMe":
		return SkillSiftTargetType_ForntOfMe, true
	case "BehindOfMe":
		return SkillSiftTargetType_BehindOfMe, true

	}
	return -1, false

}

func GetSkillSiftTargetType_ByAlias(tp string) (SkillSiftTargetType, bool) {
	switch tp {
	case "最近":
		return SkillSiftTargetType_Nearest, true
	case "最远":
		return SkillSiftTargetType_Farest, true
	case "随机":
		return SkillSiftTargetType_Random, true
	case "非召唤":
		return SkillSiftTargetType_UnSummon, true
	case "召唤":
		return SkillSiftTargetType_Summon, true
	case "伤害最高":
		return SkillSiftTargetType_MaxDamage, true
	case "除自己":
		return SkillSiftTargetType_ExculdeSelf, true
	case "血量最少":
		return SkillSiftTargetType_MinHp, true
	case "处于Buff状态":
		return SkillSiftTargetType_WithBuff, true
	case "能量最少":
		return SkillSiftTargetType_MinSp, true
	case "血量最多":
		return SkillSiftTargetType_MaxHp, true
	case "非Boss":
		return SkillSiftTargetType_NotBoss, true
	case "使用能量者":
		return SkillSiftTargetType_UseSp, true
	case "墙体":
		return SkillSiftTargetType_Wall, true
	case "位置区间":
		return SkillSiftTargetType_PosRange, true
	case "位置列表":
		return SkillSiftTargetType_PosList, true
	case "防御最高":
		return SkillSiftTargetType_MaxDefence, true
	case "防御最低":
		return SkillSiftTargetType_MinDefence, true
	case "前排优先":
		return SkillSiftTargetType_HeadFirst, true
	case "后排优先":
		return SkillSiftTargetType_TailFrist, true
	case "能量最多":
		return SkillSiftTargetType_MaxSp, true
	case "同队伍":
		return SkillSiftTargetType_SameTeam, true
	case "职业":
		return SkillSiftTargetType_Carrer, true
	case "同队伍职业筛选":
		return SkillSiftTargetType_SameTeamCarrer, true
	case "指定英雄":
		return SkillSiftTargetType_CertainHero, true
	case "固定英雄优先":
		return SkillSiftTargetType_FixedHeroPriority, true
	case "血量百分比最少":
		return SkillSiftTargetType_MinHpPrecent, true
	case "最高战力":
		return SkillSiftTargetType_MaxPower, true
	case "同队最高战力":
		return SkillSiftTargetType_SameTeamMaxPower, true
	case "同队站位靠前":
		return SkillSiftTargetType_SameTeamHeadFirst, true
	case "同队站位靠后":
		return SkillSiftTargetType_SameTeamTailFirst, true
	case "大明星前排优先":
		return SkillSiftTargetType_KofStarHeadFirst, true
	case "非Boss且血量最少":
		return SkillSiftTargetType_NotBossAndMinHpPrecent, true
	case "我及我前方":
		return SkillSiftTargetType_ForntOfMe, true
	case "我及我后方":
		return SkillSiftTargetType_BehindOfMe, true

	}
	return -1, false

}

func (e DamageType) String() string {

	if e == DamageType_All {
		return "All"
	}
	if e == DamageType_Physical {
		return "Physical"
	}
	if e == DamageType_Magic {
		return "Magic"
	}
	if e == DamageType_TrueDmg {
		return "TrueDmg"
	}
	if e == DamageType_Fire {
		return "Fire"
	}

	return ""

}

func GetDamageType(tp string) (DamageType, bool) {
	switch tp {
	case "All":
		return DamageType_All, true
	case "Physical":
		return DamageType_Physical, true
	case "Magic":
		return DamageType_Magic, true
	case "TrueDmg":
		return DamageType_TrueDmg, true
	case "Fire":
		return DamageType_Fire, true

	}
	return -1, false

}

func GetDamageType_ByAlias(tp string) (DamageType, bool) {
	switch tp {
	case "全种类":
		return DamageType_All, true
	case "物理伤害":
		return DamageType_Physical, true
	case "法术伤害":
		return DamageType_Magic, true
	case "真实伤害":
		return DamageType_TrueDmg, true
	case "火系伤害":
		return DamageType_Fire, true

	}
	return -1, false

}

func (e DamageCalculateType) String() string {

	if e == DamageCalculateType_HpPercent {
		return "HpPercent"
	}
	if e == DamageCalculateType_MaxPercent {
		return "MaxPercent"
	}
	if e == DamageCalculateType_AttrPercent {
		return "AttrPercent"
	}

	return ""

}

func GetDamageCalculateType(tp string) (DamageCalculateType, bool) {
	switch tp {
	case "HpPercent":
		return DamageCalculateType_HpPercent, true
	case "MaxPercent":
		return DamageCalculateType_MaxPercent, true
	case "AttrPercent":
		return DamageCalculateType_AttrPercent, true

	}
	return -1, false

}

func GetDamageCalculateType_ByAlias(tp string) (DamageCalculateType, bool) {
	switch tp {
	case "血量百分比":
		return DamageCalculateType_HpPercent, true
	case "血量上限百分比":
		return DamageCalculateType_MaxPercent, true
	case "属性百分比":
		return DamageCalculateType_AttrPercent, true

	}
	return -1, false

}

func (e SubDamageType) String() string {

	if e == SubDamageType_Level {
		return "Level"
	}
	if e == SubDamageType_Fixed {
		return "Fixed"
	}

	return ""

}

func GetSubDamageType(tp string) (SubDamageType, bool) {
	switch tp {
	case "Level":
		return SubDamageType_Level, true
	case "Fixed":
		return SubDamageType_Fixed, true

	}
	return -1, false

}

func GetSubDamageType_ByAlias(tp string) (SubDamageType, bool) {
	switch tp {
	case "等级":
		return SubDamageType_Level, true
	case "固定":
		return SubDamageType_Fixed, true

	}
	return -1, false

}

func (e SkillLockType) String() string {

	if e == SkillLockType_Target {
		return "Target"
	}
	if e == SkillLockType_Postion {
		return "Postion"
	}

	return ""

}

func GetSkillLockType(tp string) (SkillLockType, bool) {
	switch tp {
	case "Target":
		return SkillLockType_Target, true
	case "Postion":
		return SkillLockType_Postion, true

	}
	return -1, false

}

func GetSkillLockType_ByAlias(tp string) (SkillLockType, bool) {
	switch tp {
	case "锁目标":
		return SkillLockType_Target, true
	case "锁位置":
		return SkillLockType_Postion, true

	}
	return -1, false

}

func (e SkillQuaility) String() string {

	if e == SkillQuaility_Collect {
		return "Collect"
	}
	if e == SkillQuaility_Normal {
		return "Normal"
	}
	if e == SkillQuaility_Excellence {
		return "Excellence"
	}
	if e == SkillQuaility_Rare {
		return "Rare"
	}
	if e == SkillQuaility_Epic {
		return "Epic"
	}
	if e == SkillQuaility_Legend {
		return "Legend"
	}

	return ""

}

func GetSkillQuaility(tp string) (SkillQuaility, bool) {
	switch tp {
	case "Collect":
		return SkillQuaility_Collect, true
	case "Normal":
		return SkillQuaility_Normal, true
	case "Excellence":
		return SkillQuaility_Excellence, true
	case "Rare":
		return SkillQuaility_Rare, true
	case "Epic":
		return SkillQuaility_Epic, true
	case "Legend":
		return SkillQuaility_Legend, true

	}
	return -1, false

}

func GetSkillQuaility_ByAlias(tp string) (SkillQuaility, bool) {
	switch tp {
	case "采集":
		return SkillQuaility_Collect, true
	case "普通":
		return SkillQuaility_Normal, true
	case "卓越":
		return SkillQuaility_Excellence, true
	case "稀有":
		return SkillQuaility_Rare, true
	case "史诗":
		return SkillQuaility_Epic, true
	case "传说":
		return SkillQuaility_Legend, true

	}
	return -1, false

}

func (e BuffType) String() string {

	if e == BuffType_Stun {
		return "Stun"
	}
	if e == BuffType_Freeze {
		return "Freeze"
	}
	if e == BuffType_Silent {
		return "Silent"
	}
	if e == BuffType_Blood {
		return "Blood"
	}
	if e == BuffType_Poison {
		return "Poison"
	}
	if e == BuffType_Burn {
		return "Burn"
	}
	if e == BuffType_Positive {
		return "Positive"
	}
	if e == BuffType_Negative {
		return "Negative"
	}
	if e == BuffType_Control {
		return "Control"
	}
	if e == BuffType_Fixation {
		return "Fixation"
	}
	if e == BuffType_Decelerate {
		return "Decelerate"
	}
	if e == BuffType_Harm {
		return "Harm"
	}
	if e == BuffType_Heal {
		return "Heal"
	}
	if e == BuffType_Special {
		return "Special"
	}
	if e == BuffType_Frozen {
		return "Frozen"
	}
	if e == BuffType_Lianna {
		return "Lianna"
	}
	if e == BuffType_Shichuizhenwu {
		return "Shichuizhenwu"
	}
	if e == BuffType_TimePasue {
		return "TimePasue"
	}
	if e == BuffType_Banqiliangxiaozhaozhuangtai {
		return "Banqiliangxiaozhaozhuangtai"
	}
	if e == BuffType_Oumijia1 {
		return "Oumijia1"
	}
	if e == BuffType_Oumijia2 {
		return "Oumijia2"
	}
	if e == BuffType_AllDizzy {
		return "AllDizzy"
	}
	if e == BuffType_None {
		return "None"
	}

	return ""

}

func GetBuffType(tp string) (BuffType, bool) {
	switch tp {
	case "Stun":
		return BuffType_Stun, true
	case "Freeze":
		return BuffType_Freeze, true
	case "Silent":
		return BuffType_Silent, true
	case "Blood":
		return BuffType_Blood, true
	case "Poison":
		return BuffType_Poison, true
	case "Burn":
		return BuffType_Burn, true
	case "Positive":
		return BuffType_Positive, true
	case "Negative":
		return BuffType_Negative, true
	case "Control":
		return BuffType_Control, true
	case "Fixation":
		return BuffType_Fixation, true
	case "Decelerate":
		return BuffType_Decelerate, true
	case "Harm":
		return BuffType_Harm, true
	case "Heal":
		return BuffType_Heal, true
	case "Special":
		return BuffType_Special, true
	case "Frozen":
		return BuffType_Frozen, true
	case "Lianna":
		return BuffType_Lianna, true
	case "Shichuizhenwu":
		return BuffType_Shichuizhenwu, true
	case "TimePasue":
		return BuffType_TimePasue, true
	case "Banqiliangxiaozhaozhuangtai":
		return BuffType_Banqiliangxiaozhaozhuangtai, true
	case "Oumijia1":
		return BuffType_Oumijia1, true
	case "Oumijia2":
		return BuffType_Oumijia2, true
	case "AllDizzy":
		return BuffType_AllDizzy, true
	case "None":
		return BuffType_None, true

	}
	return -1, false

}

func GetBuffType_ByAlias(tp string) (BuffType, bool) {
	switch tp {
	case "眩晕":
		return BuffType_Stun, true
	case "冰冻":
		return BuffType_Freeze, true
	case "沉默":
		return BuffType_Silent, true
	case "流血":
		return BuffType_Blood, true
	case "中毒":
		return BuffType_Poison, true
	case "灼烧":
		return BuffType_Burn, true
	case "正面":
		return BuffType_Positive, true
	case "负面":
		return BuffType_Negative, true
	case "控制":
		return BuffType_Control, true
	case "定身":
		return BuffType_Fixation, true
	case "减速":
		return BuffType_Decelerate, true
	case "伤害":
		return BuffType_Harm, true
	case "治愈":
		return BuffType_Heal, true
	case "特殊":
		return BuffType_Special, true
	case "冻结":
		return BuffType_Frozen, true
	case "莉安娜吸血":
		return BuffType_Lianna, true
	case "矢吹真吾灼烧":
		return BuffType_Shichuizhenwu, true
	case "时停":
		return BuffType_TimePasue, true
	case "坂崎良小招状态":
		return BuffType_Banqiliangxiaozhaozhuangtai, true
	case "欧米茄充能":
		return BuffType_Oumijia1, true
	case "欧米茄充能协同":
		return BuffType_Oumijia2, true
	case "全部":
		return BuffType_AllDizzy, true
	case "无":
		return BuffType_None, true

	}
	return -1, false

}

func (e BuffGroupType) String() string {

	if e == BuffGroupType_Unique {
		return "Unique"
	}
	if e == BuffGroupType_Scattered {
		return "Scattered"
	}
	if e == BuffGroupType_Concentrated {
		return "Concentrated"
	}
	if e == BuffGroupType_Share {
		return "Share"
	}

	return ""

}

func GetBuffGroupType(tp string) (BuffGroupType, bool) {
	switch tp {
	case "Unique":
		return BuffGroupType_Unique, true
	case "Scattered":
		return BuffGroupType_Scattered, true
	case "Concentrated":
		return BuffGroupType_Concentrated, true
	case "Share":
		return BuffGroupType_Share, true

	}
	return -1, false

}

func GetBuffGroupType_ByAlias(tp string) (BuffGroupType, bool) {
	switch tp {
	case "唯一":
		return BuffGroupType_Unique, true
	case "分散":
		return BuffGroupType_Scattered, true
	case "集中":
		return BuffGroupType_Concentrated, true
	case "共享":
		return BuffGroupType_Share, true

	}
	return -1, false

}

func (e BuffFlagType) String() string {

	if e == BuffFlagType_Dizzy {
		return "Dizzy"
	}
	if e == BuffFlagType_Endure {
		return "Endure"
	}
	if e == BuffFlagType_Invincible {
		return "Invincible"
	}
	if e == BuffFlagType_BadlyHurt {
		return "BadlyHurt"
	}
	if e == BuffFlagType_Silent {
		return "Silent"
	}
	if e == BuffFlagType_Freeze {
		return "Freeze"
	}
	if e == BuffFlagType_Snare {
		return "Snare"
	}
	if e == BuffFlagType_Disarm {
		return "Disarm"
	}
	if e == BuffFlagType_Taunt {
		return "Taunt"
	}
	if e == BuffFlagType_Charm {
		return "Charm"
	}
	if e == BuffFlagType_Bokeh {
		return "Bokeh"
	}
	if e == BuffFlagType_Support {
		return "Support"
	}
	if e == BuffFlagType_ResistDeath {
		return "ResistDeath"
	}
	if e == BuffFlagType_HunterSign {
		return "HunterSign"
	}
	if e == BuffFlagType_Stealth {
		return "Stealth"
	}
	if e == BuffFlagType_Imprison {
		return "Imprison"
	}
	if e == BuffFlagType_LimtSp {
		return "LimtSp"
	}
	if e == BuffFlagType_Frozen {
		return "Frozen"
	}
	if e == BuffFlagType_Transfiguration {
		return "Transfiguration"
	}
	if e == BuffFlagType_DeTrans {
		return "DeTrans"
	}
	if e == BuffFlagType_ChangeNormalSkill {
		return "ChangeNormalSkill"
	}
	if e == BuffFlagType_DeChange {
		return "DeChange"
	}
	if e == BuffFlagType_ScaleChange {
		return "ScaleChange"
	}
	if e == BuffFlagType_ForbidenRotate {
		return "ForbidenRotate"
	}
	if e == BuffFlagType_ForbidenSpeedDown {
		return "ForbidenSpeedDown"
	}
	if e == BuffFlagType_BreakEndrue {
		return "BreakEndrue"
	}
	if e == BuffFlagType_ChargeImmunity {
		return "ChargeImmunity"
	}
	if e == BuffFlagType_ForbidenHeal {
		return "ForbidenHeal"
	}
	if e == BuffFlagType_KyoBreak {
		return "KyoBreak"
	}
	if e == BuffFlagType_NoDOT {
		return "NoDOT"
	}
	if e == BuffFlagType_NoDedcut {
		return "NoDedcut"
	}
	if e == BuffFlagType_AllTrueDamage {
		return "AllTrueDamage"
	}
	if e == BuffFlagType_IgnoreArmor {
		return "IgnoreArmor"
	}

	return ""

}

func GetBuffFlagType(tp string) (BuffFlagType, bool) {
	switch tp {
	case "Dizzy":
		return BuffFlagType_Dizzy, true
	case "Endure":
		return BuffFlagType_Endure, true
	case "Invincible":
		return BuffFlagType_Invincible, true
	case "BadlyHurt":
		return BuffFlagType_BadlyHurt, true
	case "Silent":
		return BuffFlagType_Silent, true
	case "Freeze":
		return BuffFlagType_Freeze, true
	case "Snare":
		return BuffFlagType_Snare, true
	case "Disarm":
		return BuffFlagType_Disarm, true
	case "Taunt":
		return BuffFlagType_Taunt, true
	case "Charm":
		return BuffFlagType_Charm, true
	case "Bokeh":
		return BuffFlagType_Bokeh, true
	case "Support":
		return BuffFlagType_Support, true
	case "ResistDeath":
		return BuffFlagType_ResistDeath, true
	case "HunterSign":
		return BuffFlagType_HunterSign, true
	case "Stealth":
		return BuffFlagType_Stealth, true
	case "Imprison":
		return BuffFlagType_Imprison, true
	case "LimtSp":
		return BuffFlagType_LimtSp, true
	case "Frozen":
		return BuffFlagType_Frozen, true
	case "Transfiguration":
		return BuffFlagType_Transfiguration, true
	case "DeTrans":
		return BuffFlagType_DeTrans, true
	case "ChangeNormalSkill":
		return BuffFlagType_ChangeNormalSkill, true
	case "DeChange":
		return BuffFlagType_DeChange, true
	case "ScaleChange":
		return BuffFlagType_ScaleChange, true
	case "ForbidenRotate":
		return BuffFlagType_ForbidenRotate, true
	case "ForbidenSpeedDown":
		return BuffFlagType_ForbidenSpeedDown, true
	case "BreakEndrue":
		return BuffFlagType_BreakEndrue, true
	case "ChargeImmunity":
		return BuffFlagType_ChargeImmunity, true
	case "ForbidenHeal":
		return BuffFlagType_ForbidenHeal, true
	case "KyoBreak":
		return BuffFlagType_KyoBreak, true
	case "NoDOT":
		return BuffFlagType_NoDOT, true
	case "NoDedcut":
		return BuffFlagType_NoDedcut, true
	case "AllTrueDamage":
		return BuffFlagType_AllTrueDamage, true
	case "IgnoreArmor":
		return BuffFlagType_IgnoreArmor, true

	}
	return -1, false

}

func GetBuffFlagType_ByAlias(tp string) (BuffFlagType, bool) {
	switch tp {
	case "控制":
		return BuffFlagType_Dizzy, true
	case "霸体":
		return BuffFlagType_Endure, true
	case "免伤":
		return BuffFlagType_Invincible, true
	case "重伤":
		return BuffFlagType_BadlyHurt, true
	case "沉默":
		return BuffFlagType_Silent, true
	case "冰冻":
		return BuffFlagType_Freeze, true
	case "定身":
		return BuffFlagType_Snare, true
	case "缴械":
		return BuffFlagType_Disarm, true
	case "嘲讽":
		return BuffFlagType_Taunt, true
	case "魅惑":
		return BuffFlagType_Charm, true
	case "虚化":
		return BuffFlagType_Bokeh, true
	case "援护":
		return BuffFlagType_Support, true
	case "抗死":
		return BuffFlagType_ResistDeath, true
	case "标记":
		return BuffFlagType_HunterSign, true
	case "隐身":
		return BuffFlagType_Stealth, true
	case "禁锢":
		return BuffFlagType_Imprison, true
	case "SP限制":
		return BuffFlagType_LimtSp, true
	case "冻结":
		return BuffFlagType_Frozen, true
	case "变身":
		return BuffFlagType_Transfiguration, true
	case "抗变身":
		return BuffFlagType_DeTrans, true
	case "变更普攻":
		return BuffFlagType_ChangeNormalSkill, true
	case "禁普攻改变":
		return BuffFlagType_DeChange, true
	case "放缩":
		return BuffFlagType_ScaleChange, true
	case "禁转向":
		return BuffFlagType_ForbidenRotate, true
	case "禁减速":
		return BuffFlagType_ForbidenSpeedDown, true
	case "破霸体":
		return BuffFlagType_BreakEndrue, true
	case "冲锋免疫":
		return BuffFlagType_ChargeImmunity, true
	case "禁疗":
		return BuffFlagType_ForbidenHeal, true
	case "草薙流识破":
		return BuffFlagType_KyoBreak, true
	case "DOT免疫":
		return BuffFlagType_NoDOT, true
	case "Dedcut免疫":
		return BuffFlagType_NoDedcut, true
	case "绝对真伤":
		return BuffFlagType_AllTrueDamage, true
	case "护盾失效":
		return BuffFlagType_IgnoreArmor, true

	}
	return -1, false

}

func (e ShockType) String() string {

	if e == ShockType_Rect {
		return "Rect"
	}
	if e == ShockType_Sector {
		return "Sector"
	}
	if e == ShockType_Annular {
		return "Annular"
	}

	return ""

}

func GetShockType(tp string) (ShockType, bool) {
	switch tp {
	case "Rect":
		return ShockType_Rect, true
	case "Sector":
		return ShockType_Sector, true
	case "Annular":
		return ShockType_Annular, true

	}
	return -1, false

}

func GetShockType_ByAlias(tp string) (ShockType, bool) {
	switch tp {
	case "矩形":
		return ShockType_Rect, true
	case "扇形":
		return ShockType_Sector, true
	case "环形":
		return ShockType_Annular, true

	}
	return -1, false

}

func (e ShieldBrokenType) String() string {

	if e == ShieldBrokenType_EffectTimes {
		return "EffectTimes"
	}
	if e == ShieldBrokenType_Time {
		return "Time"
	}
	if e == ShieldBrokenType_SheildValue {
		return "SheildValue"
	}
	if e == ShieldBrokenType_Puppet {
		return "Puppet"
	}

	return ""

}

func GetShieldBrokenType(tp string) (ShieldBrokenType, bool) {
	switch tp {
	case "EffectTimes":
		return ShieldBrokenType_EffectTimes, true
	case "Time":
		return ShieldBrokenType_Time, true
	case "SheildValue":
		return ShieldBrokenType_SheildValue, true
	case "Puppet":
		return ShieldBrokenType_Puppet, true

	}
	return -1, false

}

func GetShieldBrokenType_ByAlias(tp string) (ShieldBrokenType, bool) {
	switch tp {
	case "次数":
		return ShieldBrokenType_EffectTimes, true
	case "时间":
		return ShieldBrokenType_Time, true
	case "吸收值":
		return ShieldBrokenType_SheildValue, true
	case "娃娃盾":
		return ShieldBrokenType_Puppet, true

	}
	return -1, false

}

func (e EffectPlayType) String() string {

	if e == EffectPlayType_None {
		return "None"
	}
	if e == EffectPlayType_Sprite {
		return "Sprite"
	}
	if e == EffectPlayType_FramSeq {
		return "FramSeq"
	}
	if e == EffectPlayType_Skeleton {
		return "Skeleton"
	}

	return ""

}

func GetEffectPlayType(tp string) (EffectPlayType, bool) {
	switch tp {
	case "None":
		return EffectPlayType_None, true
	case "Sprite":
		return EffectPlayType_Sprite, true
	case "FramSeq":
		return EffectPlayType_FramSeq, true
	case "Skeleton":
		return EffectPlayType_Skeleton, true

	}
	return -1, false

}

func GetEffectPlayType_ByAlias(tp string) (EffectPlayType, bool) {
	switch tp {
	case "无":
		return EffectPlayType_None, true
	case "图片":
		return EffectPlayType_Sprite, true
	case "序列帧":
		return EffectPlayType_FramSeq, true
	case "动画":
		return EffectPlayType_Skeleton, true

	}
	return -1, false

}

func (e EffectOwnerType) String() string {

	if e == EffectOwnerType_Skill {
		return "Skill"
	}
	if e == EffectOwnerType_Buff {
		return "Buff"
	}

	return ""

}

func GetEffectOwnerType(tp string) (EffectOwnerType, bool) {
	switch tp {
	case "Skill":
		return EffectOwnerType_Skill, true
	case "Buff":
		return EffectOwnerType_Buff, true

	}
	return -1, false

}

func GetEffectOwnerType_ByAlias(tp string) (EffectOwnerType, bool) {
	switch tp {
	case "技能":
		return EffectOwnerType_Skill, true
	case "状态":
		return EffectOwnerType_Buff, true

	}
	return -1, false

}

func (e BulletMoveType) String() string {

	if e == BulletMoveType_Line {
		return "Line"
	}
	if e == BulletMoveType_Parabola {
		return "Parabola"
	}
	if e == BulletMoveType_ReverseParbola {
		return "ReverseParbola"
	}
	if e == BulletMoveType_Ballistics {
		return "Ballistics"
	}
	if e == BulletMoveType_FaceBallistics {
		return "FaceBallistics"
	}
	if e == BulletMoveType_GroudBallistics {
		return "GroudBallistics"
	}

	return ""

}

func GetBulletMoveType(tp string) (BulletMoveType, bool) {
	switch tp {
	case "Line":
		return BulletMoveType_Line, true
	case "Parabola":
		return BulletMoveType_Parabola, true
	case "ReverseParbola":
		return BulletMoveType_ReverseParbola, true
	case "Ballistics":
		return BulletMoveType_Ballistics, true
	case "FaceBallistics":
		return BulletMoveType_FaceBallistics, true
	case "GroudBallistics":
		return BulletMoveType_GroudBallistics, true

	}
	return -1, false

}

func GetBulletMoveType_ByAlias(tp string) (BulletMoveType, bool) {
	switch tp {
	case "直线":
		return BulletMoveType_Line, true
	case "抛物线":
		return BulletMoveType_Parabola, true
	case "反抛物线":
		return BulletMoveType_ReverseParbola, true
	case "弹道":
		return BulletMoveType_Ballistics, true
	case "面向弹道":
		return BulletMoveType_FaceBallistics, true
	case "贴地弹道":
		return BulletMoveType_GroudBallistics, true

	}
	return -1, false

}

func (e SkillMoveTargetType) String() string {

	if e == SkillMoveTargetType_TargetPos {
		return "TargetPos"
	}
	if e == SkillMoveTargetType_SelfPos {
		return "SelfPos"
	}
	if e == SkillMoveTargetType_Escape {
		return "Escape"
	}
	if e == SkillMoveTargetType_AbsoluteAngle {
		return "AbsoluteAngle"
	}

	return ""

}

func GetSkillMoveTargetType(tp string) (SkillMoveTargetType, bool) {
	switch tp {
	case "TargetPos":
		return SkillMoveTargetType_TargetPos, true
	case "SelfPos":
		return SkillMoveTargetType_SelfPos, true
	case "Escape":
		return SkillMoveTargetType_Escape, true
	case "AbsoluteAngle":
		return SkillMoveTargetType_AbsoluteAngle, true

	}
	return -1, false

}

func GetSkillMoveTargetType_ByAlias(tp string) (SkillMoveTargetType, bool) {
	switch tp {
	case "目标位置":
		return SkillMoveTargetType_TargetPos, true
	case "自身位置":
		return SkillMoveTargetType_SelfPos, true
	case "逃逸位置":
		return SkillMoveTargetType_Escape, true
	case "固定角度":
		return SkillMoveTargetType_AbsoluteAngle, true

	}
	return -1, false

}

func (e SkillMoveDirectionType) String() string {

	if e == SkillMoveDirectionType_TargetForward {
		return "TargetForward"
	}
	if e == SkillMoveDirectionType_SelfForward {
		return "SelfForward"
	}

	return ""

}

func GetSkillMoveDirectionType(tp string) (SkillMoveDirectionType, bool) {
	switch tp {
	case "TargetForward":
		return SkillMoveDirectionType_TargetForward, true
	case "SelfForward":
		return SkillMoveDirectionType_SelfForward, true

	}
	return -1, false

}

func GetSkillMoveDirectionType_ByAlias(tp string) (SkillMoveDirectionType, bool) {
	switch tp {
	case "目标正向":
		return SkillMoveDirectionType_TargetForward, true
	case "自身正向":
		return SkillMoveDirectionType_SelfForward, true

	}
	return -1, false

}

func (e EffectRoleType) String() string {

	if e == EffectRoleType_SkillTarget {
		return "SkillTarget"
	}
	if e == EffectRoleType_SkillPreform {
		return "SkillPreform"
	}
	if e == EffectRoleType_EnemyAll {
		return "EnemyAll"
	}
	if e == EffectRoleType_EnemyFront {
		return "EnemyFront"
	}
	if e == EffectRoleType_EnemyBack {
		return "EnemyBack"
	}
	if e == EffectRoleType_SelfAll {
		return "SelfAll"
	}
	if e == EffectRoleType_SelfFront {
		return "SelfFront"
	}
	if e == EffectRoleType_SelfBack {
		return "SelfBack"
	}
	if e == EffectRoleType_All {
		return "All"
	}
	if e == EffectRoleType_Self {
		return "Self"
	}
	if e == EffectRoleType_Wall {
		return "Wall"
	}
	if e == EffectRoleType_BuffOwner {
		return "BuffOwner"
	}
	if e == EffectRoleType_CheckTarget {
		return "CheckTarget"
	}
	if e == EffectRoleType_CurAttack {
		return "CurAttack"
	}
	if e == EffectRoleType_MySummon {
		return "MySummon"
	}
	if e == EffectRoleType_Master {
		return "Master"
	}
	if e == EffectRoleType_Attacker {
		return "Attacker"
	}
	if e == EffectRoleType_Support {
		return "Support"
	}
	if e == EffectRoleType_AllPlayers {
		return "AllPlayers"
	}
	if e == EffectRoleType_SelfExcudeSelf {
		return "SelfExcudeSelf"
	}
	if e == EffectRoleType_AllExcludeSelf {
		return "AllExcludeSelf"
	}
	if e == EffectRoleType_EnemyMonster {
		return "EnemyMonster"
	}
	if e == EffectRoleType_SelfMonster {
		return "SelfMonster"
	}
	if e == EffectRoleType_EnemyAllWithWall {
		return "EnemyAllWithWall"
	}
	if e == EffectRoleType_SelfWall {
		return "SelfWall"
	}
	if e == EffectRoleType_SelfMonsterWithoutHunt {
		return "SelfMonsterWithoutHunt"
	}
	if e == EffectRoleType_GlobalMan {
		return "GlobalMan"
	}
	if e == EffectRoleType_SelfGlobalMan {
		return "SelfGlobalMan"
	}
	if e == EffectRoleType_EnemyGlobalMan {
		return "EnemyGlobalMan"
	}
	if e == EffectRoleType_LockTarget {
		return "LockTarget"
	}

	return ""

}

func GetEffectRoleType(tp string) (EffectRoleType, bool) {
	switch tp {
	case "SkillTarget":
		return EffectRoleType_SkillTarget, true
	case "SkillPreform":
		return EffectRoleType_SkillPreform, true
	case "EnemyAll":
		return EffectRoleType_EnemyAll, true
	case "EnemyFront":
		return EffectRoleType_EnemyFront, true
	case "EnemyBack":
		return EffectRoleType_EnemyBack, true
	case "SelfAll":
		return EffectRoleType_SelfAll, true
	case "SelfFront":
		return EffectRoleType_SelfFront, true
	case "SelfBack":
		return EffectRoleType_SelfBack, true
	case "All":
		return EffectRoleType_All, true
	case "Self":
		return EffectRoleType_Self, true
	case "Wall":
		return EffectRoleType_Wall, true
	case "BuffOwner":
		return EffectRoleType_BuffOwner, true
	case "CheckTarget":
		return EffectRoleType_CheckTarget, true
	case "CurAttack":
		return EffectRoleType_CurAttack, true
	case "MySummon":
		return EffectRoleType_MySummon, true
	case "Master":
		return EffectRoleType_Master, true
	case "Attacker":
		return EffectRoleType_Attacker, true
	case "Support":
		return EffectRoleType_Support, true
	case "AllPlayers":
		return EffectRoleType_AllPlayers, true
	case "SelfExcudeSelf":
		return EffectRoleType_SelfExcudeSelf, true
	case "AllExcludeSelf":
		return EffectRoleType_AllExcludeSelf, true
	case "EnemyMonster":
		return EffectRoleType_EnemyMonster, true
	case "SelfMonster":
		return EffectRoleType_SelfMonster, true
	case "EnemyAllWithWall":
		return EffectRoleType_EnemyAllWithWall, true
	case "SelfWall":
		return EffectRoleType_SelfWall, true
	case "SelfMonsterWithoutHunt":
		return EffectRoleType_SelfMonsterWithoutHunt, true
	case "GlobalMan":
		return EffectRoleType_GlobalMan, true
	case "SelfGlobalMan":
		return EffectRoleType_SelfGlobalMan, true
	case "EnemyGlobalMan":
		return EffectRoleType_EnemyGlobalMan, true
	case "LockTarget":
		return EffectRoleType_LockTarget, true

	}
	return -1, false

}

func GetEffectRoleType_ByAlias(tp string) (EffectRoleType, bool) {
	switch tp {
	case "技能目标":
		return EffectRoleType_SkillTarget, true
	case "安装者":
		return EffectRoleType_SkillPreform, true
	case "敌方全体":
		return EffectRoleType_EnemyAll, true
	case "敌方前排":
		return EffectRoleType_EnemyFront, true
	case "敌方后排":
		return EffectRoleType_EnemyBack, true
	case "己方全体":
		return EffectRoleType_SelfAll, true
	case "己方前排":
		return EffectRoleType_SelfFront, true
	case "己方后排":
		return EffectRoleType_SelfBack, true
	case "场上全部":
		return EffectRoleType_All, true
	case "自身":
		return EffectRoleType_Self, true
	case "墙体":
		return EffectRoleType_Wall, true
	case "拥有者":
		return EffectRoleType_BuffOwner, true
	case "检测对象":
		return EffectRoleType_CheckTarget, true
	case "当前目标":
		return EffectRoleType_CurAttack, true
	case "自身召唤物":
		return EffectRoleType_MySummon, true
	case "召唤者":
		return EffectRoleType_Master, true
	case "攻击者":
		return EffectRoleType_Attacker, true
	case "援护对象":
		return EffectRoleType_Support, true
	case "玩家全体":
		return EffectRoleType_AllPlayers, true
	case "己方除自己":
		return EffectRoleType_SelfExcudeSelf, true
	case "全部除自己":
		return EffectRoleType_AllExcludeSelf, true
	case "敌方怪物":
		return EffectRoleType_EnemyMonster, true
	case "己方怪物":
		return EffectRoleType_SelfMonster, true
	case "敌方全体以及城墙":
		return EffectRoleType_EnemyAllWithWall, true
	case "己方墙体":
		return EffectRoleType_SelfWall, true
	case "己方怪物除狩猎":
		return EffectRoleType_SelfMonsterWithoutHunt, true
	case "机器人":
		return EffectRoleType_GlobalMan, true
	case "己方机器人":
		return EffectRoleType_SelfGlobalMan, true
	case "敌方机器人":
		return EffectRoleType_EnemyGlobalMan, true
	case "当前锁定目标":
		return EffectRoleType_LockTarget, true

	}
	return -1, false

}

func (e EffectDirType) String() string {

	if e == EffectDirType_None {
		return "None"
	}
	if e == EffectDirType_Target {
		return "Target"
	}
	if e == EffectDirType_Self {
		return "Self"
	}
	if e == EffectDirType_SelfAndFilp {
		return "SelfAndFilp"
	}
	if e == EffectDirType_LogicDir {
		return "LogicDir"
	}

	return ""

}

func GetEffectDirType(tp string) (EffectDirType, bool) {
	switch tp {
	case "None":
		return EffectDirType_None, true
	case "Target":
		return EffectDirType_Target, true
	case "Self":
		return EffectDirType_Self, true
	case "SelfAndFilp":
		return EffectDirType_SelfAndFilp, true
	case "LogicDir":
		return EffectDirType_LogicDir, true

	}
	return -1, false

}

func GetEffectDirType_ByAlias(tp string) (EffectDirType, bool) {
	switch tp {
	case "无朝向":
		return EffectDirType_None, true
	case "目标朝向":
		return EffectDirType_Target, true
	case "自身朝向":
		return EffectDirType_Self, true
	case "自身朝向与翻转":
		return EffectDirType_SelfAndFilp, true
	case "逻辑朝向":
		return EffectDirType_LogicDir, true

	}
	return -1, false

}

func (e BuffFixedDamageType) String() string {

	if e == BuffFixedDamageType_Sub {
		return "Sub"
	}
	if e == BuffFixedDamageType_Replace {
		return "Replace"
	}

	return ""

}

func GetBuffFixedDamageType(tp string) (BuffFixedDamageType, bool) {
	switch tp {
	case "Sub":
		return BuffFixedDamageType_Sub, true
	case "Replace":
		return BuffFixedDamageType_Replace, true

	}
	return -1, false

}

func GetBuffFixedDamageType_ByAlias(tp string) (BuffFixedDamageType, bool) {
	switch tp {
	case "额外附加":
		return BuffFixedDamageType_Sub, true
	case "变为":
		return BuffFixedDamageType_Replace, true

	}
	return -1, false

}

func (e SummonCampType) String() string {

	if e == SummonCampType_SameWithMaster {
		return "SameWithMaster"
	}
	if e == SummonCampType_EnemyWithMaster {
		return "EnemyWithMaster"
	}
	if e == SummonCampType_ThirdCamp {
		return "ThirdCamp"
	}
	if e == SummonCampType_None {
		return "None"
	}

	return ""

}

func GetSummonCampType(tp string) (SummonCampType, bool) {
	switch tp {
	case "SameWithMaster":
		return SummonCampType_SameWithMaster, true
	case "EnemyWithMaster":
		return SummonCampType_EnemyWithMaster, true
	case "ThirdCamp":
		return SummonCampType_ThirdCamp, true
	case "None":
		return SummonCampType_None, true

	}
	return -1, false

}

func GetSummonCampType_ByAlias(tp string) (SummonCampType, bool) {
	switch tp {
	case "同主人":
		return SummonCampType_SameWithMaster, true
	case "与主人敌对":
		return SummonCampType_EnemyWithMaster, true
	case "第三方阵营":
		return SummonCampType_ThirdCamp, true
	case "无阵营":
		return SummonCampType_None, true

	}
	return -1, false

}

func (e SummonMoveType) String() string {

	if e == SummonMoveType_Follow {
		return "Follow"
	}
	if e == SummonMoveType_Auto {
		return "Auto"
	}
	if e == SummonMoveType_None {
		return "None"
	}

	return ""

}

func GetSummonMoveType(tp string) (SummonMoveType, bool) {
	switch tp {
	case "Follow":
		return SummonMoveType_Follow, true
	case "Auto":
		return SummonMoveType_Auto, true
	case "None":
		return SummonMoveType_None, true

	}
	return -1, false

}

func GetSummonMoveType_ByAlias(tp string) (SummonMoveType, bool) {
	switch tp {
	case "跟随":
		return SummonMoveType_Follow, true
	case "自动":
		return SummonMoveType_Auto, true
	case "不动":
		return SummonMoveType_None, true

	}
	return -1, false

}

func (e ReplaceTargetType) String() string {

	if e == ReplaceTargetType_Circle {
		return "Circle"
	}
	if e == ReplaceTargetType_Sector {
		return "Sector"
	}
	if e == ReplaceTargetType_None {
		return "None"
	}

	return ""

}

func GetReplaceTargetType(tp string) (ReplaceTargetType, bool) {
	switch tp {
	case "Circle":
		return ReplaceTargetType_Circle, true
	case "Sector":
		return ReplaceTargetType_Sector, true
	case "None":
		return ReplaceTargetType_None, true

	}
	return -1, false

}

func GetReplaceTargetType_ByAlias(tp string) (ReplaceTargetType, bool) {
	switch tp {
	case "圆形范围":
		return ReplaceTargetType_Circle, true
	case "扇形范围":
		return ReplaceTargetType_Sector, true
	case "不替换":
		return ReplaceTargetType_None, true

	}
	return -1, false

}

func (e EffectShapeType) String() string {

	if e == EffectShapeType_Circle {
		return "Circle"
	}
	if e == EffectShapeType_Rect {
		return "Rect"
	}
	if e == EffectShapeType_Sector {
		return "Sector"
	}

	return ""

}

func GetEffectShapeType(tp string) (EffectShapeType, bool) {
	switch tp {
	case "Circle":
		return EffectShapeType_Circle, true
	case "Rect":
		return EffectShapeType_Rect, true
	case "Sector":
		return EffectShapeType_Sector, true

	}
	return -1, false

}

func GetEffectShapeType_ByAlias(tp string) (EffectShapeType, bool) {
	switch tp {
	case "圆形":
		return EffectShapeType_Circle, true
	case "矩形":
		return EffectShapeType_Rect, true
	case "扇形":
		return EffectShapeType_Sector, true

	}
	return -1, false

}

func (e SummonShapCopyType) String() string {

	if e == SummonShapCopyType_Self {
		return "Self"
	}
	if e == SummonShapCopyType_Target {
		return "Target"
	}
	if e == SummonShapCopyType_None {
		return "None"
	}

	return ""

}

func GetSummonShapCopyType(tp string) (SummonShapCopyType, bool) {
	switch tp {
	case "Self":
		return SummonShapCopyType_Self, true
	case "Target":
		return SummonShapCopyType_Target, true
	case "None":
		return SummonShapCopyType_None, true

	}
	return -1, false

}

func GetSummonShapCopyType_ByAlias(tp string) (SummonShapCopyType, bool) {
	switch tp {
	case "自己":
		return SummonShapCopyType_Self, true
	case "目标":
		return SummonShapCopyType_Target, true
	case "不复制":
		return SummonShapCopyType_None, true

	}
	return -1, false

}

func (e EffectLayerType) String() string {

	if e == EffectLayerType_UP {
		return "UP"
	}
	if e == EffectLayerType_Down {
		return "Down"
	}
	if e == EffectLayerType_UI {
		return "UI"
	}
	if e == EffectLayerType_Black {
		return "Black"
	}

	return ""

}

func GetEffectLayerType(tp string) (EffectLayerType, bool) {
	switch tp {
	case "UP":
		return EffectLayerType_UP, true
	case "Down":
		return EffectLayerType_Down, true
	case "UI":
		return EffectLayerType_UI, true
	case "Black":
		return EffectLayerType_Black, true

	}
	return -1, false

}

func GetEffectLayerType_ByAlias(tp string) (EffectLayerType, bool) {
	switch tp {
	case "上特效层":
		return EffectLayerType_UP, true
	case "下特效层":
		return EffectLayerType_Down, true
	case "UI层":
		return EffectLayerType_UI, true
	case "黑遮罩":
		return EffectLayerType_Black, true

	}
	return -1, false

}

func (e EffectPosType) String() string {

	if e == EffectPosType_Role {
		return "Role"
	}
	if e == EffectPosType_HitPoint {
		return "HitPoint"
	}
	if e == EffectPosType_HpPos {
		return "HpPos"
	}
	if e == EffectPosType_FirePoint {
		return "FirePoint"
	}

	return ""

}

func GetEffectPosType(tp string) (EffectPosType, bool) {
	switch tp {
	case "Role":
		return EffectPosType_Role, true
	case "HitPoint":
		return EffectPosType_HitPoint, true
	case "HpPos":
		return EffectPosType_HpPos, true
	case "FirePoint":
		return EffectPosType_FirePoint, true

	}
	return -1, false

}

func GetEffectPosType_ByAlias(tp string) (EffectPosType, bool) {
	switch tp {
	case "角色":
		return EffectPosType_Role, true
	case "受击点":
		return EffectPosType_HitPoint, true
	case "血条":
		return EffectPosType_HpPos, true
	case "开火点":
		return EffectPosType_FirePoint, true

	}
	return -1, false

}

func (e BattleRoleCamp) String() string {

	if e == BattleRoleCamp_Player {
		return "Player"
	}
	if e == BattleRoleCamp_Enemy {
		return "Enemy"
	}
	if e == BattleRoleCamp_Third {
		return "Third"
	}
	if e == BattleRoleCamp_NPC {
		return "NPC"
	}
	if e == BattleRoleCamp_FriendNpc {
		return "FriendNpc"
	}
	if e == BattleRoleCamp_Friend {
		return "Friend"
	}
	if e == BattleRoleCamp_FriendEnemy {
		return "FriendEnemy"
	}
	if e == BattleRoleCamp_SelfHero {
		return "SelfHero"
	}
	if e == BattleRoleCamp_EnemyHero {
		return "EnemyHero"
	}
	if e == BattleRoleCamp_Max {
		return "Max"
	}

	return ""

}

func GetBattleRoleCamp(tp string) (BattleRoleCamp, bool) {
	switch tp {
	case "Player":
		return BattleRoleCamp_Player, true
	case "Enemy":
		return BattleRoleCamp_Enemy, true
	case "Third":
		return BattleRoleCamp_Third, true
	case "NPC":
		return BattleRoleCamp_NPC, true
	case "FriendNpc":
		return BattleRoleCamp_FriendNpc, true
	case "Friend":
		return BattleRoleCamp_Friend, true
	case "FriendEnemy":
		return BattleRoleCamp_FriendEnemy, true
	case "SelfHero":
		return BattleRoleCamp_SelfHero, true
	case "EnemyHero":
		return BattleRoleCamp_EnemyHero, true
	case "Max":
		return BattleRoleCamp_Max, true

	}
	return -1, false

}

func GetBattleRoleCamp_ByAlias(tp string) (BattleRoleCamp, bool) {
	switch tp {
	case "玩家":
		return BattleRoleCamp_Player, true
	case "敌方":
		return BattleRoleCamp_Enemy, true
	case "第三方":
		return BattleRoleCamp_Third, true
	case "第三方Npc":
		return BattleRoleCamp_NPC, true
	case "友方NPC":
		return BattleRoleCamp_FriendNpc, true
	case "友方":
		return BattleRoleCamp_Friend, true
	case "友方敌人":
		return BattleRoleCamp_FriendEnemy, true
	case "我方英雄":
		return BattleRoleCamp_SelfHero, true
	case "敌方英雄":
		return BattleRoleCamp_EnemyHero, true
	case "":
		return BattleRoleCamp_Max, true

	}
	return -1, false

}

func (e SkillIconType) String() string {

	if e == SkillIconType_Healing {
		return "Healing"
	}
	if e == SkillIconType_Negative {
		return "Negative"
	}
	if e == SkillIconType_Ultimate {
		return "Ultimate"
	}
	if e == SkillIconType_Normal {
		return "Normal"
	}
	if e == SkillIconType_Positive {
		return "Positive"
	}
	if e == SkillIconType_Fetters {
		return "Fetters"
	}

	return ""

}

func GetSkillIconType(tp string) (SkillIconType, bool) {
	switch tp {
	case "Healing":
		return SkillIconType_Healing, true
	case "Negative":
		return SkillIconType_Negative, true
	case "Ultimate":
		return SkillIconType_Ultimate, true
	case "Normal":
		return SkillIconType_Normal, true
	case "Positive":
		return SkillIconType_Positive, true
	case "Fetters":
		return SkillIconType_Fetters, true

	}
	return -1, false

}

func GetSkillIconType_ByAlias(tp string) (SkillIconType, bool) {
	switch tp {
	case "治":
		return SkillIconType_Healing, true
	case "被":
		return SkillIconType_Negative, true
	case "奥":
		return SkillIconType_Ultimate, true
	case "普":
		return SkillIconType_Normal, true
	case "主":
		return SkillIconType_Positive, true
	case "羁":
		return SkillIconType_Fetters, true

	}
	return -1, false

}

func (e ConnectionType) String() string {

	if e == ConnectionType_Flash {
		return "Flash"
	}
	if e == ConnectionType_Move {
		return "Move"
	}
	if e == ConnectionType_Immediate {
		return "Immediate"
	}

	return ""

}

func GetConnectionType(tp string) (ConnectionType, bool) {
	switch tp {
	case "Flash":
		return ConnectionType_Flash, true
	case "Move":
		return ConnectionType_Move, true
	case "Immediate":
		return ConnectionType_Immediate, true

	}
	return -1, false

}

func GetConnectionType_ByAlias(tp string) (ConnectionType, bool) {
	switch tp {
	case "闪现连携":
		return ConnectionType_Flash, true
	case "移动连携":
		return ConnectionType_Move, true
	case "瞬时连携":
		return ConnectionType_Immediate, true

	}
	return -1, false

}

func (e HitStunOverlay) String() string {

	if e == HitStunOverlay_Cannot {
		return "Cannot"
	}
	if e == HitStunOverlay_Replace {
		return "Replace"
	}
	if e == HitStunOverlay_Overlay {
		return "Overlay"
	}
	if e == HitStunOverlay_SameOverlay {
		return "SameOverlay"
	}
	if e == HitStunOverlay_SameReplace {
		return "SameReplace"
	}

	return ""

}

func GetHitStunOverlay(tp string) (HitStunOverlay, bool) {
	switch tp {
	case "Cannot":
		return HitStunOverlay_Cannot, true
	case "Replace":
		return HitStunOverlay_Replace, true
	case "Overlay":
		return HitStunOverlay_Overlay, true
	case "SameOverlay":
		return HitStunOverlay_SameOverlay, true
	case "SameReplace":
		return HitStunOverlay_SameReplace, true

	}
	return -1, false

}

func GetHitStunOverlay_ByAlias(tp string) (HitStunOverlay, bool) {
	switch tp {
	case "不可操作":
		return HitStunOverlay_Cannot, true
	case "替换":
		return HitStunOverlay_Replace, true
	case "可叠加":
		return HitStunOverlay_Overlay, true
	case "同种可叠加":
		return HitStunOverlay_SameOverlay, true
	case "同种可替换":
		return HitStunOverlay_SameReplace, true

	}
	return -1, false

}

func (e HitStunType) String() string {

	if e == HitStunType_HitBack {
		return "HitBack"
	}
	if e == HitStunType_LockedUp {
		return "LockedUp"
	}
	if e == HitStunType_HitDown {
		return "HitDown"
	}
	if e == HitStunType_SkillHitBack {
		return "SkillHitBack"
	}

	return ""

}

func GetHitStunType(tp string) (HitStunType, bool) {
	switch tp {
	case "HitBack":
		return HitStunType_HitBack, true
	case "LockedUp":
		return HitStunType_LockedUp, true
	case "HitDown":
		return HitStunType_HitDown, true
	case "SkillHitBack":
		return HitStunType_SkillHitBack, true

	}
	return -1, false

}

func GetHitStunType_ByAlias(tp string) (HitStunType, bool) {
	switch tp {
	case "击退":
		return HitStunType_HitBack, true
	case "击飞":
		return HitStunType_LockedUp, true
	case "击倒":
		return HitStunType_HitDown, true
	case "技能特殊击退":
		return HitStunType_SkillHitBack, true

	}
	return -1, false

}

func (e SpecialStun) String() string {

	if e == SpecialStun_Guard {
		return "Guard"
	}
	if e == SpecialStun_Hurt {
		return "Hurt"
	}

	return ""

}

func GetSpecialStun(tp string) (SpecialStun, bool) {
	switch tp {
	case "Guard":
		return SpecialStun_Guard, true
	case "Hurt":
		return SpecialStun_Hurt, true

	}
	return -1, false

}

func GetSpecialStun_ByAlias(tp string) (SpecialStun, bool) {
	switch tp {
	case "防御硬直":
		return SpecialStun_Guard, true
	case "受击硬直":
		return SpecialStun_Hurt, true

	}
	return -1, false

}

func (e JumpBackType) String() string {

	if e == JumpBackType_None {
		return "None"
	}
	if e == JumpBackType_Immediate {
		return "Immediate"
	}
	if e == JumpBackType_AfterCombo {
		return "AfterCombo"
	}

	return ""

}

func GetJumpBackType(tp string) (JumpBackType, bool) {
	switch tp {
	case "None":
		return JumpBackType_None, true
	case "Immediate":
		return JumpBackType_Immediate, true
	case "AfterCombo":
		return JumpBackType_AfterCombo, true

	}
	return -1, false

}

func GetJumpBackType_ByAlias(tp string) (JumpBackType, bool) {
	switch tp {
	case "无":
		return JumpBackType_None, true
	case "立即":
		return JumpBackType_Immediate, true
	case "连击后":
		return JumpBackType_AfterCombo, true

	}
	return -1, false

}

func (e ConnectCameraType) String() string {

	if e == ConnectCameraType_None {
		return "None"
	}
	if e == ConnectCameraType_Lerp {
		return "Lerp"
	}
	if e == ConnectCameraType_Immeidate {
		return "Immeidate"
	}

	return ""

}

func GetConnectCameraType(tp string) (ConnectCameraType, bool) {
	switch tp {
	case "None":
		return ConnectCameraType_None, true
	case "Lerp":
		return ConnectCameraType_Lerp, true
	case "Immeidate":
		return ConnectCameraType_Immeidate, true

	}
	return -1, false

}

func GetConnectCameraType_ByAlias(tp string) (ConnectCameraType, bool) {
	switch tp {
	case "无":
		return ConnectCameraType_None, true
	case "缓动":
		return ConnectCameraType_Lerp, true
	case "立即":
		return ConnectCameraType_Immeidate, true

	}
	return -1, false

}

func (e GrabPointType) String() string {

	if e == GrabPointType_Check {
		return "Check"
	}
	if e == GrabPointType_Head {
		return "Head"
	}

	return ""

}

func GetGrabPointType(tp string) (GrabPointType, bool) {
	switch tp {
	case "Check":
		return GrabPointType_Check, true
	case "Head":
		return GrabPointType_Head, true

	}
	return -1, false

}

func GetGrabPointType_ByAlias(tp string) (GrabPointType, bool) {
	switch tp {
	case "胸":
		return GrabPointType_Check, true
	case "头":
		return GrabPointType_Head, true

	}
	return -1, false

}

func (e ImportSkillID) String() string {

	if e == ImportSkillID_CommonEndrue {
		return "CommonEndrue"
	}
	if e == ImportSkillID_CommoRemoteNormal {
		return "CommoRemoteNormal"
	}
	if e == ImportSkillID_CounterAttack {
		return "CounterAttack"
	}
	if e == ImportSkillID_ComboAttack {
		return "ComboAttack"
	}
	if e == ImportSkillID_FirstMoveGain {
		return "FirstMoveGain"
	}
	if e == ImportSkillID_CNT {
		return "CNT"
	}

	return ""

}

func GetImportSkillID(tp string) (ImportSkillID, bool) {
	switch tp {
	case "CommonEndrue":
		return ImportSkillID_CommonEndrue, true
	case "CommoRemoteNormal":
		return ImportSkillID_CommoRemoteNormal, true
	case "CounterAttack":
		return ImportSkillID_CounterAttack, true
	case "ComboAttack":
		return ImportSkillID_ComboAttack, true
	case "FirstMoveGain":
		return ImportSkillID_FirstMoveGain, true
	case "CNT":
		return ImportSkillID_CNT, true

	}
	return -1, false

}

func GetImportSkillID_ByAlias(tp string) (ImportSkillID, bool) {
	switch tp {
	case "通用受击霸体":
		return ImportSkillID_CommonEndrue, true
	case "远程普攻检测":
		return ImportSkillID_CommoRemoteNormal, true
	case "反击技能":
		return ImportSkillID_CounterAttack, true
	case "连击技能":
		return ImportSkillID_ComboAttack, true
	case "先手增益":
		return ImportSkillID_FirstMoveGain, true
	case "结束":
		return ImportSkillID_CNT, true

	}
	return -1, false

}

func (e SkillTagType) String() string {

	if e == SkillTagType_NormalAttack {
		return "NormalAttack"
	}
	if e == SkillTagType_SkillA {
		return "SkillA"
	}
	if e == SkillTagType_SkillB {
		return "SkillB"
	}
	if e == SkillTagType_Ultimate {
		return "Ultimate"
	}
	if e == SkillTagType_UltimateNext {
		return "UltimateNext"
	}
	if e == SkillTagType_Negative {
		return "Negative"
	}
	if e == SkillTagType_Connect {
		return "Connect"
	}
	if e == SkillTagType_SkillC {
		return "SkillC"
	}
	if e == SkillTagType_NoCDReplace {
		return "NoCDReplace"
	}
	if e == SkillTagType_AtrSkillA {
		return "AtrSkillA"
	}
	if e == SkillTagType_AtrSkillB {
		return "AtrSkillB"
	}
	if e == SkillTagType_AssistSkillA {
		return "AssistSkillA"
	}
	if e == SkillTagType_AssistSkillB {
		return "AssistSkillB"
	}
	if e == SkillTagType_AssistSkillC {
		return "AssistSkillC"
	}
	if e == SkillTagType_PetSkillA {
		return "PetSkillA"
	}
	if e == SkillTagType_PetSkillB {
		return "PetSkillB"
	}
	if e == SkillTagType_PetSkillManifest {
		return "PetSkillManifest"
	}
	if e == SkillTagType_None {
		return "None"
	}

	return ""

}

func GetSkillTagType(tp string) (SkillTagType, bool) {
	switch tp {
	case "NormalAttack":
		return SkillTagType_NormalAttack, true
	case "SkillA":
		return SkillTagType_SkillA, true
	case "SkillB":
		return SkillTagType_SkillB, true
	case "Ultimate":
		return SkillTagType_Ultimate, true
	case "UltimateNext":
		return SkillTagType_UltimateNext, true
	case "Negative":
		return SkillTagType_Negative, true
	case "Connect":
		return SkillTagType_Connect, true
	case "SkillC":
		return SkillTagType_SkillC, true
	case "NoCDReplace":
		return SkillTagType_NoCDReplace, true
	case "AtrSkillA":
		return SkillTagType_AtrSkillA, true
	case "AtrSkillB":
		return SkillTagType_AtrSkillB, true
	case "AssistSkillA":
		return SkillTagType_AssistSkillA, true
	case "AssistSkillB":
		return SkillTagType_AssistSkillB, true
	case "AssistSkillC":
		return SkillTagType_AssistSkillC, true
	case "PetSkillA":
		return SkillTagType_PetSkillA, true
	case "PetSkillB":
		return SkillTagType_PetSkillB, true
	case "PetSkillManifest":
		return SkillTagType_PetSkillManifest, true
	case "None":
		return SkillTagType_None, true

	}
	return -1, false

}

func GetSkillTagType_ByAlias(tp string) (SkillTagType, bool) {
	switch tp {
	case "普攻":
		return SkillTagType_NormalAttack, true
	case "小招A":
		return SkillTagType_SkillA, true
	case "小招B":
		return SkillTagType_SkillB, true
	case "大招":
		return SkillTagType_Ultimate, true
	case "大招后续":
		return SkillTagType_UltimateNext, true
	case "被动":
		return SkillTagType_Negative, true
	case "连携":
		return SkillTagType_Connect, true
	case "小招C":
		return SkillTagType_SkillC, true
	case "无冷却普攻替换":
		return SkillTagType_NoCDReplace, true
	case "奥义A":
		return SkillTagType_AtrSkillA, true
	case "奥义B":
		return SkillTagType_AtrSkillB, true
	case "援护A":
		return SkillTagType_AssistSkillA, true
	case "援护B":
		return SkillTagType_AssistSkillB, true
	case "援护C":
		return SkillTagType_AssistSkillC, true
	case "宠物技能A":
		return SkillTagType_PetSkillA, true
	case "宠物技能B":
		return SkillTagType_PetSkillB, true
	case "宠物技能幻化":
		return SkillTagType_PetSkillManifest, true
	case "无":
		return SkillTagType_None, true

	}
	return -1, false

}

func (e ConnectFocusType) String() string {

	if e == ConnectFocusType_Form {
		return "Form"
	}
	if e == ConnectFocusType_Target {
		return "Target"
	}
	if e == ConnectFocusType_Center {
		return "Center"
	}

	return ""

}

func GetConnectFocusType(tp string) (ConnectFocusType, bool) {
	switch tp {
	case "Form":
		return ConnectFocusType_Form, true
	case "Target":
		return ConnectFocusType_Target, true
	case "Center":
		return ConnectFocusType_Center, true

	}
	return -1, false

}

func GetConnectFocusType_ByAlias(tp string) (ConnectFocusType, bool) {
	switch tp {
	case "发起者":
		return ConnectFocusType_Form, true
	case "目标":
		return ConnectFocusType_Target, true
	case "中点":
		return ConnectFocusType_Center, true

	}
	return -1, false

}

func (e ComboFntShowType) String() string {

	if e == ComboFntShowType_AmplifyScale {
		return "AmplifyScale"
	}
	if e == ComboFntShowType_Flash {
		return "Flash"
	}
	if e == ComboFntShowType_Appear {
		return "Appear"
	}

	return ""

}

func GetComboFntShowType(tp string) (ComboFntShowType, bool) {
	switch tp {
	case "AmplifyScale":
		return ComboFntShowType_AmplifyScale, true
	case "Flash":
		return ComboFntShowType_Flash, true
	case "Appear":
		return ComboFntShowType_Appear, true

	}
	return -1, false

}

func GetComboFntShowType_ByAlias(tp string) (ComboFntShowType, bool) {
	switch tp {
	case "放缩":
		return ComboFntShowType_AmplifyScale, true
	case "位移":
		return ComboFntShowType_Flash, true
	case "渐显":
		return ComboFntShowType_Appear, true

	}
	return -1, false

}

func (e ComboFntMoveInDir) String() string {

	if e == ComboFntMoveInDir_Up {
		return "Up"
	}
	if e == ComboFntMoveInDir_Down {
		return "Down"
	}
	if e == ComboFntMoveInDir_Left {
		return "Left"
	}
	if e == ComboFntMoveInDir_Right {
		return "Right"
	}
	if e == ComboFntMoveInDir_LeftUp {
		return "LeftUp"
	}
	if e == ComboFntMoveInDir_RightUp {
		return "RightUp"
	}
	if e == ComboFntMoveInDir_LeftDown {
		return "LeftDown"
	}
	if e == ComboFntMoveInDir_RightDown {
		return "RightDown"
	}

	return ""

}

func GetComboFntMoveInDir(tp string) (ComboFntMoveInDir, bool) {
	switch tp {
	case "Up":
		return ComboFntMoveInDir_Up, true
	case "Down":
		return ComboFntMoveInDir_Down, true
	case "Left":
		return ComboFntMoveInDir_Left, true
	case "Right":
		return ComboFntMoveInDir_Right, true
	case "LeftUp":
		return ComboFntMoveInDir_LeftUp, true
	case "RightUp":
		return ComboFntMoveInDir_RightUp, true
	case "LeftDown":
		return ComboFntMoveInDir_LeftDown, true
	case "RightDown":
		return ComboFntMoveInDir_RightDown, true

	}
	return -1, false

}

func GetComboFntMoveInDir_ByAlias(tp string) (ComboFntMoveInDir, bool) {
	switch tp {
	case "上":
		return ComboFntMoveInDir_Up, true
	case "下":
		return ComboFntMoveInDir_Down, true
	case "左":
		return ComboFntMoveInDir_Left, true
	case "右":
		return ComboFntMoveInDir_Right, true
	case "左上":
		return ComboFntMoveInDir_LeftUp, true
	case "右上":
		return ComboFntMoveInDir_RightUp, true
	case "左下":
		return ComboFntMoveInDir_LeftDown, true
	case "右下":
		return ComboFntMoveInDir_RightDown, true

	}
	return -1, false

}

func (e BattleAreaType) String() string {

	if e == BattleAreaType_Left {
		return "Left"
	}
	if e == BattleAreaType_Midddle {
		return "Midddle"
	}
	if e == BattleAreaType_Right {
		return "Right"
	}
	if e == BattleAreaType_Self {
		return "Self"
	}
	if e == BattleAreaType_Enemy {
		return "Enemy"
	}

	return ""

}

func GetBattleAreaType(tp string) (BattleAreaType, bool) {
	switch tp {
	case "Left":
		return BattleAreaType_Left, true
	case "Midddle":
		return BattleAreaType_Midddle, true
	case "Right":
		return BattleAreaType_Right, true
	case "Self":
		return BattleAreaType_Self, true
	case "Enemy":
		return BattleAreaType_Enemy, true

	}
	return -1, false

}

func GetBattleAreaType_ByAlias(tp string) (BattleAreaType, bool) {
	switch tp {
	case "左区":
		return BattleAreaType_Left, true
	case "中区":
		return BattleAreaType_Midddle, true
	case "右区":
		return BattleAreaType_Right, true
	case "己方区":
		return BattleAreaType_Self, true
	case "敌方区":
		return BattleAreaType_Enemy, true

	}
	return -1, false

}

func (e CameraLockType) String() string {

	if e == CameraLockType_From {
		return "From"
	}
	if e == CameraLockType_Target {
		return "Target"
	}
	if e == CameraLockType_Center {
		return "Center"
	}

	return ""

}

func GetCameraLockType(tp string) (CameraLockType, bool) {
	switch tp {
	case "From":
		return CameraLockType_From, true
	case "Target":
		return CameraLockType_Target, true
	case "Center":
		return CameraLockType_Center, true

	}
	return -1, false

}

func GetCameraLockType_ByAlias(tp string) (CameraLockType, bool) {
	switch tp {
	case "发起者":
		return CameraLockType_From, true
	case "目标":
		return CameraLockType_Target, true
	case "中心":
		return CameraLockType_Center, true

	}
	return -1, false

}

func (e EffectAddParam) String() string {

	if e == EffectAddParam_SkillDmg {
		return "SkillDmg"
	}
	if e == EffectAddParam_Heal {
		return "Heal"
	}
	if e == EffectAddParam_CNT {
		return "CNT"
	}

	return ""

}

func GetEffectAddParam(tp string) (EffectAddParam, bool) {
	switch tp {
	case "SkillDmg":
		return EffectAddParam_SkillDmg, true
	case "Heal":
		return EffectAddParam_Heal, true
	case "CNT":
		return EffectAddParam_CNT, true

	}
	return -1, false

}

func GetEffectAddParam_ByAlias(tp string) (EffectAddParam, bool) {
	switch tp {
	case "技伤":
		return EffectAddParam_SkillDmg, true
	case "技疗":
		return EffectAddParam_Heal, true
	case "结束":
		return EffectAddParam_CNT, true

	}
	return -1, false

}

func (e DmgCalculateRateType) String() string {

	if e == DmgCalculateRateType_Crital {
		return "Crital"
	}
	if e == DmgCalculateRateType_HealCrital {
		return "HealCrital"
	}
	if e == DmgCalculateRateType_Dodge {
		return "Dodge"
	}
	if e == DmgCalculateRateType_Block {
		return "Block"
	}

	return ""

}

func GetDmgCalculateRateType(tp string) (DmgCalculateRateType, bool) {
	switch tp {
	case "Crital":
		return DmgCalculateRateType_Crital, true
	case "HealCrital":
		return DmgCalculateRateType_HealCrital, true
	case "Dodge":
		return DmgCalculateRateType_Dodge, true
	case "Block":
		return DmgCalculateRateType_Block, true

	}
	return -1, false

}

func GetDmgCalculateRateType_ByAlias(tp string) (DmgCalculateRateType, bool) {
	switch tp {
	case "暴击":
		return DmgCalculateRateType_Crital, true
	case "治疗暴击":
		return DmgCalculateRateType_HealCrital, true
	case "闪避":
		return DmgCalculateRateType_Dodge, true
	case "格挡":
		return DmgCalculateRateType_Block, true

	}
	return -1, false

}

func (e SummonAppearType) String() string {

	if e == SummonAppearType_None {
		return "None"
	}
	if e == SummonAppearType_Straight {
		return "Straight"
	}
	if e == SummonAppearType_Parabola {
		return "Parabola"
	}
	if e == SummonAppearType_Fall {
		return "Fall"
	}

	return ""

}

func GetSummonAppearType(tp string) (SummonAppearType, bool) {
	switch tp {
	case "None":
		return SummonAppearType_None, true
	case "Straight":
		return SummonAppearType_Straight, true
	case "Parabola":
		return SummonAppearType_Parabola, true
	case "Fall":
		return SummonAppearType_Fall, true

	}
	return -1, false

}

func GetSummonAppearType_ByAlias(tp string) (SummonAppearType, bool) {
	switch tp {
	case "无":
		return SummonAppearType_None, true
	case "直线冲刺":
		return SummonAppearType_Straight, true
	case "曲线进场":
		return SummonAppearType_Parabola, true
	case "上空掉落":
		return SummonAppearType_Fall, true

	}
	return -1, false

}
