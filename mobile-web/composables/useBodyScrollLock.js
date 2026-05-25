import { onBeforeUnmount, watch } from 'vue'

const BODY_CLASS = 'mw-drawer-open'

/**
 * Lock document scroll while mobile drawer is open (iOS Safari friendly).
 * @param {import('vue').Ref<boolean>} isOpen
 */
export function useBodyScrollLock(isOpen) {
  const apply = (open) => {
    document.body.classList.toggle(BODY_CLASS, Boolean(open))
  }

  watch(isOpen, apply, { immediate: true })

  onBeforeUnmount(() => {
    document.body.classList.remove(BODY_CLASS)
  })
}
