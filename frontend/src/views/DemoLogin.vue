<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const session = useSessionStore()

const roles = [
  {
    role: 'student',
    title: '学生用户',
    intro: '查看论坛、发帖求助、评论互动、体验社区主流程。',
    usernamePlaceholder: '例如 demo_student',
  },
  {
    role: 'admin',
    title: '协会管理员',
    intro: '查看内容审核、活动运营和用户封禁的最小闭环。',
    usernamePlaceholder: '例如 demo_admin',
  },
  {
    role: 'platform_admin',
    title: '中台管理员',
    intro: '查看系统配置、监管规则和数据概览。',
    usernamePlaceholder: '例如 demo_platform_admin',
  },
]

const activeRole = ref('')
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

function openLogin(role) {
  activeRole.value = role
  username.value = ''
  password.value = ''
  error.value = ''
}

function cancelLogin() {
  activeRole.value = ''
  username.value = ''
  password.value = ''
  error.value = ''
}

async function submit(role) {
  if (!username.value.trim() || !password.value) {
    error.value = '请输入用户名和密码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const user = await session.loginWithCredentials(
      username.value.trim(),
      password.value,
      role,
    )
    session.setFlash(`欢迎回来，${user.name || user.id}。`, 'success')
    cancelLogin()
    if (user.role === 'student') {
      router.push('/community')
    } else {
      router.push('/admin')
    }
  } catch (err) {
    error.value = err.message || '登录失败，请检查账号和密码'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="login-shell">
    <section class="login-hero">
      <p class="eyebrow">AI 智联论坛 MVP</p>
      <h1>一个可以直接拿去展示的学院级 AI 社区原型。</h1>
      <p class="lead">
        这版聚焦社区氛围、问答交流、公告运营和最小中台闭环。你可以从不同角色登录，顺着前台到后台把整个故事讲完整。
      </p>
      <div class="metric-row">
        <div>
          <strong>3</strong>
          <span>核心板块</span>
        </div>
        <div>
          <strong>3</strong>
          <span>Go 服务</span>
        </div>
        <div>
          <strong>1</strong>
          <span>最小中台</span>
        </div>
      </div>
    </section>

    <section class="login-panel panel">
      <div class="section-copy">
        <p class="eyebrow">演示登录</p>
        <h2>选择这次展示要切入的身份</h2>
        <p class="login-hint">
          点击「进入演示」后，输入该角色对应的用户名与密码。
          <RouterLink to="/register" class="inline-link">邀请码注册</RouterLink>
        </p>
      </div>
      <div class="role-grid">
        <article
          v-for="item in roles"
          :key="item.role"
          class="role-card"
          :class="{ 'role-card--open': activeRole === item.role }"
        >
          <h3>{{ item.title }}</h3>
          <p>{{ item.intro }}</p>

          <form
            v-if="activeRole === item.role"
            class="role-login-form"
            @submit.prevent="submit(item.role)"
          >
            <label class="role-login-field">
              <span>用户名</span>
              <input
                v-model="username"
                type="text"
                name="username"
                autocomplete="username"
                :placeholder="item.usernamePlaceholder"
                :disabled="loading"
              />
            </label>
            <label class="role-login-field">
              <span>密码</span>
              <input
                v-model="password"
                type="password"
                name="password"
                autocomplete="current-password"
                placeholder="请输入密码"
                :disabled="loading"
              />
            </label>
            <p v-if="error" class="role-login-error">{{ error }}</p>
            <div class="role-login-actions">
              <button type="button" class="ghost-button" :disabled="loading" @click="cancelLogin">
                取消
              </button>
              <button type="submit" class="primary-button" :disabled="loading">
                {{ loading ? '登录中…' : '确认登录' }}
              </button>
            </div>
          </form>

          <button v-else class="primary-button" type="button" @click="openLogin(item.role)">
            进入演示
          </button>
        </article>
      </div>
    </section>
  </main>
</template>
