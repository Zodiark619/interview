package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Task struct
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    DueDate     time.Time `json:"due_date"`
}

// Pagination struct
type Pagination struct {
    CurrentPage int `json:"current_page"`
    TotalPages  int `json:"total_pages"`
    TotalTasks  int `json:"total_tasks"`
}
type UpdateTaskResponse struct {
    Message string `json:"message"`
    Task    Task   `json:"task"`
}
// Response struct
type GetTasksResponse struct {
    Tasks      []Task     `json:"tasks"`
    Pagination Pagination `json:"pagination"`
}

// In-memory storage
var tasks []Task

func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
    status := r.URL.Query().Get("status")
    search := r.URL.Query().Get("search")
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if page == 0 { page = 1 }
    if limit == 0 { limit = 10 }

    // Build query dynamically
    query := "SELECT id, title, description, status, due_date FROM tasks WHERE 1=1"
    args := []interface{}{}
    if status != "" {
        query += " AND status=$" + strconv.Itoa(len(args)+1)
        args = append(args, status)
    }
    if search != "" {
        query += " AND (LOWER(title) LIKE $" + strconv.Itoa(len(args)+1) +
                 " OR LOWER(description) LIKE $" + strconv.Itoa(len(args)+1) + ")"
        args = append(args, "%"+strings.ToLower(search)+"%")
    }
    query += " ORDER BY id LIMIT $" + strconv.Itoa(len(args)+1) +
             " OFFSET $" + strconv.Itoa(len(args)+2)
    args = append(args, limit, (page-1)*limit)

    rows, err := db.Query(query, args...)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var t Task
        rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.DueDate)
        tasks = append(tasks, t)
    }

    // Count total tasks for pagination
    var total int
    countQuery := "SELECT COUNT(*) FROM tasks WHERE 1=1"
    countArgs := []interface{}{}
    if status != "" {
        countQuery += " AND status=$1"
        countArgs = append(countArgs, status)
    }
    if search != "" {
        if len(countArgs) == 0 {
            countQuery += " AND (LOWER(title) LIKE $1 OR LOWER(description) LIKE $1)"
        } else {
            countQuery += " AND (LOWER(title) LIKE $2 OR LOWER(description) LIKE $2)"
        }
        countArgs = append(countArgs, "%"+strings.ToLower(search)+"%")
    }
    db.QueryRow(countQuery, countArgs...).Scan(&total)

    response := GetTasksResponse{
        Tasks: tasks,
        Pagination: Pagination{
            CurrentPage: page,
            TotalPages:  (total + limit - 1) / limit,
            TotalTasks:  total,
        },
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var t Task
    err = db.QueryRow("SELECT id, title, description, status, due_date FROM tasks WHERE id=$1", id).
        Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.DueDate)
    if err == sql.ErrNoRows {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(t)
}
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var updated Task
    if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec(
        "UPDATE tasks SET title=$1, description=$2, status=$3, due_date=$4 WHERE id=$5",
        updated.Title, updated.Description, updated.Status, updated.DueDate, id,
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    updated.ID = id
    response := UpdateTaskResponse{
        Message: "Task updated successfully",
        Task:    updated,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Delete response struct
type DeleteTaskResponse struct {
    Message string `json:"message"`
}

// DELETE /tasks/:id
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    res, err := db.Exec("DELETE FROM tasks WHERE id=$1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    response := DeleteTaskResponse{Message: "Task deleted successfully"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

type CreateTaskResponse struct {
    Message string `json:"message"`
    Task    Task   `json:"task"`
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := db.QueryRow(
        "INSERT INTO tasks (title, description, status, due_date) VALUES ($1,$2,$3,$4) RETURNING id",
        task.Title, task.Description, task.Status, task.DueDate,
    ).Scan(&task.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := CreateTaskResponse{
        Message: "Task created successfully",
        Task:    task,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
