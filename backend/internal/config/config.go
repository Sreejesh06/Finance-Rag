package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenRouterAPIKey   string
	OpenRouterModel    string
	EmbeddingProvider  string
	EmbeddingModel     string
	ChromaURL          string
	DatabasePath       string
	UploadDir          string
	ChunkSize          int
	ChunkOverlap       int
	DefaultTopK        int
	RetrievalThreshold float64
	ServerPort         string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // ignore error if .env doesn't exist

	return &Config{
		OpenRouterAPIKey:   getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterModel:    getEnv("OPENROUTER_MODEL", "openai/gpt-4o-mini"),
		EmbeddingProvider:  getEnv("EMBEDDING_PROVIDER", "chroma"),
		EmbeddingModel:     getEnv("EMBEDDING_MODEL", "all-MiniLM-L6-v2"),
		ChromaURL:          getEnv("CHROMA_URL", "http://localhost:8000"),
		DatabasePath:       getEnv("DATABASE_PATH", "./rag.db"),
		UploadDir:          getEnv("UPLOAD_DIR", "./uploads"),
		ChunkSize:          getEnvAsInt("CHUNK_SIZE", 1200),
		ChunkOverlap:       getEnvAsInt("CHUNK_OVERLAP", 150),
		DefaultTopK:        getEnvAsInt("DEFAULT_TOP_K", 5),
		RetrievalThreshold: getEnvAsFloat("RETRIEVAL_THRESHOLD", 0.7),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, fallback int) int {
	valStr := getEnv(name, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

func getEnvAsFloat(name string, fallback float64) float64 {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return fallback
}
