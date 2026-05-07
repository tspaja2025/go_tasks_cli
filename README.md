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
