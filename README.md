# Supply Chain / Financial Reports RAG System

An end-to-end autonomous RAG (Retrieval-Augmented Generation) system built with Go and React, specifically designed for querying large financial and supply chain PDF reports.

## What Has Been Built

### Backend (Go / Chi)
- **RESTful API**: Handlers to upload documents, list documents, retrieve metrics, and chat with the AI.
- **SQLite Database**: Persists document metadata (filename, page count, chunk count, file size, status, and upload/indexing timestamps).
- **ChromaDB Vector Store (v2 API)**: Automatically connects to a local ChromaDB instance using the `v2` REST API. Handles embedding storage and similarity search.
- **PDF Processing**: Extracts text page-by-page from PDFs and splits them into manageable chunks (default 1200 words, 150 overlap) with page-number metadata tracking.
- **OpenRouter Integration**:
  - Uses `openai/text-embedding-3-small` (or configured model) to embed document chunks before saving them to ChromaDB.
  - Uses `openai/gpt-4o-mini` (or configured model) to generate context-grounded answers to user queries, strictly rejecting out-of-context requests.
- **Resilient Pipeline**: Tracks ingestion state (`uploaded`, `processing`, `indexed`, `failed`) and stores failure reasons (e.g., if OpenRouter auth fails) directly into the SQLite DB.

### Frontend (React / TypeScript / Vite / Tailwind)
- **Admin Dashboard** (`/admin`): 
  - Drag-and-drop or select PDF files to upload.
  - View overall system stats (Total Docs, Indexed Docs, Pages, Chunks).
  - Manage uploaded files (Index, Delete) and monitor their statuses.
- **Chat Interface** (`/chat`):
  - A sleek conversational UI to query your supply chain and financial data.
  - Displays accurate page-level citations for every piece of context retrieved from the database.
- **Modern UI**: Fully responsive, accessible, and polished using Lucide React icons and Tailwind CSS.

---

## API Endpoints

### Admin & Documents
- `GET /api/admin/stats` - Returns overall statistics (totals of documents, pages, chunks).
- `GET /api/admin/documents` - Lists all uploaded documents and their ingestion status.
- `GET /api/admin/documents/{id}` - Gets details for a specific document.
- `POST /api/admin/documents/upload` - Multipart form upload endpoint for PDF files.
- `POST /api/admin/documents/{id}/index` - Triggers the background ingestion/chunking/embedding process for a file.
- `DELETE /api/admin/documents/{id}` - Deletes a document from SQLite, ChromaDB, and the disk.

### Chat & Retrieval
- `POST /api/chat` - Submits a user query.
  - **Body**: `{"question": "What was Tesla's revenue in Q1?"}`
  - **Response**: `{"answer": "...", "sources": [{"filename": "...", "page": 3}]}`

---

## Setup Instructions & How to Run

### Prerequisites
- **Go** (1.21+)
- **Node.js** (18+)
- **ChromaDB** running locally. 
  - *Using Docker*: `docker run -p 8000:8000 chromadb/chroma`
  - *Using Python (native)*: `chroma run --path ./my_vector_data`

### 1. Backend Setup
Navigate into the backend directory:
```bash
cd backend
```
Copy the example environment file and configure it:
```bash
cp .env.example .env
```
Open `.env` and add your OpenRouter API Key:
```env
OPENROUTER_API_KEY=sk-or-v1-...
OPENROUTER_MODEL=openai/gpt-4o-mini
```
Run the Go server:
```bash
go run ./cmd/server
```
*(The server will start on `http://localhost:8080`)*

### 2. Frontend Setup
Open a new terminal window and navigate to the frontend directory:
```bash
cd frontend
```
Install dependencies:
```bash
npm install
```
Start the Vite development server:
```bash
npm run dev
```
*(The frontend will be available at `http://localhost:5173`)*

---

## How to Use the System
1. Make sure ChromaDB, your Go Backend, and your React Frontend are all running.
2. Open `http://localhost:5173/admin` in your browser.
3. Click **Upload PDF** and select a supply chain report (e.g., a Quarterly Update PDF).
4. Wait for it to upload, then click the **Play** (Index) icon next to the document in the list.
5. Once the status turns green (`Indexed`), switch to the **Chat** tab (`http://localhost:5173/chat`).
6. Ask a question about the document! The assistant will perform a similarity search and reply with the answer and page citations.

## Known Limitations
- Simple text extraction only (complex table extraction logic is not fully optimized).
- Recursive chunking is based strictly on chunk sizes rather than semantic boundaries.
- PDF extraction relies on simple parsing without advanced OCR capabilities.
