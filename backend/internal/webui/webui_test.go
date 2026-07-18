package webui

import (
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func TestMountRoutesServesStaticAssets(t *testing.T) {
	router := gin.New()
	assets := fstest.MapFS{
		"index.html":    {Data: []byte("<html>app</html>")},
		"assets/app.js": {Data: []byte("console.log('app')")},
	}

	if err := mountRoutes(router, assets); err != nil {
		t.Fatalf("mountRoutes() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "console.log('app')" {
		t.Fatalf("body = %q", recorder.Body.String())
	}
}

func TestMountRoutesFallsBackToIndexForSPARoutes(t *testing.T) {
	router := gin.New()
	assets := fstest.MapFS{
		"index.html": {Data: []byte("<html>app</html>")},
	}

	if err := mountRoutes(router, assets); err != nil {
		t.Fatalf("mountRoutes() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/poems/123", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "<html>app</html>" {
		t.Fatalf("body = %q", recorder.Body.String())
	}
}

func TestMountRoutesDoesNotHandleMissingAPIRoutes(t *testing.T) {
	router := gin.New()
	assets := fstest.MapFS{
		"index.html": {Data: []byte("<html>app</html>")},
	}

	if err := mountRoutes(router, assets); err != nil {
		t.Fatalf("mountRoutes() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/missing", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
	if recorder.Body.String() == "<html>app</html>" {
		t.Fatal("missing API route returned the frontend index")
	}
}

func TestMountRoutesRequiresIndex(t *testing.T) {
	router := gin.New()
	assets := fstest.MapFS{
		"assets/app.js": {Data: []byte("console.log('app')")},
	}

	err := mountRoutes(router, assets)
	if err == nil {
		t.Fatal("mountRoutes() error = nil, want missing index error")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("mountRoutes() error = %v, want fs.ErrNotExist", err)
	}
}
