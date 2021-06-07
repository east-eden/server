package auto

import (
	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/utils/random"
)

// RandomPicker interface
type CrystalAttRepoList struct {
	list []random.Item
}

func (l *CrystalAttRepoList) GetItemList() []random.Item {
	return l.list
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
		list: make([]random.Item, 0, GetCrystalAttRepoSize()),
	}

	// 所有副属性共用一个属性随机库，晶石位置没有区别
	if tp == define.Crystal_AttTypeVice {
		for _, entry := range crystalAttRepoEntries.Rows {
			if entry.Type == tp {
				ls.list = append(ls.list, entry)
			}
		}
	}

	// 主属性按位置区分属性库
	if tp == define.Crystal_AttTypeMain {
		for _, entry := range crystalAttRepoEntries.Rows {
			if entry.Pos == pos && entry.Type == tp {
				ls.list = append(ls.list, entry)
			}
		}
	}

	return ls
}
