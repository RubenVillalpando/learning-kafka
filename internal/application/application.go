package application

import (
	"net/http"
)

type App struct{}

func New() *App {
	return &App{}
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
