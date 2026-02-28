package expense

import (
	"context"
	"errors"
	"fmt"
	"sort"
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

func (p *Service) AddExpense(ctx context.Context, req *model.Expense) (string, error) {
	if req.Amount <= 0 {
		return "", errors.New("❌ Сумма расхода должна быть положительной")
	}
	if req.Category == "" || req.Description == "" {
		return "", errors.New("❌ Не хватает данных для добавления расхода")
	}
	if req.Created_at.IsZero() {
		req.Created_at = time.Now()
	}
	newExpense := &model.Expense{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Created_at:  req.Created_at,
	}
	expense, err := p.repo.AddExpense(ctx, newExpense)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена в базе данных") {
			return "", fmt.Errorf("❌ Категория \"%s\" не найдена", req.Category)
		}
		return "", fmt.Errorf("❌ Ошибка при создании расхода %w", err)
	}
	addExpenseMessage := fmt.Sprintf(`
	✅ Расход добавлен!

	💰 Сумма: %.2f₽
	📂 Категория: %s
	📝 Описание: %s
	📅 Дата: %s

	💵 Осталось до лимита: X
	`, expense.Amount, expense.Category, expense.Description, expense.Created_at.Format("02.01.2006"))

	return addExpenseMessage, nil
}
func (p *Service) TodayExpense(ctx context.Context, userID int64) (string, error) {
	expenses, err := p.repo.TodayExpense(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка при при получении расходов за сегодня %w", err)
	}
	today := time.Now().Format("02.01.2006")
	if len(expenses) == 0 {
		return fmt.Sprintf(`📊 Расходы за сегодня (%s)
		
		Пока нет расходов за сегодня.
		Используйте /add для добавления расхода.`, today), nil
	}
	categoriesMap := make(map[string][]model.Expense)

	for _, expense := range expenses {
		categoriesMap[expense.Category] = append(categoriesMap[expense.Category], expense)
	}

	response := fmt.Sprintf("📊 Расходы за сегодня (%s)\n\n", today)
	total := 0.0

	for category, expenseList := range categoriesMap {
		sum := 0.0
		for _, exp := range expenseList {
			sum += exp.Amount
		}
		response += fmt.Sprintf("%s: %.2f₽\n", category, sum)

		for _, exp := range expenseList {
			response += fmt.Sprintf("   • %s: %.2f₽\n", exp.Description, exp.Amount)
		}
		total += sum
	}
	response += "\n━━━━━━━━━━━━━━━━━━━━\n"
	response += fmt.Sprintf("💰 Итого: %.2f₽", total)

	return response, nil
}
func (p *Service) WeekExpense(ctx context.Context, userID int64) (string, error) {
	expenses, err := p.repo.WeekExpense(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка при при получении расходов за неделю %w", err)
	}
	now := time.Now()
	weekDay := int(now.Weekday())
	if weekDay == 0 {
		weekDay = 7
	}
	startOfWeek := now.AddDate(0, 0, -weekDay+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	if len(expenses) == 0 {
		return fmt.Sprintf("📊 Нет расходов за неделю (%s - %s). Используйте /add для добавления расхода",
			startOfWeek.Format("02.01"), endOfWeek.Format("02.01")), nil
	}

	dayNames := []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
	daySum := make(map[string]float64)
	categorySum := make(map[string]float64)

	total := 0.0
	for _, exp := range expenses {
		idx := int(exp.Created_at.Weekday())
		if idx == 0 {
			idx = 7
		}
		dayName := dayNames[idx-1]
		daySum[dayName] += exp.Amount
		categorySum[exp.Category] += exp.Amount
		total += exp.Amount
	}
	response := fmt.Sprintf("📊 Расходы за неделю (%s - %s)\n\n",
		startOfWeek.Format("02.01"), endOfWeek.Format("02.01"))

	// вывод по дням недели
	for _, day := range dayNames {
		if sum, ok := daySum[day]; ok && sum > 0 {
			response += fmt.Sprintf("%s: %.2f₽\n", day, sum)
		}
	}

	response += fmt.Sprintf("\n ━━━━━━━━━━━━━━━━━━━━\n 💰 Итого: %.2f₽\n", total)
	response += fmt.Sprintf("📈 Средний расход в день: %.2f₽\n", total/float64(len(daySum)))

	type statistics struct {
		Name    string
		Sum     float64
		Percent float64
	}

	stats := make([]statistics, 0, len(categorySum))
	for name, sum := range categorySum {
		stats = append(stats, statistics{
			Name:    name,
			Sum:     sum,
			Percent: (sum / total) * 100})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Sum > stats[j].Sum
	})
	response += "🏆 Топ категории:\n"
	for i, s := range stats {
		if i >= 3 {
			break
		}
		response += fmt.Sprintf("   %d. %s: %.0f₽ (%.0f%%)\n", i+1, s.Name, s.Sum, s.Percent)
	}

	return response, nil
}
func (p *Service) MonthExpense(ctx context.Context, userID int64) (string, error) {
	expenses, err := p.repo.MonthExpense(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка при при получении расходов за месяц %w", err)
	}
	monthNames := []string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}

	if len(expenses) == 0 {
		return fmt.Sprintf("📊 Нет расходов за месяц %s. Используйте /add для добавления расхода",
			monthNames[time.Now().Month()-1]), nil
	}

	categorySum := make(map[string]float64)

	total := 0.0
	for _, exp := range expenses {
		categorySum[exp.Category] += exp.Amount
		total += exp.Amount
	}
	response := fmt.Sprintf("📊 Расходы за месяц (%s)\n\n", monthNames[time.Now().Month()-1])

	type statistics struct {
		Name string
		Sum  float64
	}

	stats := make([]statistics, 0, len(categorySum))
	for name, sum := range categorySum {
		stats = append(stats, statistics{
			Name: name,
			Sum:  sum,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Sum > stats[j].Sum
	})
	for i, s := range stats {
		if i >= 3 {
			break
		}
		response += fmt.Sprintf("   %d. %s: %.0f₽\n", i+1, s.Name, s.Sum)
	}
	response += fmt.Sprintf("\n ━━━━━━━━━━━━━━━━━━━━\n 💰 Итого: %.2f₽\n", total)

	return response, nil
}
func (p *Service) StatsExpense(ctx context.Context, userID int64) (string, error) {
	expenses, err := p.repo.StatsExpense(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка при получении расходов за весь период %w", err)
	}

	if len(expenses) == 0 {
		return "📊 Нет данных для статистики. Используйте /add для добавления расхода", nil
	}

	categorySum := make(map[string]float64)
	dailySum := make(map[string]float64)

	total := 0.0
	var firstDate, lastDate time.Time

	for _, exp := range expenses {
		categorySum[exp.Category] += exp.Amount
		total += exp.Amount

		date := exp.Created_at.Format("2006-01-02")
		dailySum[date] += exp.Amount

		if firstDate.IsZero() || exp.Created_at.Before(firstDate) {
			firstDate = exp.Created_at
		}
		if lastDate.IsZero() || exp.Created_at.After(lastDate) {
			lastDate = exp.Created_at
		}

	}
	days := int(lastDate.Sub(firstDate).Hours()/24) + 1

	avgDay := total / float64(days)
	avgWeek := total / (float64(days) / 7)
	avgMonth := total / (float64(days) / 30.44)

	response := "📈 Статистика расходов\n\n"
	response += fmt.Sprintf("💰 Всего потрачено: %.0f₽\n", total)
	response += fmt.Sprintf("📊 Всего транзакций: %d\n\n", len(expenses))
	response += "📅 Средний расход:\n"
	response += fmt.Sprintf("   • В день: %.0f₽\n", avgDay)
	response += fmt.Sprintf("   • В неделю: %.0f₽\n", avgWeek)
	response += fmt.Sprintf("   • В месяц: %.0f₽\n\n", avgMonth)

	type statistics struct {
		Name string
		Sum  float64
	}

	stats := make([]statistics, 0, len(categorySum))
	for name, sum := range categorySum {
		stats = append(stats, statistics{
			Name: name,
			Sum:  sum,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Sum > stats[j].Sum
	})

	response += "🏆 Топ категории:\n"

	for i, s := range stats {
		if i >= 4 {
			break
		}
		percent := (s.Sum / total) * 100
		response += fmt.Sprintf("   %d. %s: %.0f₽ (%.0f%%)\n", i+1, s.Name, s.Sum, percent)
	}

	now := time.Now()
	lastMonthSum, prevMonthSum := 0.0, 0.0

	for _, exp := range expenses {
		if exp.Created_at.After(now.AddDate(0, 0, -30)) {
			lastMonthSum += exp.Amount
		} else if exp.Created_at.After(now.AddDate(0, 0, -60)) {
			prevMonthSum += exp.Amount
		}
	}

	if prevMonthSum == 0 {
		response += "\n📉 Тренд: нет данных за последние месяцы\n"
	}

	if prevMonthSum > 0 {
		percent := (lastMonthSum - prevMonthSum) / prevMonthSum * 100
		if percent >= 0 {
			response += fmt.Sprintf("\n📈 Тренд: +%.0f%% к прошлому месяцу\n", percent)
		} else {
			response += fmt.Sprintf("\n📉 Тренд: %.0f%% к прошлому месяцу\n", percent)
		}
	}

	maxDay := ""
	maxSum := 0.0

	for day, sum := range dailySum {
		if sum > maxSum {
			maxSum = sum
			maxDay = day
		}
	}
	if maxDay != "" {
		maxDate, _ := time.Parse("2006-01-02", maxDay)
		maxDayFormatted := maxDate.Format("02.01.2006")
		response += fmt.Sprintf("📅 Самый дорогой день: %s (%.0f₽)", maxDayFormatted, maxSum)
	}

	return response, nil
}
