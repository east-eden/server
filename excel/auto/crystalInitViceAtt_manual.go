package auto

import (
	"errors"

	"bitbucket.org/funplus/server/utils"
	"bitbucket.org/funplus/server/utils/random"
	"github.com/rs/zerolog/log"
)

// RandomPicker interface
func (e *CrystalInitViceAttEntry) GetItemList() []random.Item {
	ret := make([]random.Item, 0, len(e.AttNumWeight))
	for k, v := range e.AttNumWeight {
		ret = append(ret, &CrystalInitViceAttItem{
			AttNum:    k + 1,  // 副属性条数
			AttWeight: int(v), // 副属性条数权重
		})
	}
	return ret
}

// random.Item interface
type CrystalInitViceAttItem struct {
	AttNum    int
	AttWeight int
}

func (i *CrystalInitViceAttItem) GetId() int {
	return i.AttNum
}

func (i *CrystalInitViceAttItem) GetWeight() int {
	return i.AttWeight
}

// 获取副属性随机条数
func GetCrystalInitViceAttNum(quality int32) int {
	entry, ok := GetCrystalInitViceAttEntry(quality)
	if !ok {
		log.Error().Int32("quality", quality).Msg("GetCrystalInitViceAttEntry failed")
		return 0
	}

	item, err := random.PickOne(entry, func(random.Item) bool {
		return true
	})

	if errors.Is(err, random.ErrNoResult) {
		return 0
	}

	if pass := utils.ErrCheck(err, "GetCrystalInitViceAttNum failed", quality); !pass {
		return 0
	}

	return item.GetId()
}
