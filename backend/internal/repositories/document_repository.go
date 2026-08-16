package repositories

import (
	"database/sql"
	"supplychain-rag/internal/models"
)

type DocumentRepository interface {
	Create(doc *models.Document) error
	GetByID(id string) (*models.Document, error)
	GetByHash(hash string) (*models.Document, error)
	UpdateStatus(id string, status models.DocumentStatus, errorMsg string) error
	UpdateIndexStats(id string, chunkCount int, pageCount int) error
	List() ([]models.Document, error)
	Delete(id string) error
	GetStats() (*models.Stats, error)
}

type sqliteDocRepo struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) DocumentRepository {
	return &sqliteDocRepo{db: db}
}

func (r *sqliteDocRepo) Create(doc *models.Document) error {
	query := `INSERT INTO documents (id, filename, file_hash, file_size, status, uploaded_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, doc.ID, doc.Filename, doc.FileHash, doc.FileSize, doc.Status, doc.UploadedAt)
	return err
}

func (r *sqliteDocRepo) GetByID(id string) (*models.Document, error) {
	query := `SELECT id, filename, file_hash, file_size, page_count, status, chunk_count, uploaded_at, indexed_at, error_message FROM documents WHERE id = ?`
	row := r.db.QueryRow(query, id)
	return scanDoc(row)
}

func (r *sqliteDocRepo) GetByHash(hash string) (*models.Document, error) {
	query := `SELECT id, filename, file_hash, file_size, page_count, status, chunk_count, uploaded_at, indexed_at, error_message FROM documents WHERE file_hash = ?`
	row := r.db.QueryRow(query, hash)
	return scanDoc(row)
}

func (r *sqliteDocRepo) UpdateStatus(id string, status models.DocumentStatus, errorMsg string) error {
	query := `UPDATE documents SET status = ?, error_message = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, errorMsg, id)
	return err
}

func (r *sqliteDocRepo) UpdateIndexStats(id string, chunkCount int, pageCount int) error {
	query := `UPDATE documents SET chunk_count = ?, page_count = ?, status = ?, indexed_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, chunkCount, pageCount, models.StatusIndexed, id)
	return err
}

func (r *sqliteDocRepo) List() ([]models.Document, error) {
	query := `SELECT id, filename, file_hash, file_size, page_count, status, chunk_count, uploaded_at, indexed_at, error_message FROM documents ORDER BY uploaded_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := make([]models.Document, 0)
	for rows.Next() {
		doc, err := scanDocRow(rows)
		if err != nil {
			return nil, err
		}
		docs = append(docs, *doc)
	}
	return docs, nil
}

func (r *sqliteDocRepo) Delete(id string) error {
	query := `DELETE FROM documents WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *sqliteDocRepo) GetStats() (*models.Stats, error) {
	stats := &models.Stats{}

	// Total docs
	r.db.QueryRow(`SELECT COUNT(*) FROM documents`).Scan(&stats.TotalDocuments)
	r.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE status = ?`, models.StatusIndexed).Scan(&stats.IndexedDocuments)
	r.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE status = ? OR status = ?`, models.StatusUploaded, models.StatusProcessing).Scan(&stats.PendingDocuments)
	r.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE status = ?`, models.StatusFailed).Scan(&stats.FailedDocuments)

	// Total pages and chunks
	r.db.QueryRow(`SELECT COALESCE(SUM(page_count), 0) FROM documents`).Scan(&stats.TotalPages)
	r.db.QueryRow(`SELECT COALESCE(SUM(chunk_count), 0) FROM documents`).Scan(&stats.TotalChunks)

	return stats, nil
}

func scanDoc(row *sql.Row) (*models.Document, error) {
	var doc models.Document
	var indexedAt sql.NullTime
	var errorMsg sql.NullString
	err := row.Scan(&doc.ID, &doc.Filename, &doc.FileHash, &doc.FileSize, &doc.PageCount, &doc.Status, &doc.ChunkCount, &doc.UploadedAt, &indexedAt, &errorMsg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if indexedAt.Valid {
		doc.IndexedAt = &indexedAt.Time
	}
	if errorMsg.Valid {
		doc.ErrorMessage = errorMsg.String
	}
	return &doc, nil
}

func scanDocRow(rows *sql.Rows) (*models.Document, error) {
	var doc models.Document
	var indexedAt sql.NullTime
	var errorMsg sql.NullString
	err := rows.Scan(&doc.ID, &doc.Filename, &doc.FileHash, &doc.FileSize, &doc.PageCount, &doc.Status, &doc.ChunkCount, &doc.UploadedAt, &indexedAt, &errorMsg)
	if err != nil {
		return nil, err
	}
	if indexedAt.Valid {
		doc.IndexedAt = &indexedAt.Time
	}
	if errorMsg.Valid {
		doc.ErrorMessage = errorMsg.String
	}
	return &doc, nil
}
