package logger

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		cfg Config
		err error
	}{
		"no error": {
			Config{Level: "debug"}, nil,
		},
		"invalid level": {
			Config{Level: "test"}, errors.New("Unknown Level String: 'test', defaulting to NoLevel"),
		},
	}

	for k, v := range cases {
		v := v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			l, err := New(v.cfg)
			assert.Equal(t, v.err, err)

			if v.err != nil {
				return
			}

			lvl, _ := zerolog.ParseLevel(v.cfg.Level)
			log := zerolog.New(os.Stderr).Level(lvl).With().Timestamp().Logger()

			assert.Equal(t, log, l.Logger)
		})
	}
}

func TestLogSetLevel(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		cfg Config
		err error
	}{
		"no error": {
			Config{Level: "debug"}, nil,
		},
		"invalid level": {
			Config{Level: "test"}, errors.New("Unknown Level String: 'test', defaulting to NoLevel"),
		},
	}

	for k, v := range cases {
		v := v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			l, err := New(v.cfg)

			assert.Equal(t, v.err, err)
			if v.err == nil {
				assert.Equal(t, v.cfg.Level, l.GetLevel().String())
			}
		})
	}
}

func TestLogLogMiddleware(t *testing.T) {
	t.Parallel()
	isCalled := false

	l, err := New(Config{Level: "debug"})
	assert.Nil(t, err)

	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chi.NewRouteContext())
	req := httptest.
		NewRequest("GET", "http://testing", nil).
		WithContext(ctx)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		le := middleware.GetLogEntry(r)
		assert.NotNil(t, le)

		isCalled = true
	})

	l.LogMiddleware(nextHandler).ServeHTTP(httptest.NewRecorder(), req)
	assert.True(t, isCalled)
}
