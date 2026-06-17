<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
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
  const label = ['admin', 'platform_admin'].includes(user.role) ? '管理端' : '社区'
  statusHint.value = `登录成功，正在进入${label}…`
  return router.replace(target).catch(() => {
    window.location.assign(target)
  })
}

function loginWithQQ() {
  const returnTo = '/community'
  const u = new URL('/user-api/api/v1/auth/qq/start', window.location.origin)
  u.searchParams.set('return_to', returnTo)
  window.location.assign(u.toString())
}

onMounted(() => {
  if (session.token && session.currentUser) {
    redirectAfterLogin(session.currentUser)
  }
})

async function submit() {
  if (!username.value.trim() || !password.value) {
    error.value = '请输入学号/用户名和密码'
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
    headline="万千帖子，齐聚 AI智联平台。"
    form-title="学号登录"
    form-hint="学生进入社区；管理员账号将自动进入管理端"
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
        placeholder="例如 demo_student"
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
      <p>首次使用？请用学校发放的邀请码注册账号</p>
      <RouterLink to="/register">邀请码注册</RouterLink>
      <div class="gx-auth-divider" />
      <button class="gx-auth-qq" type="button" @click="loginWithQQ">QQ 一键登录</button>
    </template>
  </GxAuthShell>
</template>

<style scoped>
.gx-auth-divider {
  height: 1px;
  margin: 12px 0;
  background: rgba(255, 255, 255, 0.08);
}
.gx-auth-qq {
  width: 100%;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: rgba(255, 255, 255, 0.06);
  color: inherit;
  padding: 10px 12px;
  font-weight: 600;
  cursor: pointer;
}
.gx-auth-qq:hover {
  background: rgba(255, 255, 255, 0.1);
}
</style>
