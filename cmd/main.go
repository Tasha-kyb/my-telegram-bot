package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	usecase "github.com/Tasha-kyb/my-telegram-bot/internal/app"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/category"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/expense"
	"github.com/Tasha-kyb/my-telegram-bot/internal/domain/user"
	"github.com/Tasha-kyb/my-telegram-bot/internal/handlers"
	httphandler "github.com/Tasha-kyb/my-telegram-bot/internal/handlers/http"
	"github.com/Tasha-kyb/my-telegram-bot/internal/handlers/telegram"
	"github.com/Tasha-kyb/my-telegram-bot/internal/repository"
	"github.com/Tasha-kyb/my-telegram-bot/internal/repository/database"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Printf("Предупреждение: .env: %v", err)
	}

	required := []string{
		"TELEGRAM_BOT_TOKEN",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
	}
	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("Обязательная переменная в .env не задана %s", env)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPool(ctx)
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к БД: %v", err)
	}
	defer pool.Close()
	log.Println("Подключение к БД")

	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		timezone = "Europe/Moscow"
		log.Printf("⚠️ TIMEZONE не задан, используется %s", timezone)
	}

	repo := repository.NewRepo(pool, timezone)

	userSvc := user.NewService(repo)
	categorySvc := category.NewService(repo)
	expenseSvc := expense.NewService(repo)

	uc := usecase.New(userSvc, categorySvc, expenseSvc)

	httpHandler := httphandler.NewHandler(uc)
	router := handlers.NewRouter(httpHandler)

	tgHandler, err := telegram.NewTelegramUpdates(uc)
	if err != nil {
		return fmt.Errorf("Ошибка создания соединения с Telegram ботом: %v", err)
	}

	go tgHandler.StartUpdates(ctx)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}
	go func() {
		log.Printf("HTTP сервер запущен на %s:", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка сервера, %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Получен сигнал, завершаем работу")
	cancel()

	shutdownTimeout := os.Getenv("SHUTDOWN_TIMEOUT")
	timeoutSec, _ := strconv.Atoi(shutdownTimeout)
	if timeoutSec == 0 {
		timeoutSec = 10
	}
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeoutSec)*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	tgHandler.Wg.Wait()
	log.Println("Приложение завершено")
	return nil
}
