<script setup>
defineProps({
  title: { type: String, default: '暂无内容' },
  description: { type: String, default: '' },
})
</script>

<template>
  <div class="gx-empty-state py-12 text-center">
    <p class="text-title text-gx-primary">{{ title }}</p>
    <p v-if="description" class="mt-2 text-body text-gx-muted">{{ description }}</p>
    <div v-if="$slots.default" class="mt-4 flex justify-center gap-2">
      <slot />
    </div>
  </div>
</template>
