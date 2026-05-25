<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '../stores/session'
import { userApi } from '../api'

const router = useRouter()
const session = useSessionStore()

const form = ref({
  username: '',
  password: '',
  invitationCode: '',
})
const loading = ref(false)
const error = ref('')

async function submit() {
  if (!form.value.username.trim() || !form.value.password || !form.value.invitationCode.trim()) {
    error.value = '请填写用户名、密码和邀请码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const result = await userApi.register({
      username: form.value.username.trim(),
      password: form.value.password,
      invitationCode: form.value.invitationCode.trim(),
    })
    session.persistSession(result)
    session.setFlash('注册成功，欢迎加入 AI 智联论坛。', 'success')
    router.push('/community')
  } catch (err) {
    error.value = err.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="login-shell">
    <section class="login-hero">
      <p class="eyebrow">邀请注册</p>
      <h1>用邀请码开通你的学院论坛账号。</h1>
      <p class="lead">注册后即可体验社区发帖、评论与中台演示流程（视角色而定）。</p>
    </section>

    <section class="login-panel panel">
      <div class="section-copy">
        <p class="eyebrow">创建账号</p>
        <h2>填写注册信息</h2>
      </div>
      <form class="role-login-form" @submit.prevent="submit">
        <label class="role-login-field">
          <span>用户名</span>
          <input v-model="form.username" type="text" autocomplete="username" placeholder="3-32 位字符" />
        </label>
        <label class="role-login-field">
          <span>密码</span>
          <input v-model="form.password" type="password" autocomplete="new-password" />
        </label>
        <label class="role-login-field">
          <span>邀请码</span>
          <input
            v-model="form.invitationCode"
            type="text"
            inputmode="text"
            autocomplete="one-time-code"
            placeholder="向管理员索取"
          />
        </label>
        <p v-if="error" class="role-login-error">{{ error }}</p>
        <div class="role-login-actions">
          <RouterLink to="/" class="ghost-button">返回登录</RouterLink>
          <button type="submit" class="primary-button" :disabled="loading">
            {{ loading ? '注册中…' : '完成注册' }}
          </button>
        </div>
      </form>
    </section>
  </main>
</template>
