//go:build embed_frontend

package webui

import (
	"embed"
	"io/fs"
)

//go:embed dist
var embeddedFiles embed.FS

func embeddedFrontend() (fs.FS, bool, error) {
	frontend, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		return nil, false, err
	}
	return frontend, true, nil
}
