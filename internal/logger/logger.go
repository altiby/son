package logger

import (
	"bytes"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

type Logger interface {
	Err(err error) *zerolog.Event

	Error() *zerolog.Event
	Warn() *zerolog.Event
	Info() *zerolog.Event
	Debug() *zerolog.Event
	Trace() *zerolog.Event
}

type Log struct {
	zerolog.Logger
}

func New(cfg Config) (*Log, error) {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if cfg.Format == "text" {
		log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	l := Log{log}

	err := l.SetLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *Log) SetLevel(lvl string) error {
	level, err := zerolog.ParseLevel(lvl)
	if err != nil {
		return err
	}

	l.Logger = l.Level(level)

	return nil
}

func (l *Log) NewRequestLogFormatter() *RequestLogFormatter {
	return &RequestLogFormatter{Logger: l}
}

func (l *Log) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := l.NewRequestLogFormatter().NewLogEntry(r)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		var respBody bytes.Buffer
		ww.Tee(&respBody)

		next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))

		entry.Write(
			ww.Status(),
			ww.BytesWritten(),
			ww.Header(),
			time.Since(time.Now().UTC()),
			respBody.String(),
		)
	})
}
