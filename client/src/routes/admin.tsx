import { useEffect, useState } from 'react'
import { type Document } from '../types'
import { getAllDocuments } from '../api'

export function Admin() {
  const [docs, setDocs] = useState<Document[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getAllDocuments()
      .then((list) => setDocs(list))
      .catch((error: unknown) => setError(error instanceof Error ? error.message : 'Failed to load'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <section className="page">
      <h1>Solicitor Admin</h1>
      <p className="muted">All uploaded documents.</p>

      {error && <p className="msg err">{error}</p>}
      {loading ? (
        <p className="muted">Loading…</p>
      ) : docs.length === 0 ? (
        <p className="muted">No documents uploaded yet.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>Client ID</th>
              <th>File</th>
              <th>Type</th>
              <th>Status</th>
              <th>Uploaded</th>
            </tr>
          </thead>
          <tbody>
            {docs.map((doc) => (
              <tr key={doc.id}>
                <td className="mono">{doc.clientId}</td>
                <td>{doc.filename}</td>
                <td className="muted">{doc.contentType}</td>
                <td>{doc.status}</td>
                <td>{new Date(doc.createdAt).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}