package protocol

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Error struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

func WriteError(w http.ResponseWriter, requestID string, code int, err string) {
	w.Header().Add("Content-Type", "application/json")
	if code == http.StatusInternalServerError {
		log.Err(errors.New(err))
	}

	w.WriteHeader(code)
	bytes, _ := json.Marshal(Error{
		Message:   err,
		RequestID: requestID,
		Code:      code,
	})
	w.Write(bytes)
}

func WriteOk(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(data)
	w.Write(bytes)
}
