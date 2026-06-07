package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

func runGoTests(arguments string) string {
	var args RunGoTestsArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	if args.Path == "" {
		return "Error: Path is required."
	}

	// 1. Defensive Architecture: Enforce a strict timeout.
	// If a test hangs or infinite-loops, this kills the process after 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 2. Use CommandContext instead of Command
	cmd := exec.CommandContext(ctx, "go", "test", "-v", args.Path)

	// Optional but recommended: If you want to run the tests *inside* that directory
	// instead of passing the path as a package argument, you can set the working directory:
	// cmd.Dir = args.Path
	// cmd.Args = []string{"go", "test", "-v", "./..."}

	output, err := cmd.CombinedOutput()

	if err != nil {
		// 3. Catch the timeout specifically so the LLM knows it wrote a hanging test
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Sprintf("Error: Tests timed out and were killed after 30 seconds. Check for infinite loops or deadlocks.\nPartial Output:\n%s", string(output))
		}

		// Standard test failure
		return fmt.Sprintf("Tests failed:\n%s\nError: %v", string(output), err)
	}

	return fmt.Sprintf("Tests passed successfully:\n%s", string(output))
}

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
