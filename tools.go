package main

import (
	"bufio"
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
	{
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
	},
}

// executeTool acts as the central router for LLM function calls
func executeTool(name string, arguments string) string {
	switch name {
	case "list_files":
		return runListFiles(arguments)
	case "read_file":
		return runReadFile(arguments)
	default:
		return fmt.Sprintf("Error: Unknown tool '%s'", name)
	}
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
