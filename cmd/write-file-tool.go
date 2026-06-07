package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var writeFileTool = Tool{
	Definition: writeFileToolDef,
	Guidelines: writeToolGuideline,
	Execute: func(arguments string) string {
		return runWriteFile(arguments)
	},
}

var writeToolGuideline = `
Writes a file to the local filesystem.

Usage:
- This tool will overwrite the existing file if there is one at the provided path.
- NEVER create documentation files (*.md) or README files unless explicitly requested by the User.
- Only use emojis if the user explicitly requests it. Avoid writing emojis to files unless asked.
- Never write the line numbers in to the file`

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
