package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example/internal/app"
)

func main() {
	// Инициализируем приложение
	application := app.NewApp()
	defer application.Close()

	server := application.Server()

	// Канал для обработки сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной goroutine
	go func() {
		log.Printf("api server started on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-sigChan

	// Graceful shutdown — ждём завершения текущих запросов (max 10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("server stopped")
}
