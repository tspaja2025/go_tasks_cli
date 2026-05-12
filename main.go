package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
	TimeFormat       = "2006-01-02 15:04:05"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	tasks, err = executeCommand(command, args, tasks)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
	}
}

// Execute commands
func executeCommand(command string, args []string, tasks []Task) ([]Task, error) {
	switch command {
	case "add":
		return handleAdd(args, tasks)
	case "update":
		return handleUpdate(args, tasks)
	case "delete":
		return handleDelete(args, tasks)
	case "batch-delete":
		return handleBatchDelete(args, tasks)
	case "list":
		return handleList(tasks), nil
	case "list-status":
		return handleListStatus(args, tasks)
	case "mark-with-status":
		return handleMarkStatus(args, tasks)
	case "get-task-by-id":
		return handleGetTask(args, tasks)
	case "clear-all-tasks":
		return handleClearAll(tasks)
	case "help", "-h", "--help":
		printHelp()
		return tasks, nil
	default:
		return tasks, fmt.Errorf("Unknown command: %s", command)
	}
}

// Handle task add
func handleAdd(args []string, tasks []Task) ([]Task, error) {
	if len(args) < 2 {
		return tasks, fmt.Errorf("Usage: add <description> <status>")
	}

	currentTime := time.Now()
	newTask := Task{
		ID:          nextID(tasks),
		Description: args[0],
		Status:      args[1],
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	tasks = append(tasks, newTask)
	fmt.Println("Task added successfully")
	return tasks, nil
}

// Handle task update
func handleUpdate(args []string, tasks []Task) ([]Task, error) {
	if len(args) < 2 {
		return tasks, fmt.Errorf("Usage: update <id> <new-description>")
	}

	id, err := getID(args[0])
	if err != nil {
		return tasks, err
	}

	return updateTask(tasks, id, args[1]), nil
}

func updateTask(tasks []Task, id int, newDescription string) []Task {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now()
			fmt.Println("Task updated successfully")
			return tasks
		}
	}

	fmt.Println("Task not found")
	return tasks
}

// Handle task delete
func handleDelete(args []string, tasks []Task) ([]Task, error) {
	if len(args) < 1 {
		return tasks, fmt.Errorf("Usage: delete <id>")
	}

	id, err := getID(args[0])
	if err != nil {
		return tasks, err
	}

	return deleteTask(tasks, id), nil
}

func deleteTask(tasks []Task, id int) []Task {
	newTasks := []Task{}

	found := false
	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		} else {
			found = true
		}
	}

	if found {
		fmt.Println("Task deleted successfully.")
	} else {
		fmt.Println("Task not found.")
	}

	return newTasks
}

// Handle task batch delete
func handleBatchDelete(args []string, tasks []Task) ([]Task, error) {
	if len(args) < 1 {
		return tasks, fmt.Errorf("Usage: batch-delete <id-range>")
	}

	ids, err := parseIDRange(args[0])
	if err != nil {
		return tasks, err
	}

	return batchDeleteTasks(tasks, ids), nil
}

func batchDeleteTasks(tasks []Task, ids []int) []Task {
	idSet := make(map[int]struct{})

	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	var newTasks []Task

	for _, task := range tasks {
		if _, exists := idSet[task.ID]; !exists {
			newTasks = append(newTasks, task)
		}
	}

	fmt.Println("Batch delete completed.")

	return newTasks
}

// Handle list tasks
func handleList(tasks []Task) []Task {
	listTasks(tasks)
	return tasks
}

func listTasks(tasks []Task) {
	descriptionWidth := getMaxDescriptionWidth(tasks)

	headerFormat := fmt.Sprintf("%%-6s | %%-%ds | %%-12s | %%-20s | %%-20s\n", descriptionWidth)
	rowFormat := fmt.Sprintf("%%-6d | %%-%ds | %%-12s | %%-20s | %%-20s\n", descriptionWidth)

	fmt.Println("\n=============================================== Task List ====================================================")

	fmt.Printf(headerFormat, "ID", "Description", "Status", "Created At", "Updated At")

	fmt.Println(strings.Repeat("-", descriptionWidth+70))

	for _, task := range tasks {
		fmt.Printf(rowFormat,
			task.ID,
			task.Description,
			task.Status,
			task.CreatedAt.Format(TimeFormat),
			task.UpdatedAt.Format(TimeFormat),
		)
	}
	fmt.Println(strings.Repeat("=", descriptionWidth+70))
}

// Handle list statuses
func handleListStatus(args []string, tasks []Task) ([]Task, error) {
	if len(args) < 1 {
		return tasks, fmt.Errorf("Usage: list-status <status>")
	}

	status := args[0]
	if status != StatusTodo && status != StatusInProgress && status != StatusDone {
		return tasks, fmt.Errorf("Invalid status: %s", status)
	}

	listTasksByStatus(tasks, status)
	return tasks, nil
}

func listTasksByStatus(tasks []Task, status string) {
	for _, task := range tasks {
		if task.Status == status {
			fmt.Printf("ID: %d | Description: %s | Status: %s\n",
				task.ID, task.Description, task.Status)
		}
	}
}

// Handle mark status
func handleMarkStatus(args []string, tasks []Task) ([]Task, error) {
	if len(args) < 2 {
		return tasks, fmt.Errorf("Usage: mark-with-status <id> <status>")
	}

	id, err := getID(args[0])
	if err != nil {
		return tasks, err
	}

	return markWithStatus(tasks, id, args[1]), nil
}

func markWithStatus(tasks []Task, id int, status string) []Task {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			fmt.Println("Task status updated successfully")
			return tasks
		}
	}
	return tasks
}

// Handle get task
func handleGetTask(args []string, tasks []Task) ([]Task, error) {
	if len(args) < 1 {
		return tasks, fmt.Errorf("Usage: get-tasks-by-id <id>")
	}

	id, err := getID(args[0])
	if err != nil {
		return tasks, err
	}

	return getTaskByID(tasks, id), nil
}

func getTaskByID(tasks []Task, id int) []Task {
	for _, task := range tasks {
		if task.ID == id {
			fmt.Printf("ID: %d | Description: %s | Status: %s | CreatedAt: %s | UpdatedAt: %s\n",
				task.ID,
				task.Description,
				task.Status,
				task.CreatedAt,
				task.UpdatedAt,
			)
		}
	}
	return tasks
}

// Handle clear all
func handleClearAll(tasks []Task) ([]Task, error) {
	return clearAllTasks(tasks), nil
}

func clearAllTasks(tasks []Task) []Task {
	var confirmation string
	fmt.Print("Are you sure? This command will delete ALL tasks. (yes/no)")
	fmt.Scanln(&confirmation)
	if confirmation == "yes" || confirmation == "y" {
		fmt.Println("All tasks deleted successfully.")
		return []Task{}
	}
	fmt.Println("Task deletion cancelled.")
	return tasks
}

// Load Tasks
func loadTasks() ([]Task, error) {
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("Failed to read tasks file: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(file, &tasks); err != nil {
		fmt.Println("Warning: tasks.json is corrupted, starting fresh")
		return []Task{}, nil
	}

	return tasks, nil
}

// Save Tasks
func saveTasks(tasks []Task) error {
	byteValue, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to encode tasks: %w", err)
	}

	if err := os.WriteFile("tasks.json", byteValue, 0644); err != nil {
		return fmt.Errorf("Failed to save tasks: %w", err)
	}

	return nil
}

// Print usage
func printUsage() {
	fmt.Println("No arguments provided")
	fmt.Println("Usage: <command> [arguments]")
	fmt.Println("Run 'help' for available commands")
}

// Print help
func printHelp() {
	fmt.Println(`
Task Tracker CLI - Available Commands:
  add <description> <status>     - Add a new task
  update <id> <description>      - Update task description
  delete <id>                    - Delete a task
  batch-delete <start-end>       - Batch delete tasks
  clear-all-tasks                - Clear all tasks
  list                           - List all tasks
  list-status <status>           - List tasks by status
  get-task-by-id <id>            - Get task by id
  mark-with-status <id> <status> - Change task status
  show <id>                      - Show task details
  help                           - Show this help
    `)
}

func getID(id string) (int, error) {
	idStr, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("Invalid ID: %s", id)
	}
	return idStr, nil
}

// Find maximum amount IDs and increment
func nextID(tasks []Task) int {
	maxId := 0

	for _, task := range tasks {
		if task.ID > maxId {
			maxId = task.ID
		}
	}

	return maxId + 1
}

// Parse range of ids
func parseIDRange(input string) ([]int, error) {
	var ids []int

	if strings.Contains(input, "-") {
		parts := strings.Split(input, "-")

		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid range format")
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("Invalid start id")
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("Invalid end id")
		}

		if start > end {
			return nil, fmt.Errorf("Start id cannot be greater than end id")
		}

		for i := start; i <= end; i++ {
			ids = append(ids, i)
		}

		return ids, nil
	}

	// Single ID
	id, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("Invalid ID")
	}

	ids = append(ids, id)

	return ids, nil
}

// Validate status
func isValidStatus(status string) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

// Find maximum description width
func getMaxDescriptionWidth(tasks []Task) int {
	maxWidth := len("Description")

	for _, task := range tasks {
		if len(task.Description) > maxWidth {
			maxWidth = len(task.Description)
		}
	}

	return maxWidth
}
