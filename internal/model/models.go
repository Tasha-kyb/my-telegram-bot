package model

import (
	"time"
)

type Profile struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"-"`
}

type Category struct {
	UserID int64  `json:"user_id"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color,omitempty"`
}

type CategoryResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

type Expense struct {
	UserID      int64     `json:"user_id"`
	CategoryID  int       `json:"category_id"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
}
