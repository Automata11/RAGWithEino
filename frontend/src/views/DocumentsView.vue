<template>
  <div class="documents-view">
    <div class="page-header">
      <h1>Documents</h1>
      <el-button type="primary" @click="showUploadDialog = true">
        <el-icon><Upload /></el-icon>
        Upload Document
      </el-button>
    </div>

    <!-- Documents Table -->
    <el-card>
      <el-table
        v-loading="documentStore.loading"
        :data="filteredDocuments"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="title" label="Title" min-width="200">
          <template #default="scope">
            <div class="document-title">
              <el-icon class="file-icon"><Document /></el-icon>
              <span>{{ scope.row.title }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="original_name" label="Original Name" min-width="150" />
        
        <el-table-column prop="file_type" label="Type" width="100">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.file_type }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="file_size" label="Size" width="100">
          <template #default="scope">
            {{ formatFileSize(scope.row.file_size) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="Status" width="120">
          <template #default="scope">
            <el-tag
              :type="getStatusTagType(scope.row.status)"
              size="small"
            >
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="vector_count" label="Vectors" width="100" />
        
        <el-table-column prop="created_at" label="Created" width="160">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="Actions" width="150" fixed="right">
          <template #default="scope">
            <el-button
              v-if="scope.row.status === 'processing'"
              size="small"
              @click="checkStatus(scope.row)"
            >
              <el-icon><View /></el-icon>
            </el-button>
            
            <el-button
              size="small"
              type="danger"
              @click="deleteDocument(scope.row)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Upload Dialog -->
    <el-dialog
      v-model="showUploadDialog"
      title="Upload Document"
      width="500px"
    >
      <el-form
        ref="uploadFormRef"
        :model="uploadForm"
        :rules="uploadRules"
        label-width="100px"
      >
        <el-form-item label="Title" prop="title">
          <el-input
            v-model="uploadForm.title"
            placeholder="Enter document title (optional)"
          />
        </el-form-item>
        
        <el-form-item label="File" prop="file">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :on-change="handleFileSelect"
            :on-remove="handleFileRemove"
            :limit="1"
            accept=".pdf,.txt,.doc,.docx,.md"
            drag
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              Drop file here or <em>click to upload</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                Supported formats: PDF, TXT, DOC, DOCX, MD
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showUploadDialog = false">Cancel</el-button>
        <el-button
          type="primary"
          :loading="uploading"
          :disabled="!uploadForm.file"
          @click="handleUpload"
        >
          {{ uploading ? `Uploading... ${documentStore.uploadProgress}%` : 'Upload' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Status Dialog -->
    <el-dialog
      v-model="showStatusDialog"
      title="Document Status"
      width="400px"
    >
      <div v-if="selectedDocument">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="Title">
            {{ selectedDocument.title }}
          </el-descriptions-item>
          <el-descriptions-item label="Status">
            <el-tag :type="getStatusTagType(documentStatus.status)">
              {{ documentStatus.status }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="documentStatus.detailed_status" label="Detail">
            {{ documentStatus.detailed_status }}
          </el-descriptions-item>
          <el-descriptions-item label="Vector Count">
            {{ documentStatus.vector_count }}
          </el-descriptions-item>
          <el-descriptions-item v-if="documentStatus.error_message" label="Error">
            <span style="color: #f56c6c">{{ documentStatus.error_message }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadInstance } from 'element-plus'
import {
  Upload,
  Document,
  View,
  Delete,
  UploadFilled
} from '@element-plus/icons-vue'
import { useDocumentStore } from '@/stores/document'
import type { Document as DocumentType, ProcessingStatus } from '@/types'

const documentStore = useDocumentStore()

// Refs
const uploadFormRef = ref<FormInstance>()
const uploadRef = ref<UploadInstance>()
const showUploadDialog = ref(false)
const showStatusDialog = ref(false)
const uploading = ref(false)
const selectedDocument = ref<DocumentType | null>(null)
const documentStatus = ref<ProcessingStatus>({
  status: '',
  vector_count: 0
})

// Form data
const uploadForm = reactive({
  title: '',
  file: null as File | null
})

// Form rules
const uploadRules: FormRules = {
  file: [
    { required: true, message: 'Please select a file', trigger: 'change' }
  ]
}

// Computed
const filteredDocuments = computed(() => documentStore.documents)

// Methods
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
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

const handleFileSelect = (file: any) => {
  uploadForm.file = file.raw
  if (!uploadForm.title) {
    uploadForm.title = file.name
  }
}

const handleFileRemove = () => {
  uploadForm.file = null
}

const handleUpload = async () => {
  if (!uploadFormRef.value || !uploadForm.file) return
  
  uploading.value = true
  try {
    const document = await documentStore.uploadDocument(
      uploadForm.file,
      uploadForm.title || undefined
    )
    
    ElMessage.success('Document uploaded successfully')
    showUploadDialog.value = false
    resetUploadForm()
    
    // Start polling for processing status
    documentStore.startStatusPolling(document.id)
    
  } catch (error) {
    console.error('Upload error:', error)
  } finally {
    uploading.value = false
  }
}

const resetUploadForm = () => {
  uploadForm.title = ''
  uploadForm.file = null
  uploadRef.value?.clearFiles()
}

const checkStatus = async (document: DocumentType) => {
  try {
    selectedDocument.value = document
    documentStatus.value = await documentStore.getProcessingStatus(document.id)
    showStatusDialog.value = true
  } catch (error) {
    console.error('Error fetching status:', error)
  }
}

const deleteDocument = async (document: DocumentType) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete "${document.title}"?`,
      'Confirm Deletion',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
    
    await documentStore.deleteDocument(document.id)
    ElMessage.success('Document deleted successfully')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
    }
  }
}

// Lifecycle
onMounted(() => {
  documentStore.fetchDocuments()
})
</script>

<style scoped>
.documents-view {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.document-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  color: #409eff;
}

.el-upload__tip {
  color: #606266;
  font-size: 12px;
  line-height: 1.5;
}
</style>