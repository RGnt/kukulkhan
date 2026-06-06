package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// RunAgentLoop handles the core ReAct (Reasoning and Acting) cycle with streaming
func RunAgentLoop(history []Message) []Message {
	serverURL := "http://localhost:8080/v1/chat/completions"
	temp := 0.0
	maxSteps := 5

	client := &http.Client{Timeout: 2 * time.Minute}

	for step := 1; step <= maxSteps; step++ {
		fmt.Printf("\n[Agent Step %d]\n", step)

		reqBody := ChatRequest{
			Model:       "local-model",
			Messages:    history,
			Temperature: &temp,
			Stream:      true,
			Tools:       tools,
		}

		jsonBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Agent network error: %v", err)
		}

		responseMsg := streamAndBuildMessage(resp)
		resp.Body.Close()

		// The model invoked a tool
		if len(responseMsg.ToolCalls) > 0 {
			history = append(history, responseMsg)

			for _, tc := range responseMsg.ToolCalls {
				fmt.Printf("\n>> Executing Tool: %s\n", tc.Function.Name)
				toolResult := executeTool(tc.Function.Name, tc.Function.Arguments)
				history = append(history, Message{
					Role:       "tool",
					Content:    toolResult,
					ToolCallID: tc.ID,
				})
			}
			continue
		}

		// The model provided a final text answer
		if responseMsg.Content != "" {
			history = append(history, responseMsg)
			fmt.Println()
			return history
		}
	}

	fmt.Println("\n[Agent Error]: Reached maximum steps without returning a final answer.")
	return history
}

func streamAndBuildMessage(resp *http.Response) Message {
	scanner := bufio.NewScanner(resp.Body)

	finalMsg := Message{Role: "assistant"}
	var contentBuilder strings.Builder

	type toolCallBuilder struct {
		ID        string
		Type      string
		Name      string
		Arguments strings.Builder
	}
	toolBuilders := make(map[int]*toolCallBuilder)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta

			// If it's standard text, print it instantly and save it
			if delta.Content != "" {
				fmt.Print(delta.Content)
				contentBuilder.WriteString(delta.Content)
			}

			// If it's a tool call fragment, accumulate it silently
			if len(delta.ToolCalls) > 0 {
				tc := delta.ToolCalls[0]
				idx := tc.Index

				// Initialize the builder for this index if it doesn't exist
				if toolBuilders[idx] == nil {
					toolBuilders[idx] = &toolCallBuilder{Type: "function"}
				}

				if tc.ID != "" {
					toolBuilders[idx].ID = tc.ID
				}
				if tc.Function.Name != "" {
					toolBuilders[idx].Name = tc.Function.Name
				}
				if tc.Function.Arguments != "" {
					toolBuilders[idx].Arguments.WriteString(tc.Function.Arguments)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("\n[Stream reading error: %v]\n", err)
	}

	// Assemble the final Message struct
	finalMsg.Content = contentBuilder.String()

	// Convert our map of builders into the final ToolCall slice.
	for i := 0; i < len(toolBuilders); i++ {
		if tb, exists := toolBuilders[i]; exists {
			finalMsg.ToolCalls = append(finalMsg.ToolCalls, ToolCall{
				ID:   tb.ID,
				Type: tb.Type,
				Function: struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				}{
					Name:      tb.Name,
					Arguments: tb.Arguments.String(),
				},
			})
		}
	}

	return finalMsg
}
