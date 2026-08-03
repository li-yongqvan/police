<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSessionStore } from '../stores/session'
import GxAuthShell from '../components/gx/GxAuthShell.vue'

const route = useRoute()
const router = useRouter()
const session = useSessionStore()

const loading = ref(true)
const error = ref('')

function decodeUser(encoded) {
  if (!encoded) return null
  try {
    const b64 = encoded.replace(/-/g, '+').replace(/_/g, '/')
    const padded = b64 + '==='.slice((b64.length + 3) % 4)
    const bytes = Uint8Array.from(window.atob(padded), (c) => c.charCodeAt(0))
    const json = new TextDecoder().decode(bytes)
    return JSON.parse(json || 'null')
  } catch {
    return null
  }
}

onMounted(async () => {
  const token = String(route.query.token || '')
  const refreshToken = String(route.query.refresh_token || '')
  const returnTo = String(route.query.return_to || '')
  const user = decodeUser(String(route.query.user || ''))

  if (!token) {
    error.value = 'QQ 登录回调缺少令牌，请重试'
    loading.value = false
    return
  }

  try {
    session.persistSession({
      token,
      refreshToken: refreshToken || '',
      user: user || null,
    })
    if (!user) {
      await session.refreshMe()
    }

    const target =
      returnTo ||
      session.routeAfterLogin(session.currentUser || { role: 'student' })
    await router.replace(target)
  } catch (err) {
    error.value = err?.message || 'QQ 登录失败，请重试'
    // Cleanup invalid token
    session.logout()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <GxAuthShell
    headline="正在完成 QQ 登录…"
    form-title="请稍候"
    form-hint="若长时间无响应，请返回重新登录。"
    submit-label="返回登录"
    :loading="loading"
    :error="error"
    @submit="() => router.replace('/')"
  >
    <template #footer>
      <p v-if="!loading && error">你也可以尝试使用学号/用户名登录。</p>
    </template>
  </GxAuthShell>
</template>

