export interface User {
  id: number
  username: string
  email: string
  role: 'user' | 'admin'
  is_active: boolean
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface Document {
  id: number
  user_id: number
  user?: User
  title: string
  original_name: string
  file_path: string
  file_size: number
  file_type: string
  status: 'uploading' | 'processing' | 'completed' | 'failed'
  processed_at?: string
  vector_count: number
  collection_id?: string
  error_message?: string
  created_at: string
  updated_at: string
}

export interface ProcessingStatus {
  status: string
  detailed_status?: string
  processed_at?: string
  vector_count: number
  error_message?: string
}

export interface Conversation {
  id: number
  user_id: number
  user?: User
  title: string
  created_at: string
  updated_at: string
  messages?: Message[]
}

export interface Message {
  id: number
  conversation_id: number
  role: 'user' | 'assistant'
  content: string
  sources?: string
  created_at: string
  updated_at: string
}

export interface QueryRequest {
  question: string
  conversation_id?: number
}

export interface DocumentSource {
  document_id: number
  document_name: string
  chunk_id: number
  content: string
  similarity: number
}

export interface QueryResponse {
  answer: string
  sources: DocumentSource[]
  conversation_id: number
  message_id: number
}

export interface APIResponse<T = any> {
  data?: T
  message?: string
  error?: string
}