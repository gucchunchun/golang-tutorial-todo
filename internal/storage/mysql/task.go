package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo { return &TaskRepo{db: db} }

func rowToTask(r *sql.Row) (models.Task, error) {
	var t models.Task
	if err := r.Scan(
		&t.ID, &t.Name, &t.Status, &t.DueAt, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, apperr.E(apperr.CodeNotFound, "task not found", err)
		}
		return t, err
	}
	return t, nil
}

func scanTasks(rows *sql.Rows) (models.Tasks, error) {
	defer rows.Close()
	var out models.Tasks
	for rows.Next() {
		var t models.Task
		var due *models.Date
		if err := rows.Scan(&t.ID, &t.Name, &t.Status, &due, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		t.DueAt = due
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TaskRepo) Create(ctx context.Context, name string, dueAt *models.Date) (models.Task, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (name, status, due_at) VALUES (?, ?, ?)`,
		name, models.StatusPending, dueAt,
	)
	if err != nil {
		return models.Task{}, err
	}
	id, _ := res.LastInsertId()
	return r.GetByID(ctx, models.TaskID(id))
}

func (r *TaskRepo) GetByID(ctx context.Context, id models.TaskID) (models.Task, error) {
	return rowToTask(r.db.QueryRowContext(ctx,
		`SELECT id, name, status, due_at, created_at, updated_at FROM tasks WHERE id = ?`, id))
}

func (r *TaskRepo) List(ctx context.Context, limit, offset int) (models.Tasks, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, status, due_at, created_at, updated_at
		   FROM tasks
		   ORDER BY id DESC
		   LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	return scanTasks(rows)
}

func (r *TaskRepo) Update(ctx context.Context, id models.TaskID, upd models.TaskUpdate) (models.Task, error) {
	q := "UPDATE tasks SET "
	args := []any{}
	sep := ""

	if upd.Name != nil {
		q += sep + "name = ?"
		args = append(args, *upd.Name)
		sep = ", "
	}
	if upd.Status != nil {
		q += sep + "status = ?"
		args = append(args, *upd.Status)
		sep = ", "
	}
	if upd.ClearDue {
		q += sep + "due_at = NULL"
		sep = ", "
	} else if upd.DueAt != nil {
		if upd.DueAt == nil {
			q += sep + "due_at = NULL"
		} else {
			q += sep + "due_at = ?"
			args = append(args, *upd.DueAt)
		}
		sep = ", "
	}

	if sep == "" {
		return r.GetByID(ctx, id) // nothing to update
	}

	q += " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return models.Task{}, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return models.Task{}, apperr.E(apperr.CodeNotFound, "task not found", fmt.Errorf("id=%d", id))
	}
	return r.GetByID(ctx, id)
}

func (r *TaskRepo) Delete(ctx context.Context, id models.TaskID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return apperr.E(apperr.CodeNotFound, "task not found", sql.ErrNoRows)
	}
	return nil
}
