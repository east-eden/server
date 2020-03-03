package define

// 一级属性
const (
	Att_Str = iota // 力量
	Att_Agl        // 敏捷
	Att_Con        // 耐力
	Att_Int        // 智力
	Att_End
)

// 二级属性
const (
	AttEx_HP  = iota // 血量
	AttEx_MP         // 蓝量
	AttEx_Atn        // 物理攻击力
	AttEx_Def        // 物理防御力
	AttEx_Ats        // 魔法攻击力
	AttEx_Adf        // 魔法防御力
	AttEx_End
)

// att entry
type AttEntry struct {
	ID   int32  `json:"_id"`
	Desc string `json:"Desc"`

	Str int32 `json:"Str"`
	Agl int32 `json:"Agl"`
	Con int32 `json:"Con"`
	Int int32 `json:"Int"`

	Hp  int32 `json:"Hp"`
	Mp  int32 `json:"Mp"`
	Atn int32 `json:"Atn"`
	Def int32 `json:"Def"`
	Ats int32 `json:"Ats"`
	Adf int32 `json:"Adf"`
}
