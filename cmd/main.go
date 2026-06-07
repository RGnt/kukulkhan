package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	LoadConfig()

	fmt.Println("======================================================")
	fmt.Println(" Kukulkhan CLI Initialized ")
	fmt.Printf(" Workspace: %s\n", Config.WorkspaceDir)
	fmt.Printf(" Streaming: %v\n", Config.StreamResponse)
	fmt.Println(" Type your prompt below. Type /quit to exit.")
	fmt.Println("======================================================")

	fmt.Println("[Sandbox] Initializing secure container environment...")
	ctx := context.Background()
	sandbox, err := NewSandbox(ctx, ".")
	if err != nil {
		log.Fatalf("Fatal: Could not start sandbox. Is Docker running? Error: %v", err)
	}

	defer sandbox.Close(context.Background())

	bashTool := GenerateSandboxTool(sandbox)

	history := []Message{
		{
			Role: "system",
			Content: `
			You are senior software engineering agent.
			Before working on any multi-step task, ALWAYS call todo_write first to write your complete plant.
			Execute each step in order.
			Call todo_update after completing each step.
			`,
		},
	}

	mainAgent := NewAgent(
		"Coordinator",
		"You are a local developer agent. Use execute_bash to explore and test the project.",
		"Gemma 4 12b",
		0.7,
		[]Tool{bashTool, readFilesTool, WriteTodoTool, ReadTodoTool, UpdateTodoTool},
	)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)

		// Handle empty inputs gracefully
		if input == "" {
			continue
		}

		if input == "/quit" {
			fmt.Println("Shutting down agent. Goodbye!")
			break
		}

		history = append(history, Message{
			Role:    "user",
			Content: input,
		})

		mainAgent.Run(history)
	}
}
