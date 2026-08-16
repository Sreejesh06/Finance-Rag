import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import Admin from './pages/admin/Admin';
import Chat from './pages/chat/Chat';
import { Database, MessageSquare } from 'lucide-react';

function Navigation() {
  const location = useLocation();
  const isAdmin = location.pathname.startsWith('/admin');

  return (
    <nav className="bg-white border-b border-slate-200">
      <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
        <div className="flex items-center gap-2 font-bold text-xl text-indigo-900 tracking-tight">
          <div className="w-8 h-8 bg-indigo-600 text-white rounded-lg flex items-center justify-center">
            <Database className="w-5 h-5" />
          </div>
          SupplyChain RAG
        </div>
        <div className="flex gap-4">
          <Link
            to="/chat"
            className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition ${!isAdmin ? 'bg-indigo-50 text-indigo-700' : 'text-slate-600 hover:bg-slate-50'}`}
          >
            <MessageSquare className="w-5 h-5" />
            Chat
          </Link>
          <Link
            to="/admin"
            className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition ${isAdmin ? 'bg-indigo-50 text-indigo-700' : 'text-slate-600 hover:bg-slate-50'}`}
          >
            <Database className="w-5 h-5" />
            Admin
          </Link>
        </div>
      </div>
    </nav>
  );
}

export default function App() {
  return (
    <Router>
      <div className="min-h-screen bg-slate-50 font-sans">
        <Navigation />
        <main className="py-6">
          <Routes>
            <Route path="/admin" element={<Admin />} />
            <Route path="/chat" element={<Chat />} />
            <Route path="/" element={<Chat />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}
