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
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"` // "todo", "in-progress", "done"
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const (
	TimeFormat = "2006-01-02 15:04:05"
)

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

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("Invalid start id")
		}

		end, err := strconv.Atoi(parts[1])
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

// Add tasks
func addTask() {
	// TODO:
}

// Delete task
func deleteTask(tasks []Task, id int) []Task {
	newTasks := []Task{}

	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		}
	}
	fmt.Println("Task deleted successfully.")
	return newTasks
}

// Batch delete
func batchDeleteTasks(tasks []Task, ids []int) []Task {
	mapIds := make(map[int]bool)

	for _, id := range ids {
		mapIds[id] = true
	}

	var newTasks []Task

	for _, task := range tasks {
		if !mapIds[task.ID] {
			newTasks = append(newTasks, task)
		}
	}

	fmt.Println("Batch delete completed.")

	return newTasks
}

// Update task
func updateTask(tasks []Task, id int, newDescription string) []Task {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now().Format(TimeFormat)
			fmt.Println("Task updated successfully")
			return tasks
		}
	}

	fmt.Println("Task not found")
	return tasks
}

// List tasks
func listTasks(tasks []Task) {
	fmt.Println("\n========================================= Task List ==============================================")
	fmt.Printf("%-6s | %-48s | %-12s | %-20s | %-20s\n", "ID", "Description", "Status", "Created At", "Updated At")
	fmt.Println(strings.Repeat("-", 98))
	for _, task := range tasks {
		fmt.Printf("%-6d | %-48s | %-12s | %-20s | %-20s\n",
			task.ID,
			task.Description,
			task.Status,
			task.CreatedAt,
			task.UpdatedAt,
		)
	}
	fmt.Println(strings.Repeat("=", 98))
}

// List tasks by status
func listTasksByStatus(tasks []Task, status string) {
	for _, task := range tasks {
		if task.Status == status {
			fmt.Printf("ID: %d | Description: %s | Status: %s\n",
				task.ID, task.Description, task.Status)
		}
	}
}

// Mark task with new status
func markWithStatus(tasks []Task, id int, status string) []Task {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Format(TimeFormat)
			fmt.Println("Task status updated successfully")
			return tasks
		}
	}
	return tasks
}

// Get status by id
func getStatusById(tasks []Task, id int) []Task {
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

// Clear all tasks
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

// Show help
func showHelp() {
	fmt.Println(`
Task Tracker CLI - Available Commands:
  add <description> <status>     - Add a new task
  update <id> <description>      - Update task description
  delete <id>                    - Delete a task
  batch-delete <id> - <id>       - Batch delete tasks
  clear-all-tasks                - Clear all tasks
  list                           - List all tasks
  list-status <status>           - List tasks by status
  get-task-by-id <id>            - Get task by id
  mark-with-status <id> <status> - Change task status
  show <id>                      - Show task details
  help                           - Show this help
    `)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		fmt.Println("Usage: <command> [arguments]")
		return
	}

	currentTime := time.Now().Format(TimeFormat)

	file, err := os.ReadFile("tasks.json")
	var tasks []Task
	if err == nil {
		if err := json.Unmarshal(file, &tasks); err != nil {
			fmt.Println("Warning: tasks.json is corrupted, starting fresh")
			tasks = []Task{}
		}
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: add <description> <status>")
			return
		}
		newTask := Task{
			ID:          nextID(tasks),
			Description: os.Args[2],
			Status:      os.Args[3],
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		}
		tasks = append(tasks, newTask)
		fmt.Println("Task added successfully.")

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: update <id> <new-description>")
			return
		}
		id, err := getID(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		newDescription := os.Args[3]
		tasks = updateTask(tasks, id, newDescription)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: delete <id>")
			return
		}
		id, err := getID(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		tasks = deleteTask(tasks, id)
	case "batch-delete":
		ids, err := parseIDRange(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		tasks = batchDeleteTasks(tasks, ids)

	case "list":
		listTasks(tasks)

	case "list-status":
		if len(os.Args) < 3 {
			fmt.Println("Usage: list-status <status>")
			return
		}
		status := os.Args[2]
		if status == "todo" || status == "in-progress" || status == "done" {
			listTasksByStatus(tasks, status)
		} else {
			fmt.Println("Invalid status")
		}

	case "mark-with-status":
		if len(os.Args) < 4 {
			fmt.Println("Usage: mark-with-status <id> <status>")
			return
		}
		id, err := getID(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		newStatus := os.Args[3]
		tasks = markWithStatus(tasks, id, newStatus)

	case "get-task-by-id":
		if len(os.Args) < 3 {
			fmt.Println("Usage: get-task-by-id <id>")
			return
		}
		id, err := getID(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		tasks = getStatusById(tasks, id)

	case "clear-all-tasks":
		tasks = clearAllTasks(tasks)

	case "help", "-h", "--help":
		showHelp()

	default:
		fmt.Println("Unknown command:")
	}

	byteValue, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Failed to encode tasks:", err)
		return
	}
	if err := os.WriteFile("tasks.json", byteValue, 0644); err != nil {
		fmt.Println("Failed to save tasks:", err)
		return
	}
}
