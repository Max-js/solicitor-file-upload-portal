import type { Client, Document } from './types'

const baseUrl = '/api'

async function parseError(res: Response): Promise<string> {
  try {
    const body = await res.json()
    return body.error ?? `HTTP ${res.status}`
  } catch {
    return `HTTP ${res.status}`
  }
}

export async function getAllDocuments(): Promise<Document[]> {
  const res = await fetch(`${baseUrl}/documents`)
  if (!res.ok) throw new Error(await parseError(res))
  return res.json()
}

export async function getDocumentsByClient(clientId: string): Promise<Document[]> {
  const res = await fetch(`${baseUrl}/documents?clientId=${clientId}`)
  if (!res.ok) throw new Error(await parseError(res))
  return res.json()
}

export async function getClientByEmail(email: string): Promise<Client> {
  const res = await fetch(`${baseUrl}/client?email=${encodeURIComponent(email)}`)
  if (!res.ok) throw new Error(await parseError(res))
  return res.json()
}

export async function uploadDocument(clientId: string, file: File): Promise<Document> {
  const form = new FormData()
  form.append('clientId', clientId)
  form.append('file', file)
  const res = await fetch(`${baseUrl}/documents`, { method: 'POST', body: form })
  if (!res.ok) throw new Error(await parseError(res))
  return res.json()
}