package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
)

func ResponseWithError(logger *slog.Logger, w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	resp := errorResponse{
		Success: false,
		Error:   msg,
		Code:    code,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("failed on encoding json error: " + err.Error())
	}
}

func Response(logger *slog.Logger, w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if data != nil {
		// Обработка пустых слайсов
		if reflect.ValueOf(data).Kind() == reflect.Slice && reflect.ValueOf(data).Len() == 0 {
			data = []any{}
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Error("failed on encoding http response: " + err.Error())
		}
	} else {
		// Для методов без возвращаемых данных (например, DELETE)
		if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
			logger.Error("failed on encoding http response: " + err.Error())
		}
	}
}

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}
