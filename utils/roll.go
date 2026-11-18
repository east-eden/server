package utils

import (
	"errors"
	"math/rand"
)

// Roll rand idx
func Roll(weights []int) (idx int) {
	sum := 0
	for _, w := range weights {
		sum += w
	}
	if sum == 0 {
		return 0
	}
	dice := rand.Intn(sum)
	sum = 0
	for i := 0; i < len(weights); i++ {
		if dice >= sum && dice < sum+weights[i] {
			return i
		}
		sum += weights[i]
	}
	return 0
}

// Roll rand idx
func RollInt32(weights []int32) (idx int) {
	sum := int32(0)
	for _, w := range weights {
		sum += w
	}
	if sum == 0 {
		return 0
	}
	dice := int32(rand.Intn(int(sum)))
	sum = 0
	for i := 0; i < len(weights); i++ {
		if dice >= sum && dice < sum+weights[i] {
			return i
		}
		sum += weights[i]
	}
	return 0
}

// RandInt32InRange [a, b)
func RandInt32InRange(a int32, b int32) (c int32) {
	return a + rand.Int31()%max(b-a, 1)
}

// RandInt64InRange [a, b)
func RandInt64InRange(a int64, b int64) (c int64) {
	return a + rand.Int63()%max(b-a, 1)
}

// RandFloat32InRange [a, b)
func RandFloat32InRange(a float32, b float32) (c float32) {
	return a + rand.Float32()*(b-a)
}

func Shuffle(nums []int) []int {
	for i := len(nums); i > 0; i-- {
		last := i - 1
		idx := rand.Intn(i)
		nums[last], nums[idx] = nums[idx], nums[last]
	}
	return nums
}

// 洗牌算法-泛型
func ShuffleT[T any](list []T) []T {
	for i := len(list); i > 0; i-- {
		last := i - 1
		idx := rand.Intn(i)
		list[last], list[idx] = list[idx], list[last]
	}
	return list
}

func RandNorm32(n int) float32 {
	sum := float32(0)
	for i := 0; i < n; i++ {
		sum += rand.Float32()
	}
	return sum / float32(n)
}

func RandNorm64(n int) float64 {
	sum := float64(0)
	for i := 0; i < n; i++ {
		sum += rand.Float64()
	}
	return sum / float64(n)
}

type NumberT interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// 按权重随机一个 返回下标 没找到返回-1
func RandomByWeight[T NumberT](weights []T) int {
	if weights == nil || len(weights) <= 0 {
		return -1
	}
	if len(weights) == 1 {
		return 0
	}
	var sum T = 0
	for _, weight := range weights {
		sum += weight
	}
	if sum == 0 {
		return -1
	}
	rate := float64(rand.Float64()) * float64(sum) // 将 sum 转换为 float64
	for i, weight := range weights {
		if rate < float64(weight) {
			return i
		}
		rate -= float64(weight)
	}
	return -1
}

func RandomByWeightMap[T NumberT](weightMap map[T]T) (T, error) {
	if weightMap == nil || len(weightMap) <= 0 {
		return -1, errors.New("weightMap is empty")
	}
	if len(weightMap) == 1 {
		for k := range weightMap {
			return k, nil
		}
	}
	var sum T = 0
	for _, weight := range weightMap {
		sum += weight
	}
	if sum <= 0 {
		return -1, errors.New("weightMap all weight is zero")
	}
	ids := make([]T, len(weightMap))
	weight := make([]T, len(weightMap))
	idx := 0
	for id, w := range weightMap {
		ids[idx] = id
		weight[idx] = w
		idx++
	}
	rIndex := RandomByWeight(weight)
	if rIndex == -1 {
		return -1, errors.New("weightMap weight is zero")
	}
	return ids[rIndex], nil
}
