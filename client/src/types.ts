export type Status = 'pending' | 'verified' | 'rejected'

export interface Document {
  id: string
  clientId: string
  filename: string
  contentType: string
  sizeBytes: number
  checksumSha256: string
  status: Status
  rejectionReason?: string
  createdAt: string
  updatedAt: string
}

export interface Client {
  id: string,
  name: string,
  email: string
}

export const currentUserEmail = 'defaultuser@fake.com'
