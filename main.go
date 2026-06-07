// main.go

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var listFilesTool = Tool{
	Definition: listFilesToolDef,
	Guidelines: "Use when user asks to read a file",
	Execute: func(arguments string) string {
		return runListFiles(arguments)
	},
}

var readFilesTool = Tool{
	Definition: readFileToolDef,
	Guidelines: "Use to read contents of the file",
	Execute: func(arguments string) string {
		return runReadFile(arguments)
	},
}

var writeFileTool = Tool{
	Definition: writeFileToolDef,
	Guidelines: "Use to write to a file",
	Execute: func(arguments string) string {
		return runWriteFile(arguments)
	},
}

var revertFileTool = Tool{
	Definition: revertFileToolDef,
	Guidelines: "Use to revert the contents of a file",
	Execute: func(arguments string) string {
		return runRevertFile()
	},
}

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

	mainAgent := NewAgent(
		"Coordinator",
		"You are the main coordinator. Delegate complex tasks.",
		"Gemma 4 12b",
		0.7,
		[]Tool{listFilesTool, readFilesTool, writeFileTool, revertFileTool},
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
