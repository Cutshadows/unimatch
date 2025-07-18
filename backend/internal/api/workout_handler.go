package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"unimatch-back/internal/store"
	"unimatch-back/util"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := util.ReadIOParam(r)
	if err != nil {
		wh.logger.Printf("Error: readParams: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid workout id"})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("Error: getWorkoutByID: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusOK, util.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodingCreateWorkout %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request sent"})
		return
	}

	createWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: creatingWorkout %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request payload"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"workout": createWorkout})
}

func (wh *WorkoutHandler) HanldeUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkOutID, err := util.ReadIOParam(r)

	if err != nil {
		wh.logger.Printf("Error: readParams: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid workout id"})
		return
	}

	workoutExistance, err := wh.workoutStore.GetWorkoutByID(paramsWorkOutID)
	if err != nil {
		wh.logger.Printf("Error: getWorkoutByID: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "internal server error"})
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
		wh.logger.Printf("Error: decodingUpdateWorkout %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request sent"})
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
		wh.logger.Printf("error: updatingWorkout %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "failed to update workout"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{"workout": workoutExistance})

}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkOutID, err := util.ReadIOParam(r)
	if err != nil {
		wh.logger.Printf("Error: readParams: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid workout id"})
		return
	}

	err = wh.workoutStore.DeleteWorkoutByID(paramsWorkOutID)
	if err == sql.ErrNoRows {
		wh.logger.Printf("error: workout not found for ID %d", paramsWorkOutID)
		util.WriteJSON(w, http.StatusNotFound, util.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
