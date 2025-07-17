package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"unimatch-back/internal/api"
	"unimatch-back/internal/store"
	"unimatch-back/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDb, err := store.Open()
	// Add enhanced configuration to the connection pool settings with:
	// pgDb.SetMaxOpenConns(25)
	// db.SetMaxOpenConns(), db.SetMaxIdleConns(), and db.SetConnMaxIdleTime()
	if err != nil {
		return nil, err
	}
	logger := log.New(os.Stdout, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)

	// our stores will go here

	err = store.MigrateFs(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	workoutStore := store.NewPostgreWorkoutStore(pgDb)
	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDb,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: OK \n")
}
