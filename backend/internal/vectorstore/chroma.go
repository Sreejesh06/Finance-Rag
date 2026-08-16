package vectorstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VectorStore interface {
	AddDocuments(ctx context.Context, chunks []string, embeddings [][]float32, metadatas []map[string]interface{}, ids []string) error
	SimilaritySearch(ctx context.Context, query string, queryEmbedding []float32, topK int) ([]string, []map[string]interface{}, error)
	DeleteByDocumentID(ctx context.Context, documentID string) error
}

type chromaStore struct {
	baseURL    string
	collection string
	httpClient *http.Client
}

func NewChromaStore(chromaURL string) (VectorStore, error) {
	store := &chromaStore{
		baseURL:    chromaURL,
		collection: "rag_documents",
		httpClient: &http.Client{},
	}
	
	// Create collection if it doesn't exist
	_ = store.createCollection()

	return store, nil
}

func (s *chromaStore) createCollection() error {
	payload := map[string]interface{}{
		"name": s.collection,
	}
	data, _ := json.Marshal(payload)
	
	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections", s.baseURL)
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *chromaStore) getCollectionID() (string, error) {
	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s", s.baseURL, s.collection)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("collection not found")
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if id, ok := res["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("no id found")
}

func (s *chromaStore) AddDocuments(ctx context.Context, chunks []string, embeddings [][]float32, metadatas []map[string]interface{}, ids []string) error {
	if len(chunks) == 0 {
		return nil
	}

	colID, err := s.getCollectionID()
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"ids":       ids,
		"documents": chunks,
		"embeddings": embeddings,
		"metadatas": metadatas,
	}
	data, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s/add", s.baseURL, colID)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add documents, status: %d, %s", resp.StatusCode, string(body))
	}
	return nil
}

func (s *chromaStore) SimilaritySearch(ctx context.Context, query string, queryEmbedding []float32, topK int) ([]string, []map[string]interface{}, error) {
	colID, err := s.getCollectionID()
	if err != nil {
		return nil, nil, err
	}

	payload := map[string]interface{}{
		"query_texts":      []string{query},
		"query_embeddings": [][]float32{queryEmbedding},
		"n_results":        topK,
	}
	data, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s/query", s.baseURL, colID)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("search failed, status: %d, %s", resp.StatusCode, string(body))
	}

	var result struct {
		Documents [][]string                   `json:"documents"`
		Metadatas [][]map[string]interface{} `json:"metadatas"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	if len(result.Documents) == 0 || len(result.Documents[0]) == 0 {
		return nil, nil, nil
	}

	return result.Documents[0], result.Metadatas[0], nil
}

func (s *chromaStore) DeleteByDocumentID(ctx context.Context, documentID string) error {
	colID, err := s.getCollectionID()
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"where": map[string]interface{}{
			"document_id": documentID,
		},
	}
	data, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s/delete", s.baseURL, colID)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
