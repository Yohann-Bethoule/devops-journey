package models

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

type Todo struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
	Done  bool   `json:"isDone"`
}

func GetTodos(count int) ([]Todo, error) {

	rows, err := DB.Query("SELECT id, label, done FROM task LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := make([]Todo, 0)

	for rows.Next() {
		todo := Todo{}
		err = rows.Scan(&todo.Id, &todo.Label, &todo.Done)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return todos, err
}

func GetTodoById(id string) (Todo, error) {
	stmt, err := DB.Prepare("SELECT id, label, done FROM task WHERE id = ?")

	if err != nil {
		return Todo{}, err
	}

	todo := Todo{}

	sqlErr := stmt.QueryRow(id).Scan(&todo.Id, &todo.Label, &todo.Done)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Todo{}, nil
		}
		return Todo{}, sqlErr
	}

	return todo, nil
}

func InsertTodo(todo Todo) (bool, error) {
	trans, err := DB.Begin()

	stmt, err := trans.Prepare("INSERT INTO task (label, done) VALUES (?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(todo.Label, todo.Done)

	if err != nil {
		return false, err
	}

	trans.Commit()

	return true, nil
}

func UpdateTodo(todo Todo) (bool, error) {
	trans, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := trans.Prepare("UPDATE task SET label = ?, done = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(todo.Label, todo.Done, todo.Id)

	if err != nil {
		return false, err
	}

	trans.Commit()

	return true, nil
}

func DeleteTodo(id int) (bool, error) {
	trans, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := trans.Prepare("DELETE FROM task WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return false, err
	}

	trans.Commit()

	return true, nil
}
