<template>
  <el-container class="dashboard-container">
    <!-- Header -->
    <el-header class="header">
      <div class="header-left">
        <h2>RAG with Eino</h2>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleMenuCommand">
          <el-button text>
            <el-icon><User /></el-icon>
            {{ authStore.user?.username }}
            <el-icon><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                Profile
              </el-dropdown-item>
              <el-dropdown-item command="logout" divided>
                <el-icon><SwitchButton /></el-icon>
                Logout
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <el-container>
      <!-- Sidebar -->
      <el-aside width="250px" class="sidebar">
        <el-menu
          :default-active="$route.path"
          router
          class="sidebar-menu"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>Dashboard</span>
          </el-menu-item>
          
          <el-menu-item index="/chat">
            <el-icon><ChatLineRound /></el-icon>
            <span>Chat</span>
          </el-menu-item>
          
          <el-menu-item index="/documents">
            <el-icon><Document /></el-icon>
            <span>Documents</span>
          </el-menu-item>
          
          <el-menu-item v-if="authStore.isAdmin" index="/users">
            <el-icon><User /></el-icon>
            <span>Users</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <!-- Main Content -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  House,
  ChatLineRound,
  Document,
  User,
  ArrowDown,
  SwitchButton
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const handleMenuCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm(
          'Are you sure you want to logout?',
          'Confirm',
          {
            confirmButtonText: 'Yes',
            cancelButtonText: 'Cancel',
            type: 'warning'
          }
        )
        await authStore.logout()
        ElMessage.success('Logged out successfully')
        router.push('/login')
      } catch (error) {
        // User cancelled or error occurred
      }
      break
  }
}
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left h2 {
  margin: 0;
  color: #303133;
  font-size: 20px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
}

.sidebar {
  background-color: #f5f5f5;
  border-right: 1px solid #e4e7ed;
}

.sidebar-menu {
  border: none;
  background-color: transparent;
}

.main-content {
  background-color: #fff;
  padding: 20px;
}

.el-menu-item {
  height: 50px;
  line-height: 50px;
}

.el-menu-item.is-active {
  background-color: #ecf5ff;
  color: #409eff;
}

.el-menu-item:hover {
  background-color: #f0f0f0;
}
</style>