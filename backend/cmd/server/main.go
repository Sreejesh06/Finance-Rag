package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"supplychain-rag/internal/config"
	"supplychain-rag/internal/database"
	"supplychain-rag/internal/handlers"
	"supplychain-rag/internal/repositories"
	"supplychain-rag/internal/services"
	"supplychain-rag/internal/vectorstore"
)

func main() {
	cfg := config.LoadConfig()

	// Ensure upload directory exists
	if err := os.MkdirAll(cfg.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Init SQLite
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Init ChromaDB
	vs, err := vectorstore.NewChromaStore(cfg.ChromaURL)
	if err != nil {
		log.Fatalf("Failed to connect to ChromaDB: %v", err)
	}

	// Init Repositories & Services
	repo := repositories.NewDocumentRepository(db)
	llm := services.NewLLMService(cfg)
	embed := services.NewEmbeddingService(cfg)
	ingestionService := services.NewIngestionService(cfg, repo, vs, embed)
	retrievalService := services.NewRetrievalService(cfg, vs, llm, embed)

	// Init Handlers
	adminHandler := handlers.NewAdminHandler(cfg, repo, ingestionService, vs)
	chatHandler := handlers.NewChatHandler(retrievalService)

	// Router setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/admin", func(r chi.Router) {
		r.Post("/documents/upload", adminHandler.Upload)
		r.Post("/documents/{id}/index", adminHandler.IndexDocument)
		r.Get("/documents", adminHandler.ListDocuments)
		r.Get("/documents/{id}", adminHandler.GetDocument)
		r.Delete("/documents/{id}", adminHandler.DeleteDocument)
		r.Get("/stats", adminHandler.GetStats)
	})

	r.Route("/api/chat", func(r chi.Router) {
		r.Post("/", chatHandler.Chat)
	})

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
