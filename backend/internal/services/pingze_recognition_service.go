package services

import (
	"fmt"

	"shihai/internal/dto"
	"shihai/internal/poetry"
)

// PingzeRecognitionService 提供诗词平仄自动识别服务。
type PingzeRecognitionService struct{}

// NewPingzeRecognitionService 创建诗词平仄自动识别服务。
func NewPingzeRecognitionService() *PingzeRecognitionService {
	return &PingzeRecognitionService{}
}

// RecognizePingze 按请求文本批量识别平仄，并保持输入顺序。
func (s *PingzeRecognitionService) RecognizePingze(req *dto.PingzeRecognitionRequest) (*dto.PingzeRecognitionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("平仄识别请求不能为空")
	}
	return &dto.PingzeRecognitionResponse{
		Pingze: poetry.RecognizePingzeLines(req.Texts),
	}, nil
}
