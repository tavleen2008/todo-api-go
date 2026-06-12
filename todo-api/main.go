package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Completed bool   `json:"completed"`
}
var db *sql.DB

var todos = []Todo{
	{
		ID:    1,
		Title: "Learn Go",
	},
	{
		ID:    2,
		Title: "Learn Backend",
	},
}

// GET 
func getTodos(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo

		err := rows.Scan(&t.ID, &t.Title, &t.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		todos = append(todos, t)
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(todos)
}

func getTodoByID(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var todo Todo

	err = db.QueryRow(
		"SELECT id, title, completed FROM todos WHERE id = $1",
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Completed)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// POST 
func createTodo(w http.ResponseWriter, r *http.Request) {

	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.QueryRow(
		"INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id",
		todo.Title,
		todo.Completed,
	).Scan(&todo.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// PUT 
func updateTodo(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updatedTodo Todo

	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.QueryRow(
		"UPDATE todos SET title = $1, completed = $2 WHERE id = $3 RETURNING id, title, completed",
		updatedTodo.Title,
		updatedTodo.Completed,
		id,
	).Scan(
		&updatedTodo.ID,
		&updatedTodo.Title,
		&updatedTodo.Completed,
	)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

// DELETE 
func deleteTodo(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"DELETE FROM todos WHERE id = $1",
		id,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func todoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet &&
		r.URL.Path == "/todo" {

		getTodos(w, r)
		return
	}

	if r.Method == http.MethodPost &&
		r.URL.Path == "/todo" {

		createTodo(w, r)
		return
	}

	if r.Method == http.MethodGet &&
		strings.HasPrefix(r.URL.Path, "/todo/") {

		getTodoByID(w, r)
		return
	}

	if r.Method == http.MethodPut &&
		strings.HasPrefix(r.URL.Path, "/todo/") {

		updateTodo(w, r)
		return
	}

	if r.Method == http.MethodDelete &&
		strings.HasPrefix(r.URL.Path, "/todo/") {

		deleteTodo(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
func connectDB() {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL!")
}

func main() {
	connectDB()

	http.HandleFunc("/todo", todoHandler)

	http.HandleFunc("/todo/", todoHandler)

	http.ListenAndServe(":8080", nil)
}