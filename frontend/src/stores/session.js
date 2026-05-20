import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
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

  async function loginAs(role) {
    const result = await userApi.demoLogin(role)
    token.value = result.token
    currentUser.value = result.user
    localStorage.setItem('ai-forum-token', result.token)
    localStorage.setItem('ai-forum-user', JSON.stringify(result.user))
  }

  async function refreshMe() {
    if (!token.value) return
    currentUser.value = await userApi.me(token.value)
    localStorage.setItem('ai-forum-user', JSON.stringify(currentUser.value))
  }

  function logout() {
    token.value = ''
    currentUser.value = null
    localStorage.removeItem('ai-forum-token')
    localStorage.removeItem('ai-forum-user')
  }

  return {
    token,
    currentUser,
    flashMessage,
    flashType,
    canAccessAdmin,
    setFlash,
    loginAs,
    refreshMe,
    logout,
  }
})
