import { Link, Navigate, Route, Routes } from 'react-router'
import { Dashboard } from './routes/dashboard'
import { Admin } from './routes/admin'
import './App.css'

export default function App() {
  return (
    <div className="app">
      <header className="nav">
        <span className="brand">Solicitor Account Portal</span>
        <nav>
          <Link to="/dashboard">Dashboard</Link>
          <Link to="/admin">Admin</Link>
        </nav>
      </header>
      <main>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/admin" element={<Admin />} />
          <Route path="*" element={<p>Not found</p>} />
        </Routes>
      </main>
    </div>
  )
}