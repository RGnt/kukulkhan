package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

var readToolGuideline = `
Reads a file from the local filesystem. You can access any file directly by using this tool.
Assume this tool is able to read all files on the machine. If the User provides a path to a file assume that path is valid. It is okay to read a file that does not exist; an error will be returned.

Usage:
- The path parameter must be an absolute path, not a relative path
- By default, it reads whole file starting from the beginning of the file
- File can be read in parts by giving start_line, and end_line parameters
- This tool can only read files, not directories. To list files in a directory, use the registered list_files tool.
- If you read a file that exists but has empty contents you will receive a system reminder warning in place of file contents.`

var readFilesTool = Tool{
	Definition: readFileToolDef,
	Guidelines: readToolGuideline,
	Execute: func(arguments string) string {
		return runReadFile(arguments)
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
