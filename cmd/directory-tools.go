package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ==========================================
// Tool Definitions
// ==========================================

var AddDirectoryTool = Tool{
	Definition: APITool{
		Type: "function",
		Function: FunctionDef{
			Name:        "add_directory",
			Description: "Creates a new directory. Automatically creates any necessary parent directories.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The absolute or relative path of the directory to create.",
					},
				},
				"required": []string{"path"},
			},
		},
	},
	Guidelines: "Use this when you need to create a new folder structure before writing files into it.",
	Execute:    runAddDirectory,
}

var RemoveDirectoryTool = Tool{
	Definition: APITool{
		Type: "function",
		Function: FunctionDef{
			Name:        "remove_directory",
			Description: "Deletes a directory. By default, it only deletes empty directories.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path of the directory to delete.",
					},
					"recursive": map[string]any{
						"type":        "boolean",
						"description": "Set to true to delete the directory and all of its contents (files and subdirectories). Default is false.",
					},
				},
				"required": []string{"path", "recursive"},
			},
		},
	},
	Guidelines: "Only set recursive to true if you are absolutely sure you want to permanently delete all files inside the target directory.",
	Execute:    runRemoveDirectory,
}

// ==========================================
// Execution Handlers
// ==========================================

func runAddDirectory(arguments string) string {
	var args AddDirectoryArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	if args.Path == "" {
		return "Error: Path is required."
	}

	// 0755 gives read/write/execute to the owner, and read/execute to everyone else.
	// os.MkdirAll is equivalent to 'mkdir -p' in Linux.
	if err := os.MkdirAll(args.Path, 0755); err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}

	return fmt.Sprintf("Success: Directory created at '%s'.", args.Path)
}

func runRemoveDirectory(arguments string) string {
	var args RemoveDirectoryArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	if args.Path == "" {
		return "Error: Path is required."
	}

	// Prevent the agent from accidentally deleting the root or current working directory
	if args.Path == "." || args.Path == "/" || args.Path == "./" {
		return "Error: Deleting the current working directory or root directory is strictly prohibited."
	}

	// Check if the path actually exists and is a directory
	info, err := os.Stat(args.Path)
	if os.IsNotExist(err) {
		return fmt.Sprintf("Error: Directory '%s' does not exist.", args.Path)
	}
	if !info.IsDir() {
		return fmt.Sprintf("Error: '%s' is a file, not a directory. Use a file removal tool instead.", args.Path)
	}

	// Execute deletion based on the recursive flag
	if args.Recursive {
		// Equivalent to 'rm -rf'
		if err := os.RemoveAll(args.Path); err != nil {
			return fmt.Sprintf("Error recursively deleting directory: %v", err)
		}
		return fmt.Sprintf("Success: Directory '%s' and all its contents were permanently deleted.", args.Path)
	} else {
		// Equivalent to 'rmdir'. Will fail natively if the directory contains files.
		if err := os.Remove(args.Path); err != nil {
			return fmt.Sprintf("Error deleting directory (it may not be empty): %v", err)
		}
		return fmt.Sprintf("Success: Empty directory '%s' was deleted.", args.Path)
	}
}
