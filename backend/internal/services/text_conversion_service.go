package services

import (
	"fmt"
	"sync"

	"shihai/internal/dto"

	"github.com/longbridgeapp/opencc"
)

// TextConverter 定义单段文本转换能力。
type TextConverter interface {
	Convert(text string) (string, error)
}

// TextConverterFactory 定义按 OpenCC 配置创建转换器的工厂函数。
type TextConverterFactory func(config string) (TextConverter, error)

// TextConversionService 提供中文简繁转换服务。
type TextConversionService struct {
	factory TextConverterFactory
	mu      sync.Mutex
	cache   map[string]TextConverter
}

var supportedTextConversionModes = map[string]struct{}{
	"s2t":   {},
	"t2s":   {},
	"s2tw":  {},
	"tw2s":  {},
	"s2hk":  {},
	"hk2s":  {},
	"t2tw":  {},
	"tw2t":  {},
	"t2hk":  {},
	"hk2t":  {},
	"t2jp":  {},
	"jp2t":  {},
	"tw2sp": {},
}

// NewTextConversionService 创建中文简繁转换服务。
//
// factory 可选传入，主要用于测试；生产环境默认使用 longbridge/opencc。
func NewTextConversionService(factory ...TextConverterFactory) *TextConversionService {
	textConverterFactory := TextConverterFactory(func(config string) (TextConverter, error) {
		return opencc.New(config)
	})
	if len(factory) > 0 && factory[0] != nil {
		textConverterFactory = factory[0]
	}

	return &TextConversionService{
		factory: textConverterFactory,
		cache:   make(map[string]TextConverter),
	}
}

// ConvertTexts 按请求模式批量转换文本，并保持输入顺序。
//
// req 为转换请求，texts 允许为空数组；当 mode 不在支持范围内或转换器失败时返回错误。
func (s *TextConversionService) ConvertTexts(req *dto.TextConversionRequest) (*dto.TextConversionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("转换请求不能为空")
	}
	if !isSupportedTextConversionMode(req.Mode) {
		return nil, fmt.Errorf("不支持的转换模式: %s", req.Mode)
	}

	converter, err := s.getConverter(req.Mode)
	if err != nil {
		return nil, fmt.Errorf("初始化转换器失败: %w", err)
	}

	texts := make([]string, 0, len(req.Texts))
	for _, text := range req.Texts {
		converted, err := converter.Convert(text)
		if err != nil {
			return nil, fmt.Errorf("转换文本失败: %w", err)
		}
		texts = append(texts, converted)
	}

	return &dto.TextConversionResponse{
		Mode:  req.Mode,
		Texts: texts,
	}, nil
}

func (s *TextConversionService) getConverter(mode string) (TextConverter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if converter, ok := s.cache[mode]; ok {
		return converter, nil
	}
	converter, err := s.factory(mode)
	if err != nil {
		return nil, err
	}
	s.cache[mode] = converter
	return converter, nil
}

func isSupportedTextConversionMode(mode string) bool {
	_, ok := supportedTextConversionModes[mode]
	return ok
}
