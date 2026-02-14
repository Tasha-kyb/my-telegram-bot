package model

import (
	"time"
)

type Profile struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Created_at time.Time `json:"-"`
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
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
}

type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Massage  struct {
		From struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
		} `json:"from"`
		Text string `json:"text"`
	} `json:"massage"`
}
