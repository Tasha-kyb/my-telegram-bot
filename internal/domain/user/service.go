package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
func (p *Service) CreateProfile(ctx context.Context, req model.Profile) (string, error) {
	if req.ID == 0 || strings.TrimSpace(req.Username) == "" {
		return "", errors.New("❌ Не хватает параметров для создания профиля")
	}
	newProfile := &model.Profile{
		ID:        req.ID,
		Username:  req.Username,
		CreatedAt: time.Now(),
	}
	err := p.repo.CreateProfile(ctx, newProfile)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка при создании профиля, %w", err)
	}
	startMessage := `
	👋 Добро пожаловать в Expense Tracker!

	Я помогу вам отслеживать расходы и управлять бюджетами.

	✅ Вы зарегистрированы!
	📂 Созданы базовые категории:
   • Еда
   • Транспорт
   • Развлечения
   • Прочее
	`
	return startMessage, nil
}
