package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"unimatch-back/internal/api"
	"unimatch-back/internal/middleware"
	"unimatch-back/internal/store"
	"unimatch-back/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     *middleware.Middleware
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

	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	workoutStore := store.NewPostgresWorkoutStore(pgDb)
	userStore := store.NewPostgresUserStore(pgDb)
	tokenStore := store.NewPostgresTokenStore(pgDb)
	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := &middleware.Middleware{UserStore: userStore}
	// our middleware will go here
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,
		DB:             pgDb,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: OK \n")
}
