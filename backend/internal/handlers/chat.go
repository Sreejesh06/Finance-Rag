package handlers

import (
	"encoding/json"
	"net/http"
	"supplychain-rag/internal/models"
	"supplychain-rag/internal/services"
)

type ChatHandler struct {
	retrievalService services.RetrievalService
}

func NewChatHandler(rs services.RetrievalService) *ChatHandler {
	return &ChatHandler{retrievalService: rs}
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Question == "" {
		respondWithError(w, http.StatusBadRequest, "Question cannot be empty", nil)
		return
	}

	resp, err := h.retrievalService.Chat(r.Context(), req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to process chat", err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}
