package router

import (
	"github.com/vbelorus/go-app/v2/src/config"
	"github.com/vbelorus/go-app/v2/src/handler"
	"net/http"
)

func SetRoutes(app *config.Application) *http.ServeMux {
	mux := http.NewServeMux()

	//mux.HandleFunc("/", handler.Home(app))                  // Method Get
	mux.Handle("/events", &handler.DeviceEventHandler{App: app}) // Method Post

	return mux
}
