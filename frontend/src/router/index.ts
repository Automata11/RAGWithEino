import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/RegisterView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/',
      component: () => import('@/views/DashboardLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('@/views/DashboardView.vue')
        },
        {
          path: 'chat',
          name: 'Chat',
          component: () => import('@/views/ChatView.vue')
        },
        {
          path: 'documents',
          name: 'Documents',
          component: () => import('@/views/DocumentsView.vue')
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/UsersView.vue'),
          meta: { requiresAdmin: true }
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/ProfileView.vue')
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundView.vue')
    }
  ]
})

// Route guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Check if user is authenticated
  if (!authStore.isLoggedIn) {
    authStore.checkAuth()
  }
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    // Redirect to login if authentication required
    next('/login')
  } else if (to.meta.requiresGuest && authStore.isLoggedIn) {
    // Redirect to dashboard if already logged in
    next('/')
  } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    // Redirect to dashboard if admin access required but user is not admin
    next('/')
  } else {
    next()
  }
})

export default router