package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// GenerateSandboxTool binds the running Docker Sandbox to the standard Agent Tool format
// GenerateSandboxTool binds the running Docker Sandbox to the standard Agent Tool format
func GenerateSandboxTool(sandbox *Sandbox) Tool {
	return Tool{
		Definition: APITool{
			Type: "function",
			Function: FunctionDef{
				Name:        "execute_bash",
				Description: "Execute a bash command in a secure, isolated Linux sandbox. The sandbox contains your current project files in the /workspace directory.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"command": map[string]any{
							"type":        "string",
							"description": "The shell command to execute.",
						},
					},
					"required": []string{"command"},
				},
			},
		},
		Guidelines: "Use this to run tests, compile code, search files via grep/find, or inspect system state. Do NOT start interactive processes (like vim or top).",
		Execute: func(arguments string) string {
			var args struct {
				Command string `json:"command"`
			}
			if err := json.Unmarshal([]byte(arguments), &args); err != nil {
				return fmt.Sprintf("Error parsing arguments: %v", err)
			}

			// 1. Log the exact command sent by the agent to your host terminal
			fmt.Printf("\n➔ [Tool Exec] $ %s\n", args.Command)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			output, err := sandbox.Execute(ctx, args.Command)

			// 2. Handle failures visually for the human operator
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					timeoutMsg := fmt.Sprintf("Error: Command timed out after 30 seconds.\nPartial Output:\n%s", output)
					fmt.Printf("❌ [Tool Error]: Timeout exceeded.\n")
					return timeoutMsg
				}
				fmt.Printf("❌ [Tool Error]: %v\n", err)
				return fmt.Sprintf("Execution Error: %v\nOutput:\n%s", err, output)
			}

			// 3. Log the command's stdout/stderr to your host terminal before returning it to the LLM
			fmt.Println("--- [Sandbox Stdout/Stderr] ---")
			if output == "" {
				fmt.Println("(Command returned no output)")
				fmt.Println("-------------------------------")
				return "Command executed successfully with no output."
			}

			fmt.Print(output)
			if !hasTrailingNewline(output) {
				fmt.Println() // Ensure the closing delimiter isn't stuck to the text
			}
			fmt.Println("-------------------------------")

			return output
		},
	}
}

// Helper to prevent squashed terminal lines if stdout lacks a trailing newline
func hasTrailingNewline(s string) bool {
	if len(s) == 0 {
		return true
	}
	return s[len(s)-1] == '\n'
}
