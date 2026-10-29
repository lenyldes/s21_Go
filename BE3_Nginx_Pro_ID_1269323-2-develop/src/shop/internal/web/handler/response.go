package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"shop/internal/web/dto"
)

func writeError(w http.ResponseWriter, err error, status int) {
	response := dto.ErrorResponse{Errors: err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа с ошибкой: %v", err)
	}
}
