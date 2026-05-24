<template>
  <div class="admin-login-page">
    <el-card class="login-card">
      <h2>管理后台登录</h2>
      <el-form @submit.prevent="handleLogin">
        <el-form-item>
          <el-input v-model="username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" style="width: 100%" @click="handleLogin">登录</el-button>
        </el-form-item>
      </el-form>
      <p v-if="error" style="color: #f56c6c; font-size: 14px;">{{ error }}</p>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import axios from 'axios'

const username = ref('')
const password = ref('')
const error = ref('')
const router = useRouter()
const adminStore = useAdminStore()

async function handleLogin() {
  error.value = ''
  try {
    const res = await axios.post('/api/v1/login', {
      username: username.value,
      password: password.value,
    })
    localStorage.setItem('access_token', res.data.access_token)
    // Check admin role from token
    const token = res.data.access_token
    const payload = JSON.parse(atob(token.split('.')[1]))
    if (payload.role !== 'admin' && payload.role !== 'platform_admin') {
      localStorage.removeItem('access_token')
      error.value = '您没有管理员权限'
      return
    }
    adminStore.setAdminUser({ username: res.data.user.username, role: payload.role })
    router.push('/admin')
  } catch (e) {
    error.value = e.response?.data?.error || '登录失败'
  }
}
</script>

<style scoped>
.admin-login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
}
.login-card {
  width: 400px;
}
.login-card h2 {
  text-align: center;
  margin-bottom: 20px;
}
</style>
