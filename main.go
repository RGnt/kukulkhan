package main

import (
	"fmt"
)

func main() {
	messages := []Message{
		{Role: "user", Content: "First, list the files in the current directory ('.'). Then, calculate the speed of a car that travels 200 meters in 15 seconds."},
	}

	fmt.Println("User: ", messages[0].Content)
	fmt.Println("--------------------------------------------------")

	RunAgentLoop(messages)
}
