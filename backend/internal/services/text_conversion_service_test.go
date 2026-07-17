package services

import (
	"errors"
	"strings"
	"testing"

	"shihai/internal/dto"
)

type fakeTextConverter struct {
	replace map[string]string
}

func (c fakeTextConverter) Convert(text string) (string, error) {
	converted := text
	for from, to := range c.replace {
		converted = strings.ReplaceAll(converted, from, to)
	}
	return converted, nil
}

func TestTextConversionService_ConvertTexts_ConvertsMultipleTextsInOrder(t *testing.T) {
	service := NewTextConversionService(func(config string) (TextConverter, error) {
		if config != "s2t" {
			t.Fatalf("expected config s2t, got %s", config)
		}
		return fakeTextConverter{replace: map[string]string{"汉": "漢", "国": "國"}}, nil
	})

	resp, err := service.ConvertTexts(&dto.TextConversionRequest{
		Mode:  "s2t",
		Texts: []string{"汉家烟尘", "故国三千里"},
	})

	if err != nil {
		t.Fatalf("ConvertTexts returned error: %v", err)
	}
	expected := []string{"漢家烟尘", "故國三千里"}
	for index, text := range expected {
		if resp.Texts[index] != text {
			t.Fatalf("expected text %d to be %q, got %q", index, text, resp.Texts[index])
		}
	}
	if resp.Mode != "s2t" {
		t.Fatalf("expected mode s2t, got %s", resp.Mode)
	}
}

func TestTextConversionService_ConvertTexts_UsesDefaultOpenCCConverter(t *testing.T) {
	service := NewTextConversionService()

	resp, err := service.ConvertTexts(&dto.TextConversionRequest{
		Mode:  "s2t",
		Texts: []string{"汉"},
	})

	if err != nil {
		t.Fatalf("ConvertTexts returned error: %v", err)
	}
	if resp.Texts[0] != "漢" {
		t.Fatalf("expected default OpenCC converter to return 漢, got %q", resp.Texts[0])
	}
}

func TestTextConversionService_ConvertTexts_RejectsUnsupportedMode(t *testing.T) {
	service := NewTextConversionService(func(config string) (TextConverter, error) {
		t.Fatalf("factory should not be called for unsupported mode %s", config)
		return nil, nil
	})

	_, err := service.ConvertTexts(&dto.TextConversionRequest{
		Mode:  "bad",
		Texts: []string{"山河"},
	})

	if err == nil {
		t.Fatal("expected unsupported mode error")
	}
	if !strings.Contains(err.Error(), "不支持") {
		t.Fatalf("expected unsupported mode message, got %v", err)
	}
}

func TestTextConversionService_ConvertTexts_AllowsEmptyTextArray(t *testing.T) {
	service := NewTextConversionService(func(config string) (TextConverter, error) {
		return fakeTextConverter{}, nil
	})

	resp, err := service.ConvertTexts(&dto.TextConversionRequest{
		Mode:  "t2s",
		Texts: []string{},
	})

	if err != nil {
		t.Fatalf("ConvertTexts returned error: %v", err)
	}
	if len(resp.Texts) != 0 {
		t.Fatalf("expected empty texts, got %v", resp.Texts)
	}
}

func TestTextConversionService_ConvertTexts_ReturnsConverterError(t *testing.T) {
	service := NewTextConversionService(func(config string) (TextConverter, error) {
		return nil, errors.New("load failed")
	})

	_, err := service.ConvertTexts(&dto.TextConversionRequest{
		Mode:  "s2t",
		Texts: []string{"山河"},
	})

	if err == nil {
		t.Fatal("expected converter error")
	}
	if !strings.Contains(err.Error(), "load failed") {
		t.Fatalf("expected wrapped converter error, got %v", err)
	}
}
