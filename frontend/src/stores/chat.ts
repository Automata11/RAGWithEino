import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Conversation, Message, QueryRequest, QueryResponse } from '@/types'
import api from '@/utils/api'

export const useChatStore = defineStore('chat', () => {
  const conversations = ref<Conversation[]>([])
  const currentConversation = ref<Conversation | null>(null)
  const loading = ref(false)
  const sending = ref(false)

  const fetchConversations = async (): Promise<void> => {
    loading.value = true
    try {
      const response = await api.get<{ conversations: Conversation[] }>('/rag/conversations')
      conversations.value = response.data.conversations
    } finally {
      loading.value = false
    }
  }

  const fetchConversation = async (id: number): Promise<void> => {
    loading.value = true
    try {
      const response = await api.get<{ conversation: Conversation }>(`/rag/conversations/${id}`)
      currentConversation.value = response.data.conversation
    } finally {
      loading.value = false
    }
  }

  const sendQuery = async (query: QueryRequest): Promise<QueryResponse> => {
    sending.value = true
    try {
      const response = await api.post<QueryResponse>('/rag/query', query)
      const result = response.data

      // Update conversations list if new conversation
      if (!query.conversation_id) {
        await fetchConversations()
      }

      // Update current conversation if it's the active one
      if (currentConversation.value?.id === result.conversation_id) {
        await fetchConversation(result.conversation_id)
      }

      return result
    } finally {
      sending.value = false
    }
  }

  const startNewConversation = () => {
    currentConversation.value = null
  }

  const setCurrentConversation = (conversation: Conversation) => {
    currentConversation.value = conversation
  }

  // Add message to current conversation optimistically
  const addMessageOptimistically = (content: string, role: 'user' | 'assistant' = 'user') => {
    if (currentConversation.value) {
      const newMessage: Message = {
        id: Date.now(), // Temporary ID
        conversation_id: currentConversation.value.id,
        role,
        content,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
      
      if (!currentConversation.value.messages) {
        currentConversation.value.messages = []
      }
      
      currentConversation.value.messages.push(newMessage)
    }
  }

  return {
    conversations,
    currentConversation,
    loading,
    sending,
    fetchConversations,
    fetchConversation,
    sendQuery,
    startNewConversation,
    setCurrentConversation,
    addMessageOptimistically
  }
})