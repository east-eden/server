package cmap

func KeyHashStr(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func KeyHashUint32(key uint32) uint32 {
	return key
}

func KeyHashUint64(key uint64) uint64 {
	return key
}
