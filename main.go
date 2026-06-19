package main

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)
var db *sql.DB
func initDB() {
    var err error
    connStr := "user=postgres password=postgres dbname=Interview sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    if err = db.Ping(); err != nil {
        panic(err)
    }
}
func main() {
	  initDB()
    // Register endpoints
    http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            createTaskHandler(w, r)
        } else if r.Method == http.MethodGet {
            getAllTasksHandler(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            getTaskByIDHandler(w, r)   // SELECT by id
        } else if r.Method == http.MethodPut {
            updateTaskHandler(w, r)    // UPDATE
        } else if r.Method == http.MethodDelete {
            deleteTaskHandler(w, r)    // DELETE
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    http.ListenAndServe(":8080", nil)
}
//post http://localhost:8080/tasks
// {
//   "title": "My First Task",
//   "description": "Try out Postman with Go",
//   "status": "pending",
//   "due_date": "2026-06-20T00:00:00Z"
// }