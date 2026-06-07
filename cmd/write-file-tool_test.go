package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWriteFile(t *testing.T) {
	t.Run("HappyPath_NewFile", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "writefile_test_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		filePath := filepath.Join(tempDir, "new_file.txt")
		content := "hello world"
		args := WriteFileArgs{Path: filePath, Content: content}
		argsJSON, _ := json.Marshal(args)

		result := runWriteFile(string(argsJSON))

		if !strings.Contains(result, "Success: Wrote 11 bytes") {
			t.Errorf("Expected success message containing 'Success: Wrote 11 bytes', but got: %s", result)
		}

		gotContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}
		if string(gotContent) != content {
			t.Errorf("Expected content %q, but got %q", content, string(gotContent))
		}
	})

	t.Run("HappyPath_OverwriteExistingFile", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "writefile_test_overwrite_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		filePath := filepath.Join(tempDir, "existing_file.txt")
		originalContent := "original content"
		if err := os.WriteFile(filePath, []byte(originalContent), 0644); err != nil {
			t.Fatalf("Failed to create initial file: %v", err)
		}

		newContent := "new content"
		args := WriteFileArgs{Path: filePath, Content: newContent}
		argsJSON, _ := json.Marshal(args)

		result := runWriteFile(string(argsJSON))

		if !strings.Contains(result, "Success: Wrote 11 bytes") {
			t.Errorf("Expected success message containing 'Success: Wrote 11 bytes', but got: %s", result)
		}

		gotContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}
		if string(gotContent) != newContent {
			t.Errorf("Expected content %q, but got %q", newContent, string(gotContent))
		}
	})

	t.Run("CornerCase_InvalidJSON", func(t *testing.T) {
		invalidJSON := `{"path": "missing_content" ` // Missing closing brace
		result := runWriteFile(invalidJSON)
		if !strings.Contains(result, "Error Parsing arguments") {
			t.Errorf("Expected error message for invalid JSON, but got: %s", result)
		}
	})

	t.Run("CornerCase_PathIsDirectory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "writefile_test_dir_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		args := WriteFileArgs{Path: tempDir, Content: "content"}
		argsJSON, _ := json.Marshal(args)

		result := runWriteFile(string(argsJSON))
		// os.WriteFile returns an error if the path is a directory
		if !strings.Contains(result, "Error writing file") {
			t.Errorf("Expected error message for directory path, but got: %s", result)
		}
	})

	t.Run("CornerCase_NonExistentParentDirectories", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "writefile_test_parents_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create a file where a directory should be
		filePath := filepath.Join(tempDir, "file")
		if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}

		// Try to write to a path that would require creating a subdirectory inside that file
		targetPath := filepath.Join(filePath, "subdir", "newfile.txt")
		args := WriteFileArgs{Path: targetPath, Content: "content"}
		argsJSON, _ := json.Marshal(args)

		result := runWriteFile(string(argsJSON))
		if !strings.Contains(result, "Error creating parents directories") {
			t.Errorf("Expected error message for path where a component is a file, but got: %s", result)
		}
	})
}
