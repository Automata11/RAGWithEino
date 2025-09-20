<template>
  <div class="chat-view">
    <div class="chat-container">
      <!-- Conversation List -->
      <div class="conversation-sidebar">
        <div class="sidebar-header">
          <el-button
            type="primary"
            @click="startNewChat"
            style="width: 100%"
          >
            <el-icon><Plus /></el-icon>
            New Chat
          </el-button>
        </div>
        
        <div class="conversations-list">
          <div
            v-for="conversation in chatStore.conversations"
            :key="conversation.id"
            class="conversation-item"
            :class="{ active: selectedConversationId === conversation.id }"
            @click="selectConversation(conversation)"
          >
            <div class="conversation-title">{{ conversation.title }}</div>
            <div class="conversation-date">
              {{ formatDate(conversation.updated_at) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Chat Area -->
      <div class="chat-main">
        <div v-if="!chatStore.currentConversation && !isNewChat" class="welcome-screen">
          <div class="welcome-content">
            <el-icon size="64px" color="#409eff"><ChatLineRound /></el-icon>
            <h2>Welcome to RAG Chat</h2>
            <p>
              Start a conversation to ask questions about your uploaded documents.
              The AI will search through your document collection and provide relevant answers.
            </p>
            <el-button type="primary" @click="startNewChat">
              <el-icon><Plus /></el-icon>
              Start New Chat
            </el-button>
          </div>
        </div>

        <div v-else class="chat-content">
          <!-- Chat Header -->
          <div class="chat-header">
            <h3>
              {{ chatStore.currentConversation?.title || 'New Conversation' }}
            </h3>
            <div class="chat-actions">
              <el-button
                size="small"
                @click="startNewChat"
              >
                <el-icon><Plus /></el-icon>
                New
              </el-button>
            </div>
          </div>

          <!-- Messages -->
          <div class="messages-container" ref="messagesContainer">
            <div
              v-for="message in currentMessages"
              :key="message.id"
              class="message"
              :class="{ 'user-message': message.role === 'user' }"
            >
              <div class="message-avatar">
                <el-icon v-if="message.role === 'user'">
                  <User />
                </el-icon>
                <el-icon v-else>
                  <Service />
                </el-icon>
              </div>
              
              <div class="message-content">
                <div class="message-text">
                  {{ message.content }}
                </div>
                
                <!-- Sources for assistant messages -->
                <div v-if="message.role === 'assistant' && message.sources" class="message-sources">
                  <el-divider content-position="left">Sources</el-divider>
                  <div class="sources-list">
                    <!-- This would be populated with actual source data -->
                    <el-tag size="small" class="source-tag">
                      Document referenced
                    </el-tag>
                  </div>
                </div>
              </div>
            </div>

            <!-- Loading indicator -->
            <div v-if="chatStore.sending" class="message">
              <div class="message-avatar">
                <el-icon><Service /></el-icon>
              </div>
              <div class="message-content">
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          </div>

          <!-- Input Area -->
          <div class="input-area">
            <div class="input-container">
              <el-input
                v-model="messageText"
                type="textarea"
                :autosize="{ minRows: 1, maxRows: 4 }"
                placeholder="Ask a question about your documents..."
                @keydown.enter.exact="handleSendMessage"
                @keydown.enter.shift.exact.prevent="messageText += '\n'"
                :disabled="chatStore.sending"
              />
              <el-button
                type="primary"
                :disabled="!messageText.trim() || chatStore.sending"
                @click="handleSendMessage"
              >
                <el-icon><Promotion /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Plus,
  ChatLineRound,
  User,
  Service, // Using Service instead of Robot
  Promotion
} from '@element-plus/icons-vue'
import { useChatStore } from '@/stores/chat'
import type { Conversation } from '@/types'

const route = useRoute()
const chatStore = useChatStore()

// Refs
const messagesContainer = ref<HTMLElement>()
const messageText = ref('')
const isNewChat = ref(false)
const selectedConversationId = ref<number | null>(null)

// Computed
const currentMessages = computed(() => {
  return chatStore.currentConversation?.messages || []
})

// Methods
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  } else if (diffDays < 7) {
    return diffDays === 1 ? 'Yesterday' : `${diffDays} days ago`
  } else {
    return date.toLocaleDateString()
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const startNewChat = () => {
  chatStore.startNewConversation()
  selectedConversationId.value = null
  isNewChat.value = true
  messageText.value = ''
}

const selectConversation = async (conversation: Conversation) => {
  selectedConversationId.value = conversation.id
  isNewChat.value = false
  await chatStore.fetchConversation(conversation.id)
  scrollToBottom()
}

const handleSendMessage = async () => {
  const text = messageText.value.trim()
  if (!text || chatStore.sending) return

  const queryRequest = {
    question: text,
    conversation_id: selectedConversationId.value || undefined
  }

  // Clear input immediately
  messageText.value = ''

  try {
    const response = await chatStore.sendQuery(queryRequest)
    
    // If this was a new chat, update the selected conversation
    if (!selectedConversationId.value) {
      selectedConversationId.value = response.conversation_id
      isNewChat.value = false
      // Refresh conversations list to show the new conversation
      await chatStore.fetchConversations()
    }
    
    scrollToBottom()
  } catch (error) {
    console.error('Error sending message:', error)
    ElMessage.error('Failed to send message')
  }
}

// Watchers
watch(currentMessages, () => {
  scrollToBottom()
}, { deep: true })

// Lifecycle
onMounted(async () => {
  await chatStore.fetchConversations()
  
  // Check if there's a conversation ID in the query params
  const conversationId = route.query.conversation
  if (conversationId) {
    const id = parseInt(conversationId as string)
    const conversation = chatStore.conversations.find(c => c.id === id)
    if (conversation) {
      await selectConversation(conversation)
    }
  }
})
</script>

<style scoped>
.chat-view {
  height: calc(100vh - 140px);
  max-width: 1200px;
  margin: 0 auto;
}

.chat-container {
  display: flex;
  height: 100%;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

/* Conversation Sidebar */
.conversation-sidebar {
  width: 300px;
  border-right: 1px solid #e4e7ed;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.conversations-list {
  flex: 1;
  overflow-y: auto;
}

.conversation-item {
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background-color 0.2s;
}

.conversation-item:hover {
  background-color: #e9ecef;
}

.conversation-item.active {
  background-color: #409eff;
  color: white;
}

.conversation-title {
  font-weight: 500;
  margin-bottom: 4px;
  font-size: 14px;
}

.conversation-date {
  font-size: 12px;
  color: #6c757d;
}

.conversation-item.active .conversation-date {
  color: rgba(255, 255, 255, 0.8);
}

/* Chat Main */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
}

.welcome-screen {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.welcome-content {
  text-align: center;
  max-width: 400px;
}

.welcome-content h2 {
  margin: 20px 0 16px 0;
  color: #303133;
}

.welcome-content p {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 24px;
}

.chat-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-header h3 {
  margin: 0;
  color: #303133;
}

/* Messages */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.message {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.message.user-message {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-message .message-avatar {
  background: #409eff;
  color: white;
}

.message-content {
  flex: 1;
  max-width: 70%;
}

.user-message .message-content {
  text-align: right;
}

.message-text {
  background: #f5f5f5;
  padding: 12px 16px;
  border-radius: 18px;
  line-height: 1.5;
  white-space: pre-wrap;
}

.user-message .message-text {
  background: #409eff;
  color: white;
}

.message-sources {
  margin-top: 12px;
}

.sources-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.source-tag {
  cursor: pointer;
}

/* Typing indicator */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 18px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #999;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) {
  animation-delay: 0s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
  }
  30% {
    transform: translateY(-10px);
  }
}

/* Input Area */
.input-area {
  border-top: 1px solid #e4e7ed;
  padding: 16px 20px;
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.input-container .el-textarea {
  flex: 1;
}
</style>