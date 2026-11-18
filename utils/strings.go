package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

func String2ByteSlice(str string) []byte {
	if str == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func ByteSlice2String(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

// SimpleTitle 将字符串首字母转为大写，其余转为小写（Unicode 安全版）
func SimpleTitle(s string) string {
	if len(s) == 0 {
		return s
	}

	// 优化点：直接操作字符而非字符串，避免多次内存分配
	builder := strings.Builder{}
	builder.Grow(len(s)) // 预分配内存

	// 处理首字符
	first, remaining := splitFirstRune(s)
	if first != nil {
		builder.WriteRune(unicode.ToUpper(*first))
	}

	// 处理剩余字符
	for _, r := range remaining {
		builder.WriteRune(r)
	}

	return builder.String()
}

func SimpleUntitle(s string) string {
	if len(s) == 0 {
		return s
	}

	// 优化点：直接操作字符而非字符串，避免多次内存分配
	builder := strings.Builder{}
	builder.Grow(len(s)) // 预分配内存

	// 处理首字符
	first, remaining := splitFirstRune(s)
	if first != nil {
		builder.WriteRune(unicode.ToLower(*first))
	}

	// 处理剩余字符
	for _, r := range remaining {
		builder.WriteRune(r)
	}

	return builder.String()
}

// 分离字符串的首字符和剩余部分（Unicode 安全）
func splitFirstRune(s string) (first *rune, remaining []rune) {
	if s == "" {
		return nil, nil
	}

	runes := []rune(s)
	if len(runes) > 0 {
		return &runes[0], runes[1:]
	}
	return nil, nil
}

type IntKey interface {
	~int | ~int32 | ~int64
}

// JoinMapKeys 泛型函数：将任意 map[K]V（K 是 IntKey）
// 的所有 key 提取、升序后拼成 "k1|k2|k3" 格式
func JoinMapKeys[K IntKey, V any](m map[K]V) string {
	// 1. 收集并转成 int64 便于排序
	keys := make([]int64, 0, len(m))
	for k := range m {
		keys = append(keys, int64(k))
	}

	// 2. 升序排序
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// 3. 转为字符串切片
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = strconv.FormatInt(k, 10)
	}

	// 4. 用 "|" 连接并返回
	return strings.Join(parts, "|")
}

func String2IntSlice(s string, sep string) ([]int32, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, sep)
	result := make([]int32, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue // 可选：也可以选择 return error here
		}
		val, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid int value: '%s'", p)
		}
		result = append(result, int32(val))
	}
	return result, nil
}
