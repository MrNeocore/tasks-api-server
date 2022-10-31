package main

import (
	"os"

	"github.com/MrNeocore/tasks-api-server/internal/server"
	"github.com/MrNeocore/tasks-api-server/internal/util"
)

var HOST = util.GetOrElse(os.LookupEnv, "SERVER_HOST", "0.0.0.0")
var PORT = util.GetOrElse(os.LookupEnv, "SERVER_PORT", "8080")

func main() {
	server.Run(HOST, util.StringToInt(PORT))
}
