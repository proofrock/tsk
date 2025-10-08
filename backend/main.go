package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

// Version is set at build time via ldflags
var Version = "dev"

//go:embed frontend/dist
var frontendFS embed.FS

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id"`
	ParentID    *int   `json:"parent_id"`
	Order       int    `json:"order"`
	Completed   bool   `json:"completed"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id"`
	ParentID    *int   `json:"parent_id"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id"`
	ParentID    *int   `json:"parent_id"`
}

type ReorderTasksRequest struct {
	Tasks []TaskUpdate `json:"tasks"`
}

type TaskUpdate struct {
	ID       int  `json:"id"`
	ParentID *int `json:"parent_id"`
}

var db *sql.DB

func main() {
	dbPath := flag.String("db", "./trx.db", "Database file path")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("tsk version %s\n", Version)
		os.Exit(0)
	}

	var err error
	db, err = sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB()

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/version", getVersion).Methods("GET")
	api.HandleFunc("/categories", getCategories).Methods("GET")
	api.HandleFunc("/tasks", getTasks).Methods("GET")
	api.HandleFunc("/tasks", createTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}/complete", completeTask).Methods("POST")
	api.HandleFunc("/tasks/reorder", reorderTasks).Methods("POST")

	// Serve frontend
	frontendDist, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.FS(frontendDist)))

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}).Handler(r)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func initDB() {
	schema := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		category_id INTEGER NOT NULL,
		parent_id INTEGER,
		task_order INTEGER NOT NULL DEFAULT 0,
		completed BOOLEAN NOT NULL DEFAULT 0,
		FOREIGN KEY (category_id) REFERENCES categories(id),
		FOREIGN KEY (parent_id) REFERENCES tasks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_category ON tasks(category_id);
	CREATE INDEX IF NOT EXISTS idx_tasks_order ON tasks(category_id, task_order);
	CREATE INDEX IF NOT EXISTS idx_tasks_parent ON tasks(parent_id);
	`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	// Add parent_id column if it doesn't exist (migration for existing databases)
	_, err = db.Exec("ALTER TABLE tasks ADD COLUMN parent_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE")
	if err != nil {
		// Column already exists, ignore error
	}

	// Insert default category if not exists
	_, err = db.Exec("INSERT OR IGNORE INTO categories (id, name) VALUES (1, 'General')")
	if err != nil {
		log.Fatal(err)
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"version": Version})
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")

	var rows *sql.Rows
	var err error

	if categoryID != "" {
		rows, err = db.Query(
			"SELECT id, title, description, category_id, parent_id, task_order, completed FROM tasks WHERE category_id = ? AND completed = 0 ORDER BY task_order",
			categoryID,
		)
	} else {
		rows, err = db.Query(
			"SELECT id, title, description, category_id, parent_id, task_order, completed FROM tasks WHERE completed = 0 ORDER BY category_id, task_order",
		)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CategoryID, &t.ParentID, &t.Order, &t.Completed); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get max order for this category and parent
	var maxOrder int
	var err error
	if req.ParentID != nil {
		err = db.QueryRow("SELECT COALESCE(MAX(task_order), -1) FROM tasks WHERE category_id = ? AND parent_id = ? AND completed = 0", req.CategoryID, req.ParentID).Scan(&maxOrder)
	} else {
		err = db.QueryRow("SELECT COALESCE(MAX(task_order), -1) FROM tasks WHERE category_id = ? AND parent_id IS NULL AND completed = 0", req.CategoryID).Scan(&maxOrder)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := db.Exec(
		"INSERT INTO tasks (title, description, category_id, parent_id, task_order) VALUES (?, ?, ?, ?, ?)",
		req.Title, req.Description, req.CategoryID, req.ParentID, maxOrder+1,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	task := Task{
		ID:          int(id),
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		ParentID:    req.ParentID,
		Order:       maxOrder + 1,
		Completed:   false,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start transaction to update task and cascade category changes to subtasks
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Update the task
	_, err = tx.Exec(
		"UPDATE tasks SET title = ?, description = ?, category_id = ?, parent_id = ? WHERE id = ?",
		req.Title, req.Description, req.CategoryID, req.ParentID, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If category changed, update all subtasks to the same category
	_, err = tx.Exec(
		"UPDATE tasks SET category_id = ? WHERE parent_id = ?",
		req.CategoryID, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func completeTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Delete the task when completed
	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func reorderTasks(w http.ResponseWriter, r *http.Request) {
	var req ReorderTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE tasks SET task_order = ?, parent_id = ? WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for i, task := range req.Tasks {
		_, err := stmt.Exec(i, task.ParentID, task.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
