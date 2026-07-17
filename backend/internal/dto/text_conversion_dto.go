package dto

// TextConversionRequest 定义中文文本简繁转换请求。
type TextConversionRequest struct {
	Mode  string   `json:"mode" binding:"required"`  // Mode 转换模式，例如 s2t、t2s。
	Texts []string `json:"texts" binding:"required"` // Texts 待转换文本列表，响应会保持相同顺序。
}

// TextConversionResponse 定义中文文本简繁转换响应。
type TextConversionResponse struct {
	Mode  string   `json:"mode"`  // Mode 实际使用的转换模式。
	Texts []string `json:"texts"` // Texts 转换后的文本列表。
}
