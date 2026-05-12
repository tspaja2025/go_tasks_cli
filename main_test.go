package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// Test helper functions
func setupTestTasks() []Task {
	var testTime = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	return []Task{
		{ID: 1, Description: "Test task 1", Status: StatusTodo, CreatedAt: testTime},
		{ID: 2, Description: "Test task 2", Status: StatusInProgress, CreatedAt: testTime},
		{ID: 3, Description: "Test task 3", Status: StatusDone, CreatedAt: testTime},
		{ID: 4, Description: "Test task 4", Status: StatusTodo, CreatedAt: testTime},
	}
}

func setupTempFile(t *testing.T) string {
	tempDir := t.TempDir()
	return filepath.Join(tempDir, "tasks.json")
}

// Test nextID function
// nextID() = highest existing ID + 1
func TestNextID(t *testing.T) {
	testCases := []struct {
		name     string
		tasks    []Task
		expected int
	}{
		{
			name:     "Empty task",
			tasks:    []Task{},
			expected: 1,
		},
		{
			name:     "Single task",
			tasks:    []Task{{ID: 5}},
			expected: 6,
		},
		{
			name:     "Multiple tasks",
			tasks:    []Task{{ID: 1}, {ID: 3}, {ID: 7}, {ID: 2}},
			expected: 8,
		},
		{
			name:     "Non-sequantial IDs",
			tasks:    []Task{{ID: 100}, {ID: 999}},
			expected: 1000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := nextID(tc.tasks)
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}

// Test getID function
func TestGetID(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  int
		expectErr bool
	}{
		{"Valid ID", "42", 42, false},
		{"Zero", "0", 0, false},
		{"Negative", "-5", -5, false},
		{"Invalid string", "abc", 0, true},
		{"Empty string", "", 0, true},
		{"Float", "3.14", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getID(tc.input)
			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if tc.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}

// Test parseIDRange function
func TestParseIDRange(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  []int
		expectErr bool
	}{
		{"Single ID", "5", []int{5}, false},
		{"Range", "1-5", []int{1, 2, 3, 4, 5}, false},
		{"Range with spaces", "1 - 5", []int{1, 2, 3, 4, 5}, false},
		{"Invalid range order", "5-1", nil, true},
		{"Invalid format", "1-5-3", nil, true},
		{"Invalid start", "abc-5", nil, true},
		{"Invalid end", "1-xyz", nil, true},
		{"Invalid single", "abc", nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseIDRange(tc.input)
			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if tc.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// Test isValidStatus function
func TestIsValidStatus(t *testing.T) {
	testCases := []struct {
		status   string
		expected bool
	}{
		{StatusTodo, true},
		{StatusInProgress, true},
		{StatusDone, true},
		{"Invalid", false},
		{"TODO", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			result := isValidStatus(tc.status)
			if result != tc.expected {
				t.Errorf("Expected %v, got: %v", tc.expected, result)
			}
		})
	}
}

// Test updateTask function
func TestUpdateTask(t *testing.T) {
	testCases := []struct {
		name            string
		tasks           []Task
		taskID          int
		newDescription  string
		expectedLength  int
		expectedDesc    string
		expectTimestamp bool
	}{
		{
			name:            "Update existing task",
			tasks:           setupTestTasks(),
			taskID:          2,
			newDescription:  "Updated description",
			expectedLength:  4,
			expectedDesc:    "Updated description",
			expectTimestamp: true,
		},
		{
			name:            "Update non-existent task",
			tasks:           setupTestTasks(),
			taskID:          999,
			newDescription:  "Should not update",
			expectedLength:  4,
			expectedDesc:    "",
			expectTimestamp: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedTasks := updateTask(tc.tasks, tc.taskID, tc.newDescription)

			if len(updatedTasks) != tc.expectedLength {
				t.Errorf("Expected length %d, got %d", tc.expectedLength, len(updatedTasks))
			}

			// Find and verify the task if it should exist
			found := false
			for _, task := range updatedTasks {
				if task.ID == tc.taskID {
					found = true
					if tc.expectedDesc != "" && task.Description != tc.expectedDesc {
						t.Errorf("Expected description '%s', got '%s'", tc.expectedDesc, task.Description)
					}
				}
			}

			if tc.expectedDesc != "" && !found {
				t.Errorf("Task with ID %d not found", tc.taskID)
			}
		})
	}
}

// Test deleteTask function
func TestDeleteTask(t *testing.T) {
	testCases := []struct {
		name           string
		tasks          []Task
		deleteID       int
		expectedLength int
		shouldExist    bool
	}{
		{
			name:           "Delete existing task",
			tasks:          setupTestTasks(),
			deleteID:       2,
			expectedLength: 3,
			shouldExist:    false,
		},
		{
			name:           "Delete non-existent task",
			tasks:          setupTestTasks(),
			deleteID:       999,
			expectedLength: 4,
			shouldExist:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newTasks := deleteTask(tc.tasks, tc.deleteID)

			if len(newTasks) != tc.expectedLength {
				t.Errorf("Expected length %d, got %d", tc.expectedLength, len(newTasks))
			}

			// Verify task are removed
			for _, task := range newTasks {
				if task.ID == tc.deleteID {
					t.Error("Task should have been deelte but still exists")
				}
			}
		})
	}
}

// Test deleteTask function
func TestBatchDeleteTasks(t *testing.T) {
	testCases := []struct {
		name           string
		tasks          []Task
		deleteIDs      []int
		expectedLength int
	}{
		{
			name:           "Delete multiple existing tasks",
			tasks:          setupTestTasks(),
			deleteIDs:      []int{1, 3},
			expectedLength: 2,
		},
		{
			name:           "Delete with non-existent IDs",
			tasks:          setupTestTasks(),
			deleteIDs:      []int{999, 1000},
			expectedLength: 4,
		},
		{
			name:           "Delete mixed valid and invalid IDs",
			tasks:          setupTestTasks(),
			deleteIDs:      []int{1, 999},
			expectedLength: 3,
		},
		{
			name:           "Delete empty slice",
			tasks:          setupTestTasks(),
			deleteIDs:      []int{},
			expectedLength: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newTasks := batchDeleteTasks(tc.tasks, tc.deleteIDs)

			if len(newTasks) != tc.expectedLength {
				t.Errorf("Expected length %d, got %d", tc.expectedLength, len(newTasks))
			}

			// Verify deleted tasks are removed
			for _, task := range newTasks {
				for _, deleteID := range tc.deleteIDs {
					if task.ID == deleteID {
						t.Errorf("Task with ID %d should have been deleted", task.ID)
					}
				}
			}
		})
	}
}

// Test markWithStatus function
func TestMarkWithStatus(t *testing.T) {
	testCases := []struct {
		name           string
		tasks          []Task
		taskID         int
		newStatus      string
		expectedStatus string
		expectError    bool
	}{
		{
			name:           "Mark existing task as done",
			tasks:          setupTestTasks(),
			taskID:         1,
			newStatus:      StatusDone,
			expectedStatus: StatusDone,
			expectError:    false,
		},
		{
			name:           "Mark existing task as in progress",
			tasks:          setupTestTasks(),
			taskID:         3,
			newStatus:      StatusInProgress,
			expectedStatus: StatusInProgress,
			expectError:    false,
		},
		{
			name:           "Mark non-existent task",
			tasks:          setupTestTasks(),
			taskID:         999,
			newStatus:      StatusDone,
			expectedStatus: "",
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedTasks := markWithStatus(tc.tasks, tc.taskID, tc.newStatus)

			if tc.expectedStatus != "" {
				found := false
				// Verify deleted tasks are removed
				for _, task := range updatedTasks {
					if task.ID == tc.taskID {
						found = true
						if task.Status != tc.expectedStatus {
							t.Errorf("Expected status '%s', got '%s'", tc.expectedStatus, task.Status)
						}
					}
				}
				if !found {
					t.Errorf("Task with ID %d not found", tc.taskID)
				}
			}

		})
	}
}

// Test handleAdd function
func TestHandleAdd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		tasks          []Task
		expectedLength int
		expectError    bool
		expectedID     int
		expectedStatus string
	}{
		{
			name:           "Add task with valid arguments",
			args:           []string{"New task", StatusTodo},
			tasks:          setupTestTasks(),
			expectedLength: 5,
			expectError:    false,
			expectedID:     5,
			expectedStatus: StatusTodo,
		},
		{
			name:           "Add task with different status",
			args:           []string{"Important task", StatusInProgress},
			tasks:          setupTestTasks(),
			expectedLength: 5,
			expectError:    false,
			expectedID:     5,
			expectedStatus: StatusInProgress,
		},
		{
			name:           "Add task with unsufficient arguments",
			args:           []string{"Only one argument"},
			tasks:          setupTestTasks(),
			expectedLength: 4,
			expectError:    true,
			expectedID:     0,
			expectedStatus: "",
		},
		{
			name:           "Add task with empty description",
			args:           []string{"", StatusTodo},
			tasks:          setupTestTasks(),
			expectedLength: 5,
			expectError:    false,
			expectedID:     5,
			expectedStatus: StatusTodo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newTasks, err := handleAdd(tc.args, tc.tasks)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(newTasks) != tc.expectedLength {
				t.Errorf("Expected length %d, got %d", tc.expectedLength, len(newTasks))
			}

			if !tc.expectError && tc.expectedID > 0 {
				newTask := newTasks[len(newTasks)-1]
				if newTask.ID != tc.expectedID {
					t.Errorf("Expected ID %d, got %d", tc.expectedID, newTask.ID)
				}
				if newTask.Status != tc.expectedStatus {
					t.Errorf("Expected status %s, got %s", tc.expectedStatus, newTask.Status)
				}
			}
		})
	}
}

// Test handleDelete function
func TestHandleDelete(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		tasks          []Task
		expectedLength int
		expectError    bool
	}{
		{
			name:           "Delete with valid ID",
			args:           []string{"2"},
			tasks:          setupTestTasks(),
			expectedLength: 3,
			expectError:    false,
		},
		{
			name:           "Delete with invalid ID format",
			args:           []string{"invalid"},
			tasks:          setupTestTasks(),
			expectedLength: 4,
			expectError:    true,
		},
		{
			name:           "Delete with non-existent ID",
			args:           []string{"999"},
			tasks:          setupTestTasks(),
			expectedLength: 4,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newTasks, err := handleDelete(tc.args, tc.tasks)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(newTasks) != tc.expectedLength {
				t.Errorf("Expected length %d, got %d", tc.expectedLength, len(newTasks))
			}
		})
	}
}

// Test handleBatchDelete function
func TestHandleBatchDelete(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		tasks          []Task
		expectedLength int
		expectError    bool
	}{
		{
			name:           "Delete valid range",
			args:           []string{"1-2"},
			tasks:          setupTestTasks(),
			expectedLength: 2,
			expectError:    false,
		},
		{
			name:           "Delete invalid range format",
			args:           []string{"invalid"},
			tasks:          setupTestTasks(),
			expectedLength: 4,
			expectError:    true,
		},
		{
			name:           "Delete with missing argument",
			args:           []string{},
			tasks:          setupTestTasks(),
			expectedLength: 4,
			expectError:    true,
		},
		{
			name:           "Delete overlapping range",
			args:           []string{"2-4"},
			tasks:          setupTestTasks(),
			expectedLength: 1,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newTasks, err := handleBatchDelete(tc.args, tc.tasks)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(newTasks) != tc.expectedLength {
				t.Errorf("Expected length %d, got %d", tc.expectedLength, len(newTasks))
			}
		})
	}
}

// Test handleListStatus function
func TestHandleListStatus(t *testing.T) {
	tasks := setupTestTasks()

	// Test valid status
	args := []string{StatusTodo}
	_, err := handleListStatus(args, tasks)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test invalid status
	args = []string{"invalid"}
	_, err = handleListStatus(args, tasks)
	if err == nil {
		t.Error("Expected error for invalid status")
	}

	// Test missing arguments
	args = []string{}
	_, err = handleListStatus(args, tasks)
	if err == nil {
		t.Error("Expected error for missing argument")
	}
}

// Test handleMarkStatus function
func TestHandleMarkStatus(t *testing.T) {
	tasks := setupTestTasks()

	// Test valid marking
	args := []string{"1", StatusDone}
	newTasks, err := handleMarkStatus(args, tasks)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify status changed
	for _, task := range newTasks {
		if task.ID == 1 && task.Status != StatusDone {
			t.Errorf("Expected status '%s', got '%s'", StatusDone, task.Status)
		}
	}

	// Test invalid ID
	args = []string{"invalid", StatusDone}
	_, err = handleListStatus(args, tasks)
	if err == nil {
		t.Error("Expected error for invalid ID")
	}

	// Test missing arguments
	args = []string{"1"}
	_, err = handleListStatus(args, tasks)
	if err == nil {
		t.Error("Expected error for missing argument")
	}
}

// Test handleGetTask function
func TestHandleGetTask(t *testing.T) {
	tasks := setupTestTasks()

	// Test valid marking
	args := []string{"2"}
	_, err := handleMarkStatus(args, tasks)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test non-existent ID
	args = []string{"999"}
	_, err = handleGetTask(args, tasks)
	if err != nil {
		t.Errorf("Unexpected error for non-existent ID: %v", err)
	}

	// Test invalid ID
	args = []string{"invalid"}
	_, err = handleGetTask(args, tasks)
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
}

// Test loadTasks and saveTasks functions
func TestLoadAndSaveTasks(t *testing.T) {
	// Save original file path and restore after test
	originalFile := "tasks.json"
	defer func() {
		os.Remove(originalFile)
	}()

	tasks := setupTestTasks()

	// Test saving
	err := saveTasks(tasks)
	if err != nil {
		t.Errorf("Error saving tasks: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(originalFile); os.IsNotExist(err) {
		t.Error("Tasks file was not created")
	}

	// Test loading
	loadedTasks, err := loadTasks()
	if err != nil {
		t.Errorf("Error loading tasks: %v", err)
	}
	if len(loadedTasks) != len(tasks) {
		t.Errorf("Expected %d tasks, got %d", len(tasks), len(loadedTasks))
	}

	// Test loading from non-exitent file
	os.Remove(originalFile)
	emptyTasks, err := loadTasks()
	if err != nil {
		t.Errorf("Error loading from non-existent file: %v", err)
	}
	if len(emptyTasks) != 0 {
		t.Error("Expected empty task list for non-existent file")
	}
}

// Test handleClearAll function
func TestHandleClearAll(t *testing.T) {
	tasks := setupTestTasks()

	result := clearAllTasks(tasks)
	if result == nil {
		t.Error("clearAllTasks should return a task slice")
	}
}

// Test executeCommand function
func TestExecuteCommand(t *testing.T) {
	tasks := setupTestTasks()

	testCases := []struct {
		name      string
		command   string
		args      []string
		expectErr bool
	}{
		{"List command", "list", []string{}, false},
		{"Add command", "add", []string{"New Task", StatusTodo}, false},
		{"Update command", "update", []string{"1", "Updated"}, false},
		{"Delete command", "delete", []string{"1"}, false},
		{"Help command", "help", []string{}, false},
		{"Invalid command", "invalid", []string{}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := executeCommand(tc.command, tc.args, tasks)
			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
