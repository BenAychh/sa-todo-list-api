package main

import (
	"database/sql"
	"errors"
)

type todo struct {
	id          int
	description string
	complete    bool
}

func (t *todo) create(db *sql.DB) error {
	return errors.New("Not Implimented Yet")
}

func (t *todo) toggleCompletion(db *sql.DB) error {
	return errors.New("Not Implimented Yet")
}

func (t *todo) updateDescription(db *sql.DB) error {
	return errors.New("Not Implimented Yet")
}

func (t *todo) delete(db *sql.DB) error {
	return errors.New("Not Implimented Yet")
}

func getTodos(db *sql.DB) ([]todo, error) {
	return nil, errors.New("Not Implimented Yet")
}
