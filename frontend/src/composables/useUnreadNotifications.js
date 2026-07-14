import { ref } from 'vue'
import { forumApi } from '../api'

const unreadCount = ref(0)
const loading = ref(false)

export function useUnreadNotifications() {
  async function refreshUnreadCount() {
    const token = localStorage.getItem('ai-forum-token')
    if (!token) {
      unreadCount.value = 0
      return 0
    }
    loading.value = true
    try {
      unreadCount.value = await forumApi.getUnreadNotificationCount()
      return unreadCount.value
    } catch {
      unreadCount.value = 0
      return 0
    } finally {
      loading.value = false
    }
  }

  function setUnreadCount(value) {
    unreadCount.value = Math.max(0, Number(value) || 0)
  }

  return {
    unreadCount,
    unreadLoading: loading,
    refreshUnreadCount,
    setUnreadCount,
  }
}
