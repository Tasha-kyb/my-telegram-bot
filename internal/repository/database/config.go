package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func getConnectionString() string {
	godotenv.Load()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
}

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	// Конфигурация пула
	config, err := pgxpool.ParseConfig(getConnectionString())
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	// Настройки пула
	config.MaxConns = 10                      // Максимум соединений
	config.MinConns = 5                       // Минимум соединений
	config.MaxConnLifetime = time.Hour        // Время жизни соединения
	config.MaxConnIdleTime = 30 * time.Minute // Время простоя
	config.HealthCheckPeriod = time.Minute    // Проверка здоровья

	// Создаем пул
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// Проверяем соединение
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

/*
1. Читаем переменные окружения из .env
2. Формируем строку подключения
3. Парсим строку в объект конфигурации
4. Настраиваем параметры пула
5. Создаем пул с нашей конфигурацией
6. Проверяем, что БД доступна
7. Возвращаем готовый пул для использования*/
