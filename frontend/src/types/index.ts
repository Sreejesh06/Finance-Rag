export interface Document {
  id: string;
  filename: string;
  file_hash: string;
  file_size: number;
  page_count: number;
  status: 'uploaded' | 'processing' | 'indexed' | 'failed';
  chunk_count: number;
  uploaded_at: string;
  indexed_at: string | null;
  error_message: string;
}

export interface Stats {
  total_documents: number;
  indexed_documents: number;
  pending_documents: number;
  failed_documents: number;
  total_pages: number;
  total_chunks: number;
}

export interface Source {
  filename: string;
  page: number;
}

export interface ChatResponse {
  answer: string;
  sources: Source[];
}
