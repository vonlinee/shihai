package dto

// PingzeRecognitionRequest 定义平仄自动识别请求。
type PingzeRecognitionRequest struct {
	Texts []string `json:"texts" binding:"required"` // Texts 待识别文本列表，响应保持相同顺序。
}

// PingzeRecognitionResponse 定义平仄自动识别响应。
type PingzeRecognitionResponse struct {
	Pingze []string `json:"pingze"` // Pingze 与请求文本列表对应的平仄标记。
}
