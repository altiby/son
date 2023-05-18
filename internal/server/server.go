package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/altiby/son/internal/config"
	gometrics "github.com/altiby/son/internal/gometrics"
	"github.com/altiby/son/internal/handlers"
	"github.com/altiby/son/internal/hasher"
	"github.com/altiby/son/internal/logger"
	"github.com/altiby/son/internal/storage"
	"github.com/altiby/son/internal/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"net"
	"net/http"
	"strconv"
)

type Builder struct {
	port    int
	metrics gometrics.Metrics
	logger  *logger.Log
	ctx     context.Context
	cfg     *config.Config
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) WithMetrics(m gometrics.Metrics) *Builder {
	b.metrics = m
	return b
}

func (b *Builder) WithLogger(l *logger.Log) *Builder {
	b.logger = l
	return b
}

func (b *Builder) WithPort(port int) *Builder {
	b.port = port
	return b
}

func (b *Builder) WithContext(ctx context.Context) *Builder {
	b.ctx = ctx
	return b
}

func (b *Builder) WithConfig(cfg *config.Config) *Builder {
	b.cfg = cfg
	return b
}

func (b *Builder) initPostgresql() (*storage.Postgres, error) {
	// postgres
	storage, err := storage.New(context.TODO(), &b.cfg.Postgresql)
	if err != nil {
		return nil, err
	}

	conn, err := storage.Connection()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", b.cfg.Postgresql.MigrationDir),
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return storage, nil
}

func (b *Builder) Build() (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		b.logger.LogMiddleware,
		middleware.Recoverer,
		middleware.RequestID,
	)

	postgres, err := b.initPostgresql()
	if err != nil {
		return nil, err
	}

	h := hasher.New()

	userStorage := storage.NewUserStorage(postgres)
	userService := user.NewService(userStorage, h)
	userHandler := handlers.NewUserHandler(userService)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/user", userHandler)
	})

	server := &http.Server{
		Addr:    net.JoinHostPort("", strconv.Itoa(b.port)),
		Handler: router,
	}

	return server, nil
}
