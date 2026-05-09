package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"` // "todo", "in-progress", "done"
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func getId(id string) int {
	idStr, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Invalid ID.")
	}
	return idStr
}

// Delete task from json file
func deleteTask(tasks []Task, id int) []Task {
	newTasks := []Task{}

	for _, task := range tasks {
		if task.Id != id {
			newTasks = append(newTasks, task)
		}
	}
	fmt.Println("Task deleted successfully.")
	return newTasks
}

// Update task
func updateTask(tasks []Task, id int, newDescription string) []Task {
	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			fmt.Println("Task updated successfully")
			return tasks
		}
	}

	fmt.Println("Task not found")
	return tasks
}

// List tasks
func listTasks(tasks []Task) {
	for _, task := range tasks {
		fmt.Printf("ID: %d | Description: %s | Status: %s | CreatedAt: %s | UpdatedAt: %s\n",
			task.Id,
			task.Description,
			task.Status,
			task.CreatedAt,
			task.UpdatedAt,
		)
	}
}

// List tasks by status
func listTasksByStatus(tasks []Task, status string) {
	for _, task := range tasks {
		if task.Status == status {
			fmt.Printf("ID: %d | Description: %s | Status: %s\n",
				task.Id, task.Description, task.Status)
		}
	}
}

// Mark task with new status
func markWithStatus(tasks []Task, id int, status string) []Task {
	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			fmt.Println("Task status updated successfully")
			return tasks
		}
	}
	return tasks
}

// Get status by id
func getStatusById(tasks []Task, id int) []Task {
	for _, task := range tasks {
		if task.Id == id {
			fmt.Printf("ID: %d | Description: %s | Status: %s | CreatedAt: %s | UpdatedAt: %s\n",
				task.Id,
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

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	file, err := os.ReadFile("tasks.json")
	var tasks []Task

	if err == nil {
		json.Unmarshal(file, &tasks)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: add <description> <status>")
			return
		}
		newTask := Task{
			Id:          len(tasks) + 1,
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
		id := getId(os.Args[2])
		newDescription := os.Args[3]
		tasks = updateTask(tasks, id, newDescription)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: delete <id>")
			return
		}
		id := getId(os.Args[2])
		tasks = deleteTask(tasks, id)

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
		if len(os.Args) < 3 {
			fmt.Println("Usage: mark-with-status <id> <status>")
			return
		}
		id := getId(os.Args[2])
		newStatus := os.Args[3]
		tasks = markWithStatus(tasks, id, newStatus)

	case "get-task-by-id":
		if len(os.Args) < 3 {
			fmt.Println("Usage: get-status-by-id <id>")
			return
		}
		id := getId(os.Args[2])
		tasks = getStatusById(tasks, id)

	case "clear-all-tasks":
		tasks = clearAllTasks(tasks)

	case "help", "-h", "--help":
		showHelp()

	default:
		fmt.Println("Unknown command:")
	}

	byteValue, _ := json.MarshalIndent(tasks, "", "  ")
	err = os.WriteFile("tasks.json", byteValue, 0644)
}
