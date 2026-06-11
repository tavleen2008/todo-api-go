package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

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

	for _, todo := range todos {

		if todo.ID == id {

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(todo)

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

// POST 
func createTodo(w http.ResponseWriter, r *http.Request) {

	var todo Todo

	json.NewDecoder(r.Body).Decode(&todo)

	todo.ID = len(todos) + 1

	todos = append(todos, todo)

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

	json.NewDecoder(r.Body).Decode(&updatedTodo)

	for i, todo := range todos {

		if todo.ID == id {

			todos[i].Title = updatedTodo.Title

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(todos[i])

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
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

	for i, todo := range todos {

		if todo.ID == id {

			todos = append(
				todos[:i],
				todos[i+1:]...,
			)

			w.WriteHeader(http.StatusNoContent)

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
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

func main() {

	http.HandleFunc("/todo", todoHandler)

	http.HandleFunc("/todo/", todoHandler)

	http.ListenAndServe(":8080", nil)
}