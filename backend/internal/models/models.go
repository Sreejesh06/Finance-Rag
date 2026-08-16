package models

import "time"

type DocumentStatus string

const (
	StatusUploaded   DocumentStatus = "uploaded"
	StatusProcessing DocumentStatus = "processing"
	StatusIndexed    DocumentStatus = "indexed"
	StatusFailed     DocumentStatus = "failed"
)

type Document struct {
	ID           string         `json:"id"`
	Filename     string         `json:"filename"`
	FileHash     string         `json:"file_hash"`
	FileSize     int64          `json:"file_size"`
	PageCount    int            `json:"page_count"`
	Status       DocumentStatus `json:"status"`
	ChunkCount   int            `json:"chunk_count"`
	UploadedAt   time.Time      `json:"uploaded_at"`
	IndexedAt    *time.Time     `json:"indexed_at"`
	ErrorMessage string         `json:"error_message"`
}

type ChunkMetadata struct {
	DocumentID   string `json:"document_id"`
	Filename     string `json:"filename"`
	Page         int    `json:"page"`
	ChunkIndex   int    `json:"chunk_index"`
	DocumentHash string `json:"document_hash"`
}

type ChatRequest struct {
	Question string `json:"question"`
	TopK     int    `json:"top_k"`
}

type Source struct {
	Filename string `json:"filename"`
	Page     int    `json:"page"`
}

type ChatResponse struct {
	Answer  string   `json:"answer"`
	Sources []Source `json:"sources"`
}

type Stats struct {
	TotalDocuments   int `json:"total_documents"`
	IndexedDocuments int `json:"indexed_documents"`
	PendingDocuments int `json:"pending_documents"`
	FailedDocuments  int `json:"failed_documents"`
	TotalPages       int `json:"total_pages"`
	TotalChunks      int `json:"total_chunks"`
}
