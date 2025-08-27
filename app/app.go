package app

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/IceMAN2377/market/internal/config"
	"github.com/IceMAN2377/market/internal/repository/postgres"
	"github.com/IceMAN2377/market/internal/service/subscription"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"net/url"

	v1Http "github.com/IceMAN2377/market/internal/transport/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

const (
	migrationsPath = "db/migrations"
	dbDriverName   = "postgres"
)

type App struct {
	router *http.ServeMux
	port   int
	logger *slog.Logger
}

func NewApp(config *config.Config, logger *slog.Logger) *App {
	// Строка подключения к PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresHost, config.PostgresPort, config.PostgresUser,
		config.PostgresPassword, config.PostgresDb, config.PostgresSslMode)

	// Подключение к базе данных
	db, err := sql.Open(dbDriverName, psqlInfo)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		panic("failed to connect to db: " + err.Error())
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", "error", err)
		panic("failed to ping database: " + err.Error())
	}

	psql := sqlx.NewDb(db, dbDriverName)

	// Применение миграций
	if config.PostgresMigrate {
		dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			config.PostgresUser, url.QueryEscape(config.PostgresPassword),
			config.PostgresHost, config.PostgresPort, config.PostgresDb, config.PostgresSslMode)

		m, err := migrate.New("file://"+migrationsPath, dbURI)
		if err != nil {
			logger.Error("Failed to create migrate object", "error", err)
			panic("failed to create migrate object: " + err.Error())
		}

		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			logger.Error("Failed to apply migrations", "error", err)
			panic("failed to apply migrations: " + err.Error())
		}

		logger.Info("Database migrations applied successfully")
	}

	// Инициализация слоев приложения
	repo := postgres.NewRepository(psql)
	service := subscription.NewService(repo)
	router := http.NewServeMux()

	// Регистрация HTTP endpoints
	v1Http.RegisterEndpoints(logger, router, service)

	logger.Info("Application initialized successfully")

	return &App{
		router: router,
		port:   config.HttpPort,
		logger: logger,
	}
}

func (a *App) Run() error {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", a.port),
		Handler:        a.addLoggingMiddleware(a.router),
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	a.logger.Info("HTTP server starting", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

// addLoggingMiddleware добавляет middleware для логирования HTTP запросов
func (a *App) addLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("HTTP request",
			"method", r.Method,
			"url", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(w, r)
	})
}
