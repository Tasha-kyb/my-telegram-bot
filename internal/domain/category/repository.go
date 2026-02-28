package category

import (
	"context"

	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

//go:generate mockery --name=Repository --output=../mocks --outpkg=mocks
type Repository interface {
	AddCategory(ctx context.Context, req *model.Category) (int, error)
	GetAllCategories(ctx context.Context, userID int64) ([]model.Category, error)
	DeleteCategory(ctx context.Context, userID int64, id int) (string, error)
}
