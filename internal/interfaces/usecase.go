package interfaces

import (
	"context"

	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

type UserUseCase interface {
	CreateProfile(ctx context.Context, req model.Profile) (string, error)
}

type CategoryUseCase interface {
	AddCategory(ctx context.Context, req model.Category) (string, error)
	GetAllCategories(ctx context.Context, userID int64) (string, error)
	DeleteCategory(ctx context.Context, userID int64, id int) (string, error)
}

type ExpenseUseCase interface {
	AddExpense(ctx context.Context, expense *model.Expense) (string, error)
	TodayExpense(ctx context.Context, userID int64) (string, error)
	WeekExpense(ctx context.Context, userID int64) (string, error)
	MonthExpense(ctx context.Context, userID int64) (string, error)
	StatsExpense(ctx context.Context, userID int64) (string, error)
}
