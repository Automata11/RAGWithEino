import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Document, ProcessingStatus } from '@/types'
import api from '@/utils/api'

export const useDocumentStore = defineStore('document', () => {
  const documents = ref<Document[]>([])
  const loading = ref(false)
  const uploadProgress = ref(0)

  const fetchDocuments = async (): Promise<void> => {
    loading.value = true
    try {
      const response = await api.get<{ documents: Document[] }>('/documents')
      documents.value = response.data.documents
    } finally {
      loading.value = false
    }
  }

  const uploadDocument = async (file: File, title?: string): Promise<Document> => {
    const formData = new FormData()
    formData.append('file', file)
    if (title) {
      formData.append('title', title)
    }

    const response = await api.post<{ document: Document }>('/documents', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total) {
          uploadProgress.value = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          )
        }
      }
    })

    // Add to local state
    documents.value.unshift(response.data.document)
    uploadProgress.value = 0
    
    return response.data.document
  }

  const getDocument = async (id: number): Promise<Document> => {
    const response = await api.get<{ document: Document }>(`/documents/${id}`)
    return response.data.document
  }

  const deleteDocument = async (id: number): Promise<void> => {
    await api.delete(`/documents/${id}`)
    // Remove from local state
    documents.value = documents.value.filter(doc => doc.id !== id)
  }

  const getProcessingStatus = async (id: number): Promise<ProcessingStatus> => {
    const response = await api.get<ProcessingStatus>(`/documents/${id}/status`)
    return response.data
  }

  const updateDocumentStatus = (id: number, status: ProcessingStatus) => {
    const docIndex = documents.value.findIndex(doc => doc.id === id)
    if (docIndex !== -1) {
      documents.value[docIndex].status = status.status as Document['status']
      documents.value[docIndex].vector_count = status.vector_count
      documents.value[docIndex].error_message = status.error_message
      if (status.processed_at) {
        documents.value[docIndex].processed_at = status.processed_at
      }
    }
  }

  // Polling for document processing status
  const startStatusPolling = (documentId: number, callback?: (status: ProcessingStatus) => void) => {
    const poll = async () => {
      try {
        const status = await getProcessingStatus(documentId)
        updateDocumentStatus(documentId, status)
        
        if (callback) {
          callback(status)
        }
        
        // Continue polling if still processing
        if (status.status === 'processing' || status.status === 'uploading') {
          setTimeout(poll, 2000) // Poll every 2 seconds
        }
      } catch (error) {
        console.error('Error polling status:', error)
      }
    }
    
    poll()
  }

  return {
    documents,
    loading,
    uploadProgress,
    fetchDocuments,
    uploadDocument,
    getDocument,
    deleteDocument,
    getProcessingStatus,
    updateDocumentStatus,
    startStatusPolling
  }
})