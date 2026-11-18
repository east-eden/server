package config

// LD技能表-拳皇/LDAttrEnum
type LDAttrType int

const (
	LDAttrType_BEGIN               LDAttrType = 1    //枚举开始
	LDAttrType_HP                  LDAttrType = 1    //血量
	LDAttrType_MAXHP               LDAttrType = 2    //血量上限
	LDAttrType_ATTACK              LDAttrType = 3    //攻击力
	LDAttrType_MAGIC               LDAttrType = 4    //魔法攻击力
	LDAttrType_DEFENCE             LDAttrType = 5    //防御力
	LDAttrType_MAGICDEFENCE        LDAttrType = 6    //法术防御
	LDAttrType_IGNDEFENCE          LDAttrType = 7    //忽视防御
	LDAttrType_IGNMAGICDEFENCE     LDAttrType = 8    //忽视法防
	LDAttrType_IGNHEAL             LDAttrType = 9    //抗疗
	LDAttrType_CRITAL              LDAttrType = 10   //暴击率
	LDAttrType_IGNCRITAL           LDAttrType = 11   //忽视暴击
	LDAttrType_DODGE               LDAttrType = 12   //闪避率
	LDAttrType_IGNDODGE            LDAttrType = 13   //忽视闪避
	LDAttrType_TRUEHARM            LDAttrType = 14   //真实伤害
	LDAttrType_IGNTRUEHARM         LDAttrType = 15   //忽视真伤
	LDAttrType_FIREHARM            LDAttrType = 16   //火系伤害
	LDAttrType_IGNFIREHARM         LDAttrType = 17   //忽视火伤
	LDAttrType_MOVESPEED           LDAttrType = 18   //移动速度
	LDAttrType_ATTACKSPEED         LDAttrType = 19   //攻击速度
	LDAttrType_HITRATE             LDAttrType = 20   //命中率
	LDAttrType_BLASTS              LDAttrType = 21   //暴伤
	LDAttrType_IGNBLASTS           LDAttrType = 22   //忽视暴伤
	LDAttrType_MAXSP               LDAttrType = 23   //最大精神力
	LDAttrType_SP                  LDAttrType = 24   //精神力
	LDAttrType_DEFENCEIMMUNE       LDAttrType = 25   //防御免伤比
	LDAttrType_BLOCK               LDAttrType = 26   //格挡率
	LDAttrType_DAMAGEADDPRECENT    LDAttrType = 27   //伤害百分比加成
	LDAttrType_DAMAGEADD           LDAttrType = 28   //伤害固定加成
	LDAttrType_PENETRATE           LDAttrType = 29   //穿透率
	LDAttrType_TRICKADD            LDAttrType = 30   //大招增伤
	LDAttrType_BUFFADD             LDAttrType = 31   //状态加强
	LDAttrType_SKILLTRIGGERADD     LDAttrType = 32   //技能触发率增长
	LDAttrType_DEFENCEBLOCK        LDAttrType = 33   //忽视格挡
	LDAttrType_DEFDAMAGEADDPRECENT LDAttrType = 34   //百分比减伤
	LDAttrType_DEFDAMAGEADD        LDAttrType = 35   //固定减伤
	LDAttrType_DEFPENETRATE        LDAttrType = 36   //忽视穿透
	LDAttrType_DEFTRICKADD         LDAttrType = 37   //大招减伤
	LDAttrType_DEFBUFFADD          LDAttrType = 38   //状态抵抗
	LDAttrType_DEFSKILLTRIGGERADD  LDAttrType = 39   //技能触发率抵抗
	LDAttrType_BLOCKSUBPRECENT     LDAttrType = 40   //格挡免伤比例
	LDAttrType_DEFSTUN             LDAttrType = 41   //眩晕抵抗
	LDAttrType_DEFFREEZE           LDAttrType = 42   //冰冻抵抗
	LDAttrType_DEFSILENT           LDAttrType = 43   //沉默抵抗
	LDAttrType_DEFBLOOD            LDAttrType = 44   //流血抵抗
	LDAttrType_DEFPOISON           LDAttrType = 45   //中毒抵抗
	LDAttrType_DEFBURN             LDAttrType = 46   //灼烧抵抗
	LDAttrType_SPPRECENTADD        LDAttrType = 47   //精神力百分比加成
	LDAttrType_SPADD               LDAttrType = 48   //精神力固定加成
	LDAttrType_MonsterDefence      LDAttrType = 49   //怪物防御力
	LDAttrType_MonsterSpeedDown    LDAttrType = 50   //怪物移速
	LDAttrType_MonsterDizzyDefence LDAttrType = 51   //怪物眩晕抵抗
	LDAttrType_HeroAttackSpeedAdd  LDAttrType = 52   //英雄攻击速度
	LDAttrType_SpAdd               LDAttrType = 53   //英雄sp恢复速度
	LDAttrType_HeroMoveSpeedAdd    LDAttrType = 54   //英雄移动速度
	LDAttrType_PeopleMaxCount      LDAttrType = 55   //兵种上限
	LDAttrType_HeroAttackAdd       LDAttrType = 56   //英雄攻击力增加
	LDAttrType_BossDamageAdd       LDAttrType = 57   //首领增伤
	LDAttrType_PhysicalDamageAdd   LDAttrType = 58   //物理增伤
	LDAttrType_SkillDamageAdd      LDAttrType = 59   //技能增伤
	LDAttrType_CritalAdd           LDAttrType = 60   //暴击率增加
	LDAttrType_BlastAdd            LDAttrType = 61   //爆伤加成
	LDAttrType_SkillUseRate        LDAttrType = 62   //技能触发率
	LDAttrType_StunTimeAdd         LDAttrType = 63   //眩晕效果加成
	LDAttrType_StunDamageAdd       LDAttrType = 64   //眩晕伤害加成
	LDAttrType_MagicDamageAdd      LDAttrType = 65   //魔法伤害加成
	LDAttrType_MonsterMagicDefence LDAttrType = 66   //怪物魔法抗性
	LDAttrType_ARMOR               LDAttrType = 67   //护甲值
	LDAttrType_BossMaxHP           LDAttrType = 68   //Boss最大血量加成
	LDAttrType_MonsterMaxHp        LDAttrType = 69   //小怪最大血量加成
	LDAttrType_EpicAttackAdd       LDAttrType = 70   //神话英雄攻击加成
	LDAttrType_AvoidDamageParam    LDAttrType = 71   //免伤系数
	LDAttrType_MagicArmor          LDAttrType = 72   //魔法护甲
	LDAttrType_AvoidMagicParam     LDAttrType = 73   //魔免系数
	LDAttrType_CNT                 LDAttrType = 74   //结束标志
	LDAttrType_MAXHPPERCENT        LDAttrType = 202  //血量上限百分比
	LDAttrType_ATTACKPERCENT       LDAttrType = 203  //攻击力百分比
	LDAttrType_MAGICPERCENT        LDAttrType = 204  //魔法攻击力百分比
	LDAttrType_DEFENCEPERCENT      LDAttrType = 205  //防御力百分比
	LDAttrType_MAGICDEFENCEPERCENT LDAttrType = 206  //法术防御力百分比
	LDAttrType_MOVESPEEDPERCENT    LDAttrType = 218  //移动速度百分比
	LDAttrType_ATTACKSPEEDPERCENT  LDAttrType = 219  //攻击速度百分比
	LDAttrType_Extra               LDAttrType = 500  //额外战斗属性开始
	LDAttrType_ATTACKWARRIOR       LDAttrType = 701  //对战士加成
	LDAttrType_ATTACKMAGE          LDAttrType = 702  //对法师加成
	LDAttrType_ATTACKPRIEST        LDAttrType = 703  //对牧师加成
	LDAttrType_ATTACKKNIGHT        LDAttrType = 704  //对骑士加成
	LDAttrType_ATTACKARCHER        LDAttrType = 705  //对射手加成
	LDAttrType_ATTACKASSASIN       LDAttrType = 706  //对刺客加成
	LDAttrType_ATTACKSAVAGE        LDAttrType = 707  //对野蛮人加成
	LDAttrType_ATTACKMONK          LDAttrType = 708  //对武僧加成
	LDAttrType_ATTACKWARKLOCK      LDAttrType = 709  //对术士加成
	LDAttrType_ATTACKNEO           LDAttrType = 710  //对亡灵法师加成
	LDAttrType_ATTACKENCHATER      LDAttrType = 711  //对幻术师加成
	LDAttrType_ATTACKSACRIFICE     LDAttrType = 712  //对祭祀加成
	LDAttrType_ATTACKBARD          LDAttrType = 713  //对吟游诗人加成
	LDAttrType_ATTACKSHAMAN        LDAttrType = 714  //对萨满加成
	LDAttrType_ATTACKDRUID         LDAttrType = 715  //对德鲁伊加成
	LDAttrType_ATTACKPALADIN       LDAttrType = 716  //对圣骑士加成
	LDAttrType_ATTACKDEADKNIGHT    LDAttrType = 717  //对死骑加成
	LDAttrType_ATTACKHUNTER        LDAttrType = 718  //对猎人加成
	LDAttrType_ATTACKRANGER        LDAttrType = 719  //对游侠加成
	LDAttrType_ATTACKSTEAL         LDAttrType = 720  //对潜行者加成
	LDAttrType_ATTACKDEMONHUNTER   LDAttrType = 721  //对恶魔猎手加成
	LDAttrType_ATTACKHUMAN         LDAttrType = 801  //对人类加成
	LDAttrType_ATTACKORC           LDAttrType = 802  //对兽人加成
	LDAttrType_ATTACKTAUREN        LDAttrType = 803  //对牛头人加成
	LDAttrType_ATTACKUNDEAD        LDAttrType = 804  //对亡灵加成
	LDAttrType_ATTACKDWARF         LDAttrType = 805  //对矮人加成
	LDAttrType_ATTACKELF           LDAttrType = 806  //对精灵加成
	LDAttrType_ATTACKTROLL         LDAttrType = 807  //对巨魔加成
	LDAttrType_ATTACKGNOME         LDAttrType = 808  //对侏儒加成
	LDAttrType_ATTACKPANDA         LDAttrType = 809  //对熊猫人加成
	LDAttrType_ATTACKGOD           LDAttrType = 810  //对神加成
	LDAttrType_ATTACKDEMON         LDAttrType = 811  //对恶魔加成
	LDAttrType_ATTACKWORGEN        LDAttrType = 812  //对狼人加成
	LDAttrType_ATTACKNAGA          LDAttrType = 813  //对娜迦加成
	LDAttrType_ATTACKBIRDMAN       LDAttrType = 814  //对鸟人加成
	LDAttrType_ATTACKLIONMAN       LDAttrType = 815  //对狮身人加成
	LDAttrType_ATTACKPIGMAN        LDAttrType = 816  //对野猪人加成
	LDAttrType_ATTACKGOLEM         LDAttrType = 817  //对山顶巨人加成
	LDAttrType_ATTACKWOOD          LDAttrType = 818  //对树精加成
	LDAttrType_ATTACKPUPPET        LDAttrType = 819  //对傀儡加成
	LDAttrType_ATTACKCENTAUR       LDAttrType = 820  //对半人马加成
	LDAttrType_ATTACKSHEEPMAN      LDAttrType = 821  //对萨特加成
	LDAttrType_ATTACKDOGMAN        LDAttrType = 822  //对狗头人加成
	LDAttrType_ATTACKBEARMAN       LDAttrType = 823  //对熊人加成
	LDAttrType_ATTACKMONKEYMAN     LDAttrType = 824  //对猩猩人加成
	LDAttrType_ATTACKCATMAN        LDAttrType = 825  //对猫人加成
	LDAttrType_ATTACKVANISHMAN     LDAttrType = 826  //对隐身人加成
	LDAttrType_ATTACKUNKOWNMAN     LDAttrType = 827  //对穿越者加成
	LDAttrType_ATTACKMONSTERONE    LDAttrType = 901  //对类型1怪物加成
	LDAttrType_ATTACKMONSTERTWO    LDAttrType = 902  //对类型2怪物加成
	LDAttrType_ATTACKMONSTERTREE   LDAttrType = 903  //对类型3怪物加成
	LDAttrType_ATTACKMONSTERFOUR   LDAttrType = 904  //对类型4怪物加成
	LDAttrType_ATTACKMONSTERFIVE   LDAttrType = 905  //对类型5怪物加成
	LDAttrType_HUMAADD             LDAttrType = 1001 //人类加成
	LDAttrType_ORCADD              LDAttrType = 1002 //兽人加成
	LDAttrType_TAUREADD            LDAttrType = 1003 //牛头人加成
	LDAttrType_UNDEAADD            LDAttrType = 1004 //亡灵加成
	LDAttrType_DWARADD             LDAttrType = 1005 //矮人加成
	LDAttrType_ELFADD              LDAttrType = 1006 //精灵加成
	LDAttrType_TROLADD             LDAttrType = 1007 //巨魔加成
	LDAttrType_GNOMADD             LDAttrType = 1008 //侏儒加成
	LDAttrType_PANDADD             LDAttrType = 1009 //熊猫人加成
	LDAttrType_GODADD              LDAttrType = 1010 //神加成
	LDAttrType_DEMOADD             LDAttrType = 1011 //恶魔加成
	LDAttrType_WORGEADD            LDAttrType = 1012 //狼人加成
	LDAttrType_NAGADD              LDAttrType = 1013 //娜迦加成
	LDAttrType_BIRDMAADD           LDAttrType = 1014 //鸟人加成
	LDAttrType_LIONMAADD           LDAttrType = 1015 //狮身人加成
	LDAttrType_PIGMAADD            LDAttrType = 1016 //野猪人加成
	LDAttrType_GOLEADD             LDAttrType = 1017 //山顶巨人加成
	LDAttrType_WOOADD              LDAttrType = 1018 //树精加成
	LDAttrType_PUPPEADD            LDAttrType = 1019 //傀儡加成
	LDAttrType_CENTAUADD           LDAttrType = 1020 //半人马加成
	LDAttrType_SHEEPMAADD          LDAttrType = 1021 //萨特加成
	LDAttrType_DOGMAADD            LDAttrType = 1022 //狗头人加成
	LDAttrType_BEARMAADD           LDAttrType = 1023 //熊人加成
	LDAttrType_MONKEYMAADD         LDAttrType = 1024 //猩猩人加成
	LDAttrType_CATMAADD            LDAttrType = 1025 //猫人加成
	LDAttrType_VANISHMAADD         LDAttrType = 1026 //隐身人加成
	LDAttrType_UNKOWNMAADD         LDAttrType = 1027 //穿越者加成
	LDAttrType_WARRIORADD          LDAttrType = 1001 //战士加成
	LDAttrType_MAGEADD             LDAttrType = 1002 //法师加成
	LDAttrType_PRIESTADD           LDAttrType = 1003 //牧师加成
	LDAttrType_KNIGHTADD           LDAttrType = 1004 //骑士加成
	LDAttrType_ARCHERADD           LDAttrType = 1005 //射手加成
	LDAttrType_ASSASINADD          LDAttrType = 1006 //刺客加成
	LDAttrType_SAVAGEADD           LDAttrType = 1007 //野蛮人加成
	LDAttrType_MONKADD             LDAttrType = 1008 //武僧加成
	LDAttrType_WARKLOCKADD         LDAttrType = 1009 //术士加成
	LDAttrType_NEOADD              LDAttrType = 1010 //亡灵法师加成
	LDAttrType_ENCHATERADD         LDAttrType = 1011 //幻术师加成
	LDAttrType_SACRIFICEADD        LDAttrType = 1012 //祭祀加成
	LDAttrType_BARDADD             LDAttrType = 1013 //吟游诗人加成
	LDAttrType_SHAMANADD           LDAttrType = 1014 //萨满加成
	LDAttrType_DRUIDADD            LDAttrType = 1015 //德鲁伊加成
	LDAttrType_PALADINADD          LDAttrType = 1016 //圣骑士加成
	LDAttrType_DEADKNIGHTADD       LDAttrType = 1017 //死骑加成
	LDAttrType_HUNTERADD           LDAttrType = 1018 //猎人加成
	LDAttrType_RANGERADD           LDAttrType = 1019 //游侠加成
	LDAttrType_STEALADD            LDAttrType = 1020 //潜行者加成
	LDAttrType_DEMONHUNTERADD      LDAttrType = 1021 //恶魔猎手加成
	LDAttrType_ADDLEVEL            LDAttrType = 1028 //增加等级
	LDAttrType_ExtraCNT            LDAttrType = 1100 //额外属性结束
	LDAttrType_END                 LDAttrType = 1101 //枚举结束
)

// LD技能表-拳皇/LDAttrEnum
type LDAttrTitleType int

const (
	LDAttrTitleType_BEGIN   LDAttrTitleType = 1 //枚举开始
	LDAttrTitleType_Basic   LDAttrTitleType = 1 //基础属性
	LDAttrTitleType_Advance LDAttrTitleType = 2 //进阶属性
	LDAttrTitleType_END     LDAttrTitleType = 3 //枚举结束
)

func (e LDAttrType) String() string {

	if e == LDAttrType_HP {
		return "HP"
	}
	if e == LDAttrType_MAXHP {
		return "MAXHP"
	}
	if e == LDAttrType_ATTACK {
		return "ATTACK"
	}
	if e == LDAttrType_MAGIC {
		return "MAGIC"
	}
	if e == LDAttrType_DEFENCE {
		return "DEFENCE"
	}
	if e == LDAttrType_MAGICDEFENCE {
		return "MAGICDEFENCE"
	}
	if e == LDAttrType_IGNDEFENCE {
		return "IGNDEFENCE"
	}
	if e == LDAttrType_IGNMAGICDEFENCE {
		return "IGNMAGICDEFENCE"
	}
	if e == LDAttrType_IGNHEAL {
		return "IGNHEAL"
	}
	if e == LDAttrType_CRITAL {
		return "CRITAL"
	}
	if e == LDAttrType_IGNCRITAL {
		return "IGNCRITAL"
	}
	if e == LDAttrType_DODGE {
		return "DODGE"
	}
	if e == LDAttrType_IGNDODGE {
		return "IGNDODGE"
	}
	if e == LDAttrType_TRUEHARM {
		return "TRUEHARM"
	}
	if e == LDAttrType_IGNTRUEHARM {
		return "IGNTRUEHARM"
	}
	if e == LDAttrType_FIREHARM {
		return "FIREHARM"
	}
	if e == LDAttrType_IGNFIREHARM {
		return "IGNFIREHARM"
	}
	if e == LDAttrType_MOVESPEED {
		return "MOVESPEED"
	}
	if e == LDAttrType_ATTACKSPEED {
		return "ATTACKSPEED"
	}
	if e == LDAttrType_HITRATE {
		return "HITRATE"
	}
	if e == LDAttrType_BLASTS {
		return "BLASTS"
	}
	if e == LDAttrType_IGNBLASTS {
		return "IGNBLASTS"
	}
	if e == LDAttrType_MAXSP {
		return "MAXSP"
	}
	if e == LDAttrType_SP {
		return "SP"
	}
	if e == LDAttrType_DEFENCEIMMUNE {
		return "DEFENCEIMMUNE"
	}
	if e == LDAttrType_BLOCK {
		return "BLOCK"
	}
	if e == LDAttrType_DAMAGEADDPRECENT {
		return "DAMAGEADDPRECENT"
	}
	if e == LDAttrType_DAMAGEADD {
		return "DAMAGEADD"
	}
	if e == LDAttrType_PENETRATE {
		return "PENETRATE"
	}
	if e == LDAttrType_TRICKADD {
		return "TRICKADD"
	}
	if e == LDAttrType_BUFFADD {
		return "BUFFADD"
	}
	if e == LDAttrType_SKILLTRIGGERADD {
		return "SKILLTRIGGERADD"
	}
	if e == LDAttrType_DEFENCEBLOCK {
		return "DEFENCEBLOCK"
	}
	if e == LDAttrType_DEFDAMAGEADDPRECENT {
		return "DEFDAMAGEADDPRECENT"
	}
	if e == LDAttrType_DEFDAMAGEADD {
		return "DEFDAMAGEADD"
	}
	if e == LDAttrType_DEFPENETRATE {
		return "DEFPENETRATE"
	}
	if e == LDAttrType_DEFTRICKADD {
		return "DEFTRICKADD"
	}
	if e == LDAttrType_DEFBUFFADD {
		return "DEFBUFFADD"
	}
	if e == LDAttrType_DEFSKILLTRIGGERADD {
		return "DEFSKILLTRIGGERADD"
	}
	if e == LDAttrType_BLOCKSUBPRECENT {
		return "BLOCKSUBPRECENT"
	}
	if e == LDAttrType_DEFSTUN {
		return "DEFSTUN"
	}
	if e == LDAttrType_DEFFREEZE {
		return "DEFFREEZE"
	}
	if e == LDAttrType_DEFSILENT {
		return "DEFSILENT"
	}
	if e == LDAttrType_DEFBLOOD {
		return "DEFBLOOD"
	}
	if e == LDAttrType_DEFPOISON {
		return "DEFPOISON"
	}
	if e == LDAttrType_DEFBURN {
		return "DEFBURN"
	}
	if e == LDAttrType_SPPRECENTADD {
		return "SPPRECENTADD"
	}
	if e == LDAttrType_SPADD {
		return "SPADD"
	}
	if e == LDAttrType_MonsterDefence {
		return "MonsterDefence"
	}
	if e == LDAttrType_MonsterSpeedDown {
		return "MonsterSpeedDown"
	}
	if e == LDAttrType_MonsterDizzyDefence {
		return "MonsterDizzyDefence"
	}
	if e == LDAttrType_HeroAttackSpeedAdd {
		return "HeroAttackSpeedAdd"
	}
	if e == LDAttrType_SpAdd {
		return "SpAdd"
	}
	if e == LDAttrType_HeroMoveSpeedAdd {
		return "HeroMoveSpeedAdd"
	}
	if e == LDAttrType_PeopleMaxCount {
		return "PeopleMaxCount"
	}
	if e == LDAttrType_HeroAttackAdd {
		return "HeroAttackAdd"
	}
	if e == LDAttrType_BossDamageAdd {
		return "BossDamageAdd"
	}
	if e == LDAttrType_PhysicalDamageAdd {
		return "PhysicalDamageAdd"
	}
	if e == LDAttrType_SkillDamageAdd {
		return "SkillDamageAdd"
	}
	if e == LDAttrType_CritalAdd {
		return "CritalAdd"
	}
	if e == LDAttrType_BlastAdd {
		return "BlastAdd"
	}
	if e == LDAttrType_SkillUseRate {
		return "SkillUseRate"
	}
	if e == LDAttrType_StunTimeAdd {
		return "StunTimeAdd"
	}
	if e == LDAttrType_StunDamageAdd {
		return "StunDamageAdd"
	}
	if e == LDAttrType_MagicDamageAdd {
		return "MagicDamageAdd"
	}
	if e == LDAttrType_MonsterMagicDefence {
		return "MonsterMagicDefence"
	}
	if e == LDAttrType_ARMOR {
		return "ARMOR"
	}
	if e == LDAttrType_BossMaxHP {
		return "BossMaxHP"
	}
	if e == LDAttrType_MonsterMaxHp {
		return "MonsterMaxHp"
	}
	if e == LDAttrType_EpicAttackAdd {
		return "EpicAttackAdd"
	}
	if e == LDAttrType_AvoidDamageParam {
		return "AvoidDamageParam"
	}
	if e == LDAttrType_MagicArmor {
		return "MagicArmor"
	}
	if e == LDAttrType_AvoidMagicParam {
		return "AvoidMagicParam"
	}
	if e == LDAttrType_CNT {
		return "CNT"
	}
	if e == LDAttrType_MAXHPPERCENT {
		return "MAXHPPERCENT"
	}
	if e == LDAttrType_ATTACKPERCENT {
		return "ATTACKPERCENT"
	}
	if e == LDAttrType_MAGICPERCENT {
		return "MAGICPERCENT"
	}
	if e == LDAttrType_DEFENCEPERCENT {
		return "DEFENCEPERCENT"
	}
	if e == LDAttrType_MAGICDEFENCEPERCENT {
		return "MAGICDEFENCEPERCENT"
	}
	if e == LDAttrType_MOVESPEEDPERCENT {
		return "MOVESPEEDPERCENT"
	}
	if e == LDAttrType_ATTACKSPEEDPERCENT {
		return "ATTACKSPEEDPERCENT"
	}
	if e == LDAttrType_Extra {
		return "Extra"
	}
	if e == LDAttrType_ATTACKWARRIOR {
		return "ATTACKWARRIOR"
	}
	if e == LDAttrType_ATTACKMAGE {
		return "ATTACKMAGE"
	}
	if e == LDAttrType_ATTACKPRIEST {
		return "ATTACKPRIEST"
	}
	if e == LDAttrType_ATTACKKNIGHT {
		return "ATTACKKNIGHT"
	}
	if e == LDAttrType_ATTACKARCHER {
		return "ATTACKARCHER"
	}
	if e == LDAttrType_ATTACKASSASIN {
		return "ATTACKASSASIN"
	}
	if e == LDAttrType_ATTACKSAVAGE {
		return "ATTACKSAVAGE"
	}
	if e == LDAttrType_ATTACKMONK {
		return "ATTACKMONK"
	}
	if e == LDAttrType_ATTACKWARKLOCK {
		return "ATTACKWARKLOCK"
	}
	if e == LDAttrType_ATTACKNEO {
		return "ATTACKNEO"
	}
	if e == LDAttrType_ATTACKENCHATER {
		return "ATTACKENCHATER"
	}
	if e == LDAttrType_ATTACKSACRIFICE {
		return "ATTACKSACRIFICE"
	}
	if e == LDAttrType_ATTACKBARD {
		return "ATTACKBARD"
	}
	if e == LDAttrType_ATTACKSHAMAN {
		return "ATTACKSHAMAN"
	}
	if e == LDAttrType_ATTACKDRUID {
		return "ATTACKDRUID"
	}
	if e == LDAttrType_ATTACKPALADIN {
		return "ATTACKPALADIN"
	}
	if e == LDAttrType_ATTACKDEADKNIGHT {
		return "ATTACKDEADKNIGHT"
	}
	if e == LDAttrType_ATTACKHUNTER {
		return "ATTACKHUNTER"
	}
	if e == LDAttrType_ATTACKRANGER {
		return "ATTACKRANGER"
	}
	if e == LDAttrType_ATTACKSTEAL {
		return "ATTACKSTEAL"
	}
	if e == LDAttrType_ATTACKDEMONHUNTER {
		return "ATTACKDEMONHUNTER"
	}
	if e == LDAttrType_ATTACKHUMAN {
		return "ATTACKHUMAN"
	}
	if e == LDAttrType_ATTACKORC {
		return "ATTACKORC"
	}
	if e == LDAttrType_ATTACKTAUREN {
		return "ATTACKTAUREN"
	}
	if e == LDAttrType_ATTACKUNDEAD {
		return "ATTACKUNDEAD"
	}
	if e == LDAttrType_ATTACKDWARF {
		return "ATTACKDWARF"
	}
	if e == LDAttrType_ATTACKELF {
		return "ATTACKELF"
	}
	if e == LDAttrType_ATTACKTROLL {
		return "ATTACKTROLL"
	}
	if e == LDAttrType_ATTACKGNOME {
		return "ATTACKGNOME"
	}
	if e == LDAttrType_ATTACKPANDA {
		return "ATTACKPANDA"
	}
	if e == LDAttrType_ATTACKGOD {
		return "ATTACKGOD"
	}
	if e == LDAttrType_ATTACKDEMON {
		return "ATTACKDEMON"
	}
	if e == LDAttrType_ATTACKWORGEN {
		return "ATTACKWORGEN"
	}
	if e == LDAttrType_ATTACKNAGA {
		return "ATTACKNAGA"
	}
	if e == LDAttrType_ATTACKBIRDMAN {
		return "ATTACKBIRDMAN"
	}
	if e == LDAttrType_ATTACKLIONMAN {
		return "ATTACKLIONMAN"
	}
	if e == LDAttrType_ATTACKPIGMAN {
		return "ATTACKPIGMAN"
	}
	if e == LDAttrType_ATTACKGOLEM {
		return "ATTACKGOLEM"
	}
	if e == LDAttrType_ATTACKWOOD {
		return "ATTACKWOOD"
	}
	if e == LDAttrType_ATTACKPUPPET {
		return "ATTACKPUPPET"
	}
	if e == LDAttrType_ATTACKCENTAUR {
		return "ATTACKCENTAUR"
	}
	if e == LDAttrType_ATTACKSHEEPMAN {
		return "ATTACKSHEEPMAN"
	}
	if e == LDAttrType_ATTACKDOGMAN {
		return "ATTACKDOGMAN"
	}
	if e == LDAttrType_ATTACKBEARMAN {
		return "ATTACKBEARMAN"
	}
	if e == LDAttrType_ATTACKMONKEYMAN {
		return "ATTACKMONKEYMAN"
	}
	if e == LDAttrType_ATTACKCATMAN {
		return "ATTACKCATMAN"
	}
	if e == LDAttrType_ATTACKVANISHMAN {
		return "ATTACKVANISHMAN"
	}
	if e == LDAttrType_ATTACKUNKOWNMAN {
		return "ATTACKUNKOWNMAN"
	}
	if e == LDAttrType_ATTACKMONSTERONE {
		return "ATTACKMONSTERONE"
	}
	if e == LDAttrType_ATTACKMONSTERTWO {
		return "ATTACKMONSTERTWO"
	}
	if e == LDAttrType_ATTACKMONSTERTREE {
		return "ATTACKMONSTERTREE"
	}
	if e == LDAttrType_ATTACKMONSTERFOUR {
		return "ATTACKMONSTERFOUR"
	}
	if e == LDAttrType_ATTACKMONSTERFIVE {
		return "ATTACKMONSTERFIVE"
	}
	if e == LDAttrType_HUMAADD {
		return "HUMAADD"
	}
	if e == LDAttrType_ORCADD {
		return "ORCADD"
	}
	if e == LDAttrType_TAUREADD {
		return "TAUREADD"
	}
	if e == LDAttrType_UNDEAADD {
		return "UNDEAADD"
	}
	if e == LDAttrType_DWARADD {
		return "DWARADD"
	}
	if e == LDAttrType_ELFADD {
		return "ELFADD"
	}
	if e == LDAttrType_TROLADD {
		return "TROLADD"
	}
	if e == LDAttrType_GNOMADD {
		return "GNOMADD"
	}
	if e == LDAttrType_PANDADD {
		return "PANDADD"
	}
	if e == LDAttrType_GODADD {
		return "GODADD"
	}
	if e == LDAttrType_DEMOADD {
		return "DEMOADD"
	}
	if e == LDAttrType_WORGEADD {
		return "WORGEADD"
	}
	if e == LDAttrType_NAGADD {
		return "NAGADD"
	}
	if e == LDAttrType_BIRDMAADD {
		return "BIRDMAADD"
	}
	if e == LDAttrType_LIONMAADD {
		return "LIONMAADD"
	}
	if e == LDAttrType_PIGMAADD {
		return "PIGMAADD"
	}
	if e == LDAttrType_GOLEADD {
		return "GOLEADD"
	}
	if e == LDAttrType_WOOADD {
		return "WOOADD"
	}
	if e == LDAttrType_PUPPEADD {
		return "PUPPEADD"
	}
	if e == LDAttrType_CENTAUADD {
		return "CENTAUADD"
	}
	if e == LDAttrType_SHEEPMAADD {
		return "SHEEPMAADD"
	}
	if e == LDAttrType_DOGMAADD {
		return "DOGMAADD"
	}
	if e == LDAttrType_BEARMAADD {
		return "BEARMAADD"
	}
	if e == LDAttrType_MONKEYMAADD {
		return "MONKEYMAADD"
	}
	if e == LDAttrType_CATMAADD {
		return "CATMAADD"
	}
	if e == LDAttrType_VANISHMAADD {
		return "VANISHMAADD"
	}
	if e == LDAttrType_UNKOWNMAADD {
		return "UNKOWNMAADD"
	}
	if e == LDAttrType_WARRIORADD {
		return "WARRIORADD"
	}
	if e == LDAttrType_MAGEADD {
		return "MAGEADD"
	}
	if e == LDAttrType_PRIESTADD {
		return "PRIESTADD"
	}
	if e == LDAttrType_KNIGHTADD {
		return "KNIGHTADD"
	}
	if e == LDAttrType_ARCHERADD {
		return "ARCHERADD"
	}
	if e == LDAttrType_ASSASINADD {
		return "ASSASINADD"
	}
	if e == LDAttrType_SAVAGEADD {
		return "SAVAGEADD"
	}
	if e == LDAttrType_MONKADD {
		return "MONKADD"
	}
	if e == LDAttrType_WARKLOCKADD {
		return "WARKLOCKADD"
	}
	if e == LDAttrType_NEOADD {
		return "NEOADD"
	}
	if e == LDAttrType_ENCHATERADD {
		return "ENCHATERADD"
	}
	if e == LDAttrType_SACRIFICEADD {
		return "SACRIFICEADD"
	}
	if e == LDAttrType_BARDADD {
		return "BARDADD"
	}
	if e == LDAttrType_SHAMANADD {
		return "SHAMANADD"
	}
	if e == LDAttrType_DRUIDADD {
		return "DRUIDADD"
	}
	if e == LDAttrType_PALADINADD {
		return "PALADINADD"
	}
	if e == LDAttrType_DEADKNIGHTADD {
		return "DEADKNIGHTADD"
	}
	if e == LDAttrType_HUNTERADD {
		return "HUNTERADD"
	}
	if e == LDAttrType_RANGERADD {
		return "RANGERADD"
	}
	if e == LDAttrType_STEALADD {
		return "STEALADD"
	}
	if e == LDAttrType_DEMONHUNTERADD {
		return "DEMONHUNTERADD"
	}
	if e == LDAttrType_ADDLEVEL {
		return "ADDLEVEL"
	}
	if e == LDAttrType_ExtraCNT {
		return "ExtraCNT"
	}

	return ""

}

func GetLDAttrType(tp string) (LDAttrType, bool) {
	switch tp {
	case "HP":
		return LDAttrType_HP, true
	case "MAXHP":
		return LDAttrType_MAXHP, true
	case "ATTACK":
		return LDAttrType_ATTACK, true
	case "MAGIC":
		return LDAttrType_MAGIC, true
	case "DEFENCE":
		return LDAttrType_DEFENCE, true
	case "MAGICDEFENCE":
		return LDAttrType_MAGICDEFENCE, true
	case "IGNDEFENCE":
		return LDAttrType_IGNDEFENCE, true
	case "IGNMAGICDEFENCE":
		return LDAttrType_IGNMAGICDEFENCE, true
	case "IGNHEAL":
		return LDAttrType_IGNHEAL, true
	case "CRITAL":
		return LDAttrType_CRITAL, true
	case "IGNCRITAL":
		return LDAttrType_IGNCRITAL, true
	case "DODGE":
		return LDAttrType_DODGE, true
	case "IGNDODGE":
		return LDAttrType_IGNDODGE, true
	case "TRUEHARM":
		return LDAttrType_TRUEHARM, true
	case "IGNTRUEHARM":
		return LDAttrType_IGNTRUEHARM, true
	case "FIREHARM":
		return LDAttrType_FIREHARM, true
	case "IGNFIREHARM":
		return LDAttrType_IGNFIREHARM, true
	case "MOVESPEED":
		return LDAttrType_MOVESPEED, true
	case "ATTACKSPEED":
		return LDAttrType_ATTACKSPEED, true
	case "HITRATE":
		return LDAttrType_HITRATE, true
	case "BLASTS":
		return LDAttrType_BLASTS, true
	case "IGNBLASTS":
		return LDAttrType_IGNBLASTS, true
	case "MAXSP":
		return LDAttrType_MAXSP, true
	case "SP":
		return LDAttrType_SP, true
	case "DEFENCEIMMUNE":
		return LDAttrType_DEFENCEIMMUNE, true
	case "BLOCK":
		return LDAttrType_BLOCK, true
	case "DAMAGEADDPRECENT":
		return LDAttrType_DAMAGEADDPRECENT, true
	case "DAMAGEADD":
		return LDAttrType_DAMAGEADD, true
	case "PENETRATE":
		return LDAttrType_PENETRATE, true
	case "TRICKADD":
		return LDAttrType_TRICKADD, true
	case "BUFFADD":
		return LDAttrType_BUFFADD, true
	case "SKILLTRIGGERADD":
		return LDAttrType_SKILLTRIGGERADD, true
	case "DEFENCEBLOCK":
		return LDAttrType_DEFENCEBLOCK, true
	case "DEFDAMAGEADDPRECENT":
		return LDAttrType_DEFDAMAGEADDPRECENT, true
	case "DEFDAMAGEADD":
		return LDAttrType_DEFDAMAGEADD, true
	case "DEFPENETRATE":
		return LDAttrType_DEFPENETRATE, true
	case "DEFTRICKADD":
		return LDAttrType_DEFTRICKADD, true
	case "DEFBUFFADD":
		return LDAttrType_DEFBUFFADD, true
	case "DEFSKILLTRIGGERADD":
		return LDAttrType_DEFSKILLTRIGGERADD, true
	case "BLOCKSUBPRECENT":
		return LDAttrType_BLOCKSUBPRECENT, true
	case "DEFSTUN":
		return LDAttrType_DEFSTUN, true
	case "DEFFREEZE":
		return LDAttrType_DEFFREEZE, true
	case "DEFSILENT":
		return LDAttrType_DEFSILENT, true
	case "DEFBLOOD":
		return LDAttrType_DEFBLOOD, true
	case "DEFPOISON":
		return LDAttrType_DEFPOISON, true
	case "DEFBURN":
		return LDAttrType_DEFBURN, true
	case "SPPRECENTADD":
		return LDAttrType_SPPRECENTADD, true
	case "SPADD":
		return LDAttrType_SPADD, true
	case "MonsterDefence":
		return LDAttrType_MonsterDefence, true
	case "MonsterSpeedDown":
		return LDAttrType_MonsterSpeedDown, true
	case "MonsterDizzyDefence":
		return LDAttrType_MonsterDizzyDefence, true
	case "HeroAttackSpeedAdd":
		return LDAttrType_HeroAttackSpeedAdd, true
	case "SpAdd":
		return LDAttrType_SpAdd, true
	case "HeroMoveSpeedAdd":
		return LDAttrType_HeroMoveSpeedAdd, true
	case "PeopleMaxCount":
		return LDAttrType_PeopleMaxCount, true
	case "HeroAttackAdd":
		return LDAttrType_HeroAttackAdd, true
	case "BossDamageAdd":
		return LDAttrType_BossDamageAdd, true
	case "PhysicalDamageAdd":
		return LDAttrType_PhysicalDamageAdd, true
	case "SkillDamageAdd":
		return LDAttrType_SkillDamageAdd, true
	case "CritalAdd":
		return LDAttrType_CritalAdd, true
	case "BlastAdd":
		return LDAttrType_BlastAdd, true
	case "SkillUseRate":
		return LDAttrType_SkillUseRate, true
	case "StunTimeAdd":
		return LDAttrType_StunTimeAdd, true
	case "StunDamageAdd":
		return LDAttrType_StunDamageAdd, true
	case "MagicDamageAdd":
		return LDAttrType_MagicDamageAdd, true
	case "MonsterMagicDefence":
		return LDAttrType_MonsterMagicDefence, true
	case "ARMOR":
		return LDAttrType_ARMOR, true
	case "BossMaxHP":
		return LDAttrType_BossMaxHP, true
	case "MonsterMaxHp":
		return LDAttrType_MonsterMaxHp, true
	case "EpicAttackAdd":
		return LDAttrType_EpicAttackAdd, true
	case "AvoidDamageParam":
		return LDAttrType_AvoidDamageParam, true
	case "MagicArmor":
		return LDAttrType_MagicArmor, true
	case "AvoidMagicParam":
		return LDAttrType_AvoidMagicParam, true
	case "CNT":
		return LDAttrType_CNT, true
	case "MAXHPPERCENT":
		return LDAttrType_MAXHPPERCENT, true
	case "ATTACKPERCENT":
		return LDAttrType_ATTACKPERCENT, true
	case "MAGICPERCENT":
		return LDAttrType_MAGICPERCENT, true
	case "DEFENCEPERCENT":
		return LDAttrType_DEFENCEPERCENT, true
	case "MAGICDEFENCEPERCENT":
		return LDAttrType_MAGICDEFENCEPERCENT, true
	case "MOVESPEEDPERCENT":
		return LDAttrType_MOVESPEEDPERCENT, true
	case "ATTACKSPEEDPERCENT":
		return LDAttrType_ATTACKSPEEDPERCENT, true
	case "Extra":
		return LDAttrType_Extra, true
	case "ATTACKWARRIOR":
		return LDAttrType_ATTACKWARRIOR, true
	case "ATTACKMAGE":
		return LDAttrType_ATTACKMAGE, true
	case "ATTACKPRIEST":
		return LDAttrType_ATTACKPRIEST, true
	case "ATTACKKNIGHT":
		return LDAttrType_ATTACKKNIGHT, true
	case "ATTACKARCHER":
		return LDAttrType_ATTACKARCHER, true
	case "ATTACKASSASIN":
		return LDAttrType_ATTACKASSASIN, true
	case "ATTACKSAVAGE":
		return LDAttrType_ATTACKSAVAGE, true
	case "ATTACKMONK":
		return LDAttrType_ATTACKMONK, true
	case "ATTACKWARKLOCK":
		return LDAttrType_ATTACKWARKLOCK, true
	case "ATTACKNEO":
		return LDAttrType_ATTACKNEO, true
	case "ATTACKENCHATER":
		return LDAttrType_ATTACKENCHATER, true
	case "ATTACKSACRIFICE":
		return LDAttrType_ATTACKSACRIFICE, true
	case "ATTACKBARD":
		return LDAttrType_ATTACKBARD, true
	case "ATTACKSHAMAN":
		return LDAttrType_ATTACKSHAMAN, true
	case "ATTACKDRUID":
		return LDAttrType_ATTACKDRUID, true
	case "ATTACKPALADIN":
		return LDAttrType_ATTACKPALADIN, true
	case "ATTACKDEADKNIGHT":
		return LDAttrType_ATTACKDEADKNIGHT, true
	case "ATTACKHUNTER":
		return LDAttrType_ATTACKHUNTER, true
	case "ATTACKRANGER":
		return LDAttrType_ATTACKRANGER, true
	case "ATTACKSTEAL":
		return LDAttrType_ATTACKSTEAL, true
	case "ATTACKDEMONHUNTER":
		return LDAttrType_ATTACKDEMONHUNTER, true
	case "ATTACKHUMAN":
		return LDAttrType_ATTACKHUMAN, true
	case "ATTACKORC":
		return LDAttrType_ATTACKORC, true
	case "ATTACKTAUREN":
		return LDAttrType_ATTACKTAUREN, true
	case "ATTACKUNDEAD":
		return LDAttrType_ATTACKUNDEAD, true
	case "ATTACKDWARF":
		return LDAttrType_ATTACKDWARF, true
	case "ATTACKELF":
		return LDAttrType_ATTACKELF, true
	case "ATTACKTROLL":
		return LDAttrType_ATTACKTROLL, true
	case "ATTACKGNOME":
		return LDAttrType_ATTACKGNOME, true
	case "ATTACKPANDA":
		return LDAttrType_ATTACKPANDA, true
	case "ATTACKGOD":
		return LDAttrType_ATTACKGOD, true
	case "ATTACKDEMON":
		return LDAttrType_ATTACKDEMON, true
	case "ATTACKWORGEN":
		return LDAttrType_ATTACKWORGEN, true
	case "ATTACKNAGA":
		return LDAttrType_ATTACKNAGA, true
	case "ATTACKBIRDMAN":
		return LDAttrType_ATTACKBIRDMAN, true
	case "ATTACKLIONMAN":
		return LDAttrType_ATTACKLIONMAN, true
	case "ATTACKPIGMAN":
		return LDAttrType_ATTACKPIGMAN, true
	case "ATTACKGOLEM":
		return LDAttrType_ATTACKGOLEM, true
	case "ATTACKWOOD":
		return LDAttrType_ATTACKWOOD, true
	case "ATTACKPUPPET":
		return LDAttrType_ATTACKPUPPET, true
	case "ATTACKCENTAUR":
		return LDAttrType_ATTACKCENTAUR, true
	case "ATTACKSHEEPMAN":
		return LDAttrType_ATTACKSHEEPMAN, true
	case "ATTACKDOGMAN":
		return LDAttrType_ATTACKDOGMAN, true
	case "ATTACKBEARMAN":
		return LDAttrType_ATTACKBEARMAN, true
	case "ATTACKMONKEYMAN":
		return LDAttrType_ATTACKMONKEYMAN, true
	case "ATTACKCATMAN":
		return LDAttrType_ATTACKCATMAN, true
	case "ATTACKVANISHMAN":
		return LDAttrType_ATTACKVANISHMAN, true
	case "ATTACKUNKOWNMAN":
		return LDAttrType_ATTACKUNKOWNMAN, true
	case "ATTACKMONSTERONE":
		return LDAttrType_ATTACKMONSTERONE, true
	case "ATTACKMONSTERTWO":
		return LDAttrType_ATTACKMONSTERTWO, true
	case "ATTACKMONSTERTREE":
		return LDAttrType_ATTACKMONSTERTREE, true
	case "ATTACKMONSTERFOUR":
		return LDAttrType_ATTACKMONSTERFOUR, true
	case "ATTACKMONSTERFIVE":
		return LDAttrType_ATTACKMONSTERFIVE, true
	case "HUMAADD":
		return LDAttrType_HUMAADD, true
	case "ORCADD":
		return LDAttrType_ORCADD, true
	case "TAUREADD":
		return LDAttrType_TAUREADD, true
	case "UNDEAADD":
		return LDAttrType_UNDEAADD, true
	case "DWARADD":
		return LDAttrType_DWARADD, true
	case "ELFADD":
		return LDAttrType_ELFADD, true
	case "TROLADD":
		return LDAttrType_TROLADD, true
	case "GNOMADD":
		return LDAttrType_GNOMADD, true
	case "PANDADD":
		return LDAttrType_PANDADD, true
	case "GODADD":
		return LDAttrType_GODADD, true
	case "DEMOADD":
		return LDAttrType_DEMOADD, true
	case "WORGEADD":
		return LDAttrType_WORGEADD, true
	case "NAGADD":
		return LDAttrType_NAGADD, true
	case "BIRDMAADD":
		return LDAttrType_BIRDMAADD, true
	case "LIONMAADD":
		return LDAttrType_LIONMAADD, true
	case "PIGMAADD":
		return LDAttrType_PIGMAADD, true
	case "GOLEADD":
		return LDAttrType_GOLEADD, true
	case "WOOADD":
		return LDAttrType_WOOADD, true
	case "PUPPEADD":
		return LDAttrType_PUPPEADD, true
	case "CENTAUADD":
		return LDAttrType_CENTAUADD, true
	case "SHEEPMAADD":
		return LDAttrType_SHEEPMAADD, true
	case "DOGMAADD":
		return LDAttrType_DOGMAADD, true
	case "BEARMAADD":
		return LDAttrType_BEARMAADD, true
	case "MONKEYMAADD":
		return LDAttrType_MONKEYMAADD, true
	case "CATMAADD":
		return LDAttrType_CATMAADD, true
	case "VANISHMAADD":
		return LDAttrType_VANISHMAADD, true
	case "UNKOWNMAADD":
		return LDAttrType_UNKOWNMAADD, true
	case "WARRIORADD":
		return LDAttrType_WARRIORADD, true
	case "MAGEADD":
		return LDAttrType_MAGEADD, true
	case "PRIESTADD":
		return LDAttrType_PRIESTADD, true
	case "KNIGHTADD":
		return LDAttrType_KNIGHTADD, true
	case "ARCHERADD":
		return LDAttrType_ARCHERADD, true
	case "ASSASINADD":
		return LDAttrType_ASSASINADD, true
	case "SAVAGEADD":
		return LDAttrType_SAVAGEADD, true
	case "MONKADD":
		return LDAttrType_MONKADD, true
	case "WARKLOCKADD":
		return LDAttrType_WARKLOCKADD, true
	case "NEOADD":
		return LDAttrType_NEOADD, true
	case "ENCHATERADD":
		return LDAttrType_ENCHATERADD, true
	case "SACRIFICEADD":
		return LDAttrType_SACRIFICEADD, true
	case "BARDADD":
		return LDAttrType_BARDADD, true
	case "SHAMANADD":
		return LDAttrType_SHAMANADD, true
	case "DRUIDADD":
		return LDAttrType_DRUIDADD, true
	case "PALADINADD":
		return LDAttrType_PALADINADD, true
	case "DEADKNIGHTADD":
		return LDAttrType_DEADKNIGHTADD, true
	case "HUNTERADD":
		return LDAttrType_HUNTERADD, true
	case "RANGERADD":
		return LDAttrType_RANGERADD, true
	case "STEALADD":
		return LDAttrType_STEALADD, true
	case "DEMONHUNTERADD":
		return LDAttrType_DEMONHUNTERADD, true
	case "ADDLEVEL":
		return LDAttrType_ADDLEVEL, true
	case "ExtraCNT":
		return LDAttrType_ExtraCNT, true

	}
	return -1, false

}

func GetLDAttrType_ByAlias(tp string) (LDAttrType, bool) {
	switch tp {
	case "血量":
		return LDAttrType_HP, true
	case "血量上限":
		return LDAttrType_MAXHP, true
	case "攻击力":
		return LDAttrType_ATTACK, true
	case "魔法攻击力":
		return LDAttrType_MAGIC, true
	case "防御力":
		return LDAttrType_DEFENCE, true
	case "法术防御":
		return LDAttrType_MAGICDEFENCE, true
	case "忽视防御":
		return LDAttrType_IGNDEFENCE, true
	case "忽视法防":
		return LDAttrType_IGNMAGICDEFENCE, true
	case "抗疗":
		return LDAttrType_IGNHEAL, true
	case "暴击率":
		return LDAttrType_CRITAL, true
	case "忽视暴击":
		return LDAttrType_IGNCRITAL, true
	case "闪避率":
		return LDAttrType_DODGE, true
	case "忽视闪避":
		return LDAttrType_IGNDODGE, true
	case "真实伤害":
		return LDAttrType_TRUEHARM, true
	case "忽视真伤":
		return LDAttrType_IGNTRUEHARM, true
	case "火系伤害":
		return LDAttrType_FIREHARM, true
	case "忽视火伤":
		return LDAttrType_IGNFIREHARM, true
	case "移动速度":
		return LDAttrType_MOVESPEED, true
	case "攻击速度":
		return LDAttrType_ATTACKSPEED, true
	case "命中率":
		return LDAttrType_HITRATE, true
	case "暴伤":
		return LDAttrType_BLASTS, true
	case "忽视暴伤":
		return LDAttrType_IGNBLASTS, true
	case "最大精神力":
		return LDAttrType_MAXSP, true
	case "精神力":
		return LDAttrType_SP, true
	case "防御免伤比":
		return LDAttrType_DEFENCEIMMUNE, true
	case "格挡率":
		return LDAttrType_BLOCK, true
	case "伤害百分比加成":
		return LDAttrType_DAMAGEADDPRECENT, true
	case "伤害固定加成":
		return LDAttrType_DAMAGEADD, true
	case "穿透率":
		return LDAttrType_PENETRATE, true
	case "大招增伤":
		return LDAttrType_TRICKADD, true
	case "状态加强":
		return LDAttrType_BUFFADD, true
	case "技能触发率增长":
		return LDAttrType_SKILLTRIGGERADD, true
	case "忽视格挡":
		return LDAttrType_DEFENCEBLOCK, true
	case "百分比减伤":
		return LDAttrType_DEFDAMAGEADDPRECENT, true
	case "固定减伤":
		return LDAttrType_DEFDAMAGEADD, true
	case "忽视穿透":
		return LDAttrType_DEFPENETRATE, true
	case "大招减伤":
		return LDAttrType_DEFTRICKADD, true
	case "状态抵抗":
		return LDAttrType_DEFBUFFADD, true
	case "技能触发率抵抗":
		return LDAttrType_DEFSKILLTRIGGERADD, true
	case "格挡免伤比例":
		return LDAttrType_BLOCKSUBPRECENT, true
	case "眩晕抵抗":
		return LDAttrType_DEFSTUN, true
	case "冰冻抵抗":
		return LDAttrType_DEFFREEZE, true
	case "沉默抵抗":
		return LDAttrType_DEFSILENT, true
	case "流血抵抗":
		return LDAttrType_DEFBLOOD, true
	case "中毒抵抗":
		return LDAttrType_DEFPOISON, true
	case "灼烧抵抗":
		return LDAttrType_DEFBURN, true
	case "精神力百分比加成":
		return LDAttrType_SPPRECENTADD, true
	case "精神力固定加成":
		return LDAttrType_SPADD, true
	case "怪物防御力":
		return LDAttrType_MonsterDefence, true
	case "怪物移速":
		return LDAttrType_MonsterSpeedDown, true
	case "怪物眩晕抵抗":
		return LDAttrType_MonsterDizzyDefence, true
	case "英雄攻击速度":
		return LDAttrType_HeroAttackSpeedAdd, true
	case "英雄sp恢复速度":
		return LDAttrType_SpAdd, true
	case "英雄移动速度":
		return LDAttrType_HeroMoveSpeedAdd, true
	case "兵种上限":
		return LDAttrType_PeopleMaxCount, true
	case "英雄攻击力增加":
		return LDAttrType_HeroAttackAdd, true
	case "首领增伤":
		return LDAttrType_BossDamageAdd, true
	case "物理增伤":
		return LDAttrType_PhysicalDamageAdd, true
	case "技能增伤":
		return LDAttrType_SkillDamageAdd, true
	case "暴击率增加":
		return LDAttrType_CritalAdd, true
	case "爆伤加成":
		return LDAttrType_BlastAdd, true
	case "技能触发率":
		return LDAttrType_SkillUseRate, true
	case "眩晕效果加成":
		return LDAttrType_StunTimeAdd, true
	case "眩晕伤害加成":
		return LDAttrType_StunDamageAdd, true
	case "魔法伤害加成":
		return LDAttrType_MagicDamageAdd, true
	case "怪物魔法抗性":
		return LDAttrType_MonsterMagicDefence, true
	case "护甲值":
		return LDAttrType_ARMOR, true
	case "Boss最大血量加成":
		return LDAttrType_BossMaxHP, true
	case "小怪最大血量加成":
		return LDAttrType_MonsterMaxHp, true
	case "神话英雄攻击加成":
		return LDAttrType_EpicAttackAdd, true
	case "免伤系数":
		return LDAttrType_AvoidDamageParam, true
	case "魔法护甲":
		return LDAttrType_MagicArmor, true
	case "魔免系数":
		return LDAttrType_AvoidMagicParam, true
	case "结束标志":
		return LDAttrType_CNT, true
	case "血量上限百分比":
		return LDAttrType_MAXHPPERCENT, true
	case "攻击力百分比":
		return LDAttrType_ATTACKPERCENT, true
	case "魔法攻击力百分比":
		return LDAttrType_MAGICPERCENT, true
	case "防御力百分比":
		return LDAttrType_DEFENCEPERCENT, true
	case "法术防御力百分比":
		return LDAttrType_MAGICDEFENCEPERCENT, true
	case "移动速度百分比":
		return LDAttrType_MOVESPEEDPERCENT, true
	case "攻击速度百分比":
		return LDAttrType_ATTACKSPEEDPERCENT, true
	case "额外战斗属性开始":
		return LDAttrType_Extra, true
	case "对战士加成":
		return LDAttrType_ATTACKWARRIOR, true
	case "对法师加成":
		return LDAttrType_ATTACKMAGE, true
	case "对牧师加成":
		return LDAttrType_ATTACKPRIEST, true
	case "对骑士加成":
		return LDAttrType_ATTACKKNIGHT, true
	case "对射手加成":
		return LDAttrType_ATTACKARCHER, true
	case "对刺客加成":
		return LDAttrType_ATTACKASSASIN, true
	case "对野蛮人加成":
		return LDAttrType_ATTACKSAVAGE, true
	case "对武僧加成":
		return LDAttrType_ATTACKMONK, true
	case "对术士加成":
		return LDAttrType_ATTACKWARKLOCK, true
	case "对亡灵法师加成":
		return LDAttrType_ATTACKNEO, true
	case "对幻术师加成":
		return LDAttrType_ATTACKENCHATER, true
	case "对祭祀加成":
		return LDAttrType_ATTACKSACRIFICE, true
	case "对吟游诗人加成":
		return LDAttrType_ATTACKBARD, true
	case "对萨满加成":
		return LDAttrType_ATTACKSHAMAN, true
	case "对德鲁伊加成":
		return LDAttrType_ATTACKDRUID, true
	case "对圣骑士加成":
		return LDAttrType_ATTACKPALADIN, true
	case "对死骑加成":
		return LDAttrType_ATTACKDEADKNIGHT, true
	case "对猎人加成":
		return LDAttrType_ATTACKHUNTER, true
	case "对游侠加成":
		return LDAttrType_ATTACKRANGER, true
	case "对潜行者加成":
		return LDAttrType_ATTACKSTEAL, true
	case "对恶魔猎手加成":
		return LDAttrType_ATTACKDEMONHUNTER, true
	case "对人类加成":
		return LDAttrType_ATTACKHUMAN, true
	case "对兽人加成":
		return LDAttrType_ATTACKORC, true
	case "对牛头人加成":
		return LDAttrType_ATTACKTAUREN, true
	case "对亡灵加成":
		return LDAttrType_ATTACKUNDEAD, true
	case "对矮人加成":
		return LDAttrType_ATTACKDWARF, true
	case "对精灵加成":
		return LDAttrType_ATTACKELF, true
	case "对巨魔加成":
		return LDAttrType_ATTACKTROLL, true
	case "对侏儒加成":
		return LDAttrType_ATTACKGNOME, true
	case "对熊猫人加成":
		return LDAttrType_ATTACKPANDA, true
	case "对神加成":
		return LDAttrType_ATTACKGOD, true
	case "对恶魔加成":
		return LDAttrType_ATTACKDEMON, true
	case "对狼人加成":
		return LDAttrType_ATTACKWORGEN, true
	case "对娜迦加成":
		return LDAttrType_ATTACKNAGA, true
	case "对鸟人加成":
		return LDAttrType_ATTACKBIRDMAN, true
	case "对狮身人加成":
		return LDAttrType_ATTACKLIONMAN, true
	case "对野猪人加成":
		return LDAttrType_ATTACKPIGMAN, true
	case "对山顶巨人加成":
		return LDAttrType_ATTACKGOLEM, true
	case "对树精加成":
		return LDAttrType_ATTACKWOOD, true
	case "对傀儡加成":
		return LDAttrType_ATTACKPUPPET, true
	case "对半人马加成":
		return LDAttrType_ATTACKCENTAUR, true
	case "对萨特加成":
		return LDAttrType_ATTACKSHEEPMAN, true
	case "对狗头人加成":
		return LDAttrType_ATTACKDOGMAN, true
	case "对熊人加成":
		return LDAttrType_ATTACKBEARMAN, true
	case "对猩猩人加成":
		return LDAttrType_ATTACKMONKEYMAN, true
	case "对猫人加成":
		return LDAttrType_ATTACKCATMAN, true
	case "对隐身人加成":
		return LDAttrType_ATTACKVANISHMAN, true
	case "对穿越者加成":
		return LDAttrType_ATTACKUNKOWNMAN, true
	case "对类型1怪物加成":
		return LDAttrType_ATTACKMONSTERONE, true
	case "对类型2怪物加成":
		return LDAttrType_ATTACKMONSTERTWO, true
	case "对类型3怪物加成":
		return LDAttrType_ATTACKMONSTERTREE, true
	case "对类型4怪物加成":
		return LDAttrType_ATTACKMONSTERFOUR, true
	case "对类型5怪物加成":
		return LDAttrType_ATTACKMONSTERFIVE, true
	case "人类加成":
		return LDAttrType_HUMAADD, true
	case "兽人加成":
		return LDAttrType_ORCADD, true
	case "牛头人加成":
		return LDAttrType_TAUREADD, true
	case "亡灵加成":
		return LDAttrType_UNDEAADD, true
	case "矮人加成":
		return LDAttrType_DWARADD, true
	case "精灵加成":
		return LDAttrType_ELFADD, true
	case "巨魔加成":
		return LDAttrType_TROLADD, true
	case "侏儒加成":
		return LDAttrType_GNOMADD, true
	case "熊猫人加成":
		return LDAttrType_PANDADD, true
	case "神加成":
		return LDAttrType_GODADD, true
	case "恶魔加成":
		return LDAttrType_DEMOADD, true
	case "狼人加成":
		return LDAttrType_WORGEADD, true
	case "娜迦加成":
		return LDAttrType_NAGADD, true
	case "鸟人加成":
		return LDAttrType_BIRDMAADD, true
	case "狮身人加成":
		return LDAttrType_LIONMAADD, true
	case "野猪人加成":
		return LDAttrType_PIGMAADD, true
	case "山顶巨人加成":
		return LDAttrType_GOLEADD, true
	case "树精加成":
		return LDAttrType_WOOADD, true
	case "傀儡加成":
		return LDAttrType_PUPPEADD, true
	case "半人马加成":
		return LDAttrType_CENTAUADD, true
	case "萨特加成":
		return LDAttrType_SHEEPMAADD, true
	case "狗头人加成":
		return LDAttrType_DOGMAADD, true
	case "熊人加成":
		return LDAttrType_BEARMAADD, true
	case "猩猩人加成":
		return LDAttrType_MONKEYMAADD, true
	case "猫人加成":
		return LDAttrType_CATMAADD, true
	case "隐身人加成":
		return LDAttrType_VANISHMAADD, true
	case "穿越者加成":
		return LDAttrType_UNKOWNMAADD, true
	case "战士加成":
		return LDAttrType_WARRIORADD, true
	case "法师加成":
		return LDAttrType_MAGEADD, true
	case "牧师加成":
		return LDAttrType_PRIESTADD, true
	case "骑士加成":
		return LDAttrType_KNIGHTADD, true
	case "射手加成":
		return LDAttrType_ARCHERADD, true
	case "刺客加成":
		return LDAttrType_ASSASINADD, true
	case "野蛮人加成":
		return LDAttrType_SAVAGEADD, true
	case "武僧加成":
		return LDAttrType_MONKADD, true
	case "术士加成":
		return LDAttrType_WARKLOCKADD, true
	case "亡灵法师加成":
		return LDAttrType_NEOADD, true
	case "幻术师加成":
		return LDAttrType_ENCHATERADD, true
	case "祭祀加成":
		return LDAttrType_SACRIFICEADD, true
	case "吟游诗人加成":
		return LDAttrType_BARDADD, true
	case "萨满加成":
		return LDAttrType_SHAMANADD, true
	case "德鲁伊加成":
		return LDAttrType_DRUIDADD, true
	case "圣骑士加成":
		return LDAttrType_PALADINADD, true
	case "死骑加成":
		return LDAttrType_DEADKNIGHTADD, true
	case "猎人加成":
		return LDAttrType_HUNTERADD, true
	case "游侠加成":
		return LDAttrType_RANGERADD, true
	case "潜行者加成":
		return LDAttrType_STEALADD, true
	case "恶魔猎手加成":
		return LDAttrType_DEMONHUNTERADD, true
	case "增加等级":
		return LDAttrType_ADDLEVEL, true
	case "额外属性结束":
		return LDAttrType_ExtraCNT, true

	}
	return -1, false

}

func (e LDAttrTitleType) String() string {

	if e == LDAttrTitleType_Basic {
		return "Basic"
	}
	if e == LDAttrTitleType_Advance {
		return "Advance"
	}

	return ""

}

func GetLDAttrTitleType(tp string) (LDAttrTitleType, bool) {
	switch tp {
	case "Basic":
		return LDAttrTitleType_Basic, true
	case "Advance":
		return LDAttrTitleType_Advance, true

	}
	return -1, false

}

func GetLDAttrTitleType_ByAlias(tp string) (LDAttrTitleType, bool) {
	switch tp {
	case "基础属性":
		return LDAttrTitleType_Basic, true
	case "进阶属性":
		return LDAttrTitleType_Advance, true

	}
	return -1, false

}
