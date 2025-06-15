package walker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkFiles(t *testing.T) {
		tmpDir := t.TempDir()
	
		dirs := []string{
		"src",
		"src/subdir",
		".git",
		"node_modules",
		"vendor",
	}
	
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
	}
	
		files := []string{
		"README.md",
		"main.go",
		"src/file1.go",
		"src/file2.js",
		"src/subdir/file3.py",
		".git/config",
		"node_modules/package.json",
		"vendor/lib.go",
	}
	
	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		err := os.WriteFile(path, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}
	
		foundFiles := make(map[string]bool)
	for fileInfo := range WalkFiles(tmpDir) {
		relPath, _ := filepath.Rel(tmpDir, fileInfo.Path)
		foundFiles[relPath] = true
	}
	
		expectedFiles := []string{
		"README.md",
		"main.go",
		"src/file1.go",
		"src/file2.js",
		"src/subdir/file3.py",
	}
	
	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("expected file not found: %s", expected)
		}
	}
	
		excludedFiles := []string{
		".git/config",
		"node_modules/package.json",
		"vendor/lib.go",
	}
	
	for _, excluded := range excludedFiles {
		if foundFiles[excluded] {
			t.Errorf("excluded file was included: %s", excluded)
		}
	}
}
