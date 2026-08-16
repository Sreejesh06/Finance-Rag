package services

import (
	"context"
	"fmt"
	"strings"
	"supplychain-rag/internal/config"
	"supplychain-rag/internal/models"
	"supplychain-rag/internal/vectorstore"
)

type RetrievalService interface {
	Chat(ctx context.Context, req models.ChatRequest) (models.ChatResponse, error)
}

type retrievalService struct {
	config *config.Config
	vs     vectorstore.VectorStore
	llm    LLMService
	embed  EmbeddingService
}

func NewRetrievalService(cfg *config.Config, vs vectorstore.VectorStore, llm LLMService, embed EmbeddingService) RetrievalService {
	return &retrievalService{
		config: cfg,
		vs:     vs,
		llm:    llm,
		embed:  embed,
	}
}

func (s *retrievalService) Chat(ctx context.Context, req models.ChatRequest) (models.ChatResponse, error) {
	if req.TopK <= 0 {
		req.TopK = s.config.DefaultTopK
	}

	embeddings, err := s.embed.EmbedTexts(ctx, []string{req.Question})
	if err != nil || len(embeddings) == 0 {
		return models.ChatResponse{}, fmt.Errorf("embedding failed: %w", err)
	}

	texts, metadatas, err := s.vs.SimilaritySearch(ctx, req.Question, embeddings[0], req.TopK)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("search failed: %w", err)
	}

	if len(texts) == 0 {
		return models.ChatResponse{
			Answer:  "I could not find this information in the uploaded documents.",
			Sources: []models.Source{},
		}, nil
	}

	var contextBuilder strings.Builder
	var sources []models.Source
	sourceSet := make(map[string]bool)

	for i, text := range texts {
		meta := metadatas[i]

		filename := meta["filename"].(string)
		page := int(meta["page"].(float64)) // JSON parses numbers as float64 usually, or Chroma returns them as such

		contextBuilder.WriteString(fmt.Sprintf("--- Source: %s (Page %d) ---\n%s\n\n", filename, page, text))

		srcKey := fmt.Sprintf("%s:%d", filename, page)
		if !sourceSet[srcKey] {
			sourceSet[srcKey] = true
			sources = append(sources, models.Source{
				Filename: filename,
				Page:     page,
			})
		}
	}

	systemPrompt := `You are an expert assistant for a Supply Chain / Financial Reports system.
Answer the user's question STRICTLY using only the information in the provided context.
DO NOT use outside knowledge. DO NOT guess, fabricate, or invent numbers, dates, or citations.
Preserve units and values exactly as they appear in the text.
If the provided context does not contain the answer, reply exactly with: "I could not find this information in the uploaded documents."

Context:
` + contextBuilder.String()

	answer, err := s.llm.GenerateAnswer(ctx, systemPrompt, req.Question)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("llm failed: %w", err)
	}

	return models.ChatResponse{
		Answer:  answer,
		Sources: sources,
	}, nil
}
