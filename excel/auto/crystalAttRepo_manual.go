package auto

import "bitbucket.org/funplus/server/define"

// 获取晶石随机库列表
func GetCrystalAttRepoList(pos define.Crystal_PosType, quality define.ItemQualityType, tp define.Crystal_AttType) []*CrystalAttRepoEntry {
	list := make([]*CrystalAttRepoEntry, 0, GetCrystalAttRepoSize())
	for _, entry := range crystalAttRepoEntries.Rows {
		if entry.Pos == int32(pos) && entry.Quality == int32(quality) && entry.Type == int32(tp) {
			list = append(list, entry)
		}
	}

	return list
}
