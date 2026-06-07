package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

var runGoTestsTool = Tool{
	Definition: runGoTestsToolDef,
	Guidelines: "Tool used to run Go tests in a specified directory. It executes 'go test -v' and returns the output.",
	Execute: func(arguments string) string {
		return runGoTests(arguments)
	},
}

var runGoTestsToolDef = APITool{
	Type: "function",
	Function: FunctionDef{
		Name:        "run_go_tests",
		Description: "Run Go tests in the specified directory",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
			"required": []string{"path"},
		},
	},
}

func runGoTests(arguments string) string {
	var args RunGoTestsArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	if args.Path == "" {
		return "Error: Path is required."
	}

	// Use 'go test -v' command
	cmd := exec.Command("go", "test", "-v", args.Path)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Sprintf("Tests failed or error occurred:\n%s\nError: %v", string(output), err)
	}

	return fmt.Sprintf("Tests passed successfully:\n%s", string(output))
}
