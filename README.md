# tsk

A minimal, aesthetically pleasing task manager/todo list web application.

## Features

- ✅ Create tasks with title, description, and category
- ✏️ Edit existing tasks
- ✓ Complete tasks (automatically deleted when checked)
- 🔄 Reorder tasks via drag and drop
- 🏷️ Filter tasks by category
- 🎨 Modern dark theme UI

## Tech Stack

- **Backend**: Go with embedded SQLite database
- **Frontend**: Svelte 5 + Vite
- **Deployment**: Docker & Docker Compose

## Quick Start with Docker

The easiest way to run tsk is using Docker Compose:

```bash
docker-compose up -d
```

The application will be available at `http://localhost:8080`

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 20+
- npm

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

### Makefile Targets

- `make build` - Build frontend and backend (default)
- `make frontend` - Build only frontend
- `make backend` - Build only backend (requires frontend to be built first)
- `make clean` - Remove all build artifacts (node_modules, dist, binaries, etc.)
- `make help` - Show available targets

## Managing Categories

Categories are stored in the SQLite database. To add new categories:

1. Connect to the database:
```bash
sqlite3 tsk.db
```

2. Insert a new category:
```sql
INSERT INTO categories (name) VALUES ('Work');
INSERT INTO categories (name) VALUES ('Personal');
INSERT INTO categories (name) VALUES ('Shopping');
```

3. Restart the application to see the new categories.

## Database Schema

### Categories Table
- `id`: INTEGER PRIMARY KEY
- `name`: TEXT (unique)

### Tasks Table
- `id`: INTEGER PRIMARY KEY
- `title`: TEXT
- `description`: TEXT
- `category_id`: INTEGER (foreign key)
- `task_order`: INTEGER (for drag-drop ordering)
- `completed`: BOOLEAN

## API Endpoints

- `GET /api/categories` - List all categories
- `GET /api/tasks?category_id={id}` - List tasks by category
- `POST /api/tasks` - Create a new task
- `PUT /api/tasks/{id}` - Update a task
- `POST /api/tasks/{id}/complete` - Mark task as complete (deletes it)
- `POST /api/tasks/reorder` - Reorder tasks

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

## Configuration

The application runs on port 8080 by default. To change this, modify the `main.go` file:

```go
log.Fatal(http.ListenAndServe(":8080", handler))
```

## License

MIT
