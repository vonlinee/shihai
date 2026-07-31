package utils

import (
	"os"
	"path/filepath"
)

// WalkDirectChildFiles 遍历 rootDir 目录下的直接子文件。
//
// rootDir 待遍历的目录路径。
// handleFile 接收文件完整路径和文件目录项；当回调返回错误时，遍历立即停止并返回该错误。
func WalkDirectChildFiles(rootDir string, handleFile func(path string, entry os.DirEntry) error) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if err := handleFile(filepath.Join(rootDir, entry.Name()), entry); err != nil {
			return err
		}
	}
	return nil
}

// WalkAllFiles 递归遍历 rootDir 目录下的所有文件。
//
// rootDir 待遍历的根目录路径。
// handleFile 接收文件完整路径和文件目录项；当回调返回错误时，遍历立即停止并返回该错误。
func WalkAllFiles(rootDir string, handleFile func(path string, entry os.DirEntry) error) error {
	return filepath.WalkDir(rootDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		return handleFile(path, entry)
	})
}
