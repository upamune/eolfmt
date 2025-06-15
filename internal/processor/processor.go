// Package processor handles file processing and newline fixing.
package processor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Known binary extensions (skip processing)
var skipExtensions = map[string]bool{
	".exe":   true,
	".dll":   true,
	".so":    true,
	".dylib": true,
	".png":   true,
	".jpg":   true,
	".jpeg":  true,
	".gif":   true,
	".pdf":   true,
	".zip":   true,
	".tar":   true,
	".gz":    true,
	".pyc":   true,
	".class": true,
	".jar":   true,
}

// Known text extensions (no content check needed)
var textExtensions = map[string]bool{
	".txt":  true,
	".md":   true,
	".go":   true,
	".js":   true,
	".py":   true,
	".java": true,
	".c":    true,
	".h":    true,
	".sh":   true,
	".yml":  true,
	".json": true,
	".xml":  true,
}

// CheckAndFixFile checks if a file ends with a newline and adds one if missing
func CheckAndFixFile(path string, info os.FileInfo) (bool, error) {
	size := info.Size()
	if size == 0 {
		return false, nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	if skipExtensions[ext] {
		return false, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	buf := make([]byte, 1)
	_, err = file.ReadAt(buf, size-1)
	file.Close()

	if err != nil {
		return false, err
	}

	if buf[0] == '\n' {
		return false, nil
	}

	if !textExtensions[ext] {
		if err := checkBinary(path); err != nil {
			return false, nil
		}
	}

	return addNewline(path), nil
}

func checkBinary(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)

	if bytes.IndexByte(buf[:n], 0) != -1 {
		return fmt.Errorf("binary file")
	}

	return nil
}

func addNewline(path string) bool {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	defer file.Close()

	_, err = file.Write([]byte("\n"))
	return err == nil
}
