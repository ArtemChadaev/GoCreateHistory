package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
)

func (h *Handler) createHistory(w http.ResponseWriter, r *http.Request) {
	var input domain.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {

		logger.Error(r.Context(), "Failed to decode request body", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if uid, ok := r.Context().Value("user_id").(int); ok {
		input.UserID = uid
	} else {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	//TODO: валидацию полей сделать

	hID, err := h.service.Create(r.Context(), input)
	if err != nil {
		logger.Error(r.Context(), "Failed to create user", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"uuid": hID.String()})
}
