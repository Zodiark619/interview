menjalankan program dengan 'go run .' di command line

postman=
GET = http://localhost:8080/tasks , http://localhost:8080/tasks/5
POST = http://localhost:8080/tasks
        body = {
    "title": "example",
    "description": "example",
    "status": "pending",
    "due_date": "2026-06-20T00:00:00Z"
    }
UPDATE = http://localhost:8080/tasks/5
        body = {
    "title": "example",
    "description": "example",
    "status": "pending",
    "due_date": "2026-06-20T00:00:00Z"
    } 
DELETE = http://localhost:8080/tasks/5

//postgres table
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('pending','completed')),
    due_date DATE
);


