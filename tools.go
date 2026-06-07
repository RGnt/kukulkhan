package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var (
	hasBackup         bool
	lastBackupPath    string
	lastBackupContent []byte
	wasNewlyCreated   bool
)

var listFilesToolDef = APITool{
	Type: "function",
	Function: FunctionDef{
		Name:        "list_files",
		Description: "List all files and directories in the specified path",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
			"required": []string{"path"},
		},
	},
}

var readFileToolDef = APITool{
	Type: "function",
	Function: FunctionDef{
		Name:        "read_file",
		Description: "Read the contents of a file. You can optionally specify a start_line and end_line to read a specific chunk. Lines are 1-indexed.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The absolute or relative path to the file",
				},
				"start_line": map[string]any{
					"type":        "integer",
					"description": "Optional. The line number to start reading from (1-indexed).",
				},
				"end_line": map[string]any{
					"type":        "integer",
					"description": "Optional. The line number to stop reading at (inclusive).",
				},
			},
			"required": []string{"path"},
		},
	},
}

var writeFileToolDef = APITool{
	Type: "function",
	Function: FunctionDef{
		Name:        "write_file",
		Description: "Write entire content to a file. Overwrites existing files and creates missing directories. Automatically backs up the previous state.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The absolute or relative path to the file.",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "The complete new content of the file.",
				},
			},
			"required": []string{"path", "content"},
		},
	},
}

var revertFileToolDef = APITool{
	Type: "function",
	Function: FunctionDef{
		Name:        "revert_file",
		Description: "Reverts the last write_file operation. Use this immediately if you realize your last file write was incorrect.",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{}, // Takes no arguments
		},
	},
}

func runListFiles(arguments string) string {
	var args ListFilesArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	entries, err := os.ReadDir(args.Path)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}

	if len(entries) == 0 {
		return "Directory is empty."
	}

	var result bytes.Buffer
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if entry.IsDir() {
			result.WriteString(fmt.Sprintf("[DIR]  %s\n", entry.Name()))
		} else {
			result.WriteString(fmt.Sprintf("[FILE] %s (%d bytes)\n", entry.Name(), info.Size()))
		}
	}

	return result.String()
}

func runReadFile(arguments string) string {
	var args ReadFileArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("error parsing arguments: %v", err)
	}

	file, err := os.Open(args.Path)
	if err != nil {
		return fmt.Sprintf("error opening file: %v", err)
	}
	defer file.Close()

	start := 1
	if args.StartLine != nil && *args.StartLine > 0 {
		start = *args.StartLine
	}

	end := -1
	if args.EndLine != nil && *args.EndLine >= start {
		end = *args.EndLine
	}

	var result bytes.Buffer
	scanner := bufio.NewScanner(file)
	currentLine := 1

	for scanner.Scan() {
		if currentLine >= start {
			result.WriteString(fmt.Sprintf("%d | %s\n", currentLine, scanner.Text()))
		}

		if end != -1 && currentLine >= end {
			break
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("error reading file during scan: %v", err)
	}

	if result.Len() == 0 {
		return "Success: File was read, but it was empty."
	}

	return result.String()
}

func runWriteFile(arguments string) string {
	var args WriteFileArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error Parsing arguments: %v", err)
	}

	info, err := os.Stat(args.Path)
	if err == nil && !info.IsDir() {
		// File exists, read, and store its current content
		existingContent, readErr := os.ReadFile(args.Path)
		if readErr != nil {
			return fmt.Sprintf("Error: File exists but failed to read for backup: %v", readErr)
		}
		lastBackupContent = existingContent
		wasNewlyCreated = false
	} else {
		lastBackupContent = nil
		wasNewlyCreated = true
	}

	lastBackupPath = args.Path
	hasBackup = true

	dir := filepath.Dir(args.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Sprintf("Error creating parents directories %v", err)
	}

	if err := os.WriteFile(args.Path, []byte(args.Content), 0644); err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}

	return fmt.Sprintf("Success: Wrote %d bytes to '%s'. Previous state backed up.", len(args.Content), args.Path)
}

func runRevertFile() string {
	if !hasBackup {
		return "Error: No previous file write operations to revert."
	}
	if wasNewlyCreated {
		if err := os.Remove(lastBackupPath); err != nil {
			return fmt.Sprintf("Error deleting newly created file during revert: %v", err)
		}

		hasBackup = false
		return fmt.Sprintf("Success: Reverted creation of '%s' by deleting it", lastBackupPath)
	}

	if err := os.WriteFile(lastBackupPath, lastBackupContent, 0644); err != nil {
		return fmt.Sprintf("Error restoring previous content: %v", err)
	}
	hasBackup = false
	return fmt.Sprintf("Success: Restored previous contents of '%s'.", lastBackupPath)
}
