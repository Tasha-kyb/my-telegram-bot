package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/internal/handlers"
	"github.com/internal/repository"
	"github.com/internal/repository/database"
	"github.com/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer pool.Close()
	fmt.Println("Подключение к БД")

	repo := repository.NewRepo(pool)
	service := usecase.NewService(repo)

	httpHandler := handlers.NewHandler(service)
	router := handlers.NewRouter(httpHandler)

	tgHandler, err := handlers.NewTelegramUpdates(service)
	if err != nil {
		log.Fatalf("Ошибка создания соединения с Telegram ботом: %v", err)
	}

	go tgHandler.StartUpdates(ctx)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		log.Println("HTTP сервер запущен на :8080")
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Ошибка сервера, %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Получен сигнал, завершаем работу")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Приложение завершено")
}
