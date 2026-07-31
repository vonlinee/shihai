package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// JSONMap 是 map[string]interface{} 的别名，方便扩展方法
type JSONMap map[string]interface{}

// Parse 将 JSON 字符串解析为 JSONMap
func Parse(jsonStr string) (JSONMap, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("JSON 字符串为空")
	}

	result, err := decodeJSONMap(strings.NewReader(jsonStr), false)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return result, nil
}

// MustParse 解析 JSON，失败时 panic
func MustParse(jsonStr string) JSONMap {
	result, err := Parse(jsonStr)
	if err != nil {
		panic(err)
	}
	return result
}

// ParseBytes 从字节数组解析 JSON
func ParseBytes(data []byte) (JSONMap, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("数据为空")
	}

	result, err := decodeJSONMap(strings.NewReader(string(data)), false)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return result, nil
}

// ParseReader 从 io.Reader 解析 JSON
func ParseReader(reader io.Reader) (JSONMap, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader 为空")
	}

	result, err := decodeJSONMap(reader, false)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return result, nil
}

// ParseWithNumber 解析 JSON，保留数字精度（使用 json.Number）
func ParseWithNumber(jsonStr string) (JSONMap, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("JSON 字符串为空")
	}

	result, err := decodeJSONMap(strings.NewReader(jsonStr), true)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return result, nil
}

// ParseArray 解析 JSON 数组为 []JSONMap
func ParseArray(jsonStr string) ([]JSONMap, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("JSON 字符串为空")
	}

	var result []JSONMap
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 数组失败: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("解析 JSON 数组失败: 顶层值必须是数组")
	}
	return result, nil
}

func decodeJSONMap(reader io.Reader, useNumber bool) (JSONMap, error) {
	decoder := json.NewDecoder(reader)
	if useNumber {
		decoder.UseNumber()
	}

	var result JSONMap
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("顶层值必须是对象")
	}

	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("只能包含一个 JSON 值")
		}
		return nil, err
	}
	return result, nil
}

// Get 按路径获取值，支持 users[0].name 格式；可选参数用于指定路径不存在时的默认值。
func (m JSONMap) Get(path string, defaultValue ...interface{}) interface{} {
	value, exists := lookupValue(m, path)
	if !exists && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetString 获取字符串
func (m JSONMap) GetString(path string) string {
	value, _ := lookupValue(m, path)
	result, _ := value.(string)
	return result
}

// GetStringOrDefault 获取字符串，带默认值
func (m JSONMap) GetStringOrDefault(path string, defaultVal string) string {
	value, exists := lookupValue(m, path)
	if !exists {
		return defaultVal
	}
	result, _ := value.(string)
	return result
}

// GetInt 获取整数
func (m JSONMap) GetInt(path string) int {
	value, _ := lookupValue(m, path)
	return valueAsInt(value)
}

// GetIntOrDefault 获取整数，带默认值
func (m JSONMap) GetIntOrDefault(path string, defaultVal int) int {
	value, exists := lookupValue(m, path)
	if !exists {
		return defaultVal
	}
	return valueAsInt(value)
}

// GetInt64 获取 int64
func (m JSONMap) GetInt64(path string) int64 {
	value, _ := lookupValue(m, path)
	return valueAsInt64(value)
}

// GetFloat64 获取 float64
func (m JSONMap) GetFloat64(path string) float64 {
	value, _ := lookupValue(m, path)
	return valueAsFloat64(value)
}

// GetBool 获取布尔值
func (m JSONMap) GetBool(path string) bool {
	value, _ := lookupValue(m, path)
	result, _ := value.(bool)
	return result
}

// GetSlice 获取切片
func (m JSONMap) GetSlice(path string) []interface{} {
	value, _ := lookupValue(m, path)
	result, _ := value.([]interface{})
	return result
}

// GetMap 获取 map
func (m JSONMap) GetMap(path string) JSONMap {
	value, _ := lookupValue(m, path)
	return JSONMap(asStringMap(value))
}

// Exists 检查路径是否存在
func (m JSONMap) Exists(path string) bool {
	_, exists := lookupValue(m, path)
	return exists
}

// Keys 获取所有顶层键
func (m JSONMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ToJSON 转换为 JSON 字符串
func (m JSONMap) ToJSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("序列化 JSON 失败: %w", err)
	}
	return string(data), nil
}

// ToJSONPretty 转换为格式化 JSON 字符串
func (m JSONMap) ToJSONPretty() (string, error) {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化 JSON 失败: %w", err)
	}
	return string(data), nil
}

// Merge 合并另一个 JSONMap，返回新的 JSONMap
func (m JSONMap) Merge(other JSONMap) JSONMap {
	result := make(JSONMap)
	for k, v := range m {
		result[k] = v
	}
	for k, v := range other {
		result[k] = v
	}
	return result
}

// DeepMerge 深度合并另一个 JSONMap
func (m JSONMap) DeepMerge(other JSONMap) JSONMap {
	result := cloneJSONMap(m)
	for k, v := range other {
		if existing, exists := result[k]; exists {
			existingMap := asStringMap(existing)
			otherMap := asStringMap(v)
			if existingMap != nil && otherMap != nil {
				result[k] = JSONMap(existingMap).DeepMerge(JSONMap(otherMap))
				continue
			}
		}
		result[k] = cloneJSONValue(v)
	}
	return result
}

func cloneJSONMap(source JSONMap) JSONMap {
	result := make(JSONMap, len(source))
	for key, value := range source {
		result[key] = cloneJSONValue(value)
	}
	return result
}

func cloneJSONValue(value interface{}) interface{} {
	if mapValue := asStringMap(value); mapValue != nil {
		return cloneJSONMap(JSONMap(mapValue))
	}
	if arrayValue, ok := value.([]interface{}); ok {
		result := make([]interface{}, len(arrayValue))
		for index, item := range arrayValue {
			result[index] = cloneJSONValue(item)
		}
		return result
	}
	return value
}

// IsJSON 检查字符串是否为有效的 JSON
func IsJSON(jsonStr string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// IsJSONObject 检查字符串是否为 JSON 对象
func IsJSONObject(jsonStr string) bool {
	var m map[string]interface{}
	return json.Unmarshal([]byte(jsonStr), &m) == nil && m != nil
}

// IsJSONArray 检查字符串是否为 JSON 数组
func IsJSONArray(jsonStr string) bool {
	var arr []interface{}
	return json.Unmarshal([]byte(jsonStr), &arr) == nil && arr != nil
}
