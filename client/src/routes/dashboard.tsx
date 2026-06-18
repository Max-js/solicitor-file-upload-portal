import { useEffect, useRef, useState } from 'react'
import { currentUserEmail, type Client, type Document } from '../types'
import { getClientByEmail, getDocumentsByClient, uploadDocument } from '../api'

export function Dashboard() {
  const [client, setClient] = useState<Client | null>(null)
  const [file, setFile] = useState<File | null>(null)
  const [docs, setDocs] = useState<Document[]>([])
  const [busy, setBusy] = useState(false)
  const [msg, setMsg] = useState<{ kind: 'ok' | 'err'; text: string } | null>(null)
  const fileRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    getClientByEmail(currentUserEmail)
      .then(setClient)
      .catch((e: unknown) =>
        setMsg({ kind: 'err', text: e instanceof Error ? e.message : 'Could not load client' }))
  }, [])

  useEffect(() => {
    if (!client) return
    getDocumentsByClient(client.id).then(setDocs).catch(() => {})
  }, [client])

  const onSubmit = async (event: React.SubmitEvent) => {
    event.preventDefault()
    if (!file || !client) return
    setBusy(true)
    setMsg(null)
    try {
      await uploadDocument(client.id, file)
      setMsg({ kind: 'ok', text: 'Uploaded.' })
      setFile(null)
      if (fileRef.current) fileRef.current.value = ''
      setDocs(await getDocumentsByClient(client.id))
    } catch (err) {
      setMsg({ kind: 'err', text: err instanceof Error ? err.message : 'Upload failed' })
    } finally {
      setBusy(false)
    }
  }

  return (
    <section className="page">
      <h1>Dashboard</h1>

      <label className="field">
        <span>Logged in as: {client?.name ?? ''}</span>
      </label>

      <form className="upload" onSubmit={onSubmit}>
        <h2>Upload a document</h2>
        <input
          ref={fileRef}
          type="file"
          accept="application/pdf,image/jpeg,image/png"
          onChange={(e) => setFile(e.target.files?.[0] ?? null)}
        />
        <button
          className="upload-btn"
          type="submit"
          disabled={!file || busy}
        >
          {busy ? 'Uploading…' : 'Upload'}
        </button>
        {msg && <p className={`msg ${msg.kind}`}>{msg.text}</p>}
      </form>

      <section>
        <h2>Your documents</h2>
          {docs.length === 0 ? (
            <p className="muted">No documents yet.</p>
          ) : (
          <table>
            <thead>
              <tr>
                <th>File</th>
                <th>Status</th>
                <th>Uploaded</th>
              </tr>
            </thead>
            <tbody>
              {docs.map((doc) => (
                <tr key={doc.id}>
                  <td>{doc.filename}</td>
                  <td>{doc.status}</td>
                  <td>{new Date(doc.createdAt).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </section>
  )
}