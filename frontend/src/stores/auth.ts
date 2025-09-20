import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, RegisterRequest, LoginResponse } from '@/types'
import api from '@/utils/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!user.value && !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  // Initialize auth state from localStorage
  const checkAuth = () => {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    
    if (storedToken && storedUser) {
      token.value = storedToken
      user.value = JSON.parse(storedUser)
    }
  }

  const login = async (credentials: LoginRequest): Promise<void> => {
    loading.value = true
    try {
      const response = await api.post<LoginResponse>('/auth/login', credentials)
      const { token: authToken, user: userData } = response.data
      
      token.value = authToken
      user.value = userData
      
      // Store in localStorage
      localStorage.setItem('token', authToken)
      localStorage.setItem('user', JSON.stringify(userData))
    } finally {
      loading.value = false
    }
  }

  const register = async (userData: RegisterRequest): Promise<void> => {
    loading.value = true
    try {
      await api.post('/auth/register', userData)
    } finally {
      loading.value = false
    }
  }

  const logout = async (): Promise<void> => {
    loading.value = true
    try {
      await api.post('/auth/logout')
    } catch {
      // Even if logout fails on server, clear local state
    } finally {
      // Clear local state
      token.value = null
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      loading.value = false
    }
  }

  const getProfile = async (): Promise<void> => {
    try {
      const response = await api.get<{ user: User }>('/auth/profile')
      user.value = response.data.user
      localStorage.setItem('user', JSON.stringify(response.data.user))
    } catch (error) {
      // If profile fetch fails, logout
      logout()
      throw error
    }
  }

  return {
    user: computed(() => user.value),
    token: computed(() => token.value),
    loading: computed(() => loading.value),
    isLoggedIn,
    isAdmin,
    checkAuth,
    login,
    register,
    logout,
    getProfile
  }
})