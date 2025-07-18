package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5433 user=centra-db-admin password=12345 dbname=unimatch_test_db sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE workouts, workout_entries RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("truncate table failed: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgreWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "Valid Workout",
			workout: &Workout{
				Title:          "Morning Run",
				Description:    "5km run in the park",
				Duration:       30, // in minutes
				CaloriesBurned: 300,
				Entries: []WorkoutEntry{
					{
						ExerciseName:    "Running",
						Sets:            1,
						Reps:            IntPtr(10),
						DurationSeconds: nil,
						Weight:          FloatPtr(0.0), // No weight for running
						Notes:           "Felt great",
						OrderIndex:      1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Workout (missing title)",
			workout: &Workout{
				Description:    "No title provided",
				Duration:       20,
				CaloriesBurned: 200,
				Entries: []WorkoutEntry{
					{
						ExerciseName:    "Cycling",
						Sets:            1,
						Reps:            IntPtr(5),
						DurationSeconds: nil,
						Weight:          FloatPtr(50.0),
						Notes:           "Good pace",
						OrderIndex:      1,
					},
					{
						ExerciseName:    "Hiking",
						Sets:            4,
						Reps:            IntPtr(2),
						DurationSeconds: IntPtr(1200),
						Weight:          FloatPtr(50.0),
						Notes:           "Challenging trail",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err, "Expected error but got none")
				return
			}
			require.NoError(t, err, "Unexpected error while creating workout")
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.Duration, createdWorkout.Duration)

			retrieve, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrieve.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieve.Entries))

			for i := range retrieve.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieve.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieve.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].Reps, retrieve.Entries[i].Reps)
				assert.Equal(t, tt.workout.Entries[i].DurationSeconds, retrieve.Entries[i].DurationSeconds)
				assert.Equal(t, tt.workout.Entries[i].Weight, retrieve.Entries[i].Weight)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
