package utils

import (
	"strings"
	"unicode/utf8"
)

// TextPair 表示一次文本相似度计算的两个输入文本。
type TextPair struct {
	Left  string // Left 左侧文本。
	Right string // Right 右侧文本。
}

// TextLength 返回文本包含的 Unicode 字符数量。
func TextLength(text string) int {
	return utf8.RuneCountInString(text)
}

// ReverseText 按 Unicode 字符反转文本。
func ReverseText(text string) string {
	runes := []rune(text)
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		runes[left], runes[right] = runes[right], runes[left]
	}
	return string(runes)
}

// TruncateText 按 Unicode 字符数量截断文本。
//
// maxLength 小于或等于零时返回空字符串；文本未超过限制时返回原文本。
func TruncateText(text string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}

	runes := []rune(text)
	if len(runes) <= maxLength {
		return text
	}
	return string(runes[:maxLength])
}

// IsBlank 判断文本是否为空或只包含 Unicode 空白字符。
func IsBlank(text string) bool {
	return strings.TrimSpace(text) == ""
}

// TextSimilarity 基于归一化 Levenshtein 距离计算两个文本的相似度。
//
// 返回值范围为 0 到 1，完全相同返回 1。计算以 Unicode 字符为单位，
// 不会自动忽略空格、标点或大小写；两个空文本视为完全相同。
func TextSimilarity(left, right string) float64 {
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	maxLength := max(len(leftRunes), len(rightRunes))
	if maxLength == 0 {
		return 1
	}

	distance := levenshteinDistance(leftRunes, rightRunes)
	return 1 - float64(distance)/float64(maxLength)
}

// BatchTextSimilarity 按输入顺序批量计算文本对的相似度。
func BatchTextSimilarity(pairs []TextPair) []float64 {
	similarities := make([]float64, len(pairs))
	for i, pair := range pairs {
		similarities[i] = TextSimilarity(pair.Left, pair.Right)
	}
	return similarities
}

func levenshteinDistance(left, right []rune) int {
	if len(left) < len(right) {
		left, right = right, left
	}
	if len(right) == 0 {
		return len(left)
	}

	previous := make([]int, len(right)+1)
	current := make([]int, len(right)+1)
	for i := range previous {
		previous[i] = i
	}

	for leftIndex, leftRune := range left {
		current[0] = leftIndex + 1
		for rightIndex, rightRune := range right {
			cost := 0
			if leftRune != rightRune {
				cost = 1
			}
			current[rightIndex+1] = min(
				previous[rightIndex+1]+1,
				current[rightIndex]+1,
				previous[rightIndex]+cost,
			)
		}
		previous, current = current, previous
	}

	return previous[len(right)]
}
