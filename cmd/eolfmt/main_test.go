package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
		tmpDir := t.TempDir()

		testFiles := []struct {
		name    string
		content string
	}{
		{"test1.txt", "no newline"},
		{"test2.txt", "has newline\n"},
		{"test3.md", "markdown file"},
	}

	for _, tf := range testFiles {
		path := filepath.Join(tmpDir, tf.name)
		if err := os.WriteFile(path, []byte(tf.content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

		ctx := context.Background()
	if err := run(ctx, []string{tmpDir}); err != nil {
		t.Errorf("run failed: %v", err)
	}

		for _, tf := range testFiles {
		path := filepath.Join(tmpDir, tf.name)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

				if len(content) > 0 && content[len(content)-1] != '\n' {
			t.Errorf("file %s does not end with newline", tf.name)
		}
	}
}

func TestRunWithCancel(t *testing.T) {
		tmpDir := t.TempDir()
	
		for i := 0; i < 100; i++ {
		path := filepath.Join(tmpDir, fmt.Sprintf("test%d.txt", i))
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

		err := run(ctx, []string{tmpDir})
	if err == nil {
				t.Skip("processing completed before context cancellation")
	}
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("expected context error, got %v", err)
	}
}
