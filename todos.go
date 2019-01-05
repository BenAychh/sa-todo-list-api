package main

import (
	"database/sql"
)

type todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
}

func (t *todo) get(db *sql.DB) error {
	return db.QueryRow(
		`select * from todos where id = $1`,
		t.ID,
	).Scan(&t.ID, &t.Description, &t.Complete)
}

func (t *todo) create(db *sql.DB) error {
	return db.QueryRow(
		`insert into todos (description) VALUES ($1) returning id, description, complete`,
		t.Description,
	).Scan(&t.ID, &t.Description, &t.Complete)
}

func (t *todo) update(db *sql.DB) error {
	return db.QueryRow(
		`update todos set description=$1, complete=$2 where id=$3 returning id, description, complete`,
		t.Description, t.Complete, t.ID,
	).Scan(&t.ID, &t.Description, &t.Complete)
}

func (t *todo) delete(db *sql.DB) error {
	_, err := db.Exec(
		`delete from todos where id=$1`,
		t.ID,
	)
	return err
}

func getTodos(db *sql.DB) ([]todo, error) {
	results, err := db.Query("select * from todos")
	var todos []todo
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
