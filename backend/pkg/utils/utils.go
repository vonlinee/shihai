package utils

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

// GetValue 从 map 中按路径取值，支持点号分隔的字段和数组索引。
// 例如 address.city、users[0].name 和 matrix[0][1]。
func GetValue(data map[string]interface{}, path string) interface{} {
	value, _ := lookupValue(data, path)
	return value
}

// GetString 按路径获取字符串；路径不存在或值类型不匹配时返回空字符串。
func GetString(data map[string]interface{}, path string) string {
	value, _ := lookupValue(data, path)
	result, _ := value.(string)
	return result
}

// GetInt 按路径获取整数；路径不存在、值类型不匹配或数字无效时返回零。
func GetInt(data map[string]interface{}, path string) int {
	value, _ := lookupValue(data, path)
	return valueAsInt(value)
}

// GetFloat 按路径获取浮点数；路径不存在、值类型不匹配或数字无效时返回零。
func GetFloat(data map[string]interface{}, path string) float64 {
	value, _ := lookupValue(data, path)
	return valueAsFloat64(value)
}

// GetBool 按路径获取布尔值；路径不存在或值类型不匹配时返回 false。
func GetBool(data map[string]interface{}, path string) bool {
	value, _ := lookupValue(data, path)
	result, _ := value.(bool)
	return result
}

// GetSlice 按路径获取 JSON 数组；路径不存在或值类型不匹配时返回 nil。
func GetSlice(data map[string]interface{}, path string) []interface{} {
	value, _ := lookupValue(data, path)
	result, _ := value.([]interface{})
	return result
}

// GetMap 按路径获取 JSON 对象；路径不存在或值类型不匹配时返回 nil。
func GetMap(data map[string]interface{}, path string) map[string]interface{} {
	value, _ := lookupValue(data, path)
	return asStringMap(value)
}

// GetStringOrDefault 按路径获取字符串；仅在路径不存在时返回默认值。
func GetStringOrDefault(data map[string]interface{}, path string, defaultVal string) string {
	value, exists := lookupValue(data, path)
	if !exists {
		return defaultVal
	}
	result, _ := value.(string)
	return result
}

// GetIntOrDefault 按路径获取整数；仅在路径不存在时返回默认值。
func GetIntOrDefault(data map[string]interface{}, path string, defaultVal int) int {
	value, exists := lookupValue(data, path)
	if !exists {
		return defaultVal
	}
	return valueAsInt(value)
}

func lookupValue(data interface{}, path string) (interface{}, bool) {
	if path == "" {
		return data, true
	}

	current := data
	for _, part := range strings.Split(path, ".") {
		key, indexes, ok := parsePathPart(part)
		if !ok {
			return nil, false
		}

		currentMap := asStringMap(current)
		if currentMap == nil {
			return nil, false
		}
		current, ok = currentMap[key]
		if !ok {
			return nil, false
		}

		for _, index := range indexes {
			array, ok := current.([]interface{})
			if !ok || index >= len(array) {
				return nil, false
			}
			current = array[index]
		}
	}

	return current, true
}

func parsePathPart(part string) (string, []int, bool) {
	if part == "" {
		return "", nil, false
	}

	openBracket := strings.IndexByte(part, '[')
	if openBracket == -1 {
		if strings.ContainsRune(part, ']') {
			return "", nil, false
		}
		return part, nil, true
	}
	if openBracket == 0 {
		return "", nil, false
	}

	key := part[:openBracket]
	remainder := part[openBracket:]
	indexes := make([]int, 0, strings.Count(remainder, "["))
	for remainder != "" {
		if remainder[0] != '[' {
			return "", nil, false
		}
		closeBracket := strings.IndexByte(remainder, ']')
		if closeBracket <= 1 {
			return "", nil, false
		}

		indexText := remainder[1:closeBracket]
		if strings.IndexFunc(indexText, func(r rune) bool { return r < '0' || r > '9' }) != -1 {
			return "", nil, false
		}
		index, err := strconv.Atoi(indexText)
		if err != nil {
			return "", nil, false
		}
		indexes = append(indexes, index)
		remainder = remainder[closeBracket+1:]
	}

	return key, indexes, true
}

func asStringMap(value interface{}) map[string]interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return typed
	case JSONMap:
		return map[string]interface{}(typed)
	default:
		return nil
	}
}

func valueAsInt(value interface{}) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		if strconv.IntSize == 32 {
			return clampFloatToInt(float64(typed))
		}
		return int(typed)
	case float64:
		return clampFloatToInt(typed)
	case json.Number:
		if integer, err := typed.Int64(); err == nil {
			return valueAsInt(integer)
		}
		if floating, err := typed.Float64(); err == nil {
			return clampFloatToInt(floating)
		}
	}
	return 0
}

func clampFloatToInt(value float64) int {
	if math.IsNaN(value) {
		return 0
	}
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	if value >= float64(maxInt) {
		return maxInt
	}
	if value <= float64(minInt) {
		return minInt
	}
	return int(value)
}

func valueAsInt64(value interface{}) int64 {
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int64:
		return typed
	case float64:
		if math.IsNaN(typed) {
			return 0
		}
		if typed >= float64(math.MaxInt64) {
			return math.MaxInt64
		}
		if typed <= float64(math.MinInt64) {
			return math.MinInt64
		}
		return int64(typed)
	case json.Number:
		if integer, err := typed.Int64(); err == nil {
			return integer
		}
		if floating, err := typed.Float64(); err == nil {
			return valueAsInt64(floating)
		}
	}
	return 0
}

func valueAsFloat64(value interface{}) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case json.Number:
		result, _ := typed.Float64()
		return result
	default:
		return 0
	}
}

// SliceToMapByKey 通用转换：整个结构体作为 Value
// 用法：SliceToMapByKey(users, func(u User) int { return u.ID })
func SliceToMapByKey[T any, K comparable](items []T, keyFn func(T) K) map[K]*T {
	m := make(map[K]*T, len(items))
	for _, item := range items {
		m[keyFn(item)] = &item
	}
	return m
}

// SliceToMapByKeyValue 通用转换：自定义 Value
// 用法：SliceToMapByKeyValue(users, func(u User) int { return u.ID }, func(u User) string { return u.Name })
func SliceToMapByKeyValue[T any, K comparable, V any](items []T, keyFn func(T) K, valFn func(T) V) map[K]V {
	m := make(map[K]V, len(items))
	for _, item := range items {
		m[keyFn(item)] = valFn(item)
	}
	return m
}

// IsEmpty 判断切片是否为空（nil 或长度为 0）
func IsEmpty[T any](items []T) bool {
	return items == nil || len(items) == 0
}

// IsNotEmpty 判断切片是否非空
func IsNotEmpty[T any](items []T) bool {
	return items != nil && len(items) > 0
}

// Count 统计切片中符合条件的元素数量
func Count[T any](items []T, predicate func(T) bool) int {
	count := 0
	for _, item := range items {
		if predicate(item) {
			count++
		}
	}
	return count
}

// Exists 判断切片中符合条件的元素
func Exists[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Keys 获取 map 的所有 key
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ExtractToSlice 提取字段
func ExtractToSlice[T any, R any](items []T, fieldFn func(T) R) []R {
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, fieldFn(item))
	}
	return result
}

// ExtractToSliceUnique 提取并去重
func ExtractToSliceUnique[T any, R comparable](items []T, fieldFn func(T) R) []R {
	seen := make(map[R]struct{})
	result := make([]R, 0, len(items))
	for _, item := range items {
		val := fieldFn(item)
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}
	return result
}

// FlatToSlice 将结构体数组中指定字段（数组类型）的所有元素展平到一个切片
func FlatToSlice[T any, R any](items []T, fieldFn func(T) []R) []R {
	// 先计算总长度
	totalLen := 0
	for _, item := range items {
		totalLen += len(fieldFn(item))
	}

	result := make([]R, 0, totalLen)
	for _, item := range items {
		result = append(result, fieldFn(item)...)
	}
	return result
}

// FlatToSliceUnique 展平并去重
func FlatToSliceUnique[T any, R comparable](items []T, fieldFn func(T) []R) []R {
	seen := make(map[R]struct{})
	result := make([]R, 0)

	for _, item := range items {
		for _, val := range fieldFn(item) {
			if _, exists := seen[val]; !exists {
				seen[val] = struct{}{}
				result = append(result, val)
			}
		}
	}
	return result
}

// FlatToSliceIf 展平并过滤（根据元素值过滤）
func FlatToSliceIf[T any, R any](items []T, fieldFn func(T) []R, predicate func(R) bool) []R {
	result := make([]R, 0)
	for _, item := range items {
		for _, val := range fieldFn(item) {
			if predicate(val) {
				result = append(result, val)
			}
		}
	}
	return result
}

// FlatToSliceWhere 展平并过滤（根据结构体条件过滤）
func FlatToSliceWhere[T any, R any](items []T, fieldFn func(T) []R, condition func(T) bool) []R {
	result := make([]R, 0)
	for _, item := range items {
		if condition(item) {
			result = append(result, fieldFn(item)...)
		}
	}
	return result
}

// FlatToSliceWithCount 展平并返回每个元素出现的次数
func FlatToSliceWithCount[T any, R comparable](items []T, fieldFn func(T) []R) map[R]int {
	countMap := make(map[R]int)
	for _, item := range items {
		for _, val := range fieldFn(item) {
			countMap[val]++
		}
	}
	return countMap
}

// FlatToMapByKey 展平并分组（按指定的 key）
func FlatToMapByKey[T any, R any, K comparable](items []T, keyFn func(T) K, fieldFn func(T) []R) map[K][]R {
	result := make(map[K][]R)
	for _, item := range items {
		key := keyFn(item)
		result[key] = append(result[key], fieldFn(item)...)
	}
	return result
}
