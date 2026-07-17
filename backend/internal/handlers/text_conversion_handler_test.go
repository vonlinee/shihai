package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shihai/internal/dto"

	"github.com/gin-gonic/gin"
)

type fakeTextConversionService struct {
	request *dto.TextConversionRequest
}

func (s *fakeTextConversionService) ConvertTexts(req *dto.TextConversionRequest) (*dto.TextConversionResponse, error) {
	s.request = req
	return &dto.TextConversionResponse{
		Mode:  req.Mode,
		Texts: []string{"漢家烟尘", "故國三千里"},
	}, nil
}

func TestTextConversionHandler_ConvertTexts_ReturnsConvertedTextArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeTextConversionService{}
	handler := NewTextConversionHandler(service)
	router := gin.New()
	router.POST("/api/admin/text-conversion", handler.ConvertTexts)

	body := []byte(`{"mode":"s2t","texts":["汉家烟尘","故国三千里"]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/text-conversion", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.request == nil || len(service.request.Texts) != 2 {
		t.Fatalf("expected service to receive two texts, got %#v", service.request)
	}

	var response struct {
		Code int                        `json:"code"`
		Data dto.TextConversionResponse `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Code != 200 {
		t.Fatalf("expected response code 200, got %d", response.Code)
	}
	expected := []string{"漢家烟尘", "故國三千里"}
	for index, text := range expected {
		if response.Data.Texts[index] != text {
			t.Fatalf("expected text %d to be %q, got %q", index, text, response.Data.Texts[index])
		}
	}
}

func TestTextConversionHandler_ConvertTexts_RejectsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTextConversionHandler(&fakeTextConversionService{})
	router := gin.New()
	router.POST("/api/admin/text-conversion", handler.ConvertTexts)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/text-conversion", bytes.NewReader([]byte(`{"mode":"s2t"}`)))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}
