package handler

import (
	"fmt"
	"net/http"

	"github.com/MrNeocore/tasks-api-server/internal/task/get"
	"github.com/MrNeocore/tasks-api-server/internal/task/store"
	t "github.com/MrNeocore/tasks-api-server/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTask(c *gin.Context) {
	// TODO: Automate with Gin ?
	taskId, parsingError := uuid.Parse(c.Param("id"))

	if parsingError != nil {
		taskIdParseErrorMessage := fmt.Sprintf("Failed to parse taskId: %v", parsingError.Error())
		fmt.Println(taskIdParseErrorMessage)
		c.String(http.StatusBadRequest, taskIdParseErrorMessage)
		return
	}

	task, getTaskError := get.GetTask(c, taskId)

	if getTaskError != nil {
		fmt.Println(getTaskError.Error())
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, task)

}

func GetTasks(c *gin.Context) {
	tasks, getTasksError := get.GetTasks(c)

	if getTasksError != nil {
		fmt.Println(getTasksError.Error())
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, tasks)

}

func PostTask(c *gin.Context) {
	var task t.Task
	bindErr := c.BindJSON(&task)

	if bindErr != nil {
		fmt.Println(bindErr.Error())
		c.String(http.StatusBadRequest, bindErr.Error())
		return
	}

	task.SetInternalFields()

	store.StoreTask(c, task)

	successMessage := fmt.Sprintf("Task %v added\n", task.ID)
	fmt.Println(successMessage)
	c.String(http.StatusCreated, successMessage)
}
