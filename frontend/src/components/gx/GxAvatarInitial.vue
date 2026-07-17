<script setup>
import { computed, ref, watch } from 'vue'
import { maskStudentLabel } from '../../utils/displayName'

const props = defineProps({
  name: { type: String, default: '' },
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
  size: { type: Number, default: 40 },
})

const imageFailed = ref(false)
const label = computed(() => maskStudentLabel(props.name))
const imageSrc = computed(() => (imageFailed.value ? '' : props.src))
const initial = computed(() => {
  const t = label.value.replace(/[·*]/g, '').trim()
  return t ? t.charAt(0).toUpperCase() : '?'
})

watch(
  () => props.src,
  () => {
    imageFailed.value = false
  },
)
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
      @error="imageFailed = true"
    />
    <span v-else>{{ initial }}</span>
  </span>
</template>
