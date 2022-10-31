package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/MrNeocore/tasks-api-server/internal/api/v1/handler/create"
	"github.com/MrNeocore/tasks-api-server/internal/api/v1/handler/get"
	"github.com/google/uuid"
)

func Tasks(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		create.CreateTask(w, req)
	case http.MethodGet:
		_taskId, hasValidTaskId, taskIdParsingError := extractValidTaskIdFromPath(req.URL)
		if hasValidTaskId {
			if taskIdParsingError != nil {
				handleTaskIdParsingError(w, taskIdParsingError)
				return
			}
			get.GetTask(*_taskId, w, req)
		} else {
			get.GetTasks(w, req)
		}
	default:
		errorMessage := "Unsupported verb for route."
		fmt.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
	}
}

func extractValidTaskIdFromPath(url *url.URL) (*uuid.UUID, bool, error) {
	pathParts := strings.Split(url.Path, "/")
	if len(pathParts) != 3 {
		return nil, false, nil
	} else {
		_taskId := pathParts[2]
		taskId, taskIdParseError := uuid.Parse(_taskId)
		if taskIdParseError != nil {
			return &taskId, true, taskIdParseError
		} else {
			return &taskId, true, nil
		}
	}
}

func handleTaskIdParsingError(w http.ResponseWriter, parsingError error) {
	taskIdParseErrorMessage := fmt.Sprintf("Failed to parse taskId: %v", parsingError)
	fmt.Println(taskIdParseErrorMessage)
	http.Error(w, taskIdParseErrorMessage, http.StatusBadRequest)
}
