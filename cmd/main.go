package main

import (
	"net/http"

	"golang-rest-api-server/internal/api/v1/handler"
)

func main() {
	http.HandleFunc("/task", handler.Task)
	http.HandleFunc("/task/", handler.Task)

	http.ListenAndServe("localhost:8080", nil)
}
