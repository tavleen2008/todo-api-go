
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)
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

