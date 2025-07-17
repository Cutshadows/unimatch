-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workouts (
	id BIGSERIAL PRIMARY KEY,
	-- user_id
	title VARCHAR(255) NOT NULL,
	description TEXT,
	duration INTEGER NOT NULL, -- Duration in minutes
	calories_burned INTEGER,
	date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,	 
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workouts;
-- +goose StatementEnd