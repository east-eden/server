package config

// 基础配置-拳皇/AttrEnum
type AttrType int

const (
	AttrType_BEGIN               AttrType = 1   //枚举开始
	AttrType_HP                  AttrType = 1   //血量
	AttrType_MAXHP               AttrType = 2   //最大血量
	AttrType_ATTACK              AttrType = 3   //攻击
	AttrType_DEFENCE             AttrType = 4   //防御
	AttrType_PENETRATE           AttrType = 5   //穿透率
	AttrType_ATTACKSPEED         AttrType = 6   //攻击速度
	AttrType_MOVESPEED           AttrType = 7   //移动速度
	AttrType_IGNDODGE            AttrType = 8   //闪避抵抗
	AttrType_DODGE               AttrType = 9   //闪避
	AttrType_CRITAL              AttrType = 10  //暴击率
	AttrType_IGNCRITAL           AttrType = 11  //免暴率
	AttrType_BLASTS              AttrType = 12  //暴击上限
	AttrType_IGNBLASTS           AttrType = 13  //暴击抵抗
	AttrType_PARRY               AttrType = 14  //格挡
	AttrType_BREAKPARRY          AttrType = 15  //格挡抵抗
	AttrType_STIFFNESS           AttrType = 16  //硬直
	AttrType_STIFFNESSDEF        AttrType = 17  //硬直抵抗
	AttrType_ENDURE              AttrType = 18  //霸体免控率
	AttrType_DEFSTUN             AttrType = 19  //眩晕抵抗
	AttrType_DEFFREEZE           AttrType = 20  //冰冻抵抗
	AttrType_PHYSICALADD         AttrType = 21  //物理增伤
	AttrType_PHYSICALSUB         AttrType = 22  //物理减伤
	AttrType_MAGICADD            AttrType = 23  //魔法增伤
	AttrType_MAGICSUB            AttrType = 24  //魔法减伤
	AttrType_NEARADD             AttrType = 25  //近战增伤
	AttrType_NEARSUB             AttrType = 26  //近战减伤
	AttrType_FARADD              AttrType = 27  //远程增伤
	AttrType_FARSUB              AttrType = 28  //远程减伤
	AttrType_TRICKADD            AttrType = 29  //大招增伤
	AttrType_DEFTRICKADD         AttrType = 30  //大招减伤
	AttrType_DAMAGEADD           AttrType = 31  //最终增伤
	AttrType_DAMAGESUB           AttrType = 32  //最终减伤
	AttrType_SPINCREASE          AttrType = 33  //怒气增加
	AttrType_SPRECOVERY          AttrType = 34  //怒气回复
	AttrType_SPINCRERESPOSTIVE   AttrType = 35  //攻击回怒
	AttrType_SPINCRERESNEGIATIVE AttrType = 36  //受击回怒
	AttrType_SPINCRERESNEFAILED  AttrType = 37  //战败回怒
	AttrType_SPRETENTION         AttrType = 38  //怒气截留
	AttrType_SPDESTRUCTION       AttrType = 39  //自爆减怒
	AttrType_EFFECTHEAL          AttrType = 40  //被治疗效果
	AttrType_SPINITIAL           AttrType = 41  //初始怒气
	AttrType_MAXSP               AttrType = 42  //最大怒气
	AttrType_SP                  AttrType = 43  //怒气
	AttrType_SPINCRERESADD       AttrType = 44  //回怒增幅
	AttrType_SPINCRERESSUB       AttrType = 45  //回怒降幅
	AttrType_SPKILLRECOVER       AttrType = 46  //击杀回怒
	AttrType_SPINITRECOVER       AttrType = 47  //初始怒气回复
	AttrType_SILENTDEFENCE       AttrType = 48  //沉默抵抗
	AttrType_SKILLTRIGGERADD     AttrType = 49  //技能触发率增长
	AttrType_DEFSKILLTRIGGERADD  AttrType = 50  //技能触发率抵抗
	AttrType_HEALADD             AttrType = 51  //治疗增幅
	AttrType_COUNTERING          AttrType = 52  //反击率
	AttrType_DEFCOUNTER          AttrType = 53  //反击抵抗
	AttrType_BLOODSUCK           AttrType = 54  //吸血系数
	AttrType_DEFBLOODSUCK        AttrType = 55  //吸血抗性
	AttrType_SPRESERVES          AttrType = 56  //怒气储量
	AttrType_COMBO               AttrType = 57  //连击
	AttrType_DEFCOMBO            AttrType = 58  //连击抗性
	AttrType_ENDUREBREAK         AttrType = 59  //霸体免控抵抗
	AttrType_COUNTERADD          AttrType = 60  //反击增幅
	AttrType_HEALSUB             AttrType = 61  //治疗抑制
	AttrType_ALLDEFENCEADD       AttrType = 62  //全状态抗性
	AttrType_ALLSTATEADD         AttrType = 63  //全状态属性
	AttrType_MAXDEFENCEADD       AttrType = 64  //最高状态抗性
	AttrType_MAXSTATEADD         AttrType = 65  //最高状态属性
	AttrType_ROLEFIRSTRATE       AttrType = 66  //角色先手值
	AttrType_ATTACKRANGE         AttrType = 67  //攻击距离
	AttrType_STIFFNESSSTRENG     AttrType = 68  //硬直增强
	AttrType_STIFFNESSWEAKEN     AttrType = 69  //硬直削弱
	AttrType_SUCKSTRENGTH        AttrType = 70  //吸血增幅
	AttrType_HITRATE             AttrType = 71  //命中率
	AttrType_MAXSTATEADDEFFECT   AttrType = 72  //最高进阶效果增幅
	AttrType_COMBOADD            AttrType = 73  //连击增幅
	AttrType_PARRYADD            AttrType = 74  //格挡免伤增幅
	AttrType_HEALCRITAL          AttrType = 75  //治疗暴击
	AttrType_HEALBLASTS          AttrType = 76  //治疗暴击增幅
	AttrType_DEFENCESTENGTH      AttrType = 77  //防御免伤增量
	AttrType_HIDEDAMAGEDD        AttrType = 78  //隐藏增伤
	AttrType_ALLATTACK           AttrType = 79  //全体攻击
	AttrType_ALLDEFFENCE         AttrType = 80  //全体防御
	AttrType_ALLMAXHP            AttrType = 81  //全体血量
	AttrType_HEALDEFENCE         AttrType = 82  //被治疗效果降低
	AttrType_SKILLDMG            AttrType = 83  //技能伤害增加
	AttrType_SKILLHEAL           AttrType = 84  //技能治疗效果增加
	AttrType_DEFENCEWEAKNESS     AttrType = 85  //防御减免
	AttrType_ATTACKSUBSP         AttrType = 86  //攻击减怒
	AttrType_BEHURTSUBSP         AttrType = 87  //受击减怒
	AttrType_SKILLDMGADD         AttrType = 88  //技能增伤
	AttrType_SKILLDMGSUB         AttrType = 89  //技能减伤
	AttrType_MUILTIDAMAGEADD     AttrType = 90  //群体伤害加成
	AttrType_SINGLEDAMAGEADD     AttrType = 91  //单体伤害加成
	AttrType_MUILTIDAMAGESUB     AttrType = 92  //群体伤害减益
	AttrType_SINGLEDAMAGESUB     AttrType = 93  //单体伤害减益
	AttrType_MONSTERDMGADD       AttrType = 94  //非Boss伤害增加
	AttrType_BOSSDMGADD          AttrType = 95  //Boss伤害增加
	AttrType_CLEANSESTRENGTH     AttrType = 96  //净化增强
	AttrType_CLEANSEDEFENCE      AttrType = 97  //净化抑制
	AttrType_CNT                 AttrType = 98  //结束
	AttrType_PERCENT_PVPHP       AttrType = 151 //PVP血量百分比
	AttrType_PERCENT_PVPMAXHP    AttrType = 152 //PVP最大血量百分比
	AttrType_PERCENT_PVPATTACK   AttrType = 153 //PVP攻击百分比
	AttrType_PERCENT_PVPDEFENCE  AttrType = 154 //PVP防御百分比
	AttrType_PERCENT_ALLATTACK   AttrType = 161 //全体攻击百分比
	AttrType_PERCENT_ALLDEFFENCE AttrType = 162 //全体防御百分比
	AttrType_PERCENT_ALLMAXHP    AttrType = 163 //全体血量百分比
	AttrType_PERCENT_BEGIN       AttrType = 200 //百分比属性开始
	AttrType_PERCENT_HP          AttrType = 201 //血量百分比
	AttrType_PERCENT_MAXHP       AttrType = 202 //最大血量百分比
	AttrType_PERCENT_ATTACK      AttrType = 203 //攻击百分比
	AttrType_PERCENT_DEFENCE     AttrType = 204 //防御百分比
	AttrType_PERCENT_END         AttrType = 205 //百分比属性结束
	AttrType_FINAL_BEGIN         AttrType = 300 //最终属性开始
	AttrType_FINAL_HP            AttrType = 301 //最终血量百分比
	AttrType_FINAL_MAXHP         AttrType = 302 //最终最大血量百分比
	AttrType_FINAL_ATTACK        AttrType = 303 //最终攻击百分比
	AttrType_FINAL_DEFENCE       AttrType = 304 //最终防御百分比
	AttrType_END                 AttrType = 305 //枚举结束
)

// 基础配置-拳皇/AttrEnum
type AttrTitleType int

const (
	AttrTitleType_BEGIN   AttrTitleType = 1 //枚举开始
	AttrTitleType_Basic   AttrTitleType = 1 //基础属性
	AttrTitleType_Advance AttrTitleType = 2 //流派属性
	AttrTitleType_Ignore  AttrTitleType = 3 //忽视属性
	AttrTitleType_Special AttrTitleType = 4 //特殊属性
	AttrTitleType_END     AttrTitleType = 5 //枚举结束
)

func (e AttrType) String() string {

	if e == AttrType_HP {
		return "HP"
	}
	if e == AttrType_MAXHP {
		return "MAXHP"
	}
	if e == AttrType_ATTACK {
		return "ATTACK"
	}
	if e == AttrType_DEFENCE {
		return "DEFENCE"
	}
	if e == AttrType_PENETRATE {
		return "PENETRATE"
	}
	if e == AttrType_ATTACKSPEED {
		return "ATTACKSPEED"
	}
	if e == AttrType_MOVESPEED {
		return "MOVESPEED"
	}
	if e == AttrType_IGNDODGE {
		return "IGNDODGE"
	}
	if e == AttrType_DODGE {
		return "DODGE"
	}
	if e == AttrType_CRITAL {
		return "CRITAL"
	}
	if e == AttrType_IGNCRITAL {
		return "IGNCRITAL"
	}
	if e == AttrType_BLASTS {
		return "BLASTS"
	}
	if e == AttrType_IGNBLASTS {
		return "IGNBLASTS"
	}
	if e == AttrType_PARRY {
		return "PARRY"
	}
	if e == AttrType_BREAKPARRY {
		return "BREAKPARRY"
	}
	if e == AttrType_STIFFNESS {
		return "STIFFNESS"
	}
	if e == AttrType_STIFFNESSDEF {
		return "STIFFNESSDEF"
	}
	if e == AttrType_ENDURE {
		return "ENDURE"
	}
	if e == AttrType_DEFSTUN {
		return "DEFSTUN"
	}
	if e == AttrType_DEFFREEZE {
		return "DEFFREEZE"
	}
	if e == AttrType_PHYSICALADD {
		return "PHYSICALADD"
	}
	if e == AttrType_PHYSICALSUB {
		return "PHYSICALSUB"
	}
	if e == AttrType_MAGICADD {
		return "MAGICADD"
	}
	if e == AttrType_MAGICSUB {
		return "MAGICSUB"
	}
	if e == AttrType_NEARADD {
		return "NEARADD"
	}
	if e == AttrType_NEARSUB {
		return "NEARSUB"
	}
	if e == AttrType_FARADD {
		return "FARADD"
	}
	if e == AttrType_FARSUB {
		return "FARSUB"
	}
	if e == AttrType_TRICKADD {
		return "TRICKADD"
	}
	if e == AttrType_DEFTRICKADD {
		return "DEFTRICKADD"
	}
	if e == AttrType_DAMAGEADD {
		return "DAMAGEADD"
	}
	if e == AttrType_DAMAGESUB {
		return "DAMAGESUB"
	}
	if e == AttrType_SPINCREASE {
		return "SPINCREASE"
	}
	if e == AttrType_SPRECOVERY {
		return "SPRECOVERY"
	}
	if e == AttrType_SPINCRERESPOSTIVE {
		return "SPINCRERESPOSTIVE"
	}
	if e == AttrType_SPINCRERESNEGIATIVE {
		return "SPINCRERESNEGIATIVE"
	}
	if e == AttrType_SPINCRERESNEFAILED {
		return "SPINCRERESNEFAILED"
	}
	if e == AttrType_SPRETENTION {
		return "SPRETENTION"
	}
	if e == AttrType_SPDESTRUCTION {
		return "SPDESTRUCTION"
	}
	if e == AttrType_EFFECTHEAL {
		return "EFFECTHEAL"
	}
	if e == AttrType_SPINITIAL {
		return "SPINITIAL"
	}
	if e == AttrType_MAXSP {
		return "MAXSP"
	}
	if e == AttrType_SP {
		return "SP"
	}
	if e == AttrType_SPINCRERESADD {
		return "SPINCRERESADD"
	}
	if e == AttrType_SPINCRERESSUB {
		return "SPINCRERESSUB"
	}
	if e == AttrType_SPKILLRECOVER {
		return "SPKILLRECOVER"
	}
	if e == AttrType_SPINITRECOVER {
		return "SPINITRECOVER"
	}
	if e == AttrType_SILENTDEFENCE {
		return "SILENTDEFENCE"
	}
	if e == AttrType_SKILLTRIGGERADD {
		return "SKILLTRIGGERADD"
	}
	if e == AttrType_DEFSKILLTRIGGERADD {
		return "DEFSKILLTRIGGERADD"
	}
	if e == AttrType_HEALADD {
		return "HEALADD"
	}
	if e == AttrType_COUNTERING {
		return "COUNTERING"
	}
	if e == AttrType_DEFCOUNTER {
		return "DEFCOUNTER"
	}
	if e == AttrType_BLOODSUCK {
		return "BLOODSUCK"
	}
	if e == AttrType_DEFBLOODSUCK {
		return "DEFBLOODSUCK"
	}
	if e == AttrType_SPRESERVES {
		return "SPRESERVES"
	}
	if e == AttrType_COMBO {
		return "COMBO"
	}
	if e == AttrType_DEFCOMBO {
		return "DEFCOMBO"
	}
	if e == AttrType_ENDUREBREAK {
		return "ENDUREBREAK"
	}
	if e == AttrType_COUNTERADD {
		return "COUNTERADD"
	}
	if e == AttrType_HEALSUB {
		return "HEALSUB"
	}
	if e == AttrType_ALLDEFENCEADD {
		return "ALLDEFENCEADD"
	}
	if e == AttrType_ALLSTATEADD {
		return "ALLSTATEADD"
	}
	if e == AttrType_MAXDEFENCEADD {
		return "MAXDEFENCEADD"
	}
	if e == AttrType_MAXSTATEADD {
		return "MAXSTATEADD"
	}
	if e == AttrType_ROLEFIRSTRATE {
		return "ROLEFIRSTRATE"
	}
	if e == AttrType_ATTACKRANGE {
		return "ATTACKRANGE"
	}
	if e == AttrType_STIFFNESSSTRENG {
		return "STIFFNESSSTRENG"
	}
	if e == AttrType_STIFFNESSWEAKEN {
		return "STIFFNESSWEAKEN"
	}
	if e == AttrType_SUCKSTRENGTH {
		return "SUCKSTRENGTH"
	}
	if e == AttrType_HITRATE {
		return "HITRATE"
	}
	if e == AttrType_MAXSTATEADDEFFECT {
		return "MAXSTATEADDEFFECT"
	}
	if e == AttrType_COMBOADD {
		return "COMBOADD"
	}
	if e == AttrType_PARRYADD {
		return "PARRYADD"
	}
	if e == AttrType_HEALCRITAL {
		return "HEALCRITAL"
	}
	if e == AttrType_HEALBLASTS {
		return "HEALBLASTS"
	}
	if e == AttrType_DEFENCESTENGTH {
		return "DEFENCESTENGTH"
	}
	if e == AttrType_HIDEDAMAGEDD {
		return "HIDEDAMAGEDD"
	}
	if e == AttrType_ALLATTACK {
		return "ALLATTACK"
	}
	if e == AttrType_ALLDEFFENCE {
		return "ALLDEFFENCE"
	}
	if e == AttrType_ALLMAXHP {
		return "ALLMAXHP"
	}
	if e == AttrType_HEALDEFENCE {
		return "HEALDEFENCE"
	}
	if e == AttrType_SKILLDMG {
		return "SKILLDMG"
	}
	if e == AttrType_SKILLHEAL {
		return "SKILLHEAL"
	}
	if e == AttrType_DEFENCEWEAKNESS {
		return "DEFENCEWEAKNESS"
	}
	if e == AttrType_ATTACKSUBSP {
		return "ATTACKSUBSP"
	}
	if e == AttrType_BEHURTSUBSP {
		return "BEHURTSUBSP"
	}
	if e == AttrType_SKILLDMGADD {
		return "SKILLDMGADD"
	}
	if e == AttrType_SKILLDMGSUB {
		return "SKILLDMGSUB"
	}
	if e == AttrType_MUILTIDAMAGEADD {
		return "MUILTIDAMAGEADD"
	}
	if e == AttrType_SINGLEDAMAGEADD {
		return "SINGLEDAMAGEADD"
	}
	if e == AttrType_MUILTIDAMAGESUB {
		return "MUILTIDAMAGESUB"
	}
	if e == AttrType_SINGLEDAMAGESUB {
		return "SINGLEDAMAGESUB"
	}
	if e == AttrType_MONSTERDMGADD {
		return "MONSTERDMGADD"
	}
	if e == AttrType_BOSSDMGADD {
		return "BOSSDMGADD"
	}
	if e == AttrType_CLEANSESTRENGTH {
		return "CLEANSESTRENGTH"
	}
	if e == AttrType_CLEANSEDEFENCE {
		return "CLEANSEDEFENCE"
	}
	if e == AttrType_CNT {
		return "CNT"
	}
	if e == AttrType_PERCENT_PVPHP {
		return "PERCENT_PVPHP"
	}
	if e == AttrType_PERCENT_PVPMAXHP {
		return "PERCENT_PVPMAXHP"
	}
	if e == AttrType_PERCENT_PVPATTACK {
		return "PERCENT_PVPATTACK"
	}
	if e == AttrType_PERCENT_PVPDEFENCE {
		return "PERCENT_PVPDEFENCE"
	}
	if e == AttrType_PERCENT_ALLATTACK {
		return "PERCENT_ALLATTACK"
	}
	if e == AttrType_PERCENT_ALLDEFFENCE {
		return "PERCENT_ALLDEFFENCE"
	}
	if e == AttrType_PERCENT_ALLMAXHP {
		return "PERCENT_ALLMAXHP"
	}
	if e == AttrType_PERCENT_BEGIN {
		return "PERCENT_BEGIN"
	}
	if e == AttrType_PERCENT_HP {
		return "PERCENT_HP"
	}
	if e == AttrType_PERCENT_MAXHP {
		return "PERCENT_MAXHP"
	}
	if e == AttrType_PERCENT_ATTACK {
		return "PERCENT_ATTACK"
	}
	if e == AttrType_PERCENT_DEFENCE {
		return "PERCENT_DEFENCE"
	}
	if e == AttrType_PERCENT_END {
		return "PERCENT_END"
	}
	if e == AttrType_FINAL_BEGIN {
		return "FINAL_BEGIN"
	}
	if e == AttrType_FINAL_HP {
		return "FINAL_HP"
	}
	if e == AttrType_FINAL_MAXHP {
		return "FINAL_MAXHP"
	}
	if e == AttrType_FINAL_ATTACK {
		return "FINAL_ATTACK"
	}
	if e == AttrType_FINAL_DEFENCE {
		return "FINAL_DEFENCE"
	}

	return ""

}

func GetAttrType(tp string) (AttrType, bool) {
	switch tp {
	case "HP":
		return AttrType_HP, true
	case "MAXHP":
		return AttrType_MAXHP, true
	case "ATTACK":
		return AttrType_ATTACK, true
	case "DEFENCE":
		return AttrType_DEFENCE, true
	case "PENETRATE":
		return AttrType_PENETRATE, true
	case "ATTACKSPEED":
		return AttrType_ATTACKSPEED, true
	case "MOVESPEED":
		return AttrType_MOVESPEED, true
	case "IGNDODGE":
		return AttrType_IGNDODGE, true
	case "DODGE":
		return AttrType_DODGE, true
	case "CRITAL":
		return AttrType_CRITAL, true
	case "IGNCRITAL":
		return AttrType_IGNCRITAL, true
	case "BLASTS":
		return AttrType_BLASTS, true
	case "IGNBLASTS":
		return AttrType_IGNBLASTS, true
	case "PARRY":
		return AttrType_PARRY, true
	case "BREAKPARRY":
		return AttrType_BREAKPARRY, true
	case "STIFFNESS":
		return AttrType_STIFFNESS, true
	case "STIFFNESSDEF":
		return AttrType_STIFFNESSDEF, true
	case "ENDURE":
		return AttrType_ENDURE, true
	case "DEFSTUN":
		return AttrType_DEFSTUN, true
	case "DEFFREEZE":
		return AttrType_DEFFREEZE, true
	case "PHYSICALADD":
		return AttrType_PHYSICALADD, true
	case "PHYSICALSUB":
		return AttrType_PHYSICALSUB, true
	case "MAGICADD":
		return AttrType_MAGICADD, true
	case "MAGICSUB":
		return AttrType_MAGICSUB, true
	case "NEARADD":
		return AttrType_NEARADD, true
	case "NEARSUB":
		return AttrType_NEARSUB, true
	case "FARADD":
		return AttrType_FARADD, true
	case "FARSUB":
		return AttrType_FARSUB, true
	case "TRICKADD":
		return AttrType_TRICKADD, true
	case "DEFTRICKADD":
		return AttrType_DEFTRICKADD, true
	case "DAMAGEADD":
		return AttrType_DAMAGEADD, true
	case "DAMAGESUB":
		return AttrType_DAMAGESUB, true
	case "SPINCREASE":
		return AttrType_SPINCREASE, true
	case "SPRECOVERY":
		return AttrType_SPRECOVERY, true
	case "SPINCRERESPOSTIVE":
		return AttrType_SPINCRERESPOSTIVE, true
	case "SPINCRERESNEGIATIVE":
		return AttrType_SPINCRERESNEGIATIVE, true
	case "SPINCRERESNEFAILED":
		return AttrType_SPINCRERESNEFAILED, true
	case "SPRETENTION":
		return AttrType_SPRETENTION, true
	case "SPDESTRUCTION":
		return AttrType_SPDESTRUCTION, true
	case "EFFECTHEAL":
		return AttrType_EFFECTHEAL, true
	case "SPINITIAL":
		return AttrType_SPINITIAL, true
	case "MAXSP":
		return AttrType_MAXSP, true
	case "SP":
		return AttrType_SP, true
	case "SPINCRERESADD":
		return AttrType_SPINCRERESADD, true
	case "SPINCRERESSUB":
		return AttrType_SPINCRERESSUB, true
	case "SPKILLRECOVER":
		return AttrType_SPKILLRECOVER, true
	case "SPINITRECOVER":
		return AttrType_SPINITRECOVER, true
	case "SILENTDEFENCE":
		return AttrType_SILENTDEFENCE, true
	case "SKILLTRIGGERADD":
		return AttrType_SKILLTRIGGERADD, true
	case "DEFSKILLTRIGGERADD":
		return AttrType_DEFSKILLTRIGGERADD, true
	case "HEALADD":
		return AttrType_HEALADD, true
	case "COUNTERING":
		return AttrType_COUNTERING, true
	case "DEFCOUNTER":
		return AttrType_DEFCOUNTER, true
	case "BLOODSUCK":
		return AttrType_BLOODSUCK, true
	case "DEFBLOODSUCK":
		return AttrType_DEFBLOODSUCK, true
	case "SPRESERVES":
		return AttrType_SPRESERVES, true
	case "COMBO":
		return AttrType_COMBO, true
	case "DEFCOMBO":
		return AttrType_DEFCOMBO, true
	case "ENDUREBREAK":
		return AttrType_ENDUREBREAK, true
	case "COUNTERADD":
		return AttrType_COUNTERADD, true
	case "HEALSUB":
		return AttrType_HEALSUB, true
	case "ALLDEFENCEADD":
		return AttrType_ALLDEFENCEADD, true
	case "ALLSTATEADD":
		return AttrType_ALLSTATEADD, true
	case "MAXDEFENCEADD":
		return AttrType_MAXDEFENCEADD, true
	case "MAXSTATEADD":
		return AttrType_MAXSTATEADD, true
	case "ROLEFIRSTRATE":
		return AttrType_ROLEFIRSTRATE, true
	case "ATTACKRANGE":
		return AttrType_ATTACKRANGE, true
	case "STIFFNESSSTRENG":
		return AttrType_STIFFNESSSTRENG, true
	case "STIFFNESSWEAKEN":
		return AttrType_STIFFNESSWEAKEN, true
	case "SUCKSTRENGTH":
		return AttrType_SUCKSTRENGTH, true
	case "HITRATE":
		return AttrType_HITRATE, true
	case "MAXSTATEADDEFFECT":
		return AttrType_MAXSTATEADDEFFECT, true
	case "COMBOADD":
		return AttrType_COMBOADD, true
	case "PARRYADD":
		return AttrType_PARRYADD, true
	case "HEALCRITAL":
		return AttrType_HEALCRITAL, true
	case "HEALBLASTS":
		return AttrType_HEALBLASTS, true
	case "DEFENCESTENGTH":
		return AttrType_DEFENCESTENGTH, true
	case "HIDEDAMAGEDD":
		return AttrType_HIDEDAMAGEDD, true
	case "ALLATTACK":
		return AttrType_ALLATTACK, true
	case "ALLDEFFENCE":
		return AttrType_ALLDEFFENCE, true
	case "ALLMAXHP":
		return AttrType_ALLMAXHP, true
	case "HEALDEFENCE":
		return AttrType_HEALDEFENCE, true
	case "SKILLDMG":
		return AttrType_SKILLDMG, true
	case "SKILLHEAL":
		return AttrType_SKILLHEAL, true
	case "DEFENCEWEAKNESS":
		return AttrType_DEFENCEWEAKNESS, true
	case "ATTACKSUBSP":
		return AttrType_ATTACKSUBSP, true
	case "BEHURTSUBSP":
		return AttrType_BEHURTSUBSP, true
	case "SKILLDMGADD":
		return AttrType_SKILLDMGADD, true
	case "SKILLDMGSUB":
		return AttrType_SKILLDMGSUB, true
	case "MUILTIDAMAGEADD":
		return AttrType_MUILTIDAMAGEADD, true
	case "SINGLEDAMAGEADD":
		return AttrType_SINGLEDAMAGEADD, true
	case "MUILTIDAMAGESUB":
		return AttrType_MUILTIDAMAGESUB, true
	case "SINGLEDAMAGESUB":
		return AttrType_SINGLEDAMAGESUB, true
	case "MONSTERDMGADD":
		return AttrType_MONSTERDMGADD, true
	case "BOSSDMGADD":
		return AttrType_BOSSDMGADD, true
	case "CLEANSESTRENGTH":
		return AttrType_CLEANSESTRENGTH, true
	case "CLEANSEDEFENCE":
		return AttrType_CLEANSEDEFENCE, true
	case "CNT":
		return AttrType_CNT, true
	case "PERCENT_PVPHP":
		return AttrType_PERCENT_PVPHP, true
	case "PERCENT_PVPMAXHP":
		return AttrType_PERCENT_PVPMAXHP, true
	case "PERCENT_PVPATTACK":
		return AttrType_PERCENT_PVPATTACK, true
	case "PERCENT_PVPDEFENCE":
		return AttrType_PERCENT_PVPDEFENCE, true
	case "PERCENT_ALLATTACK":
		return AttrType_PERCENT_ALLATTACK, true
	case "PERCENT_ALLDEFFENCE":
		return AttrType_PERCENT_ALLDEFFENCE, true
	case "PERCENT_ALLMAXHP":
		return AttrType_PERCENT_ALLMAXHP, true
	case "PERCENT_BEGIN":
		return AttrType_PERCENT_BEGIN, true
	case "PERCENT_HP":
		return AttrType_PERCENT_HP, true
	case "PERCENT_MAXHP":
		return AttrType_PERCENT_MAXHP, true
	case "PERCENT_ATTACK":
		return AttrType_PERCENT_ATTACK, true
	case "PERCENT_DEFENCE":
		return AttrType_PERCENT_DEFENCE, true
	case "PERCENT_END":
		return AttrType_PERCENT_END, true
	case "FINAL_BEGIN":
		return AttrType_FINAL_BEGIN, true
	case "FINAL_HP":
		return AttrType_FINAL_HP, true
	case "FINAL_MAXHP":
		return AttrType_FINAL_MAXHP, true
	case "FINAL_ATTACK":
		return AttrType_FINAL_ATTACK, true
	case "FINAL_DEFENCE":
		return AttrType_FINAL_DEFENCE, true

	}
	return -1, false

}

func GetAttrType_ByAlias(tp string) (AttrType, bool) {
	switch tp {
	case "血量":
		return AttrType_HP, true
	case "最大血量":
		return AttrType_MAXHP, true
	case "攻击":
		return AttrType_ATTACK, true
	case "防御":
		return AttrType_DEFENCE, true
	case "穿透率":
		return AttrType_PENETRATE, true
	case "攻击速度":
		return AttrType_ATTACKSPEED, true
	case "移动速度":
		return AttrType_MOVESPEED, true
	case "闪避抵抗":
		return AttrType_IGNDODGE, true
	case "闪避":
		return AttrType_DODGE, true
	case "暴击率":
		return AttrType_CRITAL, true
	case "免暴率":
		return AttrType_IGNCRITAL, true
	case "暴击上限":
		return AttrType_BLASTS, true
	case "暴击抵抗":
		return AttrType_IGNBLASTS, true
	case "格挡":
		return AttrType_PARRY, true
	case "格挡抵抗":
		return AttrType_BREAKPARRY, true
	case "硬直":
		return AttrType_STIFFNESS, true
	case "硬直抵抗":
		return AttrType_STIFFNESSDEF, true
	case "霸体免控率":
		return AttrType_ENDURE, true
	case "眩晕抵抗":
		return AttrType_DEFSTUN, true
	case "冰冻抵抗":
		return AttrType_DEFFREEZE, true
	case "物理增伤":
		return AttrType_PHYSICALADD, true
	case "物理减伤":
		return AttrType_PHYSICALSUB, true
	case "魔法增伤":
		return AttrType_MAGICADD, true
	case "魔法减伤":
		return AttrType_MAGICSUB, true
	case "近战增伤":
		return AttrType_NEARADD, true
	case "近战减伤":
		return AttrType_NEARSUB, true
	case "远程增伤":
		return AttrType_FARADD, true
	case "远程减伤":
		return AttrType_FARSUB, true
	case "大招增伤":
		return AttrType_TRICKADD, true
	case "大招减伤":
		return AttrType_DEFTRICKADD, true
	case "最终增伤":
		return AttrType_DAMAGEADD, true
	case "最终减伤":
		return AttrType_DAMAGESUB, true
	case "怒气增加":
		return AttrType_SPINCREASE, true
	case "怒气回复":
		return AttrType_SPRECOVERY, true
	case "攻击回怒":
		return AttrType_SPINCRERESPOSTIVE, true
	case "受击回怒":
		return AttrType_SPINCRERESNEGIATIVE, true
	case "战败回怒":
		return AttrType_SPINCRERESNEFAILED, true
	case "怒气截留":
		return AttrType_SPRETENTION, true
	case "自爆减怒":
		return AttrType_SPDESTRUCTION, true
	case "被治疗效果":
		return AttrType_EFFECTHEAL, true
	case "初始怒气":
		return AttrType_SPINITIAL, true
	case "最大怒气":
		return AttrType_MAXSP, true
	case "怒气":
		return AttrType_SP, true
	case "回怒增幅":
		return AttrType_SPINCRERESADD, true
	case "回怒降幅":
		return AttrType_SPINCRERESSUB, true
	case "击杀回怒":
		return AttrType_SPKILLRECOVER, true
	case "初始怒气回复":
		return AttrType_SPINITRECOVER, true
	case "沉默抵抗":
		return AttrType_SILENTDEFENCE, true
	case "技能触发率增长":
		return AttrType_SKILLTRIGGERADD, true
	case "技能触发率抵抗":
		return AttrType_DEFSKILLTRIGGERADD, true
	case "治疗增幅":
		return AttrType_HEALADD, true
	case "反击率":
		return AttrType_COUNTERING, true
	case "反击抵抗":
		return AttrType_DEFCOUNTER, true
	case "吸血系数":
		return AttrType_BLOODSUCK, true
	case "吸血抗性":
		return AttrType_DEFBLOODSUCK, true
	case "怒气储量":
		return AttrType_SPRESERVES, true
	case "连击":
		return AttrType_COMBO, true
	case "连击抗性":
		return AttrType_DEFCOMBO, true
	case "霸体免控抵抗":
		return AttrType_ENDUREBREAK, true
	case "反击增幅":
		return AttrType_COUNTERADD, true
	case "治疗抑制":
		return AttrType_HEALSUB, true
	case "全状态抗性":
		return AttrType_ALLDEFENCEADD, true
	case "全状态属性":
		return AttrType_ALLSTATEADD, true
	case "最高状态抗性":
		return AttrType_MAXDEFENCEADD, true
	case "最高状态属性":
		return AttrType_MAXSTATEADD, true
	case "角色先手值":
		return AttrType_ROLEFIRSTRATE, true
	case "攻击距离":
		return AttrType_ATTACKRANGE, true
	case "硬直增强":
		return AttrType_STIFFNESSSTRENG, true
	case "硬直削弱":
		return AttrType_STIFFNESSWEAKEN, true
	case "吸血增幅":
		return AttrType_SUCKSTRENGTH, true
	case "命中率":
		return AttrType_HITRATE, true
	case "最高进阶效果增幅":
		return AttrType_MAXSTATEADDEFFECT, true
	case "连击增幅":
		return AttrType_COMBOADD, true
	case "格挡免伤增幅":
		return AttrType_PARRYADD, true
	case "治疗暴击":
		return AttrType_HEALCRITAL, true
	case "治疗暴击增幅":
		return AttrType_HEALBLASTS, true
	case "防御免伤增量":
		return AttrType_DEFENCESTENGTH, true
	case "隐藏增伤":
		return AttrType_HIDEDAMAGEDD, true
	case "全体攻击":
		return AttrType_ALLATTACK, true
	case "全体防御":
		return AttrType_ALLDEFFENCE, true
	case "全体血量":
		return AttrType_ALLMAXHP, true
	case "被治疗效果降低":
		return AttrType_HEALDEFENCE, true
	case "技能伤害增加":
		return AttrType_SKILLDMG, true
	case "技能治疗效果增加":
		return AttrType_SKILLHEAL, true
	case "防御减免":
		return AttrType_DEFENCEWEAKNESS, true
	case "攻击减怒":
		return AttrType_ATTACKSUBSP, true
	case "受击减怒":
		return AttrType_BEHURTSUBSP, true
	case "技能增伤":
		return AttrType_SKILLDMGADD, true
	case "技能减伤":
		return AttrType_SKILLDMGSUB, true
	case "群体伤害加成":
		return AttrType_MUILTIDAMAGEADD, true
	case "单体伤害加成":
		return AttrType_SINGLEDAMAGEADD, true
	case "群体伤害减益":
		return AttrType_MUILTIDAMAGESUB, true
	case "单体伤害减益":
		return AttrType_SINGLEDAMAGESUB, true
	case "非Boss伤害增加":
		return AttrType_MONSTERDMGADD, true
	case "Boss伤害增加":
		return AttrType_BOSSDMGADD, true
	case "净化增强":
		return AttrType_CLEANSESTRENGTH, true
	case "净化抑制":
		return AttrType_CLEANSEDEFENCE, true
	case "结束":
		return AttrType_CNT, true
	case "PVP血量百分比":
		return AttrType_PERCENT_PVPHP, true
	case "PVP最大血量百分比":
		return AttrType_PERCENT_PVPMAXHP, true
	case "PVP攻击百分比":
		return AttrType_PERCENT_PVPATTACK, true
	case "PVP防御百分比":
		return AttrType_PERCENT_PVPDEFENCE, true
	case "全体攻击百分比":
		return AttrType_PERCENT_ALLATTACK, true
	case "全体防御百分比":
		return AttrType_PERCENT_ALLDEFFENCE, true
	case "全体血量百分比":
		return AttrType_PERCENT_ALLMAXHP, true
	case "百分比属性开始":
		return AttrType_PERCENT_BEGIN, true
	case "血量百分比":
		return AttrType_PERCENT_HP, true
	case "最大血量百分比":
		return AttrType_PERCENT_MAXHP, true
	case "攻击百分比":
		return AttrType_PERCENT_ATTACK, true
	case "防御百分比":
		return AttrType_PERCENT_DEFENCE, true
	case "百分比属性结束":
		return AttrType_PERCENT_END, true
	case "最终属性开始":
		return AttrType_FINAL_BEGIN, true
	case "最终血量百分比":
		return AttrType_FINAL_HP, true
	case "最终最大血量百分比":
		return AttrType_FINAL_MAXHP, true
	case "最终攻击百分比":
		return AttrType_FINAL_ATTACK, true
	case "最终防御百分比":
		return AttrType_FINAL_DEFENCE, true

	}
	return -1, false

}

func (e AttrTitleType) String() string {

	if e == AttrTitleType_Basic {
		return "Basic"
	}
	if e == AttrTitleType_Advance {
		return "Advance"
	}
	if e == AttrTitleType_Ignore {
		return "Ignore"
	}
	if e == AttrTitleType_Special {
		return "Special"
	}

	return ""

}

func GetAttrTitleType(tp string) (AttrTitleType, bool) {
	switch tp {
	case "Basic":
		return AttrTitleType_Basic, true
	case "Advance":
		return AttrTitleType_Advance, true
	case "Ignore":
		return AttrTitleType_Ignore, true
	case "Special":
		return AttrTitleType_Special, true

	}
	return -1, false

}

func GetAttrTitleType_ByAlias(tp string) (AttrTitleType, bool) {
	switch tp {
	case "基础属性":
		return AttrTitleType_Basic, true
	case "流派属性":
		return AttrTitleType_Advance, true
	case "忽视属性":
		return AttrTitleType_Ignore, true
	case "特殊属性":
		return AttrTitleType_Special, true

	}
	return -1, false

}
