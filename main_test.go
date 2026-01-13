package main

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStore() *TaskStore {
	tmpFile := filepath.Join(os.TempDir(), "kanban_test_tasks.json")
	_ = os.Remove(tmpFile)
	return &TaskStore{
		tasks:    make(map[int]*Task),
		nextID:   1,
		filePath: tmpFile,
	}
}

func TestAddTask(t *testing.T) {
	store := newTestStore()
	task := store.AddTask("Test Task", "Test Description")
	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("Title mismatch")
	}
	if task.Status != "todo" {
		t.Errorf("Expected status 'todo', got %s", task.Status)
	}
}

func TestGetTasksByStatus(t *testing.T) {
	store := newTestStore()
	store.AddTask("A", "")
	store.AddTask("B", "")
	store.MoveTask(1, "doing")
	todo := store.GetTasksByStatus("todo")
	doing := store.GetTasksByStatus("doing")
	if len(todo) != 1 || todo[0].ID != 2 {
		t.Errorf("Expected one todo task with ID 2")
	}
	if len(doing) != 1 || doing[0].ID != 1 {
		t.Errorf("Expected one doing task with ID 1")
	}
}

func TestMoveTask(t *testing.T) {
	store := newTestStore()
	task := store.AddTask("Move Me", "")
	_, ok := store.MoveTask(task.ID, "doing")
	if !ok {
		t.Errorf("MoveTask failed")
	}
	if store.tasks[task.ID].Status != "doing" {
		t.Errorf("Task status not updated")
	}
	_, ok = store.MoveTask(999, "done")
	if ok {
		t.Errorf("Should not move non-existent task")
	}
}

func TestPersistence(t *testing.T) {
	store := newTestStore()
	store.AddTask("Persist", "Test")
	store.AddTask("Persist2", "Test2")
	store.MoveTask(1, "done")
	store.saveToFile()

	newStore := &TaskStore{
		tasks:    make(map[int]*Task),
		nextID:   1,
		filePath: store.filePath,
	}
	if err := newStore.LoadFromFile(); err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	if len(newStore.tasks) != 2 {
		t.Errorf("Expected 2 tasks after load")
	}
	if newStore.tasks[1].Status != "done" {
		t.Errorf("Status not persisted")
	}
}

func TestEmptyStore(t *testing.T) {
	store := newTestStore()
	if len(store.GetTasksByStatus("todo")) != 0 {
		t.Errorf("Expected no tasks in empty store")
	}
}
