package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// Registry of tools provided to the LLM
var tools = []Tool{
	{
		Type: "function",
		Function: FunctionDef{
			Name:        "calculate_speed",
			Description: "Calculate average speed given distance, and time",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"distance": map[string]any{"type": "number"},
					"time":     map[string]any{"type": "number"},
				},
				"required": []string{"distance", "time"},
			},
		},
	},
	{
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
	},
}

// executeTool acts as the central router for LLM function calls
func executeTool(name string, arguments string) string {
	switch name {
	case "calculate_speed":
		return runCalculateSpeed(arguments)
	case "list_files":
		return runListFiles(arguments)
	default:
		return fmt.Sprintf("Error: Unknown tool '%s'", name)
	}
}

// Individual tools
func runCalculateSpeed(arguments string) string {
	var args CalculateSpeedArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	speed := args.Distance / args.Time
	return fmt.Sprintf(`{"speed_meters_per_second": %f}`, speed)
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
