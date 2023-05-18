package logger

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNewLogEntry(t *testing.T) {
	cases := map[string]struct {
		level string
		reqID string
	}{
		"debug level": {
			"debug", "123",
		},
		"trace level": {
			"trace", "123",
		},
		"without redID": {
			"trace", "",
		},
	}

	for k, v := range cases {
		v := v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			logger, err := New(Config{Level: v.level})
			assert.Nil(t, err)

			formatter := RequestLogFormatter{logger}

			r := &http.Request{
				Body: io.NopCloser(strings.NewReader("")),
			}
			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, v.reqID))

			entry := formatter.NewLogEntry(r)
			assert.Equal(t, v.reqID, entry.reqID)
		})
	}
}
