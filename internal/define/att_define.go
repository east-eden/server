package define

// 一级属性
const (
	Att_Begin    = iota
	Att_Str      = iota - 1 // 力量
	Att_Agl                 // 敏捷
	Att_Con                 // 耐力
	Att_Int                 // 智力
	Att_AtkSpeed            // 攻击速度
	Att_End
)

// 二级属性
const (
	AttEx_Begin = iota
	AttEx_MaxHP = iota - 1 // 血量
	AttEx_MaxMP            // 蓝量
	AttEx_Atn              // 物理攻击力
	AttEx_Def              // 物理防御力
	AttEx_Ats              // 魔法攻击力
	AttEx_Adf              // 魔法防御力

	AttEx_CurHP // 当前血量
	AttEx_CurMP // 当前蓝量
	AttEx_End
)

// att entry
type AttEntry struct {
	ID   int32  `json:"_id"`
	Desc string `json:"Desc"`

	Str      int64 `json:"Str"`
	Agl      int64 `json:"Agl"`
	Con      int64 `json:"Con"`
	Int      int64 `json:"Int"`
	AtkSpeed int64 `json:"AtkSpeed"`

	MaxHP int64 `json:"MaxHP"`
	MaxMP int64 `json:"MaxMP"`
	Atn   int64 `json:"Atn"`
	Def   int64 `json:"Def"`
	Ats   int64 `json:"Ats"`
	Adf   int64 `json:"Adf"`
}
