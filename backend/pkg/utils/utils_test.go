package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func testJSONData(t *testing.T) map[string]interface{} {
	t.Helper()

	var data map[string]interface{}
	err := json.Unmarshal([]byte(`{
		"name": "张三",
		"age": 30,
		"score": 98.5,
		"active": true,
		"empty": "",
		"zero": 0,
		"nothing": null,
		"address": {"city": "北京"},
		"users": [{"name": "李四", "hobbies": ["篮球", "音乐"]}],
		"matrix": [[1, 2], [3, 4]]
	}`), &data)
	if err != nil {
		t.Fatalf("解析测试数据失败: %v", err)
	}
	return data
}

func TestGetValue(t *testing.T) {
	data := testJSONData(t)
	tests := []struct {
		name string
		path string
		want interface{}
	}{
		{name: "顶层字段", path: "name", want: "张三"},
		{name: "嵌套字段", path: "address.city", want: "北京"},
		{name: "数组对象字段", path: "users[0].name", want: "李四"},
		{name: "连续数组索引", path: "matrix[1][0]", want: float64(3)},
		{name: "数组中的数组", path: "users[0].hobbies[1]", want: "音乐"},
		{name: "null 值", path: "nothing", want: nil},
		{name: "缺失字段", path: "address.district", want: nil},
		{name: "越界索引", path: "users[2].name", want: nil},
		{name: "负数索引", path: "users[-1]", want: nil},
		{name: "非数字索引", path: "users[x]", want: nil},
		{name: "未闭合索引", path: "users[0", want: nil},
		{name: "索引后有多余字符", path: "users[0]name", want: nil},
		{name: "空路径段", path: "address..city", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValue(data, tt.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue(%q) = %#v, want %#v", tt.path, got, tt.want)
			}
		})
	}

	if got := GetValue(data, ""); !reflect.DeepEqual(got, data) {
		t.Errorf("GetValue(data, empty path) = %#v, want original map", got)
	}
}

func TestTypedGetters(t *testing.T) {
	data := testJSONData(t)
	data["integer"] = 42
	data["integer64"] = int64(64)
	data["number"] = json.Number("12.25")
	data["invalidNumber"] = json.Number("invalid")
	data["namedMap"] = JSONMap{"value": "ok"}

	if got := GetString(data, "name"); got != "张三" {
		t.Errorf("GetString() = %q, want %q", got, "张三")
	}
	if got := GetInt(data, "score"); got != 98 {
		t.Errorf("GetInt() = %d, want 98", got)
	}
	if got := GetInt(data, "integer64"); got != 64 {
		t.Errorf("GetInt() from int64 = %d, want 64", got)
	}
	if got := GetFloat(data, "number"); got != 12.25 {
		t.Errorf("GetFloat() = %v, want 12.25", got)
	}
	if got := GetFloat(data, "invalidNumber"); got != 0 {
		t.Errorf("GetFloat() from invalid number = %v, want 0", got)
	}
	if got := GetBool(data, "active"); !got {
		t.Error("GetBool() = false, want true")
	}
	if got := GetSlice(data, "users"); len(got) != 1 {
		t.Errorf("GetSlice() length = %d, want 1", len(got))
	}
	if got := GetMap(data, "address"); got["city"] != "北京" {
		t.Errorf("GetMap() = %#v, want city 北京", got)
	}
	if got := GetMap(data, "namedMap"); got["value"] != "ok" {
		t.Errorf("GetMap() from JSONMap = %#v, want value ok", got)
	}
	if got := GetString(data, "age"); got != "" {
		t.Errorf("GetString() from number = %q, want empty string", got)
	}
}

func TestDefaultGettersOnlyDefaultForMissingPaths(t *testing.T) {
	data := testJSONData(t)

	if got := GetStringOrDefault(data, "empty", "默认"); got != "" {
		t.Errorf("existing empty string = %q, want empty string", got)
	}
	if got := GetStringOrDefault(data, "nothing", "默认"); got != "" {
		t.Errorf("existing null = %q, want empty string", got)
	}
	if got := GetStringOrDefault(data, "missing", "默认"); got != "默认" {
		t.Errorf("missing string = %q, want 默认", got)
	}
	if got := GetIntOrDefault(data, "zero", 99); got != 0 {
		t.Errorf("existing zero = %d, want 0", got)
	}
	if got := GetIntOrDefault(data, "missing", 99); got != 99 {
		t.Errorf("missing integer = %d, want 99", got)
	}
}

type User struct {
	ID     int
	Name   string
	Age    int
	City   string
	Tags   []string
	Skills []string
}

func TestSliceToMap(t *testing.T) {
	users := []User{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
	}

	// 不指定 valueFn：整个结构体作为 value
	userMap := SliceToMapByKey(users, func(u User) int { return u.ID })
	fmt.Println(userMap) // map[1:{1 Alice 25} 2:{2 Bob 30}]

	// 指定 valueFn：自定义 value
	ageMap := SliceToMapByKeyValue(users,
		func(u User) int { return u.ID },
		func(u User) int { return u.Age },
	)
	fmt.Println(ageMap) // map[1:25 2:30]
}

func TestPluck(t *testing.T) {
	users := []User{
		{ID: 1, Name: "Alice", Age: 25, City: "Beijing"},
		{ID: 2, Name: "Bob", Age: 30, City: "Shanghai"},
		{ID: 3, Name: "Charlie", Age: 35, City: "Beijing"},
		{ID: 4, Name: "David", Age: 30, City: "Guangzhou"},
		{ID: 5, Name: "Eva", Age: 25, City: "Shanghai"},
	}

	// 1. 提取所有 ID
	ids := ExtractToSlice(users, func(u User) int { return u.ID })
	fmt.Println("IDs:", ids)

	// 2. 提取所有城市（去重）
	cities := ExtractToSliceUnique(users, func(u User) string { return u.City })
	fmt.Println("Cities:", cities)

	// 3. 提取所有年龄（去重）
	ages := ExtractToSliceUnique(users, func(u User) int { return u.Age })
	fmt.Println("Ages:", ages)

	// 4. 链式组合：提取 ID 后再进行其他操作
	ids2 := ExtractToSlice(users, func(u User) int { return u.ID })
	// 去重（虽然 ID 本身唯一）
	uniqueIds := func(ids []int) []int {
		seen := make(map[int]struct{})
		var result []int
		for _, id := range ids {
			if _, exists := seen[id]; !exists {
				seen[id] = struct{}{}
				result = append(result, id)
			}
		}
		return result
	}(ids2)
	fmt.Println("Unique IDs:", uniqueIds)
}

func TestFlatToMap(t *testing.T) {
	users := []User{
		{ID: 1, Name: "Alice", City: "Beijing", Tags: []string{"go", "python"}},
		{ID: 2, Name: "Bob", City: "Shanghai", Tags: []string{"java", "go"}},
		{ID: 3, Name: "Charlie", City: "Beijing", Tags: []string{"python", "rust"}},
		{ID: 4, Name: "David", City: "Shanghai", Tags: []string{"go", "rust"}},
	}

	// 1. 展平转 Set
	tagSet := FlatToSlice(users, func(u User) []string { return u.Tags })
	fmt.Println("Tag Set:", tagSet) // {go, python, java, rust}

	// 2. 展平转 Map（按城市分组）
	tagsByCity := FlatToMapByKey(users,
		func(u User) string { return u.City },
		func(u User) []string { return u.Tags },
	)
	fmt.Println("Tags by City:", tagsByCity)
	// map[Beijing:[go python python rust] Shanghai:[java go go rust]]

	// 3. 展平并去重（按城市分组后去重）
	tagsByCityUnique := make(map[string][]string)
	for city, tags := range tagsByCity {
		seen := make(map[string]struct{})
		var unique []string
		for _, tag := range tags {
			if _, exists := seen[tag]; !exists {
				seen[tag] = struct{}{}
				unique = append(unique, tag)
			}
		}
		tagsByCityUnique[city] = unique
	}
	fmt.Println("Unique Tags by City:", tagsByCityUnique)
	// map[Beijing:[go python rust] Shanghai:[java go rust]]
}
