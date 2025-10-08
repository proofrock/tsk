# tsk

A minimal, aesthetically pleasing task manager/todo list web application with subtask support.

## Features

- ✅ **Create and manage tasks** - Add tasks with title, description, and category
- 🌳 **Hierarchical subtasks** - Create one level of subtasks under parent tasks
- ✏️ **Edit tasks** - Modify any task or subtask details
- ✓ **Delete tasks** - Select and delete multiple tasks at once
- 🔄 **Drag-and-drop reordering** - Reorder tasks and subtasks with mouse or touch
- 📂 **Categories** - Organize tasks by category with configurable default
- 🎯 **Task badges** - Parent tasks show the count of their subtasks
- 👁️ **Collapse/expand** - Toggle subtask visibility individually or all at once
- 🎨 **Modern dark theme** - Clean UI with orange accent color (#f97316)
- 📱 **Fully responsive** - Mobile-first design with Bootstrap 5
- 🚀 **Single binary deployment** - Embedded frontend and SQLite database

## Tech Stack

- **Backend**: Go 1.21+ with embedded SQLite database (CGO enabled)
- **Frontend**: Svelte 5 + Vite with Bootstrap 5
- **Deployment**: Docker & Docker Compose
- **Build System**: Makefile with automated version management

## Quick Start with Docker

The easiest way to run tsk is using Docker:

```bash
docker run -d -p 8080:8080 --user 1000:1000 -v tsk-data:/db ghcr.io/proofrock/tsk:latest
```

The application will be available at `http://localhost:8080`

### Docker Images

Images are automatically built and pushed to GitHub Container Registry on each tagged release:

- `ghcr.io/proofrock/tsk:latest` - Latest release
- `ghcr.io/proofrock/tsk:v1.0.0` - Specific version

## Development Setup

### Prerequisites

- Go 1.21+ (with CGO enabled)
- Node.js 20+
- npm
- make (optional, for build automation)

### Backend Setup

```bash
cd backend
go mod download
```

### Frontend Setup

```bash
cd frontend
npm install
```

### Running in Development

1. Start the backend:
```bash
cd backend
go run main.go
```

2. In a separate terminal, start the frontend dev server:
```bash
cd frontend
npm run dev
```

The frontend will be available at `http://localhost:5173` with API proxy to the backend.

## Building for Production

Use the Makefile:

```bash
make build
# or simply
make
```

This will:
1. Build the frontend and create optimized static files
2. Build the Go backend with embedded frontend
3. Create a single `./tsk` binary

Run the binary:
```bash
./tsk
```

Check version:
```bash
./tsk --version
```

### Building with Version

The version is automatically detected from git tags. To build with a specific version:
```bash
make build VERSION=v1.0.0
```

Or create a git tag:
```bash
git tag v1.0.0
make build
```

### Makefile Targets

- `make build` - Build frontend and backend (default)
- `make frontend` - Build only frontend
- `make backend` - Build only backend (requires frontend to be built first)
- `make clean` - Remove all build artifacts (node_modules, dist, binaries, etc.)
- `make help` - Show available targets

## Managing Categories

### Adding Categories via UI

Categories must be added directly to the database. Once added, they appear in the category dropdown.

### Adding Categories via Database

1. Connect to the database:
```bash
sqlite3 tsk.db
```

2. Insert a new category:
```sql
INSERT INTO categories (name, is_default) VALUES ('Work', 0);
INSERT INTO categories (name, is_default) VALUES ('Personal', 0);
INSERT INTO categories (name, is_default) VALUES ('Shopping', 0);
```

3. Set a category as default (optional):
```sql
UPDATE categories SET is_default = 0;  -- Unset all defaults first
UPDATE categories SET is_default = 1 WHERE name = 'Work';  -- Set Work as default
```

4. Restart the application to see the new categories.

## Database Schema

### Categories Table
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `name`: TEXT NOT NULL UNIQUE
- `is_default`: BOOLEAN NOT NULL DEFAULT 0

### Tasks Table
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `title`: TEXT NOT NULL
- `description`: TEXT
- `category_id`: INTEGER (foreign key to categories)
- `parent_id`: INTEGER (nullable, foreign key to tasks with CASCADE delete)
- `order`: INTEGER (for drag-drop ordering)
- `completed`: BOOLEAN

### Task Hierarchy Rules

- **One level only**: Subtasks cannot have their own subtasks
- **Category inheritance**: Subtasks inherit their parent's category
- **Cascade delete**: Deleting a parent deletes all its subtasks
- **Cascade category changes**: Changing a parent's category updates all subtasks
- **No circular references**: Tasks with subtasks cannot become subtasks

## Database Migrations

Migration scripts are provided for upgrading existing databases:

### migrate.sql
Adds subtask support (parent_id column):
```bash
sqlite3 tsk.db < migrate.sql
```

### migrate2.sql
Adds default category support (is_default column):
```bash
sqlite3 tsk.db < migrate2.sql
```

## API Endpoints

### Categories
- `GET /api/categories` - List all categories with default flag

### Tasks
- `GET /api/tasks?category_id={id}` - List tasks by category (includes subtasks)
- `POST /api/tasks` - Create a new task
  ```json
  {
    "title": "Task title",
    "description": "Optional description",
    "category_id": 1,
    "parent_id": null  // or parent task ID for subtasks
  }
  ```
- `PUT /api/tasks/{id}` - Update a task (cascades category to subtasks)
  ```json
  {
    "title": "Updated title",
    "description": "Updated description",
    "category_id": 1,
    "parent_id": null
  }
  ```
- `POST /api/tasks/{id}/complete` - Mark task as complete (deletes it)
- `POST /api/tasks/reorder` - Reorder tasks and update parent relationships
  ```json
  {
    "tasks": [
      {"id": 1, "parent_id": null},
      {"id": 2, "parent_id": 1},
      {"id": 3, "parent_id": null}
    ]
  }
  ```

### Version
- `GET /api/version` - Get application version

## User Interface Guide

### Task Management

- **Add Task**: Click the "Add Task" button in the header
- **Edit Task**: Click the pencil icon on the right side of any task
- **Delete Tasks**: Check the boxes of tasks to delete, then click "Delete X tasks"
- **Reorder Tasks**: Drag and drop tasks with mouse or touch

### Subtasks

- **Create Subtask**: In the task modal, select a parent from the "Parent Task" dropdown
- **Make Child**: Drag a task to the right edge (1/6 width) of another task
- **Promote to Task**: Drag a subtask to the left area of any non-subtask
- **Reorder within Subtasks**: Drag subtasks within the right edge zone (green indicator)
- **Collapse/Expand**: Click the chevron button on parent tasks, or use toolbar buttons

### Visual Indicators

- **Orange line**: Normal task reordering or new child creation
- **Green line**: Reordering within subtask group
- **Badge**: Number on parent tasks shows subtask count
- **Indentation**: Subtasks are indented 3rem from left

### Drag-and-Drop Zones

- **Left 5/6 of task**: Normal reordering (before/after)
- **Right 1/6 of parent task**: Make dragged task a child
- **Right 1/6 of subtask**: Reorder as sibling within same parent

## Docker Deployment

### Build the image:
```bash
docker build -t tsk .
```

### Run with docker-compose:
```bash
docker-compose up -d
```

The application data (SQLite database) is persisted in a Docker volume named `tsk-data`.

### Docker Volume Management

To backup the database:
```bash
docker run --rm -v tsk-data:/db -v $(pwd):/backup alpine cp /db/tsk.db /backup/tsk-backup.db
```

To restore from backup:
```bash
docker run --rm -v tsk-data:/db -v $(pwd):/backup alpine cp /backup/tsk-backup.db /db/tsk.db
```

## Configuration

### Port Configuration

The application runs on port 8080 by default. To change this, modify the `main.go` file:

```go
log.Fatal(http.ListenAndServe(":8080", handler))
```

### Database Location

By default, the database is stored in `./tsk.db`. In Docker, this is mounted to `/db/tsk.db`.

## Releases

To create a new release:

1. Tag the commit with a version following semver format:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Build the frontend
   - Build the Docker image for linux/amd64
   - Push the image to `ghcr.io/proofrock/tsk:v1.0.0`
   - Update the `latest` tag

The version number will be embedded in both the binary and the UI footer.

## Troubleshooting

### Build Issues

- **CGO errors**: Ensure CGO is enabled (`CGO_ENABLED=1`) and build tools are installed
- **SQLite errors**: Install SQLite development libraries (`libsqlite3-dev` on Debian/Ubuntu)
- **Frontend build fails**: Clear `node_modules` and run `npm install` again

### Runtime Issues

- **Database locked**: Stop all instances of the application before running a new one
- **Port already in use**: Change the port in `main.go` or stop the conflicting service
- **Drag-and-drop not working on mobile**: Ensure you're not touching interactive elements (checkboxes, buttons)

## License

This project is licensed under the European Union Public Licence (EUPL) v. 1.2.

See the [LICENSE](LICENSE) file for the full license text.

The EUPL is a copyleft open-source license compatible with GPL, AGPL, MPL, LGPL, and other major open-source licenses. You can find more information about the EUPL at https://eupl.eu/
