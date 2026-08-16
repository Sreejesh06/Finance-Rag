# Project TODO

## Phase 1: Project Initialization & Architecture Specs (In Progress)
- [x] Create project directories
- [x] Init Go module
- [ ] Init React/Vite project (In progress)
- [x] Create `TODO.md`
- [ ] Create `docs/ARCHITECTURE.md`
- [ ] Create `DECISIONS.md`

## Phase 2: Backend Foundation
- [ ] Implement `internal/config`
- [ ] Implement `internal/models`
- [ ] Setup SQLite database and create `documents` table
- [ ] Setup ChromaDB client integration

## Phase 3: Core RAG Services
- [ ] Implement PDF extraction service
- [ ] Implement text chunking logic
- [ ] Implement embedding generation
- [ ] Implement ingestion orchestration

## Phase 4: LLM Integration & Q&A Logic
- [ ] Implement OpenRouter client
- [ ] Implement vector similarity search
- [ ] Implement Q&A retrieval service with context grounding

## Phase 5: API Layer
- [ ] Implement Admin HTTP handlers
- [ ] Implement Chat HTTP handler
- [ ] Setup Chi router

## Phase 6: Frontend Development
- [ ] Configure Tailwind CSS
- [ ] Build Admin dashboard
- [ ] Build Chat UI
- [ ] Integrate frontend with backend API

## Phase 7: Testing & Finalization
- [ ] Write Go unit tests
- [ ] Format and lint code
- [ ] End-to-end testing
- [ ] Finalize `README.md`
