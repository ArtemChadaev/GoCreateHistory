package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var input domain.Auth

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {

		logger.Error(r.Context(), "Failed to decode request body", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Вообще надо бы сделать валидацию типо изза чего не удоалось создать но влом
	if err := h.service.CreateUser(r.Context(), input.Email, input.Password); err != nil {
		logger.Error(r.Context(), "Failed to create user", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := h.service.GenerateToken(r.Context(), input.Email, input.Password)

	if err != nil {
		logger.Error(r.Context(), "Failed to generate token", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var input domain.Auth

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {

		logger.Error(r.Context(), "Failed to decode request body", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.service.GenerateToken(r.Context(), input.Email, input.Password)

	if err != nil {
		logger.Error(r.Context(), "Failed to generate token", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
