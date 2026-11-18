package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/constraints"

	lz4 "github.com/pierrec/lz4/v4"
)

// //////////////////////////////////////////////////////
// exception
func CaptureException(p ...any) {
	if err := recover(); err != nil {
		stack := string(debug.Stack())
		log.Error().Caller(1).Interface("exception_param", p).Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
	}
}

// //////////////////////////////////////////////////////
// Go 启动一个带 panic 保护的 goroutine
func Go(fn func()) {
	go func() {
		defer CaptureException()
		fn()
	}()
}

// //////////////////////////////////////////////////////
// relocate to project root path
func RelocatePath(filter ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("RelocatePath failed: %w", err)
	}

	wd = strings.Replace(wd, "\\", "/", -1)
	fmt.Println("work directory: ", wd)

	var newPath string = wd

	for _, path := range filter {
		n := strings.LastIndex(wd, path)
		if n == -1 {
			continue
		}

		newPath = strings.Join([]string{wd[:n], path}, "")
		err = os.Chdir(newPath)
		if err != nil {
			return err
		}

		break
	}

	fmt.Println("relocate to new_path: ", newPath)
	return nil
}

// //////////////////////////////////////////////////////
// between [a, b)
func Between[T constraints.Ordered](n, a, b T) bool {
	return (n >= a && n < b)
}

// //////////////////////////////////////////////////////
// Round
func Round(val float64) float64 {
	return math.Round((val*10 + 0.1) / 10)
}

// //////////////////////////////////////////////////////
// pack int64
func HighId(id int64) int32 {
	return int32(id >> 32)
}

func LowId(id int64) int32 {
	return int32((id) & 0xFFFFFFFF)
}

func PackId(high, low int32) int64 {
	return (int64(high)<<32 | (int64(low) & 0x00000000FFFFFFFF))
}

// //////////////////////////////////////////////////////
// float valid
func FloatValid[T constraints.Float](n T) bool {
	return n > 0.00001 || n < -0.00001
}

// IsNil 判断接口内部是否为 nil（含 typed nil）
func IsNil(i any) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Interface:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

// JsonStringfy 序列化为 json 字符串，失败时返回空串
func JsonStringfy(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(data)
}

// ToSliceOfAny 把泛型切片转换为 []any
func ToSliceOfAny[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}

	return result
}

// DiffMaps 返回两个 map 的键对称差集
func DiffMaps[K comparable, T any](map1, map2 map[K]T) map[K]T {
	diff := make(map[K]T)

	for k, v := range map1 {
		if _, found := map2[k]; !found {
			diff[k] = v
		}
	}

	for k, v := range map2 {
		if _, found := map1[k]; !found {
			diff[k] = v
		}
	}

	return diff
}

// CompressLZ4 使用 lz4 压缩数据
func CompressLZ4(data []byte) ([]byte, error) {
	var compressedData bytes.Buffer
	lz4Writer := lz4.NewWriter(&compressedData)
	if _, err := lz4Writer.Write(data); err != nil {
		return nil, err
	}

	if err := lz4Writer.Close(); err != nil {
		return nil, err
	}

	return compressedData.Bytes(), nil
}

// DecompressLZ4 解压 lz4 数据
func DecompressLZ4(data []byte) ([]byte, error) {
	var decompressedData bytes.Buffer
	lz4Reader := lz4.NewReader(bytes.NewReader(data))
	if _, err := decompressedData.ReadFrom(lz4Reader); err != nil {
		return nil, err
	}

	return decompressedData.Bytes(), nil
}

// CalcPercentAttribute 把百分比属性折算进对应基础属性。
// kof 的约定：百分比属性区间 [percentBegin, percentEnd) 与基础属性区间一一对应，
// 即 attrType-percentBegin 就是被加成的基础属性 id。
func CalcPercentAttribute(base map[int32]float32, percentAttrBegin int32, percentAttrEnd int32) map[int32]float32 {
	for attrType := percentAttrBegin; attrType < percentAttrEnd; attrType++ {
		if !FloatValid(base[attrType]) {
			continue
		}

		base[attrType-percentAttrBegin] = base[attrType-percentAttrBegin] * (1 + base[attrType])
		base[attrType] = 0
	}

	return base
}

// CalcFinalAttribute 把最终属性折算进对应基础属性（口径同 CalcPercentAttribute）
func CalcFinalAttribute(base map[int32]float32, finalAttrBegin int32, finalAttrEnd int32) map[int32]float32 {
	for attrType := finalAttrBegin; attrType < finalAttrEnd; attrType++ {
		if !FloatValid(base[attrType]) {
			continue
		}

		base[attrType-finalAttrBegin] = base[attrType-finalAttrBegin] * (1 + base[attrType])
		base[attrType] = 0
	}

	return base
}
