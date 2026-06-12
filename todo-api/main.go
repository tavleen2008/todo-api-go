package main

import (
	"fmt"
	"net/http"
)

func main() {

	connectDB()

	http.HandleFunc("/todo", todoHandler)
	http.HandleFunc("/todo/", todoHandler)

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}