# Design Decisions

- **Embedding Model**: Default to local Ollama API if configured (e.g. `nomic-embed-text`) or a simple fallback. If no local embedding model is configured via `EMBEDDING_PROVIDER` and `EMBEDDING_MODEL`, we will rely on OpenRouter if it supports embeddings, or a simple mock/in-memory local service for demonstration if needed. Actually, we will just use the ChromaDB client's built-in Sentence Transformers if available, or just HTTP requests to an embedding endpoint.
- **Frontend Routing**: We'll use `react-router-dom` for the `/admin` and `/chat` interfaces.
- **Dependency Injection**: We will use constructor injection for our services and repositories to keep the code testable, passing in the SQLite DB connection and ChromaDB client.
- **Error Handling**: Custom error types will be mapped to HTTP status codes in a centralized error middleware or handler wrapper.
- **Text Extraction**: We'll use `github.com/ledongthuc/pdf` since it's an accessible MIT-licensed PDF parser for Go.
- **ChromaDB integration**: Will use `github.com/amikos-tech/chroma-go`.
