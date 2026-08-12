import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { setRefreshToken, setTokenCookie } from '../api/http'
import { userApi } from '../api'

export const useSessionStore = defineStore('session', () => {
  const token = ref(localStorage.getItem('ai-forum-token') || '')
  const currentUser = ref(JSON.parse(localStorage.getItem('ai-forum-user') || 'null'))
  const flashMessage = ref('')
  const flashType = ref('info')
  let validatedToken = ''
  let validationPromise = null

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
    setTokenCookie(result.token)
    if (result.refreshToken) setRefreshToken(result.refreshToken)
    validatedToken = result.token
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

  function clearSession() {
    token.value = ''
    currentUser.value = null
    validatedToken = ''
    localStorage.removeItem('ai-forum-token')
    localStorage.removeItem('ai-forum-user')
    setRefreshToken('')
    setTokenCookie('')
  }

  async function ensureValidSession() {
    if (!token.value) {
      clearSession()
      return false
    }

    setTokenCookie(token.value)
    if (validatedToken === token.value && currentUser.value) return true
    if (validationPromise) return validationPromise

    validationPromise = (async () => {
      try {
        await refreshMe()
        validatedToken = token.value
        setTokenCookie(token.value)
        return true
      } catch {
        clearSession()
        return false
      } finally {
        validationPromise = null
      }
    })()

    return validationPromise
  }

  function logout() {
    userApi.logout().catch(() => {})
    clearSession()
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
    ensureValidSession,
    clearSession,
    logout,
  }
})
