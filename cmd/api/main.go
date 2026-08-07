package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	authHandlers "github.com/endorsain/neurosis-go-api/internal/auth/handlers"
	authUsecases "github.com/endorsain/neurosis-go-api/internal/auth/usecases"
	usersDomain "github.com/endorsain/neurosis-go-api/internal/users"
	usersHandlers "github.com/endorsain/neurosis-go-api/internal/users/handlers"
	usersUsecases "github.com/endorsain/neurosis-go-api/internal/users/usecases"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	dsn := databaseURLFromEnv()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepository := usersDomain.NewPostgresUserRepository(db)
	registerUseCase := authUsecases.NewRegisterUserUseCase(userRepository)
	getCurrentUserUseCase := usersUsecases.NewGetCurrentUserUseCase(userRepository)
	getUserByIDUseCase := usersUsecases.NewGetUserByIDUseCase(userRepository)
	listUsersUseCase := usersUsecases.NewListUsersUseCase(userRepository)

	registerHandler := authHandlers.NewHandler(registerUseCase)
	getCurrentUserHandler := usersHandlers.NewGetCurrentUserHandler(getCurrentUserUseCase)
	getUserByIDHandler := usersHandlers.NewGetUserByIDHandler(getUserByIDUseCase)
	listUsersHandler := usersHandlers.NewListUsersHandler(listUsersUseCase)

	r := chi.NewRouter()
	authHandlers.RegisterRoutes(r, registerHandler)
	usersHandlers.RegisterRoutes(r, getCurrentUserHandler, getUserByIDHandler, listUsersHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func databaseURLFromEnv() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	host := getenvOrDefault("DB_HOST", "localhost")
	port := getenvOrDefault("DB_PORT", "5432")
	database := getenvOrDefault("DB_NAME", "neurosis")
	user := getenvOrDefault("DB_USER", "postgres")
	password := getenvOrDefault("DB_PASSWORD", "postgres")
	sslmode := getenvOrDefault("DB_SSLMODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, database, sslmode)
}

func getenvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
