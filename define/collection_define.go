package define

// 收集品品质
const (
	Collection_Quality_Begin  int32 = iota
	Collection_Quality_Green  int32 = iota - 1 // 绿
	Collection_Quality_Blue                    // 蓝
	Collection_Quality_Purple                  // 紫
	Collection_Quality_Yellow                  // 黄
	Collection_Quality_End
)
