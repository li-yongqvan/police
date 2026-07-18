import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { setRefreshToken } from '../api/http'
import { userApi } from '../api'

export const useSessionStore = defineStore('session', () => {
  const token = ref(localStorage.getItem('ai-forum-token') || '')
  const currentUser = ref(JSON.parse(localStorage.getItem('ai-forum-user') || 'null'))
  const flashMessage = ref('')
  const flashType = ref('info')

  const canAccessAdmin = computed(() =>
    ['admin', 'platform_admin'].includes(currentUser.value?.role),
  )

  function setFlash(message, type = 'info') {
    flashMessage.value = message
    flashType.value = type
    if (message) {
      window.clearTimeout(setFlash.timer)
      setFlash.timer = window.setTimeout(() => {
        flashMessage.value = ''
      }, 3200)
    }
  }

  function persistSession(result) {
    token.value = result.token
    currentUser.value = result.user
    localStorage.setItem('ai-forum-token', result.token)
    localStorage.setItem('ai-forum-user', JSON.stringify(result.user))
    if (result.refreshToken) setRefreshToken(result.refreshToken)
  }

  async function loginWithCredentials(username, password) {
    const result = await userApi.login(username, password)
    persistSession(result)
    return result.user
  }

  function routeAfterLogin(user) {
    if (['admin', 'platform_admin'].includes(user.role)) return '/admin'
    return '/community'
  }

  async function refreshMe() {
    if (!token.value) return
    currentUser.value = await userApi.me()
    localStorage.setItem('ai-forum-user', JSON.stringify(currentUser.value))
  }

  function logout() {
    token.value = ''
    currentUser.value = null
    localStorage.removeItem('ai-forum-token')
    localStorage.removeItem('ai-forum-user')
    setRefreshToken('')
  }

  return {
    token,
    currentUser,
    flashMessage,
    flashType,
    canAccessAdmin,
    setFlash,
    persistSession,
    loginWithCredentials,
    routeAfterLogin,
    refreshMe,
    logout,
  }
})
