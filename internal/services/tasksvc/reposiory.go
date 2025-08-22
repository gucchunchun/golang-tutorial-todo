package tasksvc

import (
	"context"
	"golang/tutorial/todo/internal/models"
)

type Repository interface {
	Create(ctx context.Context, name string, dueAt *models.Date) (models.Task, error)
	GetByID(ctx context.Context, id models.TaskID) (models.Task, error)
	List(ctx context.Context, limit, offset int) (models.Tasks, error)
	Update(ctx context.Context, id models.TaskID, upd models.TaskUpdate) (models.Task, error)
	Delete(ctx context.Context, id models.TaskID) error
}
