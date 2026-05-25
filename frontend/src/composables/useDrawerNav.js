import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useBodyScrollLock } from '../../../mobile-web/composables/useBodyScrollLock.js'

const DESKTOP_BREAKPOINT = 1024

export function useDrawerNav() {
  const route = useRoute()
  const drawerOpen = ref(false)

  useBodyScrollLock(drawerOpen)

  function closeDrawer() {
    drawerOpen.value = false
  }

  function openDrawer() {
    drawerOpen.value = true
  }

  function toggleDrawer() {
    drawerOpen.value = !drawerOpen.value
  }

  function onResize() {
    if (window.innerWidth >= DESKTOP_BREAKPOINT) {
      closeDrawer()
    }
  }

  watch(
    () => route.fullPath,
    () => {
      closeDrawer()
    },
  )

  onMounted(() => {
    window.addEventListener('resize', onResize)
    window.__closeDrawer = closeDrawer
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', onResize)
    if (window.__closeDrawer === closeDrawer) {
      delete window.__closeDrawer
    }
  })

  return {
    drawerOpen,
    openDrawer,
    closeDrawer,
    toggleDrawer,
  }
}
