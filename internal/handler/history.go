package handler

import (
	"encoding/json"
	"net/http"


	"github.com/ArtemChadaev/GoCreateHistory/internal/handler/dto"
	"github.com/ArtemChadaev/GoCreateHistory/pkg/logger"
	"github.com/creasty/defaults"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (h *Handler) createHistory(w http.ResponseWriter, r *http.Request) {
	// UserID гарантируется middleware.auth
	uid := r.Context().Value("user_id").(int)

	var req dto.CreateHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(r.Context(), "Failed to decode request body", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := defaults.Set(&req); err != nil {
		logger.Error(r.Context(), "Failed to set defaults", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := req.ToDomain(uid)

	hID, err := h.service.Create(r.Context(), input)
	if err != nil {
		logger.Error(r.Context(), "Failed to create history", err)
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
