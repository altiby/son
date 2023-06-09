package server

import (
	"context"
	"github.com/altiby/son/internal/config"
	gometrics "github.com/altiby/son/internal/gometrics"
	"github.com/altiby/son/internal/handlers"
	"github.com/altiby/son/internal/hasher"
	"github.com/altiby/son/internal/logger"
	"github.com/altiby/son/internal/metrics"
	"github.com/altiby/son/internal/storage"
	"github.com/altiby/son/internal/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"net"
	"net/http"
	"strconv"
)

type Builder struct {
	port    int
	appName string
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

func (b *Builder) WithAppName(n string) *Builder {
	b.appName = n
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
	return storage.InitPostgresql(b.cfg.Postgresql)
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

	router, err = metrics.InstrumentChiRouter(router, b.appName)
	if err != nil {
		return nil, err
	}

	_, err = metrics.NewPrometheus(b.appName)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    net.JoinHostPort("", strconv.Itoa(b.port)),
		Handler: router,
	}

	return server, nil
}
