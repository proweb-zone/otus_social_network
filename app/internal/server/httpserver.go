package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ms_baskets/internal/config"
	"ms_baskets/pkg/logger"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func StartServer(router *chi.Mux, logger *logger.Logger, config *config.Config) *http.Server {
	server := &http.Server{
		Addr:         config.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		logger.Info("сервер запущен", map[string]any{
			"address": config.HTTPServer.Address,
		})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	catchServerEvents(sigCh, serverErrCh, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(errors.New(err.Error()), "сервер завершил работу с ошибкой", nil)
	}

	return server
}

func catchServerEvents(sigCh chan os.Signal, serverErrCh chan error, logger *logger.Logger) {
	select {
	case err := <-serverErrCh:
		logger.Fatal(errors.New(err.Error()), "Не удалось запустить приложение", nil)

	case sig := <-sigCh:
		logger.Info("завершение работы...", map[string]any{"event": sig.String()})
	}
}
