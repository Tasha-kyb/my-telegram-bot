package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryT struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *RepositoryT {
	return &RepositoryT{pool: pool}
}

func (p *RepositoryT) CreateProfile(ctx context.Context, profile *model.Profile) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Не удалось начать транзакцию, %w", err)
	}
	defer tx.Rollback(ctx)

	userQuery := `
	INSERT INTO users (user_id, username, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id) DO UPDATE
	SET username = EXCLUDED.username`
	_, err = tx.Exec(ctx, userQuery, profile.ID, profile.Username, profile.Created_at)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		return fmt.Errorf("Ошибка при создании пользователя, %w", err)
	}
	categoryQuery := `
	INSERT INTO categories (user_id, name) VALUES 
	($1, 'Еда'), 
	($1, 'Транспорт'), 
	($1, 'Развлечения'), 
	($1, 'Прочее')
	ON CONFLICT (user_id, name) DO NOTHING`
	_, err = tx.Exec(ctx, categoryQuery, profile.ID)
	if err != nil {
		log.Printf("Ошибка при создании категорий: %v", err)
		return fmt.Errorf("Ошибка при создании категорий, %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Ошибка при завершении транзакции: %v", err)
		return fmt.Errorf("Ошибка при завершении транзакции: %w", err)
	}

	return nil
}
func (p *RepositoryT) AddCategory(ctx context.Context, category *model.Category) (int, error) {
	query := `
	INSERT INTO categories (user_id, name, color) 
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, name) DO NOTHING
	RETURNING id`
	var id int
	err := p.pool.QueryRow(ctx, query, category.UserID, category.Name, category.Color).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("Категория %s уже существует", category.Name)
		}
		log.Printf("Ошибка при создании категории: %v", err)
		return 0, fmt.Errorf("Ошибка при создании категории: %w", err)
	}
	return id, nil
}
func (p *RepositoryT) GetAllCategories(ctx context.Context, userID int64) ([]model.Category, error) {
	query := `
	SELECT id, name, COALESCE(color, '') as color
	FROM categories WHERE user_id = $1
	ORDER BY id`
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса категорий из базы данных, %w", err)
	}
	defer rows.Close()

	var allCategories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Color,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при получении списка категорий, %w", err)
		}
		allCategories = append(allCategories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при чтении категорий, %w", err)
	}
	if allCategories == nil {
		allCategories = []model.Category{}
	}
	return allCategories, nil
}
func (p *RepositoryT) DeleteCategory(ctx context.Context, userID int64, id int) (string, error) {
	query := `
	DELETE FROM categories WHERE user_id = $1 AND id = $2
	RETURNING name`
	var name string
	err := p.pool.QueryRow(ctx, query, userID, id).Scan(&name)
	if err != nil {
		log.Printf("Ошибка при удалении категории: %v", err)
		return "", fmt.Errorf("Ошибка при удалении категории: %w", err)
	}
	return name, nil
}
func (p *RepositoryT) AddExpense(ctx context.Context, expense *model.Expense) (*model.Expense, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Не удалось начать транзакцию, %w", err)
	}
	defer tx.Rollback(ctx)

	var profileExist bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`,
		expense.UserID).Scan(&profileExist)
	if err != nil {
		return nil, fmt.Errorf("Ошибка проверки наличия пользователя в базе данных, %w", err)
	}
	if !profileExist {
		return nil, fmt.Errorf("Пользователь c ID %d еще не зарегистрирован", expense.UserID)
	}

	var categoryExist bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM categories WHERE user_id = $1 AND name = $2)`,
		expense.UserID, expense.Category).Scan(&categoryExist)
	if err != nil {
		return nil, fmt.Errorf("Ошибка проверки категории, %w", err)
	}
	if !categoryExist {
		return nil, fmt.Errorf("Указанная категория не найдена в базе данных, %w", err)
	}

	query := `
	INSERT INTO expenses (user_id, amount, category, description, created_at) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING user_id, amount, category, description, created_at
	`
	var response model.Expense
	err = tx.QueryRow(ctx, query,
		expense.UserID,
		expense.Amount,
		expense.Category,
		expense.Description,
		expense.Created_at,
	).Scan(
		&response.UserID,
		&response.Amount,
		&response.Category,
		&response.Description,
		&response.Created_at,
	)
	if err != nil {
		log.Printf("Ошибка при создании расхода, %v", err)
		return nil, fmt.Errorf("Ошибка при создании расхода, %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Ошибка при завершении транзакции: %v", err)
		return nil, fmt.Errorf("Ошибка при завершении транзакции: %w", err)
	}
	return &response, nil
}
func (p *RepositoryT) TodayExpense(ctx context.Context, userID int64) ([]model.Expense, error) {
	query := `
	SELECT amount, category, description, created_at
	FROM expenses 
	WHERE user_id = $1 
  		AND created_at >= (CURRENT_DATE::timestamptz AT TIME ZONE 'Europe/Moscow')::timestamp
  		AND created_at <  (CURRENT_DATE::timestamptz AT TIME ZONE 'Europe/Moscow' + interval '1 day')::timestamp
	ORDER BY created_at DESC;`
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Ошибка при получении расхода за текущий день, %v", err)
		return nil, fmt.Errorf("Ошибка при получении расхода за текущий день, %v", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			&expense.Amount,
			&expense.Category,
			&expense.Description,
			&expense.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании расхода за текущий день, %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при сканировании расходов, %w", err)
	}

	return expenses, nil
}
func (p *RepositoryT) WeekExpense(ctx context.Context, userID int64) ([]model.Expense, error) {
	query := `
	SELECT amount, category, created_at
	FROM expenses 
	WHERE user_id = $1
  		AND DATE_TRUNC('week', created_at) = 
      		DATE_TRUNC('week', (CURRENT_DATE::timestamptz AT TIME ZONE 'Europe/Moscow')::timestamp)
	ORDER BY created_at DESC;`
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Ошибка при получении расхода за текущую неделю, %v", err)
		return nil, fmt.Errorf("Ошибка при получении расхода за текущую неделю, %v", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			&expense.Amount,
			&expense.Category,
			&expense.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании расхода за текущую неделю, %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при сканировании расходов, %w", err)
	}

	return expenses, nil
}
func (p *RepositoryT) MonthExpense(ctx context.Context, userID int64) ([]model.Expense, error) {
	query := `
	SELECT amount, category, created_at
	FROM expenses 
	WHERE user_id = $1
  		AND DATE_TRUNC('month', created_at) = 
      		DATE_TRUNC('month', (CURRENT_DATE::timestamptz AT TIME ZONE 'Europe/Moscow')::timestamp)
	ORDER BY created_at DESC;`
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Ошибка при получении расхода за текущий месяц, %v", err)
		return nil, fmt.Errorf("Ошибка при получении расхода за текущий месяц, %v", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			&expense.Amount,
			&expense.Category,
			&expense.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании расхода за текущий месяц, %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при сканировании расходов, %w", err)
	}

	return expenses, nil
}
func (p *RepositoryT) StatsExpense(ctx context.Context, userID int64) ([]model.Expense, error) {
	query := `
	SELECT amount, category, created_at
	FROM expenses 
	WHERE user_id = $1
	ORDER BY created_at DESC`
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Ошибка при получении расхода за весь период, %v", err)
		return nil, fmt.Errorf("Ошибка при получении расхода за весь период, %v", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			&expense.Amount,
			&expense.Category,
			&expense.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании расхода за весь период, %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при сканировании расходов, %w", err)
	}

	return expenses, nil
}
