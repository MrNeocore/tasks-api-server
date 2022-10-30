package main

import (
	"net/http"

	"github.com/MrNeocore/tasks-api-server/internal/api/v1/handler"
)

func main() {
	http.HandleFunc("/tasks", handler.Tasks)
	http.HandleFunc("/tasks/", handler.Tasks)

	http.ListenAndServe("0.0.0.0:8080", nil)
}
