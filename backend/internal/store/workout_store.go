package store

import (
	"database/sql"
	"time"
)

type Workout struct {
	ID             int            `json:"id"`
	UserID         int            `json:"user_id"` // ID of the user who created the workout
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Duration       int            `json:"duration"` // Duration in seconds
	CaloriesBurned int            `json:"calories_burned"`
	Date           *time.Time     `json:"date"` // Date of the workout
	CreatedAt      *time.Time     `json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
	Entries        []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"` // Optional for exercises that are timed
	Weight          *float64 `json:"weight"`           // Optional for weight-based exercises
	Notes           string   `json:"notes"`            // Optional notes for the entry
	OrderIndex      int      `json:"order_index"`      // Order of the entry in the workout
}

type PostgreWorkoutStore struct {
	db *sql.DB
}

func NewPostgreWorkoutStore(db *sql.DB) *PostgreWorkoutStore {
	return &PostgreWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkoutByID(id int64) error
	// GetWorkoutOwner(id int64) (int64, error)
}

func (pg *PostgreWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	// Implementation for creating a workout in the database
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	query := `
	INSERT INTO workouts (title, description, duration, calories_burned)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.Duration, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    	RETURNING id
		`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return workout, nil
}
func (pg *PostgreWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
	UPDATE workouts
	SET title = $1, description = $2, duration = $3, calories_burned = $4
	WHERE id = $5`
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.Duration, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id =$1`, workout.ID)
	if err != nil {
		return err
	}
	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := tx.Exec(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (pg *PostgreWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `
	SELECT id, title, description, duration, calories_burned, date
	FROM workouts
	WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.Duration, &workout.CaloriesBurned, &workout.Date)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entryQuery := `
	SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
	FROM workout_entries
	WHERE workout_id = $1
	ORDER BY order_index
	`
	rows, err := pg.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *PostgreWorkoutStore) DeleteWorkoutByID(id int64) error {
	query := `
	DELETE from workouts
	WHERE id = $1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return nil
	}
	rwsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rwsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
