package main

import (
	"fmt"
	"os"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TaskManager struct {
	tasks    []Task
	filename string
	nextId   int
}

func NewTaskManager(filename string) *TaskManager {
	// TODO:
}

func (tm *TaskManager) loadTasks() {
	// TODO:
}

func (tm *TaskManager) saveTasks() error {
	// TODO:
}

func (tm *TaskManager) Add(description string) {
	// TODO:
	task := Task{
		Id:          tm.nextId,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tm.tasks = append(tm.tasks, task)
}

func (tm *TaskManager) List(status string) {
	// TODO:
}

func (tm *TaskManager) Update(id int, description string) {
	// TODO:
}

func (tm *TaskManager) Delete(id int) {
	// TODO:
}

func (tm *TaskManager) MarkStatus(id int, status string) {
	// TODO:
}

func truncate(s string, maxLen int) string {
	// TODO:
}

func printUsage() {
	// TODO:
	// Example usage:
	// add "Buy groceries"
	// list
	// mark-in-progress 1
	// mark-done 1
	// delete 1
}

func main() {
	// TODO:
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		// TODO:
	case "list":
		// TODO:
	case "update":
	// TODO:
	case "delete":
	// TODO:
	case "mark-todo":
	// TODO:
	case "mark-in-progress":
	// TODO:
	case "mark-done":
	// TODO:
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
	}
}

// The application should run from the command line,
// accept user actions and inputs as arguments,
// and store the tasks in a JSON file.
// User should be able to:
// Add, Update, and Delete tasks
// Mark a task as in progress or done
// List all tasks
// List all tasks that are not done
// List all tasks that are in progress

// Task properties:
// Each task should have the following properties
// id: A unique identifier for the task
// description: A short description of the task
// status: The status of the task (todo,in-progress,done)
// createdAt: The date and time when the task was created
// updatedAt: The date and time when the task last updated

// Constraints to guide the implementation:
// Use positional arguments in command line to accept user inputs
// Use a JSON file to store the tasks in the current directory
// The JSON file should be created if it does not exist
// Use the native file system module of the programming language to interact with the JSOn file
// Do not use any external libraries or frameworks for this project
// Ensure to handle errors and edge cases gracefully
