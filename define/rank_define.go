package define

import (
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
)

const (
	RankId_Begin             = iota
	RankId_LocalPlayerLevel  // 本服玩家等级榜
	RankId_GlobalPlayerLevel // 全服玩家等级榜
	RankId_End
)

// 排行榜元数据
type RankRaw struct {
	ObjId   int64   `json:"_id" bson:"_id"`         // 排行榜数据key -- 玩家id
	ObjName string  `json:"name" bson:"name"`       // 排行数据名字 -- 玩家名字
	RankId  int32   `json:"rank_id" bson:"rank_id"` // 排行榜id
	Score   float64 `json:"score" bson:"score"`     // 排行榜得分
	Date    int32   `json:"date" bson:"date"`       // 分数更新时间
}

func (r *RankRaw) FromPB(pb *pbGlobal.RankRaw) {
	r.ObjId = pb.GetObjId()
	r.ObjName = pb.GetObjName()
	r.Score = pb.GetScore()
	r.Date = pb.GetDate()
}

func (r *RankRaw) ToPB() *pbGlobal.RankRaw {
	pb := &pbGlobal.RankRaw{
		ObjId:   r.ObjId,
		ObjName: r.ObjName,
		Score:   r.Score,
		Date:    r.Date,
	}
	return pb
}
