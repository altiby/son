package logger

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

// RequestLogFormatter is a simple logger that implements a middleware.LogFormatter.
type RequestLogFormatter struct {
	Logger *Log
}

// NewLogEntry creates a new LogEntry for the request.
func (l *RequestLogFormatter) NewLogEntry(r *http.Request) *logEntry {
	reqID := middleware.GetReqID(r.Context())
	if reqID == "" {
		reqID = r.Header.Get(middleware.RequestIDHeader)
	}

	l.Logger = &Log{l.Logger.With().Str("req_id", reqID).Logger()}

	fields := map[string]interface{}{
		"time":        time.Now().UTC().Format(time.RFC1123),
		"remote_addr": r.RemoteAddr,
		"headers":     r.Header,
		"method":      r.Method,
		"request_uri": r.RequestURI,
	}

	if l.Logger.GetLevel() <= zerolog.TraceLevel && r.Method != http.MethodGet {
		var buf bytes.Buffer
		r.Body = io.NopCloser(io.TeeReader(r.Body, &buf))
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			l.Logger.Error().Err(err).Msg("cannot read body bytes")
		}
		r.Body = io.NopCloser(&buf)

		fields["body"] = bodyBytes
	}

	entry := logEntry{
		RequestLogFormatter: l,
		reqID:               reqID,
	}

	l.Logger.Debug().Fields(fields).Msg("request started")

	return &entry
}
