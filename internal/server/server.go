package server

import (
	"fmt"

	"github.com/MrNeocore/tasks-api-server/internal/handler"
	"github.com/gin-gonic/gin"
)

func Run(host string, port int) {
	r := gin.Default()
	r.GET("/tasks/:id", handler.GetTask)
	r.GET("/tasks", handler.GetTasks)
	r.POST("/tasks", handler.PostTask)

	listenOn := fmt.Sprintf("%v:%v", host, port)
	r.Run(listenOn)
}
