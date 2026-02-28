package expense

import (
	"context"

	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

//go:generate mockery --name=Repository --output=../mocks --outpkg=mocks
type Repository interface {
	AddExpense(ctx context.Context, expense *model.Expense) (*model.Expense, error)
	TodayExpense(ctx context.Context, userID int64) ([]model.Expense, error)
	WeekExpense(ctx context.Context, userID int64) ([]model.Expense, error)
	MonthExpense(ctx context.Context, userID int64) ([]model.Expense, error)
	StatsExpense(ctx context.Context, userID int64) ([]model.Expense, error)
}
