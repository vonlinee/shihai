package poetry

import (
	"errors"
	"strings"
	"unicode"

	pinyin "github.com/mozillazg/go-pinyin"
)

const (
	// PingzePing 表示平声。
	PingzePing = "平"
	// PingzeZe 表示仄声。
	PingzeZe = "仄"
	// PingzeUnknown 表示无法识别平仄。
	PingzeUnknown = "?"
)

var errInvalidPingze = errors.New("pingze can only contain 平, 仄, ?, or non-Han characters")

// RecognizePingzeLines 识别多段文本的平仄标记，并保持输入顺序。
func RecognizePingzeLines(content []string) []string {
	pingze := make([]string, len(content))
	for i, line := range content {
		pingze[i] = RecognizePingzeText(line)
	}
	return pingze
}

// RecognizePingzeText 根据汉字拼音声调识别单段文本的平仄。
//
// 现代拼音一、二声归为平，三、四声归为仄；非汉字字符原样保留，无法识别拼音或声调的汉字记录为 ?。
func RecognizePingzeText(text string) string {
	var builder strings.Builder
	for _, r := range text {
		builder.WriteString(RecognizeRunePingze(r))
	}
	return builder.String()
}

// RecognizeRunePingze 识别单个字符的平仄。
func RecognizeRunePingze(r rune) string {
	if !unicode.Is(unicode.Han, r) {
		return string(r)
	}

	args := pinyin.NewArgs()
	args.Style = pinyin.Tone3
	values := pinyin.SinglePinyin(r, args)
	if len(values) == 0 {
		return PingzeUnknown
	}
	return PingzeFromPinyin(values[0])
}

// PingzeFromPinyin 根据带数字声调的拼音返回平仄标记。
func PingzeFromPinyin(value string) string {
	for _, r := range value {
		switch r {
		case '1', '2':
			return PingzePing
		case '3', '4':
			return PingzeZe
		}
	}
	return PingzeUnknown
}

// NormalizePingzeLines 校验并对齐平仄数组，使其与正文数组长度一致。
func NormalizePingzeLines(content []string, pingze []string) ([]string, error) {
	for _, line := range pingze {
		for _, mark := range strings.TrimSpace(line) {
			if !IsValidPingzeMark(mark) {
				return nil, errInvalidPingze
			}
		}
	}
	return AlignPingzeLines(content, pingze), nil
}

// IsValidPingzeMark 判断字符是否为可存储的平仄标记。
func IsValidPingzeMark(mark rune) bool {
	return string(mark) == PingzePing ||
		string(mark) == PingzeZe ||
		string(mark) == PingzeUnknown ||
		!unicode.Is(unicode.Han, mark)
}

// AlignPingzeLines 将平仄数组按正文文本数量补齐或截断。
func AlignPingzeLines(content []string, pingze []string) []string {
	aligned := make([]string, len(content))
	for i := range content {
		if i < len(pingze) {
			aligned[i] = strings.TrimSpace(pingze[i])
		}
	}
	return aligned
}

// HasPingzeValue 判断平仄数组是否包含人工或自动识别结果。
func HasPingzeValue(pingze []string) bool {
	for _, line := range pingze {
		if strings.TrimSpace(line) != "" {
			return true
		}
	}
	return false
}

// PreparePingzeLinesForSave 返回保存诗词时应使用的平仄数组。
//
// 调用方未提供平仄，或提供的平仄全为空时，自动从正文识别；否则校验并对齐调用方提供的平仄。
func PreparePingzeLinesForSave(content []string, pingze []string) ([]string, error) {
	if !HasPingzeValue(pingze) {
		return RecognizePingzeLines(content), nil
	}
	return NormalizePingzeLines(content, pingze)
}
