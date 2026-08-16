package services

import (
	"context"
	"fmt"
	"supplychain-rag/internal/config"
	"supplychain-rag/internal/models"
	"supplychain-rag/internal/rag"
	"supplychain-rag/internal/repositories"
	"supplychain-rag/internal/utils"
	"supplychain-rag/internal/vectorstore"
)

type IngestionService interface {
	ProcessDocument(ctx context.Context, doc *models.Document, filePath string) error
}

type ingestionService struct {
	config *config.Config
	repo   repositories.DocumentRepository
	vs     vectorstore.VectorStore
	embed  EmbeddingService
}

func NewIngestionService(cfg *config.Config, repo repositories.DocumentRepository, vs vectorstore.VectorStore, embed EmbeddingService) IngestionService {
	return &ingestionService{
		config: cfg,
		repo:   repo,
		vs:     vs,
		embed:  embed,
	}
}

func (s *ingestionService) ProcessDocument(ctx context.Context, doc *models.Document, filePath string) error {
	// Update status
	_ = s.repo.UpdateStatus(doc.ID, models.StatusProcessing, "")

	// 1. Extract text
	pages, err := utils.ExtractTextFromPDF(filePath)
	if err != nil {
		_ = s.repo.UpdateStatus(doc.ID, models.StatusFailed, err.Error())
		return err
	}

	if len(pages) == 0 {
		err := fmt.Errorf("no extractable text found in PDF")
		_ = s.repo.UpdateStatus(doc.ID, models.StatusFailed, err.Error())
		return err
	}

	// 2. Chunk text
	chunks := rag.ChunkText(pages, s.config.ChunkSize, s.config.ChunkOverlap)
	if len(chunks) == 0 {
		err := fmt.Errorf("no chunks generated from text")
		_ = s.repo.UpdateStatus(doc.ID, models.StatusFailed, err.Error())
		return err
	}

	// 3. Prepare for ChromaDB
	var texts []string
	var metadatas []map[string]interface{}
	var ids []string

	for _, chunk := range chunks {
		texts = append(texts, chunk.Text)

		id := fmt.Sprintf("%s-%d-%d", doc.FileHash, chunk.PageNumber, chunk.ChunkIndex)
		ids = append(ids, id)

		metadatas = append(metadatas, map[string]interface{}{
			"document_id":   doc.ID,
			"filename":      doc.Filename,
			"page":          chunk.PageNumber,
			"chunk_index":   chunk.ChunkIndex,
			"document_hash": doc.FileHash,
		})
	}

	embeddings, err := s.embed.EmbedTexts(ctx, texts)
	if err != nil {
		_ = s.repo.UpdateStatus(doc.ID, models.StatusFailed, fmt.Sprintf("embedding error: %v", err))
		return err
	}

	// 4. Store in Vector DB
	if err := s.vs.AddDocuments(ctx, texts, embeddings, metadatas, ids); err != nil {
		_ = s.repo.UpdateStatus(doc.ID, models.StatusFailed, fmt.Sprintf("vectordb error: %v", err))
		return err
	}

	// 5. Update SQLite
	if err := s.repo.UpdateIndexStats(doc.ID, len(chunks), len(pages)); err != nil {
		return err
	}

	return nil
}
