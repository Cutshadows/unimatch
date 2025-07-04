package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"unimatch-back/internal/api"
)


type Application struct {
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "APP: ", log.Ldate| log.Ltime | log.Lshortfile)


	// our stores will go here


	// our handlers will go here

	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger: logger,
		WorkoutHandler: workoutHandler,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: OK \n")
}