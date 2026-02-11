package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemChadaev/GoCreateHistory/internal/domain"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *Handler) freezeHistory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Frozen bool `json:"frozen"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Freeze(r.Context(), id, input.Frozen); err != nil {
		logger.Error(r.Context(), "Failed to freeze history", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteHistory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		logger.Error(r.Context(), "Failed to delete history", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
