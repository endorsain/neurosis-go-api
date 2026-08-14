package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	authModule "github.com/endorsain/neurosis-go-api/internal/auth/module"
	"github.com/endorsain/neurosis-go-api/internal/config"
	"github.com/endorsain/neurosis-go-api/internal/infrastructure/postgres"
	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
	usersModule "github.com/endorsain/neurosis-go-api/internal/users/module"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	cfg := config.Load()

	logger.Printf(
		"configuration loaded (app=%s env=%s server=%s db_host=%s db_port=%s db_name=%s db_user=%s db_sslmode=%s)",
		cfg.Application.Name,
		cfg.Application.Environment,
		cfg.Server.Address(),
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.User,
		cfg.Database.SSLMode,
	)

	pgClient, err := postgres.New(context.Background(), cfg.Database, logger)
	if err != nil {
		logger.Printf("failed to connect to PostgreSQL: %v", err)
		os.Exit(1)
	}

	users := usersModule.New(pgClient.DB())
	auth, err := authModule.New(pgClient.DB(), users.UserRepository(), cfg.Auth)
	if err != nil {
		logger.Printf("failed to initialize auth module: %v", err)
		_ = pgClient.Close()
		os.Exit(1)
	}

	router := httptransport.NewRouter()
	auth.RegisterRoutes(router)
	users.RegisterRoutes(router, auth.RequireAuthentication)

	server := httptransport.NewServer(cfg.Server, router)

	logger.Println("starting server")
	logger.Printf("HTTP server starting on %s", server.Address())

	serverErrCh := make(chan error, 1)
	go func() {
		err := server.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		close(serverErrCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		if err != nil {
			logger.Printf("HTTP server stopped with error: %v", err)
		}
	case sig := <-signalCh:
		logger.Printf("shutdown signal received: %s", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("HTTP server stopped with error: %v", err)
	}

	if err := pgClient.Close(); err != nil {
		logger.Printf("failed to close PostgreSQL connection: %v", err)
	}
}
