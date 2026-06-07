package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const todoFilePath = ".agent_todo.json"

// TodoItem represents the internal structure stored in the JSON file
type TodoItem struct {
	ID     int    `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"` // "pending", "in-progress", "done"
}

// ==========================================
// Tool Definitions
// ==========================================

var WriteTodoTool = Tool{
	Definition: APITool{
		Type: "function",
		Function: FunctionDef{
			Name:        "write_todo",
			Description: "Creates a new task plan. This overwrites any existing plan. Pass an array of task descriptions. Always do this at the start of a complex objective.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"tasks": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"description": "List of task descriptions.",
					},
				},
				"required": []string{"tasks"},
			},
		},
	},
	Guidelines: todoToolGuideline,
	Execute:    runWriteTodo,
}

var ReadTodoTool = Tool{
	Definition: APITool{
		Type: "function",
		Function: FunctionDef{
			Name:        "read_todo",
			Description: "Reads the current task plan and returns the status of all tasks.",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
	},
	Guidelines: "Use this whenever you need to remember what step you are on, or before marking a task as done.",
	Execute:    runReadTodo,
}

var UpdateTodoTool = Tool{
	Definition: APITool{
		Type: "function",
		Function: FunctionDef{
			Name:        "update_todo",
			Description: "Updates the status of a specific task by its ID.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{
						"type":        "integer",
						"description": "The ID of the task to update.",
					},
					"status": map[string]any{
						"type":        "string",
						"description": "The new status of the task.",
						"enum":        []string{"pending", "in-progress", "done"},
					},
				},
				"required": []string{"id", "status"},
			},
		},
	},
	Guidelines: "Always update a task to 'in-progress' when starting it, and 'done' when you finish it.",
	Execute:    runUpdateTodo,
}

// ==========================================
// Execution Handlers
// ==========================================

func runWriteTodo(arguments string) string {
	var args WriteTodoArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	var todos []TodoItem
	for i, taskDesc := range args.Tasks {
		todos = append(todos, TodoItem{
			ID:     i + 1, // 1-indexed for readability
			Task:   taskDesc,
			Status: "pending", // Default state
		})
	}

	if err := saveTodos(todos); err != nil {
		return fmt.Sprintf("Error saving todo list: %v", err)
	}

	return fmt.Sprintf("Success: Created new task plan with %d tasks.", len(todos))
}

func runReadTodo(_ string) string {
	todos, err := loadTodos()
	if err != nil {
		if os.IsNotExist(err) {
			return "No task plan found. The file .agent_todo.json does not exist."
		}
		return fmt.Sprintf("Error reading todo list: %v", err)
	}

	if len(todos) == 0 {
		return "The task plan is empty."
	}

	var builder strings.Builder
	builder.WriteString("--- Current Task Plan ---\n")
	for _, t := range todos {
		builder.WriteString(fmt.Sprintf("[%d] %-12s | %s\n", t.ID, strings.ToUpper(t.Status), t.Task))
	}
	return builder.String()
}

func runUpdateTodo(arguments string) string {
	var args UpdateTodoArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	todos, err := loadTodos()
	if err != nil {
		return fmt.Sprintf("Error accessing task plan: %v", err)
	}

	updated := false
	for i, t := range todos {
		if t.ID == args.ID {
			todos[i].Status = args.Status
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Sprintf("Error: No task found with ID %d.", args.ID)
	}

	if err := saveTodos(todos); err != nil {
		return fmt.Sprintf("Error saving updated todo list: %v", err)
	}

	return fmt.Sprintf("Success: Task %d marked as '%s'.", args.ID, args.Status)
}

// ==========================================
// File I/O Helpers
// ==========================================

func loadTodos() ([]TodoItem, error) {
	data, err := os.ReadFile(todoFilePath)
	if err != nil {
		return nil, err
	}

	var todos []TodoItem
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func saveTodos(todos []TodoItem) error {
	data, err := json.MarshalIndent(todos, "", "  ") // Indent makes it human-readable
	if err != nil {
		return err
	}

	return os.WriteFile(todoFilePath, data, 0644)
}

var todoToolGuideline string = `Use this tool to create and manage a structured task list for your current coding session. This helps you track progress, organize complex tasks, and demonstrate thoroughness to the user.
It also helps the user understand the progress of the task and overall progress of their requests.

## When to Use This Tool
Use this tool proactively in these scenarios:

1. Complex multi-step tasks - When a task requires 3 or more distinct steps or actions
2. Non-trivial and complex tasks - Tasks that require careful planning or multiple operations
3. User explicitly requests todo list - When the user directly asks you to use the todo list
4. User provides multiple tasks - When users provide a list of things to be done (numbered or comma-separated)
5. After receiving new instructions - Immediately capture user requirements as todos
6. When you start working on a task - Mark it as in_progress BEFORE beginning work. Ideally you should only have one todo as in_progress at a time
7. After completing a task - Mark it as completed and add any new follow-up tasks discovered during implementation

## When NOT to Use This Tool

Skip using this tool when:
1. There is only a single, straightforward task
2. The task is trivial and tracking it provides no organizational benefit
3. The task can be completed in less than 3 trivial steps
4. The task is purely conversational or informational

NOTE that you should not use this tool if there is only one trivial task to do. In this case you are better off just doing the task directly.

## Examples of When to Use the Todo List

<example>
User: I want to add a dark mode toggle to the application settings. Make sure you run the tests and build when you're done!
Assistant: *Creates todo list with the following items:*
1. Creating dark mode toggle component in Settings page
2. Adding dark mode state management (context/store)
3. Implementing CSS-in-JS styles for dark theme
4. Updating existing components to support theme switching
5. Running tests and build process, addressing any failures or errors that occur
*Begins working on the first task*

<reasoning>
The assistant used the todo list because:
1. Adding dark mode is a multi-step feature requiring UI, state management, and styling changes
2. The user explicitly requested tests and build be run afterward
3. The assistant inferred that tests and build need to pass by adding "Ensure tests and build succeed" as the final task
</reasoning>
</example>

<example>
User: Help me rename the function getCwd to getCurrentWorkingDirectory across my project
Assistant: *Uses grep or search tools to locate all instances of getCwd in the codebase*
I've found 15 instances of 'getCwd' across 8 different files.
*Creates todo list with specific items for each file that needs updating*

<reasoning>
The assistant used the todo list because:
1. First, the assistant searched to understand the scope of the task
2. Upon finding multiple occurrences across different files, it determined this was a complex task with multiple steps
3. The todo list helps ensure every instance is tracked and updated systematically
4. This approach prevents missing any occurrences and maintains code consistency
</reasoning>
</example>


<example>
User: I need to implement these features for my e-commerce site: user registration, product catalog, shopping cart, and checkout flow.
Assistant: *Creates a todo list breaking down each feature into specific tasks based on the project architecture*
Let's start with user registration. This will involve creating database models, API endpoints, and frontend forms that integrate with your existing authentication system.

<reasoning>
The assistant used the todo list because:
1. The user provided multiple complex features to implement in a comma separated list
2. The todo list helps organize these large features into manageable tasks
3. This approach allows for tracking progress across the entire implementation
</reasoning>
</example>

<example>
User: Can you help optimize my React application? It's rendering slowly and has performance issues.
Assistant: *Reviews component structure, render patterns, state management, and data fetching*
After analyzing your codebase, I've identified several performance issues.
*Creates todo list with items like: 1) Implementing memoization for expensive calculations in ProductList, 2) Adding virtualization for long lists in Dashboard, 3) Optimizing image loading in Gallery component, 4) Fixing state update loops in ShoppingCart, 5) Reviewing bundle size and implementing code splitting*
Let's start by implementing memoization for the expensive calculations in your ProductList component.

<reasoning>
The assistant used the todo list because:
1. First, the assistant examined the codebase to identify specific performance issues
2. Based on this analysis, it identified multiple optimization opportunities
3. Performance optimization is a non-trivial task requiring multiple steps
4. The todo list helps methodically track improvements across different components
5. This systematic approach ensures all performance bottlenecks are addressed
</reasoning>
</example>

## Examples of When NOT to Use the Todo List

<example>
User: How do I print 'Hello World' in Python?
Assistant: In Python, you can print "Hello World" with this simple code:

python
print("Hello World")

This will output the text "Hello World" to the console when executed.

<reasoning>
The assistant did not use the todo list because this is a single, trivial task that can be completed in one step. There's no need to track multiple tasks or steps for such a straightforward request.
</reasoning>
</example>

<example>
User: What does the git status command do?
Assistant: The git status command shows the current state of your working directory and staging area. It displays which changes have been staged, which haven't, and which files aren't being tracked by Git. This is a fundamental command for understanding the current state of your repository before making commits.

<reasoning>
The assistant did not use the todo list because this is an informational request with no actual coding task to complete. The user is simply asking for an explanation, not for the assistant to perform multiple steps or tasks.
</reasoning>
</example>

<example>
User: Can you add a comment to the calculateTotal function to explain what it does?
Assistant: Sure, let me add a comment to the calculateTotal function to explain what it does.
* Uses the ${EDIT_TOOL_NAME} tool to add a comment to the calculateTotal function *

<reasoning>
The assistant did not use the todo list because this is a single, straightforward task confined to one location in the code. Adding a comment doesn't require tracking multiple steps or systematic organization.
</reasoning>
</example>

<example>
User: Run npm install for me and tell me what happens.
Assistant: I'll run the npm install command for you.

*Executes: npm install*

The command completed successfully. Here's the output:
[Output of npm install command]

All dependencies have been installed according to your package.json file.

<reasoning>
The assistant did not use the todo list because this is a single command execution with immediate results. There are no multiple steps to track or organize, making the todo list unnecessary for this straightforward task.
</reasoning>
</example>

## Task States and Management

1. **Task States**: Use these states to track progress:
   - pending: Task not yet started
   - in_progress: Currently working on (limit to ONE task at a time)
   - completed: Task finished successfully

   **IMPORTANT**: Task descriptions must have two forms:
   - content: The imperative form describing what needs to be done (e.g., "Run tests", "Build the project")
   - activeForm: The present continuous form shown during execution (e.g., "Running tests", "Building the project")

2. **Task Management**:
   - Update task status in real-time as you work
   - Mark tasks complete IMMEDIATELY after finishing (don't batch completions)
   - Exactly ONE task must be in_progress at any time (not less, not more)
   - Complete current tasks before starting new ones
   - Remove tasks that are no longer relevant from the list entirely

3. **Task Completion Requirements**:
   - ONLY mark a task as completed when you have FULLY accomplished it
   - If you encounter errors, blockers, or cannot finish, keep the task as in_progress
   - When blocked, create a new task describing what needs to be resolved
   - Never mark a task as completed if:
     - Tests are failing
     - Implementation is partial
     - You encountered unresolved errors
     - You couldn't find necessary files or dependencies

4. **Task Breakdown**:
   - Create specific, actionable items
   - Break complex tasks into smaller, manageable steps
   - Use clear, descriptive task names
   - Always provide both forms:
     - content: "Fix authentication bug"
     - activeForm: "Fixing authentication bug"

When in doubt, use this tool. Being proactive with task management demonstrates attentiveness and ensures you complete all requirements successfully.
`
