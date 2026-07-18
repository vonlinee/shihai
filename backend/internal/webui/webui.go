package webui

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes serves the embedded frontend when the binary was built with
// the embed_frontend build tag.
func RegisterRoutes(router *gin.Engine) error {
	frontend, enabled, err := embeddedFrontend()
	if err != nil {
		return fmt.Errorf("load embedded frontend: %w", err)
	}
	if !enabled {
		return nil
	}
	return mountRoutes(router, frontend)
}

func mountRoutes(router *gin.Engine, frontend fs.FS) error {
	index, err := fs.ReadFile(frontend, "index.html")
	if err != nil {
		return fmt.Errorf("read frontend index: %w", err)
	}

	fileServer := http.FileServer(http.FS(frontend))
	router.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") {
			c.Status(http.StatusNotFound)
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}

		assetPath := strings.TrimPrefix(path.Clean("/"+requestPath), "/")
		if info, statErr := fs.Stat(frontend, assetPath); statErr == nil && !info.IsDir() {
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})

	return nil
}
