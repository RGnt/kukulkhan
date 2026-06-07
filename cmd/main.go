package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("======================================================")
	fmt.Println(" Kukulkhan CLI Initialized ")
	fmt.Println(" Type your prompt below. Type /quit to exit.")
	fmt.Println("======================================================")

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
		"You are the main coordinator. Delegate complex tasks.",
		"Gemma 4 12b",
		0.7,
		[]Tool{listFilesTool, readFilesTool, writeFileTool, revertFileTool, WriteTodoTool, ReadTodoTool, UpdateTodoTool, runGoTestsTool},
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
