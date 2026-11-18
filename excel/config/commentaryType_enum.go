package config

// 竞技场-拳皇/CommentaryType
type CommentaryType int

const (
	CommentaryType_BEGIN            CommentaryType = 1 //枚举开始
	CommentaryType_Start            CommentaryType = 1 //开场
	CommentaryType_SkillTrigger     CommentaryType = 2 //技能触发
	CommentaryType_SituationChange1 CommentaryType = 3 //形势变化1
	CommentaryType_SituationChange2 CommentaryType = 4 //形势变化2
	CommentaryType_SituationChange3 CommentaryType = 5 //形势变化3
	CommentaryType_End              CommentaryType = 6 //结束
	CommentaryType_END              CommentaryType = 7 //枚举结束
)

// 竞技场-拳皇/CommentaryType
type CommentatorSpineType int

const (
	CommentatorSpineType_BEGIN   CommentatorSpineType = 1 //枚举开始
	CommentatorSpineType_Nervous CommentatorSpineType = 1 //紧张
	CommentatorSpineType_Excited CommentatorSpineType = 2 //兴奋
	CommentatorSpineType_Afraid  CommentatorSpineType = 3 //惊恐
	CommentatorSpineType_Happy   CommentatorSpineType = 4 //高兴
	CommentatorSpineType_END     CommentatorSpineType = 5 //枚举结束
)

func (e CommentaryType) String() string {

	if e == CommentaryType_Start {
		return "Start"
	}
	if e == CommentaryType_SkillTrigger {
		return "SkillTrigger"
	}
	if e == CommentaryType_SituationChange1 {
		return "SituationChange1"
	}
	if e == CommentaryType_SituationChange2 {
		return "SituationChange2"
	}
	if e == CommentaryType_SituationChange3 {
		return "SituationChange3"
	}
	if e == CommentaryType_End {
		return "End"
	}

	return ""

}

func GetCommentaryType(tp string) (CommentaryType, bool) {
	switch tp {
	case "Start":
		return CommentaryType_Start, true
	case "SkillTrigger":
		return CommentaryType_SkillTrigger, true
	case "SituationChange1":
		return CommentaryType_SituationChange1, true
	case "SituationChange2":
		return CommentaryType_SituationChange2, true
	case "SituationChange3":
		return CommentaryType_SituationChange3, true
	case "End":
		return CommentaryType_End, true

	}
	return -1, false

}

func GetCommentaryType_ByAlias(tp string) (CommentaryType, bool) {
	switch tp {
	case "开场":
		return CommentaryType_Start, true
	case "技能触发":
		return CommentaryType_SkillTrigger, true
	case "形势变化1":
		return CommentaryType_SituationChange1, true
	case "形势变化2":
		return CommentaryType_SituationChange2, true
	case "形势变化3":
		return CommentaryType_SituationChange3, true
	case "结束":
		return CommentaryType_End, true

	}
	return -1, false

}

func (e CommentatorSpineType) String() string {

	if e == CommentatorSpineType_Nervous {
		return "Nervous"
	}
	if e == CommentatorSpineType_Excited {
		return "Excited"
	}
	if e == CommentatorSpineType_Afraid {
		return "Afraid"
	}
	if e == CommentatorSpineType_Happy {
		return "Happy"
	}

	return ""

}

func GetCommentatorSpineType(tp string) (CommentatorSpineType, bool) {
	switch tp {
	case "Nervous":
		return CommentatorSpineType_Nervous, true
	case "Excited":
		return CommentatorSpineType_Excited, true
	case "Afraid":
		return CommentatorSpineType_Afraid, true
	case "Happy":
		return CommentatorSpineType_Happy, true

	}
	return -1, false

}

func GetCommentatorSpineType_ByAlias(tp string) (CommentatorSpineType, bool) {
	switch tp {
	case "紧张":
		return CommentatorSpineType_Nervous, true
	case "兴奋":
		return CommentatorSpineType_Excited, true
	case "惊恐":
		return CommentatorSpineType_Afraid, true
	case "高兴":
		return CommentatorSpineType_Happy, true

	}
	return -1, false

}
