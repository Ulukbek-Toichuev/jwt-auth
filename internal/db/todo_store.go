package db

import "database/sql"

type TodoStore struct {
	db *sql.DB
}

func NewTodoStore(db *sql.DB) *TodoStore {
	return &TodoStore{db}
}

func (ts *TodoStore) CreateTodo() {

}

func (ts *TodoStore) GetAll() {

}

func (ts *TodoStore) GetById() {

}

func (ts *TodoStore) UpdateStatus() {

}

func (ts *TodoStore) DeleteTodo() {

}
