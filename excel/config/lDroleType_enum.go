package config

// LD英雄表-拳皇/LDroleType
type RoleRaceType int

const (
	RoleRaceType_BEGIN   RoleRaceType = -1 //枚举开始
	RoleRaceType_All     RoleRaceType = -1 //全体种族
	RoleRaceType_Empty   RoleRaceType = 0  //无种族
	RoleRaceType_Human   RoleRaceType = 1  //人类
	RoleRaceType_Animal  RoleRaceType = 2  //动物
	RoleRaceType_Machine RoleRaceType = 3  //机械
	RoleRaceType_Elf     RoleRaceType = 4  //精灵
	RoleRaceType_END     RoleRaceType = 5  //枚举结束
)

// LD英雄表-拳皇/LDroleType
type LDRoleProfessionType int

const (
	LDRoleProfessionType_BEGIN   LDRoleProfessionType = -1 //枚举开始
	LDRoleProfessionType_All     LDRoleProfessionType = -1 //全体职业
	LDRoleProfessionType_Empty   LDRoleProfessionType = 0  //无职业
	LDRoleProfessionType_Warrior LDRoleProfessionType = 1  //战士
	LDRoleProfessionType_Mage    LDRoleProfessionType = 2  //法师
	LDRoleProfessionType_Archer  LDRoleProfessionType = 3  //射手
	LDRoleProfessionType_Assasin LDRoleProfessionType = 4  //刺客
	LDRoleProfessionType_Tank    LDRoleProfessionType = 5  //坦克
	LDRoleProfessionType_Assist  LDRoleProfessionType = 6  //辅助
	LDRoleProfessionType_END     LDRoleProfessionType = 7  //枚举结束
)

// LD英雄表-拳皇/LDroleType
type FetterType int

const (
	FetterType_BEGIN      FetterType = 0 //枚举开始
	FetterType_Profession FetterType = 0 //职业
	FetterType_Race       FetterType = 1 //种族
	FetterType_END        FetterType = 2 //枚举结束
)

// LD英雄表-拳皇/LDroleType
type HerosConfigId int

const (
	HerosConfigId_BEGIN          HerosConfigId = 1        //枚举开始
	HerosConfigId_BaseBallBoy    HerosConfigId = 20100042 //棒球小子
	HerosConfigId_LoserBoy       HerosConfigId = 20200001 //破碎棒球小子
	HerosConfigId_MediumBoy      HerosConfigId = 20200002 //中阶棒球小子
	HerosConfigId_SeniorBoy      HerosConfigId = 20200003 //高阶棒球小子
	HerosConfigId_PrefessionBoy  HerosConfigId = 20200004 //特街棒球小子
	HerosConfigId_GuanYu         HerosConfigId = 20100052 //关羽
	HerosConfigId_BianXingGuai   HerosConfigId = 20100053 //变形怪
	HerosConfigId_Dragon         HerosConfigId = 20100049 //龙
	HerosConfigId_LanSiLuoTe     HerosConfigId = 20100043 //兰斯洛特
	HerosConfigId_Coldy          HerosConfigId = 20100041 //寇蒂
	HerosConfigId_Shikigami      HerosConfigId = 20100060 //式神
	HerosConfigId_JiangZiYa      HerosConfigId = 20100062 //姜子牙
	HerosConfigId_Angel          HerosConfigId = 20100063 //天使
	HerosConfigId_MonkeyKing     HerosConfigId = 20100065 //黑悟空
	HerosConfigId_LittleDice     HerosConfigId = 20100066 //小小骰子
	HerosConfigId_GuanZhi        HerosConfigId = 20300001 //广智
	HerosConfigId_JinXiGui       HerosConfigId = 20300002 //精细鬼
	HerosConfigId_BenBoBa        HerosConfigId = 20300003 //奔波霸
	HerosConfigId_XiGua          HerosConfigId = 20300004 //西瓜
	HerosConfigId_YuanShiTianZun HerosConfigId = 20300005 //原始天尊
	HerosConfigId_RuLaiFoZu      HerosConfigId = 20300006 //如来佛祖
	HerosConfigId_KaKaLuoTe      HerosConfigId = 20300007 //卡卡罗特
	HerosConfigId_MaoHou         HerosConfigId = 20300008 //毛猴
	HerosConfigId_DiceOne        HerosConfigId = 20100066 //一号骰子
	HerosConfigId_DiceTwo        HerosConfigId = 20100067 //二号骰子
	HerosConfigId_DiceThree      HerosConfigId = 20100068 //三号骰子
	HerosConfigId_DiceFour       HerosConfigId = 20100069 //四号骰子
	HerosConfigId_DiceFive       HerosConfigId = 20100070 //五号骰子
	HerosConfigId_DiceSix        HerosConfigId = 20100071 //六号骰子
	HerosConfigId_DragonEgg      HerosConfigId = 20100048 //龙蛋
	HerosConfigId_Mech           HerosConfigId = 20100044 //机甲
	HerosConfigId_IronMan        HerosConfigId = 20100045 //钢铁侠
	HerosConfigId_Frog           HerosConfigId = 20100046 //青蛙
	HerosConfigId_Prince         HerosConfigId = 20100047 //小王子
	HerosConfigId_END            HerosConfigId = 20300009 //枚举结束
)

func (e RoleRaceType) String() string {

	if e == RoleRaceType_All {
		return "All"
	}
	if e == RoleRaceType_Empty {
		return "Empty"
	}
	if e == RoleRaceType_Human {
		return "Human"
	}
	if e == RoleRaceType_Animal {
		return "Animal"
	}
	if e == RoleRaceType_Machine {
		return "Machine"
	}
	if e == RoleRaceType_Elf {
		return "Elf"
	}

	return ""

}

func GetRoleRaceType(tp string) (RoleRaceType, bool) {
	switch tp {
	case "All":
		return RoleRaceType_All, true
	case "Empty":
		return RoleRaceType_Empty, true
	case "Human":
		return RoleRaceType_Human, true
	case "Animal":
		return RoleRaceType_Animal, true
	case "Machine":
		return RoleRaceType_Machine, true
	case "Elf":
		return RoleRaceType_Elf, true

	}
	return -1, false

}

func GetRoleRaceType_ByAlias(tp string) (RoleRaceType, bool) {
	switch tp {
	case "全体种族":
		return RoleRaceType_All, true
	case "无种族":
		return RoleRaceType_Empty, true
	case "人类":
		return RoleRaceType_Human, true
	case "动物":
		return RoleRaceType_Animal, true
	case "机械":
		return RoleRaceType_Machine, true
	case "精灵":
		return RoleRaceType_Elf, true

	}
	return -1, false

}

func (e LDRoleProfessionType) String() string {

	if e == LDRoleProfessionType_All {
		return "All"
	}
	if e == LDRoleProfessionType_Empty {
		return "Empty"
	}
	if e == LDRoleProfessionType_Warrior {
		return "Warrior"
	}
	if e == LDRoleProfessionType_Mage {
		return "Mage"
	}
	if e == LDRoleProfessionType_Archer {
		return "Archer"
	}
	if e == LDRoleProfessionType_Assasin {
		return "Assasin"
	}
	if e == LDRoleProfessionType_Tank {
		return "Tank"
	}
	if e == LDRoleProfessionType_Assist {
		return "Assist"
	}

	return ""

}

func GetLDRoleProfessionType(tp string) (LDRoleProfessionType, bool) {
	switch tp {
	case "All":
		return LDRoleProfessionType_All, true
	case "Empty":
		return LDRoleProfessionType_Empty, true
	case "Warrior":
		return LDRoleProfessionType_Warrior, true
	case "Mage":
		return LDRoleProfessionType_Mage, true
	case "Archer":
		return LDRoleProfessionType_Archer, true
	case "Assasin":
		return LDRoleProfessionType_Assasin, true
	case "Tank":
		return LDRoleProfessionType_Tank, true
	case "Assist":
		return LDRoleProfessionType_Assist, true

	}
	return -1, false

}

func GetLDRoleProfessionType_ByAlias(tp string) (LDRoleProfessionType, bool) {
	switch tp {
	case "全体职业":
		return LDRoleProfessionType_All, true
	case "无职业":
		return LDRoleProfessionType_Empty, true
	case "战士":
		return LDRoleProfessionType_Warrior, true
	case "法师":
		return LDRoleProfessionType_Mage, true
	case "射手":
		return LDRoleProfessionType_Archer, true
	case "刺客":
		return LDRoleProfessionType_Assasin, true
	case "坦克":
		return LDRoleProfessionType_Tank, true
	case "辅助":
		return LDRoleProfessionType_Assist, true

	}
	return -1, false

}

func (e FetterType) String() string {

	if e == FetterType_Profession {
		return "Profession"
	}
	if e == FetterType_Race {
		return "Race"
	}

	return ""

}

func GetFetterType(tp string) (FetterType, bool) {
	switch tp {
	case "Profession":
		return FetterType_Profession, true
	case "Race":
		return FetterType_Race, true

	}
	return -1, false

}

func GetFetterType_ByAlias(tp string) (FetterType, bool) {
	switch tp {
	case "职业":
		return FetterType_Profession, true
	case "种族":
		return FetterType_Race, true

	}
	return -1, false

}

func (e HerosConfigId) String() string {

	if e == HerosConfigId_BaseBallBoy {
		return "BaseBallBoy"
	}
	if e == HerosConfigId_LoserBoy {
		return "LoserBoy"
	}
	if e == HerosConfigId_MediumBoy {
		return "MediumBoy"
	}
	if e == HerosConfigId_SeniorBoy {
		return "SeniorBoy"
	}
	if e == HerosConfigId_PrefessionBoy {
		return "PrefessionBoy"
	}
	if e == HerosConfigId_GuanYu {
		return "GuanYu"
	}
	if e == HerosConfigId_BianXingGuai {
		return "BianXingGuai"
	}
	if e == HerosConfigId_Dragon {
		return "Dragon"
	}
	if e == HerosConfigId_LanSiLuoTe {
		return "LanSiLuoTe"
	}
	if e == HerosConfigId_Coldy {
		return "Coldy"
	}
	if e == HerosConfigId_Shikigami {
		return "Shikigami"
	}
	if e == HerosConfigId_JiangZiYa {
		return "JiangZiYa"
	}
	if e == HerosConfigId_Angel {
		return "Angel"
	}
	if e == HerosConfigId_MonkeyKing {
		return "MonkeyKing"
	}
	if e == HerosConfigId_LittleDice {
		return "LittleDice"
	}
	if e == HerosConfigId_GuanZhi {
		return "GuanZhi"
	}
	if e == HerosConfigId_JinXiGui {
		return "JinXiGui"
	}
	if e == HerosConfigId_BenBoBa {
		return "BenBoBa"
	}
	if e == HerosConfigId_XiGua {
		return "XiGua"
	}
	if e == HerosConfigId_YuanShiTianZun {
		return "YuanShiTianZun"
	}
	if e == HerosConfigId_RuLaiFoZu {
		return "RuLaiFoZu"
	}
	if e == HerosConfigId_KaKaLuoTe {
		return "KaKaLuoTe"
	}
	if e == HerosConfigId_MaoHou {
		return "MaoHou"
	}
	if e == HerosConfigId_DiceOne {
		return "DiceOne"
	}
	if e == HerosConfigId_DiceTwo {
		return "DiceTwo"
	}
	if e == HerosConfigId_DiceThree {
		return "DiceThree"
	}
	if e == HerosConfigId_DiceFour {
		return "DiceFour"
	}
	if e == HerosConfigId_DiceFive {
		return "DiceFive"
	}
	if e == HerosConfigId_DiceSix {
		return "DiceSix"
	}
	if e == HerosConfigId_DragonEgg {
		return "DragonEgg"
	}
	if e == HerosConfigId_Mech {
		return "Mech"
	}
	if e == HerosConfigId_IronMan {
		return "IronMan"
	}
	if e == HerosConfigId_Frog {
		return "Frog"
	}
	if e == HerosConfigId_Prince {
		return "Prince"
	}

	return ""

}

func GetHerosConfigId(tp string) (HerosConfigId, bool) {
	switch tp {
	case "BaseBallBoy":
		return HerosConfigId_BaseBallBoy, true
	case "LoserBoy":
		return HerosConfigId_LoserBoy, true
	case "MediumBoy":
		return HerosConfigId_MediumBoy, true
	case "SeniorBoy":
		return HerosConfigId_SeniorBoy, true
	case "PrefessionBoy":
		return HerosConfigId_PrefessionBoy, true
	case "GuanYu":
		return HerosConfigId_GuanYu, true
	case "BianXingGuai":
		return HerosConfigId_BianXingGuai, true
	case "Dragon":
		return HerosConfigId_Dragon, true
	case "LanSiLuoTe":
		return HerosConfigId_LanSiLuoTe, true
	case "Coldy":
		return HerosConfigId_Coldy, true
	case "Shikigami":
		return HerosConfigId_Shikigami, true
	case "JiangZiYa":
		return HerosConfigId_JiangZiYa, true
	case "Angel":
		return HerosConfigId_Angel, true
	case "MonkeyKing":
		return HerosConfigId_MonkeyKing, true
	case "LittleDice":
		return HerosConfigId_LittleDice, true
	case "GuanZhi":
		return HerosConfigId_GuanZhi, true
	case "JinXiGui":
		return HerosConfigId_JinXiGui, true
	case "BenBoBa":
		return HerosConfigId_BenBoBa, true
	case "XiGua":
		return HerosConfigId_XiGua, true
	case "YuanShiTianZun":
		return HerosConfigId_YuanShiTianZun, true
	case "RuLaiFoZu":
		return HerosConfigId_RuLaiFoZu, true
	case "KaKaLuoTe":
		return HerosConfigId_KaKaLuoTe, true
	case "MaoHou":
		return HerosConfigId_MaoHou, true
	case "DiceOne":
		return HerosConfigId_DiceOne, true
	case "DiceTwo":
		return HerosConfigId_DiceTwo, true
	case "DiceThree":
		return HerosConfigId_DiceThree, true
	case "DiceFour":
		return HerosConfigId_DiceFour, true
	case "DiceFive":
		return HerosConfigId_DiceFive, true
	case "DiceSix":
		return HerosConfigId_DiceSix, true
	case "DragonEgg":
		return HerosConfigId_DragonEgg, true
	case "Mech":
		return HerosConfigId_Mech, true
	case "IronMan":
		return HerosConfigId_IronMan, true
	case "Frog":
		return HerosConfigId_Frog, true
	case "Prince":
		return HerosConfigId_Prince, true

	}
	return -1, false

}

func GetHerosConfigId_ByAlias(tp string) (HerosConfigId, bool) {
	switch tp {
	case "棒球小子":
		return HerosConfigId_BaseBallBoy, true
	case "破碎棒球小子":
		return HerosConfigId_LoserBoy, true
	case "中阶棒球小子":
		return HerosConfigId_MediumBoy, true
	case "高阶棒球小子":
		return HerosConfigId_SeniorBoy, true
	case "特街棒球小子":
		return HerosConfigId_PrefessionBoy, true
	case "关羽":
		return HerosConfigId_GuanYu, true
	case "变形怪":
		return HerosConfigId_BianXingGuai, true
	case "龙":
		return HerosConfigId_Dragon, true
	case "兰斯洛特":
		return HerosConfigId_LanSiLuoTe, true
	case "寇蒂":
		return HerosConfigId_Coldy, true
	case "式神":
		return HerosConfigId_Shikigami, true
	case "姜子牙":
		return HerosConfigId_JiangZiYa, true
	case "天使":
		return HerosConfigId_Angel, true
	case "黑悟空":
		return HerosConfigId_MonkeyKing, true
	case "小小骰子":
		return HerosConfigId_LittleDice, true
	case "广智":
		return HerosConfigId_GuanZhi, true
	case "精细鬼":
		return HerosConfigId_JinXiGui, true
	case "奔波霸":
		return HerosConfigId_BenBoBa, true
	case "西瓜":
		return HerosConfigId_XiGua, true
	case "原始天尊":
		return HerosConfigId_YuanShiTianZun, true
	case "如来佛祖":
		return HerosConfigId_RuLaiFoZu, true
	case "卡卡罗特":
		return HerosConfigId_KaKaLuoTe, true
	case "毛猴":
		return HerosConfigId_MaoHou, true
	case "一号骰子":
		return HerosConfigId_DiceOne, true
	case "二号骰子":
		return HerosConfigId_DiceTwo, true
	case "三号骰子":
		return HerosConfigId_DiceThree, true
	case "四号骰子":
		return HerosConfigId_DiceFour, true
	case "五号骰子":
		return HerosConfigId_DiceFive, true
	case "六号骰子":
		return HerosConfigId_DiceSix, true
	case "龙蛋":
		return HerosConfigId_DragonEgg, true
	case "机甲":
		return HerosConfigId_Mech, true
	case "钢铁侠":
		return HerosConfigId_IronMan, true
	case "青蛙":
		return HerosConfigId_Frog, true
	case "小王子":
		return HerosConfigId_Prince, true

	}
	return -1, false

}
