package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unimatch-back/internal/store"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkOutID := chi.URLParam(r, "id")
	if paramsWorkOutID == "" {
		http.NotFound(w, r)
		return
	}
	workout, err := strconv.ParseInt(paramsWorkOutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)

	fmt.Fprintf(w, "Get workout by ID: %d", workout)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createWorkout)
}

func (wh *WorkoutHandler) HanldeUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkOutID := chi.URLParam(r, "id")
	if paramsWorkOutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkOutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	workoutExistance, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if workoutExistance == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title          *string              `json:"title"`
		Description    *string              `json:"description"`
		Duration       *int                 `json:"duration_minutes"`
		CaloriesBurned *int                 `json:"calories_burned"`
		Entries        []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateWorkoutRequest.Title != nil {
		workoutExistance.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		workoutExistance.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.Duration != nil {
		workoutExistance.Duration = *updateWorkoutRequest.Duration
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		workoutExistance.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		workoutExistance.Entries = updateWorkoutRequest.Entries
	}

	err = wh.workoutStore.UpdateWorkout(workoutExistance)
	if err != nil {
		fmt.Println("update workout error", err)
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workoutExistance)
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkOutID := chi.URLParam(r, "id")
	if paramsWorkOutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkOutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = wh.workoutStore.DeleteWorkoutByID(workoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error deleting workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
