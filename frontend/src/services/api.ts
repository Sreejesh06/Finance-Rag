import axios from 'axios';
import type { Document, Stats, ChatResponse } from '../types';

const API_BASE = 'http://localhost:8080/api';

const apiClient = axios.create({
  baseURL: API_BASE,
});

export const adminApi = {
  getStats: async () => {
    const response = await apiClient.get<Stats>('/admin/stats');
    return response.data;
  },
  getDocuments: async () => {
    const response = await apiClient.get<Document[]>('/admin/documents');
    return response.data;
  },
  uploadDocument: async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await apiClient.post<Document>('/admin/documents/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },
  indexDocument: async (id: string) => {
    const response = await apiClient.post(`/admin/documents/${id}/index`);
    return response.data;
  },
  deleteDocument: async (id: string) => {
    const response = await apiClient.delete(`/admin/documents/${id}`);
    return response.data;
  },
};

export const chatApi = {
  askQuestion: async (question: string) => {
    const response = await apiClient.post<ChatResponse>('/chat', { question });
    return response.data;
  },
};
