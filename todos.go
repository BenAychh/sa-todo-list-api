package main

import (
	"database/sql"
	"errors"
)

type todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
}

func (t *todo) create(db *sql.DB) error {
	return db.QueryRow(
		`insert into todos (description) VALUES ($1) returning id, description, complete`,
		t.Description,
	).Scan(&t.ID, &t.Description, &t.Complete)
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
	results, err := db.Query("select * from todos")
	var todos []todo = nil
	defer results.Close()
	if err == nil {
		todos = []todo{}
		for results.Next() {
			var t todo
			scanError := results.Scan(&t.ID, &t.Description, &t.Complete)
			if scanError != nil {
				return nil, scanError
			}
			todos = append(todos, t)
		}
	}
	return todos, err
}
