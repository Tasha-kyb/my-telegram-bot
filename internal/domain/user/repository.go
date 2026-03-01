package user

import (
	"context"

	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

//go:generate mockery --name=Repository --output=../mocks --outpkg=mocks
type Repository interface {
	CreateProfile(ctx context.Context, req *model.Profile) error
}
