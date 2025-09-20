<template>
  <div class="users-view">
    <div class="page-header">
      <h1>User Management</h1>
    </div>

    <!-- Users Table -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Users</span>
          <div class="header-actions">
            <el-input
              v-model="searchText"
              placeholder="Search users..."
              style="width: 200px"
              clearable
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="filteredUsers"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="username" label="Username" min-width="150">
          <template #default="scope">
            <div class="user-info">
              <el-icon><User /></el-icon>
              <span>{{ scope.row.username }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="email" label="Email" min-width="200" />
        
        <el-table-column prop="role" label="Role" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.role === 'admin' ? 'warning' : 'info'"
              size="small"
            >
              {{ scope.row.role }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="is_active" label="Status" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.is_active ? 'success' : 'danger'"
              size="small"
            >
              {{ scope.row.is_active ? 'Active' : 'Inactive' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="Created" width="160">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="Actions" width="200" fixed="right">
          <template #default="scope">
            <el-button
              size="small"
              @click="editUser(scope.row)"
            >
              <el-icon><Edit /></el-icon>
            </el-button>
            
            <el-button
              size="small"
              :type="scope.row.is_active ? 'warning' : 'success'"
              @click="toggleUserStatus(scope.row)"
            >
              <el-icon v-if="scope.row.is_active"><Lock /></el-icon>
              <el-icon v-else><Unlock /></el-icon>
            </el-button>
            
            <el-button
              size="small"
              type="danger"
              @click="deleteUser(scope.row)"
              :disabled="scope.row.id === authStore.user?.id"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Edit User Dialog -->
    <el-dialog
      v-model="showEditDialog"
      title="Edit User"
      width="500px"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-width="100px"
      >
        <el-form-item label="Username" prop="username">
          <el-input v-model="editForm.username" />
        </el-form-item>
        
        <el-form-item label="Email" prop="email">
          <el-input v-model="editForm.email" type="email" />
        </el-form-item>
        
        <el-form-item label="Role" prop="role">
          <el-select v-model="editForm.role" style="width: 100%">
            <el-option label="User" value="user" />
            <el-option label="Admin" value="admin" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="Status" prop="is_active">
          <el-switch
            v-model="editForm.is_active"
            active-text="Active"
            inactive-text="Inactive"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showEditDialog = false">Cancel</el-button>
        <el-button
          type="primary"
          :loading="updating"
          @click="handleUpdateUser"
        >
          Update
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Search,
  User,
  Edit,
  Lock,
  Unlock,
  Delete
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import type { User as UserType } from '@/types'

const authStore = useAuthStore()

// Refs
const editFormRef = ref<FormInstance>()
const loading = ref(false)
const updating = ref(false)
const showEditDialog = ref(false)
const searchText = ref('')
const users = ref<UserType[]>([])

// Form data
const editForm = reactive({
  id: 0,
  username: '',
  email: '',
  role: 'user',
  is_active: true
})

// Form rules
const editRules: FormRules = {
  username: [
    { required: true, message: 'Please enter username', trigger: 'blur' },
    { min: 3, max: 50, message: 'Username must be 3-50 characters', trigger: 'blur' }
  ],
  email: [
    { required: true, message: 'Please enter email', trigger: 'blur' },
    { type: 'email', message: 'Please enter valid email', trigger: 'blur' }
  ],
  role: [
    { required: true, message: 'Please select role', trigger: 'change' }
  ]
}

// Computed
const filteredUsers = computed(() => {
  if (!searchText.value) return users.value
  
  const search = searchText.value.toLowerCase()
  return users.value.filter(user => 
    user.username.toLowerCase().includes(search) ||
    user.email.toLowerCase().includes(search)
  )
})

// Methods
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await api.get<{ users: UserType[] }>('/users')
    users.value = response.data.users
  } catch (error) {
    console.error('Error fetching users:', error)
  } finally {
    loading.value = false
  }
}

const editUser = (user: UserType) => {
  editForm.id = user.id
  editForm.username = user.username
  editForm.email = user.email
  editForm.role = user.role
  editForm.is_active = user.is_active
  showEditDialog.value = true
}

const handleUpdateUser = async () => {
  if (!editFormRef.value) return
  
  await editFormRef.value.validate(async (valid) => {
    if (valid) {
      updating.value = true
      try {
        const { id, ...updateData } = editForm
        const response = await api.put<{ user: UserType }>(`/users/${id}`, updateData)
        
        // Update local data
        const index = users.value.findIndex(u => u.id === id)
        if (index !== -1) {
          users.value[index] = response.data.user
        }
        
        ElMessage.success('User updated successfully')
        showEditDialog.value = false
      } catch (error) {
        console.error('Update error:', error)
      } finally {
        updating.value = false
      }
    }
  })
}

const toggleUserStatus = async (user: UserType) => {
  const action = user.is_active ? 'deactivate' : 'activate'
  
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to ${action} user "${user.username}"?`,
      'Confirm Action',
      {
        confirmButtonText: 'Yes',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
    
    const updateData = { is_active: !user.is_active }
    const response = await api.put<{ user: UserType }>(`/users/${user.id}`, updateData)
    
    // Update local data
    const index = users.value.findIndex(u => u.id === user.id)
    if (index !== -1) {
      users.value[index] = response.data.user
    }
    
    ElMessage.success(`User ${action}d successfully`)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Toggle status error:', error)
    }
  }
}

const deleteUser = async (user: UserType) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete user "${user.username}"? This action cannot be undone.`,
      'Confirm Deletion',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'error'
      }
    )
    
    await api.delete(`/users/${user.id}`)
    
    // Remove from local data
    users.value = users.value.filter(u => u.id !== user.id)
    
    ElMessage.success('User deleted successfully')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
    }
  }
}

// Lifecycle
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.users-view {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>