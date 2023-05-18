package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/altiby/son/internal/logger"
	"github.com/altiby/son/internal/server"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/altiby/son/internal/config"
	"github.com/altiby/son/internal/metrics"
)

const (
	appName    = "son"
	confEnvKey = "CONFIG"

	defaultConfigPath = "configs"
	defaultConfig     = "local"
	defaultPort       = 8086
)

func main() {
	cf, err := getConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log, err := logger.New(cf.Server.Logging)
	if err != nil {
		log.Fatal().Msgf("Could not initialize logger: %s", err)
	}
	log.Info().Msgf("log format '%s', log level '%s'", cf.Server.Logging.Format, cf.Server.Logging.Level)

	if err = setupMetricsServer(); err != nil {
		log.Fatal().Err(err).Msg("Could not setup metrics endpoints.")
	}

	sCh := listenToSignal()
	ctx, cancel := context.WithCancel(context.Background())
	go func(sCh <-chan os.Signal, cancel func()) {
		<-sCh
		cancel()
	}(sCh, cancel)

	s, err := server.New().
		WithContext(context.Background()).
		WithConfig(cf).
		WithPort(cf.Server.Port).
		WithLogger(log).
		Build()

	if err != nil {
		log.Fatal().Msgf("Could not create server: %s", err)
	}

	go func(s *http.Server) {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Could not shutdown server.")
		}
	}(s)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, s *http.Server) {
		defer wg.Done()
		log.Info().Msgf("Starting server at: %d", defaultPort)
		if err := s.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("Could not serve endpoints.")
			}
		}
	}(wg, s)

	wg.Wait()
}

func getConfig() (*config.Config, error) {
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		configPath = defaultConfigPath
	}
	if _, ok := os.LookupEnv(confEnvKey); !ok {
		if err := os.Setenv(confEnvKey, defaultConfig); err != nil {
			return nil, fmt.Errorf("could not set env: %w", err)
		}
	}

	cfg, err := config.Get(configPath, confEnvKey)
	if err != nil {
		return nil, fmt.Errorf("could not get config: %w", err)
	}

	return cfg, nil
}

func listenToSignal() <-chan os.Signal {
	const ChanBuffer = 2
	sig := make(chan os.Signal, ChanBuffer)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return sig
}

func setupMetricsServer() error {
	_, err := metrics.NewPrometheus(appName)
	return err
}

func health(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
