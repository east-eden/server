package define

// 关卡类型
const (
	Stage_Type_Begin = iota
	Stage_Type_Main  = iota - 1 // 0 主线关卡
	Stage_Type_Elite            // 1 精英关卡
	Stage_Type_Opera            // 2 剧情关卡
	Stage_Type_End
)

const (
	Chapter_Rewards_Num = 3 // 章节奖励数
	Stage_Objective_Num = 3 // 关卡目标数
)

// 章节信息
type ChapterInfo struct {
	Id      int32                     `bson:"_id" json:"_id"`         // 章节id
	Stars   int32                     `bson:"stars" json:"stars"`     // 当前星数
	Rewards [Chapter_Rewards_Num]bool `bson:"rewards" json:"rewards"` // 是否已领取奖励
}

// 关卡信息
type StageInfo struct {
	Id             int32                     `bson:"_id" json:"_id"`                         // 关卡id
	ChallengeTimes int16                     `bson:"challenge_times" json:"challenge_times"` // 已挑战次数
	FirstReward    bool                      `bson:"first_reward" json:"first_reward"`       // 是否已获得首次通关奖励
	Objectives     [Stage_Objective_Num]bool `bson:"objectives" json:"objectives"`           // 目标是否达成
}
