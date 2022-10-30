package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MrNeocore/tasks-api-server/internal/api/v1/handler"
	"github.com/MrNeocore/tasks-api-server/internal/util"
)

var HOST = util.GetOrElse(os.LookupEnv, "SERVER_HOST", "0.0.0.0")
var PORT = util.GetOrElse(os.LookupEnv, "SERVER_PORT", "8080")

func main() {
	http.HandleFunc("/tasks", handler.Tasks)
	http.HandleFunc("/tasks/", handler.Tasks)

	listenOn := fmt.Sprintf("%v:%v", HOST, PORT)
	http.ListenAndServe(listenOn, nil)
}
