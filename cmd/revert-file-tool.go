package main

import (
	"fmt"
	"os"
)

var (
	hasBackup         bool
	lastBackupPath    string
	lastBackupContent []byte
	wasNewlyCreated   bool
)

var revertFileTool = Tool{
	Definition: revertFileToolDef,
	Guidelines: "Use to revert the contents of a file",
	Execute: func(arguments string) string {
		return runRevertFile()
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
