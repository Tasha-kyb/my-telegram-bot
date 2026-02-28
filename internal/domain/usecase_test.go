package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/category"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/expense"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/mocks"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/user"
	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfile(t *testing.T) {
	tests := []struct {
		name        string
		input       model.Profile
		setupMock   func(repo *mocks.UserRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:  "Успешное создание профиля",
			input: model.Profile{ID: 123456, Username: "user"},
			setupMock: func(repo *mocks.UserRepository) {
				repo.On("CreateProfile",
					mock.Anything,
					mock.AnythingOfType("*model.Profile"),
				).Return(nil)
			},
			wantError:   false,
			wantMessage: "👋 Добро пожаловать",
		},
		{
			name:      "Ошибка: ID = 0",
			input:     model.Profile{ID: 0, Username: "user"},
			setupMock: func(repo *mocks.UserRepository) {},
			wantError: true,
		},
		{
			name:      "Ошибка: пустое имя",
			input:     model.Profile{ID: 123456, Username: ""},
			setupMock: func(repo *mocks.UserRepository) {},
			wantError: true,
		},
		{
			name:  "Ошибка в репозитории",
			input: model.Profile{ID: 123456, Username: "user"},
			setupMock: func(repo *mocks.UserRepository) {
				repo.On("CreateProfile",
					mock.Anything,
					mock.AnythingOfType("*model.Profile"),
				).Return(errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMock := mocks.NewUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(userMock)
			}
			service := user.NewService(userMock)
			message, err := service.CreateProfile(context.Background(), tt.input)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном создании профиля, но его нет")
			}
			userMock.AssertExpectations(t)
		})
	}

	t.Log("Тест завершен")
}
func TestAddCategory(t *testing.T) {
	tests := []struct {
		name        string
		input       model.Category
		setupMock   func(repo *mocks.CategoryRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:  "Успешное создание категории",
			input: model.Category{ID: 123456, Name: "Спорт"},
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("AddCategory",
					mock.Anything,
					mock.AnythingOfType("*model.Category"),
				).Return(123, nil)
			},
			wantError:   false,
			wantMessage: "✅ Категория создана!",
		},
		{
			name:      "Ошибка: нет названия категории",
			input:     model.Category{ID: 123456, Name: ""},
			setupMock: func(repo *mocks.CategoryRepository) {},
			wantError: true,
		},
		{
			name:  "Категория уже существует",
			input: model.Category{UserID: 123, Name: "Спорт"},
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("AddCategory",
					mock.Anything,
					mock.AnythingOfType("*model.Category"),
				).Return(0, errors.New("Категория уже существует"))
			},
			wantError: true,
		},
		{
			name:  "Ошибка в репозитории",
			input: model.Category{ID: 123456, Name: "Спорт"},
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("AddCategory",
					mock.Anything,
					mock.AnythingOfType("*model.Category"),
				).Return(0, errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryMock := mocks.NewCategoryRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(categoryMock)
			}
			service := category.NewService(categoryMock)
			message, err := service.AddCategory(context.Background(), tt.input)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном создании категории, но его нет")
			}
		})
	}

	t.Log("Тест завершен")
}
func TestGetAllCategories(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(repo *mocks.CategoryRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:   "Успешное получение категорий",
			userID: 123,
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("GetAllCategories",
					mock.Anything,
					int64(123),
				).Return([]model.Category{
					{ID: 1, Name: "Еда", Color: ""},
					{ID: 2, Name: "Транспорт", Color: ""},
				}, nil)
			},
			wantError:   false,
			wantMessage: "📂 Ваши категории:",
		},
		{
			name:   "Пустой список категорий",
			userID: 123,
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("GetAllCategories",
					mock.Anything,
					int64(123),
				).Return([]model.Category{}, nil)
			},
			wantError:   false,
			wantMessage: "У вас пока нет категорий.",
		},
		{
			name:   "Ошибка в репозитории",
			userID: 123,
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("GetAllCategories",
					mock.Anything,
					int64(123),
				).Return(nil, errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryMock := mocks.NewCategoryRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(categoryMock)
			}
			service := category.NewService(categoryMock)
			message, err := service.GetAllCategories(context.Background(), tt.userID)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")
			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном получении категорий, но его нет")
			}
		})
	}

	t.Log("Тест завершен")
}
func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		id          int
		setupMock   func(repo *mocks.CategoryRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:   "Успешное удаление категории",
			userID: 123,
			id:     5,
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("DeleteCategory",
					mock.Anything,
					int64(123),
					5,
				).Return("Спорт", nil)
			},
			wantError:   false,
			wantMessage: "✅ Категория",
		},
		{
			name:      "Некорректно указан id категории",
			userID:    123,
			id:        0,
			setupMock: func(repo *mocks.CategoryRepository) {},
			wantError: true,
		},
		{
			name:   "Ошибка в репозитории",
			userID: 123,
			id:     5,
			setupMock: func(repo *mocks.CategoryRepository) {
				repo.On("DeleteCategory",
					mock.Anything,
					int64(123),
					5,
				).Return("", errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryMock := mocks.NewCategoryRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(categoryMock)
			}
			service := category.NewService(categoryMock)
			message, err := service.DeleteCategory(context.Background(), tt.userID, tt.id)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном удалении категории, но его нет")
			}
		})
	}

	t.Log("Тест завершен")
}
func TestAddExpense(t *testing.T) {
	tests := []struct {
		name        string
		input       model.Expense
		setupMock   func(repo *mocks.ExpenseRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:  "Успешное создание расхода",
			input: model.Expense{UserID: 1, Amount: 123, Category: "Транспорт", Description: "Поездка в трамвае"},
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("AddExpense",
					mock.Anything,
					mock.AnythingOfType("*model.Expense"),
				).Return(&model.Expense{UserID: 1, Amount: 123}, nil)
			},
			wantError:   false,
			wantMessage: "✅ Расход добавлен!",
		},
		{
			name:      "Расход отрицательный",
			input:     model.Expense{UserID: 1, Amount: -123, Category: "Категория", Description: "Поездка в трамвае"},
			setupMock: func(repo *mocks.ExpenseRepository) {},
			wantError: true,
		},
		{
			name:      "Расход нулевой",
			input:     model.Expense{UserID: 1, Amount: 0, Category: "Категория", Description: "Поездка в трамвае"},
			setupMock: func(repo *mocks.ExpenseRepository) {},
			wantError: true,
		},
		{
			name:      "Не указана категория",
			input:     model.Expense{UserID: 1, Amount: 123, Category: "", Description: "Поездка в трамвае"},
			setupMock: func(repo *mocks.ExpenseRepository) {},
			wantError: true,
		},
		{
			name:      "Не указано описание",
			input:     model.Expense{UserID: 1, Amount: 123, Category: "Категория", Description: ""},
			setupMock: func(repo *mocks.ExpenseRepository) {},
			wantError: true,
		},
		{
			name:  "Категория не найдена в БД",
			input: model.Expense{UserID: 1, Amount: 123, Category: "Космос", Description: "Поездка в трамвае"},
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("AddExpense",
					mock.Anything,
					mock.AnythingOfType("*model.Expense"),
				).Return(nil, errors.New("Категория не найдена"))
			},
			wantError: true,
		},
		{
			name:      "Не хватает описания расхода",
			input:     model.Expense{UserID: 1, Amount: 123, Category: "Космос", Description: ""},
			setupMock: func(repo *mocks.ExpenseRepository) {},
			wantError: true,
		},
		{
			name:  "Ошибка в репозитории",
			input: model.Expense{UserID: 1, Amount: 123, Category: "Космос", Description: "Поездка в трамвае"},
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("AddExpense",
					mock.Anything,
					mock.AnythingOfType("*model.Expense"),
				).Return(nil, errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenseMock := mocks.NewExpenseRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(expenseMock)
			}
			service := expense.NewService(expenseMock)
			message, err := service.AddExpense(context.Background(), &tt.input)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном добавлении расхода, но его нет")
			}
		})
	}

	t.Log("Тест завершен")
}
func TestTodayExpense(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(repo *mocks.ExpenseRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:   "Успешное получение расходов за сегодня",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("TodayExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{
					{Category: "Еда", Amount: 500},
					{Category: "Транспорт", Amount: 300},
				}, nil)
			},
			wantError:   false,
			wantMessage: "📊 Расходы за сегодня",
		},
		{
			name:   "Расходов за сегодня нет",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("TodayExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{}, nil)
			},
			wantError:   false,
			wantMessage: "📊 Расходы за сегодня",
		},
		{
			name:   "Расходы с одинаковой категорией",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("TodayExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{
					{Category: "Еда", Amount: 1234.56},
					{Category: "Еда", Amount: 65},
				}, nil)
			},
			wantError:   false,
			wantMessage: "Еда: 1299.56",
		},
		{
			name:   "Ошибка в репозитории",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("TodayExpense",
					mock.Anything,
					int64(123),
				).Return(nil, errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenseMock := mocks.NewExpenseRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(expenseMock)
			}
			service := expense.NewService(expenseMock)
			message, err := service.TodayExpense(context.Background(), tt.userID)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном получении расхода за сегодня, но его нет")
			}
		})
	}

	t.Log("Тест завершен")
}
func TestWeekExpense(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(repo *mocks.ExpenseRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:   "Успешное получение расходов за неделю",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("WeekExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{
					{Category: "Еда", Amount: 500},
				}, nil)
			},
			wantError:   false,
			wantMessage: "📊 Расходы за неделю",
		},
		{
			name:   "Расходов за неделю нет",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("WeekExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{}, nil)
			},
			wantError:   false,
			wantMessage: "📊 Нет расходов за неделю",
		},
		{
			name:   "Ошибка в репозитории",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("WeekExpense",
					mock.Anything,
					int64(123),
				).Return(nil, errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenseMock := mocks.NewExpenseRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(expenseMock)
			}
			service := expense.NewService(expenseMock)
			message, err := service.WeekExpense(context.Background(), tt.userID)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном получении расхода за неделю, но его нет")
			}
		})
	}

	t.Log("Тест завершен")
}
func TestMonthExpense(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(repo *mocks.ExpenseRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:   "Успешное получение расходов за месяц",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("MonthExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{
					{Category: "Еда", Amount: 5000},
				}, nil)
			},
			wantError:   false,
			wantMessage: "📊 Расходы за",
		},
		{
			name:   "Расходов за месяц нет",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("MonthExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{}, nil)
			},
			wantError:   false,
			wantMessage: "📊 Нет расходов за месяц",
		},
		{
			name:   "Ошибка в репозитории",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("MonthExpense",
					mock.Anything,
					int64(123),
				).Return(nil, errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenseMock := mocks.NewExpenseRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(expenseMock)
			}
			service := expense.NewService(expenseMock)
			message, err := service.MonthExpense(context.Background(), tt.userID)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Error("Ожидалась сообщение об успешном получении расхода за месяц, но его нет")
			}
		})
	}

	t.Log("Тест завершен")
}
func TestStatsExpense(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(repo *mocks.ExpenseRepository)
		wantError   bool
		wantMessage string
	}{
		{
			name:   "Успешное получение расходов за весь период",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("StatsExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{
					{Category: "Еда", Amount: 10000},
					{Category: "Транспорт", Amount: 5000},
				}, nil)
			},
			wantError:   false,
			wantMessage: "📈 Статистика расходов",
		},
		{
			name:   "Расходов нет",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("StatsExpense",
					mock.Anything,
					int64(123),
				).Return([]model.Expense{}, nil)
			},
			wantError:   false,
			wantMessage: "📊 Нет данных для статистики",
		},
		{
			name:   "Ошибка в репозитории",
			userID: 123,
			setupMock: func(repo *mocks.ExpenseRepository) {
				repo.On("StatsExpense",
					mock.Anything,
					int64(123),
				).Return(nil, errors.New("Ошибка БД"))
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenseMock := mocks.NewExpenseRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(expenseMock)
			}
			service := expense.NewService(expenseMock)
			message, err := service.StatsExpense(context.Background(), tt.userID)
			if !tt.wantError && err != nil {
				t.Error("Ошибка не ожидалась, но ее получили")

			}
			if tt.wantError && err == nil {
				t.Error("Ожидалась ошибка, но ее нет")
			}
			if !tt.wantError && !strings.Contains(message, tt.wantMessage) {
				t.Errorf("Ожидалось, что сообщение содержит %q, получено: %q", tt.wantMessage, message)
			}
		})
	}
	t.Log("Тест завершен")
}
