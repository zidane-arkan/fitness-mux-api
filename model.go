// model.go

package main

import (
	"database/sql"
	"errors"
)

type exercise struct {
	ID          int `json:"id"`
	Name        string `json:"name"`
	WorkoutType string `json:"workoutType"`
	Sets        int32  `json:"sets"`
}

// The purpose of this function is to interact with the database, to retrieve an exercise record.
func (e *exercise) getExercise(db *sql.DB) error {
	return errors.New("not implemented")
}

func (e *exercise) updateExercise(db *sql.DB) error {
	return errors.New("not implemented")
}

func (e *exercise) deleteExercise(db *sql.DB) error {
	return errors.New("not implemented")
}

func (e *exercise) createExercise(db *sql.DB) error {
	return errors.New("not implemented")
}

// getExercises is a method of the exercise struct that retrieves a list of exercises from the database.
// It returns a slice of exercise objects and an error if the operation fails.
func (e *exercise) getExercises() ([]exercise, error) {
	return nil, errors.New("not implemented")
}
