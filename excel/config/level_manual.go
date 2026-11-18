package config

// levelByChapterStage 章节/节序号 -> 关卡 id 的索引，供关卡前置关系判断使用
var levelByChapterStage map[[2]int32]int32

func buildLevelIndex() {
	if levelByChapterStage != nil {
		return
	}

	// 配置尚未加载时退化为空索引，避免调用方拿到 nil 指针
	if levelEntries == nil {
		levelByChapterStage = make(map[[2]int32]int32)
		return
	}

	m := make(map[[2]int32]int32, len(levelEntries.Arrays))
	for _, e := range levelEntries.Arrays {
		m[[2]int32{e.Chapter, e.Stage}] = e.Id
	}

	levelByChapterStage = m
}

// GetLevelIdByChapterStage 按章节与节序号查关卡 id
func GetLevelIdByChapterStage(chapter, stage int32) (int32, bool) {
	buildLevelIndex()
	id, ok := levelByChapterStage[[2]int32{chapter, stage}]
	return id, ok
}

// ResetLevelIndex 配置重载后重建索引
func ResetLevelIndex() {
	levelByChapterStage = nil
}
