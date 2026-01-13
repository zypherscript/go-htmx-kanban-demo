package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Task represents a single task in the kanban board
type Task struct {
	ID          int
	Title       string
	Description string
	Status      string // "todo", "doing", "done"
}

// TaskStore holds all tasks with thread-safe access
type TaskStore struct {
	mu       sync.Mutex
	tasks    map[int]*Task
	nextID   int
	filePath string
}

// getDataFilePath returns the data file path from env var or default
func getDataFilePath() string {
	// Check environment variable first
	dataFile := os.Getenv("KANBAN_DATA_FILE")
	if dataFile != "" {
		return dataFile
	}
	// Fallback to project directory
	return filepath.Join(".", "tasks.json")
}

var store = &TaskStore{
	tasks:    make(map[int]*Task),
	nextID:   1,
	filePath: getDataFilePath(),
}

// AddTask adds a new task to the store
func (s *TaskStore) AddTask(title, description string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &Task{
		ID:          s.nextID,
		Title:       title,
		Description: description,
		Status:      "todo",
	}
	s.tasks[task.ID] = task
	s.nextID++
	s.saveToFile()
	return task
}

// GetTask retrieves a task by ID
func (s *TaskStore) GetTask(id int) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	return task, ok
}

// GetTasksByStatus returns all tasks with a specific status
func (s *TaskStore) GetTasksByStatus(status string) []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []*Task
	for _, task := range s.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// MoveTask changes the status of a task
func (s *TaskStore) MoveTask(id int, newStatus string) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, false
	}
	task.Status = newStatus
	s.saveToFile()
	return task, true
}

// Persistence structures
type PersistentData struct {
	Tasks  []*Task `json:"tasks"`
	NextID int     `json:"next_id"`
}

// saveToFile saves tasks to JSON file (must be called with lock held)
func (s *TaskStore) saveToFile() {
	var taskList []*Task
	for _, task := range s.tasks {
		taskList = append(taskList, task)
	}

	data := PersistentData{
		Tasks:  taskList,
		NextID: s.nextID,
	}

	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Error creating directory: %v", err)
			return
		}
	}

	file, err := os.Create(s.filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Printf("Error encoding data: %v", err)
	}
}

// LoadFromFile loads tasks from JSON file
func (s *TaskStore) LoadFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No existing data file found, starting fresh")
			return nil
		}
		return err
	}
	defer file.Close()

	var data PersistentData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	s.tasks = make(map[int]*Task)
	for _, task := range data.Tasks {
		s.tasks[task.ID] = task
	}
	s.nextID = data.NextID

	log.Printf("Loaded %d tasks from file", len(s.tasks))
	return nil
}

// Template data structures
type PageData struct {
	TodoTasks  []*Task
	DoingTasks []*Task
	DoneTasks  []*Task
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	// Load existing data from file
	if err := store.LoadFromFile(); err != nil {
		log.Printf("Warning: Could not load data: %v", err)
	}

	// Serve static files (for htmx)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add-task", addTaskHandler)
	http.HandleFunc("/move-task", moveTaskHandler)
	http.HandleFunc("/column/", columnHandler)

	log.Println("Starting server on http://localhost:8080")
	log.Printf("Your tasks are saved to: %s\n", store.filePath)
	if os.Getenv("KANBAN_DATA_FILE") != "" {
		log.Println("Using custom data location from KANBAN_DATA_FILE environment variable")
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// indexHandler serves the main page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		TodoTasks:  store.GetTasksByStatus("todo"),
		DoingTasks: store.GetTasksByStatus("doing"),
		DoneTasks:  store.GetTasksByStatus("done"),
	}
	templates.ExecuteTemplate(w, "index.html", data)
}

// addTaskHandler handles adding a new task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	store.AddTask(title, description)

	// Return the updated "To Do" column
	tasks := store.GetTasksByStatus("todo")
	templates.ExecuteTemplate(w, "column-content.html", map[string]interface{}{
		"Status": "todo",
		"Tasks":  tasks,
	})
}

// moveTaskHandler handles moving tasks between columns
func moveTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	newStatus := r.FormValue("status")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, ok := store.MoveTask(id, newStatus)
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Return all three columns to update the board
	data := PageData{
		TodoTasks:  store.GetTasksByStatus("todo"),
		DoingTasks: store.GetTasksByStatus("doing"),
		DoneTasks:  store.GetTasksByStatus("done"),
	}
	templates.ExecuteTemplate(w, "all-columns.html", data)

	fmt.Printf("Moved task %d (%s) to %s\n", task.ID, task.Title, task.Status)
}

// columnHandler returns a single column's content
func columnHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Path[len("/column/"):]
	if status != "todo" && status != "doing" && status != "done" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	tasks := store.GetTasksByStatus(status)
	templates.ExecuteTemplate(w, "column-content.html", map[string]interface{}{
		"Status": status,
		"Tasks":  tasks,
	})
}
