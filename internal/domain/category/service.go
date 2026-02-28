package category

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (p *Service) AddCategory(ctx context.Context, req model.Category) (string, error) {
	if strings.TrimSpace(req.Name) == "" {
		return "", errors.New("❌ Не хватает параметров для создания категории")
	}
	newCategory := &model.Category{
		UserID: req.UserID,
		Name:   req.Name,
		Color:  req.Color,
	}
	id, err := p.repo.AddCategory(ctx, newCategory)
	if err != nil {
		if strings.Contains(err.Error(), "уже существует") {
			return "", fmt.Errorf("❌ Категория %s уже существует", req.Name)
		}
		return "", fmt.Errorf("❌ Ошибка при создании категории, %w", err)
	}
	addCategoryMessage := fmt.Sprintf(`
	✅ Категория создана!
	📂 Название: %s
	🎨 Цвет: %s
	🆔 ID: %d
	Используйте этот ID для удаления категории.
	`, req.Name, req.Color, id)

	return addCategoryMessage, nil
}
func (p *Service) GetAllCategories(ctx context.Context, userID int64) (string, error) {
	categoriesDB, err := p.repo.GetAllCategories(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка при получении категорий: %w", err)
	}
	if len(categoriesDB) == 0 {
		return "У вас пока нет категорий. \nСоздать категорию можно командой /category add", nil
	}
	emojiMap := map[string]string{
		"Еда":         "🍔",
		"Транспорт":   "🚗",
		"Развлечения": "🎬",
		"Прочее":      "📦",
		"Спорт":       "⚽",
		"Красота":     "💄",
		"Магазин":     "🛒",
		"Растения":    "🌿",
		"Цветы":       "🌸",
	}
	response := "📂 Ваши категории:\n\n"

	for _, category := range categoriesDB {
		emoji := emojiMap[category.Name]
		if emoji == "" {
			emoji = "📂"
		}

		response += fmt.Sprintf("%s %s\n", emoji, category.Name)
		response += fmt.Sprintf("	ID: %d\n\n", category.ID)
	}
	response += "💡 Используйте ID для удаления категории"
	return response, nil
}
func (p *Service) DeleteCategory(ctx context.Context, userID int64, id int) (string, error) {
	if id <= 0 {
		return "", errors.New("❌ Ошибка: некорректно указан id категории")
	}
	categoryName, err := p.repo.DeleteCategory(ctx, userID, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", fmt.Errorf("❌ Ошибка: некорректно указан ID категории")
		}
		return "", fmt.Errorf("❌ Ошибка при удалении категории: %w", err)
	}
	deleteCategoryMessage := fmt.Sprintf(`
	✅ Категория %s удалена
	Все расходы из этой категории перенесены в "Прочее"
	`, categoryName)
	return deleteCategoryMessage, nil
}
