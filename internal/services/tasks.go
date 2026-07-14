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

func GetTasksByUserID(userId string, r *http.Request) ([]models.Task, error) {
	query := `SELECT * FROM tasks WHERE user_id = $1`
	tasksRows, dbErr := database.DB.Pool.Query(r.Context(), query, userId)
	if dbErr != nil {
		return []models.Task{}, dbErr
	}
	defer tasksRows.Close()

	var tasks []models.Task
	for tasksRows.Next() {
		var task models.Task
		loopErr := tasksRows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.Due_Date,
			&task.User_Id,
			&task.Created_At,
		)
		if loopErr != nil {
			return []models.Task{}, loopErr
		}
		tasks = append(tasks, task)
	}

	if rowsErr := tasksRows.Err(); rowsErr != nil {
		return []models.Task{}, rowsErr
	}

	return tasks, nil
}
