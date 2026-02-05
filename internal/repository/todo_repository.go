package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ShamirMaharjan/Go-Todo-App/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var query string = `
		INSERT INTO todo (title, completed)
		VALUES($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`
	var todo models.Todo
	var err error = pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func GetAllTodos(pool *pgxpool.Pool) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var query string = `
	SELECT id, title, completed, created_at, updated_at
	FROM todo
	ORDER BY created_at DESC
	`

	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var Todos []models.Todo = []models.Todo{}

	for rows.Next() {
		var Todo models.Todo

		err = rows.Scan(
			&Todo.ID,
			&Todo.Title,
			&Todo.Completed,
			&Todo.CreatedAt,
			&Todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		Todos = append(Todos, Todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Todos, nil

}

func GetTodoByID(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var query string = `
	SELECT id, title, completed, created_at, updated_at
	FROM todo
	WHERE id = $1
	`

	var todo models.Todo
	var err error
	err = pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodoByID(pool *pgxpool.Pool, id int, title string, completed bool) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var query = `
		UPDATE todo
		SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, title, completed, created_at, updated_at
	`

	var todo models.Todo
	var err error
	err = pool.QueryRow(ctx, query, title, completed, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodoByID(pool *pgxpool.Pool, id int) error {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var query string = `
	DELETE from todo
	WHERE id = $1
	`

	commandTag, err := pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("No todo found of id %d", id)
	}

	return nil
}
