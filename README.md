# Go Task Tracker CLI

A simple command-line task tracker built with Go.
This application allows users to manage tasks directly from the terminal and stores all task data locally in a `tasks.json` file.

## Roadmap.sh beginner project
This project was created as a part of Task Tracker beginner project.
Check out the project details [roadmap.sh](https://roadmap.sh/projects/task-tracker)

## Features

* Add tasks
* Update tasks
* Delete tasks
* Mark tasks as:

  * `todo`
  * `in-progress`
  * `done`
* List all tasks
* Filter tasks by status
* Persist tasks locally using JSON storage

---

## TODO:
1. Separate concerns
2. Avoid rewriting tasks.json after every command
3. Error handling
4. Validate status values
5. Return errors instead of printing inside helpers
6. Use constants for statuses
7. Improve command parsing
8. Improve deletion/update behavior
9. Use pointers where appropriate
10. Add unit tests
11. Add Task filtering/sorting
12. Add colored terminal output
13. Replace timestamps as string with time.Time
14. Add atomic file writes
15. Use enums via custom type
16. Add graceful UX improvements
17. Rename misleading function names
18. Avoid repeated task scanning
19. Add persistence abstraction

---

## Installation

Clone the repository:

```bash
git clone https://github.com/tspaja2025/go_tasks_cli.git
cd go_tasks_cli
```

Run the application:

```bash
go run main.go
```

Or build the executable:

```bash
go build -o task-cli
```

Run the executable:

```bash
./task-cli
```

---

## Usage

### Add a task

```bash
./task-cli add "Learn Go" todo
```

### Update a task

```bash
./task-cli update 1 "Learn advanced Go"
```

### Delete a task

```bash
./task-cli delete 1
```

### Batch delete tasks

```bash
./task-cli batch-delete 1-5
```

### List all tasks

```bash
./task-cli list
```

### List tasks by status

```bash
./task-cli list-status todo
./task-cli list-status in-progress
./task-cli list-status done
```

### Mark tasks with status

```bash
./task-cli mark-with-status 1 todo
./task-cli mark-with-status 1 in-progress
./task-cli mark-with-status 1 done
```

### Get task by id

```bash
./task-cli get-task-by-id 1
```

### Clear all tasks

```bash
./task-cli clear-all-tasks
Are you sure? This command will delete ALL tasks. (yes/no)
```

### Help

```bash
./task-cli help
./task-cli -h
./task-cli --help
```

---

## Example Output

```text
ID: 1 | Description: Learn Go | Status: todo
ID: 2 | Description: Build CLI App | Status: done
```

---

## Data Storage

Tasks are stored locally in a `tasks.json` file.

Example:

```json
[
  {
    "id": 1,
    "description": "Learn Go",
    "status": "todo",
    "createdAt": "2026-05-07 10:00:00",
    "updatedAt": "2026-05-07 10:00:00"
  }
]
```

---

## Technologies Used

* Go
* JSON file storage
* Standard library packages:

  * `encoding/json`
  * `fmt`
  * `log`
  * `os`
  * `time`
  * `strconv`

## Learning Goals

This project was built to practice:

* Go structs
* JSON encoding/decoding
* File handling
* CLI argument parsing
* Working with slices
* Basic CRUD operations

---

## License

This project is open source and available under the MIT License.
