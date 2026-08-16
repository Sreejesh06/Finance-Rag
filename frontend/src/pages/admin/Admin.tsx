import React, { useState, useEffect, useRef } from 'react';
import { adminApi } from '../../services/api';
import type { Document, Stats } from '../../types';
import { UploadCloud, FileText, Trash2, RefreshCw, CheckCircle2, XCircle, Settings, Play } from 'lucide-react';

export default function Admin() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [documents, setDocuments] = useState<Document[]>([]);
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const loadData = async () => {
    try {
      const [statsData, docsData] = await Promise.all([
        adminApi.getStats(),
        adminApi.getDocuments()
      ]);
      setStats(statsData);
      setDocuments(docsData || []);
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 5000); // refresh every 5s
    return () => clearInterval(interval);
  }, []);

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    setUploading(true);
    for (let i = 0; i < files.length; i++) {
      try {
        await adminApi.uploadDocument(files[i]);
      } catch (err) {
        console.error("Upload failed for", files[i].name, err);
      }
    }
    setUploading(false);
    if (fileInputRef.current) fileInputRef.current.value = '';
    loadData();
  };

  const indexDoc = async (id: string) => {
    try {
      await adminApi.indexDocument(id);
      loadData();
    } catch (e) {
      console.error(e);
    }
  };

  const deleteDoc = async (id: string) => {
    if (!confirm("Are you sure?")) return;
    try {
      await adminApi.deleteDocument(id);
      loadData();
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="max-w-7xl mx-auto p-6 text-slate-800">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-slate-900 flex items-center gap-2">
          <Settings className="w-8 h-8 text-indigo-600" />
          RAG Admin Dashboard
        </h1>
        <button 
          onClick={() => fileInputRef.current?.click()}
          disabled={uploading}
          className="bg-indigo-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-indigo-700 transition"
        >
          <UploadCloud className="w-5 h-5" />
          {uploading ? 'Uploading...' : 'Upload PDF'}
        </button>
        <input 
          type="file" 
          accept=".pdf" 
          multiple 
          className="hidden" 
          ref={fileInputRef} 
          onChange={handleFileUpload} 
        />
      </div>

      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-100">
            <p className="text-sm text-slate-500">Total Documents</p>
            <p className="text-2xl font-semibold">{stats.total_documents}</p>
          </div>
          <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-100">
            <p className="text-sm text-slate-500">Indexed Documents</p>
            <p className="text-2xl font-semibold text-emerald-600">{stats.indexed_documents}</p>
          </div>
          <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-100">
            <p className="text-sm text-slate-500">Total Pages</p>
            <p className="text-2xl font-semibold text-blue-600">{stats.total_pages}</p>
          </div>
          <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-100">
            <p className="text-sm text-slate-500">Total Chunks</p>
            <p className="text-2xl font-semibold text-purple-600">{stats.total_chunks}</p>
          </div>
        </div>
      )}

      <div className="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-50 text-slate-600 text-sm border-b border-slate-100">
              <th className="p-4 font-medium">Filename</th>
              <th className="p-4 font-medium">Status</th>
              <th className="p-4 font-medium">Pages</th>
              <th className="p-4 font-medium">Chunks</th>
              <th className="p-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {documents.length === 0 ? (
              <tr>
                <td colSpan={5} className="p-8 text-center text-slate-400">No documents found.</td>
              </tr>
            ) : documents.map(doc => (
              <tr key={doc.id} className="border-b border-slate-50 hover:bg-slate-50/50">
                <td className="p-4 flex items-center gap-3">
                  <FileText className="w-5 h-5 text-slate-400" />
                  <span className="font-medium">{doc.filename}</span>
                </td>
                <td className="p-4">
                  {doc.status === 'indexed' && <span className="inline-flex items-center gap-1 text-emerald-600 text-sm bg-emerald-50 px-2 py-1 rounded-full"><CheckCircle2 className="w-4 h-4"/> Indexed</span>}
                  {doc.status === 'processing' && <span className="inline-flex items-center gap-1 text-blue-600 text-sm bg-blue-50 px-2 py-1 rounded-full"><RefreshCw className="w-4 h-4 animate-spin"/> Processing</span>}
                  {doc.status === 'uploaded' && <span className="inline-flex items-center gap-1 text-slate-600 text-sm bg-slate-100 px-2 py-1 rounded-full">Uploaded</span>}
                  {doc.status === 'failed' && <span className="inline-flex items-center gap-1 text-rose-600 text-sm bg-rose-50 px-2 py-1 rounded-full"><XCircle className="w-4 h-4"/> Failed</span>}
                </td>
                <td className="p-4 text-slate-600">{doc.page_count}</td>
                <td className="p-4 text-slate-600">{doc.chunk_count}</td>
                <td className="p-4 flex justify-end gap-2">
                  <button 
                    onClick={() => indexDoc(doc.id)} 
                    disabled={doc.status === 'indexed' || doc.status === 'processing'}
                    className={`p-2 rounded-lg transition ${doc.status === 'indexed' || doc.status === 'processing' ? 'text-slate-300' : 'text-indigo-600 hover:bg-indigo-50'}`}
                    title="Index Document"
                  >
                    <Play className="w-5 h-5" />
                  </button>
                  <button 
                    onClick={() => deleteDoc(doc.id)}
                    className="p-2 rounded-lg text-rose-500 hover:bg-rose-50 transition"
                    title="Delete Document"
                  >
                    <Trash2 className="w-5 h-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
