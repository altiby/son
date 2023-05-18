package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type logEntry struct {
	*RequestLogFormatter
	reqID string
}

func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, body interface{}) {
	fields := map[string]interface{}{
		"status":  status,
		"elapsed": elapsed.String(),
		"time":    time.Now().UTC().Format(time.RFC1123),
		"bytes":   bytes,
	}

	if l.Logger.GetLevel() >= zerolog.TraceLevel {
		fields["body"] = body
	}

	switch {
	case status < 500:
		l.Logger.Debug().Fields(fields).Msg("request completed")
	default:
		l.Logger.Error().Fields(fields).Msg("request error")
	}
}

func (l *logEntry) Panic(v interface{}, stack []byte) {
	l.Logger.Warn().Fields(map[string]interface{}{
		"err":   fmt.Sprintf("%#v", v),
		"stack": string(stack),
	}).Msg("Fatal error")
}
