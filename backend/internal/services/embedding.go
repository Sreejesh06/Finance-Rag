package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"supplychain-rag/internal/config"
)

type EmbeddingService interface {
	EmbedTexts(ctx context.Context, texts []string) ([][]float32, error)
}

type openRouterEmbedding struct {
	config     *config.Config
	httpClient *http.Client
}

func NewEmbeddingService(cfg *config.Config) EmbeddingService {
	return &openRouterEmbedding{
		config:     cfg,
		httpClient: &http.Client{},
	}
}

type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (s *openRouterEmbedding) EmbedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	model := s.config.EmbeddingModel
	if model == "all-MiniLM-L6-v2" || model == "" {
		model = "openai/text-embedding-3-small"
	}

	reqBody := embeddingRequest{
		Model: model,
		Input: texts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.OpenRouterAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openrouter embedding error (status %d): %s", resp.StatusCode, string(body))
	}

	var orResp embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(orResp.Data))
	for i, data := range orResp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}
