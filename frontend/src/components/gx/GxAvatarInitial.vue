<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { maskStudentLabel } from '../../utils/displayName'

const props = defineProps({
  name: { type: String, default: '' },
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
  size: { type: Number, default: 40 },
})

const imageFailed = ref(false)
const retryCount = ref(0)
const label = computed(() => maskStudentLabel(props.name))
let retryTimer = 0

function retryableSrc(src) {
  if (!retryCount.value || src.startsWith('blob:') || src.startsWith('data:')) return src
  const separator = src.includes('?') ? '&' : '?'
  return `${src}${separator}avatar_retry=${retryCount.value}`
}

const imageSrc = computed(() => (imageFailed.value || !props.src ? '' : retryableSrc(props.src)))
const initial = computed(() => {
  const t = label.value.replace(/[·*]/g, '').trim()
  return t ? t.charAt(0).toUpperCase() : '?'
})

function clearRetryTimer() {
  if (!retryTimer) return
  window.clearTimeout(retryTimer)
  retryTimer = 0
}

function resetImageState() {
  clearRetryTimer()
  imageFailed.value = false
  retryCount.value = 0
}

function handleImageError() {
  imageFailed.value = true
  if (!props.src || retryCount.value >= 2) return
  clearRetryTimer()
  retryTimer = window.setTimeout(() => {
    retryCount.value += 1
    imageFailed.value = false
    retryTimer = 0
  }, 1200)
}

watch(
  () => props.src,
  resetImageState,
)

onBeforeUnmount(clearRetryTimer)
</script>

<template>
  <span
    class="gx-avatar-initial inline-flex shrink-0 items-center justify-center rounded-full bg-gx-primary/10 font-semibold text-gx-primary"
    :style="{ width: `${size}px`, height: `${size}px`, fontSize: `${Math.round(size * 0.4)}px` }"
    :title="label"
    aria-hidden="true"
  >
    <img
      v-if="imageSrc"
      class="gx-avatar-initial__image"
      :src="imageSrc"
      :alt="alt || label"
      loading="lazy"
      @error="handleImageError"
    />
    <span v-else>{{ initial }}</span>
  </span>
</template>
