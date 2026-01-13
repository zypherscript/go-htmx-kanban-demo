# Mini Kanban Board 📋

A simple, local Kanban board web application built with Go and htmx for managing your personal tasks.

## Features

- **Three Columns**: To Do, Doing, Done
- **Add Tasks**: Create tasks with title and description
- **Move Tasks**: Seamlessly move tasks between columns with buttons
- **No Page Reloads**: Uses htmx for dynamic updates
- **Beautiful UI**: Modern, gradient design with smooth animations
- **Offline-First**: JSON file persistence - your tasks survive restarts!
- **Thread-Safe**: Concurrent access protection with mutex

## Tech Stack

- **Backend**: Go with `net/http` and `html/template`
- **Frontend**: HTML, CSS, and htmx (loaded via CDN)
- **Storage**: JSON file (`tasks.json`) - no database needed!

## Project Structure

```
go-htmx-demo/
├── main.go                        # Go server and handlers
├── go.mod                         # Go module file
├── tasks.json                     # Your tasks (auto-created)
├── .gitignore                     # Git ignore file
├── templates/
│   ├── index.html                 # Main page template
│   ├── all-columns.html           # All three columns template
│   └── column-content.html        # Single column content template
└── README.md                      # This file
```

## Quick Start

### 1. Run the Application

```bash
go run main.go
```

### 2. Open in Browser

Navigate to: http://localhost:8080

### 3. Use the Board

- **Add a Task**: Fill in the form at the top and click "Add Task"
- **Move Tasks**: Click the buttons on each task card to move them between columns
  - "Move to Doing →" - Moves from To Do to Doing
  - "Move to Done ✓" - Moves from Doing to Done
  - "← Back to..." - Moves tasks backwards

## How It Works

### htmx Integration

The application uses htmx attributes for dynamic interactions:

- `hx-post`: Makes POST requests without page reload
- `hx-target`: Specifies where to insert the response
- `hx-swap`: Defines how to swap the content (innerHTML)
- `hx-vals`: Sends additional parameters with requests

### Go Handlers

- **`/`**: Serves the main page with all tasks
- **`/add-task`**: Handles task creation (POST)
- **`/move-task`**: Handles moving tasks between columns (POST)
- **`/column/{status}`**: Returns content for a specific column

### Data Storage

Tasks are persisted to a JSON file automatically:
- **Default location**: `./tasks.json` in project directory
- **Custom location**: Set `KANBAN_DATA_FILE` environment variable
- **Auto-save**: Every add/move operation saves to disk
- **Auto-load**: Tasks reload when you restart the server
- **Human-readable**: JSON format you can view/edit directly
- **Thread-safe**: Mutex protection for concurrent operations
- **Offline-first**: Works completely locally, no internet needed

#### Custom Data Location

**Use cloud folder for sync across devices:**
```bash
export KANBAN_DATA_FILE=~/Dropbox/kanban-tasks.json
go run main.go
```

**Use different locations for different projects:**
```bash
export KANBAN_DATA_FILE=~/work-tasks.json
go run main.go
```

**Permanent setup** (add to `~/.zshrc` or `~/.bashrc`):
```bash
export KANBAN_DATA_FILE="$HOME/Dropbox/kanban-tasks.json"
```

Your tasks survive server restarts and are portable with your project!

## Example Usage

### Adding Tasks

```
Title: Build authentication system
Description: Implement JWT-based auth with user registration and login
```

### Workflow

1. Create tasks in "To Do"
2. Move active tasks to "Doing"
3. Complete tasks and move to "Done"
4. Move tasks back if needed

## Customization

### Change Data File Location

Set environment variable:
```bash
export KANBAN_DATA_FILE=/path/to/your/tasks.json
```

### Change Port

Edit `main.go`:
```go
log.Fatal(http.ListenAndServe(":8080", nil))
```

### Modify Styling

Edit the `<style>` section in `templates/index.html`:
- Colors: Change gradient, button colors
- Layout: Adjust column widths, spacing
- Fonts: Change font family

### Add More Columns

1. Add new status in `Task` struct
2. Create handlers for the new column
3. Update templates to include the new column

## Development

### Requirements

- Go 1.21 or higher
- Modern web browser (Chrome, Firefox, Safari, Edge)

### No Dependencies

The project uses only Go standard library. htmx is loaded from CDN in the HTML template.

## Keyboard Shortcuts

Since this is a web app, you can use browser shortcuts:
- `Cmd+R` / `F5`: Refresh the page
- `Cmd+T`: Open in new tab

## Tips

- **Mobile Friendly**: The board adapts to smaller screens with a single-column layout
- **Fast**: In-memory storage means instant updates
- **Simple**: No build step, just run and go
- **Visual Feedback**: Cards animate on hover for better UX

## License

Free to use for personal projects. Modify as needed!

---

Happy task managing! 🎯
