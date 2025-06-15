// Package walker provides directory traversal functionality.
package walker

import (
	"os"
	"path/filepath"
)

// Excluded directories (O(1) lookup)
var excludeDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	".idea":        true,
	"target":       true,
	"build":        true,
	"dist":         true,
	"vendor":       true,
	"__pycache__":  true,
}

// FileInfo contains file path and metadata information
type FileInfo struct {
	Path string
	Info os.FileInfo
}

// WalkFiles walks through directories and returns a channel of FileInfo
func WalkFiles(root string) <-chan FileInfo {
	out := make(chan FileInfo, 1000)

	go func() {
		defer close(out)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Ignore errors and continue
			}

			if info.IsDir() {
				if excludeDirs[info.Name()] {
					return filepath.SkipDir
				}
				return nil
			}

			if info.Mode()&os.ModeSymlink != 0 {
				return nil
			}

			out <- FileInfo{
				Path: path,
				Info: info,
			}
			return nil
		})
	}()

	return out
}
