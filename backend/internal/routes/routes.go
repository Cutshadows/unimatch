package routes

import (
	"unimatch-back/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	// Define routes for workout-related endpoints
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID) // This should
	r.Put("/workouts/{id}", app.WorkoutHandler.HanldeUpdateWorkoutByID)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)

	return r
}
