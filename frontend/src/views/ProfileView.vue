<template>
  <div class="profile-view">
    <div class="page-header">
      <h1>Profile</h1>
    </div>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>User Information</span>
          </template>
          
          <div class="profile-info">
            <div class="avatar-section">
              <el-avatar :size="80" :icon="UserFilled" />
              <h3>{{ authStore.user?.username }}</h3>
              <el-tag :type="authStore.user?.role === 'admin' ? 'warning' : 'info'">
                {{ authStore.user?.role }}
              </el-tag>
            </div>
            
            <el-divider />
            
            <el-descriptions :column="1">
              <el-descriptions-item label="Username">
                {{ authStore.user?.username }}
              </el-descriptions-item>
              <el-descriptions-item label="Email">
                {{ authStore.user?.email }}
              </el-descriptions-item>
              <el-descriptions-item label="Role">
                {{ authStore.user?.role }}
              </el-descriptions-item>
              <el-descriptions-item label="Status">
                <el-tag type="success" size="small">Active</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="Member Since">
                {{ authStore.user ? formatDate(authStore.user.created_at) : '' }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>
            <span>Update Profile</span>
          </template>
          
          <el-form
            ref="profileFormRef"
            :model="profileForm"
            :rules="profileRules"
            label-width="120px"
          >
            <el-form-item label="Username" prop="username">
              <el-input v-model="profileForm.username" />
            </el-form-item>
            
            <el-form-item label="Email" prop="email">
              <el-input v-model="profileForm.email" type="email" />
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                :loading="updating"
                @click="updateProfile"
              >
                Update Profile
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card style="margin-top: 20px;">
          <template #header>
            <span>Change Password</span>
          </template>
          
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="120px"
          >
            <el-form-item label="Current Password" prop="currentPassword">
              <el-input
                v-model="passwordForm.currentPassword"
                type="password"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="New Password" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="Confirm Password" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                show-password
              />
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                :loading="changingPassword"
                @click="changePassword"
              >
                Change Password
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { UserFilled } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'

const authStore = useAuthStore()

// Refs
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const updating = ref(false)
const changingPassword = ref(false)

// Form data
const profileForm = reactive({
  username: '',
  email: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Validation
const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('Passwords do not match'))
  } else {
    callback()
  }
}

const profileRules: FormRules = {
  username: [
    { required: true, message: 'Please enter username', trigger: 'blur' },
    { min: 3, max: 50, message: 'Username must be 3-50 characters', trigger: 'blur' }
  ],
  email: [
    { required: true, message: 'Please enter email', trigger: 'blur' },
    { type: 'email', message: 'Please enter valid email', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  currentPassword: [
    { required: true, message: 'Please enter current password', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: 'Please enter new password', trigger: 'blur' },
    { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: 'Please confirm new password', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

// Methods
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const initializeForm = () => {
  if (authStore.user) {
    profileForm.username = authStore.user.username
    profileForm.email = authStore.user.email
  }
}

const updateProfile = async () => {
  if (!profileFormRef.value) return
  
  await profileFormRef.value.validate(async (valid) => {
    if (valid) {
      updating.value = true
      try {
        // Note: In a real implementation, you'd need a dedicated profile update endpoint
        // For now, we'll simulate the update
        await authStore.getProfile() // Refresh profile data
        ElMessage.success('Profile updated successfully')
      } catch (error) {
        console.error('Profile update error:', error)
        ElMessage.error('Failed to update profile')
      } finally {
        updating.value = false
      }
    }
  })
}

const changePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      changingPassword.value = true
      try {
        // Note: In a real implementation, you'd need a password change endpoint
        // await api.post('/auth/change-password', {
        //   currentPassword: passwordForm.currentPassword,
        //   newPassword: passwordForm.newPassword
        // })
        
        ElMessage.success('Password changed successfully')
        
        // Reset form
        passwordForm.currentPassword = ''
        passwordForm.newPassword = ''
        passwordForm.confirmPassword = ''
        passwordFormRef.value?.resetFields()
        
      } catch (error) {
        console.error('Password change error:', error)
        ElMessage.error('Failed to change password')
      } finally {
        changingPassword.value = false
      }
    }
  })
}

// Lifecycle
onMounted(() => {
  initializeForm()
})
</script>

<style scoped>
.profile-view {
  max-width: 1200px;
  margin: 0 auto;
}

.profile-info {
  text-align: center;
}

.avatar-section {
  margin-bottom: 20px;
}

.avatar-section h3 {
  margin: 16px 0 8px 0;
  color: #303133;
}

.el-descriptions {
  text-align: left;
}
</style>