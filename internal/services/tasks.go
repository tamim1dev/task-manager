package services

import (
	"net/http"

	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/models"
)

func AddTask(task *models.AddTask, r *http.Request) (models.Task, error) {
	var newTask models.Task
	query := `INSERT INTO tasks (title, description, due_date, user_id) VALUES ($1, $2, $3, $4) RETURNING *`
	dbErr := database.DB.Pool.QueryRow(r.Context(), query, task.Title, task.Description, task.Due_Date, task.User_Id).Scan(
		&newTask.Id,
		&newTask.Title,
		&newTask.Description,
		&newTask.Completed,
		&newTask.Due_Date,
		&newTask.User_Id,
		&newTask.Created_At,
	)
	if dbErr != nil {
		return models.Task{}, dbErr
	}
	return newTask, nil
}
