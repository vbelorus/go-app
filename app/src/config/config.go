package config

import (
	"fmt"
	"github.com/vbelorus/go-app/v2/src/clickhouse"
	"github.com/vbelorus/go-app/v2/src/models"
	"log"
	"net/http"
	"runtime/debug"
)

type Application struct {
	InfoLog            *log.Logger
	ErrorLog           *log.Logger
	AppDB              *clickhouse.AppDB
	DeviceEventChannel chan models.DeviceEvent
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	app.LogError(err)
	http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
}

func (app *Application) LogError(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
}

func (app *Application) LogInfo(message string) {
	app.InfoLog.Output(2, fmt.Sprintf("%s", message))
}

func (app *Application) ClientError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound, "Not Found")
}
