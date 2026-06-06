// main.go

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
			Role:    "system",
			Content: "You are a helpful CLI agent running on a local machine. You have access to tools to calculate speed and list files. Answer questions concisely.",
		},
	}

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

		history = RunAgentLoop(history)
	}
}
