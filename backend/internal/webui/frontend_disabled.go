//go:build !embed_frontend

package webui

import "io/fs"

func embeddedFrontend() (fs.FS, bool, error) {
	return nil, false, nil
}
