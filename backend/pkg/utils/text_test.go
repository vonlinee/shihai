package utils

import (
	"math"
	"testing"
)

func TestTextLengthCountsUnicodeCharacters(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{name: "空文本", text: "", want: 0},
		{name: "中文", text: "床前明月光", want: 5},
		{name: "中英文混合", text: "诗词 Go", want: 5},
		{name: "四字节字符", text: "诗𠮷", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TextLength(tt.text); got != tt.want {
				t.Errorf("TextLength(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}

func TestReverseTextReversesUnicodeCharacters(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{name: "空文本", text: "", want: ""},
		{name: "中文", text: "诗词歌赋", want: "赋歌词诗"},
		{name: "中英文混合", text: "诗Go", want: "oG诗"},
		{name: "四字节字符", text: "诗𠮷", want: "𠮷诗"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseText(tt.text); got != tt.want {
				t.Errorf("ReverseText(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestTruncateTextUsesUnicodeCharacterLimit(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		maxLength int
		want      string
	}{
		{name: "截断中文", text: "床前明月光", maxLength: 3, want: "床前明"},
		{name: "不超过限制", text: "明月", maxLength: 3, want: "明月"},
		{name: "零限制", text: "明月", maxLength: 0, want: ""},
		{name: "负数限制", text: "明月", maxLength: -1, want: ""},
		{name: "四字节字符", text: "诗𠮷文", maxLength: 2, want: "诗𠮷"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateText(tt.text, tt.maxLength); got != tt.want {
				t.Errorf("TruncateText(%q, %d) = %q, want %q", tt.text, tt.maxLength, got, tt.want)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "空文本", text: "", want: true},
		{name: "普通空白", text: " \t\r\n", want: true},
		{name: "Unicode 空白", text: "\u3000", want: true},
		{name: "包含文字", text: " 诗 ", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.text); got != tt.want {
				t.Errorf("IsBlank(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestTextSimilarity(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  float64
	}{
		{name: "相同中文", left: "床前明月光", right: "床前明月光", want: 1},
		{name: "单字替换", left: "床前明月光", right: "床前明月霜", want: 0.8},
		{name: "单字插入", left: "明月", right: "明月光", want: 2.0 / 3.0},
		{name: "完全不同", left: "春风", right: "秋雨", want: 0},
		{name: "两个空文本", left: "", right: "", want: 1},
		{name: "单个空文本", left: "", right: "诗词", want: 0},
		{name: "按字符而非字节", left: "你a", right: "你b", want: 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TextSimilarity(tt.left, tt.right)
			if math.Abs(got-tt.want) > 1e-12 {
				t.Errorf("TextSimilarity(%q, %q) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

func TestTextSimilarityIsSymmetric(t *testing.T) {
	left := "海上生明月"
	right := "海上升明月"

	forward := TextSimilarity(left, right)
	backward := TextSimilarity(right, left)

	if forward != backward {
		t.Errorf("similarity is not symmetric: forward=%v backward=%v", forward, backward)
	}
}

func TestBatchTextSimilarityPreservesInputOrder(t *testing.T) {
	pairs := []TextPair{
		{Left: "相同", Right: "相同"},
		{Left: "春风", Right: "秋雨"},
		{Left: "明月", Right: "明月光"},
	}

	got := BatchTextSimilarity(pairs)
	want := []float64{1, 0, 2.0 / 3.0}

	if len(got) != len(want) {
		t.Fatalf("BatchTextSimilarity() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if math.Abs(got[i]-want[i]) > 1e-12 {
			t.Errorf("BatchTextSimilarity()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestBatchTextSimilaritySupportsEmptyInput(t *testing.T) {
	got := BatchTextSimilarity(nil)

	if len(got) != 0 {
		t.Errorf("BatchTextSimilarity(nil) length = %d, want 0", len(got))
	}
}
