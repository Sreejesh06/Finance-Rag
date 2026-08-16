package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"supplychain-rag/internal/config"
	"supplychain-rag/internal/models"
	"supplychain-rag/internal/repositories"
	"supplychain-rag/internal/services"
	"supplychain-rag/internal/utils"
	"supplychain-rag/internal/vectorstore"
)

type AdminHandler struct {
	config           *config.Config
	repo             repositories.DocumentRepository
	ingestionService services.IngestionService
	vs               vectorstore.VectorStore
}

func NewAdminHandler(cfg *config.Config, repo repositories.DocumentRepository, is services.IngestionService, vs vectorstore.VectorStore) *AdminHandler {
	return &AdminHandler{
		config:           cfg,
		repo:             repo,
		ingestionService: is,
		vs:               vs,
	}
}

func (h *AdminHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read file", err)
		return
	}
	defer file.Close()

	if filepath.Ext(header.Filename) != ".pdf" {
		respondWithError(w, http.StatusBadRequest, "Only PDF files are allowed", nil)
		return
	}

	id := uuid.New().String()
	safeFilename := fmt.Sprintf("%s_%s", id, header.Filename)
	destPath := filepath.Join(h.config.UploadDir, safeFilename)

	destFile, err := os.Create(destPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save file", err)
		return
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to copy file", err)
		return
	}

	fileHash, err := utils.HashFile(destPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash file", err)
		return
	}

	existingDoc, err := h.repo.GetByHash(fileHash)
	if err == nil && existingDoc != nil {
		os.Remove(destPath)
		respondWithError(w, http.StatusConflict, "Document already exists", nil)
		return
	}

	doc := &models.Document{
		ID:         id,
		Filename:   header.Filename,
		FileHash:   fileHash,
		FileSize:   header.Size,
		Status:     models.StatusUploaded,
		UploadedAt: time.Now(),
	}

	if err := h.repo.Create(doc); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create document record", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, doc)
}

func (h *AdminHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := h.repo.List()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list documents", err)
		return
	}
	respondWithJSON(w, http.StatusOK, docs)
}

func (h *AdminHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := h.repo.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get document", err)
		return
	}
	if doc == nil {
		respondWithError(w, http.StatusNotFound, "Document not found", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, doc)
}

func (h *AdminHandler) IndexDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := h.repo.GetByID(id)
	if err != nil || doc == nil {
		respondWithError(w, http.StatusNotFound, "Document not found", nil)
		return
	}

	if doc.Status == models.StatusIndexed {
		respondWithError(w, http.StatusBadRequest, "Document already indexed", nil)
		return
	}

	safeFilename := fmt.Sprintf("%s_%s", doc.ID, doc.Filename)
	filePath := filepath.Join(h.config.UploadDir, safeFilename)

	go func() {
		_ = h.ingestionService.ProcessDocument(context.Background(), doc, filePath)
	}()

	respondWithJSON(w, http.StatusAccepted, map[string]string{"message": "Indexing started"})
}

func (h *AdminHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := h.repo.GetByID(id)
	if err != nil || doc == nil {
		respondWithError(w, http.StatusNotFound, "Document not found", nil)
		return
	}

	// Delete from ChromaDB
	_ = h.vs.DeleteByDocumentID(r.Context(), id)

	// Delete from SQLite
	_ = h.repo.Delete(id)

	// Delete from disk
	safeFilename := fmt.Sprintf("%s_%s", doc.ID, doc.Filename)
	_ = os.Remove(filepath.Join(h.config.UploadDir, safeFilename))

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Document deleted"})
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetStats()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get stats", err)
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}

// Helpers

func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	respondWithJSON(w, code, map[string]interface{}{
		"error": map[string]string{
			"code":    http.StatusText(code),
			"message": message,
			"details": errStr,
		},
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
