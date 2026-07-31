<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import GxAuthShell from '../components/gx/GxAuthShell.vue'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const session = useSessionStore()
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const statusHint = ref('')

function redirectAfterLogin(user) {
  const target = session.routeAfterLogin(user)
  const label = ['admin', 'platform_admin'].includes(user.role) ? '管理端' : '论坛'
  statusHint.value = `登录成功，正在进入${label}...`
  return router.replace(target).catch(() => {
    window.location.assign(target)
  })
}

onMounted(() => {
  if (session.token && session.currentUser) {
    redirectAfterLogin(session.currentUser)
  }
})

async function submit() {
  if (!username.value.trim() || !password.value) {
    error.value = '请输入学号或用户名和密码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const user = await session.loginWithCredentials(username.value.trim(), password.value)
    await redirectAfterLogin(user)
    session.setFlash(`欢迎回来，${user.name || user.username}`, 'success')
  } catch (err) {
    error.value = err.message || '登录失败，请检查账号和密码'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <GxAuthShell
    headline="AI 智联论坛"
    form-title="账号登录"
    form-hint="使用已发放的学号或用户名登录。"
    submit-label="登录"
    :loading="loading"
    :error="error"
    :status-hint="statusHint"
    @submit="submit"
  >
    <div class="gx-auth-field">
      <label for="username">学号 / 用户名</label>
      <input
        id="username"
        v-model="username"
        type="text"
        autocomplete="username"
        placeholder="请输入学号或用户名"
      />
    </div>
    <div class="gx-auth-field">
      <label for="password">密码</label>
      <input
        id="password"
        v-model="password"
        type="password"
        autocomplete="current-password"
        placeholder="请输入密码"
      />
    </div>

    <template #footer>
      <p>账号由管理员统一发放。如无法登录，请联系运营同学处理。</p>
    </template>
  </GxAuthShell>
</template>
