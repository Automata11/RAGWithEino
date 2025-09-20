import axios, { type AxiosInstance } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// Create axios instance
const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle common errors
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    const { response } = error
    
    if (response?.status === 401) {
      // Unauthorized - redirect to login
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
      ElMessage.error('Session expired. Please login again.')
    } else if (response?.status === 403) {
      ElMessage.error('Access denied')
    } else if (response?.status >= 500) {
      ElMessage.error('Server error. Please try again later.')
    } else if (response?.data?.error) {
      ElMessage.error(response.data.error)
    } else {
      ElMessage.error('An error occurred. Please try again.')
    }
    
    return Promise.reject(error)
  }
)

export default api