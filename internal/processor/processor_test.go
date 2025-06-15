package processor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckAndFixFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		content      string
		expectedFix  bool
		expectedErr  bool
		skipCreate   bool
		isBinary     bool
	}{
		{
			name:        "file without newline",
			content:     "Hello, World!",
			expectedFix: true,
		},
		{
			name:        "file with newline",
			content:     "Hello, World!\n",
			expectedFix: false,
		},
		{
			name:        "empty file",
			content:     "",
			expectedFix: false,
		},
		{
			name:        "binary file",
			content:     "binary\x00data",
			expectedFix: false,
			isBinary:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "test.txt"
			if tt.isBinary {
				filename = "test"
			}
			testFile := filepath.Join(tmpDir, filename)

			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			info, err := os.Stat(testFile)
			if err != nil {
				t.Fatalf("failed to stat file: %v", err)
			}

			fixed, err := CheckAndFixFile(testFile, info)

			if tt.expectedErr && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if fixed != tt.expectedFix {
				t.Errorf("expected fix: %v, actual: %v", tt.expectedFix, fixed)
			}

			if tt.expectedFix && !tt.isBinary {
				content, err := os.ReadFile(testFile)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if len(content) == 0 || content[len(content)-1] != '\n' {
					t.Errorf("newline was not added")
				}
			}

			os.Remove(testFile)
		})
	}
}

func TestCheckBinary(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		content  string
		isBinary bool
	}{
		{
			name:     "text file",
			content:  "This is a text file\nwith multiple lines",
			isBinary: false,
		},
		{
			name:     "binary file (contains NULL)",
			content:  "Binary\x00Data",
			isBinary: true,
		},
		{
			name:     "UTF-8 text",
			content:  "こんにちは、世界！",
			isBinary: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test.bin")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			err = checkBinary(testFile)
			if tt.isBinary && err == nil {
				t.Errorf("should be detected as binary file")
			} else if !tt.isBinary && err != nil {
				t.Errorf("should be treated as text file: %v", err)
			}

			os.Remove(testFile)
		})
	}
}
