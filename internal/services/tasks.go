package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

func GetTasksByUserID(
	userId string, limit, offset int,
	completed *bool, sortBy, sortOrder string, searchKey *string,
	r *http.Request) ([]models.Task, error) {
	query := `
		SELECT * FROM tasks
		WHERE user_id = $1 
		AND ($2::boolean IS NULL OR completed = $2)
		AND ($3::text IS NULL OR title ILIKE '%' || $3 || '%' OR description ILIKE '%' || $3 || '%')
		ORDER BY ` + sortBy + ` ` + sortOrder + ` LIMIT $4 OFFSET $5`

	tasksRows, dbErr := database.DB.Pool.Query(r.Context(), query, userId, completed, searchKey, limit, offset)
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

func GetTaskById(task_id, user_id string, r *http.Request) (models.Task, error) {
	var task models.Task
	query := `SELECT * FROM tasks WHERE id = $1 AND user_id = $2`
	dbErr := database.DB.Pool.QueryRow(r.Context(), query, task_id, user_id).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.Due_Date,
		&task.User_Id,
		&task.Created_At,
	)
	if dbErr != nil {
		return models.Task{}, dbErr
	}
	return task, nil
}

func EditTask(edits models.UpdateTask, taskId, userId string, r *http.Request) (models.Task, error) {
	var (
		editedTask models.Task
		argPos     = 1
		args       []any
		setClauses []string
	)

	if edits.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, *edits.Title)
		argPos++
	}
	if edits.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *edits.Description)
		argPos++
	}
	if edits.Completed != nil {
		setClauses = append(setClauses, fmt.Sprintf("completed = $%d", argPos))
		args = append(args, *edits.Completed)
		argPos++
	}
	if edits.Due_Date != nil {
		newDate := time.Now().Add(time.Duration(*edits.Due_Date) * time.Hour * 24)
		setClauses = append(setClauses, fmt.Sprintf("due_date = $%d", argPos))
		args = append(args, newDate)
		argPos++
	}

	if len(setClauses) == 0 {
		return models.Task{}, fmt.Errorf("Nothing to update")
	}

	query := fmt.Sprintf(
		"UPDATE tasks SET %s WHERE id = $%d AND user_id = $%d RETURNING *",
		strings.Join(setClauses, ", "),
		argPos,
		argPos+1,
	)

	args = append(args, taskId, userId)

	dbErr := database.DB.Pool.QueryRow(r.Context(), query, args...).Scan(
		&editedTask.Id,
		&editedTask.Title,
		&editedTask.Description,
		&editedTask.Completed,
		&editedTask.Due_Date,
		&editedTask.User_Id,
		&editedTask.Created_At,
	)
	if dbErr != nil {
		return models.Task{}, dbErr
	}

	return editedTask, nil
}

func DeleteTaskById(task_id, user_id string, r *http.Request) (string, error) {
	var deletedId uuid.UUID
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2 RETURNING id`
	dbErr := database.DB.Pool.QueryRow(r.Context(), query, task_id, user_id).Scan(&deletedId)
	if dbErr != nil {
		return "", dbErr
	}

	return deletedId.String(), nil
}
