package config

// 世界BOSS-拳皇/bossEnum
type DifficultyType int

const (
	DifficultyType_BEGIN     DifficultyType = 1 //枚举开始
	DifficultyType_Normal    DifficultyType = 1 //普通
	DifficultyType_Hard      DifficultyType = 2 //困难
	DifficultyType_Epic      DifficultyType = 3 //史诗
	DifficultyType_Nightmare DifficultyType = 4 //噩梦
	DifficultyType_Hell      DifficultyType = 5 //地狱
	DifficultyType_Inferno   DifficultyType = 6 //炼狱
	DifficultyType_Endless   DifficultyType = 7 //无尽
	DifficultyType_END       DifficultyType = 8 //枚举结束
)

func (e DifficultyType) String() string {

	if e == DifficultyType_Normal {
		return "Normal"
	}
	if e == DifficultyType_Hard {
		return "Hard"
	}
	if e == DifficultyType_Epic {
		return "Epic"
	}
	if e == DifficultyType_Nightmare {
		return "Nightmare"
	}
	if e == DifficultyType_Hell {
		return "Hell"
	}
	if e == DifficultyType_Inferno {
		return "Inferno"
	}
	if e == DifficultyType_Endless {
		return "Endless"
	}

	return ""

}

func GetDifficultyType(tp string) (DifficultyType, bool) {
	switch tp {
	case "Normal":
		return DifficultyType_Normal, true
	case "Hard":
		return DifficultyType_Hard, true
	case "Epic":
		return DifficultyType_Epic, true
	case "Nightmare":
		return DifficultyType_Nightmare, true
	case "Hell":
		return DifficultyType_Hell, true
	case "Inferno":
		return DifficultyType_Inferno, true
	case "Endless":
		return DifficultyType_Endless, true

	}
	return -1, false

}

func GetDifficultyType_ByAlias(tp string) (DifficultyType, bool) {
	switch tp {
	case "普通":
		return DifficultyType_Normal, true
	case "困难":
		return DifficultyType_Hard, true
	case "史诗":
		return DifficultyType_Epic, true
	case "噩梦":
		return DifficultyType_Nightmare, true
	case "地狱":
		return DifficultyType_Hell, true
	case "炼狱":
		return DifficultyType_Inferno, true
	case "无尽":
		return DifficultyType_Endless, true

	}
	return -1, false

}
