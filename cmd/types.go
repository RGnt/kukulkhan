package main

// API Communication types
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature *float64  `json:"temperature,omitempty"`
	Tools       []APITool `json:"tools,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// API Tool types
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type FunctionDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type APITool struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

// Streaming Types
type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content   string `json:"content,omitempty"`
			ToolCalls []struct {
				Index    int    `json:"index"`
				ID       string `json:"id,omitempty"`
				Type     string `json:"type,omitempty"`
				Function struct {
					Name      string `json:"name,omitempty"`
					Arguments string `json:"arguments,omitempty"`
				} `json:"function,omitempty"`
			} `json:"tool_calls,omitempty"`
		} `json:"delta"`
	} `json:"choices"`
}

// Internal agent types
type Tool struct {
	Definition APITool
	Guidelines string
	Execute    func(arguments string) string
}

type Agent struct {
	Name        string
	Role        string
	Model       string
	Temperature float64
	Tools       map[string]Tool
}

// Tool argument types

type ListFilesArgs struct {
	Path string `json:"path"`
}

type ReadFileArgs struct {
	Path      string `json:"path"`
	StartLine *int   `json:"start_line,omitempty"`
	EndLine   *int   `json:"end_line,omitempty"`
}

type WriteFileArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type WriteTodoArgs struct {
	Tasks []string `json:"tasks"`
}

type UpdateTodoArgs struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type RunGoTestsArgs struct {
	Path string `json:"path"`
}

type AddDirectoryArgs struct {
	Path string `json:"path"`
}

type RemoveDirectoryArgs struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"` // Forces the LLM to explicitly request deep deletion
}
