package utils

import (
	"fmt"
)

// MultiValueMap 多值映射，一个 key 对应多个 value
// Not Thread-Safe
type MultiValueMap[K comparable, V any] struct {
	data map[K][]V
}

// NewMultiValueMap 创建新的 MultiValueMap
func NewMultiValueMap[K comparable, V any]() *MultiValueMap[K, V] {
	return &MultiValueMap[K, V]{
		data: make(map[K][]V),
	}
}

// NewMultiValueMapWithCapacity 创建带容量的 MultiValueMap
func NewMultiValueMapWithCapacity[K comparable, V any](capacity int) *MultiValueMap[K, V] {
	return &MultiValueMap[K, V]{
		data: make(map[K][]V, capacity),
	}
}

// ===== 核心方法 =====

// Add 添加一个值到指定的 key
func (m *MultiValueMap[K, V]) Add(key K, value V) {
	m.data[key] = append(m.data[key], value)
}

// AddAll 添加多个值到指定的 key
func (m *MultiValueMap[K, V]) AddAll(key K, values ...V) {
	m.data[key] = append(m.data[key], values...)
}

// AddAllMap 将另一个 MultiValueMap 的所有数据合并进来
func (m *MultiValueMap[K, V]) AddAllMap(other *MultiValueMap[K, V]) {
	for key, values := range other.data {
		m.data[key] = append(m.data[key], values...)
	}
}

// Get 获取指定 key 的所有元素（返回切片副本）
func (m *MultiValueMap[K, V]) Get(key K) []V {
	if values, exists := m.data[key]; exists {
		// 返回副本，防止外部修改内部数据
		result := make([]V, len(values))
		copy(result, values)
		return result
	}
	return []V{} // 返回空切片，不是 nil
}

// GetByIndex 获取指定 key 的第 index 个元素（从 0 开始）
func (m *MultiValueMap[K, V]) GetByIndex(key K, index int) (V, error) {
	values, exists := m.data[key]
	if !exists {
		var zero V
		return zero, fmt.Errorf("key '%v' does not exist", key)
	}

	if index < 0 || index >= len(values) {
		var zero V
		return zero, fmt.Errorf("index %d out of bounds, length: %d", index, len(values))
	}

	return values[index], nil
}

// GetFirst 获取指定 key 的第一个元素
func (m *MultiValueMap[K, V]) GetFirst(key K) (V, error) {
	return m.GetByIndex(key, 0)
}

// GetLast 获取指定 key 的最后一个元素
func (m *MultiValueMap[K, V]) GetLast(key K) (V, error) {
	values, exists := m.data[key]
	if !exists {
		var zero V
		return zero, fmt.Errorf("key '%v' does not exist", key)
	}
	return m.GetByIndex(key, len(values)-1)
}

// ===== 查询方法 =====

// ContainsKey 检查是否包含某个 key
func (m *MultiValueMap[K, V]) ContainsKey(key K) bool {
	_, exists := m.data[key]
	return exists
}

// ContainsValue 检查是否包含某个 value（在任意 key 下）
func (m *MultiValueMap[K, V]) ContainsValue(value V) bool {
	for _, values := range m.data {
		for _, v := range values {
			if any(v) == any(value) {
				return true
			}
		}
	}
	return false
}

// ContainsEntry 检查是否包含指定的 key-value 对
func (m *MultiValueMap[K, V]) ContainsEntry(key K, value V) bool {
	values, exists := m.data[key]
	if !exists {
		return false
	}
	for _, v := range values {
		if any(v) == any(value) {
			return true
		}
	}
	return false
}

// Size 获取所有 key 的 value 总数
func (m *MultiValueMap[K, V]) Size() int {
	total := 0
	for _, values := range m.data {
		total += len(values)
	}
	return total
}

// KeySize 获取指定 key 的元素数量
func (m *MultiValueMap[K, V]) KeySize(key K) int {
	return len(m.data[key])
}

// Keys 获取所有 key
func (m *MultiValueMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有 value（展平）
func (m *MultiValueMap[K, V]) Values() []V {
	total := m.Size()
	result := make([]V, 0, total)
	for _, values := range m.data {
		result = append(result, values...)
	}
	return result
}

// ===== 修改方法 =====

// Remove 移除指定 key 的所有元素
func (m *MultiValueMap[K, V]) Remove(key K) []V {
	if values, exists := m.data[key]; exists {
		delete(m.data, key)
		return values
	}
	return []V{}
}

// RemoveEntry 移除指定 key-value 对（只移除第一个匹配的）
func (m *MultiValueMap[K, V]) RemoveEntry(key K, value V) bool {
	values, exists := m.data[key]
	if !exists {
		return false
	}

	for i, v := range values {
		if any(v) == any(value) {
			// 移除该元素
			m.data[key] = append(values[:i], values[i+1:]...)
			// 如果切片变为空，删除 key
			if len(m.data[key]) == 0 {
				delete(m.data, key)
			}
			return true
		}
	}
	return false
}

// RemoveAllEntries 移除所有匹配的 key-value 对
func (m *MultiValueMap[K, V]) RemoveAllEntries(key K, value V) int {
	values, exists := m.data[key]
	if !exists {
		return 0
	}

	count := 0
	newValues := make([]V, 0, len(values))
	for _, v := range values {
		if any(v) == any(value) {
			count++
		} else {
			newValues = append(newValues, v)
		}
	}

	if len(newValues) == 0 {
		delete(m.data, key)
	} else {
		m.data[key] = newValues
	}
	return count
}

// ReplaceValues 替换指定 key 的所有值
func (m *MultiValueMap[K, V]) ReplaceValues(key K, values ...V) {
	if len(values) == 0 {
		delete(m.data, key)
	} else {
		m.data[key] = values
	}
}

// Clear 清空所有数据
func (m *MultiValueMap[K, V]) Clear() {
	m.data = make(map[K][]V)
}

// ===== 遍历方法 =====

// ForEachFlatten 遍历所有 key-value 对
func (m *MultiValueMap[K, V]) ForEachFlatten(fn func(K, V)) {
	for key, values := range m.data {
		for _, value := range values {
			fn(key, value)
		}
	}
}

// ForEach 遍历所有 key-[]value 对
func (m *MultiValueMap[K, V]) ForEach(fn func(K, []V)) {
	for key, values := range m.data {
		fn(key, values)
	}
}

// ForEachKey 遍历指定 key 的所有 value
func (m *MultiValueMap[K, V]) ForEachKey(key K, fn func(V)) {
	if values, exists := m.data[key]; exists {
		for _, value := range values {
			fn(value)
		}
	}
}
