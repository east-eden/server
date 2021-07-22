package define

import (
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
)

// 排行榜元数据
type RankRaw struct {
	ObjId  int64   `json:"_id" bson:"_id"`         // 排行榜数据key
	RankId int32   `json:"rank_id" bson:"rank_id"` // 排行榜id
	Score  float64 `json:"score" bson:"score"`     // 排行榜得分
	Date   int32   `json:"date" bson:"date"`       // 分数更新时间
}

func (r *RankRaw) FromPB(pb *pbGlobal.RankRaw) {
	r.ObjId = pb.GetObjId()
	r.Score = pb.GetScore()
	r.Date = pb.GetDate()
}

func (r *RankRaw) ToPB() *pbGlobal.RankRaw {
	pb := &pbGlobal.RankRaw{
		ObjId: r.ObjId,
		Score: r.Score,
		Date:  r.Date,
	}
	return pb
}
