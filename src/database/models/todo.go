package models

import (
	"context"
	"database/sql"
	"time"
)

type Todo struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
}

type TodoModel struct {
	Database *sql.DB
}

func (m *TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todo (task, completed)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.Database.QueryRowContext(
		ctx,
		query,
		todo.Task,
		todo.Completed,
	).Scan(&todo.ID, &todo.CreatedAt)
}

func (m *TodoModel) GetAll() ([]*Todo, error) {
	query := `
		SELECT id, created_at, task, completed
		FROM todo
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []*Todo{}

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.CreatedAt, &todo.Task, &todo.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
