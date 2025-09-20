<template>
  <div class="dashboard-view">
    <div class="page-header">
      <h1>Dashboard</h1>
    </div>

    <div class="stats-cards">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-content">
              <div class="stats-icon documents">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stats-text">
                <h3>{{ documents.length }}</h3>
                <p>Documents</p>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-content">
              <div class="stats-icon conversations">
                <el-icon><ChatLineRound /></el-icon>
              </div>
              <div class="stats-text">
                <h3>{{ conversations.length }}</h3>
                <p>Conversations</p>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-content">
              <div class="stats-icon processing">
                <el-icon><Loading /></el-icon>
              </div>
              <div class="stats-text">
                <h3>{{ processingDocuments }}</h3>
                <p>Processing</p>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-content">
              <div class="stats-icon completed">
                <el-icon><Check /></el-icon>
              </div>
              <div class="stats-text">
                <h3>{{ completedDocuments }}</h3>
                <p>Completed</p>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <div class="content-sections">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="content-card">
            <template #header>
              <div class="card-header">
                <span>Recent Documents</span>
                <router-link to="/documents">View All</router-link>
              </div>
            </template>
            
            <div v-if="recentDocuments.length === 0" class="empty-state">
              <el-empty description="No documents yet" />
              <el-button type="primary" @click="$router.push('/documents')">
                Upload Document
              </el-button>
            </div>
            
            <div v-else>
              <div
                v-for="doc in recentDocuments"
                :key="doc.id"
                class="document-item"
              >
                <div class="document-info">
                  <h4>{{ doc.title }}</h4>
                  <p>{{ formatDate(doc.created_at) }}</p>
                </div>
                <el-tag
                  :type="getStatusTagType(doc.status)"
                  size="small"
                >
                  {{ doc.status }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="12">
          <el-card class="content-card">
            <template #header>
              <div class="card-header">
                <span>Recent Conversations</span>
                <router-link to="/chat">View All</router-link>
              </div>
            </template>
            
            <div v-if="recentConversations.length === 0" class="empty-state">
              <el-empty description="No conversations yet" />
              <el-button type="primary" @click="$router.push('/chat')">
                Start Chat
              </el-button>
            </div>
            
            <div v-else>
              <div
                v-for="conv in recentConversations"
                :key="conv.id"
                class="conversation-item"
                @click="goToConversation(conv.id)"
              >
                <div class="conversation-info">
                  <h4>{{ conv.title }}</h4>
                  <p>{{ formatDate(conv.updated_at) }}</p>
                </div>
                <el-icon><ArrowRight /></el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- Quick Actions -->
    <div class="quick-actions">
      <el-card>
        <template #header>
          <span>Quick Actions</span>
        </template>
        
        <el-row :gutter="20">
          <el-col :span="8">
            <el-button
              type="primary"
              size="large"
              @click="$router.push('/documents')"
              style="width: 100%"
            >
              <el-icon><Upload /></el-icon>
              Upload Document
            </el-button>
          </el-col>
          
          <el-col :span="8">
            <el-button
              type="success"
              size="large"
              @click="$router.push('/chat')"
              style="width: 100%"
            >
              <el-icon><ChatLineRound /></el-icon>
              Start Chat
            </el-button>
          </el-col>
          
          <el-col :span="8">
            <el-button
              v-if="authStore.isAdmin"
              type="warning"
              size="large"
              @click="$router.push('/users')"
              style="width: 100%"
            >
              <el-icon><User /></el-icon>
              Manage Users
            </el-button>
          </el-col>
        </el-row>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Document,
  ChatLineRound,
  Loading,
  Check,
  ArrowRight,
  Upload,
  User
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useDocumentStore } from '@/stores/document'
import { useChatStore } from '@/stores/chat'

const router = useRouter()
const authStore = useAuthStore()
const documentStore = useDocumentStore()
const chatStore = useChatStore()

const documents = computed(() => documentStore.documents)
const conversations = computed(() => chatStore.conversations)

const processingDocuments = computed(() => 
  documents.value.filter(doc => doc.status === 'processing' || doc.status === 'uploading').length
)

const completedDocuments = computed(() => 
  documents.value.filter(doc => doc.status === 'completed').length
)

const recentDocuments = computed(() => 
  [...documents.value]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5)
)

const recentConversations = computed(() => 
  [...conversations.value]
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 5)
)

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const getStatusTagType = (status: string) => {
  switch (status) {
    case 'completed': return 'success'
    case 'processing': return 'warning'
    case 'uploading': return 'info'
    case 'failed': return 'danger'
    default: return 'info'
  }
}

const goToConversation = (id: number) => {
  router.push('/chat?conversation=' + id)
}

onMounted(async () => {
  // Load data for dashboard
  if (documents.value.length === 0) {
    documentStore.fetchDocuments()
  }
  if (conversations.value.length === 0) {
    chatStore.fetchConversations()
  }
})
</script>

<style scoped>
.dashboard-view {
  max-width: 1200px;
  margin: 0 auto;
}

.stats-cards {
  margin-bottom: 30px;
}

.stats-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.stats-card:hover {
  transform: translateY(-2px);
}

.stats-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stats-icon {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stats-icon.documents {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stats-icon.conversations {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stats-icon.processing {
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
  color: #e6a23c;
}

.stats-icon.completed {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #67c23a;
}

.stats-text h3 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.stats-text p {
  margin: 5px 0 0 0;
  color: #909399;
  font-size: 14px;
}

.content-sections {
  margin-bottom: 30px;
}

.content-card {
  height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header a {
  color: #409eff;
  text-decoration: none;
  font-size: 14px;
}

.card-header a:hover {
  text-decoration: underline;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
}

.document-item,
.conversation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.document-item:last-child,
.conversation-item:last-child {
  border-bottom: none;
}

.conversation-item {
  cursor: pointer;
  transition: background-color 0.2s;
}

.conversation-item:hover {
  background-color: #f5f7fa;
  margin: 0 -16px;
  padding-left: 16px;
  padding-right: 16px;
}

.document-info h4,
.conversation-info h4 {
  margin: 0 0 5px 0;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.document-info p,
.conversation-info p {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

.quick-actions {
  margin-bottom: 30px;
}
</style>