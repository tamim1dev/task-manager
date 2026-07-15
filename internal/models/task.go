package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Due_Date    time.Time `json:"due_date"`
	User_Id     uuid.UUID `json:"user_id"`
	Created_At  time.Time `json:"created_at"`
}

type AddTask struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Due_Date    time.Time `json:"due_date"`
	User_Id     uuid.UUID `json:"user_id"`
}

type AddTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Days        int    `json:"days"`
}

type UpdateTask struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
	Due_Date    *int    `json:"due_date,omitempty"`
}
