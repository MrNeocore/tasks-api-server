package server

import (
	"fmt"
	"net/http"

	"github.com/MrNeocore/tasks-api-server/internal/server/v1/handler"
)

func Run(host string, port int) {
	http.HandleFunc("/tasks", handler.Tasks)
	http.HandleFunc("/tasks/", handler.Tasks)

	listenOn := fmt.Sprintf("%v:%v", host, port)
	http.ListenAndServe(listenOn, nil)

}

// func main() {
// 	r := gin.Default()
// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "pong",
// 		})
// 	})
// 	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
// }
