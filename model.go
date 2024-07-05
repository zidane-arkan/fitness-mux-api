// model.go

package main

import (
	"database/sql"
)

type exercise struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	WorkoutType string `json:"workoutType"`
	Sets        int32  `json:"sets"`
}

// The purpose of this function is to interact with the database, to retrieve an exercise record.
func (e *exercise) getExercise(db *sql.DB) error {
	// return errors.New("not implemented")
	return db.QueryRow("SELECT name, workoutType, sets FROM exercises WHERE id=$1",
		e.ID).Scan(&e.Name, &e.WorkoutType, &e.Sets)
}

func (e *exercise) updateExercise(db *sql.DB) error {
	_, err := db.Exec("UPDATE exercises SET name=$1, workoutType=$2, sets=$3 WHERE id=$4",
		e.Name, e.WorkoutType, e.Sets, e.ID)
	return err
}

func (e *exercise) deleteExercise(db *sql.DB) error {
	// return errors.New("not implemented")
	_, err := db.Exec("DELETE FROM exercises WHERE id=$1", e.ID)
	return err
}

func (e *exercise) createExercise(db *sql.DB) error {
	// return errors.New("not implemented")
	err := db.QueryRow(
		"INSERT INTO exercises(name, workoutType, sets) VALUES($1, $2, $3) RETURNING id",
		e.Name, e.WorkoutType, e.Sets).Scan(&e.ID)

	if err != nil {
		return err
	}

	return nil
}

// getExercises is a method of the exercise struct that retrieves a list of exercises from the database.
// It returns a slice of exercise objects and an error if the operation fails.
func getExercises(db *sql.DB, start, count int) ([]exercise, error) {
	rows, err := db.Query(
		"SELECT id, name, workoutType, sets FROM exercises LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	exercises := []exercise{}

	for rows.Next() {
		var p exercise
		if err := rows.Scan(&p.ID, &p.Name, &p.WorkoutType, &p.Sets); err != nil {
			return nil, err
		}
		exercises = append(exercises, p)
	}

	return exercises, nil
}
