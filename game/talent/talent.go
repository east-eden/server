package talent

import "github.com/yokaiio/yokai_server/internal/define"

type Talent struct {
	ID    int32 `json:"id"`
	entry *define.TalentEntry
}
