package usecase

import (
	"context"

	"github.com/internal/model"
)

type Repository interface {
	CreateProfile(ctx context.Context, req *model.Profile) error
	AddCategory(ctx context.Context, req *model.Category) (int, error)
	GetAllCategories(ctx context.Context, userID int64) ([]model.Category, error)
	DeleteCategory(ctx context.Context, userID int64, id int) (string, error)
	AddExpense(ctx context.Context, expence *model.Expense) (*model.Expense, error)
	TodayExpense(ctx context.Context, userID int64) ([]model.Expense, error)
	WeekExpense(ctx context.Context, userID int64) ([]model.Expense, error)
	MonthExpense(ctx context.Context, userID int64) ([]model.Expense, error)
	StatsExpense(ctx context.Context, userID int64) ([]model.Expense, error)
}
