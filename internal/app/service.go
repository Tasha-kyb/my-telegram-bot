package app

import (
	"context"

	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/category"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/expense"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/user"
	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

type Service struct {
	UserSvc     *user.Service
	CategorySvc *category.Service
	ExpenseSvc  *expense.Service
}

func New(userSvc *user.Service, categorySvc *category.Service, expenseSvc *expense.Service) *Service {
	return &Service{
		UserSvc:     userSvc,
		CategorySvc: categorySvc,
		ExpenseSvc:  expenseSvc,
	}
}

func (s *Service) CreateProfile(ctx context.Context, req model.Profile) (string, error) {
	return s.UserSvc.CreateProfile(ctx, req)
}

func (s *Service) AddCategory(ctx context.Context, req model.Category) (string, error) {
	return s.CategorySvc.AddCategory(ctx, req)
}

func (s *Service) GetAllCategories(ctx context.Context, userID int64) (string, error) {
	return s.CategorySvc.GetAllCategories(ctx, userID)
}

func (s *Service) DeleteCategory(ctx context.Context, userID int64, id int) (string, error) {
	return s.CategorySvc.DeleteCategory(ctx, userID, id)
}

func (s *Service) AddExpense(ctx context.Context, expense *model.Expense) (string, error) {
	return s.ExpenseSvc.AddExpense(ctx, expense)
}

func (s *Service) TodayExpense(ctx context.Context, userID int64) (string, error) {
	return s.ExpenseSvc.TodayExpense(ctx, userID)
}

func (s *Service) WeekExpense(ctx context.Context, userID int64) (string, error) {
	return s.ExpenseSvc.WeekExpense(ctx, userID)
}

func (s *Service) MonthExpense(ctx context.Context, userID int64) (string, error) {
	return s.ExpenseSvc.MonthExpense(ctx, userID)
}

func (s *Service) StatsExpense(ctx context.Context, userID int64) (string, error) {
	return s.ExpenseSvc.StatsExpense(ctx, userID)
}
