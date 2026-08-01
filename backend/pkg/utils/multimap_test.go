package utils

import (
	"fmt"
	"testing"
)

// ===== 使用示例 =====

func TestNewMultiValueMap(t *testing.T) {
	// 1. 创建 MultiValueMap
	m := NewMultiValueMap[string, string]()

	// 2. 添加数据
	m.Add("fruit", "apple")
	m.Add("fruit", "banana")
	m.Add("fruit", "apple") // 允许重复
	m.Add("fruit", "orange")
	m.Add("vegetable", "carrot")
	m.Add("vegetable", "broccoli")

	// 3. 获取所有元素
	fruits := m.Get("fruit")
	fmt.Println("All fruits:", fruits) // [apple banana apple orange]

	// 4. GetByIndex 方法
	firstFruit, err := m.GetByIndex("fruit", 0)
	if err == nil {
		fmt.Println("First fruit:", firstFruit) // apple
	}

	secondFruit, err := m.GetByIndex("fruit", 1)
	if err == nil {
		fmt.Println("Second fruit:", secondFruit) // banana
	}

	// 索引越界
	_, err = m.GetByIndex("fruit", 10)
	fmt.Println("Error:", err) // index 10 out of bounds, length: 4

	// 5. GetFirst 和 GetLast
	first, _ := m.GetFirst("fruit")
	last, _ := m.GetLast("fruit")
	fmt.Println("First:", first, "Last:", last) // First: apple Last: orange

	// 6. 获取不存在的 key
	empty := m.Get("unknown")
	fmt.Println("Unknown key:", empty) // []

	_, err = m.GetFirst("unknown")
	fmt.Println("Error:", err) // key 'unknown' does not exist

	// 7. 查看所有内容
	fmt.Println("\n所有数据:")
	m.ForEachFlatten(func(key string, value string) {
		fmt.Printf("  %v -> %v\n", key, value)
	})

	// 8. 统计信息
	fmt.Println("\n统计信息:")
	fmt.Println("总元素数:", m.Size())                 // 6
	fmt.Println("Key 数量:", len(m.Keys()))          // 2
	fmt.Println("fruit 的元素数:", m.KeySize("fruit")) // 4
	fmt.Println("所有 key:", m.Keys())               // [fruit vegetable]
	fmt.Println("所有 value:", m.Values())           // [apple banana apple orange carrot broccoli]

	// 9. 删除操作
	fmt.Println("\n删除操作:")
	m.RemoveEntry("fruit", "apple")
	fmt.Println("删除一个 apple 后:", m.Get("fruit")) // [banana apple orange]

	m.RemoveAllEntries("fruit", "apple")
	fmt.Println("删除所有 apple 后:", m.Get("fruit")) // [banana orange]

	// 10. 清空
	m.Clear()
	fmt.Println("清空后大小:", m.Size()) // 0
}
