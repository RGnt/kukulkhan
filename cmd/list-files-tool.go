package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

var listFilesTool = Tool{
	Definition: listFilesToolDef,
	Guidelines: "Tool used to list files, and folders of the files system.",
	Execute: func(arguments string) string {
		return runListFiles(arguments)
	},
}
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
