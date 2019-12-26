package talent

import "github.com/yokaiio/yokai_server/internal/define"

type Talent struct {
	ID    int32               `json:"id" bson:"talent_id"`
	entry *define.TalentEntry `bson:"-"`
}
