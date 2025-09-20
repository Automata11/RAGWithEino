<template>
  <div class="not-found-view">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1>Page Not Found</h1>
      <p>The page you are looking for doesn't exist or has been moved.</p>
      <el-button type="primary" @click="$router.push('/')">
        <el-icon><House /></el-icon>
        Go Home
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { House } from '@element-plus/icons-vue'
</script>

<style scoped>
.not-found-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.not-found-content {
  text-align: center;
  max-width: 500px;
  padding: 40px;
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  margin-bottom: 20px;
  opacity: 0.8;
}

h1 {
  font-size: 48px;
  margin-bottom: 16px;
  font-weight: 600;
}

p {
  font-size: 18px;
  margin-bottom: 32px;
  opacity: 0.9;
  line-height: 1.6;
}
</style>