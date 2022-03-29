package define

// type Number int32

const (
	PercentBase            = 10000 // 百分比基数
	ConsistentNodeReplicas = 256   // 一致性哈希预节点数
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 |
		~int64 | ~uint | ~uint8 | ~uint16 |
		~uint32 | ~uint64
}

type Number interface {
	Integer | ~float32 | ~float64
}
