package auto

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/utils/random"
)

// RandomPicker interface
type CrystalAttRepoList struct {
	list []*CrystalAttRepoEntry
}

func (l *CrystalAttRepoList) GetItemList() []random.Item {
	ret := make([]random.Item, 0, len(l.list))
	for _, v := range l.list {
		ret = append(ret, v)
	}
	return ret
}

// random.Item interface
func (e *CrystalAttRepoEntry) GetId() int {
	return int(e.Id)
}

func (e *CrystalAttRepoEntry) GetWeight() int {
	return int(e.AttWeight)
}

// 获取晶石随机库列表
func GetCrystalAttRepoList(pos int32, tp int32) *CrystalAttRepoList {
	ls := &CrystalAttRepoList{
		list: make([]*CrystalAttRepoEntry, 0, GetCrystalAttRepoSize()),
	}

	// 所有副属性共用一个属性随机库，晶石位置没有区别
	if tp == define.Crystal_AttTypeVice {
		pos = -1
	}

	for _, entry := range crystalAttRepoEntries.Rows {
		if entry.Pos == pos && entry.Type == tp {
			ls.list = append(ls.list, entry)
		}
	}

	return ls
}
