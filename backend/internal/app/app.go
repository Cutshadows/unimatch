package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
)


type Application struct {
	Logger *log.Logger
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "APP: ", log.Ldate| log.Ltime | log.Lshortfile)

	app := &Application{
		Logger: logger,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: OK \n")
}