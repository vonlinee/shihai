package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterSwaggerRoutesServesDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerSwaggerRoutes(router)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var document struct {
		Swagger string                     `json:"swagger"`
		Paths   map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &document); err != nil {
		t.Fatalf("decode Swagger document: %v", err)
	}
	if document.Swagger != "2.0" {
		t.Fatalf("swagger version = %q, want %q", document.Swagger, "2.0")
	}
	for _, path := range []string{"/api/auth/login", "/api/admin/poems"} {
		if _, ok := document.Paths[path]; !ok {
			t.Fatalf("Swagger document does not contain path %q", path)
		}
	}
}

func TestRegisterSwaggerRoutesServesUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerSwaggerRoutes(router)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), "Swagger UI") {
		t.Fatal("Swagger UI response does not contain its page title")
	}
}
