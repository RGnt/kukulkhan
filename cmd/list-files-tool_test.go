package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRunListFiles(t *testing.T) {
	// Helper to create a temporary directory and files
	setupTestDir := func(t *testing.T, files map[string]string) string {
		tempDir, err := os.MkdirTemp("", "listfiles_test_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}

		for name, content := range files {
			filePath := filepath.Join(tempDir, name)
			if content == "" {
				if err := os.MkdirAll(filePath, 0755); err != nil {
					t.Fatalf("Failed to create subdirectory %s: %v", filePath, err)
				}
			} else {
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", filePath, err)
				}
			}
		}
		return tempDir
	}

	t.Run("HappyPath_ValidDirectory", func(t *testing.T) {
		files := map[string]string{
			"file1.txt": "content1",
			"subdir":     "", // directory
			"file2.log":  "content2",
		}
		tempDir := setupTestDir(t, files)
		defer os.RemoveAll(tempDir)

		args := ListFilesArgs{Path: tempDir}
		argsJSON, _ := json.Marshal(args)

		result := runListFiles(string(argsJSON))

		// Check if it contains expected file names and directory indicator
		expectedParts := []string{"[FILE] file1.txt", "[DIR]  subdir", "[FILE] file2.log"}
		for _, part := range expectedParts {
			if !contains(result, part) {
				t.Errorf("Expected result to contain %q, but got: %s", part, result)
			}
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJSON := `{"path": ` 
		result := runListFiles(invalidJSON)
		if !contains(result, "Error parsing arguments") {
			t.Errorf("Expected error message for invalid JSON, but got: %s", result)
		}
	})

	t.Run("NonExistentDirectory", func(t *testing.T) {
		args := ListFilesArgs{Path: "/non/existent/path/12345"}
		argsJSON, _ := json.Marshal(args)
		result := runListFiles(string(argsJSON))
		if !contains(result, "Error reading directory") {
			t.Errorf("Expected error message for non-existent directory, but got: %s", result)
		}
	})

	t.Run("EmptyDirectory", func(t *testing.T) {
		tempDir := setupTestDir(t, map[string]string{})
		defer os.RemoveAll(tempDir)

		args := ListFilesArgs{Path: tempDir}
		argsJSON, _ := json.Marshal(args)
		result := runListFiles(string(argsJSON))

		if result != "Directory is empty." {
			t.Errorf("Expected 'Directory is empty.', but got: %s", result)
		}
	})

	t.Run("PathIsFileNotDirectory", func(t *testing.T) {
		tempDir := setupTestDir(t, map[string]string{"file.txt": "content"})
		defer os.RemoveAll(tempDir)
		
		filePath := filepath.Join(tempDir, "file.txt")
		args := ListFilesArgs{Path: filePath}
		argsJSON, _ := json.Marshal(args)
		
		result := runListFiles(string(argsJSON))
		if !contains(result, "Error reading directory") {
			t.Errorf("Expected error message for path being a file, but got: %s", result)
		}
	})
}

func contains(s, substr string) bool {
	return fmt.Sprintf("%v", s) != "" && (len(s) >= len(substr)) && (s == substr || (len(substr) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
