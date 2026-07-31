package utils

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		wantErr bool
	}{
		{
			name:    "有效的 JSON 对象",
			jsonStr: `{"name":"张三","age":30}`,
			wantErr: false,
		},
		{
			name:    "空的 JSON 对象",
			jsonStr: `{}`,
			wantErr: false,
		},
		{
			name:    "空字符串",
			jsonStr: "",
			wantErr: true,
		},
		{
			name:    "无效的 JSON",
			jsonStr: `{invalid}`,
			wantErr: true,
		},
		{
			name:    "JSON 数组（应该失败）",
			jsonStr: `[{"name":"张三"}]`,
			wantErr: true,
		},
		{
			name:    "null 不是对象",
			jsonStr: `null`,
			wantErr: true,
		},
		{
			name:    "包含多个 JSON 值",
			jsonStr: `{"name":"张三"} {"name":"李四"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Errorf("Parse() 返回 nil，期望非 nil")
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	t.Run("有效 JSON", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustParse() 不应该 panic")
			}
		}()
		result := MustParse(`{"name":"张三"}`)
		if result == nil {
			t.Errorf("MustParse() 返回 nil")
		}
	})

	t.Run("无效 JSON", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustParse() 应该 panic")
			}
		}()
		MustParse(`{invalid}`)
	})
}

func TestParseBytes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "有效的字节数据",
			data:    []byte(`{"name":"张三"}`),
			wantErr: false,
		},
		{
			name:    "空字节数组",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "无效的 JSON",
			data:    []byte(`{invalid}`),
			wantErr: true,
		},
		{
			name:    "null 不是对象",
			data:    []byte(`null`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBytes(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Errorf("ParseBytes() 返回 nil")
			}
		})
	}
}

func TestParseReader(t *testing.T) {
	t.Run("有效的 Reader", func(t *testing.T) {
		reader := strings.NewReader(`{"name":"张三"}`)
		result, err := ParseReader(reader)
		if err != nil {
			t.Errorf("ParseReader() error = %v", err)
		}
		if result == nil {
			t.Errorf("ParseReader() 返回 nil")
		}
		if result.GetString("name") != "张三" {
			t.Errorf("ParseReader() 解析错误，期望 张三，得到 %v", result.GetString("name"))
		}
	})

	t.Run("nil Reader", func(t *testing.T) {
		_, err := ParseReader(nil)
		if err == nil {
			t.Errorf("ParseReader() 应该返回错误")
		}
	})

	t.Run("无效的 Reader", func(t *testing.T) {
		reader := strings.NewReader(`{invalid}`)
		_, err := ParseReader(reader)
		if err == nil {
			t.Errorf("ParseReader() 应该返回错误")
		}
	})

	t.Run("Reader 包含多个 JSON 值", func(t *testing.T) {
		reader := strings.NewReader(`{"name":"张三"} {"name":"李四"}`)
		if _, err := ParseReader(reader); err == nil {
			t.Error("ParseReader() 应拒绝多个 JSON 值")
		}
	})
}

func TestParseWithNumber(t *testing.T) {
	t.Run("保留数字精度", func(t *testing.T) {
		jsonStr := `{"age":30,"score":98.5,"big":9007199254740993}`
		result, err := ParseWithNumber(jsonStr)
		if err != nil {
			t.Errorf("ParseWithNumber() error = %v", err)
		}

		// 使用 Get 获取值，然后验证类型
		age := result.Get("age")
		if age == nil {
			t.Errorf("age 为 nil")
		}
		// 在 ParseWithNumber 中，数字应该是 json.Number 类型
		if _, ok := age.(json.Number); !ok {
			t.Errorf("age 应该是 json.Number 类型，实际类型: %T", age)
		}
		// 使用 GetInt 转换
		ageInt := result.GetInt("age")
		if ageInt != 30 {
			t.Errorf("age = %v, want 30", ageInt)
		}

		score := result.GetFloat64("score")
		if score != 98.5 {
			t.Errorf("score = %v, want 98.5", score)
		}

		big := result.Get("big")
		if _, ok := big.(json.Number); !ok {
			t.Errorf("big 应该是 json.Number 类型，实际类型: %T", big)
		}
		if got := result.GetInt64("big"); got != 9007199254740993 {
			t.Errorf("big = %d, want 9007199254740993", got)
		}
	})
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		wantLen int
		wantErr bool
	}{
		{
			name:    "有效的 JSON 数组",
			jsonStr: `[{"name":"张三"},{"name":"李四"}]`,
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "空数组",
			jsonStr: `[]`,
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "空字符串",
			jsonStr: "",
			wantErr: true,
		},
		{
			name:    "JSON 对象（应该失败）",
			jsonStr: `{"name":"张三"}`,
			wantErr: true,
		},
		{
			name:    "null 不是数组",
			jsonStr: `null`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseArray(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("ParseArray() len = %v, want %v", len(result), tt.wantLen)
			}
		})
	}
}

func TestJSONMap_Get(t *testing.T) {
	jsonStr := `{
		"name": "张三",
		"age": 30,
		"score": 98.5,
		"active": true,
		"address": {
			"city": "北京",
			"district": "海淀",
			"detail": {
				"street": "中关村大街",
				"number": 1
			}
		},
		"hobbies": ["编程", "阅读", "游泳"],
		"matrix": [[1, 2], [3, 4]],
		"users": [
			{
				"name": "李四",
				"age": 25,
				"address": {
					"city": "上海"
				},
				"hobbies": ["篮球", "音乐"]
			},
			{
				"name": "王五",
				"age": 28,
				"address": {
					"city": "深圳"
				}
			}
		],
		"nullValue": null
	}`

	data := MustParse(jsonStr)

	tests := []struct {
		name     string
		path     string
		expected interface{}
	}{
		{"简单字段", "name", "张三"},
		{"数字字段", "age", 30.0},
		{"浮点数字段", "score", 98.5},
		{"布尔字段", "active", true},
		{"嵌套对象", "address.city", "北京"},
		{"多层嵌套", "address.detail.street", "中关村大街"},
		{"数组索引", "hobbies[0]", "编程"},
		{"数组索引2", "hobbies[1]", "阅读"},
		{"数组中的对象", "users[0].name", "李四"},
		{"数组中的对象嵌套", "users[0].address.city", "上海"},
		{"深层嵌套数组", "users[0].hobbies[0]", "篮球"},
		{"连续数组索引", "matrix[1][0]", 3.0},
		{"不存在的路径", "nonexistent", nil},
		{"不存在的嵌套路径", "address.province", nil},
		{"不存在的数组索引", "hobbies[10]", nil},
		{"非法数组索引", "hobbies[0]extra", nil},
		{"空路径段", "address..city", nil},
		{"null 值", "nullValue", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.Get(tt.path)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Get(%s) = %v, want nil", tt.path, result)
				}
			} else if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Get(%s) = %v (type %T), want %v (type %T)",
					tt.path, result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetWithDefault(t *testing.T) {
	data := MustParse(`{
		"name": "张三",
		"zero": 0,
		"inactive": false,
		"nothing": null,
		"users": [{"name": "李四"}]
	}`)

	tests := []struct {
		name         string
		path         string
		defaultValue interface{}
		want         interface{}
	}{
		{name: "已有顶层属性", path: "name", defaultValue: "默认", want: "张三"},
		{name: "已有嵌套属性", path: "users[0].name", defaultValue: "默认", want: "李四"},
		{name: "零值不使用默认值", path: "zero", defaultValue: 99, want: float64(0)},
		{name: "false 不使用默认值", path: "inactive", defaultValue: true, want: false},
		{name: "null 不使用默认值", path: "nothing", defaultValue: "默认", want: nil},
		{name: "属性不存在", path: "missing", defaultValue: "默认", want: "默认"},
		{name: "嵌套属性不存在", path: "users[0].age", defaultValue: float64(18), want: float64(18)},
		{name: "数组越界", path: "users[1].name", defaultValue: "默认", want: "默认"},
		{name: "表达式无效", path: "users[0]name", defaultValue: "默认", want: "默认"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := data.Get(tt.path, tt.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get(%q, %#v) = %#v, want %#v", tt.path, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestJSONMap_GetString(t *testing.T) {
	data := MustParse(`{"name":"张三","address":{"city":"北京"}}`)

	tests := []struct {
		path     string
		expected string
	}{
		{"name", "张三"},
		{"address.city", "北京"},
		{"nonexistent", ""},
		{"address.province", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetString(tt.path)
			if result != tt.expected {
				t.Errorf("GetString(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetStringOrDefault(t *testing.T) {
	data := MustParse(`{"name":"张三","address":{"city":"北京"},"empty":""}`)

	tests := []struct {
		path       string
		defaultVal string
		expected   string
	}{
		{"name", "默认", "张三"},
		{"address.city", "默认", "北京"},
		{"empty", "默认", ""}, // 空字符串存在，应该返回空字符串
		{"nonexistent", "默认", "默认"},
		{"address.province", "未知", "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetStringOrDefault(tt.path, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("GetStringOrDefault(%s, %s) = %v, want %v",
					tt.path, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetInt(t *testing.T) {
	data := MustParse(`{"age":30,"score":98.5,"users":[{"age":25}]}`)

	tests := []struct {
		path     string
		expected int
	}{
		{"age", 30},
		{"score", 98},
		{"users[0].age", 25},
		{"nonexistent", 0},
		{"address.age", 0},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetInt(tt.path)
			if result != tt.expected {
				t.Errorf("GetInt(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetIntOrDefault(t *testing.T) {
	data := MustParse(`{"age":30,"zero":0}`)

	tests := []struct {
		path       string
		defaultVal int
		expected   int
	}{
		{"age", 99, 30},
		{"zero", 99, 0}, // 存在但值为0，应该返回0
		{"nonexistent", 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetIntOrDefault(tt.path, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("GetIntOrDefault(%s, %d) = %v, want %v",
					tt.path, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetInt64(t *testing.T) {
	data := MustParse(`{"age":30,"big":9223372036854775807}`)

	tests := []struct {
		path     string
		expected int64
	}{
		{"age", 30},
		{"big", 9223372036854775807},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetInt64(tt.path)
			if result != tt.expected {
				t.Errorf("GetInt64(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetFloat64(t *testing.T) {
	data := MustParse(`{"age":30,"score":98.5,"users":[{"score":95.5}]}`)

	tests := []struct {
		path     string
		expected float64
	}{
		{"age", 30.0},
		{"score", 98.5},
		{"users[0].score", 95.5},
		{"nonexistent", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetFloat64(tt.path)
			if result != tt.expected {
				t.Errorf("GetFloat64(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetBool(t *testing.T) {
	data := MustParse(`{"active":true,"inactive":false,"users":[{"active":true}]}`)

	tests := []struct {
		path     string
		expected bool
	}{
		{"active", true},
		{"inactive", false},
		{"users[0].active", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetBool(tt.path)
			if result != tt.expected {
				t.Errorf("GetBool(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_GetSlice(t *testing.T) {
	data := MustParse(`{"hobbies":["编程","阅读","游泳"],"users":[{"name":"张三"}]}`)

	tests := []struct {
		path    string
		wantLen int
		wantNil bool
	}{
		{"hobbies", 3, false},
		{"users", 1, false},
		{"nonexistent", 0, true},
		{"hobbies[0]", 0, true}, // 不是数组
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetSlice(tt.path)
			if tt.wantNil {
				if result != nil {
					t.Errorf("GetSlice(%s) = %v, want nil", tt.path, result)
				}
			} else {
				if result == nil {
					t.Errorf("GetSlice(%s) 返回 nil", tt.path)
				}
				if len(result) != tt.wantLen {
					t.Errorf("GetSlice(%s) len = %v, want %v", tt.path, len(result), tt.wantLen)
				}
			}
		})
	}
}

func TestJSONMap_GetMap(t *testing.T) {
	data := MustParse(`{"address":{"city":"北京","district":"海淀"},"users":[{"address":{"city":"上海"}}]}`)

	tests := []struct {
		path     string
		wantNil  bool
		checkKey string
		checkVal string
	}{
		{"address", false, "city", "北京"},
		{"users[0].address", false, "city", "上海"},
		{"nonexistent", true, "", ""},
		{"hobbies", true, "", ""}, // 不是 map
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.GetMap(tt.path)
			if tt.wantNil {
				if result != nil {
					t.Errorf("GetMap(%s) = %v, want nil", tt.path, result)
				}
			} else {
				if result == nil {
					t.Errorf("GetMap(%s) 返回 nil", tt.path)
				}
				if result.GetString(tt.checkKey) != tt.checkVal {
					t.Errorf("GetMap(%s).%s = %v, want %v",
						tt.path, tt.checkKey, result.GetString(tt.checkKey), tt.checkVal)
				}
			}
		})
	}
}

func TestJSONMap_Exists(t *testing.T) {
	data := MustParse(`{"name":"张三","address":{"city":"北京"},"hobbies":["编程"],"nothing":null}`)

	tests := []struct {
		path     string
		expected bool
	}{
		{"name", true},
		{"address.city", true},
		{"hobbies[0]", true},
		{"nothing", true},
		{"nonexistent", false},
		{"address.province", false},
		{"hobbies[10]", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := data.Exists(tt.path)
			if result != tt.expected {
				t.Errorf("Exists(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestJSONMap_Keys(t *testing.T) {
	data := MustParse(`{"name":"张三","age":30,"address":{"city":"北京"}}`)
	keys := data.Keys()

	expectedKeys := []string{"address", "age", "name"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Keys() = %v, want %v", keys, expectedKeys)
	}
}

func TestJSONMap_ToJSON(t *testing.T) {
	data := MustParse(`{"name":"张三","age":30}`)
	jsonStr, err := data.ToJSON()
	if err != nil {
		t.Errorf("ToJSON() error = %v", err)
	}

	// 解析回来验证
	result, err := Parse(jsonStr)
	if err != nil {
		t.Errorf("解析 ToJSON 结果失败: %v", err)
	}
	if result.GetString("name") != "张三" {
		t.Errorf("ToJSON() 结果错误")
	}
}

func TestJSONMap_ToJSONPretty(t *testing.T) {
	data := MustParse(`{"name":"张三","age":30}`)
	jsonStr, err := data.ToJSONPretty()
	if err != nil {
		t.Errorf("ToJSONPretty() error = %v", err)
	}

	// 检查是否包含换行和缩进
	if !strings.Contains(jsonStr, "\n") {
		t.Errorf("ToJSONPretty() 结果应该包含换行")
	}
	if !strings.Contains(jsonStr, "  ") {
		t.Errorf("ToJSONPretty() 结果应该包含缩进")
	}
}

func TestJSONMap_ToJSONErrors(t *testing.T) {
	data := JSONMap{"unsupported": make(chan int)}
	if _, err := data.ToJSON(); err == nil {
		t.Error("ToJSON() 应在包含不支持的值时返回错误")
	}
	if _, err := data.ToJSONPretty(); err == nil {
		t.Error("ToJSONPretty() 应在包含不支持的值时返回错误")
	}
}

func TestJSONMap_Merge(t *testing.T) {
	m1 := JSONMap{"name": "张三", "age": 30}
	m2 := JSONMap{"age": 31, "city": "北京"}

	result := m1.Merge(m2)

	if result.GetString("name") != "张三" {
		t.Errorf("Merge() name = %v, want 张三", result.GetString("name"))
	}
	if result.GetInt("age") != 31 {
		t.Errorf("Merge() age = %v, want 31", result.GetInt("age"))
	}
	if result.GetString("city") != "北京" {
		t.Errorf("Merge() city = %v, want 北京", result.GetString("city"))
	}

	// 验证原对象未被修改
	if m1.GetInt("age") != 30 {
		t.Errorf("原始 m1 被修改")
	}
}

func TestJSONMap_DeepMerge(t *testing.T) {
	m1 := JSONMap{
		"name": "张三",
		"address": JSONMap{
			"city":     "北京",
			"district": "海淀",
		},
		"hobbies": []interface{}{"编程", "阅读"},
	}
	m2 := JSONMap{
		"age": 30,
		"address": JSONMap{
			"city": "上海",
			"detail": JSONMap{
				"street": "南京路",
			},
		},
		"hobbies": []interface{}{"游泳", "音乐"},
	}

	result := m1.DeepMerge(m2)

	if result.GetString("name") != "张三" {
		t.Errorf("DeepMerge() name = %v, want 张三", result.GetString("name"))
	}
	if result.GetInt("age") != 30 {
		t.Errorf("DeepMerge() age = %v, want 30", result.GetInt("age"))
	}
	// 嵌套合并
	if result.GetString("address.city") != "上海" {
		t.Errorf("DeepMerge() address.city = %v, want 上海", result.GetString("address.city"))
	}
	if result.GetString("address.district") != "海淀" {
		t.Errorf("DeepMerge() address.district = %v, want 海淀", result.GetString("address.district"))
	}
	if result.GetString("address.detail.street") != "南京路" {
		t.Errorf("DeepMerge() address.detail.street = %v, want 南京路", result.GetString("address.detail.street"))
	}
	// 非 map 类型应该覆盖
	hobbies := result.GetSlice("hobbies")
	if len(hobbies) != 2 || hobbies[0] != "游泳" {
		t.Errorf("DeepMerge() hobbies = %v, want [游泳 音乐]", hobbies)
	}

	result.GetMap("address")["city"] = "杭州"
	if got := m1.GetString("address.city"); got != "北京" {
		t.Errorf("DeepMerge() 修改了原始对象: address.city = %q", got)
	}
	result.GetSlice("hobbies")[0] = "跑步"
	if got := m2.GetSlice("hobbies")[0]; got != "游泳" {
		t.Errorf("DeepMerge() 修改了覆盖来源数组: hobbies[0] = %v", got)
	}
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		jsonStr  string
		expected bool
	}{
		{`{"name":"张三"}`, true},
		{`[{"name":"张三"}]`, true},
		{`"hello"`, true},
		{`123`, true},
		{`true`, true},
		{`null`, true},
		{`{invalid}`, false},
		{``, false},
	}

	for _, tt := range tests {
		t.Run(tt.jsonStr, func(t *testing.T) {
			result := IsJSON(tt.jsonStr)
			if result != tt.expected {
				t.Errorf("IsJSON(%s) = %v, want %v", tt.jsonStr, result, tt.expected)
			}
		})
	}
}

func TestIsJSONObject(t *testing.T) {
	tests := []struct {
		jsonStr  string
		expected bool
	}{
		{`{"name":"张三"}`, true},
		{`{}`, true},
		{`[{"name":"张三"}]`, false},
		{`"hello"`, false},
		{`123`, false},
		{`null`, false},
	}

	for _, tt := range tests {
		t.Run(tt.jsonStr, func(t *testing.T) {
			result := IsJSONObject(tt.jsonStr)
			if result != tt.expected {
				t.Errorf("IsJSONObject(%s) = %v, want %v", tt.jsonStr, result, tt.expected)
			}
		})
	}
}

func TestIsJSONArray(t *testing.T) {
	tests := []struct {
		jsonStr  string
		expected bool
	}{
		{`[{"name":"张三"}]`, true},
		{`[]`, true},
		{`{"name":"张三"}`, false},
		{`"hello"`, false},
		{`123`, false},
		{`null`, false},
	}

	for _, tt := range tests {
		t.Run(tt.jsonStr, func(t *testing.T) {
			result := IsJSONArray(tt.jsonStr)
			if result != tt.expected {
				t.Errorf("IsJSONArray(%s) = %v, want %v", tt.jsonStr, result, tt.expected)
			}
		})
	}
}

// 基准测试
func BenchmarkParse(b *testing.B) {
	jsonStr := `{"name":"张三","age":30,"address":{"city":"北京"}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(jsonStr)
	}
}

func BenchmarkGetString(b *testing.B) {
	data := MustParse(`{"name":"张三","age":30,"address":{"city":"北京"}}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.GetString("address.city")
	}
}

func BenchmarkGetInt(b *testing.B) {
	data := MustParse(`{"name":"张三","age":30,"address":{"city":"北京"}}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.GetInt("age")
	}
}

func BenchmarkGetNestedPath(b *testing.B) {
	data := MustParse(`{"users":[{"address":{"city":"北京"}}]}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.GetString("users[0].address.city")
	}
}
