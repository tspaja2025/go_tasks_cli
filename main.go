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
	idStr, _ := strconv.Atoi(id)
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
		fmt.Printf("ID: %d | Description: %s | Status: %s\n", task.Id, task.Description, task.Status)
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
		}

		status := os.Args[2]

		if status == "todo" || status == "in-progress" || status == "done" {
			listTasksByStatus(tasks, status)
		} else {
			fmt.Println("Invalid status")
		}
	default:
		fmt.Println("Unknown command:")
	}

	byteValue, _ := json.MarshalIndent(tasks, "", "  ")
	err = os.WriteFile("tasks.json", byteValue, 0644)
}

// The application should run from the command line, - done
// accept user actions and inputs as arguments, - done
// and store the tasks in a JSON file. - done
// User should be able to:
// Add, Update, and Delete tasks - done
// Mark a task as in progress or done
// List all tasks - done
// List all tasks that have status "todo"
// List all tasks that have status "in-progress"
// List all tasks that have status "done"

// Task properties:
// Each task should have the following properties
// id: A unique identifier for the task - done
// description: A short description of the task - done
// status: The status of the task (todo,in-progress,done) - done
// createdAt: The date and time when the task was created - done
// updatedAt: The date and time when the task last updated - done

// Constraints to guide the implementation:
// Use positional arguments in command line to accept user inputs
// Use a JSON file to store the tasks in the current directory
// The JSON file should be created if it does not exist
// Use the native file system module of the programming language to interact with the JSOn file
// Do not use any external libraries or frameworks for this project
// Ensure to handle errors and edge cases gracefully
