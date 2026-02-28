package telegram

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Tasha-kyb/my-telegram-bot/internal/app"
	"github.com/Tasha-kyb/my-telegram-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramHandler struct {
	usecase *app.Service
	bot     *tgbotapi.BotAPI
	Wg      sync.WaitGroup
}

func NewTelegramUpdates(usecase *app.Service) (*TelegramHandler, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return nil, err
	}
	return &TelegramHandler{
		usecase: usecase,
		bot:     bot,
	}, nil
}

func (t *TelegramHandler) StartUpdates(ctx context.Context) {
	log.Println("Бот с воркерами запущен")

	updatesChan := make(chan tgbotapi.Update, 100)

	for i := 0; i < 10; i++ {
		t.Wg.Add(1)
		go func(worker int) {
			defer t.Wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case update, ok := <-updatesChan:
					if !ok {
						return
					}
					t.handleMessage(update)
				}
			}
		}(i)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			close(updatesChan)
			t.Wg.Wait()
			log.Println("Бот с воркерами остановлен")
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			updatesChan <- update
		}

	}
}
func (t *TelegramHandler) handleMessage(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if update.Message.From == nil {
		log.Println("⚠️ Данные о пользователе отсутствуют")
		return
	}
	switch {
	case update.Message.Text == "/start":
		profile := model.Profile{
			ID:        int64(update.Message.From.ID),
			Username:  update.Message.From.UserName,
			CreatedAt: time.Now(),
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := t.usecase.CreateProfile(ctx, profile)

		if err != nil {
			log.Printf("❌Ошибка при создании профиля, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка при создании профиля"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	case strings.HasPrefix(update.Message.Text, "/category add"):
		parts := strings.Fields(update.Message.Text)
		if len(parts) < 3 {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка: вы не указали название категории"))
			return
		}
		categoryName := parts[2]
		color := ""
		if len(parts) >= 4 {
			color = parts[3]
		}
		newCategory := model.Category{
			UserID: int64(update.Message.From.ID),
			Name:   categoryName,
			Color:  color,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := t.usecase.AddCategory(ctx, newCategory)

		if err != nil {
			log.Printf("❌Ошибка создания категории, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка создания категории"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	case update.Message.Text == "/categories":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		categories, err := t.usecase.GetAllCategories(ctx, update.Message.From.ID)

		if err != nil {
			log.Printf("❌Ошибка получения категорий, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка получения категорий"))
			return
		}
		t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, categories))
	case strings.HasPrefix(update.Message.Text, "/category delete"):
		parts := strings.Fields(update.Message.Text)
		if len(parts) < 3 {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌Ошибка: Вы не указали id категории для удаления"))
			return
		}
		idstr := parts[2]
		id, err := strconv.Atoi(idstr)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌Ошибка: некорректно указан id категории"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		response, err := t.usecase.DeleteCategory(ctx, update.Message.From.ID, id)
		if err != nil {
			log.Printf("❌Ошибка при удалении категории, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка при удалении категории"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	case update.Message.Text == "/help":
		helpText := `
			📖 Доступные команды:

			💰 Расходы:
			/add <сумма> <категория> <описание> — добавить расход
			/today — расходы за сегодня
			/week — расходы за неделю
			/month — расходы за месяц
			/stats — общая статистика

			📂 Категории:
			/category add <название> <цвет> — создать категорию
			/categories — список категорий
			/category delete <id> — удалить категорию`
		t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, helpText))
	case strings.HasPrefix(update.Message.Text, "/add"):
		parts := strings.Fields(update.Message.Text)
		if len(parts) < 4 {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌Ошибка: вы не указали все параметры (сумма, категория и описание)"))
			return
		}
		amount, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Сумма должна быть числом"))
			return
		}
		category := parts[2]
		description := strings.Join(parts[3:], " ")
		newExpense := model.Expense{
			UserID:      int64(update.Message.From.ID),
			Category:    category,
			Amount:      amount,
			Description: description,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := t.usecase.AddExpense(ctx, &newExpense)

		if err != nil {
			log.Printf("❌Ошибка создания расхода, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка создания расхода"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	case update.Message.Text == "/today":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := t.usecase.TodayExpense(ctx, update.Message.From.ID)
		if err != nil {
			log.Printf("❌Ошибка получения расходов за сегодня, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка получения расходов за сегодня"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	case update.Message.Text == "/week":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := t.usecase.WeekExpense(ctx, update.Message.From.ID)
		if err != nil {
			log.Printf("❌Ошибка получения расходов за неделю, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка получения расходов за неделю"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	case update.Message.Text == "/month":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := t.usecase.MonthExpense(ctx, update.Message.From.ID)
		if err != nil {
			log.Printf("❌Ошибка получения расходов за месяц, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка получения расходов за месяц"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	case update.Message.Text == "/stats":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := t.usecase.StatsExpense(ctx, update.Message.From.ID)
		if err != nil {
			log.Printf("❌Ошибка получения расходов за весь период, %v", err)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка получения расходов за весь период"))
			return
		}
		if _, err := t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response)); err != nil {
			log.Printf("❌Ошибка отправки сообщения: %v", err)
		}
	default:
		t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌Неизвестная команда, используйте /help"))
	}
}
