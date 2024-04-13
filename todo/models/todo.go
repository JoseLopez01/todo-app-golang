package models

import "time"

type Todo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	StartDate   time.Time `json:"start_date"`
	Completed   bool      `json:"completed"`
}
