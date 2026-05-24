import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAdminStore = defineStore('admin', () => {
  const adminUser = ref(null)
  const pendingAuditCount = ref(0)

  function setAdminUser(user) {
    adminUser.value = user
  }

  function clearAdminUser() {
    adminUser.value = null
  }

  function updatePendingCount(count) {
    pendingAuditCount.value = count
  }

  return {
    adminUser,
    pendingAuditCount,
    setAdminUser,
    clearAdminUser,
    updatePendingCount,
  }
})
